package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day22)
}

func Day22(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	cubes := []Cube{}
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "..", ",")
		var op string
		x0, y0, z0 := 0, 0, 0
		x1, y1, z1 := 0, 0, 0
		Sscanf(line,
			"%s x=%d,%d,y=%d,%d,z=%d,%d",
			&op, &x0, &x1, &y0, &y1, &z0, &z1,
		)

		cubes = append(cubes, Cube{
			On: op == "on",
			S:  Point{x0, y0, z0},
			E:  Point{x1, y1, z1},
		})
	}

	a1, a2 := GetVolume(cubes)
	return fmt.Sprintf("%d, %d", a1, a2)
}

// Cuboid made out of cubes... technically
type Cube struct {
	On bool
	S  Point // Start
	E  Point // End
}

func GetVolume(cubes []Cube) (int, int) {
	all := []Cube{}

	for _, c := range cubes {
		processed := []Cube{}
		for _, a := range all {
			// Add the intersection if found
			if intersect, ok := a.Intersect(c); ok {
				processed = append(processed, intersect)
			}
		}
		// Add "on" cubes
		if c.On {
			processed = append(processed, c)
		}
		all = append(all, processed...)
	}

	initialize := 0
	min, max := Point{-50, -50, -50}, Point{50, 50, 50}
	volume := 0

	for _, c := range all {
		if c.S.Gte(min) && max.Gte(c.E) {
			initialize += c.Volume()
		}
		volume += c.Volume()
	}

	return initialize, volume
}

func (c *Cube) Intersect(c2 Cube) (Cube, bool) {
	// Create a cube with a the Max start and Min end
	in := Cube{
		On: c2.On,
		S: Point{
			Max(c.S.X, c2.S.X),
			Max(c.S.Y, c2.S.Y),
			Max(c.S.Z, c2.S.Z),
		},
		E: Point{
			Min(c.E.X, c2.E.X),
			Min(c.E.Y, c2.E.Y),
			Min(c.E.Z, c2.E.Z),
		},
	}

	// Cubes can intersect on a point
	// For example ON(x=11..13,y=11..13,z=11..13) + OFF(x=9..11,y=9..11,z=9..11) will turn {11,11,11} off
	// If start is larger than end, there is no intersection
	if in.S.X > in.E.X || in.S.Y > in.E.Y || in.S.Z > in.E.Z {
		return Cube{}, false
	}

	// This determines the "operation", two "on"s make an off
	if c.On && c2.On {
		in.On = false
	} else if !c.On && !c2.On {
		in.On = true
	}

	return in, true
}

func (c *Cube) Volume() int {
	// We add one to account for a cube *at* a point
	vol := (c.E.X - c.S.X + 1) * (c.E.Y - c.S.Y + 1) * (c.E.Z - c.S.Z + 1)
	if c.On {
		return vol
	}

	// If its off, negate the volume
	return -vol
}

func (p *Point) Gte(p2 Point) bool {
	return p.X >= p2.X && p.Y >= p2.Y && p.Z >= p2.Z
}
