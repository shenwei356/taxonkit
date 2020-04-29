// Copyright © 2016-2020 Wei Shen <shenwei356@gmail.com>
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
	"runtime"
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
	Short: "query lineage of given taxids",
	Long: `query lineage of given taxids

Input:
  - List of taxids, one taxid per line.
  - Or tab-delimited format, please specify taxid field with flag -i/--taxid-field.
  - Supporting (gzipped) file or STDIN.

Output:
  0. Input line.
  1. Status code (optional with flag -c/--show-status-code)
     - "-1" for queries not found in whole database.
     - "0" for deleted taxids, provided by "delnodes.dmp".
     - New taxids for merged taxids, provided by "merged.dmp".
     - Taxids for these found in "nodes.dmp".
  2. Lineage, delimiter can be changed with flag -d/--delimiter.
  3. Lineage taxids (optional with flag -t/--show-lineage-taxids)
  4. Rank (optional with flag -r/--show-rank)

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		delimiter := getFlagString(cmd, "delimiter")
		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")
		printRank := getFlagBool(cmd, "show-rank")
		field := getFlagPositiveInt(cmd, "taxid-field") - 1
		showCode := getFlagBool(cmd, "show-status-code")

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string
		var delnodes map[int32]struct{}
		var merged map[int32]int32

		if config.Verbose {
			log.Infof("parsing delnodes file: %s", config.NamesFile)
		}

		delnodes = getDelnodesMap(config.DelNodesFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d delnodes parsed", len(delnodes))

			log.Infof("parsing merged file: %s", config.NamesFile)
		}

		merged = getMergedNodesMap(config.MergedFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d merged nodes parsed", len(merged))

			log.Infof("parsing names file: %s", config.NamesFile)
		}

		names = getTaxonNames(config.NamesFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d names parsed", len(names))

			log.Infof("parsing nodes file: %s", config.NodesFile)
		}
		reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
		checkError(err)

		tree := make(map[int32]int32)
		ranks := make(map[int32]string)
		var taxon Taxon
		var child, parent int32
		var n int64
		var data interface{}
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data = range chunk.Data {
				taxon = data.(Taxon)
				child, parent = taxon.Taxid, taxon.Parent

				tree[child] = parent
				ranks[child] = taxon.Rank
				n++
			}
		}

		if config.Verbose {
			log.Infof("%d nodes parsed", n)
		}

		type taxid2lineage struct {
			line           string
			taxid          int32
			lineage        string
			lineageInTaxid string
		}

		fn := func(line string) (interface{}, bool, error) {
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}

			data := strings.Split(line, "\t")
			if len(data) < field+1 {
				field = len(data) - 1
			}

			if data[field] == "" {
				return taxid2lineage{line, -1, "", ""}, true, nil
			}
			id, e := strconv.Atoi(data[field])
			if e != nil {
				return taxid2lineage{line, -1, "", ""}, true, nil
			}

			lineage := []string{}
			lineageInTaxid := []string{}
			var child, parent, newtaxid int32
			var ok bool
			child = int32(id)
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
						id = -1
						log.Warningf("taxid %d not found", child)
						break
					}
				}
				lineage = append(lineage, names[child])
				if printLineageInTaxid {
					lineageInTaxid = append(lineageInTaxid, strconv.Itoa(int(child)))
				}
				if parent == 1 {
					break
				}
				child = parent
			}
			child = int32(id)
			if printLineageInTaxid {
				return taxid2lineage{line, child,
					strings.Join(stringutil.ReverseStringSlice(lineage), delimiter),
					strings.Join(stringutil.ReverseStringSlice(lineageInTaxid), delimiter),
				}, true, nil
			}
			return taxid2lineage{line, child,
				strings.Join(stringutil.ReverseStringSlice(lineage), delimiter),
				"",
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
						buf.WriteString(fmt.Sprintf("\t%d", t2l.taxid))
					}
					buf.WriteString("\t" + t2l.lineage)

					if printLineageInTaxid {
						buf.WriteString("\t" + t2l.lineageInTaxid)
					}
					if printRank {
						buf.WriteString("\t" + ranks[t2l.taxid])
					}

					buf.WriteString("\n")

					outfh.WriteString(buf.String())
					if config.LineBuffered {
						outfh.Flush()
					}
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(lineageCmd)
	lineageCmd.Flags().BoolP("show-status-code", "c", false, "show status code between lineage")
	lineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `appending lineage consisting of taxids`)
	lineageCmd.Flags().BoolP("show-rank", "r", false, `appending rank of taxids`)
	lineageCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. data should be tab-separated")
	lineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in lineage")
}
