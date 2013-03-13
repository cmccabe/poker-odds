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
	"testing"
)

func TestResultSet1(t *testing.T) {
	r1 := new(ResultSet)
	r1.AddHand( MakeHand(CardSlice { &Card{2, DIAMONDS}, &Card{2, CLUBS},
			&Card{2, HEARTS}, &Card{QUEEN_VAL, SPADES}, &Card{KING_VAL, SPADES} }))
	if (r1.handTyCnt[THREE_OF_A_KIND] != 1) {
		t.Errorf("expected ResultSet r1 to have 1 THREE_OF_A_KIND hand in it.")
	}

	r2 := new(ResultSet)
	r2.AddHand( MakeHand(CardSlice { &Card{8, DIAMONDS}, &Card{8, CLUBS},
			&Card{8, HEARTS}, &Card{QUEEN_VAL, SPADES}, &Card{KING_VAL, SPADES} }) )
	if (r2.handTyCnt[THREE_OF_A_KIND] != 1) {
		t.Errorf("expected ResultSet r2 to have 1 THREE_OF_A_KIND hand in it.")
	}

	r1.MergeResultSet(r2)

	if (r1.handTyCnt[THREE_OF_A_KIND] != 2) {
		t.Errorf("expected ResultSet r1 to have 2 THREE_OF_A_KIND hands in it.")
	}
}
