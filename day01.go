package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day01)
}

var numbers = [...]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func Day01(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	a1, a2 := 0, 0
	for scanner.Scan() {
		n1, n2 := getEdgeInt(scanner.Text())
		a1 += n1
		a2 += n2
	}

	return fmt.Sprintf("%d, %d", a1, a2)
}

func getEdgeInt(line string) (int, int) {
	var left, right int
	// Score with number words included
	var leftW, rightW int
	for i, j := 0, len(line)-1; i < len(line); i, j = i+1, j-1 {
		if left == 0 {
			if n, err := byteToInt(line[i]); err == nil {
				left = n
				if leftW == 0 {
					leftW = n
				}
			}
		}
		if leftW == 0 {
			for v, word := range numbers {
				if strings.HasPrefix(line[i:], word) {
					leftW = v
				}
			}
		}

		if right == 0 {
			if n, err := byteToInt(line[j]); err == nil {
				right = n
				if rightW == 0 {
					rightW = n
				}
			}
		}
		if rightW == 0 {
			for v, word := range numbers {
				if strings.HasPrefix(line[j:], word) {
					rightW = v
				}
			}
		}
	}
	return left*10 + right, leftW*10 + rightW
}

func byteToInt(b byte) (int, error) {
	if b < '0' || b > '9' {
		return 0, errors.New("invalid int")

	}
	return int(b - 48), nil
}
