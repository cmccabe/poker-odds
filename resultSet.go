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

type ResultSet struct {
	handTyCnt [MAX_HANDS] int64
}

func (res *ResultSet) AddHand(h *Hand) {
	res.handTyCnt[h.ty] = res.handTyCnt[h.ty] + 1
}

func (res *ResultSet) AddHandTy(h int) {
	res.handTyCnt[h] = res.handTyCnt[h] + 1
}

func (res *ResultSet) GetBestHandTy() int {
	for i := MAX_HANDS - 1; i >= 0; i-- {
		if (res.handTyCnt[i] > 0) {
			return i
		}
	}
	return HIGH_CARD
}

func (res *ResultSet) MergeResultSet(rhs *ResultSet) {
	for t := HIGH_CARD; t < MAX_HANDS; t++ {
		res.handTyCnt[t] = res.handTyCnt[t] + rhs.handTyCnt[t]
	}
}

func (res *ResultSet) String() string {
	var totalHands int64
	totalHands = 0
	for i := range(res.handTyCnt) {
		totalHands = totalHands + res.handTyCnt[i]
	}

	ret := ""
	for i := range(res.handTyCnt) {
		percent := float32(res.handTyCnt[i])
		percent *= 100.0
		percent /= float32(totalHands);
		if (percent > 0.0) {
			ret += fmt.Sprintf("%03.2f%% chance of %s\n", percent, HandTyToStr(i))
		}
	}
	return ret
}

