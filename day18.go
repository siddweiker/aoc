package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func init() {
	Register(Day18)
}

func Day18(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	pairs := []*Pair{}
	for scanner.Scan() {
		pairs = append(pairs, NewPairs(scanner.Bytes()))
	}

	a1, a2 := AddSnailFish(pairs)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func AddSnailFish(pairs []*Pair) (int, int) {
	var root *Pair
	for _, p := range pairs {
		clone := Clone(p)
		root = root.Add(clone)
	}
	magnitude := Magnitude(root)

	largest := 0
	for i := 0; i < len(pairs); i++ {
		for j := 0; j < len(pairs); j++ {
			if i == j {
				continue
			}

			pi := Clone(pairs[i])
			pj := Clone(pairs[j])
			m := Magnitude(pi.Add(pj))
			if m > largest {
				largest = m
			}
		}
	}

	return magnitude, largest
}

type Pair struct {
	Left   *Pair
	Right  *Pair
	Parent *Pair
	Value  int
}

func NewPairs(b []byte) *Pair {
	var root Pair
	err := json.Unmarshal(b, &root)
	if err != nil {
		log.Printf("Error parsing line '%s': %v", string(b), err)
		return nil
	}
	root.SetParents()
	return &root
}

func (p *Pair) Add(p2 *Pair) *Pair {
	if p == nil {
		return p2
	}
	root := &Pair{
		Left:  p,
		Right: p2,
	}
	root.Left.Parent = root
	root.Right.Parent = root
	root.SetParents()
	for exp, spl := true, true; exp || spl; {
		exp = root.Explode()
		spl = root.Split()
	}
	return root
}

// Returns true if it exploded
func (root *Pair) Explode() bool {
	end := GetBranch(root, 4)
	if end == nil {
		return false
	}

	// Make a list of neighbors that match the input
	// Increment the left and right neighbors by end.L and end.R
	neighbors := GetValues(root)
	for i, n := range neighbors {
		if end.Left == n && i > 0 {
			neighbors[i-1].Value += end.Left.Value
		}
		if end.Right == n && i < len(neighbors)-1 {
			neighbors[i+1].Value += end.Right.Value
		}
	}

	// Pairs are replaced by 0
	parent := end.Parent
	if parent.Left == end {
		parent.Left = &Pair{Value: 0, Parent: parent}
	}
	if parent.Right == end {
		parent.Right = &Pair{Value: 0, Parent: parent}
	}

	// See if we can explode once more
	root.Explode()

	return true
}

// Returns true if it split
func (root *Pair) Split() bool {
	split := false
	values := GetValues(root)
	for _, n := range values {
		if n.Value >= 10 {
			split = true
			*n = Pair{
				Left: &Pair{
					Value: n.Value / 2,
				},
				Right: &Pair{
					Value: n.Value/2 + (n.Value % 2),
				},
				Parent: n.Parent,
			}
			n.Left.Parent = n
			n.Right.Parent = n
			break
		}
	}

	return split
}

func (root *Pair) IsVal() bool {
	return root.Left == nil && root.Right == nil
}

func (root *Pair) IsBranch() bool {
	return root != nil && root.Left != nil && root.Right != nil
}

func (root *Pair) SetParents() {
	if root.Left != nil {
		root.Left.Parent = root
		root.Left.SetParents()
	}
	if root.Right != nil {
		root.Right.Parent = root
		root.Right.SetParents()
	}
}

func (root Pair) String() string {
	if root.Left == nil && root.Right == nil {
		return fmt.Sprintf("%d", root.Value)
	}

	return fmt.Sprintf("[%s,%s]", root.Left, root.Right)
}

func (root *Pair) UnmarshalJSON(b []byte) error {
	var tmp [2]json.RawMessage
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	p := Pair{}
	// Parse Left
	if err := json.Unmarshal(tmp[0], &p.Left); err != nil {
		var l Pair
		if err := json.Unmarshal(tmp[0], &l.Value); err != nil {
			return err
		}
		p.Left = &l
	}

	// Parse Right
	if err := json.Unmarshal(tmp[1], &p.Right); err != nil {
		var r Pair
		if err := json.Unmarshal(tmp[1], &r.Value); err != nil {
			return err
		}
		p.Right = &r
	}

	*root = p

	return nil
}

func Magnitude(root *Pair) int {
	if root == nil {
		return 0
	}

	if root.IsVal() {
		return root.Value
	}

	return 3*Magnitude(root.Left) + 2*Magnitude(root.Right)
}

func GetBranch(root *Pair, depth int) *Pair {
	if root == nil || depth == 0 {
		if root.IsBranch() {
			return root
		}
		return nil
	}

	if l := GetBranch(root.Left, depth-1); l != nil {
		return l
	}

	return GetBranch(root.Right, depth-1)
}

func GetValues(root *Pair) (vals []*Pair) {
	if root == nil {
		return nil
	}
	if root.IsVal() {
		vals = append(vals, root)
	}

	vals = append(vals, GetValues(root.Left)...)
	vals = append(vals, GetValues(root.Right)...)

	return
}

func Clone(root *Pair) *Pair {
	if root == nil {
		return nil
	}

	n := Pair{
		Left:  Clone(root.Left),
		Right: Clone(root.Right),
		Value: root.Value,
	}

	return &n
}
