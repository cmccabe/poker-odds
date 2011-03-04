/* assumptions: we are the only players in the game
 * later this can be refined if you have an idea of what other players have or don't have.
 * 
 * 1. get inputs
 * a. your hand (required)
 * b. the board (0 cards, 3 , 4, or 5 cards)
 *         Other numbers of cards represent errors
 * 
 * 2. for all possible poker hands that can be formed by your hand:
 * calculate the odds of getting that hand (100% if you already have it)
 * 
 * 3. print all odds in a nice format
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, `
%s: the Texas Hold Em' poker odds calculator.

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
-h this help message

Usage Example:
%s -a \"KS QS\"
Find the outs you have pre-flop with a king and queen of spades.

`, os.Args[0], os.Args[0])
}

const (
A Ty1 = iota;
B;
C
)

func strToCardList(str *string) {
	var parse_state
	for _, r := range *str {
		switch(r) {
		case parse_state
		}
	}
	fmt.Printf("str='%s'\n", str)
}

func main() {
	flag.Usage = usage
	var verbose = flag.Bool("v", false, "verbose")
	var help = flag.Bool("h", false, "help")
	var hand = flag.String("a", "", "your hand")
	flag.Parse()

	if (*help) {
		usage()
		os.Exit(0)
	}
	if (*verbose) {
		fmt.Println("Hello, 世界")
	} else {
		fmt.Println("exiting")
	}
	if (*hand == "") {
		fmt.Printf("You must give a hand with -a\n")
		usage()
		os.Exit(1)
	}
	fmt.Printf("your hand: '%s'\n", *hand)
	strToCardList(hand)
}
