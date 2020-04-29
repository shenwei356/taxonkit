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
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "taxonkit",
	Short: "NCBI Taxonomy Toolkit",
	Long:  "",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// defaulDataDir is the directory of TaxonKit
var defaulDataDir string

func init() {
	var err error
	defaulDataDir, err = homedir.Expand("~/.taxonkit/")
	checkError(errors.Wrap(err, "get data directory"))

	RootCmd.Long = fmt.Sprintf(`TaxonKit - A Cross-platform and Efficient NCBI Taxonomy Toolkit

Version: %s

Author: Wei Shen <shenwei356@gmail.com>

Source code: https://github.com/shenwei356/taxonkit
Documents  : https://bioinf.shenwei.me/taxonkit
Citation   : https://www.biorxiv.org/content/early/2019/01/08/513523

Dataset:

    Please download and decompress "taxdump.tar.gz":
    ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz

    and copy "names.dmp", "nodes.dmp", "delnodes.dmp" and "merged.dmp" to data directory:
    "%s"

    or some other directory, and later you can refer to using flag --data-dir,
    or environment variable TAXONKIT_DB

`, VERSION, defaulDataDir)

	defaultThreads := runtime.NumCPU()
	if defaultThreads > 2 {
		defaultThreads = 2
	}

	RootCmd.PersistentFlags().IntP("threads", "j", defaultThreads, "number of CPUs. 2 is enough (default value: 1 for single-CPU PC, 2 for others)")
	RootCmd.PersistentFlags().StringP("out-file", "o", "-", `out file ("-" for stdout, suffix .gz for gzipped out)`)
	RootCmd.PersistentFlags().StringP("data-dir", "", defaulDataDir, "directory containing nodes.dmp and names.dmp")
	RootCmd.PersistentFlags().BoolP("verbose", "", false, "print verbose information")
	RootCmd.PersistentFlags().BoolP("line-buffered", "", false, "use line buffering on output, i.e., immediately writing to stdin/file for every line of output")
}
