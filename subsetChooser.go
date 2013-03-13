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
)

/*
 * A subsetChooser can be used to step through all the possible subarrays of an
 * array.
 * So if our starting array was [ "a", "b", "c" ], and subset size was 2, we'd
 * step through [ "a", "b" ], [ "a", "c" ], and [ "b", "c" ]. Etc.
 */
type SubsetChooser struct {
	maxIdx uint
	subsetSize uint
	comb int64
}

// Creates a new SubsetChooser
func NewSubsetChooser(maxIdx uint, subsetSize uint) *SubsetChooser {
	ch := &SubsetChooser {maxIdx, subsetSize, 0}
	if (subsetSize > 63) {
		panic("sorry, this class can't handle subset sizes greater than 63" +
			"due to its use of 64-bit numbers to represent subsets.")
	}
	ch.comb = pow64(2, uint(subsetSize)) - 1
	return ch
}

// Gets the current subset
func (ch *SubsetChooser) Cur() []uint {
	var i uint
	ret := make([]uint, ch.subsetSize)
	var j uint = 0
	for i = 0; i < ch.maxIdx; i++ {
		if (((1 << i) & ch.comb) != 0) {
			ret[j] = i
			j++
		}
	}
	if (j != ch.subsetSize) {
		panic(fmt.Sprintf("logic error: failed to return a subset of size %d",
							ch.subsetSize))
	}
	return ret
}

// Advance to the next subset.
// Based on HAKMEM item 175.
// Returns false if there are no more subsets to view, true otherwise
func (ch *SubsetChooser) Next() bool {
	if (ch.comb == 0) {
		return false
	}
	u := ch.comb & -ch.comb
	v := u + ch.comb
	if (v==0) {
		ch.comb = 0
		return false
	}
	ch.comb = v + (((v^ch.comb)/u)>>2);
	if (ch.comb >= (1<<ch.maxIdx)) {
		ch.comb = 0
		return false
	}
	return true
}

func pow64(a int64, b uint) int64 {
	var ret int64
	ret = 1
	var i uint
	for i = 0; i < b; i++ {
		ret = ret * a
	}
	return ret
}
