package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

func init() {
	Register(Day07)
}

func Day07(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	a1, a2 := 0, 0
	game := []Hand{}
	for scanner.Scan() {
		line := scanner.Text()
		c, bid, found := strings.Cut(line, " ")
		if !found {
			log.Printf("error parsing line '%s': could not split on space", line)
			continue
		}
		game = append(game, Hand{c, Atoi(bid)})
	}

	slices.SortFunc(game, func(a, b Hand) int {
		return CompareHands(a, b, false)
	})
	for i, h := range game {
		a1 += (i + 1) * h.Bid
	}

	slices.SortFunc(game, func(a, b Hand) int {
		return CompareHands(a, b, true)
	})
	for i, h := range game {
		a2 += (i + 1) * h.Bid
	}

	return fmt.Sprintf("%d, %d", a1, a2)
}

type Hand struct {
	Cards string
	Bid   int
}

func (h Hand) Type() int {
	count := [len(Cards)]int{}
	maxCard := 0
	pairs := 0
	for _, r := range h.Cards {
		i := CardIndex(byte(r))
		count[i]++
		maxCard = max(maxCard, count[i])
		if count[i] == 2 {
			pairs++
		} else if count[i] > 2 {
			pairs--
		}
	}

	switch maxCard {
	case 5:
		return 6
	case 4:
		return 5
	case 3:
		if pairs == 1 {
			return 4
		}
		return 3
	default:
		return pairs
	}
}

func (h Hand) TypeWithJoker() int {
	if !strings.ContainsRune(h.Cards, 'J') {
		return h.Type()
	}

	foundCards := []rune{}
	for _, r := range h.Cards {
		if r != 'J' && !slices.Contains(foundCards, r) {
			foundCards = append(foundCards, r)
		}
	}

	// Sort from most valuable to least
	slices.SortFunc(foundCards, func(a, b rune) int {
		if a == b {
			return 0
		}
		aVal, bVal := CardIndex(byte(a)), CardIndex(byte(b))
		if aVal > bVal {
			return -1
		} else if aVal < bVal {
			return 1
		}
		return 0
	})

	maxType := h.Type()
	for _, c := range foundCards {
		newCard := Hand{strings.ReplaceAll(h.Cards, "J", string(c)), h.Bid}
		maxType = max(maxType, newCard.Type())
	}
	return maxType
}

func CompareHands(a, b Hand, joker bool) int {
	aType, bType := 0, 0
	if joker {
		aType, bType = a.TypeWithJoker(), b.TypeWithJoker()
	} else {
		aType, bType = a.Type(), b.Type()
	}

	if aType == bType {
		// Compare each card in the hand
		if a.Cards == b.Cards {
			return 0
		}
		for i := range a.Cards {
			aVal, bVal := CardIndex(a.Cards[i]), CardIndex(b.Cards[i])
			if joker {
				if a.Cards[i] == 'J' {
					aVal = len(Cards)
				}
				if b.Cards[i] == 'J' {
					bVal = len(Cards)
				}
			}

			if aVal > bVal {
				return -1
			} else if aVal < bVal {
				return 1
			}
		}
		return 0
	}

	if aType < bType {
		return -1
	}
	return 1
}

var Cards = [...]byte{
	'A',
	'K',
	'Q',
	'J',
	'T',
	'9',
	'8',
	'7',
	'6',
	'5',
	'4',
	'3',
	'2',
}

func CardIndex(r byte) int {
	for i, card := range Cards {
		if r == card {
			return i
		}
	}
	return len(Cards)
}
