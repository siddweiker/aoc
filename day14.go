package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(Day14)
}

func Day14(r io.Reader) string {
	template := ""
	insertions := map[string]rune{}
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		template = scanner.Text()
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var rule string
		var char rune
		Sscanf(line, "%s -> %c", &rule, &char)
		insertions[rule] = char
	}

	a1, a2 := Polymerize(template, insertions, 10, 40)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func Polymerize(start string, rules map[string]rune, early, steps int) (int, int) {
	counts := map[rune]int{}
	for _, c := range start {
		counts[c]++
	}

	pairs := map[string]int{}
	for i := 0; i < len(start)-1; i++ {
		pairs[start[i:i+2]]++
	}

	earlyResult := 0
	for i := 0; i < steps; i++ {
		if i == early {
			earlyResult = MostMinusLeastRunes(counts)
		}
		pairs = poly(rules, pairs, counts)
	}

	return earlyResult, MostMinusLeastRunes(counts)
}

func poly(rules map[string]rune, pairs map[string]int, counts map[rune]int) map[string]int {
	newPairs := map[string]int{}
	for p, count := range pairs {
		c := rules[p]
		counts[c] += count
		newPairs[string(p[0])+string(c)] += count
		newPairs[string(c)+string(p[1])] += count
	}
	return newPairs
}

func MostMinusLeastRunes(counts map[rune]int) int {
	min, max := 0, 0
	for _, tot := range counts {
		if min == 0 && max == 0 {
			min, max = tot, tot
			continue
		}
		if max < tot {
			max = tot
		}
		if min > tot {
			min = tot
		}
	}
	return max - min
}
