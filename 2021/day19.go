package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day19)
}

// Not my best code, I will revisit to improve the algorithm and code
func Day19(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	beacons := []*Scanner{}
	i := -1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else if strings.Contains(line, "---") {
			beacons = append(beacons, &Scanner{})
			i++
			continue
		}

		var p Point
		Sscanf(line, "%d,%d,%d", &p.X, &p.Y, &p.Z)
		beacons[i].Points = append(beacons[i].Points, p)
	}

	a1, a2 := AssembleScanners(beacons)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func AssembleScanners(beacons []*Scanner) (int, int) {
	links := map[*Scanner]*Scanner{}
	root := beacons[0]
	beacons = beacons[1:]
	// Find connections to root
	for i, b := range beacons {
		if Align(root, b) {
			links[b] = root
			// Pop beacon
			beacons[i] = beacons[len(beacons)-1]
			beacons = beacons[:len(beacons)-1]
		}
	}

	// Find connections to links
	for len(beacons) > 0 {
		for root := range links {
			for i, b := range beacons {
				if Align(root, b) {
					links[b] = root
					// Pop beacon
					beacons[i] = beacons[len(beacons)-1]
					beacons = beacons[:len(beacons)-1]
					break
				}
			}
		}
	}

	all := []Point{}
	all = append(all, root.Points...)
	sLocations := make([]Point, 0, len(links))
	// Traverse links and align beacons
	for k, v := range links {
		pts := k.Aligned
		sLoc := k.Diff

		for v != root {
			pts = v.AlignPoints(pts)
			sLoc = sLoc.Rotations()[v.Alignment].Add(v.Diff)
			v = links[v]
		}

		sLocations = append(sLocations, sLoc)
		all = append(all, pts...)
	}

	// Get unique points
	uniq := map[Point]struct{}{}
	allUniq := []Point{}
	for _, p := range all {
		if _, ok := uniq[p]; !ok {
			allUniq = append(allUniq, p)
			uniq[p] = struct{}{}
		}
	}

	// Find largest manhatten distance
	maxManhatten := 0
	for i := range sLocations {
		for j := range sLocations {
			if i == j {
				continue
			}
			p := sLocations[i].Sub(sLocations[j])
			dist := Abs(p.X) + Abs(p.Y) + Abs(p.Z)
			if dist > maxManhatten {
				maxManhatten = dist
			}
		}
	}

	return len(allUniq), maxManhatten
}

type Scanner struct {
	Points    []Point
	Aligned   []Point
	Diff      Point
	Alignment int
}

type Point struct {
	X, Y, Z int
}

func Align(sa *Scanner, sb *Scanner) bool {
	rotations := sb.Rotations()
	for r := 0; r < 24; r++ {
		diff := map[Point][]Point{}
		for i := 0; i < len(sa.Points); i++ {
			for j := 0; j < len(sb.Points); j++ {
				a := sa.Points[i]
				b := rotations[r][j]
				d := a.Sub(b)
				diff[d] = append(diff[d], a)
			}
		}

		for d, n := range diff {
			if len(n) >= 12 {
				newp := []Point{}
				for _, p := range rotations[r] {
					newp = append(newp, p.Add(d))
				}
				sb.Aligned = newp
				sb.Alignment = r
				sb.Diff = d
				return true
			}
		}
	}
	return false
}

func (s *Scanner) AlignPoints(pts []Point) []Point {
	for i, p := range pts {
		pts[i] = p.Rotations()[s.Alignment].Add(s.Diff)
	}
	return pts
}

func (s *Scanner) Rotations() (pts [24][]Point) {
	for _, p := range s.Points {
		for r, rp := range p.Rotations() {
			pts[r] = append(pts[r], rp)
		}
	}
	return pts
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func (p Point) Rotations() (pts [24]Point) {
	// Rotate the first point about the X,Y plane to each quadrant
	// For each quadrant add the inverse point
	c := p
	for i := 0; i < 12; i += 3 {
		pts[i].X, pts[i].Y, pts[i].Z = -c.Y, c.X, c.Z
		pts[i+12].X, pts[i+12].Y, pts[i+12].Z = -pts[i].X, -pts[i].Z, -pts[i].Y
		c = pts[i]
	}
	// Populate the remaining two values in each quadrant
	// Simply shift the values Left to Right
	// X2,Y2,Z2 = Y,Z,X ; X3,Y3,Z3 = Y2,Z2,X2
	for q := 1; q < 24; q += 3 {
		pts[q].X, pts[q].Y, pts[q].Z = pts[q-1].Y, pts[q-1].Z, pts[q-1].X
		pts[q+1].X, pts[q+1].Y, pts[q+1].Z = pts[q].Y, pts[q].Z, pts[q].X
	}

	return
}

func (p Point) Add(p2 Point) Point {
	p.X += p2.X
	p.Y += p2.Y
	p.Z += p2.Z
	return p
}

func (p Point) Sub(p2 Point) Point {
	p.X = p.X - p2.X
	p.Y = p.Y - p2.Y
	p.Z = p.Z - p2.Z
	return p
}
