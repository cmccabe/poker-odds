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
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr,
`%s: the Texas Hold Em' poker odds calculator.

This program calculates your 'outs' for a Texas Hold Em' poker hand.
Texas Hold Em' is a popular version of poker where each player receives
exactly two secret cards. Then there are five rounds of betting.

The format used to specify cards is as follows:
[type][suit]
There are 4 suits:
C = clubs, D = diamonds, H = hearts, S = spades
There are 13 different card types:
1 = A = ace, K = king, Q = queen, J = jack, 2 = 2, ... 10 = 10

Usage:
-a [your hand as a whitespace-separated list of cards]
-b [the board as a whitespace-separated list of cards]
If no -b is given, it will be assumed that no cards are on the board.

-g [num_goroutines]               Set the number of goroutines to use.

-h this help message

Usage Example:
%s -a KS\ QS
Find the outs you have pre-flop with a king and queen of spades.
`, os.Args[0], os.Args[0])
}

func checkHoleLength(hlen int) {
	if (hlen == 2) {
		return
	}
	fmt.Printf("illegal hole length. Expected a length of 2, " +
		"but you gave %d hole cards.\n", hlen)
	os.Exit(1)
}

func checkBoardLength(blen int) {
	var validLens = []int { 0, 3, 4, 5 }
	for i := range(validLens) {
		if (blen == validLens[i]) {
			if (blen == 0) {
				fmt.Printf("Now calculating ALL possible hands that can be " +
					"made starting with these hole cards. This will take a " +
					"while! You may want to set GOMAXPROCS and use -g.\n" +
					"Note: It is much faster to calculate your odds " +
					"AFTER the flop.\n")
			}
			return
		}
	}

	fmt.Printf("illegal board length. Expected a length of %s, " +
		"but your board length was %d.\n", intsToStr(validLens), blen)
	os.Exit(1)
}

func intsToStr(s []int) (string) {
	ret := ""
	sep := ""
	for i := range(s) {
		ret += fmt.Sprintf("%s%d", sep, s[i])
		sep = ", "
	}
	return ret
}

func processHand(h *Hand) {
	fmt.Printf("%s\n", h.String())
}

/* Assumptions: we are the only players in the game
 * (Future enhancement: allow the user to specify cards that other players hold!)
 * 
 * 1. Get inputs
 * a. your hand (required)
 * b. the board (0 cards, 3 , 4, or 5 cards)
 *         Other numbers of cards represent errors
 *         (Future enhancement: support other poker games besides Texas Hold em')
 * 
 * 2. for all possible final boards:
 *        Determine the best type of hand we can make with this board and the
 *        hole cards.
 *
 * 3. Print out the odds of getting each type of hand. We can use the fact that
 *        each distinct final board is equally likely.
 *
 */
func main() {
	///// Parse and validate user input ///// 
	flag.Usage = usage
	var verbose = flag.Bool("v", false, "verbose")
	var help = flag.Bool("h", false, "help")
	var holeStr = flag.String("a", "", "your two hole cards")
	var boardStr = flag.String("b", "", "the board")
	var numCsp = flag.Int("g", 3, "number of goprocs")

	flag.Parse()
	if (*help) {
		usage()
		os.Exit(0)
	}
	if (*holeStr == "") {
		fmt.Printf("You must give two hole cards with -a\n")
		usage()
		os.Exit(1)
	}
	var hole CardSlice
	var errIdx int
	hole, errIdx = StrToCards(*holeStr)
	if (errIdx != -1) {
		fmt.Printf("Error parsing your hole cards: parse error at character %d\n",
					errIdx)
		os.Exit(1)
	}
	checkHoleLength(len(hole))
	if (*verbose) {
		fmt.Printf("Your hole cards: '%s'\n", hole.String());
	}
	var board CardSlice
	board, errIdx = StrToCards(*boardStr)
	if (errIdx != -1) {
		fmt.Printf("Error parsing the board: parse error at character %d\n",
					errIdx)
		os.Exit(1)
	}
	checkBoardLength(len(board))
	if (*verbose) {
		fmt.Printf("The board: '%s'\n", board.String());
	}
	base := make(CardSlice, len(board) + len(hole))
	copy(base, board)
	copy(base[len(board):], hole)
	dupe := base.HasDuplicates()
	if (dupe != nil) {
		fmt.Printf("The card %s appears more than once in your input! " +
			"That is not possible.\n", dupe)
		os.Exit(1)
	}

	///// Process cards ///// 
	csps := make([]*CardSliceProcessor, *numCsp)
	for i := range(csps) {
		csps[i] = NewCardSliceProcessor(base)
		go csps[i].GoCardSliceProcessor()
	}

	future := Make52CardBag()
	for i := range(base) {
		future.Subtract(base[i])
	}
	numFutureCards := SPREAD_MAX - len(base)
	futureChooser := NewSubsetChooser(uint(future.Len()), uint(numFutureCards))
	cspIdx := 0
	for ;; {
		futureC := futureChooser.Cur()
		for i := 0; i < numFutureCards; i++ {
			csps[cspIdx].Card <- future.Get(futureC[i])
		}
		cspIdx++
		if (cspIdx >= *numCsp) {
			cspIdx = 0
		}
		if (!futureChooser.Next()) {
			break
		}
	}

	// Tell cardSliceProcessors to finish
	for i := range(csps) {
		csps[i].Quit <- true
	}

	// Once each cardSliceProcessor is finished, get its results
	// Merge all results together
	allResults := new(ResultSet)
	for i := range(csps) {
		<-csps[i].Finished
		allResults.MergeResultSet(&csps[i].Results)
	}

	// Now print the final results
	fmt.Printf("results:\n%s", allResults.String())
}
