package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strings"
)

func init() {
	Register(Day13)
}

func Day13(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	vals := []image.Point{}
	folds := []image.Point{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		p := image.Point{}
		Sscanf(line, "%d,%d", &p.X, &p.Y)
		vals = append(vals, p)
	}
	// Scan folds
	for scanner.Scan() {
		line := scanner.Text()

		var axis rune
		var n int
		Sscanf(line, "fold along %c=%d", &axis, &n)

		if axis == 'x' {
			folds = append(folds, image.Point{X: n})
		} else if axis == 'y' {
			folds = append(folds, image.Point{Y: n})
		}
	}

	firstFold := 0
	for i, f := range folds {
		vals = Fold(vals, f)
		if i == 0 {
			firstFold = len(vals)
		}
	}

	return fmt.Sprintf("%d,\n%s", firstFold, PrintPoints(vals))
}

func Fold(pts []image.Point, fold image.Point) []image.Point {
	unique := map[image.Point]struct{}{}
	if fold.X == 0 {
		// Fold up
		for _, p := range pts {
			if p.Y > fold.Y {
				p.Y = (fold.Y * 2) % p.Y
			}
			unique[p] = struct{}{}
		}
	} else if fold.Y == 0 {
		// Fold left
		for _, p := range pts {
			if p.X > fold.X {
				p.X = (fold.X * 2) % p.X
			}
			unique[p] = struct{}{}
		}
	}

	ret := []image.Point{}
	for p := range unique {
		ret = append(ret, p)
	}
	return ret
}

func PrintPoints(pts []image.Point) string {
	maxX, maxY := 0, 0
	for _, p := range pts {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	canvas := make([][]bool, maxY+1)
	for i := range canvas {
		canvas[i] = make([]bool, maxX+1)
	}

	for _, p := range pts {
		canvas[p.Y][p.X] = true
	}

	var out strings.Builder
	for i := range canvas {
		for j := range canvas[i] {
			if canvas[i][j] {
				out.WriteRune('#')
			} else {
				out.WriteRune(' ')
			}
		}
		out.WriteRune('\n')
	}

	return out.String()
}
