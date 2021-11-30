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
	"sync"

	"github.com/shenwei356/bio/taxdump"
	"github.com/shenwei356/util/pathutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
	"github.com/twotwotwo/sorts"
)

// profile2camiCmd represents the fx2tab command
var profile2camiCmd = &cobra.Command{
	Use:   "profile2cami",
	Short: "Convert metagenomic profile table to CAMI format",
	Long: `Convert metagenomic profile table to CAMI format

Input format: 
  1. The input file should be tab-delimited
  2. At least two columns needed:
     a) TaxId of taxon at species or lower rank.
     b) Abundance (could be percentage, automatically detected).

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		sampleID := getFlagString(cmd, "sample-id")
		taxonomyID := getFlagString(cmd, "taxonomy-id")
		fieldTaxid := getFlagPositiveInt(cmd, "taxid-field") - 1
		fieldAbd := getFlagPositiveInt(cmd, "abundance-field") - 1
		keepZero := getFlagBool(cmd, "keep-zero")
		var usePercentage bool
		// usePercentage := getFlagBool(cmd, "percentage")

		showRanks := getFlagStringSlice(cmd, "show-rank")

		maxField := fieldTaxid + 1
		if fieldAbd > fieldTaxid {
			maxField = fieldAbd + 1
		}

		files := getFileList(args)

		if len(files) > 1 {
			checkError(fmt.Errorf("only one input file allowed"))
		}

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		// ----------------------------------------------------------------

		var err error

		var taxdb *taxdump.Taxonomy

		if config.Verbose {
			log.Infof("loading Taxonomy from: %s", config.DataDir)
		}

		taxdb, err = taxdump.NewTaxonomyWithRankFromNCBI(config.NodesFile)
		if err != nil {
			checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
		}

		if config.Verbose {
			log.Infof("  %d nodes in %d ranks loaded", len(taxdb.Nodes), len(taxdb.Ranks))
		}

		var existed bool

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			err = taxdb.LoadNamesFromNCBI(config.NamesFile)
			if err != nil {
				checkError(fmt.Errorf("err on loading Taxonomy names: %s", err))
			}

			if config.Verbose {
				log.Infof("  %d names loaded", len(taxdb.Names))
			}
		}()

		go func() {
			defer wg.Done()
			existed, err = pathutil.Exists(config.DelNodesFile)
			if err != nil {
				checkError(fmt.Errorf("err on checking file merged.dmp: %s", err))
			}
			if existed {
				err = taxdb.LoadDeletedNodesFromNCBI(config.DelNodesFile)
				if err != nil {
					checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
				}
			}
			if config.Verbose {
				log.Infof("  %d deleted nodes loaded", len(taxdb.DelNodes))
			}
		}()

		go func() {
			defer wg.Done()
			existed, err = pathutil.Exists(config.MergedFile)
			if err != nil {
				checkError(fmt.Errorf("err on checking file merged.dmp: %s", err))
			}
			if existed {
				err = taxdb.LoadMergedNodesFromNCBI(config.MergedFile)
				if err != nil {
					checkError(fmt.Errorf("err on loading Taxonomy merged nodes: %s", err))
				}
			}
			if config.Verbose {
				log.Infof("  %d merged nodes loaded", len(config.MergedFile))
			}
		}()

		wg.Wait()
		taxdb.CacheLCA()

		// ----------------------------------------------------------------

		showRanksMap := make(map[string]interface{}, 128)
		for _, _rank := range showRanks {
			showRanksMap[_rank] = struct{}{}
		}
		rankOrder := make(map[string]int, len(showRanks))
		for _i, _r := range showRanks {
			rankOrder[_r] = _i
		}

		// ----------------------

		targets := make([]*Target, 0, 512)

		items := make([]string, maxField)
		// var line string
		var _taxid int
		var taxid uint32
		var abd float64
		var sum float64

		file := files[0]

		fh, err := xopen.Ropen(file)
		checkError(err)

		scanner := bufio.NewScanner(fh)

		for scanner.Scan() {
			stringSplitN(scanner.Text(), "\t", maxField, &items)
			if len(items) < maxField {
				continue
			}

			_taxid, err = strconv.Atoi(items[fieldTaxid])
			if err != nil {
				checkError(fmt.Errorf("failt to parse taxid: %s", items[fieldTaxid]))
			}
			taxid = uint32(_taxid)

			abd, err = strconv.ParseFloat(items[fieldAbd], 64)
			if err != nil {
				checkError(fmt.Errorf("failt to parse abundance: %s", items[fieldAbd]))
			}

			if !keepZero && abd == 0 {
				continue
			}

			targets = append(targets, &Target{Taxid: taxid, Abundance: abd})
			sum += abd
		}

		if err := scanner.Err(); err != nil {
			checkError(err)
		}
		checkError(fh.Close())

		if config.Verbose {
			if sum > 1 {
				log.Infof("%d taxons given, sum of abundance in percentage: %.6f", len(targets), sum)
				usePercentage = true
			} else {
				log.Infof("%d taxons given, sum of abundance : %.6f", len(targets), sum)
			}
		}

		sorts.Quicksort(Targets(targets))

		// add taxonomy info
		for _, target := range targets {
			target.AddTaxonomy(taxdb, showRanksMap, target.Taxid)
		}

		profile := generateProfile(taxdb, targets)

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

		// ----------------------------------------------------------------

		// cami format
		// https://github.com/bioboxes/rfc/blob/master/data-format/profiling.mkd

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		outfh.WriteString(fmt.Sprintf("@SampleID:%s\n", sampleID))
		outfh.WriteString("@Version:0.10.0\n")
		outfh.WriteString("@Ranks:superkingdom|phylum|class|order|family|genus|species|strain\n")
		outfh.WriteString(fmt.Sprintf("@TaxonomyID:%s\n", taxonomyID))
		outfh.WriteString("@@TAXID\tRANK\tTAXPATH\tTAXPATHSN\tPERCENTAGE\n")

		var lineageTaxids, lineageNames string
		filterByRank := len(showRanksMap) > 0
		names := make([]string, 0, 8)
		taxids := make([]string, 0, 8)
		var ok bool
		var percentage float64
		for _, node := range nodes {
			if filterByRank {
				if _, ok = showRanksMap[taxdb.Rank(node.Taxid)]; !ok {
					continue
				}

				names = names[:0]
				taxids = taxids[:0]
				for i, taxid := range node.LineageTaxids {
					if _, ok = showRanksMap[taxdb.Rank(taxid)]; ok {
						taxids = append(taxids, strconv.Itoa(int(taxid)))
						names = append(names, node.LineageNames[i])
					}
				}
				lineageTaxids = strings.Join(taxids, "|")
				lineageNames = strings.Join(names, "|")
			} else {
				taxids = taxids[:0]
				for _, taxid := range node.LineageTaxids {
					taxids = append(taxids, strconv.Itoa(int(taxid)))
				}
				lineageTaxids = strings.Join(taxids, "|")
				lineageNames = strings.Join(node.LineageNames, "|")
			}

			if usePercentage {
				percentage = node.Abundance
			} else {
				percentage = node.Abundance * 100
			}

			outfh.WriteString(fmt.Sprintf("%d\t%s\t%s\t%s\t%.15f\n",
				node.Taxid, node.Rank, lineageTaxids, lineageNames, percentage))
		}
	},
}

func init() {
	RootCmd.AddCommand(profile2camiCmd)

	profile2camiCmd.Flags().StringP("sample-id", "s", "", `sample ID in result file`)
	profile2camiCmd.Flags().StringP("taxonomy-id", "t", "", `taxonomy ID in result file`)
	profile2camiCmd.Flags().IntP("taxid-field", "i", 1, "field index of taxid. input data should be tab-separated")
	profile2camiCmd.Flags().IntP("abundance-field", "a", 2, "field index of abundance. input data should be tab-separated")
	profile2camiCmd.Flags().StringSliceP("show-rank", "r", []string{"superkingdom", "phylum", "class", "order", "family", "genus", "species", "strain"}, "only show TaxIds and names of these ranks")
	profile2camiCmd.Flags().BoolP("keep-zero", "0", false, "keep taxon with abundance of zero")
	// profile2camiCmd.Flags().BoolP("percentage", "p", false, "abundance is in percentage")
}
