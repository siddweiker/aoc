package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func init() {
	Register(Day04)
}

func Day04(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	a1, a2 := 0, 0
	cards := []int{}

	for scanner.Scan() {
		line := scanner.Text()

		_, after, found := strings.Cut(line, ": ")
		if !found {
			log.Printf("failed to parse line: %s", line)
			continue
		}
		winning, numbers, found := strings.Cut(after, " | ")
		if !found {
			log.Printf("failed to parse line: %s", line)
			continue
		}
		wins := findWinning(winning, numbers)
		cards = append(cards, wins)
		if wins > 0 {
			a1 += 1 << (wins - 1)
		}
	}
	a2 = processCards(cards)

	return fmt.Sprintf("%d, %d", a1, a2)
}

func findWinning(winning, numbers string) int {
	wins := 0
	for _, w := range strings.Fields(winning) {
		for _, n := range strings.Fields(numbers) {
			if w == n {
				wins += 1
			}
		}
	}
	return wins
}

func processCards(l []int) int {
	if len(l) == 0 {
		return 0
	}
	return 1 + processCards(l[1:l[0]+1]) + processCards(l[1:])
}
