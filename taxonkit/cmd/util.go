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
	"strings"
	"unsafe"
)

func stringSplitN(s string, sep string, n int, a *[]string) {
	if a == nil {
		tmp := make([]string, n)
		a = &tmp
	}

	n--
	i := 0
	for i < n {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		(*a)[i] = s[:m]
		s = s[m+len(sep):]
		i++
	}
	(*a)[i] = s

	(*a) = (*a)[:i+1]
}

func stringSplitNByByte(s string, sep byte, n int, a *[]string) {
	if a == nil {
		tmp := make([]string, n)
		a = &tmp
	}

	n--
	i := 0
	for i < n {
		m := strings.IndexByte(s, sep)
		if m < 0 {
			break
		}
		(*a)[i] = s[:m]
		s = s[m+1:]
		i++
	}
	(*a)[i] = s

	(*a) = (*a)[:i+1]
}

func reverseUint32s(s []uint32) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func MaxInts(a int, vals ...int) int {
	max := a
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}

func CopyString(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}
