package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

func init() {
	Register(Day10)
}

func Day10(r io.Reader) string {
	syntaxErrors := map[rune]int{}
	incompletes := []int{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		corrupt, incomplete := syntaxError(line)
		if corrupt != 0 {
			syntaxErrors[corrupt]++
		}
		if incomplete != 0 {
			incompletes = append(incompletes, incomplete)
		}
	}

	total := 0
	for r, n := range syntaxErrors {
		if r != rune(0) {
			total += n * scores[r]
		}
	}
	// Take the middle score from the sorted results
	sort.Ints(incompletes)
	total2 := incompletes[len(incompletes)/2]

	return fmt.Sprintf("%d, %d", total, total2)
}

var scores = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}
var scoresComplete = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

func syntaxError(line string) (corrupt rune, incomplete int) {
	stack := []rune{}
	for _, r := range line {
		if strings.ContainsRune("([{<", r) {
			stack = append(stack, r)
		} else if strings.ContainsRune(")]}>", r) {
			currOpen := opposite(r)
			lastIndex := len(stack) - 1
			if stack[lastIndex] != currOpen {
				corrupt = r
				break
			}
			stack = stack[:lastIndex]
		}
	}

	if corrupt == 0 && len(stack) != 0 {
		// Get the opposite of the stack
		for i := len(stack) - 1; i >= 0; i-- {
			// Calculate score: score * 5 + value of incomplete rune
			incomplete = incomplete*5 + scoresComplete[opposite(stack[i])]
		}
	}

	return

}

func opposite(r rune) rune {
	open := strings.IndexRune("([{<", r)
	close := strings.IndexRune(")]}>", r)

	if close > -1 {
		return rune("([{<"[close])
	} else if open > -1 {
		return rune(")]}>"[open])
	}
	return 0
}
