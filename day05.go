package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
)

func init() {
	Register(Day05)
}

func Day05(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	seeds := []int{}
	conv := []converter{}
	for scanner.Scan() {
		line := scanner.Text()
		if after, found := strings.CutPrefix(line, "seeds: "); found {
			for _, s := range strings.Fields(after) {
				seeds = append(seeds, Atoi(s))
			}
		} else if strings.Contains(line, "map:") {
			conv = append(conv, converter{struct {
				start  int
				offset int
			}{}})
		} else if line != "" {
			var dest, start, offset int
			Sscanf(line, "%d %d %d", &dest, &start, &offset)
			conv[len(conv)-1] = conv[len(conv)-1].add(start, start+offset-1, dest-start)
		}
	}

	seedRanges := [][2]int{}
	for i := range seeds {
		if i%2 == 0 {
			seedRanges = append(seedRanges, [2]int{seeds[i], seeds[i] + seeds[i+1]})
		}
	}

	a1, a2 := math.MaxInt, math.MaxInt
	for idx, c := range conv {
		for i := range seeds {
			seeds[i] = c.convert(seeds[i])
			if idx == len(conv)-1 {
				a1 = min(a1, seeds[i])
			}
		}
	}

	for _, seed := range seedRanges {
		for i := seed[0]; i <= seed[1]; i++ {
			n := i
			for _, c := range conv {
				n = c.convert(n)
			}
			a2 = min(a2, n)
		}
	}

	return fmt.Sprintf("%d, %d", a1, a2)

	// TODO: Fix this methodology
	// This should work, it gets really close to the answer
	fmt.Println(conv)
	fmt.Println(seedRanges)
	for idx, converter := range conv {
		for i := 0; i < len(seedRanges); i++ {
			minS, maxS := seedRanges[i][0], seedRanges[i][1]
			for j := 0; j < len(converter)-1; j++ {
				if minS >= converter[j].start && minS <= converter[j+1].start {
					// Within this range
					if maxS <= converter[j+1].start {
						break
					} else if maxS > converter[j+1].start {
						copySeed := seedRanges[i]
						seedRanges[i][1] = converter[j+1].start - 1
						copySeed[0] = converter[j+1].start

						seedRanges = slices.Insert(seedRanges, i+1, copySeed)
						// fmt.Println("  +", copySeed, seedRanges, "i: ", i)
						i += 1
					}
				}
			}
		}

		for i := range seedRanges {
			seedRanges[i][0] = converter.convert(seedRanges[i][0])
			seedRanges[i][1] = converter.convert(seedRanges[i][1])
			if idx == len(conv)-1 {
				a2 = min(a2, seedRanges[i][0], seedRanges[i][1])
			}
		}
		// fmt.Println(seedRanges, converter)
	}
	fmt.Println(seedRanges)

	return fmt.Sprintf("%d, %d", a1, a2)
}

type converter []struct {
	start, offset int
}

func (c converter) String() string {
	var out strings.Builder
	out.WriteString("[ ")
	for i := 0; i < len(c)-1; i++ {
		fmt.Fprintf(&out, "%d-%d:%d, ", c[i].start, c[i+1].start-1, c[i].offset)
	}
	fmt.Fprintf(&out, "%d+:%d", c[len(c)-1].start, c[len(c)-1].offset)
	out.WriteString(" ]")
	return out.String()
}

func (c converter) convert(item int) int {
	conv := item
	for _, v := range c {
		if v.start > item {
			break
		}
		conv = item + v.offset
	}
	return conv
}

func (c converter) add(start, end, offset int) converter {
	foundStart := false
	foundEnd := false
	// Insert start
	for i, v := range c {
		if start == v.start {
			foundStart = true
			break
		} else if start < v.start {
			c = slices.Insert(c, i, struct {
				start  int
				offset int
			}{
				start, 0,
			})
			foundStart = true
			break
		}
	}
	// Insert end
	for i, v := range c {
		if i == 0 {
			continue
		}
		if end+1 == v.start {
			foundEnd = true
			break
		} else if end+1 < v.start {
			c = slices.Insert(c, i, struct {
				start  int
				offset int
			}{
				end + 1, c[i-1].offset,
			})
			foundEnd = true
			break
		}
	}
	// Add offsets
	for i, v := range c {
		if v.start == start || (v.start > start && v.start < end+1) {
			c[i].offset += offset
		}
	}

	if !foundStart {
		c = append(c, struct {
			start  int
			offset int
		}{start, offset})
	}
	if !foundEnd {
		c = append(c, struct {
			start  int
			offset int
		}{end + 1, 0})
	}

	// Combine equal consecutive offsets
	for i := 0; i < len(c); i++ {
		if i == 0 {
			continue
		}
		if c[i-1].offset == c[i].offset {
			c = slices.Delete(c, i, i+1)
		}
	}

	return c
}
