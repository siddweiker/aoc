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
	vals := []int{}
	rollingVals := []int{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		var num int
		Sscanf(line, "%d", &num)
		vals = append(vals, num)

		if l := len(vals); l > 2 {
			rollingVals = append(rollingVals, vals[l-3]+vals[l-2]+vals[l-1])
		}
	}

	return fmt.Sprintf("%d, %d", countGrowth(vals), countGrowth(rollingVals))
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
