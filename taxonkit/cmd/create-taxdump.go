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
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/shenwei356/bio/taxdump"
	"github.com/shenwei356/breader"
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
  0. For GTDB taxonomy file, just use --gtdb
  1. The input file should be tab-delimited
  2. At least one column is needed, please specify the filed index:
     1) Kingdom/Superkingdom/Domain,     -K/--field-kingdom
     2) Phylum,                          -P/--field-phylum
     3) Class,                           -C/--field-class
     4) Order,                           -O/--field-order
     5) Family,                          -F/--field-family
     6) Genus,                           -G/--field-genus
     7) Species (needed),                -S/--field-species
     8) Subspecies,                      -T/--field-subspecies
        For GTDB, we use the numeric assembly accession
        (without the prefix GCA_ and GCF_, and version number).
  3. The column containing the genome/assembly accession is recommended to
     generate TaxId mapping file (taxid.map, id -> taxid).
     -A/--field-accession,    field contaning genome/assembly accession        
     --field-accession-re,    regular expression to extract the accession 

Attentions:
  1. Names should be distinct in taxa of different rank.
     But for these missing some taxon nodes, using names of parent nodes is allowed:

       GB_GCA_018897955.1      d__Archaea;p__EX4484-52;c__EX4484-52;o__EX4484-52;f__LFW-46;g__LFW-46;s__LFW-46 sp018897155

     It can also detect duplicate names with different ranks, e.g.,
     the Class and Genus have the same name B47-G6, and the Order and Family between them have different names.
     In this case, we reassign a new TaxId by increasing the TaxId until it being distinct.

       GB_GCA_003663585.1      d__Archaea;p__Thermoplasmatota;c__B47-G6;o__B47-G6B;f__47-G6;g__B47-G6;s__B47-G6 sp003663585

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
		fKingdom := getFlagNonNegativeInt(cmd, "field-kingdom")
		fPhylum := getFlagNonNegativeInt(cmd, "field-phylum")
		fClass := getFlagNonNegativeInt(cmd, "field-class")
		fOrder := getFlagNonNegativeInt(cmd, "field-order")
		fFamily := getFlagNonNegativeInt(cmd, "field-family")
		fGenus := getFlagNonNegativeInt(cmd, "field-genus")
		fSpecies := getFlagNonNegativeInt(cmd, "field-species")
		fSubspe := getFlagNonNegativeInt(cmd, "field-subspecies")

		rankNames := getFlagStringSlice(cmd, "rank-names")
		if len(rankNames) != 8 {
			checkError(fmt.Errorf(`the number of --rank-names should be 8`))
		}

		var err error

		reGenomeIDStr := getFlagString(cmd, "field-accession-re")

		var reGenomeID *regexp.Regexp
		if reGenomeIDStr != "" {
			if !regexp.MustCompile(`\(.+\)`).MatchString(reGenomeIDStr) {
				checkError(fmt.Errorf(`value of --field-accession-re must contains "(" and ")"`))
			}

			reGenomeID, err = regexp.Compile(reGenomeIDStr)
			if err != nil {
				checkError(fmt.Errorf("fail to compile regular expression: %s", reGenomeIDStr))
			}
		}

		isGTDB := getFlagBool(cmd, "gtdb")
		reGTDBStr := getFlagString(cmd, "gtdb-re-subs")

		var reGTDBsubspe *regexp.Regexp

		if reGTDBStr != "" {
			if !regexp.MustCompile(`\(.+\)`).MatchString(reGTDBStr) {
				checkError(fmt.Errorf(`value of --gtdb-re-subs must contains "(" and ")"`))
			}

			reGTDBsubspe, err = regexp.Compile(reGTDBStr)
			if err != nil {
				checkError(fmt.Errorf("fail to compile regular expression: %s", reGTDBStr))
			}
		}

		nulls := getFlagStringSlice(cmd, "null")

		var hasAccession bool
		var numFields int
		if isGTDB {
			numFields = 2
			hasAccession = true
		} else if fSpecies == 0 {
			checkError(fmt.Errorf("flag -S/--field-species needed"))
		} else {
			numFields = MaxInts(fAccession, fSubspe, fKingdom, fPhylum, fClass, fOrder, fFamily, fGenus, fSpecies)
			hasAccession = fAccession > 0
		}

		outDir := getFlagString(cmd, "out-dir")
		force := getFlagBool(cmd, "force")
		if outDir == "" {
			checkError(fmt.Errorf("flag --out-dir is needed"))
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

		hasKingdom := fKingdom > 0
		hasPhylum := fPhylum > 0
		hasClass := fClass > 0
		hasOrder := fOrder > 0
		hasFamily := fFamily > 0
		hasGenus := fGenus > 0
		hasSubspe := fSubspe > 0

		nullMap := make(map[string]interface{})
		for _, k := range nulls {
			nullMap[k] = struct{}{}
		}

		files := getFileList(args)

		pool := &sync.Pool{New: func() interface{} {
			tmp := make([]string, numFields)
			return &tmp
		}}
		pool7 := &sync.Pool{New: func() interface{} {
			tmp := make([]string, 7)
			return &tmp
		}}

		var reGTDBsubspeNotCaptured bool
		var reGTDBsubspeNotCapturedExample string

		var reGenomeIDNotCaptured bool
		var reGenomeIDNotCapturedExample string

		fn := func(line string) (interface{}, bool, error) {
			line = strings.Trim(line, "\r\n")
			if line == "" {
				return nil, false, nil
			}

			items := pool.Get().(*[]string)
			defer pool.Put(items)

			stringSplitNByByte(line, '\t', numFields, items)
			if len(*items) < numFields {
				return nil, false, nil
			}

			t := _Taxon{}
			var val string

			if isGTDB {
				if reGTDBsubspe != nil {
					found := reGTDBsubspe.FindAllStringSubmatch((*items)[0], 1)
					if len(found) == 0 {
						t.Subspe = (*items)[0]
						// checkError(fmt.Errorf("invalid GTDB assembly accession: %s", (*items)[0]))
						reGTDBsubspeNotCaptured = true
						reGTDBsubspeNotCapturedExample = (*items)[0]
					} else {
						t.Subspe = found[0][1]
					}
				}

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

				items7 := pool7.Get().(*[]string)
				defer pool7.Put(items7)

				stringSplitNByByte(CopyString((*items)[1]), ';', 7, items7)
				if len(*items7) < 7 {
					checkError(fmt.Errorf("invalid GTDB taxonomy record: %s", line))
				}

				val = (*items7)[0]
				if len(val) < 3 || val[0:3] != "d__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (kingdom): %s", val))
				}
				t.Kingdom = val[3:]

				val = (*items7)[1]
				if len(val) < 3 || val[0:3] != "p__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (phylum): %s", val))
				}
				t.Phylum = val[3:]

				val = (*items7)[2]
				if len(val) < 3 || val[0:3] != "c__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (class): %s", val))
				}
				t.Class = val[3:]

				val = (*items7)[3]
				if len(val) < 3 || val[0:3] != "o__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (order): %s", val))
				}
				t.Order = val[3:]

				val = (*items7)[4]
				if len(val) < 3 || val[0:3] != "f__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (family): %s", val))
				}
				t.Family = val[3:]

				val = (*items7)[5]
				if len(val) < 3 || val[0:3] != "g__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (genus): %s", val))
				}
				t.Genus = val[3:]

				val = (*items7)[6]
				if len(val) < 3 || val[0:3] != "s__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (species): %s", val))
				}
				t.Species = val[3:]

				// ----------------------------------------
				t.Names[0] = t.Kingdom
				t.TaxIds[0] = uint32(xxhash.Sum64String(strings.ToLower(t.Kingdom)))

				if t.Phylum != t.Kingdom {
					t.Names[1] = t.Phylum
					t.TaxIds[1] = uint32(xxhash.Sum64String(strings.ToLower(t.Phylum)))
				}

				if t.Class != t.Phylum {
					t.Names[2] = t.Class
					t.TaxIds[2] = uint32(xxhash.Sum64String(strings.ToLower(t.Class)))
				}

				if t.Order != t.Class {
					t.Names[3] = t.Order
					t.TaxIds[3] = uint32(xxhash.Sum64String(strings.ToLower(t.Order)))
				}

				if t.Family != t.Order {
					t.Names[4] = t.Family
					t.TaxIds[4] = uint32(xxhash.Sum64String(strings.ToLower(t.Family)))
				}
				if t.Genus != t.Family {
					t.Names[5] = t.Genus
					t.TaxIds[5] = uint32(xxhash.Sum64String(strings.ToLower(t.Genus)))
				}
				if t.Species != t.Genus {
					t.Names[6] = t.Species
					t.TaxIds[6] = uint32(xxhash.Sum64String(strings.ToLower(t.Species)))
				}

				t.Names[7] = t.Subspe
				t.TaxIds[7] = uint32(xxhash.Sum64String(strings.ToLower(t.Subspe)))

				return &t, true, nil
			}

			var ok bool

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

			if hasKingdom {
				val = (*items)[fKingdom-1]
				if _, ok = nullMap[val]; !ok {
					t.Kingdom = val

					t.Names[0] = t.Kingdom
					t.TaxIds[0] = uint32(xxhash.Sum64String(strings.ToLower(t.Kingdom)))
				}
			}

			if hasPhylum {
				val = (*items)[fPhylum-1]
				if _, ok = nullMap[val]; !ok {
					t.Phylum = val

					t.Names[1] = t.Phylum
					t.TaxIds[1] = uint32(xxhash.Sum64String(strings.ToLower(t.Phylum)))
				}
			}

			if hasClass {
				val = (*items)[fClass-1]
				if _, ok = nullMap[val]; !ok {
					t.Class = val

					t.Names[3] = t.Order
					t.TaxIds[3] = uint32(xxhash.Sum64String(strings.ToLower(t.Order)))
				}
			}

			if hasOrder {
				val = (*items)[fOrder-1]
				if _, ok = nullMap[val]; !ok {
					t.Order = val

					t.Names[3] = t.Order
					t.TaxIds[3] = uint32(xxhash.Sum64String(strings.ToLower(t.Order)))
				}
			}

			if hasFamily {
				val = (*items)[fFamily-1]
				if _, ok = nullMap[val]; !ok {
					t.Family = val

					t.Names[4] = t.Family
					t.TaxIds[4] = uint32(xxhash.Sum64String(strings.ToLower(t.Family)))
				}
			}

			if hasGenus {
				val = (*items)[fGenus-1]
				if _, ok = nullMap[val]; !ok {
					t.Genus = val

					t.Names[5] = t.Genus
					t.TaxIds[5] = uint32(xxhash.Sum64String(strings.ToLower(t.Genus)))
				}
			}

			val = (*items)[fSpecies-1]
			if _, ok = nullMap[val]; !ok {
				t.Species = val

				t.Names[6] = t.Species
				t.TaxIds[6] = uint32(xxhash.Sum64String(strings.ToLower(t.Species)))
			}

			if hasSubspe {
				val = (*items)[fSubspe-1]
				t.Subspe = val

				t.Names[7] = t.Subspe
				t.TaxIds[7] = uint32(xxhash.Sum64String(strings.ToLower(t.Subspe)))
			}

			return &t, true, nil
		}

		// child -> parent
		tree := make(map[uint32]uint32, 1<<16)

		// child -> rank-id
		ranks := make(map[uint32]uint8, 1<<16)

		// child -> name
		names := make(map[uint32]string, 1<<16)

		// accession -> taxid
		acc2taxid := make(map[string]uint32, 1<<16)
		accIdx := make(map[string]int, 1<<16)
		var idx int

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var t *_Taxon
			var data interface{}
			var i int
			var prev int
			var first bool
			var taxid uint32
			var _name string
			var _rank uint8
			var ok bool
			var reAssignTaxid bool

			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data = range chunk.Data {
					t = data.(*_Taxon)
					// fmt.Println((*t).String2())
					// fmt.Println(*t)

					first = true
					for i = 7; i >= 0; i-- {
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
							if _name != t.Names[i] { // two names hashed to the same uint32
								if config.Verbose {
									log.Infof(`"%s" and "%s" having the same taxId: %d`, _name, t.Names[i], taxid)
								}
								reAssignTaxid = true
							}
						} else {
							names[taxid] = t.Names[i]
						}

						if _rank, ok = ranks[taxid]; ok {
							if int(_rank) != i {
								if config.Verbose {
									log.Debug(*t)
									log.Infof(`duplicate name (%s) with different ranks: "%s" and "%s"`, t.Names[i], rankNames[_rank], rankNames[i])
								}
								reAssignTaxid = true
							}
						} else {
							ranks[taxid] = uint8(i)
						}

						if reAssignTaxid {
							if config.Verbose {
								log.Infof(`assign a new TaxId for "%s" (rank: %s): %d -> %d`, names[taxid], rankNames[i], taxid, taxid+1)
							}
							taxid++
							t.TaxIds[i] = taxid
							goto REASSIGNTAXID
						}

						if first {
							if hasAccession {
								idx++
								accIdx[t.Accession] = idx
								acc2taxid[t.Accession] = taxid // the lowest taxid
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
			}

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

			for _, acc := range accs {
				fmt.Fprintf(outfhAcc2Taxid, "%s\t%d\n", acc, acc2taxid[acc])
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
			fmt.Fprintf(outfhNames, "%d\t|\t%s\t|\t\t|\tscientific name\t|\n", child, names[child])
		}
		log.Infof("%d records saved to %s", len(names)+1, fileNames)

		// ------------------------------- merged.dmp -------------------------

		fileMerged := filepath.Join(outDir, "merged.dmp")
		outfhMerged, err := xopen.Wopen(fileMerged)
		checkError(err)
		defer outfhMerged.Close()

		var merged map[uint32]uint32
		var ok bool

		if taxdb != nil {
			merged = make(map[uint32]uint32, len(taxdb.MergeNodes))
			var _parent uint32
			for child, parent := range tree {
				if _parent, ok = taxdb.Nodes[child]; ok { // not new taxid
					if parent != _parent && // its parent changed
						tree[parent] == taxdb.Nodes[_parent] { // while parents of the parents not changed
						if _, ok = tree[_parent]; !ok { // and the old parent disappeared

							// R80
							// GB_GCA_002299865.1	d__Bacteria;p__Desulfobacteraeota_A;c__Desulfovibrionia;o__Desulfovibrionales;f__Desulfonatronaceae;g__UBA663;s__
							// R83
							// GB_GCA_002299865.1	d__Bacteria;p__Desulfobacterota;c__Desulfovibrionia;o__Desulfovibrionales;f__Desulfonatronaceae;g__Desulfonatronum;s__

							merged[_parent] = parent
						}
					}
				}
			}

			// append old delnodes.dmp
			for from, to := range taxdb.MergeNodes {
				if _, ok = merged[from]; !ok {
					merged[from] = to
				}
			}

			taxids := make([]uint32, 0, len(merged))
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

		// ------------------------------- delnodes.dmp -------------------------

		fileDelNodes := filepath.Join(outDir, "delnodes.dmp")
		outfhDelNodes, err := xopen.Wopen(fileDelNodes)
		checkError(err)
		defer outfhDelNodes.Close()

		var delnodes []uint32

		if taxdb != nil {
			delnodes = make([]uint32, 0, len(taxdb.DelNodes))

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

				delnodes = append(delnodes, child)
			}

			// append old delnodes.dmp
			for child := range taxdb.DelNodes {
				if _, ok = tree[child]; !ok { // some deleted taxids may be reused
					delnodes = append(delnodes, child)
				}
			}

			sort.Slice(delnodes, func(i, j int) bool {
				return delnodes[i] > delnodes[j]
			})

			for _, child := range delnodes {
				fmt.Fprintf(outfhDelNodes, "%d\t|\n", child)
			}
		}

		log.Infof("%d records saved to %s", len(delnodes), fileDelNodes)
	},
}

func init() {
	RootCmd.AddCommand(createTaxDumpCmd)

	createTaxDumpCmd.Flags().IntP("field-kingdom", "K", 0, "field index of kingdom")
	createTaxDumpCmd.Flags().IntP("field-phylum", "P", 0, "field index of phylum")
	createTaxDumpCmd.Flags().IntP("field-class", "C", 0, "field index of class")
	createTaxDumpCmd.Flags().IntP("field-order", "O", 0, "field index of order")
	createTaxDumpCmd.Flags().IntP("field-family", "F", 0, "field index of family")
	createTaxDumpCmd.Flags().IntP("field-genus", "G", 0, "field index of genus")
	createTaxDumpCmd.Flags().IntP("field-species", "S", 0, "field index of species (needed)")
	createTaxDumpCmd.Flags().IntP("field-subspecies", "T", 0, "field index of subspecies")

	createTaxDumpCmd.Flags().IntP("field-accession", "A", 0, "field index of assembly accession (genome ID), for outputting taxid.map")
	createTaxDumpCmd.Flags().StringP("field-accession-re", "", `^\w\w_(.+)$`, `regular expression to extract assembly accession`)

	// -------------------------------------------------------------------

	createTaxDumpCmd.Flags().BoolP("gtdb", "", false, "input files are GTDB taxonomy file")
	createTaxDumpCmd.Flags().StringP("gtdb-re-subs", "", `^\w\w_GC[AF]_(.+)\.\d+$`, `regular expression to extract assembly accession as the subspecies`)

	// --------------

	createTaxDumpCmd.Flags().StringSliceP("null", "", []string{"", "NULL", "NA"}, "null value of taxa")
	createTaxDumpCmd.Flags().StringSliceP("rank-names", "", []string{"superkingdom", "phylum", "class", "order", "family", "genus", "species", "no rank"}, "names of the 8 ranks, the order maters")

	// --------------

	createTaxDumpCmd.Flags().StringP("out-dir", "", "", `output directory`)
	createTaxDumpCmd.Flags().BoolP("force", "", false, `overwrite existed output directory`)

	// --------------

	createTaxDumpCmd.Flags().IntP("line-chunk-size", "", 5000, `number of lines to process for each thread, and 4 threads is fast enough.`)

	// --------------
	createTaxDumpCmd.Flags().StringP("old-taxdump-dir", "x", "", `taxdump directory of the previous version, for generating merged.dmp and delnodes.dmp`)

}

type _Taxon struct {
	Accession string

	Subspe  string
	Kingdom string
	Phylum  string
	Class   string
	Order   string
	Family  string
	Genus   string
	Species string

	Names  [8]string // Kingdom, phylum ...
	TaxIds [8]uint32 // Kingdom, phylum ...
}

func (t _Taxon) String() string {
	return fmt.Sprintf("%s, %v, %v", t.Subspe, t.Names, t.TaxIds)
}

func (t _Taxon) String2() string {
	return fmt.Sprintf("%s, k: %s, p: %s, c: %s, o: %s, f: %s, g: %s, s: %s",
		t.Subspe, t.Kingdom, t.Phylum, t.Class, t.Order, t.Family, t.Genus, t.Species)
}
