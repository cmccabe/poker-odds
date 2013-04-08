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
)

type CardBag struct {
	allCards CardSlice
}

func Make52CardBag() *CardBag {
	// generate a CardBag that has every possible card
	bag := &CardBag { make(CardSlice, 52) }
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
	return ret
}

func (bag *CardBag) Subtract(c *Card) {
	nextAllCards := make(CardSlice, len(bag.allCards) - 1)
	i := 0
	for i = 0; i < len(bag.allCards); i++ {
		if (bag.allCards[i].Compare(c) == 0) {
			break
		}
	}
	if (i == len(bag.allCards)) {
		panic(fmt.Sprintf("tried to subtract %v from this cardbag, but " +
			"it doesn't currently contain that card.", c))
	}
	copy(nextAllCards[0:i], bag.allCards[0:i])
	copy(nextAllCards[i:], bag.allCards[i+1:])
	bag.allCards = nextAllCards
}

func (bag *CardBag) Get(num uint) *Card {
	return bag.allCards[num]
}

func (bag *CardBag) Len() int {
	return len(bag.allCards)
}

func (bag *CardBag) String() string {
	return fmt.Sprintf("CardBag{allCards=%v}",
		bag.allCards)
}
