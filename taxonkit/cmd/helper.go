// Copyright © 2016-2020 Wei Shen <shenwei356@gmail.com>
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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shenwei356/breader"
	"github.com/shenwei356/util/pathutil"
	"github.com/spf13/cobra"
)

// VERSION of taxonkit
const VERSION = "0.6.0"

// Config is the struct containing all global flags
type Config struct {
	Threads      int
	OutFile      string
	DataDir      string
	NodesFile    string
	NamesFile    string
	DelNodesFile string
	MergedFile   string
	Verbose      bool
	LineBuffered bool
}

func getConfigs(cmd *cobra.Command) Config {
	var val, dataDir string
	if val = os.Getenv("TAXONKIT_DB"); val != "" {
		dataDir = val
	} else {
		dataDir = getFlagString(cmd, "data-dir")
	}

	existed, err := pathutil.DirExists(dataDir)
	checkError(err)
	if !existed {
		log.Errorf(`data directory not created. please download and decompress ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz, and copy "names.dmp", "nodes.dmp", "delnodes.dmp", and "merged.dmp" to %s`, dataDir)
	}

	return Config{Threads: getFlagPositiveInt(cmd, "threads"),
		OutFile:      getFlagString(cmd, "out-file"),
		DataDir:      dataDir,
		NodesFile:    filepath.Join(dataDir, "nodes.dmp"),
		NamesFile:    filepath.Join(dataDir, "names.dmp"),
		DelNodesFile: filepath.Join(dataDir, "delnodes.dmp"),
		MergedFile:   filepath.Join(dataDir, "merged.dmp"),
		Verbose:      getFlagBool(cmd, "verbose"),
		LineBuffered: getFlagBool(cmd, "line-buffered"),
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

// Taxon represents a taxonomic node
type Taxon struct {
	Taxid  int32
	Parent int32
	Name   string
	Rank   string
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
	return Taxon{Taxid: int32(id), Name: items[2]}, true, nil
}

var nameParseFunc2 = func(line string) (interface{}, bool, error) {
	items := strings.Split(line, "\t")
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

func getTaxonNames(file string, bufferSize int, chunkSize int) map[int32]string {
	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, nameParseFunc)
	checkError(err)

	var taxon Taxon
	taxid2name := make(map[int32]string)
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			taxon = data.(Taxon)
			taxid2name[taxon.Taxid] = taxon.Name
		}
	}
	return taxid2name
}

func getTaxonName2Taxids(file string, limit2SciName bool, bufferSize int, chunkSize int) map[string][]int32 {
	var reader *breader.BufferedReader
	var err error
	if limit2SciName {
		reader, err = breader.NewBufferedReader(file, bufferSize, chunkSize, nameParseFunc)
	} else {
		reader, err = breader.NewBufferedReader(file, bufferSize, chunkSize, nameParseFunc2)
	}
	checkError(err)

	var taxon Taxon
	name2taxids := make(map[string][]int32)
	var ok bool
	var name string
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			taxon = data.(Taxon)
			name = strings.ToLower(taxon.Name)
			if _, ok = name2taxids[name]; !ok {
				name2taxids[name] = make([]int32, 0, 1)
			}
			name2taxids[name] = append(name2taxids[name], taxon.Taxid)
		}
	}
	return name2taxids
}

func getName2Parent2Taxid(fileNodes string,
	fileNames string,
	bufferSize int,
	chunkSize int) (
	taxid2taxon map[int32]*Taxon,
	name2parent2taxid map[string]map[string]int32,
	name2taxid map[string]int32) {

	// taxid -> rank, taxid -> parentid
	reader, err := breader.NewBufferedReader(fileNodes, bufferSize, chunkSize, taxonParseFunc)
	checkError(err)

	// 2000000 comes from $(wc -l nodes.dmp)
	taxid2taxon = make(map[int32]*Taxon, 2000000)
	name2parent2taxid = make(map[string]map[string]int32, 2000000)
	name2taxid = make(map[string]int32, 2000000) // not accurate

	// taxid -> rank, taxid -> parentid
	var taxon Taxon
	var data interface{}
	var ok bool
	for chunk := range reader.Ch {
		checkError(chunk.Err)
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
		for _, data := range chunk.Data {
			taxon = data.(Taxon)
			taxid2taxon[taxon.Taxid].Name = taxon.Name
		}
	}

	// name -> parent-name -> taxid
	var name, pname string
	for taxid, taxon := range taxid2taxon {
		name = strings.ToLower(taxon.Name)
		pname = strings.ToLower(taxid2taxon[taxid2taxon[taxid].Parent].Name)
		if _, ok = name2parent2taxid[name]; !ok {
			name2parent2taxid[name] = make(map[string]int32, 1)
		}
		name2parent2taxid[name][pname] = taxid
		name2taxid[name] = taxid
	}
	return
}

func getTaxid2Lineage(fileNodes string, fileNames string, bufferSize int, chunkSize int) (map[int32][]int32, map[int32]string, map[int32]string) {
	var names map[int32]string
	names = getTaxonNames(fileNames, bufferSize, chunkSize)

	reader, err := breader.NewBufferedReader(fileNodes, bufferSize, chunkSize, taxonParseFunc)
	checkError(err)

	tree := make(map[int32]int32)
	ranks := make(map[int32]string)
	var taxon Taxon
	var child, parent int32
	var n int64
	var data interface{}
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data = range chunk.Data {
			taxon = data.(Taxon)
			child, parent = taxon.Taxid, taxon.Parent

			tree[child] = parent
			ranks[child] = taxon.Rank
			n++
		}
	}

	taxid2lineageTaxids := make(map[int32][]int32, 10000)

	var ok bool
	var i, j int
	for taxid := range tree {
		lineageTaxids := []int32{}
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
		// reverse lineageTaxids
		for i = len(lineageTaxids)/2 - 1; i >= 0; i-- {
			j = len(lineageTaxids) - 1 - i
			lineageTaxids[i], lineageTaxids[j] = lineageTaxids[j], lineageTaxids[i]
		}
		taxid2lineageTaxids[taxid] = lineageTaxids
	}
	return taxid2lineageTaxids, names, ranks
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
	return Taxon{Taxid: int32(child), Parent: int32(parent), Rank: items[4]}, true, nil
}

type delNode int32

var delnodesParseFunc = func(line string) (interface{}, bool, error) {
	items := strings.Split(line, "\t")
	if len(items) < 2 {
		return nil, false, nil
	}
	id, e := strconv.Atoi(items[0])
	if e != nil {
		return nil, false, e
	}
	return delNode(id), true, nil
}

func getDelnodes(file string, bufferSize int, chunkSize int) []int32 {
	taxids := make([]int32, 0, 100000)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("delnodes file not found: %s, deleted taxids will not be checked", file)
		return taxids
	}

	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, delnodesParseFunc)
	checkError(err)
	var taxid delNode
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			taxid = data.(delNode)
			taxids = append(taxids, int32(taxid))
		}
	}
	return taxids
}

func getDelnodesMap(file string, bufferSize int, chunkSize int) map[int32]struct{} {
	taxids := make(map[int32]struct{}, 100000)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("delnodes file not found: %s, deleted taxids will not be checked", file)
		return taxids
	}

	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, delnodesParseFunc)
	checkError(err)

	var taxid delNode
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			taxid = data.(delNode)
			taxids[int32(taxid)] = struct{}{}
		}
	}

	return taxids
}

type mergedNodes [2]int32

var mergedParseFunc = func(line string) (interface{}, bool, error) {
	items := strings.Split(line, "\t")
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

func getMergedNodes(file string, bufferSize int, chunkSize int) [][2]int32 {
	merges := make([][2]int32, 0, 10000)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("merged file not found: %s, merged taxids will not be checked", file)
		return merges
	}

	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, mergedParseFunc)
	checkError(err)

	var merge mergedNodes
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			merge = data.(mergedNodes)
			merges = append(merges, [2]int32(merge))
		}
	}

	return merges
}

func getMergedNodesMap(file string, bufferSize int, chunkSize int) map[int32]int32 {
	merges := make(map[int32]int32, 10000)

	existed, err := pathutil.Exists(file)
	if err != nil {
		checkError(err)
	}
	if !existed {
		log.Warningf("merged file not found: %s, merged taxids will not be checked", file)
		return merges
	}

	reader, err := breader.NewBufferedReader(file, bufferSize, chunkSize, mergedParseFunc)
	checkError(err)

	var merge mergedNodes
	for chunk := range reader.Ch {
		checkError(chunk.Err)

		for _, data := range chunk.Data {
			merge = data.(mergedNodes)
			merges[merge[0]] = merge[1]
		}
	}

	return merges
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
