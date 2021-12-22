package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func init() {
	Register(Day2)
}

func Day2(r io.Reader) string {
	pos, depth := 0, 0
	aimdepth, aim := 0, 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		x := 0
		direction := ""

		line := scanner.Text()
		_, err := fmt.Sscanf(line, "%s %d", &direction, &x)
		if err != nil {
			log.Printf("error parsing line '%s': %v", line, err)
			continue
		}

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
