// Copyright © 2016-2021 Wei Shen <shenwei356@gmail.com>
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
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// lineageCmd represents the lineage command
var lineageCmd = &cobra.Command{
	Use:   "lineage",
	Short: "Query taxonomic lineage of given taxIDs",
	Long: `Query taxonomic lineage of given taxIDs

Input:

  - List of taxIDs, one taxID per line.
  - Or tab-delimited format, please specify taxID field 
    with flag -i/--taxid-field (default 1).
  - Supporting (gzipped) file or STDIN.

Output:

  1. Input line data.
  2. (Optional) Status code (-c/--show-status-code), values:
     - "-1" for queries not found in whole database.
     - "0" for deleted taxIDs, provided by "delnodes.dmp".
     - New taxIDs for merged taxIDs, provided by "merged.dmp".
     - Taxids for these found in "nodes.dmp".
  3. Lineage, delimiter can be changed with flag -d/--delimiter.
  4. (Optional) Lineage in taxIDs (-t/--show-lineage-taxids)
  5. (Optional) Rank (-r/--show-rank)

Filter out invalid and deleted taxids, and replace merged 
taxids with new ones:
    
    # input is one-column-taxid
    $ taxonkit lineage -c taxids.txt \
        | awk '$2>0' \
        | cut -f 2-
        
    # taxids are in 3rd field in a 4-columns tab-delimited file,
    # for $5, where 5 = 4 + 1.
    $ cat input.txt \
        | taxonkit lineage -c -i 3 \
        | csvtk filter2 -H -t -f '$5>0' \
        | csvtk -H -t cut -f -3

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		delimiter := getFlagString(cmd, "delimiter")
		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")
		printLineageInRank := getFlagBool(cmd, "show-lineage-ranks")
		printRank := getFlagBool(cmd, "show-rank")
		printName := getFlagBool(cmd, "show-name")
		field := getFlagPositiveInt(cmd, "taxid-field") - 1
		showCode := getFlagBool(cmd, "show-status-code")
		noLineage := getFlagBool(cmd, "no-lineage")

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		if noLineage && !printRank && !printName {
			checkError(fmt.Errorf("when given -L/--no-lineage, -n/--show-name or/and -r/--show-rank needed"))
		}

		// -------------------- load data ----------------------

		var tree map[uint32]uint32
		var ranks map[uint32]string
		var names map[uint32]string
		var delnodes map[uint32]struct{}
		var merged map[uint32]uint32
		tree, ranks, names, delnodes, merged = loadData(config, true, printRank || printLineageInRank)

		// -------------------- load data ----------------------

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		type taxid2lineage struct {
			line           string
			taxid          uint32
			lineage        string
			lineageInTaxid string
			lineageInRank  string
			notFound       bool
		}

		fn := func(line string) (interface{}, bool, error) {
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}

			data := strings.Split(line, "\t")
			if len(data) <= field {
				field = len(data) - 1
			}

			if data[field] == "" {
				return taxid2lineage{line, 0, "", "", "", false}, true, nil
			}
			id, e := strconv.Atoi(data[field])
			if e != nil {
				return taxid2lineage{line, 0, "", "", "", false}, true, nil
			}

			lineage := make([]string, 0, 16)
			lineageInTaxid := make([]string, 0, 16)
			lineageInRank := make([]string, 0, 16)
			var child, parent, newtaxid uint32
			var ok bool
			child = uint32(id)
			var notFound bool
			for true {
				parent, ok = tree[child]
				if !ok { // taxid not found
					// check if it was deleted
					if _, ok = delnodes[child]; ok {
						// log
						log.Warningf("taxid %d was deleted", child)
						id = 0
						break
					}
					// check if it was merged
					if newtaxid, ok = merged[child]; ok {
						// log
						log.Warningf("taxid %d was merged into %d", child, newtaxid)
						child = newtaxid
						parent = tree[child]
						id = int(child)
					} else {
						id = 0
						log.Warningf("taxid %d not found", child)
						notFound = true
						break
					}
				}

				lineage = append(lineage, names[child])
				if noLineage {
					break
				}

				if printLineageInTaxid {
					lineageInTaxid = append(lineageInTaxid, strconv.Itoa(int(child)))
				}
				if printLineageInRank {
					lineageInRank = append(lineageInRank, ranks[child])
				}

				if parent == 1 {
					break
				}
				child = parent
			}
			child = uint32(id)

			var lineageInTaxidS, lineageInRankS string
			if printLineageInTaxid {
				lineageInTaxidS = strings.Join(stringutil.ReverseStringSlice(lineageInTaxid), delimiter)
			}

			if printLineageInRank {
				lineageInRankS = strings.Join(stringutil.ReverseStringSlice(lineageInRank), delimiter)
			}

			return taxid2lineage{line, child,
				strings.Join(stringutil.ReverseStringSlice(lineage), delimiter),
				lineageInTaxidS,
				lineageInRankS,
				notFound,
			}, true, nil
		}

		var buf bytes.Buffer
		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var t2l taxid2lineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					t2l = data.(taxid2lineage)

					buf.Reset()
					buf.WriteString(t2l.line)

					if showCode {
						if t2l.notFound {
							buf.WriteString("\t-1")
						} else {
							buf.WriteString("\t" + strconv.Itoa(int(t2l.taxid)))
						}
					}
					if !noLineage {
						buf.WriteString("\t" + t2l.lineage)
					}

					if printLineageInTaxid && !noLineage {
						buf.WriteString("\t" + t2l.lineageInTaxid)
					}

					if printName {
						buf.WriteString("\t" + names[t2l.taxid])
					}
					if printRank {
						buf.WriteString("\t" + ranks[t2l.taxid])
					}

					if printLineageInRank && !noLineage {
						buf.WriteString("\t" + t2l.lineageInRank)
					}

					buf.WriteString("\n")

					outfh.WriteString(buf.String())
					if config.LineBuffered {
						outfh.Flush()
					}
				}
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(lineageCmd)
	lineageCmd.Flags().BoolP("show-status-code", "c", false, "show status code before lineage")
	lineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `appending lineage consisting of taxids`)
	lineageCmd.Flags().BoolP("show-lineage-ranks", "R", false, `appending ranks of all levels`)
	lineageCmd.Flags().BoolP("show-rank", "r", false, `appending rank of taxids`)
	lineageCmd.Flags().BoolP("show-name", "n", false, `appending scientific name`)
	lineageCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. input data should be tab-separated")
	lineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in lineage")
	lineageCmd.Flags().BoolP("no-lineage", "L", false, "do not show lineage, when user just want names or/and ranks")
}
