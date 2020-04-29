// Copyright © 2016-2020 Wei Shen <shenwei356@gmail.com>
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
	"runtime"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// flineageCmd represents the fx2tab command
var flineageCmd = &cobra.Command{
	Use:   "reformat",
	Short: "reformat lineage",
	Long: `reformat lineage

Output format can be formated by flag --format, available placeholders:

    {k}: superkingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {S}: subspecies

Output format can contains some escape charactors like "\t".

This command appends reformated lineage to the input line.
The corresponding taxids of reformated lineage can be provided as another
column by flag "-t/--show-lineage-taxids".

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "miss-rank-repl")
		prefix := getFlagString(cmd, "miss-rank-repl-prefix")
		iblank := getFlagString(cmd, "miss-taxid-repl")
		fill := getFlagBool(cmd, "fill-miss-rank")
		field := getFlagPositiveInt(cmd, "lineage-field") - 1
		printLineageInTaxid := getFlagBool(cmd, "show-lineage-taxids")

		// check format
		if !reRankPlaceHolder.MatchString(format) {
			checkError(fmt.Errorf("placeholder of simplified rank not found in output format: %s", format))
		}
		matches := reRankPlaceHolder.FindAllStringSubmatch(format, -1)
		outSranks := make(map[string]struct{})
		outSranksList := []string{}
		var currentWeight float32
		var currentSymbol string
		for _, match := range matches {
			if weight, ok := symbol2weight[match[1]]; !ok {
				checkError(fmt.Errorf("invalid placeholder: %s", match[0]))
			} else {
				if weight < currentWeight {
					checkError(fmt.Errorf(`invalid placeholder order: {%s} {%s}. "%s" should be behind of "%s"`,
						currentSymbol, match[1], symbol2rank[currentSymbol], symbol2rank[match[1]]))
				}
				outSranks[match[1]] = struct{}{}
				outSranksList = append(outSranksList, match[1])
				currentWeight = weight
				currentSymbol = match[1]
			}
		}

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		if config.Verbose {
			log.Infof("parsing names (%s) and nodes file (%s)", config.NamesFile, config.NodesFile)
		}

		taxid2taxon, name2parent2taxid, name2taxid := getName2Parent2Taxid(config.NodesFile, config.NamesFile, config.Threads, 20)

		if config.Verbose {
			log.Infof("%d nodes parsed", len(taxid2taxon))
		}

		type line2flineage struct {
			line      string
			flineage  string
			iflineage string
		}

		unescape := stringutil.UnEscaper()
		fn := func(line string) (interface{}, bool, error) {
			if len(line) == 0 || line[0] == '#' {
				return nil, false, nil
			}
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}
			data := strings.Split(line, "\t")
			if len(data) < field+1 {
				return nil, false, fmt.Errorf("lineage-field (%d) out of range (%d):%s", field+1, len(data), line)
			}

			// names and weights
			names2 := strings.Split(data[field], delimiter)
			weights := make([]float32, len(names2))
			var rank, srank string
			var ok bool
			srank2name := make(map[string]string)
			var currentWeight float32
			name2Name := make(map[string]string, len(names2)) // lower case of name : name
			var lname string
			var plname string
			for i, name := range names2 {
				if name == "" {
					continue
				}

				lname = strings.ToLower(name)
				name2Name[lname] = name
				name = lname

				if _, ok = name2taxid[name]; !ok { // unofficial name
					log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, name)
					return line2flineage{line, "", ""}, true, nil
				}

				if i == 0 {
					rank = taxid2taxon[name2taxid[name]].Rank
				} else {
					plname = strings.ToLower(names2[i-1])
					if _, ok = name2parent2taxid[name]; !ok {
						log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, name)
						return line2flineage{line, "", ""}, true, nil
					} else if _, ok = name2parent2taxid[name][plname]; !ok {
						log.Warningf(`unofficial taxon name detected: %s. Possible reasons: 1) lineages were produced with different taxonomy data files, please re-run taxonkit lineage; 2) some taxon names contain semicolon (";"), please re-run taxonkit lineage and taxonkit reformat with different flag value of -d, e.g., -d /`, plname)
						return line2flineage{line, "", ""}, true, nil
					}
					rank = taxid2taxon[name2parent2taxid[name][plname]].Rank
				}
				if rank != norank {
					if srank, ok = rank2symbol[rank]; ok {
						srank2name[srank] = name
						weights[i] = symbol2weight[srank]
						currentWeight = weights[i]
					} else {
						// log.Warningf("please contact author to add this rank to code: %s", rank)
					}
				} else {
					currentWeight += 0.1
					weights[i] = currentWeight
				}
			}
			// preprare replacements.
			// find the orphan names and missing ranks
			replacements := make(map[string]string, len(matches))
			ireplacements := make(map[string]string, len(matches))
			if !fill {
				for _, match := range matches {
					replacements[match[1]] = blank
					ireplacements[match[1]] = blank
				}
			}

			orphans := make(map[string]float32)
			orphansList := make([][2]string, 0, 20)
			existedSranks := make(map[string]struct{})
			for i, name := range names2 {
				if name == "" {
					continue
				}
				name = strings.ToLower(name)

				if i == 0 {
					rank = taxid2taxon[name2taxid[name]].Rank
				} else {
					plname = strings.ToLower(names2[i-1])
					rank = taxid2taxon[name2parent2taxid[name][plname]].Rank
				}

				if rank == norank {
					orphans[name] = weights[i]
					orphansList = append(orphansList, [2]string{name, plname})
				} else {
					if _, ok = outSranks[rank2symbol[rank]]; ok { // to be outputted
						replacements[rank2symbol[rank]] = name2Name[name]
						if i == 0 {
							ireplacements[rank2symbol[rank]] = fmt.Sprintf("%d",
								name2taxid[name])
						} else {
							ireplacements[rank2symbol[rank]] = fmt.Sprintf("%d",
								name2parent2taxid[name][plname])
						}

						existedSranks[rank2symbol[rank]] = struct{}{}
					} else if rank == "" {
						orphans[name] = weights[i]
						orphansList = append(orphansList, [2]string{name, plname})
					}
				}
			}

			if fill {
				jj := -1
				var hit bool
				var lastRank string
				var name, pname string
				for i, srank := range outSranksList {
					if _, ok = existedSranks[srank]; ok {
						lastRank = replacements[srank]
						continue
					}
					hit = false
					for j, n2p := range orphansList {
						if j <= jj {
							continue
						}
						name, pname = n2p[0], n2p[1]
						if i == 0 {
							if orphans[name] < symbol2weight[outSranksList[i]] {
								hit = true
							}
						} else if i == len(outSranksList)-1 {

						} else if orphans[name] > symbol2weight[outSranksList[i-1]] &&
							orphans[name] < symbol2weight[outSranksList[i+1]] {
							hit = true
						}

						if hit {
							replacements[srank] = name2Name[name]
							ireplacements[srank] = fmt.Sprintf("%d", name2parent2taxid[name][pname])
							jj = j
							break
						}
					}
					hit = false
					if !hit {
						if blank == "" {
							replacements[srank] = fmt.Sprintf("%s%s %s", prefix, lastRank, symbol2rank[srank])
						} else {
							replacements[srank] = blank
						}
						ireplacements[srank] = iblank
					}
				}
			}

			flineage := format
			iflineage := format
			for srank, re := range reRankPlaceHolders {
				flineage = re.ReplaceAllString(flineage, replacements[srank])
				iflineage = re.ReplaceAllString(iflineage, ireplacements[srank])
			}

			return line2flineage{line, unescape(flineage), unescape(iflineage)}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var l2s line2flineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					l2s = data.(line2flineage)

					if printLineageInTaxid {
						outfh.WriteString(fmt.Sprintf("%s\t%s\t%s\n", l2s.line, l2s.flineage, l2s.iflineage))
					} else {
						outfh.WriteString(fmt.Sprintf("%s\t%s\n", l2s.line, l2s.flineage))
					}
					if config.LineBuffered {
						outfh.Flush()
					}
				}
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(flineageCmd)

	flineageCmd.Flags().StringP("format", "f", "{k};{p};{c};{o};{f};{g};{s}", "output format, placeholders of rank are needed")
	flineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in input lineage")
	flineageCmd.Flags().StringP("miss-rank-repl", "r", "", `replacement string for missing rank, if given "", "unclassified xxx xxx" will used, where "unclassified " is settable by flag -p/--miss-rank-repl-prefix`)
	flineageCmd.Flags().StringP("miss-rank-repl-prefix", "p", "unclassified ", `prefix for estimated taxon level`)
	flineageCmd.Flags().StringP("miss-taxid-repl", "R", "", `replacement string for missing taxid`)
	flineageCmd.Flags().BoolP("fill-miss-rank", "F", false, "fill missing rank with original lineage information (experimental)")
	flineageCmd.Flags().IntP("lineage-field", "i", 2, "field index of lineage. data should be tab-separated")
	flineageCmd.Flags().BoolP("show-lineage-taxids", "t", false, `show corresponding taxids of reformated lineage`)
}
