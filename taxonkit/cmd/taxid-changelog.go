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
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/shenwei356/util/pathutil"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// taxidlogCmd represents the fx2tab command
var taxidlogCmd = &cobra.Command{
	Use:   "taxid-changelog",
	Short: "create taxid changelog from dump archive",
	Long: `create taxid changelog from dump archive

Steps

	mkdir -p archive;
	cd archive;

	# --------- download ---------

	# option 1
	# for fast network connection
	wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump_archive/taxdmp*.zip

	# option 2
	# for bad network connection like mine
	url=https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump_archive/
	wget $url -O - -o /dev/null \
		| grep taxdmp | perl -ne '/(taxdmp_.+?.zip)/; print "$1\n";' \
		| rush -j 2 -v url=$url 'axel -n 5 {url}/{}' \
			--immediate-output  -c -C download.rush

	# --------- unzip ---------

	ls taxdmp*.zip | rush -j 1 'unzip {} names.dmp nodes.dmp merged.dmp delnodes.dmp -d {@_(.+)\.}'

	# --------- create log ---------

	taxonkit taxid-changelog -i archive -o log.csv.gz

`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		archivePath := getFlagString(cmd, "archive")
		dirs := checkArchives(config, archivePath)
		createChangelog(config, archivePath, dirs)
	},
}

func init() {
	RootCmd.AddCommand(taxidlogCmd)
	taxidlogCmd.Flags().StringP("archive", "i", ";", "directory containing decompressed dumped archives")
}

type TaxidChangeCode int32

const (
	TaxidUnchanged TaxidChangeCode = iota
	TaxidNew
	TaxidDelete
	TaxidMerge
	TaxidAbsorb
	TaxidLineageChanged
)

func (c TaxidChangeCode) String() string {
	switch c {
	case TaxidUnchanged:
		return "UNCHANGE"
	case TaxidLineageChanged:
		return "LINEAGE_CHANGED"
	case TaxidNew:
		return "NEW"
	case TaxidDelete:
		return "DELETE"
	case TaxidMerge:
		return "MERGE"
	case TaxidAbsorb:
		return "ABSORB"
	}
	return "UNDEFINED TaxidChangeCode"
}

type TaxidChange struct {
	// Taxid         int32
	Version       string
	Lineage       []string
	LineageTaxids []int32
	Change        TaxidChangeCode
	ChangeValue   []int32
}

type TaxidChanges []TaxidChange

func (changes TaxidChanges) Len() int { return len(changes) }

func (changes TaxidChanges) Swap(i, j int) { changes[i], changes[j] = changes[j], changes[i] }

func (changes TaxidChanges) Less(i, j int) bool {
	return strings.Compare(changes[i].Version, changes[j].Version) < 0 || changes[i].Change < changes[j].Change
}

func (c TaxidChange) String() string {
	var buf bytes.Buffer
	var tmp []string
	var tid int32
	var i int

	// version
	buf.WriteString(c.Version)

	// change
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%s", c.Change))

	// change value
	buf.WriteString(",")
	switch c.Change {
	case TaxidUnchanged:
	case TaxidLineageChanged:
	case TaxidNew:
	case TaxidDelete:
	case TaxidMerge:
		buf.WriteString(fmt.Sprintf("%d", c.ChangeValue[0]))
	case TaxidAbsorb:
		tmp = make([]string, len(c.ChangeValue))
		for i, tid = range c.ChangeValue {
			tmp[i] = fmt.Sprintf("%d", tid)
		}
		buf.WriteString(strings.Join(tmp, ";"))
	}

	// lineage
	buf.WriteString(",")
	if c.Lineage != nil {
		buf.WriteString(strings.Join(c.Lineage, ";"))
	}

	// lineage-taxid
	buf.WriteString(",")
	if c.LineageTaxids != nil {
		tmp = make([]string, len(c.LineageTaxids))
		for i, tid = range c.LineageTaxids {
			tmp[i] = fmt.Sprintf("%d", tid)
		}
		buf.WriteString(strings.Join(tmp, ";"))
	}

	return buf.String()
}

func createChangelog(config Config, path string, dirs []string) {

	outfh, err := xopen.Wopen(config.OutFile)
	checkError(err)

	defer outfh.Close()

	writer := csv.NewWriter(outfh)
	defer func() {
		writer.Flush()
		checkError(writer.Error())
	}()

	// taxid -> change-code -> []changes
	data := make(map[int32]map[TaxidChangeCode][]TaxidChange, 10000)
	var ok bool
	for j, dir := range dirs {
		if j == 3 {
			break
		}
		if config.Verbose {
			log.Infof("parsing archive (%3d/%3d): %s", j+1, len(dirs), dir)
		}

		// -------------- checking newly added and lineage-changed taxids --------------

		if config.Verbose {
			log.Infof("  parsing names.dmp & nodes.dmp")
		}

		taxid2lineageTaxids, taxid2name, _ := getTaxid2Lineage(
			filepath.Join(path, dir, "nodes.dmp"), filepath.Join(path, dir, "names.dmp"),
			config.Threads, 10,
		)

		var i int
		var tid int32
		var prevChange TaxidChange
		for taxid, lineageTaxids := range taxid2lineageTaxids {
			// fmt.Printf("%d: %s\n", taxid, strings.Join(lineage, "; "))

			if _, ok = data[taxid]; !ok { // newly added
				data[taxid] = make(map[TaxidChangeCode][]TaxidChange, 1)

				if _, ok = data[taxid][TaxidNew]; !ok {
					data[taxid][TaxidNew] = make([]TaxidChange, 0, 1)
				}

				lineage := make([]string, len(lineageTaxids))
				for i, tid = range lineageTaxids {
					lineage[i] = taxid2name[tid]
				}

				data[taxid][TaxidNew] = append(data[taxid][TaxidNew], TaxidChange{
					Version:       dir,
					Lineage:       lineage,
					LineageTaxids: lineageTaxids,
					Change:        TaxidNew,
					ChangeValue:   nil,
				})
			} else { // existed
				// newly added
				if _, ok = data[taxid][TaxidNew]; !ok {
					// checkError(fmt.Errorf("%d should have a NEW record", taxid))
					// !!!!!!!
					// deleted nodes can be used again ...
					// !!!!!!!
					data[taxid] = make(map[TaxidChangeCode][]TaxidChange, 1)

					if _, ok = data[taxid][TaxidNew]; !ok {
						data[taxid][TaxidNew] = make([]TaxidChange, 0, 1)
					}

					lineage := make([]string, len(lineageTaxids))
					for i, tid = range lineageTaxids {
						lineage[i] = taxid2name[tid]
					}

					data[taxid][TaxidNew] = append(data[taxid][TaxidNew], TaxidChange{
						Version:       dir,
						Lineage:       lineage,
						LineageTaxids: lineageTaxids,
						Change:        TaxidNew,
						ChangeValue:   nil,
					})
				} else {
					// check if lineage changed

					if _, ok = data[taxid][TaxidLineageChanged]; !ok { // it had already been changed
						prevChange = data[taxid][TaxidNew][0] // check NEW
					} else { // first change
						prevChange = data[taxid][TaxidLineageChanged][len(data[taxid][TaxidLineageChanged])-1] // check NEW
					}

					if !linegeUnchange(lineageTaxids, prevChange.LineageTaxids) {
						if _, ok = data[taxid][TaxidLineageChanged]; !ok {
							data[taxid][TaxidLineageChanged] = make([]TaxidChange, 0, 1)
						}

						lineage := make([]string, len(lineageTaxids))
						for i, tid = range lineageTaxids {
							lineage[i] = taxid2name[tid]
						}
						data[taxid][TaxidLineageChanged] = append(data[taxid][TaxidLineageChanged], TaxidChange{
							Version:       dir,
							Lineage:       lineage,
							LineageTaxids: lineageTaxids,
							Change:        TaxidLineageChanged,
							ChangeValue:   nil,
						})
					}
					// unchanged
					// leave it alone
				}
			}
		}

		// -------------- checking deleted taxids --------------

		if config.Verbose {
			log.Infof("  parsing delnodes.dmp")
		}
		var changes []TaxidChange
		taxids := getDelnodes(filepath.Join(path, dir, "delnodes.dmp"), config.Threads, 10)
		for _, taxid := range taxids {
			if _, ok = data[taxid]; !ok {
				data[taxid] = make(map[TaxidChangeCode][]TaxidChange, 1)
			}
			if _, ok = data[taxid][TaxidDelete]; !ok {
				data[taxid][TaxidDelete] = make([]TaxidChange, 0, 1)

				// extract lineage in NEW / CHANGED record
				if changes, ok = data[taxid][TaxidLineageChanged]; ok {
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       dir,
						Lineage:       changes[len(changes)-1].Lineage,
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						Change:        TaxidDelete,
						ChangeValue:   nil,
					})
				} else if changes, ok = data[taxid][TaxidNew]; ok {
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       dir,
						Lineage:       changes[len(changes)-1].Lineage,
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						Change:        TaxidDelete,
						ChangeValue:   nil,
					})
				} else {
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       dir,
						Lineage:       nil,
						LineageTaxids: nil,
						Change:        TaxidDelete,
						ChangeValue:   nil,
					})
				}
			}
		}

		// -------------- checking merged taxids --------------

		if config.Verbose {
			log.Infof("  parsing merged.dmp")
		}
		merges := getMergedNodes(filepath.Join(path, dir, "merged.dmp"), config.Threads, 10)
		var from, to int32
		var toRecord bool
		for _, merge := range merges {
			from, to = merge[0], merge[1]

			if _, ok = data[from]; !ok {
				data[from] = make(map[TaxidChangeCode][]TaxidChange, 1)
			}

			toRecord = false
			if _, ok = data[from][TaxidMerge]; !ok {
				data[from][TaxidMerge] = make([]TaxidChange, 0, 1)

				toRecord = true
			} else {
				// check if merged to another different taxid
				prevChange = data[from][TaxidMerge][len(data[from][TaxidMerge])-1]
				if prevChange.ChangeValue[len(prevChange.ChangeValue)-1] != to {
					toRecord = true
				}
			}

			if toRecord {
				if changes, ok = data[from][TaxidLineageChanged]; ok {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       dir,
						Lineage:       changes[len(changes)-1].Lineage,
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				} else if changes, ok = data[from][TaxidNew]; ok {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       dir,
						Lineage:       changes[len(changes)-1].Lineage,
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				} else {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       dir,
						Lineage:       nil,
						LineageTaxids: nil,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				}

				if _, ok = data[to]; !ok {
					data[to] = make(map[TaxidChangeCode][]TaxidChange, 1)
				}
				if _, ok = data[to][TaxidAbsorb]; !ok {
					data[to][TaxidAbsorb] = make([]TaxidChange, 0, 1)
				}
				if len(data[to][TaxidAbsorb]) == 0 {
					data[to][TaxidAbsorb] = append(data[to][TaxidAbsorb], TaxidChange{
						Version:       dir,
						Lineage:       nil,
						LineageTaxids: taxid2lineageTaxids[to],
						Change:        TaxidAbsorb,
						ChangeValue:   []int32{from},
					})
				} else {
					data[to][TaxidAbsorb][0].ChangeValue = append(data[to][TaxidAbsorb][0].ChangeValue, from)
				}
			}
		}
	}

	// output

	header := strings.Split("taxid,version,change,change-value,lineage,lineage-taxids", ",")
	writer.Write(header)

	var changes, cs []TaxidChange
	var c TaxidChange
	var tmp, items []string
	var tid int32
	var i int
	for taxid, v := range data {
		changes = make([]TaxidChange, 0, 1)
		for _, cs = range v {
			changes = append(changes, cs...)
		}
		sort.Sort(TaxidChanges(changes))
		for _, c = range changes {
			items = make([]string, 0, len(header))

			// taxid
			items = append(items, fmt.Sprintf("%d", taxid))

			// version
			items = append(items, c.Version)

			// change
			items = append(items, fmt.Sprintf("%s", c.Change))

			// change value
			switch c.Change {
			case TaxidUnchanged:
				items = append(items, "")
			case TaxidLineageChanged:
				items = append(items, "")
			case TaxidNew:
				items = append(items, "")
			case TaxidDelete:
				items = append(items, "")
			case TaxidMerge:
				items = append(items, fmt.Sprintf("%d", c.ChangeValue[0]))
			case TaxidAbsorb:
				tmp = make([]string, len(c.ChangeValue))
				for i, tid = range c.ChangeValue {
					tmp[i] = fmt.Sprintf("%d", tid)
				}
				items = append(items, strings.Join(tmp, ";"))
			}

			// lineage
			if c.Lineage != nil {
				items = append(items, strings.Join(c.Lineage, ";"))
			} else {
				items = append(items, "")
			}

			// lineage-taxid
			if c.LineageTaxids != nil {
				tmp = make([]string, len(c.LineageTaxids))
				for i, tid = range c.LineageTaxids {
					tmp[i] = fmt.Sprintf("%d", tid)
				}
				items = append(items, strings.Join(tmp, ";"))
			} else {
				items = append(items, "")
			}

			writer.Write(items)
		}
	}
}

func checkArchives(config Config, path string) []string {
	checkFile(path)

	_, err := ioutil.ReadFile(path)
	// it's a file
	if err == nil {
		checkError(fmt.Errorf("archive path should be a directory: %s", path))
	}

	// it's a directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Warning(err)
	}

	var filename string
	dirs := make([]string, 0, len(files))
	for _, file := range files {
		filename = file.Name()

		if filename[0] == '.' {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, filename)
		}
	}
	if len(dirs) == 0 {
		checkError(fmt.Errorf("no unzipeed directories found in: %s", path))
	}
	for _, dir := range dirs {
		checkFile(filepath.Join(path, dir, "names.dmp"))
		checkFile(filepath.Join(path, dir, "nodes.dmp"))
		checkFile(filepath.Join(path, dir, "delnodes.dmp"))
		checkFile(filepath.Join(path, dir, "merged.dmp"))
	}

	if config.Verbose {
		log.Infof("%d archieves found in %s", len(dirs), path)
	}

	return dirs
}

func checkFile(file string) {
	if exists, err := pathutil.Exists(file); err != nil {
		checkError(fmt.Errorf("checking %s: %s", file, err))
	} else if !exists {
		checkError(fmt.Errorf("path not found: %s", file))
	}
}

func linegeUnchange(a, b []int32) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
