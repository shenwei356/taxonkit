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
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/shenwei356/breader"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// createTaxDumpCmd represents the create-TaxDump command
var createTaxDumpCmd = &cobra.Command{
	Use:   "create-taxdump",
	Short: "Create NCBI-style taxdump files for custom taxonomy, e.g., GTDB",
	Long: `Create NCBI-style taxdump files for custom taxonomy, e.g., GTDB

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
        For GTDB, we use the assembly accession (without version number).

Attentions:
  1. Names should be distinct in taxa of different rank.
     But for these missing some taxon nodes, using names of parent nodes is allowed:

       GB_GCA_018897955.1      d__Archaea;p__EX4484-52;c__EX4484-52;o__EX4484-52;f__LFW-46;g__LFW-46;s__LFW-46 sp018897155

     It can also detect duplicate names with different ranks, e.g.,
     The Class and Genus have the same name B47-G6, and the Order and Family between them have different names.
     In this case, we reassign TaxId by increasing the TaxId until it being distinct.

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

		isGTDB := getFlagBool(cmd, "gtdb")
		reGTDBStr := getFlagString(cmd, "gtdb-re-subs")

		if !regexp.MustCompile(`\(.+\)`).MatchString(reGTDBStr) {
			checkError(fmt.Errorf(`value of --gtdb-re-subs must contains "(" and ")"`))
		}

		reGTDBsubspe, err := regexp.Compile(reGTDBStr)
		if err != nil {
			checkError(fmt.Errorf("fail to compile regular expression: %s", reGTDBStr))
		}

		nulls := getFlagStringSlice(cmd, "null")

		var numFields int
		if isGTDB {
			numFields = 2
		} else if fSubspe == 0 || fSpecies == 0 {
			checkError(fmt.Errorf("flags -A/--field-accession and -S/--field-species needed"))
		} else {
			numFields = MaxInts(fSubspe, fKingdom, fPhylum, fClass, fOrder, fFamily, fGenus, fSpecies)
		}

		outDir := getFlagString(cmd, "out-dir")
		force := getFlagBool(cmd, "force")
		if outDir == "" {
			checkError(fmt.Errorf("flag -O/--out-dir is needed"))
		}

		makeOutDir(outDir, force)

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
				found := reGTDBsubspe.FindAllStringSubmatch((*items)[0], 1)
				if len(found) == 0 {
					checkError(fmt.Errorf("invalid GTDB assembly accession: %s", (*items)[0]))
				}
				t.Subspe = found[0][1]

				items7 := pool7.Get().(*[]string)
				defer pool7.Put(items7)

				stringSplitNByByte(CopyString((*items)[1]), ';', 7, items7)
				if len(*items7) < 7 {
					checkError(fmt.Errorf("invalid GTDB taxonomy record: %s", line))
				}

				val = (*items7)[0]
				if len(val) < 4 || val[0:3] != "d__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (kingdom): %s", val))
				}
				t.Kingdom = val[3:]

				val = (*items7)[1]
				if len(val) < 4 || val[0:3] != "p__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (phylum): %s", val))
				}
				t.Phylum = val[3:]

				val = (*items7)[2]
				if len(val) < 4 || val[0:3] != "c__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (class): %s", val))
				}
				t.Class = val[3:]

				val = (*items7)[3]
				if len(val) < 4 || val[0:3] != "o__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (order): %s", val))
				}
				t.Order = val[3:]

				val = (*items7)[4]
				if len(val) < 4 || val[0:3] != "f__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (family): %s", val))
				}
				t.Family = val[3:]

				val = (*items7)[5]
				if len(val) < 4 || val[0:3] != "g__" {
					checkError(fmt.Errorf("invalid GTDB taxonomy format (genus): %s", val))
				}
				t.Genus = val[3:]

				val = (*items7)[6]
				if len(val) < 4 || val[0:3] != "s__" {
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
						if taxid == 0 {
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

		// ------------------------------- nodes.dmp -------------------------

		fileNodes := filepath.Join(outDir, "nodes.dmp")
		outfhNodes, err := xopen.Wopen(fileNodes)
		checkError(err)
		defer outfhNodes.Close()

		fmt.Fprintf(outfhNodes, "%d\t|\t%d\t|\t%s\t|\t%s\t|\t8\t|\t0\t|\t1\t|\t0\t|\t0\t|\t0\t|\t0\t|\t0\t|\t\t|\n", 1, 1, "no rank", "")
		for child, parent := range tree {
			fmt.Fprintf(outfhNodes, "%d\t|\t%d\t|\t%s\t|\t%s\t|\t0\t|\t1\t|\t11\t|\t1\t|\t0\t|\t1\t|\t1\t|\t0\t|\t\t|\n", child, parent, rankNames[ranks[child]], "XX")
		}
		log.Infof("%d records saved to %s", len(tree)+1, fileNodes)

		// ------------------------------- names.dmp -------------------------

		fileNames := filepath.Join(outDir, "names.dmp")
		outfhNames, err := xopen.Wopen(fileNames)
		checkError(err)
		defer outfhNames.Close()

		fmt.Fprintf(outfhNames, "%d\t|\t%s\t|\t\t|\tscientific name\t|\n", 1, "root")
		for child, name := range names {
			fmt.Fprintf(outfhNames, "%d\t|\t%s\t|\t\t|\tscientific name\t|\n", child, name)
		}
		log.Infof("%d records saved to %s", len(names)+1, fileNames)

		// ------------------------------- merged.dmp -------------------------
		// TODO: compare to a old version

		fileMerged := filepath.Join(outDir, "merged.dmp")
		outfhMerged, err := xopen.Wopen(fileMerged)
		checkError(err)
		defer outfhMerged.Close()

		// fmt.Fprintf(outfhMerged, "\n")

		log.Infof("%d records saved to %s", 0, fileMerged)

		// ------------------------------- delnodes.dmp -------------------------
		// TODO: compare to a old version

		fileDelNodes := filepath.Join(outDir, "delnodes.dmp")
		outfhDelNodes, err := xopen.Wopen(fileDelNodes)
		checkError(err)
		defer outfhDelNodes.Close()

		// fmt.Fprintf(outfhDelNodes, "\n")

		log.Infof("%d records saved to %s", 0, fileDelNodes)
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

	// -------------------------------------------------------------------

	createTaxDumpCmd.Flags().BoolP("gtdb", "", false, "input files are GTDB taxonomy file")
	createTaxDumpCmd.Flags().StringP("gtdb-re-subs", "", `^\w\w_(.+)\.\d+$`, `regular expression to extract accession as the subspecies from the assembly ID`)

	createTaxDumpCmd.Flags().IntP("line-chunk-size", "", 5000, `number of lines to process for each thread, and 4 threads is fast enough.`)

	// --------------

	createTaxDumpCmd.Flags().StringSliceP("null", "", []string{"", "NULL", "NA"}, "null value of taxa")
	createTaxDumpCmd.Flags().StringSliceP("rank-names", "", []string{"superkingdom", "phylum", "class", "order", "family", "genus", "species", "no rank"}, "names of the 8 ranks, order maters")

	// --------------

	createTaxDumpCmd.Flags().StringP("out-dir", "", "", `output directory`)
	createTaxDumpCmd.Flags().BoolP("force", "", false, `overwrite existed output directory`)
}

type _Taxon struct {
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