package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day11)
}

func Day11(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	data := Cavern{}
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		for i, s := range line {
			data.octopi[lineNum][i] = uint8(Atoi(string(s)))
		}
		lineNum++
	}

	flashes := 0
	maxFlashes := len(data.octopi) * len(data.octopi[0])
	for i := 1; ; i++ {
		flashed := data.Step()
		if i == 100 {
			flashes = data.flashes
		}
		if flashed == maxFlashes {
			maxFlashes = i
			break
		}
	}
	return fmt.Sprintf("%d, %d", flashes, maxFlashes)
}

type Cavern struct {
	octopi  [10][10]uint8
	flashes int
}

func (c *Cavern) Step() int {
	// Add 1 energy to all
	for i := range c.octopi {
		for j := range c.octopi[i] {
			c.octopi[i][j]++
		}
	}
	// Propogate flashes
	for {
		if !c.flash() {
			break
		}
	}
	// Set flashed to 0
	flashes := 0
	for i := range c.octopi {
		for j := range c.octopi[i] {
			if c.octopi[i][j] > 9 {
				c.octopi[i][j] = 0
				flashes++
			}
		}
	}
	return flashes
}

func (c *Cavern) flash() bool {
	flashes := 0
	for i := range c.octopi {
		for j, octo := range c.octopi[i] {
			if octo == 10 {
				flashes++
				// Set to 11 so it doesn't flash again
				c.octopi[i][j]++

				if i > 0 {
					c.increment(i-1, j)
					// Top diagonals
					if j > 0 {
						c.increment(i-1, j-1)
					}
					if j < len(c.octopi[i])-1 {
						c.increment(i-1, j+1)
					}
				}
				if i < len(c.octopi)-1 {
					c.increment(i+1, j)
					// Bot diagonals
					if j > 0 {
						c.increment(i+1, j-1)
					}
					if j < len(c.octopi[i])-1 {
						c.increment(i+1, j+1)
					}
				}
				if j > 0 {
					c.increment(i, j-1)
				}
				if j < len(c.octopi[i])-1 {
					c.increment(i, j+1)
				}
			}
		}
	}

	c.flashes += flashes
	return flashes > 0
}

// increment adds one up to 10 max
func (c *Cavern) increment(i, j int) {
	if c.octopi[i][j] < 10 {
		c.octopi[i][j]++
	}
}

func (c Cavern) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "Flashes: %d\n", c.flashes)
	for i := range c.octopi {
		out.WriteRune('[')
		for j := range c.octopi[i] {
			flash := ' '
			if c.octopi[i][j] == 0 || c.octopi[i][j] > 9 {
				flash = '*'
			}

			fmt.Fprintf(&out, "%2d%c ", c.octopi[i][j], flash)
		}
		out.WriteString(" ]\n")
	}
	return out.String()
}
