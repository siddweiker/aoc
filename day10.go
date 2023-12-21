package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"slices"
	"strings"
)

func init() {
	Register(Day10)
}

const SaveDay10Gif = false

func Day10(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	tiles := [][]byte{}
	starting := image.Point{}
	for scanner.Scan() {
		line := scanner.Text()
		tiles = append(tiles, []byte(line))
		if i := strings.Index(line, "S"); i >= 0 {
			starting = image.Point{X: len(tiles) - 1, Y: i}
		}
	}

	a1, a2 := parseField(tiles, starting)
	return fmt.Sprintf("%d, %d", a1, a2)
}

const (
	Unknown Direction = -1
	North   Direction = iota - 1
	East
	South
	West
)

var (
	fieldGif *gif.GIF
)

type Direction int

func (d Direction) String() string {
	switch d {
	case 0:
		return "North"
	case 1:
		return "East"
	case 2:
		return "South"
	case 3:
		return "West"
	}
	return "Unknown"
}

func ZeroPoint(p image.Point) bool {
	return p == image.Point{-1, -1}
}

func ComparePoints(a, b image.Point) int {
	if a.Eq(b) {
		return 0
	}
	if a.X > b.X {
		return 1
	} else if a.X == b.X {
		if a.Y > b.Y {
			return 1
		}
	}
	return -1
}

func parseField(tiles [][]byte, starting image.Point) (int, int) {
	createGif(len(tiles[0]), len(tiles))

	loop := []image.Point{starting}
	outer, inner := []image.Point{}, []image.Point{}

	// Find loop
	curr := image.Point{-1, -1}
	dir := Unknown
	for {
		if ZeroPoint(curr) {
			curr = starting
			appendGif(starting.Y, starting.X, color.RGBA{255, 255, 0, 255}, 10)
		}
		next, d := getNextTile(tiles, curr, dir)
		if next == starting {
			break
		}
		loop = append(loop, next)
		appendGif(next.Y, next.X, color.RGBA{0, 0, 0, 255}, 0)

		l, r, corner, bothLeft := getSides(tiles, next, d)
		if !corner {
			if !ZeroPoint(l) && !slices.Contains(outer, l) {
				outer = append(outer, l)
			}
			if !ZeroPoint(r) && !slices.Contains(inner, r) {
				inner = append(inner, r)
			}
		} else if bothLeft {
			if !ZeroPoint(l) && !slices.Contains(outer, l) {
				outer = append(outer, l)
			}
			if !ZeroPoint(r) && !slices.Contains(outer, r) {
				outer = append(outer, r)
			}
		} else {
			if !ZeroPoint(l) && !slices.Contains(inner, l) {
				inner = append(inner, l)
			}
			if !ZeroPoint(r) && !slices.Contains(inner, r) {
				inner = append(inner, r)
			}
		}

		curr, dir = next, d
	}

	// Fill false loops
	for x := range tiles {
		for y := range tiles[x] {
			if !slices.Contains(loop, image.Point{x, y}) {
				tiles[x][y] = '.'
			}
		}
	}

	// Sort and remove duplicates
	slices.SortFunc(outer, ComparePoints)
	slices.SortFunc(inner, ComparePoints)
	outer = slices.Compact(outer)
	inner = slices.Compact(inner)
	for _, n := range loop {
		outer = slices.DeleteFunc(outer, func(p image.Point) bool {
			return p == n
		})
		inner = slices.DeleteFunc(inner, func(p image.Point) bool {
			return p == n
		})
	}

	// Determine if "outer" is outside, else swap
	leftIsOuter := false
	for _, n := range outer {
		if n.X == 0 || n.Y == 0 || n.X == len(tiles)-1 || n.Y == len(tiles[0])-1 {
			leftIsOuter = true
			break
		}
	}
	if !leftIsOuter {
		outer, inner = inner, outer
	}

	// Set outer to 'O'
	for _, n := range outer {
		tiles[n.X][n.Y] = 'O'
	}

	// Set inner to 'I'
	for _, n := range inner {
		tiles[n.X][n.Y] = 'I'
		appendGif(n.Y, n.X, color.RGBA{255, 0, 0, 255}, 0)
	}

	// Populate empty inner spaces with a flood fill
	tiles = populateInner(tiles, loop, inner)

	inside := 0
	for x := range tiles {
		for y := range tiles[x] {
			if tiles[x][y] == 'I' {
				inside++
			}
		}
	}

	saveGif()

	return len(loop) / 2, inside
}

func getNeighbors(tiles [][]byte, curr image.Point) [4]image.Point {
	neighbors := [4]image.Point{
		{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1},
	}

	if curr.X > 0 {
		neighbors[North].X = curr.X - 1
		neighbors[North].Y = curr.Y
	}
	if curr.Y < len(tiles[curr.X])-1 {
		neighbors[East].X = curr.X
		neighbors[East].Y = curr.Y + 1
	}
	if curr.X < len(tiles)-1 {
		neighbors[South].X = curr.X + 1
		neighbors[South].Y = curr.Y
	}
	if curr.Y > 0 {
		neighbors[West].X = curr.X
		neighbors[West].Y = curr.Y - 1
	}

	return neighbors
}

func getNextTile(tiles [][]byte, curr image.Point, d Direction) (image.Point, Direction) {
	currTile := tiles[curr.X][curr.Y]
	next := image.Point{-1, -1}
	nextDir := d
	neighbors := getNeighbors(tiles, curr)

	// Starting position
	if d == Unknown {
		if n := neighbors[North]; !ZeroPoint(n) {
			if t := tiles[n.X][n.Y]; t == '|' || t == '7' || t == 'F' {
				return n, North
			}
		}
		if n := neighbors[East]; !ZeroPoint(n) {
			if t := tiles[n.X][n.Y]; t == '-' || t == 'J' || t == '7' {
				return n, East
			}
		}
		if n := neighbors[South]; !ZeroPoint(n) {
			if t := tiles[n.X][n.Y]; t == '|' || t == 'J' || t == 'L' {
				return n, South
			}
		}
		if n := neighbors[West]; !ZeroPoint(n) {
			if t := tiles[n.X][n.Y]; t == '-' || t == 'L' || t == 'F' {
				return n, West
			}
		}
	} else if d == South {
		// From above
		switch currTile {
		case '|':
			next = image.Point{curr.X + 1, curr.Y}
		case 'L':
			next = image.Point{curr.X, curr.Y + 1}
			nextDir = East
		case 'J':
			next = image.Point{curr.X, curr.Y - 1}
			nextDir = West
		}
	} else if d == North {
		// From below
		switch currTile {
		case '|':
			next = image.Point{curr.X - 1, curr.Y}
		case 'F':
			next = image.Point{curr.X, curr.Y + 1}
			nextDir = East
		case '7':
			next = image.Point{curr.X, curr.Y - 1}
			nextDir = West
		}
	} else if d == East {
		// From left
		switch currTile {
		case '-':
			next = image.Point{curr.X, curr.Y + 1}
		case '7':
			next = image.Point{curr.X + 1, curr.Y}
			nextDir = South
		case 'J':
			next = image.Point{curr.X - 1, curr.Y}
			nextDir = North
		}
	} else if d == West {
		// From right
		switch currTile {
		case '-':
			next = image.Point{curr.X, curr.Y - 1}
		case 'F':
			next = image.Point{curr.X + 1, curr.Y}
			nextDir = South
		case 'L':
			next = image.Point{curr.X - 1, curr.Y}
			nextDir = North
		}
	}

	return next, nextDir
}

func getSides(tiles [][]byte, curr image.Point, d Direction) (image.Point, image.Point, bool, bool) {
	neighbors := getNeighbors(tiles, curr)
	t := tiles[curr.X][curr.Y]
	first, second := Unknown, Unknown
	corner, bothLeft := false, false

	switch {
	case t == '|' && d == North:
		first, second = West, East
	case t == '|' && d == South:
		first, second = East, West
	case t == '-' && d == East:
		first, second = North, South
	case t == '-' && d == West:
		first, second = South, North
	case t == '7':
		first, second = North, East
		corner = true
		bothLeft = d == East
	case t == 'J':
		first, second = East, South
		corner = true
		bothLeft = d == South
	case t == 'L':
		first, second = South, West
		corner = true
		bothLeft = d == West
	case t == 'F':
		first, second = West, North
		corner = true
		bothLeft = d == North
	}

	return neighbors[first], neighbors[second], corner, bothLeft
}

func populateInner(tiles [][]byte, loop, inner []image.Point) [][]byte {
	seen := slices.Clone(loop)
	for _, in := range inner {
		if slices.Contains(seen, in) {
			continue
		}

		// Flood fill
		queue := []image.Point{in}
		for len(queue) > 0 {
			n := queue[0]
			queue = queue[1:]

			if slices.Contains(seen, n) {
				continue
			}
			seen = append(seen, n)

			if tiles[n.X][n.Y] == '.' {
				tiles[n.X][n.Y] = 'I'
				appendGif(n.Y, n.X, color.RGBA{255, 0, 0, 255}, 0)
			}

			for _, n := range getNeighbors(tiles, n) {
				if !ZeroPoint(n) {
					queue = append(queue, n)
				}
			}
		}
	}

	return tiles
}

func createGif(x, y int) {
	if !SaveDay10Gif {
		return
	}
	// Gif
	fieldGif = &gif.GIF{
		Image: []*image.Paletted{},
		Delay: []int{},
	}
	palette := []color.Color{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{128, 128, 128, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{173, 216, 230, 255},
		color.RGBA{255, 255, 255, 255},
	}
	img := image.NewPaletted(image.Rect(0, 0, x, y), palette)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)
	fieldGif.Image = append(fieldGif.Image, img)
	fieldGif.Delay = append(fieldGif.Delay, 0)
}

func appendGif(x, y int, c color.RGBA, delay int) {
	if !SaveDay10Gif {
		return
	}
	last := fieldGif.Image[len(fieldGif.Image)-1]
	img := image.NewPaletted(last.Bounds(), last.Palette)
	draw.Draw(img, img.Bounds(), last, image.Point{}, draw.Src)
	img.Set(x, y, c)
	fieldGif.Image = append(fieldGif.Image, img)
	fieldGif.Delay = append(fieldGif.Delay, delay)
}

func saveGif() {
	if !SaveDay10Gif {
		return
	}
	newImage := []*image.Paletted{fieldGif.Image[0]}
	newDelay := []int{fieldGif.Delay[0]}
	width := fieldGif.Image[len(fieldGif.Image)-1].Bounds().Max.X
	if width > 10 {
		for i := 1; i < len(fieldGif.Image); i += ((width * 2) / 10) {
			newImage = append(newImage, fieldGif.Image[i])
			newDelay = append(newDelay, fieldGif.Delay[i])
		}
	}
	newImage = append(newImage, fieldGif.Image[len(fieldGif.Image)-1])
	newDelay = append(newDelay, fieldGif.Delay[len(fieldGif.Delay)-1])
	fieldGif.Image = newImage
	fieldGif.Delay = newDelay
	fieldGif.Delay[len(fieldGif.Delay)-1] = 500

	f, err := os.Create("day10.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gif.EncodeAll(f, fieldGif)
}
