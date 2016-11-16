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

		nodesFile := getFlagString(cmd, "nodes")
		namesFile := getFlagString(cmd, "names")
		format := getFlagString(cmd, "format")
		delimiter := getFlagString(cmd, "delimiter")
		blank := getFlagString(cmd, "blank")

		if !reRankPlaceHolder.MatchString(format) {
			checkError(fmt.Errorf("placeholder of simplified rank not found in output format: %s", format))
		}

		files := getFileList(args)

		if len(files) == 1 && isStdin(files[0]) && !xopen.IsStdin() {
			checkError(fmt.Errorf("warning: stdin not detected"))
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string

		log.Infof("parsing names file: %s", namesFile)
		names = getTaxonNames(namesFile, config.Threads, 10)
		log.Infof("%d names parsed", len(names))

		log.Infof("parsing nodes file: %s", nodesFile)

		reader, err := breader.NewBufferedReader(nodesFile, config.Threads, 10, taxonParseFunc)
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

			var rank, srank, name string
			var ok bool
			srank2name := make(map[string]string)
			for _, name := range strings.Split(line, delimiter) {
				if name == "" {
					continue
				}
				if rank, ok = name2rank[name]; ok && rank != norank {
					if srank, ok = rank2symbol[rank]; ok {
						srank2name[srank] = name
					}
				}
			}

			flineage := format
			for srank, re := range reRankPlaceHolders {
				if name, ok = srank2name[srank]; ok {
					flineage = re.ReplaceAllString(flineage, name)
				} else {
					flineage = re.ReplaceAllString(flineage, blank)
				}
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

	flineageCmd.Flags().StringP("nodes", "", "nodes.dmp", "nodes.dmp file")
	flineageCmd.Flags().StringP("names", "", "names.dmp", "names.dmp file")
	flineageCmd.Flags().StringP("format", "f", "{k};{p};{c};{o};{f};{g};{s}", "output format, placeholder of is need")
	flineageCmd.Flags().StringP("delimiter", "d", ";", "field delimiter in input lineage")
	flineageCmd.Flags().StringP("blank", "b", "__", "blank string for missing level")
}
