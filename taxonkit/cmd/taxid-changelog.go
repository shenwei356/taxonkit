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

Steps:

    # dependencies:
    # 	rush - https://github.com/shenwei356/rush/

    mkdir -p archive; cd archive;

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

    cd ..
    taxonkit taxid-changelog -i archive -o log.csv.gz --verbose

Output format (CSV):

    # fields        comments
    taxid           # taxid
    version         # version / time of archive, e.g, 2019-07-01
    change          # change, values:
                    #   NEW             newly added
                    #   DELETE          deleted
                    #   MERGE           merged into another taxid
                    #   ABSORB          other taxids merged into this one
                    #   L_CHANGE_LIN    lineage taxids remain but lineage remain
                    #   L_CHANGE_TAX    lineage taxids changed
                    #   L_CHANGE_LEN    lineage length changed
    change-value    # variable values for changes: 
                    #   1) empty for NEW, DELETE, L_CHANGE_LIN, L_CHANGE_TAX, L_CHANGE_LEN
                    #   2) new taxid for MERGE,
                    #   3) merged taxids for ABSORB
    name            # scientific name
    rank            # rank
    lineage         # full lineage of the taxid
    lineage-taxids  # taxids of the lineage

    # you can use csvtk to investigate them. e.g.,
    csvtk grep -f taxid -p 1390515 log.csv.gz
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfigs(cmd)
		runtime.GOMAXPROCS(config.Threads)

		archivePath := getFlagString(cmd, "archive")
		if archivePath == "" {
			checkError(fmt.Errorf("flag -i/--archive needed"))
		}

		dirs := checkArchives(config, archivePath)
		createChangelog(config, archivePath, dirs)
	},
}

func init() {
	RootCmd.AddCommand(taxidlogCmd)
	taxidlogCmd.Flags().StringP("archive", "i", "", "directory containing decompressed dumped archives")
}

type TaxidChangeCode int32

const (
	TaxidUnchanged TaxidChangeCode = iota
	TaxidNew
	TaxidDelete
	TaxidMerge
	TaxidAbsorb
	TaxidLineageChanged
	TaxidLineageChangedLin // lineage taxids remain, but lineage changed
	TaxidLineageChangedTax // lineage taxids changed
	TaxidLineageChangedLen // number of lineage taxids changed
)

func (c TaxidChangeCode) String() string {
	switch c {
	case TaxidUnchanged:
		return "UNCHANGE"
	case TaxidNew:
		return "NEW"
	case TaxidDelete:
		return "DELETE"
	case TaxidMerge:
		return "MERGE"
	case TaxidAbsorb:
		return "ABSORB"
	case TaxidLineageChanged:
		return "L_CHANGE"
	case TaxidLineageChangedLin:
		return "L_CHANGE_LIN"
	case TaxidLineageChangedTax:
		return "L_CHANGE_TAX"
	case TaxidLineageChangedLen:
		return "L_CHANGE_LEN"
	}
	return "UNDEFINED TaxidChangeCode"
}

func linegeChangeType(a, b []int32, taxid2names map[string]map[int32]string, va, vb string) TaxidChangeCode {
	if (a == nil) != (b == nil) {
		return TaxidLineageChangedLen
	}

	if len(a) != len(b) {
		return TaxidLineageChangedLen
	}

	for i, v := range a {
		if v != b[i] {
			return TaxidLineageChangedTax
		}
	}

	for i, v := range a {
		if taxid2names[va][v] != taxid2names[vb][b[i]] {
			return TaxidLineageChangedLin
		}
	}

	return TaxidUnchanged
}

type TaxidChange struct {
	Version       int16
	TaxidVersion  int16
	Change        TaxidChangeCode
	LineageTaxids []int32
	ChangeValue   []int32
}

type TaxidChanges []TaxidChange

func (changes TaxidChanges) Len() int { return len(changes) }

func (changes TaxidChanges) Swap(i, j int) { changes[i], changes[j] = changes[j], changes[i] }

func (changes TaxidChanges) Less(i, j int) bool {
	if changes[i].Version < changes[j].Version {
		return true
	}
	if changes[i].Version > changes[j].Version {
		return false
	}
	if changes[i].Change < changes[j].Change {
		return true
	}
	return false
}

func (c TaxidChange) String() string {
	var buf bytes.Buffer
	var tmp []string
	var tid int32
	var i int

	// version
	buf.WriteString(fmt.Sprintf("%d", c.Version))

	// taxid version
	buf.WriteString(fmt.Sprintf(",%d", c.TaxidVersion))

	// change
	buf.WriteString(",")
	buf.WriteString(fmt.Sprintf("%s", c.Change))

	// change value
	buf.WriteString(",")
	switch c.Change {
	case TaxidMerge:
		buf.WriteString(fmt.Sprintf("%d", c.ChangeValue[0]))
	case TaxidAbsorb:
		tmp = make([]string, len(c.ChangeValue))
		for i, tid = range c.ChangeValue {
			tmp[i] = fmt.Sprintf("%d", tid)
		}
		buf.WriteString(strings.Join(tmp, ";"))
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

	// version -> taxid -> name
	taxid2names := make(map[string]map[int32]string, len(dirs))

	// version -> taxid -> rank
	taxid2ranks := make(map[string]map[int32]string, len(dirs))

	// versions
	versions := dirs

	sort.Strings(dirs)

	var ok bool
	for version, dir := range dirs {
		if config.Verbose {
			log.Infof("parsing archive (%3d/%3d): %s", version+1, len(dirs), dir)
		}

		// -------------- checking newly added and lineage-changed taxids --------------

		if config.Verbose {
			log.Infof("  parsing names.dmp & nodes.dmp")
		}

		taxid2lineageTaxids, taxid2name, taxid2rank := getTaxid2Lineage(
			filepath.Join(path, dir, "nodes.dmp"), filepath.Join(path, dir, "names.dmp"),
			config.Threads, 10,
		)
		taxid2names[dir] = taxid2name
		taxid2ranks[dir] = taxid2rank

		var prevChange TaxidChange
		var changeCode TaxidChangeCode
		for taxid, lineageTaxids := range taxid2lineageTaxids {
			if _, ok = data[taxid]; !ok { // newly added
				data[taxid] = make(map[TaxidChangeCode][]TaxidChange, 1)

				if _, ok = data[taxid][TaxidNew]; !ok {
					data[taxid][TaxidNew] = make([]TaxidChange, 0, 1)
				}

				data[taxid][TaxidNew] = append(data[taxid][TaxidNew], TaxidChange{
					Version:       int16(version),
					LineageTaxids: lineageTaxids,
					TaxidVersion:  int16(version),
					Change:        TaxidNew,
					ChangeValue:   nil,
				})
			} else { // existed
				// newly added
				if _, ok = data[taxid][TaxidNew]; !ok {
					// !!!!!!!
					// deleted nodes can be re-used ...
					// e.g., 1390515 , deleted in 2014-08-01, added in 2014-09-01
					// !!!!!!!
					// log.Infof(fmt.Sprintf("%d, v%s, %s", taxid, versions[version], data[taxid]))

					if _, ok = data[taxid][TaxidNew]; !ok {
						data[taxid][TaxidNew] = make([]TaxidChange, 0, 1)
					}
					data[taxid][TaxidNew] = append(data[taxid][TaxidNew], TaxidChange{
						Version:       int16(version),
						LineageTaxids: lineageTaxids,
						TaxidVersion:  int16(version),
						Change:        TaxidNew,
						ChangeValue:   nil,
					})
				} else {
					// check if lineage changed

					if _, ok = data[taxid][TaxidLineageChanged]; !ok { // it had already been changed
						prevChange = data[taxid][TaxidNew][0] // check NEW
					} else { // last change
						prevChange = data[taxid][TaxidLineageChanged][len(data[taxid][TaxidLineageChanged])-1] // check NEW
					}

					changeCode = linegeChangeType(lineageTaxids, prevChange.LineageTaxids, taxid2names, dir, versions[prevChange.TaxidVersion])
					if changeCode > 0 { // changed
						if _, ok = data[taxid][TaxidLineageChanged]; !ok {
							data[taxid][TaxidLineageChanged] = make([]TaxidChange, 0, 1)
						}

						data[taxid][TaxidLineageChanged] = append(data[taxid][TaxidLineageChanged], TaxidChange{
							Version:       int16(version),
							LineageTaxids: lineageTaxids,
							TaxidVersion:  int16(version),
							Change:        changeCode,
							ChangeValue:   nil,
						})
					}
					// unchanged
					// leave it alone
				}
			}
		}

		var changes []TaxidChange

		// -------------- checking deleted taxids --------------

		if config.Verbose {
			log.Infof("  parsing delnodes.dmp")
		}

		delTaxids := getDelnodes(filepath.Join(path, dir, "delnodes.dmp"), config.Threads, 10)
		for _, taxid := range delTaxids {
			if _, ok = data[taxid]; !ok {
				data[taxid] = make(map[TaxidChangeCode][]TaxidChange, 1)
			}

			if _, ok = data[taxid][TaxidDelete]; !ok {
				data[taxid][TaxidDelete] = make([]TaxidChange, 0, 1)

				// extract lineage in NEW / CHANGED record
				if changes, ok = data[taxid][TaxidLineageChanged]; ok {
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       int16(version),
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						TaxidVersion:  changes[len(changes)-1].TaxidVersion,
						Change:        TaxidDelete,
						ChangeValue:   nil,
					})
				} else if changes, ok = data[taxid][TaxidNew]; ok {
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       int16(version),
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						TaxidVersion:  changes[len(changes)-1].TaxidVersion,
						Change:        TaxidDelete,
						ChangeValue:   nil,
					})
				} else { // firt archive
					data[taxid][TaxidDelete] = append(data[taxid][TaxidDelete], TaxidChange{
						Version:       int16(version),
						LineageTaxids: nil,
						TaxidVersion:  -1,
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
				// recording merged taxid
				if changes, ok = data[from][TaxidLineageChanged]; ok {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       int16(version),
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						TaxidVersion:  changes[len(changes)-1].TaxidVersion,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				} else if changes, ok = data[from][TaxidNew]; ok {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       int16(version),
						LineageTaxids: changes[len(changes)-1].LineageTaxids,
						TaxidVersion:  changes[len(changes)-1].TaxidVersion,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				} else {
					data[from][TaxidMerge] = append(data[from][TaxidMerge], TaxidChange{
						Version:       int16(version),
						LineageTaxids: nil,
						TaxidVersion:  -1,
						Change:        TaxidMerge,
						ChangeValue:   []int32{to},
					})
				}

				// add change to "to"
				if _, ok = data[to]; !ok {
					data[to] = make(map[TaxidChangeCode][]TaxidChange, 1)
				}
				if _, ok = data[to][TaxidAbsorb]; !ok {
					data[to][TaxidAbsorb] = make([]TaxidChange, 0, 1)
				}
				if len(data[to][TaxidAbsorb]) == 0 {
					data[to][TaxidAbsorb] = append(data[to][TaxidAbsorb], TaxidChange{
						Version:       int16(version),
						LineageTaxids: taxid2lineageTaxids[to],
						TaxidVersion:  int16(version),
						Change:        TaxidAbsorb,
						ChangeValue:   []int32{from},
					})
				} else { // append to
					data[to][TaxidAbsorb][0].ChangeValue = append(data[to][TaxidAbsorb][0].ChangeValue, from)
				}
			}
		}
	}

	// -------------- output --------------

	header := strings.Split("taxid,version,change,change-value,name,rank,lineage,lineage-taxids", ",")
	writer.Write(header)

	var cs []TaxidChange
	var c TaxidChange
	var tmp, items []string
	var tid int32
	var i int
	var taxid2name map[int32]string

	if config.Verbose {
		log.Infof("sorting %d taxids", len(data))
	}
	taxids := make([]int, len(data))
	i = 0
	for taxid := range data {
		taxids[i] = int(taxid)
		i++
	}
	sort.Ints(taxids)

	if config.Verbose {
		log.Infof("write to file: %s", config.OutFile)
	}
	for _, taxid := range taxids {
		// sort by version and then change
		changes := make([]TaxidChange, 0, len(data[int32(taxid)]))
		for _, cs = range data[int32(taxid)] {
			changes = append(changes, cs...)
		}

		sort.Sort(TaxidChanges(changes))

		for _, c = range changes {
			items = make([]string, 0, len(header))

			// taxid
			items = append(items, fmt.Sprintf("%d", taxid))

			// version
			items = append(items, versions[int(c.Version)])

			// change
			items = append(items, fmt.Sprintf("%s", c.Change))

			// change value
			switch c.Change {
			case TaxidMerge:
				items = append(items, fmt.Sprintf("%d", c.ChangeValue[0]))
			case TaxidAbsorb:
				tmp = make([]string, len(c.ChangeValue))
				for i, tid = range c.ChangeValue {
					tmp[i] = fmt.Sprintf("%d", tid)
				}
				items = append(items, strings.Join(tmp, ";"))
			default:
				items = append(items, "")
			}

			// name

			if c.TaxidVersion >= 0 {
				items = append(items, taxid2names[versions[int(c.TaxidVersion)]][int32(taxid)])
			} else {
				items = append(items, "")
			}

			// rank
			if c.TaxidVersion >= 0 {
				items = append(items, taxid2ranks[versions[int(c.TaxidVersion)]][int32(taxid)])
			} else {
				items = append(items, "")
			}

			// lineage
			if c.TaxidVersion >= 0 {
				taxid2name = taxid2names[versions[int(c.TaxidVersion)]]
				tmp = make([]string, len(c.LineageTaxids))
				for i, tid = range c.LineageTaxids {
					tmp[i] = taxid2name[tid]
				}

				items = append(items, strings.Join(tmp, ";"))
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
		checkError(fmt.Errorf("no unzipped directories found in path: %s", path))
	}
	for _, dir := range dirs {
		checkFile(filepath.Join(path, dir, "names.dmp"))
		checkFile(filepath.Join(path, dir, "nodes.dmp"))
		checkFile(filepath.Join(path, dir, "delnodes.dmp"))
		checkFile(filepath.Join(path, dir, "merged.dmp"))
	}

	if config.Verbose {
		log.Infof("%d archieves found in path: %s", len(dirs), path)
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
