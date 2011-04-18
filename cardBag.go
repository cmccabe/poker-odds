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

package pokerodds

import (
	"sort"
)

type CardBag struct {
	allCards CardSlice
	skipped CardSlice
}

func Make52CardBag() *CardBag {
	// generate a CardBag that has every possible card
	bag := &CardBag { make(CardSlice, 52), CardSlice {} }
	idx := 0
	for val := 2; val <= 14; val++ {
		for suit := range([]int { DIAMONDS, CLUBS, HEARTS, SPADES }) {
			bag.allCards[idx] = new(Card)
			bag.allCards[idx].suit = suit
			bag.allCards[idx].val = val
			idx++
		}
	}

	sort.Sort(bag.allCards)
	return bag
}

func (bag *CardBag) Clone() *CardBag {
	ret := new(CardBag)
	ret.allCards = bag.allCards
	ret.skipped = make(CardSlice, 0)
	copy(ret.skipped, bag.skipped)
	return ret
}

func (bag *CardBag) Subtract(c *Card) {
	bag.skipped = append(bag.skipped, c)
}

func (bag *CardBag) Get(num uint) *Card {
	onum := num
	for i := range(bag.skipped) {
		if (bag.skipped[i].Compare(bag.allCards[num]) <= 0) {
			onum++
		}
	}
	return bag.allCards[onum]
}

func (bag *CardBag) Len() int {
	return len(bag.allCards) - len(bag.skipped)
}
