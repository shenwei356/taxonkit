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
	"fmt"
	"strconv"
	"strings"

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

Warning:

  - 'taxonkit reformat2' is recommended since Match 2025 when NCBI made
    big changes to ranks.
    See more: https://ncbiinsights.ncbi.nlm.nih.gov/2025/02/27/new-ranks-ncbi-taxonomy/

Input:

  - List of TaxIds or lineages, one record per line.
    The lineage can be a complete lineage or only one taxonomy name.
  - Or tab-delimited format.
    Plese specify the lineage field with flag -i/--lineage-field (default 2).
    Or specify the TaxId field with flag -I/--taxid-field (default 0),
    which overrides -i/--lineage-field.
  - Supporting (gzipped) file or STDIN.

Output:

  1. Input line data.
  2. Reformated lineage.
  3. (Optional) TaxIds taxons in the lineage (-t/--show-lineage-taxids)
  
Ambiguous names:

  - Some TaxIds have the same complete lineage, empty result is returned 
    by default. You can use the flag -a/--output-ambiguous-result to
    return one possible result

Output format can be formated by flag --format, available placeholders:

    {C}: cellular root
    {a}: acellular root
    {r}: realm
    {d}: domain
    {k}: superkingdom
    {K}: kingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {t}: subspecies/strain
    
    {S}: subspecies
    {T}: strain

When these're no nodes of rank "subspecies" nor "strain",
you can switch on -S/--pseudo-strain to use the node with lowest rank
as subspecies/strain name, if which rank is lower than "species". 
This flag affects {t}, {S}, {T}.
    
Output format can contains some escape charactors like "\t".

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)

		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "miss-rank-repl")
		prefix := getFlagString(cmd, "miss-rank-repl-prefix")
		suffix := getFlagString(cmd, "miss-rank-repl-suffix")

		iblank := getFlagString(cmd, "miss-taxid-repl")
		fill := getFlagBool(cmd, "fill-miss-rank")
		pseudoStrain := getFlagBool(cmd, "pseudo-strain")

		taxIdField := getFlagNonNegativeInt(cmd, "taxid-field")
		field := getFlagPositiveInt(cmd, "lineage-field")
		outputAmbigous := getFlagBool(cmd, "output-ambiguous-result")

		var parsingTaxId bool
		if taxIdField > 0 {
			if config.Verbose {
				log.Infof("parsing TaxIds from field %d", taxIdField)
			}
			parsingTaxId = true
			taxIdField--
		} else if field > 0 {
			if config.Verbose {
				log.Infof("parsing complete lineages from field %d", field)
			}
			field--
		}

		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")

		addPrefix := getFlagBool(cmd, "add-prefix")
		prefixR := getFlagString(cmd, "prefix-r")
		prefixC2 := getFlagString(cmd, "prefix-C")
		prefixD := getFlagString(cmd, "prefix-d")
		prefixA := getFlagString(cmd, "prefix-a")
		prefixK := getFlagString(cmd, "prefix-k")
		prefixK2 := getFlagString(cmd, "prefix-K")
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
			"C": prefixC2,
			"a": prefixA,
			"d": prefixD,
			"r": prefixR,
			"k": prefixK,
			"K": prefixK2,
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
			// do not require this.
			// if pseudoStrain && !fill {
			// 	fill = true
			// 	log.Infof("-F/--fill-miss-rank is switched on when giving flag -S/--pseudo-strain")
			// }
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

		var tree0 map[uint32]uint32
		var ranks0 map[uint32]string
		var names0 map[uint32]string
		var delnodes0 map[uint32]struct{}
		var merged0 map[uint32]uint32

		tree0, ranks0, names0, delnodes0, merged0 = loadData(config, true, true)

		// for querying taxid from lineage
		var name2parent2taxid map[string]map[string]uint32
		var name2taxids map[string]*[]uint32
		var ambigous map[string][]uint32

		if !parsingTaxId {
			name2parent2taxid, name2taxids, ambigous = generateName2Parent2Taxid(config, tree0, names0)
		}

		// --------------------------------------------------------

		type line2flineage struct {
			line      string
			flineage  string
			iflineage string
		}

		unescape := stringutil.UnEscaper()

		weightOfSpecies := symbol2weight["s"]

		blankS := format
		iblankS := format
		blankS = reRankPlaceHolder.ReplaceAllString(blankS, blank)
		iblankS = reRankPlaceHolder.ReplaceAllString(iblankS, iblank)

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

			var ok bool

			var taxid uint32
			var taxidInt int

			var names []string
			var ranks []string
			var taxids []uint32

			// -----------------------------------------------
			// get the taxid

			if parsingTaxId { // directly from field

				taxidInt, err = strconv.Atoi(data[taxIdField])
				if err != nil || taxidInt < 0 {
					// checkError(fmt.Errorf("invalid TaxId: %s", data[taxIdField]))
					log.Warningf("invalid TaxId: %s", data[taxIdField])
					return line2flineage{line, "", ""}, true, nil
				}
				taxid = uint32(taxidInt)

			} else { // query taxid by taxon names

				if strings.Trim(data[field], " ") == "" { // empty, returns empty result
					return line2flineage{line, "", ""}, true, nil
				}

				// names
				names = strings.Split(data[field], delimiter)
				n := len(names)

				// name and name of its parent
				var name, pname string

				if n == 1 { // single name

					// direct query via name2taxids
					_taxids := name2taxids[strings.ToLower(names[0])]
					if _taxids == nil {
						log.Warningf(`failed to query the TaxId of: %s. Possible reasons: `, data[field])
						log.Warningf(`  1) the lineage were produced with different taxonomy data files, please re-run taxonkit lineage;`)
						log.Warningf(`  2) some taxon names contain delimiter (%s), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d "/"`, delimiter)
						return line2flineage{line, "", ""}, true, nil
					}

					if len(*_taxids) == 1 { // found
						taxid = (*_taxids)[0]
					} else { // ambiguous name
						tmp := make([]string, len(*_taxids))
						for _i, _taxid := range *_taxids {
							tmp[_i] = strconv.Itoa(int(_taxid))
						}
						log.Warningf(`we can't distinguish the TaxIds (%s) for lineage: %s. But you can use -a/--output-ambiguous-result to return one possible result`,
							strings.Join(tmp, ", "), data[field])

						if !outputAmbigous {
							return line2flineage{line, "", ""}, true, nil
						}
					}

				} else { // multiple names

					name = strings.ToLower(names[n-1])  // name
					pname = strings.ToLower(names[n-2]) // parent name
					var tmp map[string]uint32
					found := true

					if tmp, ok = name2parent2taxid[name]; !ok {
						found = false
					}
					if taxid, ok = tmp[pname]; !ok {
						found = false
					}
					if !found { // try use the last node
						// direct query via name2taxids
						_taxids := name2taxids[strings.ToLower(names[n-1])]
						if _taxids == nil {
							log.Warningf(`failed to query the TaxId of: %s. Possible reasons: `, data[field])
							log.Warningf(`  1) the lineage were produced with different taxonomy data files, please re-run taxonkit lineage;`)
							log.Warningf(`  2) some taxon names contain delimiter (%s), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d "/"`, delimiter)
							return line2flineage{line, "", ""}, true, nil
						}

						if len(*_taxids) == 1 { // found
							taxid = (*_taxids)[0]
						} else { // ambiguous name
							tmp := make([]string, len(*_taxids))
							for _i, _taxid := range *_taxids {
								tmp[_i] = strconv.Itoa(int(_taxid))
							}
							log.Warningf(`we can't distinguish the TaxIds (%s) for lineage: %s. But you can use -a/--output-ambiguous-result to return one possible result`,
								strings.Join(tmp, ", "), data[field])

							if !outputAmbigous {
								return line2flineage{line, "", ""}, true, nil
							}
						}
					} else {
						// for cases where child-parent pairs are shared by multiple taxids.
						pair := name + "__" + pname
						var _ambids []uint32
						if _ambids, ok = ambigous[pair]; ok {
							tmp := make([]string, len(_ambids))
							for _i, _taxid := range _ambids {
								tmp[_i] = strconv.Itoa(int(_taxid))
							}
							log.Warningf("we can't distinguish the TaxIds (%s) for lineage: %s. But you can use -a/--output-ambiguous-result to return one possible result",
								strings.Join(tmp, ", "), data[field])

							if !outputAmbigous {
								return line2flineage{line, "", ""}, true, nil
							}
						}
					}
				}
			}

			// -----------------------------------------------
			// query complete lineage with the taxid

			names, ranks, taxids, ok = queryNamesRanksTaxids(tree0, ranks0, names0, delnodes0, merged0, taxid)
			if !ok { // taxid not found
				// return line2flineage{line, "", ""}, true, nil
				return line2flineage{line, unescape(blankS), unescape(iblankS)}, true, nil
			}

			sranks := poolStringsN16.Get().([]string)

			srank2idx := make(map[string]int) // srank: index

			var maxRankWeight float32
			var rank, srank string // lower case of name : name

			// preprare replacements.
			// find the orphan names and missing ranks
			replacements := make(map[string]string, len(matches))

			var ireplacements map[string]string
			if printLineageInTaxid {
				ireplacements = make(map[string]string, len(matches))
			}

			for i, name := range names {
				rank = ranks[i]
				taxid = taxids[i]

				if srank, ok = rank2symbol[rank]; ok {
					// special symbol "{t}"
					switch rank {
					case "strain":
						replacements["t"] = name
						if printLineageInTaxid {
							ireplacements["t"] = strconv.Itoa(int(taxid))
						}
						srank2idx["t"] = i
					case "subspecies":
						replacements["t"] = name
						if printLineageInTaxid {
							ireplacements["t"] = strconv.Itoa(int(taxid))
						}
						srank2idx["t"] = i
					}

					replacements[srank] = name
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

					if suffix == "rank" {
						replacements[srank] = prefix + names[lastI] + " " + symbol2rank[srank]
						// replacements[srank] = fmt.Sprintf("%s%s %s", prefix, names[lastI], symbol2rank[srank])
					} else {
						replacements[srank] = prefix + names[lastI] + suffix
					}
				}
			}

			if pseudoStrain {
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

					if symbol2weight[srank] > weightOfSpecies && // lower than species
						!(hasRankSubspecies || hasRankStrain) && // does not have strain or subspecies
						lastI < len(names)-1 { // not itself
						replacements[srank] = names[len(names)-1]
						continue
					}
				}
			}

			flineage := format
			var iflineage string

			if printLineageInTaxid {
				iflineage = format
			}

			for srank, re := range reRankPlaceHolders {
				if replacements[srank] == "" {
					if trim && symbol2weight[srank] > maxRankWeight {
						flineage = re.ReplaceAllString(flineage, "")

						if printLineageInTaxid {
							iflineage = re.ReplaceAllString(iflineage, "")
						}
					} else {
						flineage = re.ReplaceAllString(flineage, blank)

						if printLineageInTaxid {
							iflineage = re.ReplaceAllString(iflineage, iblank)
						}
					}
				} else {
					if addPrefix {
						flineage = re.ReplaceAllString(flineage, prefixes[srank]+replacements[srank])
					} else {
						flineage = re.ReplaceAllString(flineage, replacements[srank])
					}

					if printLineageInTaxid {
						iflineage = re.ReplaceAllString(iflineage, ireplacements[srank])
					}
				}
			}

			// recycle
			ranks = ranks[:0]
			poolStringsN16.Put(ranks)
			sranks = sranks[:0]
			poolStringsN16.Put(sranks)

			names = names[:0]
			poolStringsN16.Put(names)
			taxids = taxids[:0]
			poolUint32N16.Put(taxids)

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
	flineageCmd.Flags().StringP("miss-rank-repl-prefix", "p", "unclassified ", `prefix for estimated taxon names`)
	flineageCmd.Flags().StringP("miss-rank-repl-suffix", "s", "rank", `suffix for estimated taxon names. "rank" for rank name, "" for no suffix`)
	flineageCmd.Flags().StringP("miss-taxid-repl", "R", "", `replacement string for missing taxid`)

	flineageCmd.Flags().BoolP("fill-miss-rank", "F", false, "fill missing rank with lineage information of the next higher rank")
	flineageCmd.Flags().BoolP("pseudo-strain", "S", false, `use the node with lowest rank as strain name, only if which rank is lower than "species" and not "subpecies" nor "strain". It affects {t}, {S}, {T}. This flag needs flag -F`)

	flineageCmd.Flags().IntP("lineage-field", "i", 2, "field index of lineage. data should be tab-separated")
	flineageCmd.Flags().IntP("taxid-field", "I", 0, "field index of taxid. input data should be tab-separated. it overrides -i/--lineage-field")
	flineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `show corresponding taxids of reformated lineage`)
	flineageCmd.Flags().BoolP("output-ambiguous-result", "a", false, `output one of the ambigous result`)

	flineageCmd.Flags().BoolP("add-prefix", "P", false, `add prefixes for all ranks, single prefix for a rank is defined by flag --prefix-X`)
	flineageCmd.Flags().StringP("prefix-r", "", "r__", `prefix for realm, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-d", "", "d__", `prefix for domain, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-a", "", "d__", `prefix for acellular root, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-C", "", "d__", `prefix for cellular root, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-k", "", "k__", `prefix for superkingdom, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-K", "", "K__", `prefix for kingdom, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-p", "", "p__", `prefix for phylum, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-c", "", "c__", `prefix for class, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-o", "", "o__", `prefix for order, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-f", "", "f__", `prefix for family, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-g", "", "g__", `prefix for genus, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-s", "", "s__", `prefix for species, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-t", "", "t__", `prefix for subspecies/strain, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-S", "", "S__", `prefix for subspecies, used along with flag -P/--add-prefix`)
	flineageCmd.Flags().StringP("prefix-T", "", "T__", `prefix for strain, used along with flag -P/--add-prefix`)

	flineageCmd.Flags().BoolP("trim", "T", false, "do not fill or add prefix for missing rank lower than current rank")
}
