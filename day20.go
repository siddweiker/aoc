package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day20)
}

func Day20(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	enhance := [512]bool{}
	grid := Floor{}
	if scanner.Scan() {
		for i, c := range scanner.Text() {
			if c == '#' {
				enhance[i] = true
			}
		}
		scanner.Scan()
	}

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]bool, len(line))
		for i, c := range line {
			if c == '#' {
				row[i] = true
			}
		}

		grid.lights = append(grid.lights, row)
	}

	a1, a2 := Enhance(grid, enhance)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func Enhance(grid Floor, enhance [512]bool) (int, int) {
	nextLit := enhance[0]
	litTwice := 0
	for i := 1; i <= 50; i++ {
		grid = grid.Enhance(enhance)
		grid.infinitLit = nextLit

		// The end values are important as the infinite grid not lit is index 0
		// and all lit is index 511 (2^9-1). This converts the infinite image based
		// on the algo. The infinite space starts as !lit, examples:
		// Algo = [#to.]; !lit -> lit, lit -> !lit - alternates
		// Algo = [.to#]; !lit -> !lit, !lit -> !lit - stayes not lit
		if !nextLit {
			nextLit = enhance[0]
		} else {
			nextLit = enhance[511]
		}

		if i == 2 {
			litTwice = grid.Lit()
		}
	}

	return litTwice, grid.Lit()
}

type Floor struct {
	lights     [][]bool
	infinitLit bool
}

func (f Floor) Enhance(algo [512]bool) Floor {
	// Grow the grid outwards by 1 unit
	lights := f.lights
	rowLen := len(lights[0]) + 2
	for i, row := range lights {
		row = append([]bool{f.infinitLit}, row...)
		row = append(row, f.infinitLit)
		lights[i] = row
	}
	// Insert an empty row at the start
	empty := []bool{}
	for i := 0; i < rowLen; i++ {
		empty = append(empty, f.infinitLit)
	}
	lights = append([][]bool{empty}, lights...)
	// Append an empty row at the end
	empty = []bool{}
	for i := 0; i < rowLen; i++ {
		empty = append(empty, f.infinitLit)
	}
	f.lights = append(lights, empty)

	// Our new Floor
	new := make([][]bool, len(f.lights))
	for i := range new {
		new[i] = make([]bool, rowLen)
	}

	// For each light, determine its new value
	for x := 0; x < len(new); x++ {
		for y := 0; y < len(new[x]); y++ {
			// Convert a 9x9 grid into a binary number
			index := [9]bool{
				f.Get(x-1, y-1), f.Get(x-1, y), f.Get(x-1, y+1),
				f.Get(x, y-1), f.Get(x, y), f.Get(x, y+1),
				f.Get(x+1, y-1), f.Get(x+1, y), f.Get(x+1, y+1),
			}
			var i uint16
			for _, lit := range index {
				i <<= 1
				if lit {
					i |= 1
				}
			}

			// Use the binary number as decimal to get the new value
			new[x][y] = algo[int(i)]
		}
	}

	return Floor{lights: new, infinitLit: f.infinitLit}
}

func (f Floor) Lit() int {
	total := 0
	for _, row := range f.lights {
		for _, lit := range row {
			if lit {
				total++
			}
		}
	}
	return total
}

func (f Floor) Get(x, y int) bool {
	if x < 0 || x >= len(f.lights) || y < 0 || y >= len(f.lights[x]) {
		return f.infinitLit
	}
	return f.lights[x][y]
}

func (f Floor) String() string {
	var out strings.Builder
	for _, row := range f.lights {
		for _, light := range row {
			if light {
				out.WriteRune('#')
			} else {
				out.WriteRune('.')
			}
		}
		out.WriteRune('\n')
	}
	return out.String()
}
