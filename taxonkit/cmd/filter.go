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
	Short: "Filter TaxIds by taxonomic rank range",
	Long: `Filter TaxIds by taxonomic rank range

Attentions:

  1. Flag -L/--lower-than and -H/--higher-than are exclusive, and can be
     used along with -E/--equal-to which values can be different.
  2. A list of pre-ordered ranks is in ~/.taxonkit/ranks.txt, you can use
     your list by -r/--rank-file, the format specification is below.
  3. All ranks in taxonomy database should be defined in rank file.
  4. Ranks can be removed with black list via -B/--black-list.
  5. TaxIDs with no rank can be optionally discarded by -N/--discard-noranks.
  6. But when filtering with -L/--lower-than, you can use
    -n/--save-predictable-norank to save some special ranks without order,
    where rank of the closest higher node is still lower than rank cutoff.

Rank file:

  1. Blank lines or lines starting with "#" are ignored.
  2. Ranks are in decending order and case ignored.
  3. Ranks with same order should be in one line separated with comma (",", no space).
  4. Ranks without order should be assigned a prefix symbol "!" for each rank.

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		var err error

		files := getFileList(args)

		listOrder := getFlagBool(cmd, "list-order")
		listRanks := getFlagBool(cmd, "list-ranks")

		if !(listOrder || listRanks) && len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		rankFile := getFlagString(cmd, "rank-file")

		discardNoRank := getFlagBool(cmd, "discard-noranks")
		saveNorank := getFlagBool(cmd, "save-predictable-norank")

		blackListRanks := getFlagStringSlice(cmd, "black-list")

		rootTaxid := getFlagUint32(cmd, "root-taxid")
		discardRoot := getFlagBool(cmd, "discard-root")

		higher := strings.ToLower(getFlagString(cmd, "higher-than"))
		lower := strings.ToLower(getFlagString(cmd, "lower-than"))
		equalsS := getFlagStringSlice(cmd, "equal-to")
		equals := make([]string, 0, len(equalsS))
		for _, val := range equalsS {
			equals = append(equals, strings.ToLower(val))
		}

		field := getFlagPositiveInt(cmd, "taxid-field") - 1

		if higher != "" && lower != "" {
			checkError(fmt.Errorf("-H/--higher-than and -L/--lower-than can't be simultaneous given"))
		}

		if saveNorank {
			discardNoRank = true

			if lower == "" {
				checkError(fmt.Errorf("flag -n/--save-predictable-norank only works along with -L/--lower-than"))
			}
		}

		rankOrder, noRanks, err := readRankOrder(config, rankFile)
		checkError(errors.Wrap(err, rankFile))

		noRanksList := make([]string, 0, len(noRanks))
		for r := range noRanks {
			noRanksList = append(noRanksList, r)
		}

		if config.Verbose {
			log.Infof("ranks without order: %s", strings.Join(noRanksList, ", "))
		}

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

		tmp := make([]string, 0, len(blackListRanks))
		for _, r := range blackListRanks {
			if r == "" {
				continue
			}
			tmp = append(tmp, strings.ToLower(r))
			blackListRanks = tmp
		}

		if config.Verbose {
			if discardNoRank {
				log.Debugf("ranks without order will be discarded: %s", strings.Join(noRanksList, ", "))
			}
			if discardRoot {
				log.Debugf("root rank without order will be discarded")
			}
			if len(blackListRanks) > 0 {
				log.Debugf("ranks in black list will be discarded: %s", strings.Join(blackListRanks, ", "))
			}
		}

		filter, err := newRankFilter(taxondb, rankOrder, noRanks, lower, higher, equals, blackListRanks, discardNoRank, saveNorank)
		checkError(err)

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		for _, file := range files {
			fh, err := xopen.Ropen(file)
			checkError(err)

			var line string
			var items []string

			scanner := bufio.NewScanner(fh)
			var _taxid int
			var taxid uint32
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

				pass, err = filter.isPassed(taxid)
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
	filterCmd.Flags().BoolP("list-order", "", false, `list user defined ranks in order, from "$HOME/.taxonkit/ranks.txt"`)
	filterCmd.Flags().BoolP("list-ranks", "", false, `list ordered ranks in taxonomy database, sorted in user defined order`)

	filterCmd.Flags().BoolP("discard-noranks", "N", false, `discard all ranks without order, type "taxonkit filter --help" for details`)
	filterCmd.Flags().BoolP("save-predictable-norank", "n", false, `do not discard some special ranks without order when using -L, where rank of the closest higher node is still lower than rank cutoff`)
	filterCmd.Flags().StringSliceP("black-list", "B", []string{}, `black list of ranks to discard, e.g., '-B "no rank" -B "clade"`)

	filterCmd.Flags().BoolP("discard-root", "R", false, `discard root taxid, defined by --root-taxid`)
	filterCmd.Flags().Uint32P("root-taxid", "", 1, `root taxid`)

	filterCmd.Flags().StringP("lower-than", "L", "", "output TaxIds with rank lower than a rank, exclusive with --higher-than")
	filterCmd.Flags().StringP("higher-than", "H", "", "output TaxIds with rank higher than a rank, exclusive with --lower-than")
	filterCmd.Flags().StringSliceP("equal-to", "E", []string{}, `output TaxIds with rank equal to some ranks, multiple values can be separated with comma "," (e.g., -E "genus,species"), or give multiple times (e.g., -E genus -E species)`)

	filterCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. input data should be tab-separated")
}
