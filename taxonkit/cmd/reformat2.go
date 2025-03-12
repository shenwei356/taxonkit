// Copyright © 2016-2022 Wei Shen <shenwei356@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// reformat2Cmd represents the fx2tab command
var reformat2Cmd = &cobra.Command{
	Use:   "reformat2",
	Short: "Reformat lineage in chosen ranks, allowing more ranks than 'reformat'",
	Long: `Reformat lineage in chosen ranks, allowing more ranks than 'reformat'

Input:

  - List of TaxIds, one record per line.
  - Or tab-delimited format.
    Please specify the TaxId field with flag -I/--taxid-field (default 1)
  - Supporting (gzipped) file or STDIN.

Output:

  1. Input line data.
  2. Reformated lineage.
  3. (Optional) TaxIds taxons in the lineage (-t/--show-lineage-taxids)

Output format:

  1. it can contain some escape characters like "\t".
  2. For subspecies nodes, the rank might be "subpecies", "strain", or "no rank".
     You can use "|" to set multiple ranks, and the first valid one will be outputted.
     For example,
       $ echo -ne "562\n83333\n2697049\n" \
          | taxonkit lineage -L -r \
          | taxonkit reformat2 -f "{species};{strain|subspecies|no rank}"
       562     species Escherichia coli;
       83333   strain  Escherichia coli;Escherichia coli K-12
       2697049 no rank Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2

Differences from 'taxonkit reformat':

  - [input] only accept TaxIDs
  - [format] accept more rank place holders, not just the seven canonical ones.
  - [format] use the full name of ranks, such as "{species}", rather than "{s}"
  - [format] support multiple ranks in one place holder, such as "{subspecies|strain}"
  - do not automatically add prefixes, but you can set in the format

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		format := getFlagString(cmd, "format")
		blank := getFlagString(cmd, "miss-rank-repl")
		iblank := getFlagString(cmd, "miss-taxid-repl")
		taxIdField := getFlagPositiveInt(cmd, "taxid-field")
		noRanks := getFlagStringSlice(cmd, "no-ranks")
		trim := getFlagBool(cmd, "trim")

		if config.Verbose {
			log.Infof("parsing TaxIds from field %d", taxIdField)
		}
		taxIdField--

		noRanksMap := make(map[string]interface{}, 8)
		if len(noRanks) > 0 {
			for _, rank := range noRanks {
				noRanksMap[strings.ToLower(rank)] = struct{}{}
			}
		}

		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")

		// check format
		if !reRankPlaceHolder2.MatchString(format) {
			checkError(fmt.Errorf("placeholder of simplified rank not found in output format: %s", format))
		}
		matches := reRankPlaceHolder2.FindAllStringSubmatch(format, -1)
		if len(matches) == 0 {
			checkError(fmt.Errorf("no placeholder given %s", format))
		}

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		// --------------------------------------------------------
		// load data

		var tree0 map[uint32]uint32
		var ranks0 map[uint32]string
		var names0 map[uint32]string
		var delnodes0 map[uint32]struct{}
		var merged0 map[uint32]uint32

		tree0, ranks0, names0, delnodes0, merged0 = loadData(config, true, true)

		// --------------------------------------------------------

		type line2flineage struct {
			line      string
			flineage  string
			iflineage string
		}

		unescape := stringutil.UnEscaper()

		blankS := format
		iblankS := format
		for _, re := range reRankPlaceHolders {
			blankS = re.ReplaceAllString(blankS, blank)
		}
		for _, re := range reRankPlaceHolders {
			iblankS = re.ReplaceAllString(iblankS, iblank)
		}

		fn := func(line string) (interface{}, bool, error) {
			if len(line) == 0 || line[0] == '#' {
				return nil, false, nil
			}
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}
			data := strings.Split(line, "\t")

			if len(data) < taxIdField+1 {
				return nil, false, fmt.Errorf("taxid-field (%d) out of range (%d):%s", taxIdField+1, len(data), line)
			}

			// -----------------------------------------------

			var ok bool

			var taxid uint32
			var taxidInt int

			var names []string
			var ranks []string
			var taxids []uint32

			// -----------------------------------------------
			// get the taxid

			taxidInt, err = strconv.Atoi(data[taxIdField])
			if err != nil || taxidInt < 0 {
				// checkError(fmt.Errorf("invalid TaxId: %s", data[taxIdField]))
				log.Warningf("invalid TaxId: %s", data[taxIdField])
				return line2flineage{line, "", ""}, true, nil
			}
			taxid = uint32(taxidInt)

			// -----------------------------------------------
			// query complete lineage with the taxid

			names, ranks, taxids, ok = queryNamesRanksTaxids(tree0, ranks0, names0, delnodes0, merged0, taxid)
			if !ok { // taxid not found
				// return line2flineage{line, "", ""}, true, nil
				return line2flineage{line, unescape(blankS), unescape(iblankS)}, true, nil
			}

			rank2idx := poolRank2idx.Get().(*map[string]int)
			clear(*rank2idx)
			var meetKnownRanks bool
			var lastKnownRank string
			for i, rank := range ranks {
				rank = strings.ToLower(rank)

				if _, ok = noRanksMap[rank]; ok {
					if meetKnownRanks {
						(*rank2idx)[rank] = i
					}
				} else {
					meetKnownRanks = true
					(*rank2idx)[rank] = i
					lastKnownRank = rank
				}
			}

			flineage := format
			var iflineage string
			if printLineageInTaxid {
				iflineage = format
			}

			var i int
			var _matches []string
			var _match string
			var matched bool
			var foundLastKnownRank bool
			var repl, irepl string
			for _, match := range matches {
				_matches = strings.Split(match[1], "|")

				matched = false
				for _, _match = range _matches {
					if _match == "" {
						continue
					}

					_match = strings.ToLower(_match)
					if i, ok = (*rank2idx)[_match]; !ok {
						continue
					}

					flineage = strings.ReplaceAll(flineage, match[0], names[i])
					if printLineageInTaxid {
						iflineage = strings.ReplaceAll(iflineage, match[0], strconv.Itoa(int(taxids[i])))
					}
					matched = true

					if _match == lastKnownRank {
						foundLastKnownRank = true
					}
					break
				}

				if !matched {
					if !foundLastKnownRank {
						repl, irepl = blank, iblank
					} else if trim {
						repl, irepl = "", ""
					} else {
						repl, irepl = blank, iblank
					}

					flineage = strings.ReplaceAll(flineage, match[0], repl)
					if printLineageInTaxid {
						iflineage = strings.ReplaceAll(iflineage, match[0], irepl)
					}
				}
			}

			// recycle
			ranks = ranks[:0]
			poolStringsN16.Put(ranks)

			names = names[:0]
			poolStringsN16.Put(names)
			taxids = taxids[:0]
			poolUint32N16.Put(taxids)

			poolRank2idx.Put(rank2idx)

			return line2flineage{line, unescape(flineage), unescape(iflineage)}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 64, fn)
			checkError(err)

			var l2s line2flineage
			var data interface{}
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data = range chunk.Data {
					l2s = data.(line2flineage)

					if printLineageInTaxid {
						outfh.WriteString(l2s.line + "\t" + l2s.flineage + "\t" + l2s.iflineage + "\n")
					} else {
						outfh.WriteString(l2s.line + "\t" + l2s.flineage + "\n")
					}
					if config.LineBuffered {
						outfh.Flush()
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(reformat2Cmd)

	reformat2Cmd.Flags().StringP("format", "f", "{superkingdom};{phylum};{class};{order};{family};{genus};{species}", "output format, placeholders of rank are needed")
	reformat2Cmd.Flags().StringP("miss-rank-repl", "r", "", `replacement string for missing rank`)
	reformat2Cmd.Flags().StringP("miss-taxid-repl", "R", "", `replacement string for missing taxid`)
	reformat2Cmd.Flags().BoolP("trim", "T", false, "do not replace missing ranks lower than the rank of the current node")

	reformat2Cmd.Flags().IntP("taxid-field", "I", 1, "field index of taxid. input data should be tab-separated. it overrides -i/--lineage-field")
	reformat2Cmd.Flags().BoolP("show-lineage-taxids", "t", false, `show corresponding taxids of reformated lineage`)

	reformat2Cmd.Flags().StringSliceP("no-ranks", "B", []string{"no rank", "clade"}, `rank names of no-rank. A lineage might have many "no rank" ranks, we only keep the last one below known ranks`)

}

var reRankPlaceHolder2 = regexp.MustCompile(`\{([^\{\}]+?)\}`)

var poolRank2idx = &sync.Pool{New: func() interface{} {
	tmp := make(map[string]int, 64)
	return &tmp
}}
