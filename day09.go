package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

func init() {
	Register(Day09)
}

func Day09(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	a1, a2 := 0, 0
	for scanner.Scan() {
		vals := []int{}
		for _, n := range strings.Fields(scanner.Text()) {
			vals = append(vals, Atoi(n))
		}
		a1 += predictNextValue(vals)
		slices.Reverse(vals)
		a2 += predictNextValue(vals)
	}

	return fmt.Sprintf("%d, %d", a1, a2)
}

func predictNextValue(vals []int) int {
	if slices.Max(vals) == slices.Min(vals) {
		return vals[0]
	}

	diff := make([]int, len(vals)-1)
	for i, n := range vals {
		if i == len(vals)-1 {
			break
		}
		diff[i] = vals[i+1] - n
	}

	return vals[len(vals)-1] + predictNextValue(diff)
}
