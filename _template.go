//go:build exclude

package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(DayXX)
}

func DayXX(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	vals := []int{}
	for scanner.Scan() {
		line := scanner.Text()

		var num int
		Sscanf(line, "%d", &num)
		vals = append(vals, num)
	}

	a1, a2 := 0, 0
	return fmt.Sprintf("%d, %d", a1, a2)
}
