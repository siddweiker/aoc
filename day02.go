package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func init() {
	Register(Day02)
}

var possibleCubes = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func Day02(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	a1, a2 := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		n1, n2 := parseCubeGame(line)
		a1 += n1
		a2 += n2
	}
	return fmt.Sprintf("%d, %d", a1, a2)
}

func parseCubeGame(line string) (int, int) {
	var score int
	maxCubes := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	game, pulls, found := strings.Cut(line, ": ")
	if !found {
		log.Printf("error parsing line: %s: could not split on ': '", line)
		return 0, 0
	}
	_, err := fmt.Sscanf(game, "Game %d", &score)
	if err != nil {
		log.Printf("error parsing line: %s: %s", line, err)
		return 0, 0
	}

	for _, pull := range strings.Split(pulls, "; ") {
		for _, cube := range strings.Split(pull, ", ") {
			var n int
			var color string
			_, err := fmt.Sscanf(cube, "%d %s", &n, &color)
			if err != nil {
				log.Printf("error parsing cube: %s: %s", cube, err)
				continue
			}

			if v, ok := possibleCubes[color]; ok {
				if n > v {
					score = 0
				}
				maxCubes[color] = max(n, maxCubes[color])
			} else {
				log.Printf("error color not found: %s: %s", color, err)
				continue
			}
		}
	}

	minScore := 1
	for _, v := range maxCubes {
		minScore *= v
	}
	return score, minScore
}
