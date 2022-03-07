package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

func init() {
	Register(Day12)
}

func Day12(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	data := CaveGraph{CavePath{}, map[string]CavePath{}}
	for scanner.Scan() {
		line := scanner.Text()

		n1, n2, ok := Cut(line, "-")
		if !ok {
			log.Printf("error parsing line '%s'", line)
			continue
		}

		data.Add(n1, n2)
	}
	a1, a2 := data.AllPaths()
	return fmt.Sprintf("%d, %d", a1, a2)
}

type CaveGraph struct {
	nodes CavePath
	edges map[string]CavePath
}

type CavePath []string

type visit int

const (
	Unlimited visit = iota - 1
	NotVisited
	Once
	Twice
)

func (g *CaveGraph) Add(node, connectedNode string) {
	if connectedNode == "start" {
		node, connectedNode = connectedNode, node
	} else if node == "end" {
		node, connectedNode = connectedNode, node
	}
	n1, n2 := "", ""
	for _, s := range g.nodes {
		if s == node {
			n1 = s
		}
		if s == connectedNode {
			n2 = s
		}
	}
	if n1 == "" {
		n1 = node
		g.nodes = append(g.nodes, n1)
	}
	if n2 == "" {
		n2 = connectedNode
		g.nodes = append(g.nodes, n2)
	}

	g.edges[node] = append(g.edges[node], n2)
	if node != "start" && connectedNode != "end" {
		g.edges[connectedNode] = append(g.edges[connectedNode], n1)
	}
}

func (g CaveGraph) AllPaths() (int, int) {
	visited := map[string]visit{}
	for _, node := range g.nodes {
		if !IsLower(node) {
			visited[node] = Unlimited
		}
	}

	simple, twice := 0, 0
	f := func(n string, seenTwice bool) {
		if n == "end" {
			if !seenTwice {
				simple++
			}
			twice++
		}
	}
	g.DFS("start", visited, f)
	return simple, twice
}

func (g CaveGraph) DFS(start string, visited map[string]visit, f func(n string, t bool)) {
	if visited[start] != Unlimited {
		visited[start]++
	}

	seenTwice := false
	for k, v := range visited {
		if g.small(k) && v == Twice {
			seenTwice = true
			break
		}
	}
	f(start, seenTwice)

	for _, node := range g.edges[start] {
		if visited[node] == NotVisited || visited[node] == Unlimited {
			g.DFS(node, CopyVisited(visited), f)
		} else if visited[node] == Once && !seenTwice {
			g.DFS(node, CopyVisited(visited), f)
		}
	}
}

func CopyVisited(visited map[string]visit) map[string]visit {
	vc := map[string]visit{}
	for k, v := range visited {
		vc[k] = v
	}
	return vc
}

func (g CaveGraph) small(cave string) bool {
	return cave != "end" && cave != "start" && IsLower(cave)
}

func (g CaveGraph) String() string {
	var out strings.Builder
	for _, n := range g.nodes {
		fmt.Fprintf(&out, "%s ", n)
	}
	out.WriteRune('\n')
	for n1, conns := range g.edges {
		for _, e := range conns {
			fmt.Fprintf(&out, "%s -> %s\n", n1, e)
		}
	}
	return out.String()
}

func (p CavePath) String() string {
	var out strings.Builder
	out.WriteRune('[')
	for i, str := range p {
		if i == 0 {
			fmt.Fprintf(&out, "%s", str)

		} else {
			fmt.Fprintf(&out, " %s", str)
		}
	}
	out.WriteRune(']')
	return out.String()
}
