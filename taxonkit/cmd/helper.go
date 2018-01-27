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
	"regexp"
	"strconv"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/spf13/cobra"
)

// VERSION of csvtk
const VERSION = "0.2.4-dev3"

// Config is the struct containing all global flags
type Config struct {
	Threads   int
	OutFile   string
	NodesFile string
	NamesFile string
	Verbose   bool
}

func getConfigs(cmd *cobra.Command) Config {
	return Config{
		Threads:   getFlagPositiveInt(cmd, "threads"),
		OutFile:   getFlagString(cmd, "out-file"),
		NodesFile: getFlagString(cmd, "nodes-file"),
		NamesFile: getFlagString(cmd, "names-file"),
		Verbose:   getFlagBool(cmd, "verbose"),
	}
}

var reTaxIDs = regexp.MustCompile(`(\d+)(,\d+)*`)

func getFlagTaxonIDs(cmd *cobra.Command, flag string) []int {
	s, err := cmd.Flags().GetString(flag)
	checkError(err)
	if s == "" {
		checkError(fmt.Errorf("flag --%s needed", flag))
	}
	if !reTaxIDs.MatchString(s) {
		checkError(fmt.Errorf("invalid value of flag %s. comma-separated integers needed", flag))
	}
	idStrSlice := strings.Split(s, ",")
	ids := make([]int, len(idStrSlice))
	var id int
	for i, s := range idStrSlice {
		id, _ = strconv.Atoi(s)
		ids[i] = id
	}
	return ids
}

type Taxid2Name struct {
	id   int32
	name string
}

var nameParseFunc = func(line string) (interface{}, bool, error) {
	items := strings.Split(line, "\t")
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
	return Taxid2Name{int32(id), items[2]}, true, nil
}

func getTaxonNames(file string, bufferSize int, chunkSize int) map[int32]string {
	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, nameParseFunc)
	checkError(err)

	var rel Taxid2Name
	m := make(map[int32]string)
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			rel = data.(Taxid2Name)
			m[rel.id] = rel.name
		}
	}
	return m
}

func getTaxonNames2Taxid(file string, bufferSize int, chunkSize int) map[string][]int32 {
	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, nameParseFunc)
	checkError(err)

	var rel Taxid2Name
	m := make(map[string][]int32)
	var ok bool
	var name string
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			rel = data.(Taxid2Name)
			name = strings.ToLower(rel.name)
			if _, ok = m[name]; !ok {
				m[name] = make([]int32, 0, 1)
			}
			m[name] = append(m[name], rel.id)
		}
	}
	return m
}

func getTaxonNames2TaxidAndRank(fileNodes string,
	fileNames string,
	bufferSize int,
	chunkSize int) (taxid2parent map[int32]int32,
	taxid2rank map[int32]string,
	taxid2name map[int32]string,
	name2parent2taxid map[string]map[string]int32,
	name2taxid map[string]int32,
	fuzzyNames map[string]struct{}) {
	reader, err := breader.NewBufferedReader(fileNodes, bufferSize, chunkSize, taxonParseFunc)
	checkError(err)

	taxid2rank = make(map[int32]string, 10000)
	taxid2parent = make(map[int32]int32, 10000)
	var info taxonInfo
	var data interface{}
	var name string
	var ok bool
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data = range chunk.Data {
			info = data.(taxonInfo)
			taxid2rank[info.child] = info.rank
			taxid2parent[info.child] = info.parent
		}
	}

	reader, err = breader.NewBufferedReader(fileNames, bufferSize, chunkSize, nameParseFunc)
	checkError(err)

	var rel Taxid2Name
	taxid2name = make(map[int32]string, 10000)
	name2parent2taxid = make(map[string]map[string]int32, 10000)
	name2taxid = make(map[string]int32, 10000) // not accurate
	fuzzyNames = make(map[string]struct{}, 1000)
	var pname string

	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			rel = data.(Taxid2Name)
			taxid2name[rel.id] = rel.name
		}
	}

	// have to read names.dmp twice.
	reader, err = breader.NewBufferedReader(fileNames, bufferSize, chunkSize, nameParseFunc)
	checkError(err)

	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			rel = data.(Taxid2Name)

			name = strings.ToLower(rel.name)
			pname = strings.ToLower(taxid2name[taxid2parent[rel.id]])
			if _, ok = name2parent2taxid[name]; !ok {
				name2parent2taxid[name] = make(map[string]int32, 1)
			} else {
				fuzzyNames[name] = struct{}{}
			}
			name2parent2taxid[name][pname] = rel.id
			name2taxid[name] = rel.id
		}
	}

	return
}

type taxonInfo struct {
	child, parent int32
	rank          string
}

var taxonParseFunc = func(line string) (interface{}, bool, error) {
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
	return taxonInfo{int32(child), int32(parent), items[4]}, true, nil
}

var rank2symbol = map[string]string{
	"superkingdom": "k",
	"phylum":       "p",
	"class":        "c",
	"order":        "o",
	"family":       "f",
	"genus":        "g",
	"species":      "s",
	"subspecies":   "S",
}

var symbol2rank = map[string]string{
	"k": "superkingdom",
	"p": "phylum",
	"c": "class",
	"o": "order",
	"f": "family",
	"g": "genus",
	"s": "species",
	"S": "subspecies",
}
var symbol2weight = map[string]float32{
	"k": 1,
	"p": 2,
	"c": 3,
	"o": 4,
	"f": 5,
	"g": 6,
	"s": 7,
	"S": 8,
}

var reRankPlaceHolder = regexp.MustCompile(`\{(\w)\}`)

var reRankPlaceHolders = map[string]*regexp.Regexp{
	"k": regexp.MustCompile(`\{k\}`),
	"p": regexp.MustCompile(`\{p\}`),
	"c": regexp.MustCompile(`\{c\}`),
	"o": regexp.MustCompile(`\{o\}`),
	"f": regexp.MustCompile(`\{f\}`),
	"g": regexp.MustCompile(`\{g\}`),
	"s": regexp.MustCompile(`\{s\}`),
	"S": regexp.MustCompile(`\{S\}`),
}

const norank = "no rank"

func isStdin(file string) bool {
	return file == "-"
}
