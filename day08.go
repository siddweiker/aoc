package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

func init() {
	Register(Day08)
}

func Day08(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	lines := []string{}
	instructions := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else if instructions == "" {
			instructions = line
		} else {
			lines = append(lines, line)
		}
	}

	a1, a2 := calculateSteps(instructions, lines)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func calculateSteps(instructions string, lines []string) (int, int) {
	nodes := make([]string, len(lines))
	left := make([]int, len(lines))
	right := make([]int, len(lines))
	ghosts := []int{}

	// Important: Sorting puts AAA at the start
	slices.Sort(lines)

	// Parse lines into nodes and ghosts
	for i, line := range lines {
		node, _, found := strings.Cut(line, " = ")
		if !found {
			log.Printf("error parsing line: '%s'", line)
			continue
		}
		nodes[i] = node
		if node[2] == 'A' {
			ghosts = append(ghosts, i)
		}
	}

	// Parse left right instructions as a list of node indexes
	for i, line := range lines {
		_, pair, found := strings.Cut(line, " = ")
		if !found {
			log.Printf("error parsing line: '%s'", line)
			continue
		}

		l, r, found := strings.Cut(pair, ", ")
		if !found {
			log.Printf("error parsing line: '%s'", line)
			continue
		}

		left[i] = slices.Index(nodes, strings.Trim(l, "()"))
		right[i] = slices.Index(nodes, strings.Trim(r, "()"))
	}

	// Minimum steps until each ghost hits a Z
	ghostSteps := make([]int, len(ghosts))
	// Amount of steps for AAA to get to ZZZ
	aSteps := 0
	for i, g := range ghosts {
		gSteps := 0
		for nodes[g][2] != 'Z' {
			if instructions[gSteps%len(instructions)] == 'L' {
				g = left[g]
			} else {
				g = right[g]
			}
			gSteps++
		}
		ghostSteps[i] = gSteps

		// Double check that AAA also ends in ZZZ
		if i == 0 {
			aSteps = gSteps
			for nodes[g] != "ZZZ" {
				if instructions[aSteps%len(instructions)] == 'L' {
					g = left[g]
				} else {
					g = right[g]
				}
				aSteps++
			}
		}
	}

	// Find a number that is divisible by every ghost step
	finalStep := lcmm(ghostSteps)
	return ghostSteps[0], finalStep
}

// lcmm returns the lowest common multiple between a list of ints
func lcmm(nums []int) int {
	if len(nums) == 2 {
		return lcm(nums[0], nums[1])
	} else {
		//nums = slices.Delete(nums, 0, 1)
		return lcm(nums[0], lcmm(nums[1:]))
	}
}

// lcm returns the lowest common multiple between a and b
func lcm(a, b int) int {
	return (a * b / gcd(a, b))
}

// gcd returns the greates common divisor using Euclid's Algorithm
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
