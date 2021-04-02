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
	"strconv"
	"strings"
	"sync"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// flineageCmd represents the fx2tab command
var flineageCmd = &cobra.Command{
	Use:   "reformat",
	Short: "Reformat lineage in canonical ranks",
	Long: `Reformat lineage in canonical ranks

Output format can be formated by flag --format, available placeholders:

    {k}: superkingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {t}: subspecies/strain
    
    {S}: subspecies
    {T}: strain

When these's no nodes of rank "subspecies" nor "stain",
you can switch -S/--pseudo-strain to use the node with lowest rank
as subspecies/strain name, if which rank is lower than "species". 
This flag affects {t}, {S}, {T}.
    
Output format can contains some escape charactors like "\t".

This command appends reformated lineage to the input line.
The corresponding TaxIds of reformated lineage can be provided as another
column by flag "-t/--show-lineage-taxids".

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "miss-rank-repl")
		prefix := getFlagString(cmd, "miss-rank-repl-prefix")
		iblank := getFlagString(cmd, "miss-taxid-repl")
		fill := getFlagBool(cmd, "fill-miss-rank")
		pseudoStrain := getFlagBool(cmd, "pseudo-strain")

		taxIdField := getFlagNonNegativeInt(cmd, "taxid-field")
		field := getFlagPositiveInt(cmd, "lineage-field")

		var parsingTaxId bool
		if taxIdField > 0 {
			log.Infof("parsing TaxIds from field %d", taxIdField)
			parsingTaxId = true
			taxIdField--
		} else if field > 0 {
			log.Infof("parsing complete lineages from field %d", field)
			field--
		}

		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")

		addPrefix := getFlagBool(cmd, "add-prefix")
		prefixK := getFlagString(cmd, "prefix-k")
		prefixP := getFlagString(cmd, "prefix-p")
		prefixC := getFlagString(cmd, "prefix-c")
		prefixO := getFlagString(cmd, "prefix-o")
		prefixF := getFlagString(cmd, "prefix-f")
		prefixG := getFlagString(cmd, "prefix-g")
		prefixs := getFlagString(cmd, "prefix-s")
		prefixS := getFlagString(cmd, "prefix-S")
		prefixt := getFlagString(cmd, "prefix-t")
		prefixT := getFlagString(cmd, "prefix-T")

		trim := getFlagBool(cmd, "trim")

		prefixes := map[string]string{
			"k": prefixK,
			"p": prefixP,
			"c": prefixC,
			"o": prefixO,
			"f": prefixF,
			"g": prefixG,
			"s": prefixs,
			"S": prefixS,
			"t": prefixt,
			"T": prefixT,
		}

		// check format
		if !reRankPlaceHolder.MatchString(format) {
			checkError(fmt.Errorf("placeholder of simplified rank not found in output format: %s", format))
		}
		matches := reRankPlaceHolder.FindAllStringSubmatch(format, -1)
		flag := false
		for _, match := range matches {
			if _, ok := symbol2rank[match[1]]; !ok {
				checkError(fmt.Errorf("invalid placeholder: %s", match[0]))
			}
			switch match[1] {
			case "t", "S", "T":
				flag = true
			}
		}
		if flag {
			if pseudoStrain && !fill {
				fill = true
				log.Infof("-F/--fill-miss-rank is switched on when giving flag -S/--pseudo-strain")
			}
		} else if pseudoStrain {
			log.Warningf(`flag -S/--pseudo-strain will not work because none of "{t}", "{S}", "{T}" is found in -f/--format`)
		}

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)
		defer outfh.Close()

		// --------------------------------------------------------
		// load data

		// for parsing lineage
		var taxid2taxon map[uint32]*Taxon
		var name2parent2taxid map[string]map[string]uint32
		var name2taxid map[string]uint32
		var ambigous map[string][]uint32

		// for parsing taxid
		var tree0 map[uint32]uint32
		var ranks0 map[uint32]string
		var names0 map[uint32]string
		var delnodes0 map[uint32]struct{}
		var merged0 map[uint32]uint32

		if parsingTaxId {
			tree0, ranks0, names0, delnodes0, merged0 = loadData(config, true, true)
		} else {
			taxid2taxon, name2parent2taxid, name2taxid, ambigous = getName2Parent2Taxid(config)
		}

		// --------------------------------------------------------

		type line2flineage struct {
			line      string
			flineage  string
			iflineage string
		}

		unescape := stringutil.UnEscaper()

		var poolStrings = &sync.Pool{New: func() interface{} {
			return make([]string, 0, 32)
		}}

		weightOfSpecies := symbol2weight["s"]

		fn := func(line string) (interface{}, bool, error) {
			if len(line) == 0 || line[0] == '#' {
				return nil, false, nil
			}
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}
			data := strings.Split(line, "\t")

			if parsingTaxId {
				if len(data) < taxIdField+1 {
					return nil, false, fmt.Errorf("taxid-field (%d) out of range (%d):%s", taxIdField+1, len(data), line)
				}
			} else if len(data) < field+1 {
				return nil, false, fmt.Errorf("lineage-field (%d) out of range (%d):%s", field+1, len(data), line)
			}

			// -----------------------------------------------

			var rank, srank string   // lower case of name : name
			var lname, plname string // lower case of name, rank and it's one-letter symbol
			var ok bool

			// preprare replacements.
			// find the orphan names and missing ranks
			replacements := make(map[string]string, len(matches))

			var ireplacements map[string]string
			if printLineageInTaxid {
				ireplacements = make(map[string]string, len(matches))
			}

			for _, match := range matches {
				replacements[match[1]] = blank
				if printLineageInTaxid {
					ireplacements[match[1]] = iblank
				}
			}

			// -----------------------------------------------

			var taxid uint32
			var taxidInt int

			var names []string
			var ranks []string
			var taxids []uint32

			if parsingTaxId {
				taxidInt, err = strconv.Atoi(data[taxIdField])
				if err != nil || taxidInt < 0 {
					checkError(fmt.Errorf("invalid TaxId: %s", data[taxIdField]))
				}
				taxid = uint32(taxidInt)
				names, ranks, taxids, ok = namesRanksTaxids(tree0, ranks0, names0, delnodes0, merged0, taxid)
				if !ok {
					return line2flineage{line, "", ""}, false, nil
				}
			} else {
				// names
				names = strings.Split(data[field], delimiter) // all names of full lineage

				// ranks := make([]string, len(names))
				ranks = poolStrings.Get().([]string)
			}

			// sranks := make([]string, len(names))
			sranks := poolStrings.Get().([]string)

			name2Name := make(map[string]string, len(names)) // lower case of name of parent

			srank2idx := make(map[string]int) // srank: index

			var maxRankWeight float32

			var pair string
			var _ambids []uint32
			for i, name := range names {
				if name == "" {
					continue
				}

				lname = strings.ToLower(name)
				name2Name[lname] = name
				name = lname

				// -----------------------------------------------

				if parsingTaxId {
					rank = ranks[i]
					taxid = taxids[i]
				} else {
					if _, ok = name2taxid[name]; !ok { // unofficial name
						log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, name)
						return line2flineage{line, "", ""}, true, nil
					}

					if i == 0 { // root node
						taxid = name2taxid[name]
					} else {
						plname = strings.ToLower(names[i-1])
						if _, ok = name2parent2taxid[name]; !ok {
							log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, name)
							return line2flineage{line, "", ""}, true, nil
						}
						if taxid, ok = name2parent2taxid[name][plname]; !ok {
							log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, plname)
							return line2flineage{line, "", ""}, true, nil
						}

						// for cases where child-parent pairs are shared by multiple taxids.
						pair = name + "__" + plname

						if _ambids, ok = ambigous[pair]; ok {

							var _lineage string
							lineage0 := strings.Join(names[:i], delimiter)
							var _taxids2 []uint32
							var _taxid uint32

							_taxids2 = make([]uint32, 0, 2) // possible taxids

							for _, _taxid = range _ambids {
								_lineage = lineageFromTaxid2Taxon(taxid2taxon, _taxid, delimiter)
								if _lineage == lineage0 {
									_taxids2 = append(_taxids2, _taxid)
								}
							}
							switch len(_taxids2) { // cool
							case 0:
								log.Warningf("it's a bug, please report: '%s'. %s", pair, line)
							case 1: // we correct it
								taxid = _taxid
							default:
								tmp := make([]string, len(_taxids2))
								for _i, _taxid := range _taxids2 {
									tmp[_i] = strconv.Itoa(int(_taxid))
								}
								log.Warningf("we can't distinguish the TaxId (%s) for lineage: %s",
									strings.Join(tmp, ", "), lineage0)
							}
						}
					}
					// note that code below is computing rank of current name, not its parent.
					rank = taxid2taxon[taxid].Rank

					if rank == norank {
						ranks = append(ranks, rank)
						sranks = append(sranks, "")
						continue
					}

					// ranks[i] = rank
					ranks = append(ranks, rank)
				}

				// -----------------------------------------------

				if srank, ok = rank2symbol[rank]; ok {
					// special symbol "{t}"
					switch rank {
					case "strain":
						replacements["t"] = name2Name[name]
						if printLineageInTaxid {
							ireplacements["t"] = strconv.Itoa(int(taxid))
						}
						srank2idx["t"] = i
					case "subspecies":
						replacements["t"] = name2Name[name]
						if printLineageInTaxid {
							ireplacements["t"] = strconv.Itoa(int(taxid))
						}
						srank2idx["t"] = i
					}

					replacements[srank] = name2Name[name]
					if printLineageInTaxid {
						ireplacements[srank] = strconv.Itoa(int(taxid))
					}
					srank2idx[srank] = i
					// sranks[i] = srank
					sranks = append(sranks, srank)

					if trim && symbol2weight[srank] > maxRankWeight {
						maxRankWeight = symbol2weight[srank]
					}
				} else {
					sranks = append(sranks, "")
				}
			}

			if fill {
				_, hasRankSubspecies := srank2idx["S"]
				_, hasRankStrain := srank2idx["T"]

				var j, lastI int
				var srank2 string
				for _, srank = range srankList {
					if srank == "" {
						continue
					}

					if _, ok = srank2idx[srank]; ok {
						continue
					}

					if trim && symbol2weight[srank] > maxRankWeight {
						continue
					}

					// missing some ranks.
					// find the nearst higher formal rank
					for j, rank = range ranks {
						srank2 = sranks[j]
						if _, ok = srank2idx[srank2]; ok {
							if symbol2weight[srank2] < symbol2weight[srank] {
								lastI = j
							} else {
								break
							}
						}
					}

					if pseudoStrain {
						if symbol2weight[srank] > weightOfSpecies && // lower than species
							!(hasRankSubspecies || hasRankStrain) && // does not have strain or subspecies
							lastI < len(names)-1 { // not itself
							replacements[srank] = names[len(names)-1]
							continue
						}
					}

					replacements[srank] = prefix + names[lastI] + " " + symbol2rank[srank]
					// replacements[srank] = fmt.Sprintf("%s%s %s", prefix, names[lastI], symbol2rank[srank])
				}
			}

			flineage := format
			var iflineage string

			if printLineageInTaxid {
				iflineage = format
			}

			for srank, re := range reRankPlaceHolders {
				if addPrefix {
					flineage = re.ReplaceAllString(flineage, prefixes[srank]+replacements[srank])
				} else {
					flineage = re.ReplaceAllString(flineage, replacements[srank])
				}

				if printLineageInTaxid {
					iflineage = re.ReplaceAllString(iflineage, ireplacements[srank])
				}
			}

			// recycle
			ranks = ranks[:0]
			poolStrings.Put(ranks)
			sranks = sranks[:0]
			poolStrings.Put(sranks)

			if parsingTaxId {
				names = names[:0]
				poolStrings.Put(names)
				taxids = taxids[:0]
				poolUint32.Put(taxids)
			}

			return line2flineage{line, unescape(flineage), unescape(iflineage)}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 64, fn)
			checkError(err)

			var l2s line2flineage
			var data interface{}
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data = range chunk.Data {
					l2s = data.(line2flineage)

					if printLineageInTaxid {
						outfh.WriteString(l2s.line + "\t" + l2s.flineage + "\t" + l2s.iflineage + "\n")
					} else {
						outfh.WriteString(l2s.line + "\t" + l2s.flineage + "\n")
					}
					if config.LineBuffered {
						outfh.Flush()
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(flineageCmd)

	flineageCmd.Flags().StringP("format", "f", "{k};{p};{c};{o};{f};{g};{s}", "output format, placeholders of rank are needed")
	flineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in input lineage")
	flineageCmd.Flags().StringP("miss-rank-repl", "r", "", `replacement string for missing rank`)
	flineageCmd.Flags().StringP("miss-rank-repl-prefix", "p", "unclassified ", `prefix for estimated taxon level`)
	flineageCmd.Flags().StringP("miss-taxid-repl", "R", "", `replacement string for missing taxid`)

	flineageCmd.Flags().BoolP("fill-miss-rank", "F", false, "fill missing rank with lineage information of the next higher rank")
	flineageCmd.Flags().BoolP("pseudo-strain", "S", false, `use the node with lowest rank as strain name, only if which rank is lower than "species" and not "subpecies" nor "strain". It affects {t}, {S}, {T}. This flag needs flag -F`)

	flineageCmd.Flags().IntP("lineage-field", "i", 2, "field index of lineage. data should be tab-separated")
	flineageCmd.Flags().IntP("taxid-field", "I", 0, "field index of taxid. input data should be tab-separated")
	flineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `show corresponding taxids of reformated lineage`)

	flineageCmd.Flags().BoolP("add-prefix", "P", false, `add prefixes for all ranks, single prefix for a rank is defined by flag --prefix-X`)
	flineageCmd.Flags().StringP("prefix-k", "", "k__", `prefix for superkingdom, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-p", "", "p__", `prefix for phylum, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-c", "", "c__", `prefix for class, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-o", "", "o__", `prefix for order, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-f", "", "f__", `prefix for family, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-g", "", "g__", `prefix for genus, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-s", "", "s__", `prefix for species, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-t", "", "t__", `prefix for subspecies/strain, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-S", "", "S__", `prefix for subspecies, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-T", "", "T__", `prefix for strain, used along with flag -P/--add-prefix`)

	flineageCmd.Flags().BoolP("trim", "T", false, "do not fill missing rank lower than current rank")
}
