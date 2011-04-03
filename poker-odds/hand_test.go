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
	"testing"
)

func expectHand(t *testing.T, c CardSlice, eTy int, eVal [2]int, eFlushSuit int) {
	aHand := MakeHand(c)
	eHand := &Hand { eTy, eVal, eFlushSuit, c }
	if (!aHand.Identical(eHand)) {
		t.Errorf("expected MakeHand to create: %s.\n" +
				"Instead, it created: %s", eHand, aHand)
	}
}

func cvalToString(cval int) string {
	switch (cval) {
	case -1:
		return "less than"
	case 0:
		return "equal to"
	case 1:
		return "greater than"
	}
	panic(fmt.Sprintf("unknown cval %d", cval))
}

func compareHands(t *testing.T, eRes int, a *Hand, b *Hand) {
	aRes := a.Compare(b)
	if (aRes != eRes) {
		t.Errorf("expected hand:\n%s\nto be %s hand:\n%s, but it was " +
				 "marked as %s", a, cvalToString(eRes), b, cvalToString(aRes))
	}
}

func TestHand1(t *testing.T) {
	c1 := CardSlice { &Card{2, DIAMONDS}, &Card{2, CLUBS}, &Card{2, HEARTS},
					&Card{QUEEN_VAL, SPADES}, &Card{KING_VAL, SPADES} }
	expectHand(t, c1, THREE_OF_A_KIND, [2]int{2, -1}, -1)

	c2 := CardSlice { &Card{8, DIAMONDS}, &Card{9, DIAMONDS}, &Card{10, DIAMONDS},
					&Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }
	expectHand(t, c2, STRAIGHT_FLUSH, [2]int{QUEEN_VAL, -1}, DIAMONDS)

	c3 := CardSlice { &Card{8, DIAMONDS}, &Card{KING_VAL, DIAMONDS}, &Card{10, DIAMONDS},
					&Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }
	expectHand(t, c3, FLUSH, [2]int{-1, -1}, DIAMONDS)

	c4 := CardSlice { &Card{8, CLUBS}, &Card{8, DIAMONDS}, &Card{10, HEARTS},
					&Card{4, DIAMONDS}, &Card{10, DIAMONDS} }
	expectHand(t, c4, TWO_PAIR, [2]int{8, 10}, -1)

	c5 := CardSlice { &Card{8, CLUBS}, &Card{9, DIAMONDS}, &Card{2, HEARTS},
					&Card{KING_VAL, DIAMONDS}, &Card{ACE_VAL, CLUBS} }
	expectHand(t, c5, HIGH_CARD, [2]int{-1, -1}, -1)

	compareHands(t, 1, MakeHand(CardSlice { &Card{8, DIAMONDS}, &Card{9, DIAMONDS},
		&Card{10, DIAMONDS}, &Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }),
			MakeHand(CardSlice { &Card{8, DIAMONDS}, &Card{7, DIAMONDS},
		&Card{10, DIAMONDS}, &Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }))

	compareHands(t, -1, MakeHand(CardSlice { &Card{10, DIAMONDS}, &Card{KING_VAL, DIAMONDS},
		&Card{10, CLUBS}, &Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }),
			MakeHand(CardSlice { &Card{8, DIAMONDS}, &Card{7, DIAMONDS},
		&Card{10, DIAMONDS}, &Card{QUEEN_VAL, DIAMONDS}, &Card{JACK_VAL, DIAMONDS} }))
}
