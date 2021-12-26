package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(Day9)
}

func Day9(r io.Reader) string {
	caves := heatmap{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		nums := []uint8{}
		for _, s := range line {
			var num uint8
			Sscanf(string(s), "%d", &num)
			nums = append(nums, num)
		}
		caves.nums = append(caves.nums, nums)
	}

	level, basinArea := caves.risk()
	return fmt.Sprintf("%d, %d", level, basinArea)
}

type heatmap struct {
	nums [][]uint8
	seen [][]bool
}

func (h heatmap) risk() (int, int) {
	level := 0
	areas := &[3]int{}
	// Find the lowest point in each basin
	for i := range h.nums {
		for j := range h.nums[i] {
			if i > 0 && h.nums[i][j] >= h.nums[i-1][j] {
				continue
			}
			if i < len(h.nums)-1 && h.nums[i][j] >= h.nums[i+1][j] {
				continue
			}
			if j > 0 && h.nums[i][j] >= h.nums[i][j-1] {
				continue
			}
			if j < len(h.nums[i])-1 && h.nums[i][j] >= h.nums[i][j+1] {
				continue
			}
			// Add lowest level score and calculate basin area
			level += int(h.nums[i][j]) + 1
			addLarger(areas, h.basins(i, j))
		}
	}

	return level, areas[0] * areas[1] * areas[2]
}

func (h heatmap) area(i, j int) {
	if h.nums[i][j] == 9 {
		return
	}
	seen := false
	// Up
	if i > 0 && h.nums[i-1][j] != 9 {
		seen = true
		if !h.seen[i-1][j] {
			defer h.area(i-1, j)
		}
	}
	// Down
	if i < len(h.nums)-1 && h.nums[i+1][j] != 9 {
		seen = true
		if !h.seen[i+1][j] {
			defer h.area(i+1, j)
		}
	}
	// Left
	if j > 0 && h.nums[i][j-1] != 9 {
		seen = true
		if !h.seen[i][j-1] {
			defer h.area(i, j-1)
		}
	}
	// Right
	if j < len(h.nums[i])-1 && h.nums[i][j+1] != 9 {
		seen = true
		if !h.seen[i][j+1] {
			defer h.area(i, j+1)
		}
	}

	h.seen[i][j] = seen
}

func (h heatmap) basins(x, y int) int {
	h.seen = make([][]bool, len(h.nums))
	for i := range h.nums {
		l := make([]bool, len(h.nums[i]))
		h.seen[i] = l
	}

	h.seen[x][y] = true
	h.area(x, y)

	area := 0
	for i := range h.seen {
		for j := range h.seen[i] {
			if h.seen[i][j] {
				area += 1
			}
		}
	}
	return area
}

func addLarger(dst *[3]int, n int) {
	if n > dst[0] {
		dst[0], dst[1], dst[2] = n, dst[0], dst[1]
	} else if n > dst[1] {
		dst[1], dst[2] = n, dst[1]
	} else if n > dst[2] {
		dst[2] = n
	}
}
