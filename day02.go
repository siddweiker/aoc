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

var possibleCubes = cubes{12, 13, 14}

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
	maxCubes := cubes{}
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
			idx, err := possibleCubes.index(color)
			if err != nil {
				log.Printf("error: %s", err)
				continue
			}

			if n > possibleCubes[idx] {
				score = 0
			}
			maxCubes[idx] = max(n, maxCubes[idx])
		}
	}

	minScore := 1
	for _, v := range maxCubes {
		minScore *= v
	}
	return score, minScore
}

type cubes [3]int

func (c *cubes) index(color string) (int, error) {
	switch color {
	case "red":
		return 0, nil
	case "green":
		return 1, nil
	case "blue":
		return 2, nil
	}
	return -1, fmt.Errorf("color not found: %s", color)
}
