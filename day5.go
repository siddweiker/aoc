package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
)

func init() {
	Register(Day5)
}

func Day5(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	d := diagram{}
	d2 := diagram{}
	for scanner.Scan() {
		line := scanner.Text()

		x1, y1, x2, y2 := 0, 0, 0, 0
		Sscanf(line, "%d,%d -> %d,%d", &x1, &y1, &x2, &y2)

		d.StraightLine(x1, y1, x2, y2)
		d2.Line(x1, y1, x2, y2)
	}

	return fmt.Sprintf("%d, %d", d.Intersections(), d2.Intersections())
}

type diagram map[image.Point]int

func (d diagram) Set(x, y int) {
	d[image.Point{x, y}]++
}

func (d diagram) Intersections() int {
	total := 0
	for _, v := range d {
		if v > 1 {
			total++
		}
	}
	return total
}

func (d diagram) Line(x1, y1, x2, y2 int) {
	if x1 == x2 && y1 == y2 {
		d.Set(x1, y1)
		return
	}

	steps := Max(Abs(x1-x2), Abs(y1-y2))
	for i := 0; i <= steps; i++ {
		x, y := x1, y1
		if x1 < x2 {
			x += i
		} else if x1 > x2 {
			x -= i
		}
		if y1 < y2 {
			y += i
		} else if y1 > y2 {
			y -= i
		}
		d.Set(x, y)
	}
}

func (d diagram) StraightLine(x1, y1, x2, y2 int) {
	steps := Max(Abs(x1-x2), Abs(y1-y2))
	if x1 == x2 {
		if y1 > y2 {
			y1 = y2
		}
		for i := 0; i <= steps; i++ {
			d.Set(x1, y1+i)
		}
	} else if y1 == y2 {
		if x1 > x2 {
			x1 = x2
		}
		for i := 0; i <= steps; i++ {
			d.Set(x1+i, y1)
		}
	}
}
