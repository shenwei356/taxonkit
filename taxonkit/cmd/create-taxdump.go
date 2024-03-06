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
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/shenwei356/bio/taxdump"
	"github.com/shenwei356/util/pathutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// createTaxDumpCmd represents the create-TaxDump command
var createTaxDumpCmd = &cobra.Command{
	Use:   "create-taxdump",
	Short: "Create NCBI-style taxdump files for custom taxonomy, e.g., GTDB and ICTV",
	Long: `Create NCBI-style taxdump files for custom taxonomy, e.g., GTDB and ICTV

Input format: 
  0. For GTDB taxonomy file, just use --gtdb.
     We use the numeric assembly accession as the taxon at subspecies rank.
     (without the prefix GCA_ and GCF_, and version number).
  1. The input file should be tab-delimited, at least one column is needed.
  2. Ranks can be given either via the first row or the flag --rank-names.
  3. The column containing the genome/assembly accession is recommended to
     generate TaxId mapping file (taxid.map, id -> taxid).
       -A/--field-accession,    field contaning genome/assembly accession      
       --field-accession-re,    regular expression to extract the accession
     Note that mutiple TaxIds pointing to the same accession are listed as
     comma-seperated integers. 

Attention:
  1. Duplicated taxon names wit different ranks are allowed since v0.16.0, since
     the rank and taxon name are contatenated for generating the TaxId.
  2. The generated TaxIds are not consecutive numbers, however some tools like MMSeqs2
     required this, you can use the script below for convertion:
     
     https://github.com/apcamargo/ictv-mmseqs2-protein-database/blob/master/fix_taxdump.py

  3. We only check and eliminate taxid collision within a single version of taxonomy data.
     Therefore, if you create taxid-changelog with "taxid-changelog", different taxons
     in multiple versions might have the same TaxIds and some change events might be wrong.

     So a single version of taxonomic data created by "taxonkit create-taxdump" has no problem,
     it's just the changelog might not be perfect.

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		if config.Threads > 4 {
			if config.Verbose {
				log.Infof("using a lot of threads does not always accelerate processing, 4-threads is fast enough")
			}
			config.Threads = 4
		}
		runtime.GOMAXPROCS(config.Threads)

		fAccession := getFlagNonNegativeInt(cmd, "field-accession")
		accAssubspe := getFlagBool(cmd, "field-accession-as-subspecies")

		rankNames := getFlagStringSlice(cmd, "rank-names")

		var err error

		isGTDB := getFlagBool(cmd, "gtdb")

		reGenomeIDStr := getFlagString(cmd, "field-accession-re")

		if isGTDB && !cmd.Flags().Lookup("field-accession-re").Changed {
			reGenomeIDStr = `^\w\w_(.+)$`
		}

		var reGenomeID *regexp.Regexp
		if reGenomeIDStr != "" {
			if !regexp.MustCompile(`\(.+\)`).MatchString(reGenomeIDStr) {
				checkError(fmt.Errorf(`value of --field-accession-re must contains "(" and ")"`))
			}

			reGenomeID, err = regexp.Compile(reGenomeIDStr)
			if err != nil {
				checkError(fmt.Errorf("failed to compile regular expression: %s", reGenomeIDStr))
			}
		}

		reGTDBStr := getFlagString(cmd, "gtdb-re-subs")

		var reGTDBsubspe *regexp.Regexp

		if reGTDBStr != "" {
			if !regexp.MustCompile(`\(.+\)`).MatchString(reGTDBStr) {
				checkError(fmt.Errorf(`value of --gtdb-re-subs must contains "(" and ")"`))
			}

			reGTDBsubspe, err = regexp.Compile(reGTDBStr)
			if err != nil {
				checkError(fmt.Errorf("failed to compile regular expression: %s", reGTDBStr))
			}
		}

		nulls := getFlagStringSlice(cmd, "null")

		var hasAccession bool
		var numFields int
		var numRanks int
		var useFirstRow bool
		if isGTDB {
			numFields = 2
			hasAccession = true
			rankNames = []string{"superkingdom", "phylum", "class", "order", "family", "genus", "species", "no rank"}
		} else {
			hasAccession = fAccession > 0

			if accAssubspe && !hasAccession {
				checkError(fmt.Errorf("flag -S/--field-accession-as-subspecies should be used along with -A/--field-accession "))
			}

			if len(rankNames) == 0 {
				log.Infof("I will use the first row of input as rank names")
				useFirstRow = true
			} else {
				numRanks = len(rankNames)

				if hasAccession {
					if accAssubspe {
						numFields = numRanks
					} else {
						numFields = numRanks + 1
					}
					if fAccession > numFields {
						checkError(fmt.Errorf("value of -A/--field-accession (%d) is out of range (%d columns), or inadequate rank names (%d)", fAccession, numFields, numRanks))
					}
				} else {
					numFields = numRanks
				}
			}
		}

		outDir := getFlagString(cmd, "out-dir")
		force := getFlagBool(cmd, "force")
		if outDir == "" {
			checkError(fmt.Errorf("flag -O/--out-dir is needed"))
		}
		makeOutDir(outDir, force)

		// ------------------------------------------------------------

		oldTaxdumpDir := getFlagString(cmd, "old-taxdump-dir")

		var taxdb *taxdump.Taxonomy

		if oldTaxdumpDir != "" {
			log.Infof("loading Taxonomy from: %s", oldTaxdumpDir)

			taxdb, err = taxdump.NewTaxonomyWithRankFromNCBI(filepath.Join(oldTaxdumpDir, "nodes.dmp"))
			if err != nil {
				checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
			}
			log.Infof("  %d nodes in %d ranks loaded", len(taxdb.Nodes), len(taxdb.Ranks))

			var existed bool

			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				err = taxdb.LoadNamesFromNCBI(filepath.Join(oldTaxdumpDir, "names.dmp"))
				if err != nil {
					checkError(fmt.Errorf("err on loading Taxonomy names: %s", err))
				}
				log.Infof("  %d names loaded", len(taxdb.Names))
			}()

			go func() {
				defer wg.Done()
				file := filepath.Join(oldTaxdumpDir, "delnodes.dmp")
				existed, err = pathutil.Exists(file)
				if err != nil {
					checkError(fmt.Errorf("err on checking file delnodes.dmp: %s", err))
				}
				if existed {
					err = taxdb.LoadDeletedNodesFromNCBI(file)
					if err != nil {
						checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
					}
				}
				log.Infof("  %d deleted nodes loaded", len(taxdb.DelNodes))
			}()

			go func() {
				defer wg.Done()
				file := filepath.Join(oldTaxdumpDir, "merged.dmp")
				existed, err = pathutil.Exists(file)
				if err != nil {
					checkError(fmt.Errorf("err on checking file merged.dmp: %s", err))
				}
				if existed {
					err = taxdb.LoadMergedNodesFromNCBI(file)
					if err != nil {
						checkError(fmt.Errorf("err on loading Taxonomy merged nodes: %s", err))
					}
				}
				log.Infof("  %d merged nodes loaded", len(taxdb.MergeNodes))
			}()

			wg.Wait()
			log.Info()
		}

		// ------------------------------------------------------------

		nullMap := make(map[string]interface{})
		for _, k := range nulls {
			nullMap[k] = struct{}{}
		}

		files := getFileList(args)

		var reGTDBsubspeNotCaptured bool
		var reGTDBsubspeNotCapturedExample string

		var reGenomeIDNotCaptured bool
		var reGenomeIDNotCapturedExample string

		// child -> parent
		tree := make(map[uint32]uint32, 1<<16)

		// child -> rank-id
		ranks := make(map[uint32]uint8, 1<<16)

		// child -> name
		names0 := make(map[uint32]string, 1<<16)
		// child -> rank+name
		names := make(map[uint32]string, 1<<16)

		// accession -> taxid
		acc2taxid := make(map[string]*map[uint32]interface{}, 1<<16)
		var _taxids *map[uint32]interface{}
		accIdx := make(map[string]int, 1<<16)
		var idx int

		var firstLine string

		for ifile, file := range files {
			var i int
			var prev int
			var first bool
			var taxid uint32
			var _name string
			var ok bool
			var reAssignTaxid bool

			var n int
			isFirstLine := true

			var items *[]string
			if isGTDB || len(rankNames) > 0 {
				_items := make([]string, numFields)
				items = &_items
			}

			var items7 *[]string
			if isGTDB {
				_items7 := make([]string, 7)
				items7 = &_items7
			}

			fh, err := xopen.Ropen(file)
			checkError(err)

			scanner := bufio.NewScanner(fh)
			var line string
			var items0 []string
			var j int

			for scanner.Scan() {
				n++ // line number
				line = strings.Trim(scanner.Text(), "\r\n")
				if line == "" {
					continue
				}

				if isFirstLine {
					if useFirstRow {
						items0 = strings.Split(line, "\t")
						numFields = len(items0)

						if fAccession > numFields {
							checkError(fmt.Errorf("value of -A/--field-accession (%d) is out of range (%d columns)", fAccession, numFields))
						}

						if hasAccession {
							if accAssubspe {
								numRanks = numFields

								rankNames = rankNames[:0]
								for i, val := range items0 {
									if i+1 != fAccession {
										rankNames = append(rankNames, strings.ToLower(val))
									} else {
										rankNames = append(rankNames, "subspecies")
									}
								}
							} else {
								numRanks = numFields - 1

								if numRanks == 0 {
									checkError(fmt.Errorf("at least 2 columns needed for -A/--field-accession when -S/--field-accession-as-subspecies is not given"))
								}

								rankNames = rankNames[:0]
								for i, val := range items0 {
									if i+1 != fAccession {
										rankNames = append(rankNames, strings.ToLower(val))
									}
								}
							}
						} else {
							numRanks = numFields

							rankNames = make([]string, len(items0))
							for i, val := range items0 {
								rankNames[i] = strings.ToLower(val)
							}
						}

						_items := make([]string, numFields)
						items = &_items

						if ifile > 0 { // later files, need to check whether first row match in multiple files
							if firstLine != line {
								checkError(fmt.Errorf("inconsistent rank names at the first line: %s", file))
							}
						} else {
							firstLine = line
						}

						isFirstLine = false

						continue
					}
				}

				// efficient but can't handle cases where len(items) > numFields
				// stringSplitNByByte(line, '\t', numFields, items)
				*items = strings.Split(line, "\t")

				if !isGTDB {
					if len(*items) != numFields {
						if hasAccession && !accAssubspe {
							checkError(fmt.Errorf("the number (%d, expect %d) of columns at line %d does not match #rank-names + 1 (%d+1): %s", len(*items), numFields, n, len(rankNames), file))
						} else {
							checkError(fmt.Errorf("the number (%d, expect %d) of columns at line %d does not match that of rank names (%d): %s", len(*items), numFields, n, len(rankNames), file))
						}
					}
				} else {
					if len(*items) != numFields {
						checkError(fmt.Errorf("expect %d columns, while only %d given at line %d ", numFields, len(*items), n))
					}
				}

				t := _Taxon{}
				var val string

				if isGTDB {
					if reGenomeID != nil {
						found := reGenomeID.FindAllStringSubmatch((*items)[0], 1)
						if len(found) == 0 {
							t.Accession = (*items)[0]
							// checkError(fmt.Errorf("invalid GTDB assembly accession: %s", (*items)[0]))
							reGTDBsubspeNotCaptured = true
							reGTDBsubspeNotCapturedExample = (*items)[0]
						} else {
							t.Accession = found[0][1]
						}
					}

					stringSplitNByByte(CopyString((*items)[1]), ';', 7, items7)
					if len(*items7) < 7 {
						checkError(fmt.Errorf("invalid GTDB taxonomy record: %s", line))
					}

					val = (*items7)[0]
					if len(val) < 3 || val[0:3] != "d__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (kingdom): %s", val))
					}
					kingdom := val[3:]

					val = (*items7)[1]
					if len(val) < 3 || val[0:3] != "p__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (phylum): %s", val))
					}
					phylum := val[3:]

					val = (*items7)[2]
					if len(val) < 3 || val[0:3] != "c__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (class): %s", val))
					}
					class := val[3:]

					val = (*items7)[3]
					if len(val) < 3 || val[0:3] != "o__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (order): %s", val))
					}
					order := val[3:]

					val = (*items7)[4]
					if len(val) < 3 || val[0:3] != "f__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (family): %s", val))
					}
					family := val[3:]

					val = (*items7)[5]
					if len(val) < 3 || val[0:3] != "g__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (genus): %s", val))
					}
					genus := val[3:]

					val = (*items7)[6]
					if len(val) < 3 || val[0:3] != "s__" {
						checkError(fmt.Errorf("invalid GTDB taxonomy format (species): %s", val))
					}
					species := val[3:]

					var subspe string

					if reGTDBsubspe != nil {
						found := reGTDBsubspe.FindAllStringSubmatch((*items)[0], 1)
						if len(found) == 0 {
							subspe = (*items)[0]
							// checkError(fmt.Errorf("invalid GTDB assembly accession: %s", (*items)[0]))
							reGTDBsubspeNotCaptured = true
							reGTDBsubspeNotCapturedExample = (*items)[0]
						} else {
							subspe = found[0][1]
						}
					}

					// ----------------------------------------
					t.Names = []string{kingdom, phylum, class, order, family, genus, species, subspe}
					t.TaxIds = make([]uint32, 8)

					for j = 0; j < 8; j++ {
						if _, ok = nullMap[t.Names[j]]; ok {
							continue
						}

						t.TaxIds[j] = uint32(xxhash.Sum64String(rankNames[j]+strings.ToLower(t.Names[j])) & 2147483647)
					}
				} else {
					// var ok bool

					if hasAccession {
						val = (*items)[fAccession-1]

						if reGenomeID != nil {
							found := reGenomeID.FindAllStringSubmatch(val, 1)
							if len(found) == 0 {
								t.Accession = val
								reGenomeIDNotCaptured = true
								reGenomeIDNotCapturedExample = val
							} else {
								t.Accession = found[0][1]
							}
						}

					}

					t.Names = make([]string, numRanks)
					t.TaxIds = make([]uint32, numRanks)

					if hasAccession && !accAssubspe {
						copy((*items)[fAccession-1:len(*items)-1], (*items)[fAccession:])
					}

					for j = 0; j < numRanks; j++ {
						t.Names[j] = (*items)[j]

						if _, ok = nullMap[t.Names[j]]; ok {
							continue
						}

						t.TaxIds[j] = uint32(xxhash.Sum64String(rankNames[j]+strings.ToLower(t.Names[j])) & 2147483647)
					}
				}

				// ------------------------------------

				first = true
				for i = len(t.TaxIds) - 1; i >= 0; i-- {
					taxid = t.TaxIds[i]
					if taxid == 1 { // just in case
						taxid = 2
					}

					if taxid == 0 || t.Names[i] == "" {
						continue
					}

				REASSIGNTAXID:

					reAssignTaxid = false

					if _name, ok = names[taxid]; ok { // check name
						if _name != rankNames[i]+t.Names[i] { // two names hashed to the same uint32
							if config.Verbose {
								log.Infof(`"%s" and "%s" have the same taxId: %d`, _name, rankNames[i]+t.Names[i], taxid)
							}
							reAssignTaxid = true
						} else if i > 0 { // taxa with different parents may have the same names, many cases in ICTV
							for j = i - 1; j >= 0; j-- { // find the parent
								if t.Names[j] != "" && t.TaxIds[j] != 0 {
									break
								}
							}
							if j >= 0 { // have a non-root parent
								if tree[taxid] != t.TaxIds[j] && names[tree[taxid]] != rankNames[j]+t.Names[j] { // not the same parent
									if config.Verbose {
										log.Infof(`"%s" (%d) and "%s" (%d) having the same child: %s`,
											names[tree[taxid]], tree[taxid], t.Names[j], t.TaxIds[j], _name)
									}
									reAssignTaxid = true
								}
							}
						}
					} else {
						names[taxid] = rankNames[i] + t.Names[i]
						names0[taxid] = t.Names[i]
					}

					if _, ok = ranks[taxid]; !ok {
						ranks[taxid] = uint8(i)
					}

					if reAssignTaxid {
						if config.Verbose {
							log.Infof(`  assign a new TaxId for "%s" (rank: %s): %d -> %d`, names[taxid], rankNames[i], taxid, taxid+1)
						}
						taxid++
						t.TaxIds[i] = taxid
						goto REASSIGNTAXID
					}

					if first {
						if hasAccession {
							idx++
							accIdx[t.Accession] = idx

							if _taxids, ok = acc2taxid[t.Accession]; !ok {
								acc2taxid[t.Accession] = &map[uint32]interface{}{taxid: struct{}{}}
							} else {
								(*_taxids)[taxid] = struct{}{}
							}
						}

						prev = i
						first = false
						continue
					}

					tree[t.TaxIds[prev]] = taxid
					prev = i
				}

				// the highest node
				tree[t.TaxIds[prev]] = 1
			}

			if err = scanner.Err(); err != nil {
				checkError(err)
			}
			checkError(fh.Close())
		}

		if isGTDB && reGTDBsubspeNotCaptured {
			log.Warningf("--gtdb-re-subs failed to extract ID for subspecies, the origninal value is used instead. e.g., %s", reGTDBsubspeNotCapturedExample)
		}

		// ------------------------------- taxid.map -------------------------

		if hasAccession {
			fileAcc2Taxid := filepath.Join(outDir, "taxid.map")

			if reGenomeIDNotCaptured {
				log.Warningf("--field-accession-re failed to extract genome accession, the origninal value is used instead. e.g., %s", reGenomeIDNotCapturedExample)
			}

			outfhAcc2Taxid, err := xopen.Wopen(fileAcc2Taxid)
			checkError(err)
			defer outfhAcc2Taxid.Close()

			accs := make([]string, 0, len(accIdx))
			for acc := range accIdx {
				accs = append(accs, acc)
			}
			sort.Slice(accs, func(i, j int) bool {
				return accIdx[accs[i]] < accIdx[accs[j]]
			})

			_taxidsInt := make([]int, 0, 128)
			_taxids := make([]string, 0, 128)
			var taxid uint32
			var taxidInt int
			for _, acc := range accs {
				_taxids = _taxids[:0]
				_taxidsInt = _taxidsInt[:0]
				for taxid = range *acc2taxid[acc] {
					_taxidsInt = append(_taxidsInt, int(taxid))
				}
				sort.Ints(_taxidsInt)

				for _, taxidInt = range _taxidsInt {
					_taxids = append(_taxids, strconv.Itoa(taxidInt))
				}
				fmt.Fprintf(outfhAcc2Taxid, "%s\t%s\n", acc, strings.Join(_taxids, ","))
			}

			log.Infof("%d records saved to %s", len(acc2taxid), fileAcc2Taxid)
		}

		// ------------------------------- nodes.dmp -------------------------

		fileNodes := filepath.Join(outDir, "nodes.dmp")
		outfhNodes, err := xopen.Wopen(fileNodes)
		checkError(err)
		defer outfhNodes.Close()

		taxids := make([]uint32, 0, len(tree))
		for child := range names {
			taxids = append(taxids, child)
		}
		sort.Slice(taxids, func(i, j int) bool {
			return taxids[i] < taxids[j]
		})

		fmt.Fprintf(outfhNodes, "%d\t|\t%d\t|\t%s\t|\t%s\t|\t8\t|\t0\t|\t1\t|\t0\t|\t0\t|\t0\t|\t0\t|\t0\t|\t\t|\n", 1, 1, "no rank", "")
		for _, child := range taxids {
			fmt.Fprintf(outfhNodes, "%d\t|\t%d\t|\t%s\t|\t%s\t|\t0\t|\t1\t|\t11\t|\t1\t|\t0\t|\t1\t|\t1\t|\t0\t|\t\t|\n", child, tree[child], rankNames[ranks[child]], "XX")
		}
		log.Infof("%d records saved to %s", len(tree)+1, fileNodes)

		// ------------------------------- names.dmp -------------------------

		fileNames := filepath.Join(outDir, "names.dmp")
		outfhNames, err := xopen.Wopen(fileNames)
		checkError(err)
		defer outfhNames.Close()

		fmt.Fprintf(outfhNames, "%d\t|\t%s\t|\t\t|\tscientific name\t|\n", 1, "root")
		for _, child := range taxids {
			fmt.Fprintf(outfhNames, "%d\t|\t%s\t|\t\t|\tscientific name\t|\n", child, names0[child])
		}
		log.Infof("%d records saved to %s", len(names)+1, fileNames)

		// ------------------------------- merged.dmp -------------------------

		fileMerged := filepath.Join(outDir, "merged.dmp")
		outfhMerged, err := xopen.Wopen(fileMerged)
		checkError(err)
		defer outfhMerged.Close()

		var merged map[uint32]uint32
		var ok bool

		// ------------------------------- delnodes.dmp -------------------------

		fileDelNodes := filepath.Join(outDir, "delnodes.dmp")
		outfhDelNodes, err := xopen.Wopen(fileDelNodes)
		checkError(err)
		defer outfhDelNodes.Close()

		var delnodes map[uint32]interface{}

		if taxdb != nil {
			// --------------------- newly merged --------------------
			merged = make(map[uint32]uint32, len(taxdb.MergeNodes))
			var _parent uint32
			for child, parent := range tree {
				if _parent, ok = taxdb.Nodes[child]; ok { // not new taxid
					if parent != _parent && // its parent changed
						tree[parent] == taxdb.Nodes[_parent] { // while parents of the parents not changed
						if _, ok = tree[_parent]; !ok { // and the old parent disappeared
							// then the old parent is merged into the new parent

							// For example, phylum 'Desulfobacteraeota' was merged into 'Desulfobacterota':
							// R80
							// GB_GCA_002299865.1	d__Bacteria;p__Desulfobacteraeota_A;c__Desulfovibrionia;o__Desulfovibrionales;f__Desulfonatronaceae;g__UBA663;s__
							// R83
							// GB_GCA_002299865.1	d__Bacteria;p__Desulfobacterota;c__Desulfovibrionia;o__Desulfovibrionales;f__Desulfonatronaceae;g__Desulfonatronum;s__

							// multiple nodes may merged to one node.
							merged[_parent] = parent
						}
					}
				}
			}

			// --------------------- newly deleted --------------------

			delnodes = make(map[uint32]interface{}, len(taxdb.DelNodes))

			for child := range taxdb.Nodes {
				if child == 1 {
					continue
				}

				if _, ok = tree[child]; ok { // still there
					continue
				}

				if _, ok = merged[child]; ok { // merged
					continue
				}

				delnodes[child] = struct{}{}
			}

			// --------------------- append old merged --------------------

			newMerged := make(map[uint32]uint32, 256)
			toDeleteMerged := make(map[uint32]interface{}, 256)
			toDelete := make(map[uint32]interface{}, 256)

			// append old merged.dmp
			var toNew uint32
			for from, to := range taxdb.MergeNodes {
				// previoulsly merged reads may be reused again, it happens both in GTDB and NCBI Taxonomy
				if _, ok = tree[from]; ok {
					// discard the old record
					continue
				}

				// need to check the new taxid

				// the new taxid is merged
				if toNew, ok = merged[to]; ok {
					if toNew == from { // changed back to the old taxid
						// https://github.com/shenwei356/gtdb-taxdump/issues/2#issuecomment-1226186877
						// https://gtdb.ecogenomic.org/genome?gid=GCF_001405015.1
						// The history of GCF_001405015.1 showed Clostridium disporicum was renamed to
						// Clostridium disporicum_A in R95, and changed back in R207.

						// discard the old record
						continue
					} else { // merged to a new taxid
						// https://github.com/shenwei356/gtdb-taxdump/issues/2#issuecomment-1226728018
						// detect chaining merging:
						// previous: A -> B
						// current : B -> C
						// how     : change A -> C, and keep B -> C
						if _, ok = newMerged[from]; !ok {
							newMerged[from] = toNew
						}

						continue
					}
				}

				// the new taxid is deleted
				if _, ok = delnodes[to]; ok {
					// if the new taxid has been deleted, mark the old taxid too
					toDelete[from] = struct{}{}
					continue
				}

				// the old taxid is merged to a new taxid, e.g.,
				// gtdb-taxdump/R089/nodes.dmp: from: 1996042274, oldto: 1424454668, newto: 1975216569
				// previous: A -> B
				// current : A -> C
				// how     : mark B as deleted
				if toNew, ok = merged[from]; ok {
					if to != toNew {
						toDelete[to] = struct{}{}
					}
					continue
				}

				// append
				// merged[from] = to
				if _, ok = tree[to]; !ok {
					newMerged[from] = to
				}
			}

			// apply changes
			for child := range toDeleteMerged {
				delete(merged, child)
			}
			for from, to := range newMerged {
				merged[from] = to
			}
			for child := range toDelete {
				if _, ok = tree[child]; ok {
					continue
				}
				if _, ok = merged[child]; ok {
					continue
				}

				delnodes[child] = struct{}{}
			}

			// --------------------- append old delnodes.dmp ---------------------
			for child := range taxdb.DelNodes {
				if _, ok = tree[child]; ok { // some deleted taxids may be reused
					continue
				}

				if _, ok = merged[child]; ok {
					continue
				}

				delnodes[child] = struct{}{}
			}

			// --------------------------------- write -----------------------------------

			// -------------- write delnodes.dmp ------------

			taxids := make([]uint32, 0, len(delnodes))
			for child := range delnodes {
				taxids = append(taxids, child)
			}
			sort.Slice(taxids, func(i, j int) bool {
				return taxids[i] > taxids[j]
			})

			for _, child := range taxids {
				fmt.Fprintf(outfhDelNodes, "%d\t|\n", child)
			}

			// -------------- write merged.dmp ------------

			taxids = taxids[:0]
			for child := range merged {
				taxids = append(taxids, child)
			}
			sort.Slice(taxids, func(i, j int) bool {
				return taxids[i] < taxids[j]
			})

			for _, child := range taxids {
				fmt.Fprintf(outfhMerged, "%d\t|\t%d\t|\n", child, merged[child])
			}
		}

		log.Infof("%d records saved to %s", len(merged), fileMerged)
		log.Infof("%d records saved to %s", len(delnodes), fileDelNodes)
	},
}

func init() {
	RootCmd.AddCommand(createTaxDumpCmd)

	createTaxDumpCmd.Flags().IntP("field-accession", "A", 0, "field index of assembly accession (genome ID), for outputting taxid.map")
	createTaxDumpCmd.Flags().StringP("field-accession-re", "", `^(.+)$`, `regular expression to extract assembly accession`)
	createTaxDumpCmd.Flags().BoolP("field-accession-as-subspecies", "S", false, "treate the accession as subspecies rank")
	// -------------------------------------------------------------------

	createTaxDumpCmd.Flags().BoolP("gtdb", "", false, "input files are GTDB taxonomy file")
	createTaxDumpCmd.Flags().StringP("gtdb-re-subs", "", `^\w\w_GC[AF]_(.+)\.\d+$`, `regular expression to extract assembly accession as the subspecies`)

	// --------------

	createTaxDumpCmd.Flags().StringSliceP("null", "", []string{"", "NULL", "NA"}, "null value of taxa")
	createTaxDumpCmd.Flags().StringSliceP("rank-names", "R", []string{}, "names of all ranks, leave it empty to use the (lowercase) first row of input as rank names")

	// --------------

	createTaxDumpCmd.Flags().StringP("out-dir", "O", "", `output directory`)
	createTaxDumpCmd.Flags().BoolP("force", "", false, `overwrite existing output directory`)

	// --------------

	createTaxDumpCmd.Flags().IntP("line-chunk-size", "", 5000, `number of lines to process for each thread, and 4 threads is fast enough.`)

	// --------------
	createTaxDumpCmd.Flags().StringP("old-taxdump-dir", "x", "", `taxdump directory of the previous version, for generating merged.dmp and delnodes.dmp`)

}

type _Taxon struct {
	Accession string

	Names  []string
	TaxIds []uint32
}

func (t _Taxon) String() string {
	vs := make([]string, len(t.TaxIds))
	for i, v := range t.TaxIds {
		vs[i] = strconv.Itoa(int(v))
	}
	return fmt.Sprintf("%s, %s, %s", t.Accession, strings.Join(t.Names, ";"), strings.Join(vs, ";"))
}
