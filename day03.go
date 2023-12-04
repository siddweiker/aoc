package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

func init() {
	Register(Day03)
}

func Day03(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	a1, a2 := engineSchematic(lines)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func isSymbol(r byte) bool {
	return r != '.' && (r < '0' || r > '9')
}

func notNumber(r byte) bool {
	return r < '0' || r > '9'
}

// Lines should be ascii
func engineSchematic(lines []string) (int, int) {
	sum := 0
	gears := map[int][]int{}
	for i, line := range lines {
		num := 0
		valid := false
		for j, c := range line {
			// Number end
			if notNumber(byte(c)) {
				if valid {
					sum += num
					findGears(lines, i, j-1, num, gears)
				}
				if num > 0 {
					num = 0
					valid = false
				}
				continue
			}

			// In number
			num = num*10 + int(c-48)
			if num > 0 && symbolNearby(lines, i, j) {
				valid = true
			}
		}

		// Check last number
		if valid {
			findGears(lines, i, len(line)-1, num, gears)
			sum += num
		}
	}

	// Find common gears with two numbers only
	gearSum := 0
	for _, nums := range gears {
		nums = slices.Compact(nums)
		if len(nums) == 2 {
			gearSum += nums[0] * nums[1]
		}
	}

	return sum, gearSum
}

func symbolNearby(lines []string, i, j int) bool {
	if len(lines) == 0 {
		return false
	}
	width, length := len(lines[0])-1, len(lines)-1
	// Above
	if x := i - 1; i > 0 {
		// NW
		if j > 0 && isSymbol(lines[x][j-1]) {
			return true
		}
		// N
		if isSymbol(lines[x][j]) {
			return true
		}
		// NE
		if j < width && isSymbol(lines[x][j+1]) {
			return true
		}
	}
	// Below
	if x := i + 1; i < length {
		// SW
		if j > 0 && isSymbol(lines[x][j-1]) {
			return true
		}
		// S
		if isSymbol(lines[x][j]) {
			return true
		}
		// SE
		if j < width && isSymbol(lines[x][j+1]) {
			return true
		}
	}
	// W
	if j > 0 && isSymbol(lines[i][j-1]) {
		return true
	}
	// E
	if j < width && isSymbol(lines[i][j+1]) {
		return true
	}

	return false
}

func findGears(lines []string, i, j, num int, gears map[int][]int) {
	width, length := len(lines[0])-1, len(lines)-1
	pos := j
	// E
	if pos < width && lines[i][pos+1] == '*' {
		gears[i*length+pos+1] = append(gears[i*length+pos+1], num)
	}
	for ; pos > 0; pos-- {
		// We have reached the start of the number
		if notNumber(lines[i][pos]) {
			pos += 1
			break
		}
		// Above
		if idx := (i-1)*length + pos; i > 0 {
			// NW
			if pos > 0 && lines[i-1][pos-1] == '*' {
				gears[idx-1] = append(gears[idx-1], num)
			}
			// N
			if lines[i-1][pos] == '*' {
				gears[idx] = append(gears[idx], num)
			}
			// NE
			if pos < width && lines[i-1][pos+1] == '*' {
				gears[idx+1] = append(gears[idx+1], num)
			}
		}
		// Below
		if idx := (i+1)*length + pos; i < length {
			// SW
			if pos > 0 && lines[i+1][pos-1] == '*' {
				gears[idx-1] = append(gears[idx-1], num)
			}
			// S
			if lines[i+1][pos] == '*' {
				gears[idx] = append(gears[idx], num)
			}
			// SE
			if pos < width && lines[i+1][pos+1] == '*' {
				gears[idx+1] = append(gears[idx+1], num)
			}
		}
	}
	// W
	if pos > 0 && lines[i][pos-1] == '*' {
		gears[i*length+pos-1] = append(gears[i*length+pos-1], num)
	}
}
