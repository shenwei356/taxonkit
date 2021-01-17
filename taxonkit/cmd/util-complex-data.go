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
	"bufio"
	"strconv"
	"strings"

	"github.com/shenwei356/xopen"
)

// names -> []taxid
func getTaxonName2Taxids(file string, limit2SciName bool) map[string][]int32 {
	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	name2taxids := make(map[string][]int32, mapInitialSize)

	items := make([]string, 8)
	scanner := bufio.NewScanner(fh)
	var id int
	var name string
	var ok bool
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 8, &items)
		if len(items) < 7 {
			continue
		}
		if limit2SciName {
			if items[6] != "scientific name" {
				continue
			}
		} else {
			if !(items[6] == "scientific name" || items[6] == "synonym") {
				continue
			}
		}
		name = items[2]

		id, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		// -------------

		name = strings.ToLower(name)
		if _, ok = name2taxids[name]; !ok {
			name2taxids[name] = []int32{int32(id)}
		} else {
			name2taxids[name] = append(name2taxids[name], int32(id))
		}

	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return name2taxids
}

// taxid -> lineage
func getTaxid2LineageTaxids(fileNodes string) (
	map[int32][]int32, // taxid2lineageTaxids
	map[int32]string, // taxid2rank
) {
	var tree map[int32]int32
	var ranks map[int32]string

	tree, ranks = getNodes(fileNodes, true)

	taxid2lineageTaxids := make(map[int32][]int32, mapInitialSize)

	var ok bool
	var i, j int
	var child, parent int32
	for taxid := range tree {
		lineageTaxids := make([]int32, 0, 8)
		child = taxid
		for true {
			parent, ok = tree[child]
			if !ok {
				break
			}

			lineageTaxids = append(lineageTaxids, child)

			if parent == 1 {
				break
			}
			child = parent
		}

		// reverse lineageTaxids in place
		for i = len(lineageTaxids)/2 - 1; i >= 0; i-- {
			j = len(lineageTaxids) - 1 - i
			lineageTaxids[i], lineageTaxids[j] = lineageTaxids[j], lineageTaxids[i]
		}
		taxid2lineageTaxids[taxid] = lineageTaxids
	}
	return taxid2lineageTaxids, ranks
}
