package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(Day02)
}

func Day02(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	pos, depth := 0, 0
	aimdepth, aim := 0, 0
	for scanner.Scan() {
		x := 0
		direction := ""

		line := scanner.Text()
		Sscanf(line, "%s %d", &direction, &x)

		switch direction {
		case "forward":
			pos += x
			aimdepth += aim * x
		case "up":
			depth -= x
			aim -= x
		case "down":
			depth += x
			aim += x
		}
	}

	return fmt.Sprintf("%d, %d", pos*depth, pos*aimdepth)
}
