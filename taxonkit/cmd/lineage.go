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
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// lineageCmd represents the fx2tab command
var lineageCmd = &cobra.Command{
	Use:   "lineage",
	Short: "query lineage of given taxids from file",
	Long: `query lineage of given taxids from file

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		nodesFile := getFlagString(cmd, "nodes")
		namesFile := getFlagString(cmd, "names")
		formatedRank := getFlagBool(cmd, "formated-rank")

		files := getFileList(args)

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string

		log.Infof("parsing names file: %s", namesFile)
		names = getTaxonNames(namesFile, config.Threads, 10)
		log.Infof("%d names parsed", len(names))

		log.Infof("parsing nodes file: %s", nodesFile)

		type info struct {
			child, parent int32
			rank          string
		}
		fn := func(line string) (interface{}, bool, error) {
			items := strings.SplitN(line, "\t", 6)
			if len(items) < 6 {
				return nil, false, nil
			}
			child, e := strconv.Atoi(items[0])
			if e != nil {
				return nil, false, e
			}
			parent, e := strconv.Atoi(items[2])
			if e != nil {
				return nil, false, e
			}
			return info{int32(child), int32(parent), items[4]}, true, nil
		}

		reader, err := breader.NewBufferedReader(nodesFile, config.Threads, 10, fn)
		checkError(err)

		tree := make(map[int32]int32)
		ranks := make(map[int32]string)
		var rel info
		var child, parent int32
		var n int64
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data := range chunk.Data {
				rel = data.(info)
				child, parent = rel.child, rel.parent

				tree[child] = parent
				ranks[child] = rel.rank
				n++
			}
		}

		log.Infof("%d nodes parsed", n)

		type taxid2lineage struct {
			taxid   int32
			lineage string
		}
		rankPrefix := map[string]string{
			"superkingdom": "k__",
			"phylum":       "p__",
			"class":        "c__",
			"order":        "o__",
			"family":       "f__",
			"genus":        "g__",
			"species":      "s__",
			"subspecies":   "S__",
		}
		fn = func(line string) (interface{}, bool, error) {
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}
			id, e := strconv.Atoi(line)
			if e != nil {
				return nil, false, e
			}

			lineage := []string{}
			var l string
			var child, parent int32
			var ok bool
			child = int32(id)
			for true {
				parent, ok = tree[child]
				if !ok {
					break
				}
				if !formatedRank || ranks[child] != "no rank" {
					l = names[child]
					if formatedRank {
						l = rankPrefix[ranks[child]] + l
					}
					lineage = append(lineage, l)
				}
				if parent == 1 && child != 1 {
					if !formatedRank || ranks[parent] != "no rank" {
						l = names[child]
						if formatedRank {
							l = rankPrefix[ranks[parent]] + l
						}
						lineage = append(lineage, l)
					}
					break
				}
				child = parent
			}
			child = int32(id)
			return taxid2lineage{child, strings.Join(ReverseStringSlice(lineage), ";")}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var t2l taxid2lineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					t2l = data.(taxid2lineage)
					outfh.WriteString(fmt.Sprintf("%d\t%s\n", t2l.taxid, t2l.lineage))
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(lineageCmd)

	lineageCmd.Flags().StringP("nodes", "", "nodes.dmp", "nodes.dmp file")
	lineageCmd.Flags().StringP("names", "", "names.dmp", "names.dmp file")
	lineageCmd.Flags().BoolP("formated-rank", "f", false, "show formated rank")
}
