package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

func init() {
	Register(Day17)
}

func Day17(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	var target image.Rectangle
	if scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "..", ",")
		x0, y0, x1, y1 := 0, 0, 0, 0
		Sscanf(line, "target area: x=%d,%d, y=%d,%d", &x0, &x1, &y1, &y0)
		target = image.Rect(x0, y0+1, x1+1, y1)
	}

	highest, hits := 0, 0
	for x := 1; x <= target.Max.X; x++ {
		for y := target.Min.Y * 2; y < -target.Min.Y*2; y++ {
			_, maxY, miss := CalculatePath(image.Point{x, y}, target)

			if !miss {
				hits++
				if maxY > highest {
					highest = maxY
				}
			}

		}
	}

	return fmt.Sprintf("%d, %d", highest, hits)
}

func CalculatePath(vel image.Point, target image.Rectangle) ([]image.Point, int, bool) {
	if vel.In(target) {
		return []image.Point{{0, 0}, vel}, vel.Y, false
	}

	path := []image.Point{{0, 0}}

	for {
		new := path[len(path)-1].Add(vel)
		path = append(path, new)

		if new.In(target) {
			_, maxY, _ := MaxMin(path)
			return path, maxY, false
		}

		// Generous bounding box here
		if new.X > target.Max.X+target.Dx() || new.Y < target.Max.Y-target.Dy() {
			return path, 0, true
		}

		// Update velocity
		vel.Y -= 1
		if vel.X > 0 {
			vel.X -= 1
		} else if vel.X < 0 {
			vel.X += 1
		}
	}
}

func MaxMin(points []image.Point) (maxX, maxY, minY int) {
	for i, p := range points {
		if i == 0 {
			maxX, maxY, minY = p.X, p.Y, p.Y
			continue
		}

		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	return
}

func PointInPath(p image.Point, path []image.Point) bool {
	for _, pa := range path {
		if p == pa {
			return true
		}
	}
	return false
}

func PrintGrid(path []image.Point, target image.Rectangle) string {
	maxX, maxY, minY := MaxMin(path)
	maxX = Max(maxX, target.Max.X-1)
	minY = Min(minY, target.Min.Y)

	var out strings.Builder
	for x := 0; x <= maxX; x++ {
		fmt.Fprintf(&out, "%d", x%10)
	}
	out.WriteRune('\n')
	numLine := out.String()
	for y := maxY; y >= minY; y-- {
		for x := 0; x <= maxX; x++ {
			p := image.Point{x, y}
			if x == path[0].X && y == path[0].Y {
				out.WriteRune('S')
			} else if PointInPath(p, path) {
				out.WriteRune('#')
			} else if p.In(target) {
				out.WriteRune('T')
			} else {
				out.WriteRune('.')
			}
		}
		fmt.Fprintf(&out, " %2d", y)
		out.WriteRune('\n')
	}
	out.WriteString(numLine)

	return out.String()
}
