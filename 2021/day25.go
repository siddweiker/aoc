package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

func init() {
	Register(Day25)
}

func Day25(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	vals := []string{}
	for scanner.Scan() {
		vals = append(vals, scanner.Text())
	}

	sf := NewSeaFloor(vals)
	a1 := 1
	for ; sf.Move(); a1++ {
	}

	return fmt.Sprintf("%d", a1)
}

type Seafloor struct {
	w, h int
	grid []*bool
	// Store all east and south values for easy lookup
	east  []*image.Point
	south []*image.Point
}

func NewSeaFloor(floor []string) *Seafloor {
	w := len(floor)
	h := len(floor[0])
	s := &Seafloor{
		w:    w,
		h:    h,
		grid: make([]*bool, w*h),
	}

	for x, line := range floor {
		for y, v := range line {
			switch v {
			case '>':
				b := true // east = true
				s.grid[s.index(x, y)] = &b
				s.east = append(s.east, &image.Point{x, y})
			case 'v':
				b := false // south = false
				s.grid[s.index(x, y)] = &b
				s.south = append(s.south, &image.Point{x, y})
			}
		}
	}

	return s
}

func (s *Seafloor) Move() bool {
	moved := false
	swap := []image.Point{}

	// Check all cucumbers moving east
	for _, curr := range s.east {
		next := s.next(true, *curr)
		nextI := s.index(next.X, next.Y)
		if s.grid[nextI] == nil {
			swap = append(swap, image.Point{
				s.index(curr.X, curr.Y),
				nextI,
			})
			curr.Y = next.Y
		}
	}

	if len(swap) > 0 {
		moved = true
	}
	// Perform swaps
	for _, p := range swap {
		s.grid[p.X], s.grid[p.Y] = s.grid[p.Y], s.grid[p.X]
	}

	swap = []image.Point{}

	// Check all cucumbers moving south
	for _, curr := range s.south {
		next := s.next(false, *curr)
		nextI := s.index(next.X, next.Y)
		if s.grid[nextI] == nil {
			swap = append(swap, image.Point{
				s.index(curr.X, curr.Y),
				nextI,
			})
			curr.X = next.X
		}
	}

	if len(swap) > 0 {
		moved = true
	}
	// Perform swaps
	for _, p := range swap {
		s.grid[p.X], s.grid[p.Y] = s.grid[p.Y], s.grid[p.X]
	}

	return moved
}

// Use a 1d array to represent a 2d grid
func (s *Seafloor) index(x, y int) int {
	return x*s.h + y
}

// Move a cucumber in the given direction, if at the end, loop back around
func (s *Seafloor) next(east bool, p image.Point) image.Point {
	if east {
		p.Y++
	} else {
		p.X++
	}
	if p.Y > s.h-1 {
		p.Y = 0
	}
	if p.X > s.w-1 {
		p.X = 0
	}
	return p
}

func (s Seafloor) String() string {
	var out strings.Builder
	out.WriteRune('\n')
	for x := 0; x < s.w; x++ {
		for y := 0; y < s.h; y++ {
			if s.grid[s.index(x, y)] == nil {
				out.WriteRune('.')
				continue
			}
			v := *s.grid[s.index(x, y)]
			if v {
				out.WriteRune('>')
			} else {
				out.WriteRune('v')
			}
		}
		out.WriteRune('\n')
	}
	return out.String()
}
