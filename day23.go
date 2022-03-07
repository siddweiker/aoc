package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day23)
}

func Day23(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	board := Amphipods{}
	i := 0
	for scanner.Scan() {
		if i == 2 || i == 3 {
			line := scanner.Text()
			line = strings.Trim(line, "# ")
			for j, s := range strings.Split(line, "#") {
				board.Rooms[j] = append(board.Rooms[j], Burrow(s[0]))
			}
		}
		i++
	}

	a1 := FindLeastEnergy(&board)
	board.Unfold()
	a2 := FindLeastEnergy(&board)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func FindLeastEnergy(d *Amphipods) int {
	if d.Sorted() {
		return d.Energy
	}

	movesIn, movesOut := d.Moves()
	if len(movesIn) == 0 && len(movesOut) == 0 {
		return 0
	}

	lowest := 0
	for _, m := range movesIn {
		new := d.Copy()
		new.Enter(m[0], m[1])
		e := FindLeastEnergy(new)
		if lowest == 0 || (e != 0 && e < lowest) {
			lowest = e
		}
	}
	for _, m := range movesOut {
		new := d.Copy()
		new.Leave(m[0], m[1])
		e := FindLeastEnergy(new)
		if lowest == 0 || (e != 0 && e < lowest) {
			lowest = e
		}
	}

	return lowest
}

type Burrow rune

type Amphipods struct {
	Hallway [11]Burrow
	Rooms   [4][]Burrow
	Energy  int
}

func (d *Amphipods) Moves() (in, out [][2]int) {
	// Move in
	for i, a := range d.Hallway {
		if a == 0 || d.AtDoor(i) {
			continue
		}

		// Move left and see if we can move into a room
		for left := i - 1; left > 1; left-- {
			if d.Hallway[left] != 0 {
				break
			}
			if !d.AtDoor(left) {
				continue
			}

			r := d.ToRoom(left)
			if d.EnterAllowed(r, a) {
				in = append(in, [2]int{i, r})
			}
		}
		// Move right and see if we can move into a room
		for right := i + 1; right < len(d.Hallway)-2; right++ {
			if d.Hallway[right] != 0 {
				break
			}
			if !d.AtDoor(right) {
				continue
			}

			r := d.ToRoom(right)
			if d.EnterAllowed(r, a) {
				in = append(in, [2]int{i, r})
			}
		}
	}
	if len(in) > 0 {
		return
	}

	// Move out
	for r, rooms := range d.Rooms {
		if d.RoomSorted(r) || !d.LeaveAllowed(r) {
			continue
		}

		for _, a := range rooms {
			if a == 0 {
				continue
			}

			// Move out
			o := d.ToHallway(r)
			// Check left positions
			for left := o - 1; left >= 0 && d.Hallway[left] == 0; left-- {
				if d.AtDoor(left) {
					continue
				}
				out = append(out, [2]int{r, left})
			}
			// Check right positions
			for right := o + 1; right < len(d.Hallway) && d.Hallway[right] == 0; right++ {
				if d.AtDoor(right) {
					continue
				}
				out = append(out, [2]int{r, right})
			}
			break
		}
	}

	return
}

func (d *Amphipods) Leave(r, h int) {
	cost := 1
	var amphi Burrow
	for i, a := range d.Rooms[r] {
		if a != 0 {
			amphi = a
			d.Rooms[r][i] = 0
			break
		}
		cost++
	}

	cost += Abs(d.ToHallway(r) - h)
	d.Hallway[h] = amphi
	d.Energy += cost * amphi.Cost()
}

func (d *Amphipods) Enter(h, r int) {
	cost := 1
	amphi := d.Hallway[h]
	d.Hallway[h] = 0
	cost += Abs(d.ToHallway(r) - h)

	for i := len(d.Rooms[r]) - 1; i >= 0; i-- {
		if d.Rooms[r][i] == 0 {
			d.Rooms[r][i] = amphi
			cost += i
			break
		}
	}

	d.Energy += cost * amphi.Cost()
}

func (d *Amphipods) ToHallway(r int) int {
	return (r + 1) * 2
}

func (d *Amphipods) ToRoom(h int) int {
	return h/2 - 1
}

func (d *Amphipods) EnterAllowed(r int, a Burrow) bool {
	if !d.IsHome(r, a) {
		return false
	}
	for _, ra := range d.Rooms[r] {
		if ra == 0 {
			continue
		} else if !d.IsHome(r, ra) {
			return false
		}
	}
	return true
}

func (d *Amphipods) LeaveAllowed(r int) bool {
	if d.RoomSorted(r) {
		return false
	}

	for _, ra := range d.Rooms[r] {
		if ra != 0 && !d.IsHome(r, ra) {
			return true
		}
	}
	return false
}

func (d *Amphipods) IsHome(r int, a Burrow) bool {
	switch a {
	case 'A':
		return r == 0
	case 'B':
		return r == 1
	case 'C':
		return r == 2
	case 'D':
		return r == 3
	}
	return false
}

func (d *Amphipods) AtDoor(h int) bool {
	return h == 2 || h == 4 || h == 6 || h == 8
}

func (d *Amphipods) Sorted() bool {
	for r := range d.Rooms {
		if !d.RoomSorted(r) {
			return false
		}
	}
	return true
}

func (d *Amphipods) RoomSorted(r int) bool {
	for _, a := range d.Rooms[r] {
		if !d.IsHome(r, a) {
			return false
		}
	}
	return true
}

// Surprise!
func (d *Amphipods) Unfold() {
	// Insert these amphipods in between existing
	// #D#C#B#A#
	// #D#B#A#C#
	d.Rooms[0] = append(d.Rooms[0][:1], append([]Burrow{'D', 'D'}, d.Rooms[0][1:]...)...)
	d.Rooms[1] = append(d.Rooms[1][:1], append([]Burrow{'C', 'B'}, d.Rooms[1][1:]...)...)
	d.Rooms[2] = append(d.Rooms[2][:1], append([]Burrow{'B', 'A'}, d.Rooms[2][1:]...)...)
	d.Rooms[3] = append(d.Rooms[3][:1], append([]Burrow{'A', 'C'}, d.Rooms[3][1:]...)...)
}

func (d *Amphipods) Copy() *Amphipods {
	n := &Amphipods{
		Hallway: d.Hallway,
		Energy:  d.Energy,
	}
	for r := range d.Rooms {
		n.Rooms[r] = make([]Burrow, len(d.Rooms[r]))
		copy(n.Rooms[r], d.Rooms[r])
	}
	return n
}

func (d Amphipods) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "\n%s\n#", strings.Repeat("#", 13))
	for _, a := range d.Hallway {
		fmt.Fprintf(&out, "%s", a)
	}

	fmt.Fprintf(&out, "#\n")
	for i := 0; i < len(d.Rooms[0]); i++ {
		if i == 0 {
			fmt.Fprintf(&out, "###")
		} else {
			fmt.Fprintf(&out, "  #")
		}
		for r := range d.Rooms {
			fmt.Fprintf(&out, "%s#", d.Rooms[r][i])
		}
		if i == 0 {
			fmt.Fprintf(&out, "##")
		}
		fmt.Fprintf(&out, "\n")
	}

	fmt.Fprintf(&out, "  %s", strings.Repeat("#", 9))
	return out.String()
}

func (b Burrow) Cost() int {
	switch b {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	}
	return 0
}

func (b Burrow) String() string {
	if b == 0 {
		return "."
	}
	return fmt.Sprintf("%c", b)
}
