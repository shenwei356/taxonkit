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
	"sort"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// listCmd represents the fx2tab command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list taxon tree of given taxids",
	Long: `list taxon tree of given taxids

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		ids := getFlagTaxonIDs(cmd, "ids")
		indent := getFlagString(cmd, "indent")
		jsonFormat := getFlagBool(cmd, "json")

		files := getFileList(args)
		if len(files) > 1 || (len(files) == 1 && files[0] == "stdin") {
			log.Warningf("no positional arguments needed")
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		printName := getFlagBool(cmd, "show-name")
		printRank := getFlagBool(cmd, "show-rank")

		var names map[int32]string
		var delnodes map[int32]struct{}
		var merged map[int32]int32

		if config.Verbose {
			log.Infof("parsing delnodes file: %s", config.NamesFile)
		}

		delnodes = getDelnodesMap(config.DelNodesFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d delnodes parsed", len(delnodes))

			log.Infof("parsing merged file: %s", config.NamesFile)
		}

		merged = getMergedNodesMap(config.MergedFile, config.Threads, 10)

		if config.Verbose {
			log.Infof("%d merged nodes parsed", len(merged))
		}

		if printName {
			if config.Verbose {
				log.Infof("parsing names file: %s", config.NamesFile)
			}
			names = getTaxonNames(config.NamesFile, config.Threads, 10)
			if config.Verbose {
				log.Infof("%d names parsed", len(names))
			}
		}

		if config.Verbose {
			log.Infof("parsing nodes file: %s", config.NodesFile)
		}

		reader, err := breader.NewBufferedReader(config.NodesFile, config.Threads, 10, taxonParseFunc)
		checkError(err)

		tree := make(map[int32]map[int32]bool)
		ranks := make(map[int32]string)
		var taxon Taxon
		var child, parent int32
		var ok bool
		var n int64
		var data interface{}
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data = range chunk.Data {
				taxon = data.(Taxon)
				child, parent = taxon.Taxid, taxon.Parent

				if _, ok = tree[parent]; !ok {
					tree[parent] = make(map[int32]bool)
				}
				tree[parent][child] = false
				if _, ok = tree[child]; !ok {
					tree[child] = make(map[int32]bool)
				}
				if printRank {
					ranks[child] = taxon.Rank
				}
				n++
			}
		}

		if config.Verbose {
			log.Infof("%d nodes parsed", n)
		}

		var level int
		if jsonFormat {
			outfh.WriteString("{\n")
		}
		var newtaxid int32
		for i, id := range ids {
			if _, ok := tree[int32(id)]; !ok {
				// check if it was deleted
				if _, ok = delnodes[int32(id)]; ok {
					log.Warningf("taxid %d was deleted", child)
					continue
				}
				// check if it was merged
				if newtaxid, ok = merged[int32(id)]; ok {
					log.Warningf("taxid %d was merged into %d", child, newtaxid)
					id = int(newtaxid)
				} else {
					log.Warningf("taxid %d not found", child)
					continue
				}
			}

			level = 0
			if jsonFormat {
				level = 1
			}

			outfh.WriteString(strings.Repeat(indent, level))

			if jsonFormat {
				outfh.WriteString(`"`)
			}
			outfh.WriteString(fmt.Sprintf("%d", id))

			if printRank {
				outfh.WriteString(fmt.Sprintf(" [%s]", ranks[int32(id)]))
			}
			if printName {
				outfh.WriteString(fmt.Sprintf(" %s", names[int32(id)]))
			}

			level = 0
			if jsonFormat {
				outfh.WriteString(`": {`)
				level = 1
			}
			outfh.WriteString("\n")
			if config.LineBuffered {
				outfh.Flush()
			}

			traverseTree(tree, int32(id), outfh, indent, level+1, names,
				printName, ranks, printRank, jsonFormat, config)

			if jsonFormat {
				outfh.WriteString(fmt.Sprintf("%s}", strings.Repeat(indent, level)))
			}
			if jsonFormat && i < len(ids)-1 {
				outfh.WriteString(",")
			}
			outfh.WriteString("\n")
			if config.LineBuffered {
				outfh.Flush()
			}
		}

		if jsonFormat {
			outfh.WriteString("}\n")
			if config.LineBuffered {
				outfh.Flush()
			}
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("ids", "", "", "taxid(s), multiple values should be separated by comma")
	listCmd.Flags().StringP("indent", "", "  ", "indent")
	listCmd.Flags().BoolP("show-rank", "r", false, `output rank`)
	listCmd.Flags().BoolP("show-name", "n", false, `output scientific name`)
	listCmd.Flags().BoolP("json", "", false, `output in JSON format. you can save the result in file with suffix ".json" and open with modern text editor`)
}

func traverseTree(tree map[int32]map[int32]bool, parent int32,
	outfh *xopen.Writer, indent string, level int,
	names map[int32]string, printName bool,
	ranks map[int32]string, printRank bool,
	jsonFormat bool, config Config) {
	if _, ok := tree[parent]; !ok {
		return
	}

	// sort children by taxid
	children := make([]int, len(tree[parent]))
	i := 0
	for child := range tree[parent] {
		children[i] = int(child)
		i++
	}
	sort.Ints(children)

	var child int32
	for i, c := range children {
		child = int32(c)
		if tree[parent][child] {
			continue
		}

		outfh.WriteString(strings.Repeat(indent, level))

		if jsonFormat {
			outfh.WriteString(`"`)
		}
		outfh.WriteString(fmt.Sprintf("%d", child))
		if printRank {
			outfh.WriteString(fmt.Sprintf(" [%s]", ranks[child]))
		}
		if printName {
			outfh.WriteString(fmt.Sprintf(" %s", names[child]))
		}

		var ok bool
		if jsonFormat {
			_, ok = tree[child]
			if ok {
				outfh.WriteString(`": {`)
			} else {
				outfh.WriteString(`": {}`)
				if i < len(children)-1 {
					outfh.WriteString(",")
				}
			}
		}
		outfh.WriteString("\n")
		if config.LineBuffered {
			outfh.Flush()
		}

		tree[parent][child] = true

		traverseTree(tree, child, outfh, indent, level+1, names, printName,
			ranks, printRank, jsonFormat, config)

		if jsonFormat && ok {
			outfh.WriteString(fmt.Sprintf("%s}", strings.Repeat(indent, level)))
			if level > 1 && i < len(children)-1 {
				outfh.WriteString(",")
			}
			outfh.WriteString("\n")
			if config.LineBuffered {
				outfh.Flush()
			}
		}
	}
}
