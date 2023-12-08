package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func init() {
	Register(Day06)
}

func Day06(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	raceTimes, raceDistances := []int{}, []int{}
	fullTime, fullDistance := 0, 0
	for scanner.Scan() {
		line := scanner.Text()

		before, after, found := strings.Cut(line, ":")
		if !found {
			log.Printf("error parsing line '%s': ':' not found", line)
			continue
		}

		nums := []int{}
		for _, n := range strings.Fields(after) {
			nums = append(nums, Atoi(n))
		}
		full := Atoi(strings.ReplaceAll(after, " ", ""))
		if before == "Time" {
			raceTimes = nums
			fullTime = full
		} else {
			raceDistances = nums
			fullDistance = full
		}
	}

	a1, a2 := 1, 1
	for i := range raceTimes {
		a1 *= numOfWins(raceTimes[i], raceDistances[i])
	}
	a2 = numOfWins(fullTime, fullDistance)

	return fmt.Sprintf("%d, %d", a1, a2)
}

// TODO: I know theres a trick / formula to this...
// For now I figured out that start and end are not fast enough to win,
// so I find the slowest and fastest times that win, then subtract
// them for the number of wins
func numOfWins(time, distance int) int {
	wins := 0
	for i := time - 1; i > 0; i-- {
		if i*(time-i) > distance {
			wins = i
			break
		}
	}
	min := 0
	for i := 1; i < time; i++ {
		if i*(time-i) > distance {
			min = i
			break
		}
	}
	return wins - min + 1
}
