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

import "regexp"

var rankList = []string{
	"",
	"superkingdom",
	"kingdom",
	"phylum",
	"class",
	"order",
	"family",
	"genus",
	"species",
	"subspecies",
	"subspecies/strain",
	"strain",
}

var srankList = []string{
	"",
	"k",
	"K",
	"p",
	"c",
	"o",
	"f",
	"g",
	"s",
	"t",
	"S",
	"T",
}

var rank2symbol = map[string]string{
	"superkingdom":      "k",
	"kingdom":           "K",
	"phylum":            "p",
	"class":             "c",
	"order":             "o",
	"family":            "f",
	"genus":             "g",
	"species":           "s",
	"subspecies/strain": "t",
	"subspecies":        "S",
	"strain":            "T",
}

var symbol2rank = map[string]string{
	"k": "superkingdom",
	"K": "kingdom",
	"p": "phylum",
	"c": "class",
	"o": "order",
	"f": "family",
	"g": "genus",
	"s": "species",
	"t": "subspecies/strain",
	"S": "subspecies",
	"T": "strain",
}
var symbol2weight = map[string]float32{
	"k": 1,
	"K": 1.5,
	"p": 2,
	"c": 3,
	"o": 4,
	"f": 5,
	"g": 6,
	"s": 7,
	"t": 8,
	"S": 9,
	"T": 10,
}

var reRankPlaceHolder = regexp.MustCompile(`\{(\w)\}`)

var reRankPlaceHolders = map[string]*regexp.Regexp{
	"k": regexp.MustCompile(`\{k\}`),
	"K": regexp.MustCompile(`\{K\}`),
	"p": regexp.MustCompile(`\{p\}`),
	"c": regexp.MustCompile(`\{c\}`),
	"o": regexp.MustCompile(`\{o\}`),
	"f": regexp.MustCompile(`\{f\}`),
	"g": regexp.MustCompile(`\{g\}`),
	"s": regexp.MustCompile(`\{s\}`),
	"t": regexp.MustCompile(`\{t\}`),
	"S": regexp.MustCompile(`\{S\}`),
	"T": regexp.MustCompile(`\{T\}`),
}

const norank = "no rank"
