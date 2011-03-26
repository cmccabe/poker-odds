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
	"testing"
)

func TestCardBag1(t *testing.T) {
	bag := Make52CardBag()
	test0 := Card { 2, DIAMONDS }
	bag0 := bag.Get(0)
	if (test0.Compare(bag0) != 0) {
		t.Errorf("expected:%s. got: %s\n", test0, bag0)
	}

	test1 := Card { 3, DIAMONDS }
	bag1 := bag.Get(4)
	if (test1.Compare(bag1) != 0) {
		t.Errorf("expected:%s. got: %s\n", test1, bag1)
	}

	bag.Subtract( &Card {2, DIAMONDS} )

	test2 := &Card { 2, CLUBS }
	bag2 := bag.Get(0)
	if (test2.Compare(bag2) != 0) {
		t.Errorf("expected:%s. got: %s\n", test2, bag2)
	}

	test3 := &Card { 3, CLUBS }
	bag3 := bag.Get(4)
	if (test3.Compare(bag3) != 0) {
		t.Errorf("expected:%s. got: %s\n", test3, bag3)
	}

}
