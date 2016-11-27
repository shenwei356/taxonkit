// Copyright © 2016 Wei Shen <shenwei356@gmail.com>
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

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "miss-rank-repl")
		fill := getFlagBool(cmd, "fill-miss-rank")

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
			checkError(fmt.Errorf("warning: stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string

		log.Infof("parsing names file: %s", config.NamesFile)
		names = getTaxonNames(config.NamesFile, config.Threads, 10)
		log.Infof("%d names parsed", len(names))

		log.Infof("parsing nodes file: %s", config.NodesFile)

		reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
		checkError(err)

		name2rank := make(map[string]string)
		var info taxonInfo
		var n int64
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data := range chunk.Data {
				info = data.(taxonInfo)
				name2rank[names[info.child]] = info.rank
				n++
			}
		}

		log.Infof("%d nodes parsed", n)

		type lineage2flineage struct {
			lineage  string
			flineage string
		}

		fn := func(line string) (interface{}, bool, error) {
			line = strings.Trim(line, "\r\n ")
			if line == "" {
				return nil, false, nil
			}

			// names and weights
			names2 := strings.Split(line, delimiter)
			weights := make([]float32, len(names2))
			var rank, srank string
			var ok bool
			srank2name := make(map[string]string)
			var currentWeight float32
			for i, name := range names2 {
				if name == "" {
					continue
				}
				rank, ok = name2rank[name]
				if !ok { // unofficial name
					currentWeight += 0.1
					weights[i] = currentWeight
					continue
				}
				if rank != norank {
					if srank, ok = rank2symbol[rank]; ok {
						srank2name[srank] = name
						weights[i] = symbol2weight[srank]
						currentWeight = weights[i]
					} else {
						log.Warningf("please contact author to add this rank to code: %s", rank)
					}
				} else {
					currentWeight += 0.1
					weights[i] = currentWeight
				}
			}

			// preprare replacements.
			// find the orphan names and missing ranks
			replacements := make(map[string]string, len(matches))
			if !fill {
				for _, match := range matches {
					replacements[match[1]] = blank
				}
			}

			orphans := make(map[string]float32)
			orphansList := []string{}
			existedSranks := make(map[string]struct{})
			for i, name := range names2 {
				if name == "" {
					continue
				}
				if name2rank[name] == norank {
					orphans[name] = weights[i]
					orphansList = append(orphansList, name)
				} else {
					if _, ok = outSranks[rank2symbol[name2rank[name]]]; ok { // to be outputted
						replacements[rank2symbol[name2rank[name]]] = name
						existedSranks[rank2symbol[name2rank[name]]] = struct{}{}
					} else if name2rank[name] == "" {
						orphans[name] = weights[i]
						orphansList = append(orphansList, name)
					}
				}
			}

			if fill {
				jj := -1
				var hit bool
				var lastRank string
				for i, srank := range outSranksList {
					if _, ok = existedSranks[srank]; ok {
						lastRank = replacements[srank]
						continue
					}
					hit = false
					for j, name := range orphansList {
						if j <= jj {
							continue
						}
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
							replacements[srank] = name
							jj = j
							break
						}
					}
					if !hit {
						if blank == "" {
							replacements[srank] = fmt.Sprintf("unclassified %s %s", lastRank, symbol2rank[srank])
						} else {
							replacements[srank] = blank
						}
					}
				}
			}

			flineage := format
			for srank, re := range reRankPlaceHolders {
				flineage = re.ReplaceAllString(flineage, replacements[srank])
			}

			return lineage2flineage{line, flineage}, true, nil
		}

		for _, file := range files {
			reader, err := breader.NewBufferedReader(file, config.Threads, 10, fn)
			checkError(err)

			var l2s lineage2flineage
			for chunk := range reader.Ch {
				checkError(chunk.Err)

				for _, data := range chunk.Data {
					l2s = data.(lineage2flineage)
					outfh.WriteString(fmt.Sprintf("%s\t%s\n", l2s.lineage, l2s.flineage))
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
	flineageCmd.Flags().StringP("miss-rank-repl", "r", "", `replacement string for missing rank, if given "", "unclassified xxx xxx" will used`)
	flineageCmd.Flags().BoolP("fill-miss-rank", "F", false, "estimate and fill missing rank with original lineage information (recommended)")
}
