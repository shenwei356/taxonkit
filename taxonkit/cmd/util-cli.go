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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shenwei356/util/pathutil"
	"github.com/spf13/cobra"
	"github.com/twotwotwo/sorts"
)

func isStdin(file string) bool {
	return file == "-"
}

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

func errDataNotFound(dataDir string) {
	checkError(fmt.Errorf(`taxonomy data not found, please download and uncompress ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz, and copy "names.dmp", "nodes.dmp", "delnodes.dmp", and "merged.dmp" to %s`, dataDir))
}

func getConfigs(cmd *cobra.Command) Config {
	threads := getFlagPositiveInt(cmd, "threads")

	runtime.GOMAXPROCS(threads)
	sorts.MaxProcs = threads

	var val, dataDir string
	if val = os.Getenv("TAXONKIT_DB"); val != "" {
		dataDir = val
	} else {
		dataDir = getFlagString(cmd, "data-dir")
	}

	existed, err := pathutil.DirExists(dataDir)
	checkError(err)
	if !existed {
		checkError(os.MkdirAll(dataDir, 0777))
		errDataNotFound(dataDir)
	}

	nodesFile := filepath.Join(dataDir, "nodes.dmp")
	existed, err = pathutil.Exists(nodesFile)
	checkError(err)
	if !existed {
		errDataNotFound(dataDir)
	}

	namesFile := filepath.Join(dataDir, "names.dmp")
	existed, err = pathutil.Exists(namesFile)
	checkError(err)
	if !existed {
		errDataNotFound(dataDir)
	}

	delNodesFile := filepath.Join(dataDir, "delnodes.dmp")
	mergedFile := filepath.Join(dataDir, "merged.dmp")

	return Config{
		Threads:      threads,
		OutFile:      getFlagString(cmd, "out-file"),
		DataDir:      dataDir,
		NodesFile:    nodesFile,
		NamesFile:    namesFile,
		DelNodesFile: delNodesFile,
		MergedFile:   mergedFile,

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

func makeOutDir(outDir string, force bool) {
	pwd, _ := os.Getwd()
	if outDir != "./" && outDir != "." && pwd != filepath.Clean(outDir) {
		existed, err := pathutil.DirExists(outDir)
		checkError(errors.Wrap(err, outDir))
		if existed {
			empty, err := pathutil.IsEmpty(outDir)
			checkError(errors.Wrap(err, outDir))
			if !empty {
				if force {
					checkError(os.RemoveAll(outDir))
				} else {
					checkError(fmt.Errorf("out-dir not empty: %s, use --force to overwrite", outDir))
				}
			} else {
				checkError(os.RemoveAll(outDir))
			}
		}
		checkError(os.MkdirAll(outDir, 0777))
	}
}
