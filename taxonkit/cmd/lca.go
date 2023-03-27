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
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shenwei356/util/bytesize"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// lcaCmd represents
var lcaCmd = &cobra.Command{
	Use:   "lca",
	Short: "Compute lowest common ancestor (LCA) for TaxIds",
	Long: `Compute lowest common ancestor (LCA) for TaxIds

Attention:

  1. This command computes LCA TaxId for a list of TaxIds 
     in a field ("-i/--taxids-field) of tab-delimited file or STDIN.
  2. TaxIDs should have the same separator ("-s/--separator"),
     single charactor separator is prefered.
  3. Empty lines or lines without valid TaxIds in the field are omitted.
  4. If some TaxIds are not found in database, it returns 0.
  
Examples:

    $ echo 239934, 239935, 349741 | taxonkit lca  -s ", "
    239934, 239935, 349741  239934

    $ time echo 239934  239935  349741 9606  | taxonkit lca
    239934 239935 349741 9606       131567

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		var err error

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		field := getFlagPositiveInt(cmd, "taxids-field") - 1

		separater := getFlagString(cmd, "separater")
		separator := getFlagString(cmd, "separator")
		if separater == "" && separator == "" {
			checkError(fmt.Errorf("flag -s (--separator) needed"))
		}

		if cmd.Flags().Lookup("separater").Changed && !cmd.Flags().Lookup("separator").Changed { // using --separater
			separator = separater
		}

		skipDeleted := getFlagBool(cmd, "skip-deleted")
		skipUnfound := getFlagBool(cmd, "skip-unfound")

		bufferSizeS := getFlagString(cmd, "buffer-size")
		if bufferSizeS == "" {
			checkError(fmt.Errorf("value of buffer size. supported unit: K, M, G"))
		}
		bufferSize, err := bytesize.ParseByteSize(bufferSizeS)
		if err != nil {
			checkError(fmt.Errorf("invalid value of buffer size. supported unit: K, M, G"))
		}

		taxondb := loadTaxonomy(&config, false)
		nodes := taxondb.Nodes
		merged := taxondb.MergeNodes
		delnodes := taxondb.DelNodes

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		buf := make([]byte, bufferSize)

		taxids := make([]uint32, 0, 128)
		for _, file := range files {
			fh, err := xopen.Ropen(file)
			checkError(err)

			scanner := bufio.NewScanner(fh)
			scanner.Buffer(buf, int(bufferSize))

			var _taxid int
			var line, item string
			var items []string
			var lca, taxid, taxid2 uint32
			var ok, flag bool
			for scanner.Scan() {
				line = strings.Trim(scanner.Text(), "\r\n ")
				if line == "" {
					continue
				}

				lca = 0

				items = strings.Split(line, "\t")
				if len(items) <= field {
					field = len(items) - 1
				}

				if items[field] == "" {
					continue
				}

				items = strings.Split(items[field], separator)

				taxids = taxids[:0]

				flag = false
				for _, item = range items {
					item = reNonTaxid.ReplaceAllString(item, "")
					if item == "" {
						continue
					}

					_taxid, _ = strconv.Atoi(item)
					taxid = uint32(_taxid)

					_, ok = nodes[taxid]
					if ok {
						taxids = append(taxids, taxid)
						continue
					}

					if _, ok = delnodes[taxid]; ok {
						log.Warningf("taxid %d was deleted", taxid)
						if !skipDeleted {
							flag = true
							break
						}
						continue
					}
					if taxid2, ok = merged[taxid]; ok {
						log.Warningf("taxid %d was merged into %d", taxid, taxid2)
						taxid = taxid2
						taxids = append(taxids, taxid)
					} else {
						log.Warningf("taxid %d not found", taxid)
						if !skipUnfound {
							flag = true
							break
						}
					}
				}
				if flag {
					outfh.WriteString(fmt.Sprintf("%s\t%d\n", line, 0))
					continue
				}

				switch len(taxids) {
				case 0:
					continue
				case 1:
					lca = taxids[0]
				default:
					lca = taxids[0]
					for _, taxid = range taxids {
						lca = taxondb.LCA(lca, taxid)
					}
				}

				outfh.WriteString(fmt.Sprintf("%s\t%d\n", line, lca))
			}
			if err := scanner.Err(); err != nil {
				checkError(err)
			}

			checkError(fh.Close())
		}

	},
}

func init() {
	RootCmd.AddCommand(lcaCmd)

	lcaCmd.Flags().IntP("taxids-field", "i", 1, "field index of TaxIds. Input data should be tab-separated")

	lcaCmd.Flags().StringP("separater", "", " ", "separater for TaxIds. This flag is same to --separator.")
	lcaCmd.Flags().StringP("separator", "s", " ", "separator for TaxIds")
	lcaCmd.Flags().BoolP("skip-deleted", "D", false, "skip deleted TaxIds and compute with left ones")
	lcaCmd.Flags().BoolP("skip-unfound", "U", false, "skip unfound TaxIds and compute with left ones")
	lcaCmd.Flags().StringP("buffer-size", "b", "1M", `size of line buffer, supported unit: K, M, G. You need increase the value when "bufio.Scanner: token too long" error occured`)

}

var reTaxid = regexp.MustCompile(`^\d+$`)
var reNonTaxid = regexp.MustCompile(`\D+`)
