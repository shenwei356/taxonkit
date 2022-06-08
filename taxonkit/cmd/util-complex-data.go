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
	"bufio"
	"strconv"
	"strings"
	"sync"

	"github.com/shenwei356/util/stringutil"
	"github.com/shenwei356/xopen"
)

// ----------------------------------  name2taxid ---------------------------

// names -> []taxid
func getTaxonName2Taxids(file string, limit2SciName bool) map[string][]uint32 {
	fh, err := xopen.Ropen(file)
	checkError(err)
	defer func() {
		checkError(fh.Close())
	}()

	name2taxids := make(map[string][]uint32, mapInitialSize)

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
			name2taxids[name] = []uint32{uint32(id)}
		} else {
			name2taxids[name] = append(name2taxids[name], uint32(id))
		}

	}
	if err := scanner.Err(); err != nil {
		checkError(err)
	}

	return name2taxids
}

// ----------------------------------  taxid-changelog ---------------------------

// taxid -> lineageTaxids
func getTaxid2LineageTaxids(fileNodes string) (
	map[uint32][]uint32, // taxid2lineageTaxids
	map[uint32]string, // taxid2rank
) {
	var tree map[uint32]uint32
	var ranks map[uint32]string

	tree, ranks = getNodes(fileNodes, true)

	taxid2lineageTaxids := make(map[uint32][]uint32, mapInitialSize)

	var ok bool
	var i, j int
	var child, parent uint32
	for taxid := range tree {
		lineageTaxids := make([]uint32, 0, 8)
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

// ----------------------------------  reformat ---------------------------

func generateName2Parent2Taxid(
	config Config,
	tree map[uint32]uint32,
	names map[uint32]string,
) (
	map[string]map[string]uint32,
	map[string]*[]uint32,
	map[string][]uint32,
) {
	if config.Verbose {
		log.Infof("creating links: child name -> parent name -> taxid")
	}

	// name -> parent-name -> taxid
	name2parent2taxid := make(map[string]map[string]uint32, mapInitialSize)

	// name -> taxids
	name2taxids := make(map[string]*[]uint32, mapInitialSize)

	// name__prent-name -> taxids
	ambigous := make(map[string][]uint32, 128)

	var name, pname string
	var _n2i map[string]uint32
	var ok bool
	var pair string
	var taxids *[]uint32
	for child, parent := range tree {
		name = strings.ToLower(names[child])
		pname = strings.ToLower(names[parent])

		if _n2i, ok = name2parent2taxid[name]; !ok {
			name2parent2taxid[name] = map[string]uint32{pname: child}
		} else {
			if _, ok = _n2i[pname]; ok {
				// log.Warningf("ambigous name pair: (%s, %s). TaxIds: %d, %d", _name, taxid2name[taxid2taxon[taxid].Parent], _n2i[pname], taxid)
				pair = name + "__" + pname
				if _, ok = ambigous[pair]; !ok {
					ambigous[pair] = []uint32{_n2i[pname], child}
				} else {
					ambigous[pair] = append(ambigous[pair], child)
				}
			} else {
				_n2i[pname] = child
			}
		}

		if taxids, ok = name2taxids[name]; !ok {
			name2taxids[name] = &[]uint32{child}
		} else {
			*taxids = append(*taxids, child)
		}
	}

	if config.Verbose {
		log.Infof("created links: child name -> parent name -> taxid")
	}
	return name2parent2taxid, name2taxids, ambigous
}

var poolStringsN16 = &sync.Pool{New: func() interface{} {
	return make([]string, 0, 16)
}}

var poolUint32N16 = &sync.Pool{New: func() interface{} {
	return make([]uint32, 0, 16)
}}

// only for reformat.
// remember to recyle return values
func queryNamesRanksTaxids(
	tree map[uint32]uint32,
	ranks map[uint32]string,
	names map[uint32]string,
	delnodes map[uint32]struct{},
	merged map[uint32]uint32,
	id uint32,
) ([]string, []string, []uint32, bool) {

	lineage := poolStringsN16.Get().([]string)
	lineageInRank := poolStringsN16.Get().([]string)
	lineageInTaxid := poolUint32N16.Get().([]uint32)

	var child, parent, newtaxid uint32
	var ok bool
	child = id
	for {
		parent, ok = tree[child]
		if !ok { // taxid not found
			// check if it was deleted
			if _, ok = delnodes[child]; ok {
				// log
				log.Warningf("taxid %d was deleted", child)
				id = 0
				return nil, nil, nil, false
			}
			// check if it was merged
			if newtaxid, ok = merged[child]; ok {
				// log
				log.Warningf("taxid %d was merged into %d", child, newtaxid)
				child = newtaxid
				parent = tree[child]
				id = child
			} else {
				id = 0
				log.Warningf("taxid %d not found", child)
				return nil, nil, nil, false
			}
		}

		lineage = append(lineage, names[child])
		lineageInRank = append(lineageInRank, ranks[child])
		lineageInTaxid = append(lineageInTaxid, child)

		if parent == 1 {
			break
		}

		child = parent
	}

	stringutil.ReverseStringSliceInplace(lineage)
	stringutil.ReverseStringSliceInplace(lineageInRank)
	reverseUint32s(lineageInTaxid)

	return lineage, lineageInRank, lineageInTaxid, true
}
