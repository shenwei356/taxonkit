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
	"strconv"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// lineageCmd represents the fx2tab command
var lineageCmd = &cobra.Command{
	Use:   "lineage",
	Short: "query lineage of given taxids",
	Long: `query lineage of given taxids

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")
		field := getFlagPositiveInt(cmd, "taxid-field") - 1

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string

		log.Infof("parsing names file: %s", config.NamesFile)
		names = getTaxonNames(config.NamesFile, config.Threads, 10)
		log.Infof("%d names parsed", len(names))

		log.Infof("parsing nodes file: %s", config.NodesFile)

		reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
		checkError(err)

		tree := make(map[int32]int32)
		ranks := make(map[int32]string)
		var info taxonInfo
		var child, parent int32
		var n int64
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data := range chunk.Data {
				info = data.(taxonInfo)
				child, parent = info.child, info.parent

				tree[child] = parent
				ranks[child] = info.rank
				n++
			}
		}

		log.Infof("%d nodes parsed", n)

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
			var child, parent int32
			var ok bool
			child = int32(id)
			for true {
				parent, ok = tree[child]
				if !ok {
					break
				}
				lineage = append(lineage, names[child])
				if printLineageInTaxid {
					lineageInTaxid = append(lineageInTaxid, strconv.Itoa(int(child)))
				}
				if parent == 1 && child != 1 {
					break
				}
				child = parent
			}
			child = int32(id)
			if printLineageInTaxid {
				return taxid2lineage{line, child,
					strings.Join(stringutil.ReverseStringSlice(lineage), ";"),
					strings.Join(stringutil.ReverseStringSlice(lineageInTaxid), ";"),
				}, true, nil
			}
			return taxid2lineage{line, child,
				strings.Join(stringutil.ReverseStringSlice(lineage), ";"),
				"",
			}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var t2l taxid2lineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					t2l = data.(taxid2lineage)
					if printLineageInTaxid {
						outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\n", t2l.line, t2l.lineage, t2l.lineageInTaxid))
					} else {
						outfh.WriteString(fmt.Sprintf("%s\t%s\n", t2l.line, t2l.lineage))
					}
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(lineageCmd)
	lineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `show lineage consisting of taxids`)
	lineageCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. data should be tab-separated")
}
