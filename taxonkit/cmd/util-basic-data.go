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
	"sync"

	"github.com/shenwei356/util/pathutil"
	"github.com/shenwei356/xopen"
)

var mapInitialSize = 8 << 10

func loadData(config Config, loadTree bool, recordRank bool) (
	map[int32]int32,
	map[int32]string,
	map[int32]string,
	map[int32]struct{},
	map[int32]int32,
) {

	var tree map[int32]int32
	var ranks map[int32]string
	var names map[int32]string
	var delnodes map[int32]struct{}
	var merged map[int32]int32

	var wg sync.WaitGroup

	// tree
	if loadTree {
		wg.Add(1)
		go func() {
			if config.Verbose {
				log.Infof("parsing nodes file: %s", config.NodesFile)
			}
			tree, ranks = getNodes(config.NodesFile, recordRank)
			if config.Verbose {
				log.Infof("%d nodes parsed", len(tree))
			}
			wg.Done()
		}()
	}

	// names
	wg.Add(1)
	go func() {
		if config.Verbose {
			log.Infof("parsing names file: %s", config.NamesFile)
		}
		names = getTaxonNames(config.NamesFile)
		if config.Verbose {
			log.Infof("%d names parsed", len(names))
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if config.Verbose {
			log.Infof("parsing delnodes file: %s", config.NamesFile)
		}
		delnodes = getDelnodesMap(config.DelNodesFile)
		if config.Verbose {
			log.Infof("%d delnodes parsed", len(delnodes))
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if config.Verbose {
			log.Infof("parsing merged file: %s", config.NamesFile)
		}
		merged = getMergedNodesMap(config.MergedFile)
		if config.Verbose {
			log.Infof("%d merged nodes parsed", len(merged))
		}
		wg.Done()
	}()

	wg.Wait()

	return tree, ranks, names, delnodes, merged
}

// taxid -> name
func getTaxonNames(file string) map[int32]string {
	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	taxid2name := make(map[int32]string, mapInitialSize)

	items := make([]string, 8)
	scanner := bufio.NewScanner(fh)
	var id int
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 8, &items)
		if len(items) < 7 {
			continue
		}
		if items[6] != "scientific name" {
			continue
		}
		id, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		taxid2name[int32(id)] = items[2]
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return taxid2name
}

func getNodes(file string, recordRank bool) (map[int32]int32, map[int32]string) {
	tree := make(map[int32]int32, mapInitialSize)
	var ranks map[int32]string
	if recordRank {
		ranks = make(map[int32]string, mapInitialSize)
	}

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 6)
	scanner := bufio.NewScanner(fh)
	var _child, _parent int
	var child, parent int32
	var rank string
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 6, &items)
		if len(items) < 6 {
			continue
		}

		_child, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		_parent, err = strconv.Atoi(items[2])
		if err != nil {
			continue
		}
		child, parent, rank = int32(_child), int32(_parent), items[4]

		// ----------------------------------

		tree[child] = parent
		if recordRank {
			ranks[child] = rank
		}
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return tree, ranks
}

func getRanks(file string) map[int32]string {
	ranks := make(map[int32]string, mapInitialSize)

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 6)
	scanner := bufio.NewScanner(fh)
	var _child int
	var child int32
	var rank string
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 6, &items)
		if len(items) < 6 {
			continue
		}

		_child, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		child, rank = int32(_child), items[4]

		// ----------------------------------

		ranks[child] = rank
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return ranks
}

func getDelnodes(file string) []int32 {
	taxids := make([]int32, 0, 1<<10)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("delnodes file not found: %s, deleted taxids will not be checked", file)
		return taxids
	}

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 2)

	scanner := bufio.NewScanner(fh)
	var id int
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 2, &items)
		if len(items) < 2 {
			continue
		}
		id, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		taxids = append(taxids, int32(id))
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return taxids
}

func getDelnodesMap(file string) map[int32]struct{} {
	taxids := make(map[int32]struct{}, 1<<10)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("delnodes file not found: %s, deleted taxids will not be checked", file)
		return taxids
	}

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 2)

	scanner := bufio.NewScanner(fh)
	var id int
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 2, &items)
		if len(items) < 2 {
			continue
		}
		id, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}

		taxids[int32(id)] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return taxids
}

func getMergedNodes(file string) [][2]int32 {
	merges := make([][2]int32, 0, 1<<10)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("merged file not found: %s, merged taxids will not be checked", file)
		return merges
	}

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 4)

	scanner := bufio.NewScanner(fh)
	var from, to int
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 4, &items)
		if len(items) < 4 {
			continue
		}
		from, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}
		to, err = strconv.Atoi(items[2])
		if err != nil {
			continue
		}

		merges = append(merges, [2]int32{int32(from), int32(to)})
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return merges
}

func getMergedNodesMap(file string) map[int32]int32 {
	merges := make(map[int32]int32, 1<<10)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("merged file not found: %s, merged taxids will not be checked", file)
		return merges
	}

	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	items := make([]string, 4)

	scanner := bufio.NewScanner(fh)
	var from, to int
	for scanner.Scan() {
		stringSplitN(scanner.Text(), "\t", 4, &items)
		if len(items) < 4 {
			continue
		}
		from, err = strconv.Atoi(items[0])
		if err != nil {
			continue
		}
		to, err = strconv.Atoi(items[2])
		if err != nil {
			continue
		}

		merges[int32(from)] = int32(to)
	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return merges
}
