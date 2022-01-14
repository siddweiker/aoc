package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

func init() {
	Register(Day15)
}

// TODO: A big improvement would be to replace Graph Node maps by an array of length h*w
// Using the index as a key instead of "x,y" strings
func Day15(r io.Reader) string {
	cave := [][]uint8{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		l := make([]uint8, len(line))
		for i, r := range line {
			l[i] = uint8(r - 48)
		}
		cave = append(cave, l)
	}

	_, a1 := NewGraph(cave).ShortestPath()
	cave = ScaleCave(cave, 5)
	_, a2 := NewGraph(cave).ShortestPath()

	return fmt.Sprintf("%d, %d", a1, a2)
}

type Node struct {
	Key       string
	Value     uint8
	Neighbors map[*Node]uint8
}

type Graph struct {
	Start, End string
	Nodes      map[string]*Node
}

// *Min* Priority Queue, taken from container/heap and modified to be min instead of max
// An Item is something we manage in a priority queue.
type Item struct {
	node     *Node
	prev     *Node
	risk     int
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (graph *Graph) ShortestPath() ([]string, int) {
	start := graph.Nodes[graph.Start]
	end := graph.Nodes[graph.End]
	endX, endY := 0, 0
	Sscanf(graph.End, "%d,%d", &endX, &endY)
	heuristic := func(key string) int {
		x, y := 0, 0
		Sscanf(key, "%d,%d", &x, &y)
		first := math.Pow(float64(endX-x), 2)
		second := math.Pow(float64(endY-y), 2)
		return int(math.Sqrt(first + second))
	}

	queue := &PriorityQueue{}
	openItems := map[*Node]*Item{}
	seenItems := map[*Node]*Item{}

	item := &Item{start, nil, 0, 0, 0}
	openItems[start] = item
	heap.Push(queue, item)

	for queue.Len() > 0 {
		// Pop and move node from open to seen
		currentItem := heap.Pop(queue).(*Item)
		current := currentItem.node
		seenItems[current] = openItems[current]
		delete(openItems, current)

		// We reached the end
		if current == end {
			path := []string{}

			for current != nil {
				path = append(path, current.Key)
				current = seenItems[current].prev
			}

			// Reverse path
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path, currentItem.risk
		}

		risk := seenItems[current].risk

		for neighbor, val := range current.Neighbors {
			// Neighbor seen already
			if _, ok := seenItems[neighbor]; ok {
				continue
			}

			riskToNeighbor := risk + int(val)

			// Check to see if a seen node should be removed
			if notVisited, ok := openItems[neighbor]; ok {
				if notVisited.risk < riskToNeighbor {
					continue
				} else {
					heap.Remove(queue, notVisited.index)
				}
			}

			// Add neighbor
			item := &Item{
				neighbor,
				current,
				riskToNeighbor,
				riskToNeighbor + heuristic(neighbor.Key),
				0,
			}
			openItems[neighbor] = item
			heap.Push(queue, item)
		}
	}

	return nil, 0
}

func NewGraph(grid [][]uint8) *Graph {
	g := &Graph{
		"0,0",
		fmt.Sprintf("%d,%d", len(grid)-1, len(grid[len(grid)-1])-1),
		make(map[string]*Node),
	}

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			key := fmt.Sprintf("%d,%d", i, j)
			// Up, skip if we are on the last row
			if i > 0 && j != len(grid[i])-1 {
				g.Connect(key, fmt.Sprintf("%d,%d", i-1, j), grid[i-1][j])
			}
			// Down
			if i < len(grid)-1 {
				g.Connect(key, fmt.Sprintf("%d,%d", i+1, j), grid[i+1][j])
			}
			// Left, skip if we are on the last column
			if j > 0 && i != len(grid)-1 {
				g.Connect(key, fmt.Sprintf("%d,%d", i, j-1), grid[i][j-1])
			}
			// Right
			if j < len(grid[i])-1 {
				g.Connect(key, fmt.Sprintf("%d,%d", i, j+1), grid[i][j+1])
			}
		}
	}

	return g
}

func (g *Graph) Connect(a, b string, weight uint8) {
	if _, ok := g.Nodes[a]; !ok {
		g.Nodes[a] = &Node{Key: a, Neighbors: make(map[*Node]uint8)}
	}
	if _, ok := g.Nodes[b]; !ok {
		g.Nodes[b] = &Node{Key: b, Neighbors: make(map[*Node]uint8)}
	}

	g.Nodes[b].Value = weight
	g.Nodes[a].Neighbors[g.Nodes[b]] = weight
}

// ScaleCave will scale a grid by a scale factor
// The cave will be repeated x times in all directions growing each time
//
// Given the following cave:
// 123
// 456
// 789
// Scaled by 2 becomes:
// 123 | 234
// 456 | 567
// 789 | 891
// ----|----
// 234 | 345
// 567 | 678
// 891 | 912
func ScaleCave(grid [][]uint8, scale int) [][]uint8 {
	lenI := len(grid)
	lenJ := len(grid[0])
	newgrid := make([][]uint8, lenI*scale)
	for i := range newgrid {
		newgrid[i] = make([]uint8, lenJ*scale)
	}

	// Add a and b together.
	// The resulting number can only be from 1 to 9
	// Example: 9+1 = 1, 8+4 = 3
	add := func(a uint8, b int) uint8 {
		return (a+uint8(b)-1)%9 + 1
	}

	for si := 0; si < scale; si++ {
		for sj := 0; sj < scale; sj++ {
			for i := 0; i < lenI; i++ {
				for j := 0; j < lenJ; j++ {
					newgrid[i+si*lenI][j+sj*lenJ] = add(grid[i][j], si+sj)
				}
			}
		}
	}

	return newgrid
}

func (n Node) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s -> {", n.Key)
	for n, v := range n.Neighbors {
		fmt.Fprintf(&out, " %s:%d", n.Key, v)
	}
	out.WriteString(" }")
	return out.String()
}

func (g Graph) String() string {
	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var out strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&out, "%s\n", g.Nodes[k])
	}
	return out.String()
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use less than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
