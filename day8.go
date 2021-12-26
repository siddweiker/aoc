package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

func init() {
	Register(Day8)
}

func Day8(r io.Reader) string {
	vals := []signalEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		in, out, ok := Cut(line, " | ")
		if !ok {
			log.Printf("Failed to split line on ' | ': %s", line)
			continue
		}
		se := signalEntry{}
		inF := strings.Fields(in)
		outF := strings.Fields(out)
		for i, str := range inF {
			s := []rune(str)
			sort.Slice(s, func(i int, j int) bool { return s[i] < s[j] })
			inF[i] = string(s)
		}
		for i, str := range outF {
			s := []rune(str)
			sort.Slice(s, func(i int, j int) bool { return s[i] < s[j] })
			outF[i] = string(s)
		}
		copy(se.unique[:], inF)
		copy(se.output[:], outF)
		vals = append(vals, se)
	}

	total, total2 := 0, 0
	for _, v := range vals {
		total += v.countEasyNums()
		total2 += v.outputNum()
	}

	return fmt.Sprintf("%d, %d", total, total2)
}

type signalEntry struct {
	unique [10]string
	output [4]string
}

func (se signalEntry) countEasyNums() int {
	total := 0
	for _, s := range se.output {
		switch len(s) {
		case 2, 4, 3, 7:
			total++
		}
	}
	return total
}

func (se signalEntry) outputNum() int {
	nums := [10]string{}
	for _, s := range se.unique {
		switch len(s) {
		case 2:
			nums[1] = s
		case 4:
			nums[4] = s
		case 3:
			nums[7] = s
		case 7:
			nums[8] = s
		}
	}

	// Find strings that don't contain each ones character
	onesFirstChar := 0
	onesSecondChar := 0
	for _, s := range se.unique {
		if !strings.ContainsRune(s, rune(nums[1][0])) {
			onesFirstChar++
		} else if !strings.ContainsRune(s, rune(nums[1][1])) {
			onesSecondChar++
		}
	}
	c := rune(nums[1][0])
	f := rune(nums[1][1])
	// Only 1 word does not have f, 2 do not have c
	if onesFirstChar < onesSecondChar {
		// swap
		c, f = f, c
	}

	for _, s := range se.unique {
		if len(s) == 5 {
			if !strings.ContainsRune(s, f) {
				nums[2] = s
			} else if !strings.ContainsRune(s, c) {
				nums[5] = s
			} else {
				nums[3] = s
			}
		} else if len(s) == 6 {
			if !strings.ContainsRune(s, c) {
				nums[6] = s
			} else if ContainsAll(s, nums[4]) {
				nums[9] = s
			} else {
				nums[0] = s
			}
		}
	}

	var out strings.Builder
	for _, s := range se.output {
		for i, n := range nums {
			if s == n {
				fmt.Fprintf(&out, "%d", i)
			}
		}
	}
	return Atoi(out.String())
}
