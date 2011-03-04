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
)

const (
	PARSE_STATE_EAT_VAL = iota
	PARSE_STATE_EAT_VAL_SAW_1
	PARSE_STATE_EAT_SUIT
)

const (
	CLUBS = iota
	DIAMONDS
	HEARTS
	SPADES
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
-h this help message

Usage Example:
%s -a KS\ QS
Find the outs you have pre-flop with a king and queen of spades.
`, os.Args[0], os.Args[0])
}

type card struct {
	val int
	suit int
}

func valToStr(v int) (string) {
	switch {
	case v == 1:
		return "A"
	case v == 11:
		return "J"
	case v == 12:
		return "Q"
	case v == 13:
		return "K"
	}
	return fmt.Sprintf("%d", v)
}

func suitToStr(s int) (string) {
	switch {
	case s == CLUBS:
		return "♣C"
	case s == DIAMONDS:
		return "♦D"
	case s == HEARTS:
		return "♥H"
	case s == SPADES:
		return "♠S"
	}
	return ""
}

func (p *card) toStr() string {
	return fmt.Sprintf("%s%s", valToStr(p.val), suitToStr(p.suit))
}

func strToCard(str *string) (myCard *card) {
	myCard = new(card)
	var parseState = PARSE_STATE_EAT_VAL
	for _, c := range *str {
		switch {
		case parseState == PARSE_STATE_EAT_VAL:
			switch {
			case c == ' ' || c == '\t':
				continue
			case c == '1':
				parseState = PARSE_STATE_EAT_VAL_SAW_1
			case c >= '2' && c <= '9':
				myCard.val = c - '0'
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'J':
				myCard.val = 11
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'Q':
				myCard.val = 12
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'K':
				myCard.val = 13
				parseState = PARSE_STATE_EAT_SUIT
			case c == 'A':
				myCard.val = 1
				parseState = PARSE_STATE_EAT_SUIT
			default:
				return nil
			}
		case parseState == PARSE_STATE_EAT_VAL_SAW_1:
			switch {
			case c == '0':
				myCard.val = 10
				parseState = PARSE_STATE_EAT_SUIT
			default:
				return nil
			}
		case parseState == PARSE_STATE_EAT_SUIT:
			switch {
			case c == 'C':
				myCard.suit = CLUBS
				return myCard
			case c == 'D':
				myCard.suit = DIAMONDS
				return myCard
			case c == 'H':
				myCard.suit = HEARTS
				return myCard
			case c == 'S':
				myCard.suit = SPADES
				return myCard
			default:
				return nil
			}
		}
	}
    return nil
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
	var card = strToCard(hand)
	fmt.Printf("card = %s\n", card.toStr())
}
