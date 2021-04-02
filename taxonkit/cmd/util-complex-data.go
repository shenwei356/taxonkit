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

// Taxon represents a taxonomic node
type Taxon struct {
	Taxid  uint32
	Parent uint32
	Name   string
	Rank   string
}

// some taxons from different rank may have same names,
// so we use name and name of its parent to point to the right taxid.
func getName2Parent2Taxid(config Config) (
	map[uint32]*Taxon,
	map[string]map[string]uint32,
	map[string]uint32,
	map[string][]uint32, // save ambigous child-parent pairs
) {

	taxid2taxon := make(map[uint32]*Taxon, mapInitialSize)

	// ------------------------------------------------------

	var wg sync.WaitGroup
	wg.Add(2)

	// taxid -> name

	var taxid2name map[uint32]string
	go func() {

		if config.Verbose {
			log.Infof("parsing names file: %s", config.NamesFile)
		}
		taxid2name = getTaxonNames(config.NamesFile)
		if config.Verbose {
			log.Infof("%d names parsed", len(taxid2name))
		}
		wg.Done()
	}()

	// taxid -> taxon

	go func() {
		if config.Verbose {
			log.Infof("parsing nodes file: %s", config.NodesFile)
		}
		fh, err := xopen.Ropen(config.NodesFile)
		checkError(err)

		items := make([]string, 6)
		scanner := bufio.NewScanner(fh)
		var _child, _parent int
		var child, parent uint32
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
			child, parent, rank = uint32(_child), uint32(_parent), items[4]

			// ----------------------------------

			taxid2taxon[child] = &Taxon{
				Taxid:  child,
				Parent: parent,
				// Name:   "",
				Rank: rank,
			}
		}
		if err := scanner.Err(); err != nil {
			checkError(err)
		}

		checkError(fh.Close())

		if config.Verbose {
			log.Infof("%d nodes parsed", len(taxid2taxon))
		}

		wg.Done()
	}()

	wg.Wait()

	// ------------------------------------------------------

	if config.Verbose {
		log.Infof(" create links: child name -> parent name -> taxid")
	}

	name2parent2taxid := make(map[string]map[string]uint32, len(taxid2taxon))

	name2taxid := make(map[string]uint32, len(taxid2taxon)) // not accurate

	ambigous := make(map[string][]uint32, 128)

	// name -> parent-name -> taxid

	var _name, name, pname string
	var _n2i map[string]uint32
	var ok bool
	var pair string
	for taxid, taxon := range taxid2taxon {
		_name = taxid2name[taxid]
		taxon.Name = _name

		name = strings.ToLower(_name)
		pname = strings.ToLower(taxid2name[taxid2taxon[taxid].Parent])

		if _n2i, ok = name2parent2taxid[name]; !ok {
			name2parent2taxid[name] = map[string]uint32{pname: taxid}
		} else {
			if _, ok = _n2i[pname]; ok {
				// log.Warningf("ambigous name pair: (%s, %s). TaxIds: %d, %d", _name, taxid2name[taxid2taxon[taxid].Parent], _n2i[pname], taxid)
				pair = name + pname
				if _, ok = ambigous[pair]; !ok {
					ambigous[pair] = []uint32{_n2i[pname], taxid}
				} else {
					ambigous[pair] = append(ambigous[pair], taxid)
				}
			} else {
				_n2i[pname] = taxid
			}
		}

		name2taxid[name] = taxid
	}

	if config.Verbose {
		log.Infof("created links: child name -> parent name -> taxid")
	}
	return taxid2taxon, name2parent2taxid, name2taxid, ambigous
}

// this function is only used in reformat command, with
func lineageFromTaxid2Taxon(taxid2taxon map[uint32]*Taxon, id uint32, delimiter string) string {
	lineage := make([]string, 0, 16)

	var taxon *Taxon
	var child, parent uint32
	child = uint32(id)

	for {
		taxon = taxid2taxon[child]
		parent = taxon.Parent

		lineage = append(lineage, taxon.Name)

		if parent == 1 {
			break
		}
		child = parent
	}
	child = uint32(id)

	return strings.Join(stringutil.ReverseStringSlice(lineage), delimiter)
}
