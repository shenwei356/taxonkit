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
	"strconv"

	"github.com/shenwei356/bio/taxdump"
)

type Target struct {
	Taxid     uint32
	Abundance float64

	Rank          string
	TaxonName     string
	LineageNames  []string
	LineageTaxids []string

	CompleteLineageNames  []string
	CompleteLineageTaxids []uint32
}

func (t *Target) AddTaxonomy(taxdb *taxdump.Taxonomy, showRanksMap map[string]interface{}, taxid uint32) bool {
	var ok bool
	t.Taxid, ok = taxdb.TaxId(taxid)
	if !ok {
		return ok
	}
	t.Rank = taxdb.Rank(taxid)
	t.TaxonName = taxdb.Name(taxid)

	_taxids := taxdb.LineageTaxIds(taxid)

	t.CompleteLineageTaxids = _taxids
	t.CompleteLineageNames = taxdb.LineageNames(taxid)

	var _taxids2 []uint32
	if len(showRanksMap) > 0 {
		_taxids2 = make([]uint32, 0, len(_taxids))
		for _, _taxid := range _taxids {
			if _, ok = showRanksMap[taxdb.Rank(_taxid)]; ok {
				_taxids2 = append(_taxids2, _taxid)
			}
		}
		_taxids = _taxids2
	}

	t.LineageTaxids = make([]string, len(_taxids))
	for i, _taxid := range _taxids {
		t.LineageTaxids[i] = strconv.Itoa(int(_taxid))
	}

	t.LineageNames = make([]string, len(_taxids))
	for i, _taxid := range _taxids {
		t.LineageNames[i] = taxdb.Names[_taxid]
	}

	return true
}

type Targets []*Target

func (t Targets) Len() int { return len(t) }
func (t Targets) Less(i, j int) bool {
	return t[i].Abundance > t[j].Abundance
}
func (t Targets) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type ProfileNode struct {
	Taxid         uint32
	Rank          string
	TaxonName     string
	LineageNames  []string // complete lineage
	LineageTaxids []uint32
	Abundance     float64
}

func generateProfile(taxdb *taxdump.Taxonomy, targets []*Target) map[uint32]*ProfileNode {

	profile := make(map[uint32]*ProfileNode, len(targets))

	for _, target := range targets {
		for _, taxid := range target.CompleteLineageTaxids {
			if node, ok := profile[taxid]; !ok {
				profile[taxid] = &ProfileNode{
					Taxid:         taxid,
					Rank:          taxdb.Rank(taxid),
					TaxonName:     taxdb.Names[taxid],
					LineageNames:  taxdb.LineageNames(taxid),
					LineageTaxids: taxdb.LineageTaxIds(taxid),

					Abundance: target.Abundance,
				}
			} else {
				node.Abundance += target.Abundance
			}
		}
	}

	return profile
}
