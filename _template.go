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

	a1, a2 := 0, 0
	vals := []int{}
	for scanner.Scan() {
		line := scanner.Text()

		var num int
		Sscanf(line, "%d", &num)
		vals = append(vals, num)
	}

	return fmt.Sprintf("%d, %d", a1, a2)
}
