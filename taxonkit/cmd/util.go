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

import "sync"

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
