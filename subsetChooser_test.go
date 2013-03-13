/*
 * Copyright 2011 Colin Patrick McCabe
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 2.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"sort"
	"testing"
)

type uintSlice []uint
type uintSliceSlice []uintSlice

func (arr uintSliceSlice) Len() int {
	return len(arr)
}

func (arr uintSliceSlice) Less(i, j int) bool {
	c := arr[i].Compare(arr[j])
	return (c < 0)
}

func (arr uintSliceSlice) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (a uintSlice) Compare(b uintSlice) int {
	if (len(a) < len(b)) {
		return -1
	} else if (len(a) > len(b)) {
		return 1;
	}
	for i := 0; i < len(a); i++ {
		if (a[i] < b[i]) {
			return -1
		} else if (a[i] > b[i]) {
			return 1
		}
	}
	return 0
}

func (a uintSliceSlice) Compare(b uintSliceSlice) int {
	if (len(a) < len(b)) {
		return -1
	} else if (len(a) > len(b)) {
		return 1;
	}
	for i := 0; i < len(a); i++ {
		c := a[i].Compare(b[i])
		if (c != 0) {
			return c
		}
	}
	return 0
}

func (arr uintSlice) String() string {
	var ret string
	ret += "["
	sep := ""
	for i := 0; i < len(arr); i++ {
		ret += sep
		ret += fmt.Sprintf("%d", arr[i])
		sep = ", "
	}
	ret += "]"
	return ret
}

func (arr uintSliceSlice) String() string {
	var ret string
	ret += "["
	sep := ""
	for i := 0; i < len(arr); i++ {
		ret += sep
		ret += arr[i].String()
		sep = ", "
	}
	ret += "]"
	return ret
}

func test(t *testing.T, maxIdx uint, subsetSize uint,
			expected *uintSliceSlice) {
	var all uintSliceSlice
	ch := NewSubsetChooser(maxIdx, subsetSize)
	for ;; {
		s := ch.Cur()
		all = append(all, s)
		if (!ch.Next()) {
			break
		}
	}
	sort.Sort(all)
	sort.Sort(expected)
	if (all.Compare(*expected) != 0) {
		t.Errorf("expected:%s. got: %s\n",
			expected.String(), all.String())
	}
}

func TestSubsetChooser1(t *testing.T) {
	test(t, 3, 2, &uintSliceSlice{ {0, 1}, {0, 2}, {1, 2} } )
	test(t, 2, 1, &uintSliceSlice{ {0}, {1} } )
	test(t, 3, 1, &uintSliceSlice{ {0}, {1}, {2} } )
	test(t, 3, 3, &uintSliceSlice{ {0, 1, 2} } )
	test(t, 4, 3, &uintSliceSlice{ {0, 1, 2}, {0, 1, 3}, {0, 2, 3}, {1, 2, 3} } )

	test(t, 100, 0, &uintSliceSlice{ {} })
}
