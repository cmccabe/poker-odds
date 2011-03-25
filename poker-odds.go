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
	"sort"
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

type cardArray []*card

func (cards cardArray) Len() int {
	return len(cards)
}

func (cards cardArray) Less(i, j int) bool {
	if (cards[i].val < cards[j].val) {
		return true
	}
	if (cards[i].val > cards[j].val) {
		return false
	}
	if (cards[i].suit < cards[j].suit) {
		return true
	}
	return false
}

func (cards cardArray) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
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

func (p *card) String() string {
	return fmt.Sprintf("%s%s", valToStr(p.val), suitToStr(p.suit))
}

func strToCard(str string, cnt *int) (myCard *card) {
	myCard = new(card)
	var parseState = PARSE_STATE_EAT_VAL
	for ;*cnt < len(str); {
		var c = str[*cnt]
		*cnt++
		switch {
		case parseState == PARSE_STATE_EAT_VAL:
			switch {
			case c == ' ' || c == '\t':
				continue
			case c == '1':
				parseState = PARSE_STATE_EAT_VAL_SAW_1
			case c >= '2' && c <= '9':
				myCard.val = (int)(c - '0')
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
			case c == 'D':
				myCard.suit = DIAMONDS
			case c == 'H':
				myCard.suit = HEARTS
			case c == 'S':
				myCard.suit = SPADES
			default:
				return nil
			}
			return myCard
		}
	}
	*cnt = -1
    return nil
}

func (p *card) Compare(rhs *card) int {
	if (p.suit < rhs.suit) {
		return -1;
	}
	if (p.suit > rhs.suit) {
		return 1;
	}
	if (p.val < rhs.val) {
		return -1;
	}
	if (p.val > rhs.val) {
		return 1;
	}
	return 0;
}

func hasDuplicates(c []*card) *card {
	for i := range(c) {
		for j := range(c) {
			if i == j {
				continue
			}
			if (c[i].Compare(c[j]) == 0) {
				return c[i]
			}
		}
	}
	return nil
}

func strToCards(str string) (cards []*card, cnt int) {
	for cnt = 0; cnt != -1; {
		var c = strToCard(str, &cnt)
		if (c == nil) {
			return
		}
		cards = append(cards,c)
	}
	return
}

func cardsToStr(c []*card) (string) {
	ret := ""
	sep := ""
	for i := range(c) {
		ret += fmt.Sprintf("%s%s", sep, c[i].String())
		sep = ", "
	}
	return ret
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

func checkBoardLength(l int) {
	var validLens = []int { 0, 3, 4, 5 }
	for i := range(validLens) {
		if (int(l) == validLens[i]) {
			return
		}
	}

	fmt.Printf("illegal board length. Expected a length of %s, " +
	"but your board length was %d.\n", intsToStr(validLens), l)
	os.Exit(1)
}

func generateAllVals(i *int, cards *[52]*card, suit int) {
	for val := 1; val <= 13; val++ {
		cards[*i] = new(card)
		cards[*i].suit = suit
		cards[*i].val = val
		*i++
	}
}

func generateAllCards(cards *[52]*card) {
	i := 0
	generateAllVals(&i, cards, CLUBS)
	generateAllVals(&i, cards, DIAMONDS)
	generateAllVals(&i, cards, HEARTS)
	generateAllVals(&i, cards, SPADES)
}

// find next k-combination
// assume x has form x'01^a10^b in binary
func nextCombination(x *int64) bool {
	u := *x & -*x // extract rightmost bit 1; u =  0'00^a10^b
	v := u + *x // set last non-trailing bit 0, and clear to the right; v=x'10^a00^b
	if (v==0) { // then overflow in v, or x==0
		return false; // signal that next k-combination cannot be represented
	}
	*x = v + (((v^*x)/u)>>2); // v^x = 0'11^a10^b, (v^x)/u = 0'0^b1^{a+2}, and x ← x'100^b1^a
	if (*x >= (1<<52)) {
		return false; // too big
	}
	return true; // successful completion
}

func combinationToCards(comb int64, allCards *[52]*card, holeC *[]*card,
						boardC *[]*card) ([]*card) {
	var ret []*card = make([]*card, 7)
	copy(ret, *holeC)
	copy(ret[len(*holeC):], *boardC)
	var n = len(*holeC) + len(*boardC)
	for i := range(allCards) {
		if (((1 << uint(i)) & comb) != 0) {
			ret[n] = &(*allCards[i])
			n++
			if (n >= 5) {
				return ret
			}
		}
	}
	fmt.Printf("combinationToHand: logic error: got to unreachable point\n")
	os.Exit(1)
	return ret
}

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

type hand struct {
	ty int
	val [2]int
	flushSuit int
	cards []*card
}

func (h hand) String() string {
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

func pow64(a int64, b uint) int64 {
	var ret int64
	ret = 1
	var i uint
	for i = 0; i < b; i++ {
		ret = ret * a
	}
	return ret
}

func remainingCardsToComb(cur uint) int64 {
	fmt.Printf("remainingCardsToComb, cur=%d\n", cur)
	var rem uint
	rem = 5 - cur
	if (rem < 0) {
		panic("remainingCardsToComb: invalid argument")
	}
	return pow64(2, uint(rem)) - 1
}

func makeHand(cards cardArray) *hand {
	// Sort the cards appropriately to make straight detection easier.
	sort.Sort(cards)

	ret := new(hand)
	ret.cards = cards
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
			ret.flushSuit = i
		}
	}
	// check for straight flush
	runEnd := -1
	runLen := 0
	prev := -1
	for i := range(cards) {
		if (prev + 1 == cards[i].val) {
			runEnd = cards[i].val
			runLen++
		} else {
			runLen = 0
		}
	}
	if ((runLen >= 5) && (ret.flushSuit != 0)) {
		ret.val[0] = runEnd
		ret.ty = STRAIGHT_FLUSH
		return ret
	}

	freqs := make(map[int] []int)
	for k,v := range(vals) {
		if (v > 4) {
			fmt.Printf("got %d of a kind for value %d (max is 4)\n", v, k)
			os.Exit(0)
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
		ret.ty = FOUR_OF_A_KIND
		ret.val[0] = freqs[4][0]
		return ret
	}

	// full house
	if (len(freqs[3]) > 0) {
		if (len(freqs[3]) > 1) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[3][1]
			ret.ty = FULL_HOUSE
		} else if (len(freqs[2]) > 0) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[2][0]
			ret.ty = FULL_HOUSE
		}
	}

	// flush
	if (ret.flushSuit != 0) {
		ret.ty = FLUSH
		return ret
	}

	// straight
	if (runLen >= 5) {
		ret.val[0] = runEnd
		ret.ty = STRAIGHT
		return ret
	}

	// three of a kind
	if (len(freqs[3]) > 0) {
		ret.val[0] = freqs[3][0]
		ret.ty = THREE_OF_A_KIND
		return ret
	}

	// two pairs
	if (len(freqs[2]) >= 2) {
		ret.val[0] = freqs[2][0]
		ret.val[1] = freqs[2][1]
		ret.ty = TWO_PAIR
		return ret
	}

	// a pair
	if (len(freqs[2]) >= 1) {
		ret.val[0] = freqs[2][0]
		ret.ty = PAIR
		return ret
	}

	// I guess not.
	ret.ty = HIGH_CARD
	return ret
}

func main() {
	flag.Usage = usage
	var verbose = flag.Bool("v", false, "verbose")
	var help = flag.Bool("h", false, "help")
	var hole = flag.String("a", "", "your two hole cards")
	var board = flag.String("b", "", "the board")
	flag.Parse()

	if (*help) {
		usage()
		os.Exit(0)
	}
	if (*hole == "") {
		fmt.Printf("You must give two hole cards with -a\n")
		usage()
		os.Exit(1)
	}
	holeC, errIdx := strToCards(*hole)
	if (errIdx != -1) {
		fmt.Printf("Error parsing your hole cards: parse error at character %d\n",
					errIdx)
		os.Exit(1)
	}

	if (*verbose) {
		fmt.Printf("Your hole cards: %s\n", cardsToStr(holeC));
	}

	fmt.Printf("board = %s\n", *board)
	boardC, bErrIdx := strToCards(*board)
	if (bErrIdx != -1) {
		fmt.Printf("Error parsing the board: parse error at character %d\n",
					bErrIdx)
		os.Exit(1)
	}
	checkBoardLength(len(boardC))
	if (*verbose) {
		fmt.Printf("The board: %s\n", cardsToStr(boardC));
	}

	var c = make([]*card, len(boardC) + len(holeC))
	copy(c, boardC)
	copy(c[len(boardC):], holeC)
	dupe := hasDuplicates(c)
	if (dupe != nil) {
		fmt.Printf("The card %s appears more than once!\n", dupe)
		os.Exit(1)
	}

	// generate all cards
	var allCards [52]*card
	generateAllCards(&allCards)

	remCards := 5 - len(boardC)
	if (remCards < 0) {
		fmt.Printf("invalid board length %d\n", len(boardC))
		os.Exit(1)
	}

	choose 5 cards out of the hole + board + filler

	// There will be 2 hole cards and 5 cards on the board, making 7 cards to
	// choose from.
	card_chooser := pow64(2, 7) - 1

func do_combination(fixedCards []*card) {
	// what can we choose from in our 'random picks'?
	futureCards = allCards.deepcopy
	futureCards.subtract(fixedCards)

	numFutureCards = 5 - len(fixedCards)
	iterateThroughFutureChoices {
		fiveCards := make([]*card, 5)
		copy(fiveCards, fixedCards)
		for i := 0; i < numFutureCards; i++ {
			for j := range(futureCards) {
				if (((1 << uint(j)) & comb) != 0) {
					ret[n] = &(*allCards[i])
					n++
					if (n >= 5) {
						return ret
					}
				}
			}
		}
	}
	fmt.Printf("combinationToHand: logic error: got to unreachable point\n")
	os.Exit(1)
	return ret
}

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

type hand struct {
	ty int
	val [2]int
	flushSuit int
	cards []*card
}

func (h hand) String() string {
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

func pow64(a int64, b uint) int64 {
	var ret int64
	ret = 1
	var i uint
	for i = 0; i < b; i++ {
		ret = ret * a
	}
	return ret
}

func remainingCardsToComb(cur uint) int64 {
	fmt.Printf("remainingCardsToComb, cur=%d\n", cur)
	var rem uint
	rem = 5 - cur
	if (rem < 0) {
		panic("remainingCardsToComb: invalid argument")
	}
	return pow64(2, uint(rem)) - 1
}

func makeHand(cards cardArray) *hand {
	// Sort the cards appropriately to make straight detection easier.
	sort.Sort(cards)

	ret := new(hand)
	ret.cards = cards
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
			ret.flushSuit = i
		}
	}
	// check for straight flush
	runEnd := -1
	runLen := 0
	prev := -1
	for i := range(cards) {
		if (prev + 1 == cards[i].val) {
			runEnd = cards[i].val
			runLen++
		} else {
			runLen = 0
		}
	}
	if ((runLen >= 5) && (ret.flushSuit != 0)) {
		ret.val[0] = runEnd
		ret.ty = STRAIGHT_FLUSH
		return ret
	}

	freqs := make(map[int] []int)
	for k,v := range(vals) {
		if (v > 4) {
			fmt.Printf("got %d of a kind for value %d (max is 4)\n", v, k)
			os.Exit(0)
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
		ret.ty = FOUR_OF_A_KIND
		ret.val[0] = freqs[4][0]
		return ret
	}

	// full house
	if (len(freqs[3]) > 0) {
		if (len(freqs[3]) > 1) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[3][1]
			ret.ty = FULL_HOUSE
		} else if (len(freqs[2]) > 0) {
			ret.val[0] = freqs[3][0]
			ret.val[1] = freqs[2][0]
			ret.ty = FULL_HOUSE
		}
	}

	// flush
	if (ret.flushSuit != 0) {
		ret.ty = FLUSH
		return ret
	}

	// straight
	if (runLen >= 5) {
		ret.val[0] = runEnd
		ret.ty = STRAIGHT
		return ret
	}

	// three of a kind
	if (len(freqs[3]) > 0) {
		ret.val[0] = freqs[3][0]
		ret.ty = THREE_OF_A_KIND
		return ret
	}

	// two pairs
	if (len(freqs[2]) >= 2) {
		ret.val[0] = freqs[2][0]
		ret.val[1] = freqs[2][1]
		ret.ty = TWO_PAIR
		return ret
	}

	// a pair
	if (len(freqs[2]) >= 1) {
		ret.val[0] = freqs[2][0]
		ret.ty = PAIR
		return ret
	}

	// I guess not.
	ret.ty = HIGH_CARD
	return ret
}

func main() {
	flag.Usage = usage
	var verbose = flag.Bool("v", false, "verbose")
	var help = flag.Bool("h", false, "help")
	var hole = flag.String("a", "", "your two hole cards")
	var board = flag.String("b", "", "the board")
	flag.Parse()

	if (*help) {
		usage()
		os.Exit(0)
	}
	if (*hole == "") {
		fmt.Printf("You must give two hole cards with -a\n")
		usage()
		os.Exit(1)
	}
	holeC, errIdx := strToCards(*hole)
	if (errIdx != -1) {
		fmt.Printf("Error parsing your hole cards: parse error at character %d\n",
					errIdx)
		os.Exit(1)
	}

	if (*verbose) {
		fmt.Printf("Your hole cards: %s\n", cardsToStr(holeC));
	}

	fmt.Printf("board = %s\n", *board)
	boardC, bErrIdx := strToCards(*board)
	if (bErrIdx != -1) {
		fmt.Printf("Error parsing the board: parse error at character %d\n",
					bErrIdx)
		os.Exit(1)
	}
	checkBoardLength(len(boardC))
	if (*verbose) {
		fmt.Printf("The board: %s\n", cardsToStr(boardC));
	}

	var c = make([]*card, len(boardC) + len(holeC))
	copy(c, boardC)
	copy(c[len(boardC):], holeC)
	dupe := hasDuplicates(c)
	if (dupe != nil) {
		fmt.Printf("The card %s appears more than once!\n", dupe)
		os.Exit(1)
	}

	// generate all cards
	var allCards [52]*card
	generateAllCards(&allCards)

	remCards := 5 - len(boardC)
	if (remCards < 0) {
		fmt.Printf("invalid board length %d\n", len(boardC))
		os.Exit(1)
	}

	choose 5 cards out of the hole + board + filler

	// There will be 2 hole cards and 5 cards on the board, making 7 cards to
	// choose from.
	card_chooser := pow64(2, 7) - 1

func do_combination(fixedCards []*card) {
	// what can we choose from in our 'random picks'?
	futureCards = allCards.deepcopy
	futureCards.subtract(fixedCards)

	numFutureCards = 5 - len(fixedCards)
	fiveCards := make([]*card, 5)
	copy(fiveCards, fixedCards)
	for i := 0; i < numFutureCards; i++ {
		
		print hand
	}
}

	do_combination(...)

	comb := pow64(2, uint(remCards)) - 1
	fmt.Printf("comb = %d, len(boardC)=%d\n", comb, len(boardC))
	for ;nextCombination(&comb); {
		handC := combinationToCards(comb, &allCards, &holeC, &boardC)
		var h = makeHand(handC)
		if (h != nil) {
			fmt.Printf("%s\n", h.String())
		}
	}

}
