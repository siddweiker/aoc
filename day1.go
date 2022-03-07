package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(Day1)
}

func Day1(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	vals := []int{}
	rollingVals := []int{}
	for scanner.Scan() {
		line := scanner.Text()
		var num int
		Sscanf(line, "%d", &num)
		vals = append(vals, num)

		if l := len(vals); l > 2 {
			rollingVals = append(rollingVals, vals[l-3]+vals[l-2]+vals[l-1])
		}
	}

	a1, a2 := countGrowth(vals), countGrowth(rollingVals)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func countGrowth(l []int) int {
	incr, prev := 0, 0
	for i, n := range l {
		if i != 0 && n > prev {
			incr++
		}
		prev = n
	}
	return incr
}
