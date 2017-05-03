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
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// name2taxidCmd represents the fx2tab command
var name2taxidCmd = &cobra.Command{
	Use:   "name2taxid",
	Short: "query taxid by taxon scientific name",
	Long: `query taxid by taxon scientific name

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		printRank := getFlagBool(cmd, "show-rank")
		field := getFlagPositiveInt(cmd, "name-field") - 1

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var m map[string][]int32

		log.Infof("parsing names file: %s", config.NamesFile)
		m = getTaxonNames2Taxid(config.NamesFile, config.Threads, 10)
		log.Infof("%d names parsed", len(m))

		var ranks map[int32]string
		if printRank {
			log.Infof("parsing nodes file: %s", config.NodesFile)

			reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
			checkError(err)

			ranks = make(map[int32]string)
			var info taxonInfo
			var n int64
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					info = data.(taxonInfo)
					ranks[info.child] = info.rank
					n++
				}
			}

			log.Infof("%d nodes parsed", n)
		}

		type line2taxids struct {
			line   string
			taxids []int32
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
			return line2taxids{line, m[strings.ToLower(data[field])]}, true, nil
		}

		var taxid int32
		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var l2t line2taxids
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					l2t = data.(line2taxids)
					if len(l2t.taxids) == 0 {
						if printRank {
							outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\n", l2t.line, "", ""))
						} else {
							outfh.WriteString(fmt.Sprintf("%s\t%s\n", l2t.line, ""))
						}
					} else {
						for _, taxid = range l2t.taxids {
							if printRank {
								outfh.WriteString(fmt.Sprintf("%s\t%d\t%s\n", l2t.line, taxid, ranks[taxid]))
							} else {
								outfh.WriteString(fmt.Sprintf("%s\t%d\n", l2t.line, taxid))
							}
						}
					}
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(name2taxidCmd)
	name2taxidCmd.Flags().BoolP("show-rank", "", false, `show rank`)
	name2taxidCmd.Flags().IntP("name-field", "i", 1, "field index of name. data should be tab-separated")
}
