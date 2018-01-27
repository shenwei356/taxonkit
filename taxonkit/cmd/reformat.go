// Copyright © 2016 Wei Shen <shenwei356@gmail.com>
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
	"runtime"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// flineageCmd represents the fx2tab command
var flineageCmd = &cobra.Command{
	Use:   "reformat",
	Short: "reformat lineage",
	Long: `reformat lineage

Output format can be formated by flag --format, available placeholders:

    {k}: superkingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {S}: subspecies

Output format can contains some escape charactors like "\t".

This command appends reformated lineage to the input line, along with an extra
flag-column to indicate the reliability of the result.

Note that lots of taxids share same taxon name like "diastema" and "solieria".
This command does not consider the context for a given taxon name, which may
bring potential bug. If the reformating is reliable, the value of flag column
is "OK", otherwise it list the taxon names which may cause error.

The corresponding taxids of reformated lineage can be provided as another
column by flag "-t/--show-lineage-taxids".

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "miss-rank-repl")
		iblank := getFlagString(cmd, "miss-taxid-repl")
		fill := getFlagBool(cmd, "fill-miss-rank")
		field := getFlagPositiveInt(cmd, "lineage-field") - 1
		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")

		// check format
		if !reRankPlaceHolder.MatchString(format) {
			checkError(fmt.Errorf("placeholder of simplified rank not found in output format: %s", format))
		}
		matches := reRankPlaceHolder.FindAllStringSubmatch(format, -1)
		outSranks := make(map[string]struct{})
		outSranksList := []string{}
		var currentWeight float32
		var currentSymbol string
		for _, match := range matches {
			if weight, ok := symbol2weight[match[1]]; !ok {
				checkError(fmt.Errorf("invalid placeholder: %s", match[0]))
			} else {
				if weight < currentWeight {
					checkError(fmt.Errorf(`invalid placeholder order: {%s} {%s}. "%s" should be behind of "%s"`,
						currentSymbol, match[1], symbol2rank[currentSymbol], symbol2rank[match[1]]))
				}
				outSranks[match[1]] = struct{}{}
				outSranksList = append(outSranksList, match[1])
				currentWeight = weight
				currentSymbol = match[1]
			}
		}

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		if config.Verbose {
			log.Infof("parsing names file: %s", config.NamesFile)
		}

		taxid2name, name2taxids, fuzzyNamesMap := parseTaxonNames(config.NamesFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d names parsed", len(taxid2name))

			log.Infof("parsing nodes file: %s", config.NodesFile)
		}

		reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
		checkError(err)

		name2rank := make(map[string]string)
		var info taxonInfo
		var n int64
		var data interface{}
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data = range chunk.Data {
				info = data.(taxonInfo)
				name2rank[strings.ToLower(taxid2name[info.child])] = info.rank
				n++
			}
		}

		if config.Verbose {
			log.Infof("%d nodes parsed", n)
		}

		type line2flineage struct {
			line       string
			flineage   string
			iflineage  string
			fuzzyNames []string
		}

		unescape := stringutil.UnEscaper()
		fn := func(line string) (interface{}, bool, error) {
			if len(line) == 0 || line[0] == '#' {
				return nil, false, nil
			}
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}
			data := strings.Split(line, "\t")
			if len(data) < field+1 {
				return nil, false, fmt.Errorf("lineage-field (%d) out of range (%d)", field+1, len(data))
			}

			// names and weights
			names2 := strings.Split(data[field], delimiter)
			weights := make([]float32, len(names2))
			var rank, srank string
			var ok bool
			srank2name := make(map[string]string)
			var currentWeight float32
			var fuzzyNames []string
			name2Name := make(map[string]string, len(names2)) // lower case of name : name
			var lname string
			for i, name := range names2 {
				if name == "" {
					continue
				}

				if _, ok = fuzzyNamesMap[name]; ok {
					fuzzyNames = append(fuzzyNames, name)
				}

				lname = strings.ToLower(name)
				name2Name[lname] = name
				name = lname

				rank, ok = name2rank[name]
				if !ok { // unofficial name
					currentWeight += 0.1
					weights[i] = currentWeight
					continue
				}
				if rank != norank {
					if srank, ok = rank2symbol[rank]; ok {
						srank2name[srank] = name
						weights[i] = symbol2weight[srank]
						currentWeight = weights[i]
					} else {
						// log.Warningf("please contact author to add this rank to code: %s", rank)
					}
				} else {
					currentWeight += 0.1
					weights[i] = currentWeight
				}
			}
			// preprare replacements.
			// find the orphan names and missing ranks
			replacements := make(map[string]string, len(matches))
			ireplacements := make(map[string]string, len(matches))
			if !fill {
				for _, match := range matches {
					replacements[match[1]] = blank
					ireplacements[match[1]] = blank
				}
			}

			orphans := make(map[string]float32)
			orphansList := []string{}
			existedSranks := make(map[string]struct{})
			for i, name := range names2 {
				if name == "" {
					continue
				}
				name = strings.ToLower(name)

				if name2rank[name] == norank {
					orphans[name] = weights[i]
					orphansList = append(orphansList, name)
				} else {
					if _, ok = outSranks[rank2symbol[name2rank[name]]]; ok { // to be outputted
						replacements[rank2symbol[name2rank[name]]] = name2Name[name]
						ireplacements[rank2symbol[name2rank[name]]] = fmt.Sprintf("%d", name2taxids[name][0])
						existedSranks[rank2symbol[name2rank[name]]] = struct{}{}
					} else if name2rank[name] == "" {
						orphans[name] = weights[i]
						orphansList = append(orphansList, name)
					}
				}
			}

			if fill {
				jj := -1
				var hit bool
				var lastRank string
				for i, srank := range outSranksList {
					if _, ok = existedSranks[srank]; ok {
						lastRank = replacements[srank]
						continue
					}
					hit = false
					for j, name := range orphansList {
						if j <= jj {
							continue
						}
						if i == 0 {
							if orphans[name] < symbol2weight[outSranksList[i]] {
								hit = true
							}
						} else if i == len(outSranksList)-1 {

						} else if orphans[name] > symbol2weight[outSranksList[i-1]] &&
							orphans[name] < symbol2weight[outSranksList[i+1]] {
							hit = true
						}

						if hit {
							replacements[srank] = name2Name[name]
							ireplacements[srank] = fmt.Sprintf("%d", name2taxids[name][0])
							jj = j
							break
						}
					}
					if !hit {
						if blank == "" {
							replacements[srank] = fmt.Sprintf("unclassified %s %s", lastRank, symbol2rank[srank])
						} else {
							replacements[srank] = blank
						}
						ireplacements[srank] = iblank
					}
				}
			}

			flineage := format
			iflineage := format
			for srank, re := range reRankPlaceHolders {
				flineage = re.ReplaceAllString(flineage, replacements[srank])
				iflineage = re.ReplaceAllString(iflineage, ireplacements[srank])
			}

			return line2flineage{line, unescape(flineage), unescape(iflineage), fuzzyNames}, true, nil
		}

		var flag string
		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var l2s line2flineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					l2s = data.(line2flineage)
					if len(l2s.fuzzyNames) > 0 {
						flag = strings.Join(l2s.fuzzyNames, delimiter)
					} else {
						flag = "OK"
					}
					if printLineageInTaxid {
						outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\n", l2s.line, l2s.flineage, flag, l2s.iflineage))
					} else {
						outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\n", l2s.line, l2s.flineage, flag))
					}
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(flineageCmd)

	flineageCmd.Flags().StringP("format", "f", "{k};{p};{c};{o};{f};{g};{s}", "output format, placeholders of rank are needed")
	flineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in input lineage")
	flineageCmd.Flags().StringP("miss-rank-repl", "r", "", `replacement string for missing rank, if given "", "unclassified xxx xxx" will used`)
	flineageCmd.Flags().StringP("miss-taxid-repl", "R", "0", `replacement string for missing taxid`)
	flineageCmd.Flags().BoolP("fill-miss-rank", "F", false, "estimate and fill missing rank with original lineage information (recommended)")
	flineageCmd.Flags().IntP("lineage-field", "i", 2, "field index of lineage. data should be tab-separated")
	flineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `show reformated taxids`)
}
