package main

import (
	"bufio"
	"fmt"
	"io"
)

func init() {
	Register(Day21)
}

func Day21(r io.Reader) string {
	dice := Dice{
		p1: Player{},
		p2: Player{},
	}
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		var num int
		Sscanf(line, "Player 1 starting position: %d", &num)
		dice.p1.pos = num
	}
	if scanner.Scan() {
		line := scanner.Text()
		var num int
		Sscanf(line, "Player 2 starting position: %d", &num)
		dice.p2.pos = num
	}

	a1 := dice.Deterministic()
	p1w, p2w := dice.Quantum()
	a2 := Max(p1w, p2w)

	return fmt.Sprintf("%d, %d", a1, a2)
}

type Player struct {
	pos   int
	score int
}

type Dice struct {
	p1, p2 Player
	// Used for deterministic
	last  int
	rolls int
	// Used for quantum
	p2turn bool
}

var allDimensionRolls = calculateAllDimensions()

func (p *Player) Move(steps int) {
	p.pos = ((p.pos + steps - 1) % 10) + 1
	p.score += p.pos
}

func (p *Player) Won() bool {
	return p.score >= 1000
}

func (p *Player) WonDirac() bool {
	return p.score >= 21
}

func (d Dice) Deterministic() int {
	losing := 0
	for {
		d.p1.Move(d.RollThrice())
		if d.p1.Won() {
			losing = d.p2.score
			break
		}
		d.p2.Move(d.RollThrice())
		if d.p2.Won() {
			losing = d.p1.score
			break
		}
	}
	return losing * d.rolls
}

// Rolls 3 times and returns the sum
func (d *Dice) RollThrice() int {
	return d.Roll() + d.Roll() + d.Roll()
}

func (d *Dice) Roll() int {
	d.rolls++
	if d.last == 100 {
		d.last = 1
	} else {
		d.last++
	}

	return d.last
}

func (d Dice) Quantum() (int, int) {
	p1wins, p2wins := 0, 0
	for roll, dims := range allDimensionRolls {
		// Copy the dice, including players
		new := Dice(d)

		// Check if the current player has won
		if !new.p2turn {
			new.p1.Move(roll)
			if new.p1.WonDirac() {
				p1wins += dims
				continue
			}
		} else {
			new.p2.Move(roll)
			if new.p2.WonDirac() {
				p2wins += dims
				continue
			}
		}

		// Next players turn
		new.p2turn = !new.p2turn
		// Continue playing (recurse)
		p1w, p2w := new.Quantum()

		// If a player won, multiply it by the number of this rolls' dimensions
		p1wins += p1w * dims
		p2wins += p2w * dims
	}

	return p1wins, p2wins
}

// All possibilities summed up of rolling a 3 sided dice 3 times
func calculateAllDimensions() map[int]int {
	dimensions := map[int]int{}
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				dimensions[i+j+k]++
			}
		}
	}
	return dimensions
}
