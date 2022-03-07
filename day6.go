package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day6)
}

func Day6(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	fish := LanternFish{}
	if scanner.Scan() {
		line := scanner.Text()

		for _, s := range strings.Split(line, ",") {
			fish.ages[Atoi(s)]++
		}
	}

	a1, a2 := Simulate(fish, 80, 256)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func Simulate(fish LanternFish, day1, day2 int) (int, int) {
	firstTotal := 0
	for d := 0; d < day2; d++ {
		if d+1 == day1 {
			firstTotal = fish.Total()
		}
		fish.Age()
	}
	return firstTotal, fish.Total()
}

type LanternFish struct {
	ages     [9]int
	spawning int
}

func (lf *LanternFish) Age() {
	// Advance the spawning day to "age" everyone down 1 day
	lf.spawning = (lf.spawning + 1) % 9
	// Calculate respawn and new(max) age based on the spawning day
	respawn := (lf.spawning + 6) % 9
	new := (lf.spawning + 8) % 9
	// Add new ages into the respawn age, the respawns(6) 'produced' the new ages
	// but also include previous days' age 7
	lf.ages[respawn] += lf.ages[new]
}

func (lf LanternFish) Total() int {
	total := 0
	for _, num := range lf.ages {
		total += num
	}
	return total
}

func (lf LanternFish) String() string {
	var out strings.Builder
	for i := 0; i < len(lf.ages); i++ {
		j := (lf.spawning + i) % 9
		fmt.Fprintf(&out, "%d,", lf.ages[j])
	}
	return fmt.Sprintf("Spawning Day: %d, [%s]", lf.spawning, out.String()[:out.Len()-1])
}
