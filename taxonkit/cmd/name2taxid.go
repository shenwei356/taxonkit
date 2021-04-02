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
	"fmt"
	"strings"
	"sync"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// name2taxidCmd represents the fx2tab command
var name2taxidCmd = &cobra.Command{
	Use:   "name2taxid",
	Short: "Convert scientific names to TaxIds",
	Long: `Convert scientific names to TaxIds

Attention:

  1. Some TaxIds share the same scientific names, e.g, Drosophila.
     These input lines are duplicated with multiple TaxIds.

    $ echo Drosophila | taxonkit name2taxid | taxonkit lineage -i 2 -r -L
    Drosophila      7215    genus
    Drosophila      32281   subgenus
    Drosophila      2081351 genus

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		printRank := getFlagBool(cmd, "show-rank")
		field := getFlagPositiveInt(cmd, "name-field") - 1
		limite2SciName := getFlagBool(cmd, "sci-name")

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		var m map[string][]uint32

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			if config.Verbose {
				log.Infof("parsing names file: %s", config.NamesFile)
			}
			m = getTaxonName2Taxids(config.NamesFile, limite2SciName)
			if config.Verbose {
				log.Infof("%d names parsed", len(m))
			}
			wg.Done()
		}()

		var ranks map[uint32]string

		if printRank {
			wg.Add(1)
			go func() {
				if config.Verbose {
					log.Infof("parsing nodes file: %s", config.NodesFile)
				}
				ranks = getRanks(config.NodesFile)
				if config.Verbose {
					log.Infof("%d nodes parsed", len(ranks))
				}
				wg.Done()
			}()
		}

		wg.Wait()

		// ----------------------------------------------------------

		type line2taxids struct {
			line   string
			taxids []uint32
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

		var taxid uint32
		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var l2t line2taxids
			var data interface{}
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data = range chunk.Data {
					l2t = data.(line2taxids)
					if len(l2t.taxids) == 0 {
						if printRank {
							outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\n", l2t.line, "", ""))
						} else {
							outfh.WriteString(fmt.Sprintf("%s\t%s\n", l2t.line, ""))
						}
						if config.LineBuffered {
							outfh.Flush()
						}
					} else {
						for _, taxid = range l2t.taxids {
							if printRank {
								outfh.WriteString(fmt.Sprintf("%s\t%d\t%s\n", l2t.line, taxid, ranks[taxid]))
							} else {
								outfh.WriteString(fmt.Sprintf("%s\t%d\n", l2t.line, taxid))
							}
							if config.LineBuffered {
								outfh.Flush()
							}
						}
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(name2taxidCmd)
	name2taxidCmd.Flags().BoolP("show-rank", "r", false, `show rank`)
	name2taxidCmd.Flags().IntP("name-field", "i", 1, "field index of name. data should be tab-separated")
	name2taxidCmd.Flags().BoolP("sci-name", "s", false, "only searching scientific names")
}
