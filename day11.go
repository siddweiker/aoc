package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"slices"
)

func init() {
	Register(Day11)
}

func Day11(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	height, width := 0, 0
	galaxies := []image.Point{}
	for scanner.Scan() {
		line := scanner.Text()

		if width == 0 {
			width = len(line)
		}
		for x, c := range line {
			if c == '#' {
				galaxies = append(galaxies, image.Point{x, height})
			}
		}
		height++
	}

	older := expand(slices.Clone(galaxies), height, width, 1_000_000)
	galaxies = expand(galaxies, height, width, 1)

	a1, a2 := 0, 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			d := galaxies[i].Sub(galaxies[j])
			a1 += Abs(d.X) + Abs(d.Y)
			o := older[i].Sub(older[j])
			a2 += Abs(o.X) + Abs(o.Y)
		}
	}

	return fmt.Sprintf("%d, %d", a1, a2)
}

func expand(galaxies []image.Point, height, width, rate int) []image.Point {
	if rate > 1 {
		rate--
	}
	// Expand horizontally
	for x := 1; x < width; x++ {
		if !slices.ContainsFunc(galaxies, func(g image.Point) bool {
			return g.X == x
		}) {
			for i, g := range galaxies {
				if g.X > x {
					galaxies[i].X += rate
				}
			}
			width += rate
			x += rate
		}
	}
	// Expand vertically
	for y := 1; y < height; y++ {
		if !slices.ContainsFunc(galaxies, func(g image.Point) bool {
			return g.Y == y
		}) {
			for i, g := range galaxies {
				if g.Y > y {
					galaxies[i].Y += rate
				}
			}
			height += rate
			y += rate
		}
	}

	return galaxies
}
