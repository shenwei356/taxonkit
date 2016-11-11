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
	"sort"
	"strconv"
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
		nodesFile := getFlagString(cmd, "nodes")
		namesFile := getFlagString(cmd, "names")
		indent := getFlagString(cmd, "indent")
		jsonFormat := getFlagBool(cmd, "json")

		files := getFileList(args)
		if len(files) > 1 || (len(files) == 1 && files[0] == "stdin") {
			log.Warningf("no positional arguments needed")
		}

		outfh, err := xopen.Wopen(config.OutFile)
		checkError(err)

		var names map[int32]string
		var printName bool
		var printRank bool
		if namesFile != "" {
			printName = true
			printRank = true
			log.Infof("parsing names file: %s", namesFile)
			names = getTaxonNames(namesFile, config.Threads, 10)
			log.Infof("%d names parsed", len(names))
		}

		log.Infof("parsing nodes file: %s", nodesFile)

		type info struct {
			child, parent int32
			rank          string
		}
		fn := func(line string) (interface{}, bool, error) {
			items := strings.SplitN(line, "\t", 6)
			if len(items) < 6 {
				return nil, false, nil
			}
			child, e := strconv.Atoi(items[0])
			if e != nil {
				return nil, false, e
			}
			parent, e := strconv.Atoi(items[2])
			if e != nil {
				return nil, false, e
			}
			return info{int32(child), int32(parent), items[4]}, true, nil
		}

		reader, err := breader.NewBufferedReader(nodesFile, config.Threads, 10, fn)
		checkError(err)

		tree := make(map[int32]map[int32]bool)
		ranks := make(map[int32]string)
		var rel info
		var child, parent int32
		var ok bool
		var n int64
		for chunk := range reader.Ch {
			checkError(chunk.Err)

			for _, data := range chunk.Data {
				rel = data.(info)
				child, parent = rel.child, rel.parent

				if _, ok = tree[parent]; !ok {
					tree[parent] = make(map[int32]bool)
				}
				tree[parent][child] = false
				if printRank {
					ranks[child] = rel.rank
				}
				n++
			}
		}

		log.Infof("%d nodes parsed", n)

		var level int
		if jsonFormat {
			outfh.WriteString("{\n")
		}
		for i, id := range ids {
			if _, ok := tree[int32(id)]; !ok {
				continue
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

			traverseTree(tree, int32(id), outfh, indent, level+1, names,
				printName, ranks, printRank, jsonFormat)

			if jsonFormat {
				outfh.WriteString(fmt.Sprintf("%s}", strings.Repeat(indent, level)))
			}
			if jsonFormat && i < len(ids)-1 {
				outfh.WriteString(",")
			}
			outfh.WriteString("\n")
		}

		if jsonFormat {
			outfh.WriteString("}\n")
		}

		defer outfh.Close()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("nodes", "", "nodes.dmp", "nodes.dmp file")
	listCmd.Flags().StringP("ids", "", "1", "taxid(s), multiple values should be seperated by comma")

	listCmd.Flags().StringP("indent", "", "  ", "indent")
	listCmd.Flags().StringP("names", "", "", "names.dmp file, when it given taxid will be followed by its scientific name")
	listCmd.Flags().BoolP("json", "", false, `output in JSON format. you can save the result in file with suffix ".json" and open with modern text editor`)
}

func traverseTree(tree map[int32]map[int32]bool, parent int32,
	outfh *xopen.Writer, indent string, level int,
	names map[int32]string, printName bool,
	ranks map[int32]string, printRank bool,
	jsonFormat bool) {
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

		tree[parent][child] = true

		traverseTree(tree, child, outfh, indent, level+1, names, printName,
			ranks, printRank, jsonFormat)

		if jsonFormat && ok {
			outfh.WriteString(fmt.Sprintf("%s}", strings.Repeat(indent, level)))
			if level > 2 && i < len(children)-1 {
				outfh.WriteString(",")
			}
			outfh.WriteString("\n")
		}
	}
}
