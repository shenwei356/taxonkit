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
	"strings"
	"sync"

	"github.com/shenwei356/breader"
)

// Taxon represents a taxonomic node
type Taxon struct {
	Taxid  int32
	Parent int32
	Name   string
	Rank   string
}

// ChunkSize is the chunk Size for breader
var ChunkSize = 100

var poolStringSliceL8 = &sync.Pool{New: func() interface{} {
	return make([]string, 8)
}}
var nameParseFunc = func(line string) (interface{}, bool, error) {
	// items := strings.SplitN(line, "\t", 8)
	items := poolStringSliceL8.Get().([]string)
	stringSplitN(line, "\t", 8, &items)
	if len(items) < 7 {
		return nil, false, nil
	}
	if items[6] != "scientific name" {
		return nil, false, nil
	}
	id, e := strconv.Atoi(items[0])
	if e != nil {
		return nil, false, e
	}
	taxon := Taxon{Taxid: int32(id), Name: items[2]}
	poolStringSliceL8.Put(items)
	return taxon, true, nil
}

var nameParseFunc2 = func(line string) (interface{}, bool, error) {
	items := strings.SplitN(line, "\t", 8)
	if len(items) < 7 {
		return nil, false, nil
	}
	if !(items[6] == "scientific name" || items[6] == "synonym") {
		return nil, false, nil
	}
	id, e := strconv.Atoi(items[0])
	if e != nil {
		return nil, false, e
	}
	return Taxon{Taxid: int32(id), Name: items[2]}, true, nil
}

// some taxons from different rank may have same names,
// so we use name and name of its parent to point to the right taxid.
func getName2Parent2Taxid(
	fileNodes string,
	fileNames string,
	bufferSize int,
	chunkSize int,
) (
	taxid2taxon map[int32]*Taxon,
	name2parent2taxid map[string]map[string]int32,
	name2taxid map[string]int32,
) {

	taxid2taxon = make(map[int32]*Taxon, mapInitialSize)

	name2parent2taxid = make(map[string]map[string]int32, mapInitialSize)

	name2taxid = make(map[string]int32, mapInitialSize) // not accurate

	reader, err := breader.NewBufferedReader(fileNodes, bufferSize, chunkSize, taxonParseFunc)
	checkError(err)

	// taxid -> rank, taxid -> parentid
	var taxon Taxon
	var data interface{}
	var ok bool
	for chunk := range reader.Ch {
		for _, data = range chunk.Data {
			taxon = data.(Taxon)
			// clone
			taxid2taxon[taxon.Taxid] = &Taxon{
				Taxid:  taxon.Taxid,
				Parent: taxon.Parent,
				// Name:   taxon.Name,
				Rank: taxon.Rank,
			}
		}
	}

	// taxid -> name
	reader, err = breader.NewBufferedReader(fileNames, bufferSize, chunkSize, nameParseFunc)
	checkError(err)
	for chunk := range reader.Ch {
		checkError(chunk.Err)
		for _, data = range chunk.Data {
			taxon = data.(Taxon)
			taxid2taxon[taxon.Taxid].Name = taxon.Name
		}
	}

	// name -> parent-name -> taxid
	var name, pname string
	var _n2i map[string]int32
	for taxid, taxon := range taxid2taxon {
		name = strings.ToLower(taxon.Name)
		pname = strings.ToLower(taxid2taxon[taxid2taxon[taxid].Parent].Name)
		if _n2i, ok = name2parent2taxid[name]; !ok {
			name2parent2taxid[name] = map[string]int32{pname: taxid}
		} else {
			_n2i[pname] = taxid
		}

		name2taxid[name] = taxid
	}
	return
}

var poolStringSliceL6 = &sync.Pool{New: func() interface{} {
	return make([]string, 6)
}}
var taxonParseFunc = func(line string) (interface{}, bool, error) {
	// items := strings.SplitN(line, "\t", 6)
	items := poolStringSliceL6.Get().([]string)
	stringSplitN(line, "\t", 6, &items)
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

	taxon := Taxon{Taxid: int32(child), Parent: int32(parent), Rank: items[4]}
	poolStringSliceL6.Put(items)
	return taxon, true, nil
}

// [from, to]
type mergedNodes [2]int32

var mergedParseFunc = func(line string) (interface{}, bool, error) {
	items := strings.SplitN(line, "\t", 4)
	if len(items) < 4 {
		return nil, false, nil
	}
	from, e := strconv.Atoi(items[0])
	if e != nil {
		return nil, false, e
	}
	to, e := strconv.Atoi(items[2])
	if e != nil {
		return nil, false, e
	}
	return mergedNodes([2]int32{int32(from), int32(to)}), true, nil
}
