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
	"sort"
	"strconv"
	"strings"

	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// camiFilterCmd represents the fx2tab command
var camiFilterCmd = &cobra.Command{
	Use:   "cami-filter",
	Short: "Remove taxa of given TaxIds and their descendants in CAMI metagenomic profile table",
	Long: `Remove taxa of given TaxIds and their descendants in CAMI metagenomic profile table

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		taxidsStr := getFlagStringSlice(cmd, "taxids")
		// if len(taxidsStr) == 0 {
		// 	checkError(fmt.Errorf("flag --taxids needed"))
		// }

		filter := make(map[string]interface{}, len(taxidsStr))
		for _, t := range taxidsStr {
			filter[t] = struct{}{}
		}

		fieldTaxid := getFlagPositiveInt(cmd, "field-taxid") - 1
		fieldRank := getFlagPositiveInt(cmd, "field-rank") - 1
		fieldTaxpath := getFlagPositiveInt(cmd, "field-taxpath") - 1
		fieldTaxpathSN := getFlagPositiveInt(cmd, "field-taxpathsn") - 1
		fieldPercentage := getFlagPositiveInt(cmd, "field-percentage") - 1

		taxidSep := getFlagString(cmd, "taxid-sep")
		if taxidSep == "" {
			checkError(fmt.Errorf("flag --taxid-sep needed and should not be empty"))
		}

		files := getFileList(args)

		if len(files) > 1 {
			checkError(fmt.Errorf("only one input file allowed"))
		}

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		showRanks := getFlagStringSlice(cmd, "show-rank")

		showRanksMap := make(map[string]interface{}, 128)
		for _, _rank := range showRanks {
			showRanksMap[_rank] = struct{}{}
		}
		rankOrder := make(map[string]int, len(showRanks))
		for _i, _r := range showRanks {
			rankOrder[_r] = _i
		}

		leavesRanks := getFlagStringSlice(cmd, "leave-ranks")
		leavesRanksMap := make(map[string]interface{}, len(leavesRanks))
		for _, r := range leavesRanks {
			leavesRanksMap[r] = struct{}{}
		}

		// ----------------------------------------------------------------

		idx := []int{fieldTaxid, fieldRank, fieldTaxpath, fieldTaxpathSN, fieldPercentage}

		maxField := fieldTaxid
		for _, i := range idx[1:] {
			if i > maxField {
				maxField = i
			}
		}
		maxField++

		n := maxField + 1

		file := files[0]

		fh, err := xopen.Ropen(file)
		checkError(err)

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		scanner := bufio.NewScanner(fh)

		items := make([]string, n)
		var _line, line string
		var _taxid int
		var taxid uint32
		var taxids []string
		var taxidsUint []uint32
		var percenage float64
		var rank, taxpath, taxpathsn string
		var hasData bool

		var taxidS string
		var ok bool
		var skipThis bool

		rankMap := make(map[uint32]string, 1024)
		meta := make([]string, 0, 8)

		targets := make([]*Target, 0, 512)

		for scanner.Scan() {
			line = scanner.Text()
			if line == "" {
				continue
			}

			// meta date
			if line[0] == '#' || line[0] == '@' {
				if hasData { // new record, need to summarize and output
					targets1 := filterLeaves(rankMap, leavesRanksMap, targets)

					profile := generateProfile2(targets, targets1)

					nodes := make([]*ProfileNode, 0, len(profile))
					for _, node := range profile {
						nodes = append(nodes, node)
					}

					sort.Slice(nodes, func(i, j int) bool {
						if rankOrder[nodes[i].Rank] < rankOrder[nodes[j].Rank] {
							return true
						}
						if rankOrder[nodes[i].Rank] == rankOrder[nodes[j].Rank] {
							return nodes[i].Abundance > nodes[j].Abundance
						}
						return false
					})

					for _, _line = range meta {
						outfh.WriteString(_line + "\n")
					}
					for _, node := range nodes {
						taxids = taxids[:0]
						for _, taxid = range node.LineageTaxids {
							taxids = append(taxids, strconv.Itoa(int(taxid)))
						}

						fmt.Fprintf(outfh, "%d\t%s\t%s\t%s\t%.15f\n",
							node.Taxid,
							node.Rank,
							strings.Join(taxids, taxidSep),
							strings.Join(node.LineageNames, taxidSep),
							node.Abundance,
						)
					}

					rankMap = make(map[uint32]string, 1024)
					meta = meta[:0]
					meta = append(meta, line)
					targets = targets[:0]
					hasData = false
					continue
				}

				meta = append(meta, line)
				continue
			}

			stringSplitN(line, "\t", n, &items)
			if len(items) < maxField {
				continue
			}

			percenage, err = strconv.ParseFloat(items[fieldPercentage], 64)
			if err != nil {
				checkError(fmt.Errorf("failed to parse abundance: %s", items[fieldPercentage]))
			}

			if percenage == 0 {
				continue
			}

			hasData = true

			_taxid, err = strconv.Atoi(items[fieldTaxid])
			if err != nil {
				checkError(fmt.Errorf("failed to parse taxid: %s", items[fieldTaxid]))
			}
			taxid = uint32(_taxid)

			rank = items[fieldRank]
			rankMap[taxid] = rank
			taxpath = items[fieldTaxpath]
			taxpathsn = items[fieldTaxpathSN]

			taxids = strings.Split(taxpath, taxidSep)
			skipThis = false
			for _, taxidS = range taxids {
				if _, ok = filter[taxidS]; ok {
					skipThis = true
					break
				}
			}
			if skipThis {
				continue
			}

			taxidsUint = make([]uint32, 0, len(taxids))
			for _, taxidS = range taxids {
				_taxid, err = strconv.Atoi(taxidS)
				if err != nil {
					checkError(fmt.Errorf("failed to parse taxid: %s", taxidS))
				}
				taxidsUint = append(taxidsUint, uint32(_taxid))
			}

			targets = append(targets, &Target{
				Taxid:     taxid,
				Abundance: percenage,

				Rank:          rank,
				TaxonName:     "",
				LineageNames:  strings.Split(taxpathsn, taxidSep),
				LineageTaxids: taxids,

				CompleteLineageTaxids: taxidsUint,
			})

		}

		if err := scanner.Err(); err != nil {
			checkError(err)
		}
		checkError(fh.Close())

		if hasData { // new record, need to summarize and output
			targets1 := filterLeaves(rankMap, leavesRanksMap, targets)

			profile := generateProfile2(targets, targets1)

			nodes := make([]*ProfileNode, 0, len(profile))
			for _, node := range profile {
				nodes = append(nodes, node)
			}

			sort.Slice(nodes, func(i, j int) bool {
				if rankOrder[nodes[i].Rank] < rankOrder[nodes[j].Rank] {
					return true
				}
				if rankOrder[nodes[i].Rank] == rankOrder[nodes[j].Rank] {
					return nodes[i].Abundance > nodes[j].Abundance
				}
				return false
			})

			for _, _line = range meta {
				outfh.WriteString(_line + "\n")
			}
			for _, node := range nodes {
				taxids = taxids[:0]
				for _, taxid = range node.LineageTaxids {
					taxids = append(taxids, strconv.Itoa(int(taxid)))
				}

				fmt.Fprintf(outfh, "%d\t%s\t%s\t%s\t%.15f\n",
					node.Taxid,
					node.Rank,
					strings.Join(taxids, taxidSep),
					strings.Join(node.LineageNames, taxidSep),
					node.Abundance,
				)
			}
		}

	},
}

func generateProfile2(targets0, targets []*Target) map[uint32]*ProfileNode {

	targetsMap := make(map[uint32]*Target, len(targets0))
	for _, target := range targets0 {
		targetsMap[target.Taxid] = target
	}

	profile := make(map[uint32]*ProfileNode, len(targets))

	var target0 *Target
	for _, target := range targets {
		for _, taxid := range target.CompleteLineageTaxids {
			if node, ok := profile[taxid]; !ok {
				target0 = targetsMap[taxid]

				profile[taxid] = &ProfileNode{
					Taxid:         taxid,
					Rank:          target0.Rank,
					TaxonName:     target0.TaxonName,
					LineageNames:  target0.LineageNames,
					LineageTaxids: target0.CompleteLineageTaxids,

					Abundance: target.Abundance,
				}
			} else {
				node.Abundance += target.Abundance
			}
		}
	}

	return profile
}

func filterLeaves(rankMap map[uint32]string, leavesRanksMap map[string]interface{}, targets []*Target) []*Target {

	targetsMap := make(map[uint32]*Target, len(targets))
	tree := make(map[uint32]map[uint32]uint32, 1024)

	var i int
	var taxidP, taxid uint32
	var ok bool
	for _, target := range targets {
		targetsMap[target.Taxid] = target

		for i, taxid = range target.CompleteLineageTaxids {
			if i == 0 {
				taxidP = 1
			} else {
				taxidP = target.CompleteLineageTaxids[i-1]
			}

			if _, ok = tree[taxidP]; !ok {
				tree[taxidP] = make(map[uint32]uint32, 8)
			}
			tree[taxidP][taxid] = target.Taxid
		}
	}
	// leaves := make([]uint32, 0, 1024)
	leaves := make([]*Target, 0, 1024)
	for _, m := range tree {
		for taxid = range m {
			if _, ok = tree[taxid]; !ok {
				if _, ok = leavesRanksMap[rankMap[taxid]]; ok {
					// leaves = append(leaves, taxid)
					leaves = append(leaves, targetsMap[m[taxid]])
				}
			}
		}
	}

	// recompute abundance
	var sum float64
	for _, target := range leaves {
		sum += target.Abundance
	}

	for _, target := range leaves {
		target.Abundance = target.Abundance / sum * 100
	}

	return leaves
}

func init() {
	RootCmd.AddCommand(camiFilterCmd)

	camiFilterCmd.Flags().IntP("field-taxid", "", 1, "field index of taxid")
	camiFilterCmd.Flags().IntP("field-rank", "", 2, "field index of taxid")
	camiFilterCmd.Flags().IntP("field-taxpath", "", 3, "field index of TAXPATH")
	camiFilterCmd.Flags().IntP("field-taxpathsn", "", 4, "field index of TAXPATHSN")
	camiFilterCmd.Flags().IntP("field-percentage", "", 5, "field index of PERCENTAGE")

	camiFilterCmd.Flags().StringP("taxid-sep", "", "|", "separator of taxid in TAXPATH and TAXPATHSN")

	camiFilterCmd.Flags().StringSliceP("taxids", "t", []string{}, "the parent taxid(s) to filter out")

	camiFilterCmd.Flags().StringSliceP("show-rank", "", []string{"superkingdom", "phylum", "class", "order", "family", "genus", "species", "strain"}, "only show TaxIds and names of these ranks")
	camiFilterCmd.Flags().StringSliceP("leave-ranks", "", []string{"species", "strain"}, "only consider leaves at these ranks")
}
