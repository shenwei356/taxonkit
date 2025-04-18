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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shenwei356/bio/taxdump"
	"github.com/shenwei356/util/pathutil"
)

type rankFilter struct {
	taxondb *taxdump.Taxonomy

	dbRanks   map[string]interface{}
	rankOrder map[string]int

	lower  string
	higher string
	equals []string

	oLower  int
	oHigher int
	oEquals map[int]interface{}

	limitLower  bool
	limitHigher bool
	limitEqual  bool

	noRanks    map[string]interface{}
	blackLists map[string]interface{}

	discardNorank   bool
	saveKnownNoRank bool

	cache map[uint32]bool
}

func loadTaxonomy(opt *Config, withRank bool) *taxdump.Taxonomy {

	if opt.Verbose {
		log.Infof("loading Taxonomy from: %s", opt.DataDir)
	}
	var t *taxdump.Taxonomy
	var err error
	if withRank {
		t, err = taxdump.NewTaxonomyWithRankFromNCBI(filepath.Join(opt.DataDir, "nodes.dmp"))
	} else {
		t, err = taxdump.NewTaxonomyFromNCBI(filepath.Join(opt.DataDir, "nodes.dmp"))
	}
	if err != nil {
		checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
	}
	if opt.Verbose {
		if withRank {
			log.Infof("%d nodes in %d ranks loaded", len(t.Nodes), len(t.Ranks))
		} else {
			log.Infof("%d nodes loaded", len(t.Nodes))
		}
	}

	var existed bool

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		existed, err = pathutil.Exists(filepath.Join(opt.DataDir, "delnodes.dmp"))
		if err != nil {
			checkError(fmt.Errorf("err on checking file merged.dmp: %s", err))
		}
		if existed {
			err = t.LoadDeletedNodesFromNCBI(filepath.Join(opt.DataDir, "delnodes.dmp"))
			if err != nil {
				checkError(fmt.Errorf("err on loading Taxonomy nodes: %s", err))
			}
		}
		if opt.Verbose {
			log.Infof("%d deleted nodes loaded", len(t.DelNodes))
		}
	}()

	go func() {
		defer wg.Done()
		existed, err = pathutil.Exists(filepath.Join(opt.DataDir, "merged.dmp"))
		if err != nil {
			checkError(fmt.Errorf("err on checking file merged.dmp: %s", err))
		}
		if existed {
			err = t.LoadMergedNodesFromNCBI(filepath.Join(opt.DataDir, "merged.dmp"))
			if err != nil {
				checkError(fmt.Errorf("err on loading Taxonomy merged nodes: %s", err))
			}
		}
		if opt.Verbose {
			log.Infof("%d merged nodes loaded", len(t.MergeNodes))
		}
	}()

	wg.Wait()

	t.CacheLCA()

	return t
}

func newRankFilter(taxondb *taxdump.Taxonomy, rankOrder map[string]int, noRanks map[string]interface{},
	lower string, higher string, equals []string, blackList []string, discardNorank bool, saveKnownNoRank bool) (*rankFilter, error) {

	if lower != "" && higher != "" {
		return nil, fmt.Errorf("higher and lower can't be simultaneous given")
	}

	blackListMap := make(map[string]interface{})
	for _, r := range blackList {
		blackListMap[r] = struct{}{}
	}
	dbRanks := taxondb.Ranks
	f := &rankFilter{
		taxondb:         taxondb,
		dbRanks:         dbRanks,
		rankOrder:       rankOrder,
		lower:           lower,
		higher:          higher,
		equals:          equals,
		noRanks:         noRanks,
		blackLists:      blackListMap,
		discardNorank:   discardNorank,
		saveKnownNoRank: saveKnownNoRank,
		cache:           make(map[uint32]bool, 1024),
	}

	var err error
	if lower != "" {
		f.oLower, err = getRankOrder(dbRanks, rankOrder, lower)
		if err != nil {
			return nil, err
		}
		f.limitLower = true
	}
	if higher != "" {
		f.oHigher, err = getRankOrder(dbRanks, rankOrder, higher)
		if err != nil {
			return nil, err
		}
		f.limitHigher = true
	}
	if len(equals) > 0 {
		f.oEquals = make(map[int]interface{}, len(equals))
		var oe int
		for _, equal := range equals {
			oe, err = getRankOrder(dbRanks, rankOrder, equal)
			if err != nil {
				return nil, err
			}
			f.oEquals[oe] = struct{}{}
		}
		f.limitEqual = true
	}
	return f, nil
}

func getRankOrder(dbRanks map[string]interface{}, rankOrder map[string]int, rank string) (int, error) {
	var ok bool
	if _, ok = rankOrder[rank]; !ok {
		return -1, fmt.Errorf("rank order not defined in rank file: %s", rank)
	}
	if _, ok = dbRanks[rank]; !ok {
		return -1, fmt.Errorf("rank order not found in taxonomy database: %s", rank)
	}

	return rankOrder[rank], nil
}

func (f *rankFilter) isPassed(taxid uint32) (bool, error) {
	rank := f.taxondb.Rank(taxid)
	if rank == "" {
		return false, nil
	}

	rank = strings.ToLower(rank)

	if v, ok := f.cache[taxid]; ok {
		return v, nil
	}

	if _, ok := f.blackLists[rank]; ok {
		f.cache[taxid] = false
		return false, nil
	}

	var isNoRank bool
	_, ok := f.noRanks[rank]
	if ok {
		if f.discardNorank {
			isNoRank = true
			if !f.saveKnownNoRank {
				f.cache[taxid] = false
				return false, nil
			}
		} else { // all nonrank will be outputted if !discardNorank
			f.cache[taxid] = true
			return true, nil
		}
	}

	// checking taxid
	if _, ok := f.taxondb.Nodes[taxid]; !ok {
		if _, ok = f.taxondb.DelNodes[taxid]; ok {
			log.Warningf("taxid %d was deleted", taxid)
			return false, nil
		} else if newtaxid, ok := f.taxondb.MergeNodes[taxid]; ok {
			log.Warningf("taxid %d was merged into %d", taxid, newtaxid)
			taxid = newtaxid
		} else {
			log.Warningf("taxid %d not found", taxid)
			return false, nil
		}
	}

	var pass bool

	if isNoRank && f.limitLower && f.saveKnownNoRank {
		nodes := f.taxondb.Nodes
		var _rank string
		var _ok bool
		var _order int

		parent := nodes[taxid]
		for {
			if parent == 1 {
				f.cache[taxid] = false
				return false, nil
			}

			_rank = f.taxondb.Rank(parent)
			_order, _ok = f.rankOrder[_rank]
			if _ok {
				pass = _order <= f.oLower

				f.cache[taxid] = pass
				return pass, nil
			}
			parent = nodes[parent]
		}
	}

	order := f.rankOrder[rank]
	// order, ok := f.rankOrder[rank]
	// if !ok {
	// 	return false, fmt.Errorf("rank order not defined in rank file: %s", rank)
	// }

	if f.limitEqual {
		if _, pass = f.oEquals[order]; pass {
			// pass = true
		} else if f.limitLower {
			pass = order < f.oLower
		} else if f.limitHigher {
			pass = order > f.oHigher
		} else {
			pass = false
		}
	} else if f.limitLower {
		pass = order < f.oLower
	} else if f.limitHigher {
		pass = order > f.oHigher
	} else {
		pass = true // no any filter
	}

	f.cache[taxid] = pass
	return pass, nil
}

func readRankOrderFromFile(file string) (map[string]int, map[string]interface{}, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("read rank order list from '%s': %s", file, err)
	}

	ranks := make([][]string, 0, 128)
	noranks := make(map[string]interface{}, 10)

	scanner := bufio.NewScanner(fh)
	var record, item string
	for scanner.Scan() {
		record = strings.TrimSpace(scanner.Text())
		if record == "" || record[0] == '#' {
			continue
		}

		items := make([]string, 0, 1)

		for _, item = range strings.Split(record, ",") {
			if len(item) == 0 {
				continue
			}
			item = strings.ToLower(strings.TrimSpace(item))

			if item[0] == '!' {
				noranks[item[1:]] = struct{}{}
			} else {
				items = append(items, item)
			}
		}

		if len(items) > 0 {
			ranks = append(ranks, items)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read rank order list from '%s': %s", file, err)
	}

	if len(ranks) == 0 {
		return nil, nil, fmt.Errorf("no ranks found in file: %s", file)
	}

	rankOrder := make(map[string]int, len(ranks))
	order := 1
	var ok bool
	var rank string
	for i := len(ranks) - 1; i >= 0; i-- {
		for _, rank = range ranks[i] {
			if _, ok = rankOrder[rank]; ok {
				return nil, nil, fmt.Errorf("duplicated rank: %s", ranks[i])
			}
			rankOrder[rank] = order
		}
		order++
	}
	return rankOrder, noranks, nil
}

func readRankOrder(opt Config, rankFile string) (map[string]int, map[string]interface{}, error) {
	if rankFile != "" {
		if opt.Verbose {
			log.Infof("read rank order from: %s", rankFile)
		}
		return readRankOrderFromFile(rankFile)
	}

	defaultRankFile := filepath.Join(opt.DataDir, defaultRanksFile)
	existed, err := pathutil.Exists(defaultRankFile)
	if err != nil {
		return nil, nil, fmt.Errorf("check default rank file: %s", defaultRankFile)
	}
	if !existed {
		if opt.Verbose {
			log.Infof("write default rank order to: %s", defaultRankFile)
		}
		err = writeDefaltRankOrderFile(defaultRankFile)
		if err != nil {
			return nil, nil, fmt.Errorf("write default rank file: %s", defaultRankFile)
		}
	}

	if opt.Verbose {
		log.Infof("read rank order from: %s", defaultRankFile)
	}
	return readRankOrderFromFile(defaultRankFile)
}

func writeDefaltRankOrderFile(file string) error {
	return ioutil.WriteFile(file, []byte(defaultRanksText), 0644)
}

const defaultRanksFile = "ranks.txt"
const defaultRanksText = `
# This file defines taxonomic rank order for taxdump/taxonkit.
# 
# Here'are the rules:
#     1. Blank lines or lines starting with "#" are ignored.
#     2. Ranks are in decending order and case ignored.
#     3. Ranks with same order should be in one line separated with comma (",", no space).
#     4. Ranks without order should be assigned a prefix symbol "!" for each rank.
# 
# Deault ranks reference from https://en.wikipedia.org/wiki/Taxonomic_rank ,
# and contains some ranks from NCIB Taxonomy database.
#

!no rank
!clade


life

acellular root,cellular root
domain,superkingdom,realm,empire

kingdom
subkingdom
infrakingdom
parvkingdom

superphylum,superdivision
phylum,division
subphylum,subdivision
infraphylum,infradivision
microphylum,microdivision

superclass
class
subclass
infraclass
parvclass

superlegion
legion
sublegion
infralegion

supercohort
cohort
subcohort
infracohort

gigaorder
magnorder,megaorder
grandorder,capaxorder
mirorder,hyperorder
superorder
# series
order
# parvorder
nanorder
hypoorder
minorder
suborder
infraorder
parvorder

# section
# subsection

gigafamily
megafamily
grandfamily
hyperfamily
superfamily
epifamily
# series
group
family
subfamily
infrafamily

supertribe
tribe
subtribe
infratribe

genus
subgenus
section
subsection
series
subseries


superspecies,species group
species subgroup
species

subspecies,forma specialis,pathovar

pathogroup,serogroup
biotype,serotype,genotype

variety,varietas,morph,aberration
subvariety,subvarietas,submorph,subaberration
form,forma
subform,subforma

strain
isolate
`
