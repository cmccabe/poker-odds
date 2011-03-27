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
	"fmt"
	"sort"
)

const (
	HAND_SZ = 5
)

const (
	JACK_VAL = 11
	QUEEN_VAL = 12
	KING_VAL = 13
	ACE_VAL = 14
)

const (
	HIGH_CARD = iota
	PAIR
	TWO_PAIR
	THREE_OF_A_KIND
	STRAIGHT
	FLUSH
	FULL_HOUSE
	FOUR_OF_A_KIND
	STRAIGHT_FLUSH
)

type Hand struct {
	ty int
	val [2]int
	flushSuit int
	cards CardSlice
}

func MakeHand(cards CardSlice) *Hand {
	// Sort the cards appropriately to make straight detection easier.
	sort.Sort(cards)

	h := &Hand{ -1, [2]int {-1, -1}, -1, cards }
	var vals = make(map[int] int)
	var suits = make(map[int] int)
	for i := range(cards) {
		c := cards[i]
		vals[c.val] = vals[c.val] + 1
		suits[c.suit] = vals[c.suit] + 1
	}

	// check for flush
	for i := range(suits) {
		if (suits[i] >= 4) {
			h.flushSuit = i
		}
	}
	// check for straight flush
	runEnd := -1
	runLen := 0
	prev := -1
	if (cards[len(cards)-1].val == ACE_VAL) {
		// Aces play both low and high in straights.
		//
		// This is a special case where we have a bunch of cards where the
		// highest card is a king, and we also have an ace.
		// In this situation, the last run of cards is actually one card longer
		// than it might seem without taking the ace into account.
		runLen++
	}
	for i := range(cards) {
		if (prev + 1 == cards[i].val) {
			runEnd = cards[i].val
			runLen++
			if (runLen >= 5) {
				// I know it may seem like by breaking out here, we may miss
				// some potential straights. But remember that a poker hand is
				// exactly 5 cards. So to have gotten to this point, all 5
				// cards must have played, so we must have discovered the best
				// possible straight.
				//
				// It's the Ace playing both high and low that makes this
				// weird.
				break
			}
		} else if (prev == cards[i].val) {
			// We have more than one card with the same value.
			// The duplicate cards don't help us get a straight, but they also
			// don't mean we don't have one.
		} else {
			// Clear the straight counter.
			runLen = 0
		}
	}
	if ((runLen >= 5) && (h.flushSuit != -1)) {
		h.val[0] = runEnd
		h.ty = STRAIGHT_FLUSH
		return h
	}

	freqs := make(map[int] []int)
	for k,v := range(vals) {
		if (v > 4) {
			panic(fmt.Sprintf("got %d of a kind for value %d (max is 4)\n",
				v, k))
		}
		curFreqs := freqs[v]
		m := 0
		for m = 0; m < len(curFreqs); m++ {
			if (curFreqs[m] >= k) {
				break
			}
		}
		newFreqs := make([]int, len(curFreqs) + 1)
		copy(newFreqs, curFreqs[:m])
		newFreqs[m] = k
		copy(newFreqs[m+1:], curFreqs[m:])
		freqs[v] = newFreqs
	}

	// four of a kind
	if (len(freqs[4]) > 0) {
		h.ty = FOUR_OF_A_KIND
		h.val[0] = freqs[4][0]
		return h
	}

	// full house
	if (len(freqs[3]) > 0) {
		if (len(freqs[3]) > 1) {
			h.val[0] = freqs[3][0]
			h.val[1] = freqs[3][1]
			h.ty = FULL_HOUSE
		} else if (len(freqs[2]) > 0) {
			h.val[0] = freqs[3][0]
			h.val[1] = freqs[2][0]
			h.ty = FULL_HOUSE
		}
	}

	// flush
	if (h.flushSuit != -1) {
		h.ty = FLUSH
		return h
	}

	// straight
	if (runLen >= 5) {
		h.val[0] = runEnd
		h.ty = STRAIGHT
		return h
	}

	// three of a kind
	if (len(freqs[3]) > 0) {
		h.val[0] = freqs[3][0]
		h.ty = THREE_OF_A_KIND
		return h
	}

	// two pairs
	if (len(freqs[2]) >= 2) {
		h.val[0] = freqs[2][0]
		h.val[1] = freqs[2][1]
		h.ty = TWO_PAIR
		return h
	}

	// a pair
	if (len(freqs[2]) >= 1) {
		h.val[0] = freqs[2][0]
		h.ty = PAIR
		return h
	}

	// I guess not.
	h.ty = HIGH_CARD
	return h
}

func (h *Hand) String() string {
	ret := "Hand(ty:"
	switch (h.ty) {
	case HIGH_CARD:
		ret += "HIGH CARD"
	case PAIR:
		ret += "PAIR of "
		ret += valToStr(h.val[0])
	case TWO_PAIR:
		ret += "TWO PAIR of "
		ret += valToStr(h.val[0])
		ret += " and "
		ret += valToStr(h.val[1])
	case THREE_OF_A_KIND:
		ret += "THREE OF A KIND of "
		ret += valToStr(h.val[0])
	case STRAIGHT:
		ret += "STRAIGHT with high of "
		ret += valToStr(h.val[0])
	case FLUSH:
		ret += "FLUSH in "
		ret += suitToStr(h.flushSuit)
	case FULL_HOUSE:
		ret += "FULL HOUSE of "
		ret += valToStr(h.val[0])
		ret += " full of "
		ret += valToStr(h.val[1])
	case FOUR_OF_A_KIND:
		ret += "FOUR OF A KIND of "
		ret += valToStr(h.val[0])
	case STRAIGHT_FLUSH:
		ret += "STRAIGHT FLUSH with high of "
		ret += valToStr(h.val[0])
		ret += " in "
		ret += suitToStr(h.flushSuit)
	}

	ret += ", cards:"
	sep := ""
	for c := range(h.cards) {
		ret += sep
		ret += h.cards[c].String()
		sep = ", "
	}
	ret += ")"

	return ret
}
