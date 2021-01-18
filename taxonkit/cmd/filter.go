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
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
	"github.com/twotwotwo/sorts"
)

// filterCmd represents
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter taxIDs by taxonomic rank range",
	Long: `Filter taxIDs by taxonomic rank range

Attentions:
  1. Flag -L/--lower-than and -H/--higher-than are exclusive, and can be
     used along with -E/--equal-to which values can be different.
  2. A list of pre-ordered ranks is in ~/.taxonkit/ranks.txt, you can give
     your list by -r/--rank-file, the format specification is below.
  3. TaxIDss with no rank will be discarded.

Rank file:
  1. Blank lines or lines starting with "#" are ignored.
  2. Ranks are in decending order and case ignored.
  3. Ranks with same order should be in one line separated with comma (",", no space).
  4. Ranks without order should be assigning a prefix symbol "!" for each rank.

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		var err error

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		rankFile := getFlagString(cmd, "rank-file")

		discardNoRank := getFlagBool(cmd, "discard-noranks")
		blackListRanks := getFlagStringSlice(cmd, "black-list")

		rootTaxid := getFlagUint32(cmd, "root-taxid")
		discardRoot := getFlagBool(cmd, "discard-root")

		higher := strings.ToLower(getFlagString(cmd, "higher-than"))
		lower := strings.ToLower(getFlagString(cmd, "lower-than"))
		equal := strings.ToLower(getFlagString(cmd, "equal-to"))

		listOrder := getFlagBool(cmd, "list-order")
		listRanks := getFlagBool(cmd, "list-ranks")

		field := getFlagPositiveInt(cmd, "taxid-field") - 1

		if higher != "" && lower != "" {
			checkError(fmt.Errorf("-H/--higher-than and -L/--lower-than can't be simultaneous given"))
		}

		rankOrder, noRanks, err := readRankOrder(config, rankFile)
		checkError(errors.Wrap(err, rankFile))

		if listOrder {
			orders := make([]stringutil.StringCount, 0, len(rankOrder))
			for r, o := range rankOrder {
				orders = append(orders, stringutil.StringCount{Key: r, Count: o})
			}
			sorts.Quicksort(stringutil.ReversedStringCountList{orders})
			preOrder := -1
			for _, order := range orders {
				// fmt.Printf("%d\t%s\n", order.Count, order.Key)
				if order.Count == preOrder {
					fmt.Printf(",%s", order.Key)
				} else {
					if preOrder != -1 {
						fmt.Println()
					}
					fmt.Printf("%s", order.Key)
					preOrder = order.Count
				}
			}
			fmt.Println()
			return
		}

		taxondb := loadTaxonomy(&config, true)

		if config.Verbose {
			log.Infof("checking defined taxonomic rank order")
		}
		notDefined := make([]string, 0, 10)
		for rank := range taxondb.Ranks {
			if _, ok := rankOrder[rank]; !ok {
				if _, ok := noRanks[rank]; !ok {
					notDefined = append(notDefined, rank)
				}
			}
		}
		if len(notDefined) > 0 {
			checkError(fmt.Errorf("rank order not defined in rank file: %s", strings.Join(notDefined, ", ")))
		}
		if config.Verbose {
			log.Infof("checking defined taxonomic rank order passed")
		}

		if listRanks {
			orders := make([]stringutil.StringCount, 0, len(taxondb.Ranks))
			var ok bool
			for rank := range taxondb.Ranks {
				if _, ok = rankOrder[rank]; !ok {
					if _, ok := noRanks[rank]; !ok {
						checkError(fmt.Errorf("rank order not defined: %s", rank))
					}
				}
				orders = append(orders, stringutil.StringCount{Key: rank, Count: rankOrder[rank]})
			}
			sorts.Quicksort(stringutil.ReversedStringCountList{orders})
			for _, order := range orders {
				// fmt.Printf("%d\t%s\n", order.Count, order.Key)
				fmt.Printf("%s\n", order.Key)
			}
			return
		}

		filter, err := newRankFilter(taxondb.Ranks, rankOrder, noRanks, lower, higher, equal, blackListRanks, discardNoRank)
		checkError(err)

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		for _, file := range files {
			fh, err := xopen.Ropen(file)
			checkError(err)

			var line string
			items := make([]string, 8)

			scanner := bufio.NewScanner(fh)
			var _taxid int
			var taxid uint32
			var rank string
			var pass bool
			for scanner.Scan() {
				line = strings.Trim(scanner.Text(), "\r\n ")
				if line == "" {
					continue
				}

				items = strings.Split(line, "\t")
				if len(items) <= field {
					field = len(items) - 1
				}

				if items[field] == "" {
					continue
				}

				_taxid, err = strconv.Atoi(items[field])
				if err != nil {
					continue
				}

				taxid = uint32(_taxid)

				// ----------------------------------

				if discardRoot && taxid == rootTaxid {
					continue
				}

				rank = taxondb.Rank(taxid)
				if rank == "" {
					continue
				}

				pass, err = filter.isPassed(rank)
				if err != nil {
					checkError(err)
				}

				if !pass {
					continue
				}

				outfh.WriteString(line + "\n")
			}
			if err := scanner.Err(); err != nil {
				checkError(err)
			}

			checkError(fh.Close())
		}

	},
}

func init() {
	RootCmd.AddCommand(filterCmd)

	filterCmd.Flags().StringP("rank-file", "r", "", `user-defined ordered taxonomic ranks, type "taxonkit filter --help" for details`)
	filterCmd.Flags().BoolP("list-order", "", false, "list defined ranks in order")
	filterCmd.Flags().BoolP("list-ranks", "", false, "list ordered ranks in taxonomy database")

	filterCmd.Flags().BoolP("discard-noranks", "N", false, `discard ranks without order, type "taxonkit filter --help" for details`)
	filterCmd.Flags().StringSliceP("black-list", "B", []string{"no rank", "clade"}, `black list of ranks to discard`)
	filterCmd.Flags().BoolP("discard-root", "R", false, `discard root taxid, defined by --root-taxid`)
	filterCmd.Flags().Uint32P("root-taxid", "", 1, `root taxid`)

	filterCmd.Flags().StringP("lower-than", "L", "", "output taxIDs with rank lower than a rank, exclusive with --higher-than")
	filterCmd.Flags().StringP("higher-than", "H", "", "output taxIDs with rank higher than a rank, exclusive with --lower-than")
	filterCmd.Flags().StringP("equal-to", "E", "", "output taxIDs with rank equal to a rank")

	filterCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. input data should be tab-separated")
}
