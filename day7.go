package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day7)
}

func Day7(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	max := 0
	vals := []int{}
	if scanner.Scan() {
		line := scanner.Text()

		for _, s := range strings.Split(line, ",") {
			n := Atoi(s)
			vals = append(vals, n)
			if n > max {
				max = n
			}
		}
	}

	fuelCosts := []int{0}
	for i, curr := 0, 1; i <= max; i, curr = i+1, curr+i+1 {
		fuelCosts = append(fuelCosts, i+curr)
	}

	min := max * len(vals)
	min2 := fuelCosts[len(fuelCosts)-1] * len(vals)
	for i := 0; i <= len(vals); i++ {
		tot, tot2 := 0, 0
		for _, n := range vals {
			tot += Abs(n - i)
			tot2 += fuelCosts[Abs(n-i)]
		}
		if tot < min {
			min = tot
		}
		if tot2 < min2 {
			min2 = tot2
		}
	}

	return fmt.Sprintf("%d, %d", min, min2)
}
