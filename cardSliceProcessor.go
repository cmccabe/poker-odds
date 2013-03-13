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

type CardSliceProcessor struct {
	Card chan *Card
	Quit chan bool
	Finished chan bool
	base CardSlice
	Results ResultSet
}

func NewCardSliceProcessor(base_ CardSlice) *CardSliceProcessor {
	ret := new(CardSliceProcessor)
	ret.Card = make(chan *Card)
	ret.Quit = make(chan bool)
	ret.Finished = make(chan bool)
	ret.base = base_.Copy()
	return ret
}

func (csp *CardSliceProcessor) processSpread(spread CardSlice) {
	setupChooser := NewSubsetChooser(SPREAD_MAX, HAND_SZ)
	var tmpRes ResultSet
	for ;; {
		setup := setupChooser.Cur()
		setupC := make(CardSlice, len(setup))
		for i := range(setup) {
			setupC[i] = spread[setup[i]]
		}
		h := MakeHand(setupC)
		tmpRes.AddHand(h)
		if (!setupChooser.Next()) {
			break
		}
	}
	csp.Results.AddHandTy(tmpRes.GetBestHandTy())
}

func (csp *CardSliceProcessor) GoCardSliceProcessor() {
	spread := make(CardSlice, SPREAD_MAX)
	copy(spread, csp.base)
	j := len(csp.base)
	for {
		if (j == len(spread)) {
			csp.processSpread(spread)
			j = len(csp.base)
		}
		select {
		case c := <-csp.Card:
			//fmt.Printf("%p: received card %s\n", csp, c.String())
			spread[j] = c
			j++
		case <-csp.Quit:
			csp.Finished <-true
			return
		}
	}
}
