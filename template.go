//go:build exclude
package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(DayX)
}

func DayX(r io.Reader) string {
	vals := []int{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		var num int
		Sscanf(line, "%d", &num)
		vals = append(vals, num)
	}

	return fmt.Sprintf("%d, %d", 0, 0)
}
