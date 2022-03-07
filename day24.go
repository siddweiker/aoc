package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

func init() {
	Register(Day24)
}

func Day24(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	alu := &ALU{
		cmds: [][]Command{},
	}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) == 2 {
			alu.AddCmd(line[0], rune(line[1][0]), 0, 0)
		} else if d, err := strconv.Atoi(line[2]); err == nil {
			alu.AddCmd(line[0], rune(line[1][0]), 0, d)
		} else {
			alu.AddCmd(line[0], rune(line[1][0]), rune(line[2][0]), 0)
		}
	}

	a1, a2 := alu.Solve()
	return fmt.Sprintf("%d, %d", a2, a1)
}

type Command struct {
	cmd, r1, r2, val int
}

type ALU struct {
	cmds [][]Command
}

type Range struct {
	Min, Max int
}

type Variables [4]int

func (v *Variables) Get(i int) int {
	return v[i-1]
}

func (v *Variables) Set(i, val int) {
	v[i-1] = val
}

func (a *ALU) Solve() (min, max int) {
	// Go backwards through the commands, finding solutions for each step
	// We know the solution to the last step is z = 0
	solutions := a.Recurse(
		len(a.cmds)-1,
		map[int]*Range{0: {0, 0}},
	)

	for _, v := range solutions {
		if min == 0 || min > v.Min {
			min = v.Min
		}
		if max < v.Max {
			max = v.Max
		}
	}

	// Check answers
	vars := a.Run(max)
	if vars.Get(4) != 0 {
		log.Printf("Found solution does not validate: %d %s", max, vars)
	}
	vars = a.Run(min)
	if vars.Get(4) != 0 {
		log.Printf("Found solution does not validate: %d %s", min, vars)
	}

	return min, max
}

func (a *ALU) Recurse(step int, want map[int]*Range) map[int]*Range {
	// Potential optimization: at step 1 (or midway?) check all 9 numbers
	if step == -1 {
		return want
	}

	// Find possible answers
	found := map[int]*Range{}
	vars := Variables{}
	// Check each number 1-9
	for w := 1; w <= 9; w++ {
		// Find solutions for w
		// NOTE: This is a guess and may need to be increased
		for z := 0; z <= 1000000; z++ {
			vars.Set(4, z) // 4 == z
			a.RunStep(w, step, &vars)

			// Check if z is a previous solution
			if got, ok := want[vars.Get(4)]; ok {
				// Since w is from 1-9, first set min, the last value will be max
				if _, ok := found[z]; !ok {
					found[z] = &Range{
						Min: PadZeroes(w, 14-step) + got.Min,
						Max: PadZeroes(w, 14-step) + got.Max,
					}
				} else {
					found[z].Max = PadZeroes(w, 14-step) + got.Max
				}
			}
		}
	}

	return a.Recurse(step-1, found)
}

func (a *ALU) Run(input int) *Variables {
	digits := ToDigits(input)
	vars := &Variables{}
	for i := range a.cmds {
		a.RunStep(digits[i], i, vars)
	}
	return vars
}

func (a *ALU) RunStep(input, i int, vars *Variables) {
	for _, cmd := range a.cmds[i] {
		v1 := vars.Get(cmd.r1)
		v2 := cmd.val
		if cmd.r2 != 0 {
			v2 = vars.Get(cmd.r2)
		}

		switch cmd.cmd {
		case 0:
			vars.Set(cmd.r1, input)
		case 1:
			vars.Set(cmd.r1, v1+v2)
		case 2:
			vars.Set(cmd.r1, v1*v2)
		case 3:
			vars.Set(cmd.r1, v1/v2)
		case 4:
			vars.Set(cmd.r1, v1%v2)
		case 5:
			if v1 == v2 {
				vars.Set(cmd.r1, 1)
			} else {
				vars.Set(cmd.r1, 0)
			}
		}
	}
}

func (a *ALU) AddCmd(cmd string, r1, r2 rune, val int) {
	cmdI := cmdToIndex(cmd)
	// inp
	if cmdI == 0 {
		a.cmds = append(a.cmds, []Command{{cmdI, varToIndex(r1), 0, 0}})
		return
	}
	a.cmds[len(a.cmds)-1] = append(a.cmds[len(a.cmds)-1], Command{
		cmdI, varToIndex(r1), varToIndex(r2), val,
	})
}

func cmdToIndex(cmd string) int {
	switch cmd {
	case "inp":
		return 0
	case "add":
		return 1
	case "mul":
		return 2
	case "div":
		return 3
	case "mod":
		return 4
	case "eql":
		return 5
	default:
		log.Panicf("Invalid cmd: %s", cmd)
		return -1
	}
}

func varToIndex(r rune) int {
	switch r {
	case 'w':
		return 1
	case 'x':
		return 2
	case 'y':
		return 3
	case 'z':
		return 4
	default:
		return 0
	}
}

func ToDigits(n int) (digits []int) {
	for n > 0 {
		digits = append(digits, n%10)
		n = n / 10
	}
	// Reverse order
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return
}

func PadZeroes(n, index int) int {
	for ; index > 1; index-- {
		n *= 10
	}
	return n
}

func (v Variables) String() string {
	return fmt.Sprintf(
		"{ w:%d x:%d y:%d z:%d }",
		v[0], v[1], v[2], v[3],
	)
}

func (c Command) String() string {
	cmd := ""
	switch c.cmd {
	case 0:
		cmd = "inp"
	case 1:
		cmd = "add"
	case 2:
		cmd = "mul"
	case 3:
		cmd = "div"
	case 4:
		cmd = "mod"
	case 5:
		cmd = "eql"
	}
	reg := [5]rune{' ', 'w', 'x', 'y', 'z'}
	return fmt.Sprintf("{%s r1:%c r2:%c val:%d}", cmd, reg[c.r1], reg[c.r2], c.val)
}
