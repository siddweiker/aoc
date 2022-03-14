package main

import (
	"bufio"
	"container/heap"
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
				board.Rooms[j] = append(board.Rooms[j], rune(s[0]))
			}
		}
		i++
	}

	a1 := FindLeastEnergyAStar(&board)
	board.Unfold()
	a2 := FindLeastEnergyAStar(&board)
	return fmt.Sprintf("%d, %d", a1, a2)
}

// A* Solution
func FindLeastEnergyAStar(start *Amphipods) int {
	queue := &PriorityQueueA{}
	openItems := map[string]*ItemA{}
	seenItems := map[string]*ItemA{}

	item := &ItemA{start, 0, 0}
	openItems[start.ID()] = item
	heap.Push(queue, item)

	for queue.Len() > 0 {
		currentItem := heap.Pop(queue).(*ItemA)
		current := currentItem.board
		seenItems[current.ID()] = openItems[current.ID()]
		delete(openItems, current.ID())

		if current.Sorted() {
			return currentItem.board.Energy
		}

		movesIn, movesOut := current.Moves()
		if len(movesIn) == 0 && len(movesOut) == 0 {
			continue
		}

		// Create neighbors
		neighbors := []*Amphipods{}
		for _, m := range movesIn {
			new := current.Copy()
			new.Enter(m[0], m[1])
			neighbors = append(neighbors, new)
		}
		for _, m := range movesOut {
			new := current.Copy()
			new.Leave(m[0], m[1])
			neighbors = append(neighbors, new)
		}

		for _, neighbor := range neighbors {
			if _, ok := seenItems[neighbor.ID()]; ok {
				continue
			}

			if notVisited, ok := openItems[neighbor.ID()]; ok {
				if notVisited.board.Energy < neighbor.Energy {
					continue
				} else {
					heap.Remove(queue, notVisited.index)
				}
			}

			item := &ItemA{neighbor, neighbor.Energy, 0}
			openItems[neighbor.ID()] = item
			heap.Push(queue, item)
		}
	}

	return 0
}

// Naive Solution, checks every path
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

type Amphipods struct {
	Hallway [11]rune
	Rooms   [4][]rune
	Energy  int
}

func (d *Amphipods) ID() string {
	var out strings.Builder
	for r := 0; r < 4; r++ {
		for _, c := range d.Rooms[r] {
			out.WriteRune(c)
		}
	}
	for _, c := range d.Hallway {
		out.WriteRune(c)
	}
	return out.String()
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
	var amphi rune
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
	d.Energy += cost * Cost(amphi)
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

	d.Energy += cost * Cost(amphi)
}

func (d *Amphipods) ToHallway(r int) int {
	return (r + 1) * 2
}

func (d *Amphipods) ToRoom(h int) int {
	return h/2 - 1
}

func (d *Amphipods) EnterAllowed(r int, a rune) bool {
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

func (d *Amphipods) IsHome(r int, a rune) bool {
	return int(a)-65-r == 0
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
	d.Rooms[0] = append(d.Rooms[0][:1], append([]rune{'D', 'D'}, d.Rooms[0][1:]...)...)
	d.Rooms[1] = append(d.Rooms[1][:1], append([]rune{'C', 'B'}, d.Rooms[1][1:]...)...)
	d.Rooms[2] = append(d.Rooms[2][:1], append([]rune{'B', 'A'}, d.Rooms[2][1:]...)...)
	d.Rooms[3] = append(d.Rooms[3][:1], append([]rune{'A', 'C'}, d.Rooms[3][1:]...)...)
}

func (d *Amphipods) Copy() *Amphipods {
	n := &Amphipods{
		Hallway: d.Hallway,
		Energy:  d.Energy,
	}
	for r := range d.Rooms {
		n.Rooms[r] = make([]rune, len(d.Rooms[r]))
		copy(n.Rooms[r], d.Rooms[r])
	}
	return n
}

func (d Amphipods) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "\n%s\n#", strings.Repeat("#", 13))
	for _, a := range d.Hallway {
		if a == 0 {
			out.WriteRune('.')
		} else {
			out.WriteRune(a)
		}
	}

	fmt.Fprintf(&out, "#\n")
	for i := 0; i < len(d.Rooms[0]); i++ {
		if i == 0 {
			fmt.Fprintf(&out, "###")
		} else {
			fmt.Fprintf(&out, "  #")
		}
		for r := range d.Rooms {
			fmt.Fprintf(&out, "%c#", d.Rooms[r][i])
		}
		if i == 0 {
			fmt.Fprintf(&out, "##")
		}
		fmt.Fprintf(&out, "\n")
	}

	fmt.Fprintf(&out, "  %s", strings.Repeat("#", 9))
	return out.String()
}

func Cost(b rune) int {
	switch b {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	default:
		return 0
	}
}

// *Min* Priority Queue, taken from container/heap and modified to be min instead of max
// An Item is something we manage in a priority queue.
type ItemA struct {
	board    *Amphipods
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueueA implements heap.Interface and holds Items.
type PriorityQueueA []*ItemA

func (pq PriorityQueueA) Len() int { return len(pq) }

func (pq PriorityQueueA) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use less than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueueA) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueA) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ItemA)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueueA) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
