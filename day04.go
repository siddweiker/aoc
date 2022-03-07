package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func init() {
	Register(Day4)
}

func Day4(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	boards := []*Board{}
	drawed := []int{}
	for scanner.Scan() {
		line := scanner.Text()

		if len(drawed) == 0 {
			for _, s := range strings.Split(line, ",") {
				drawed = append(drawed, Atoi(s))
			}
			continue
		}

		if line == "" {
			b := [5][5]int{}
			total := 0
			for i := 0; i < 5; i++ {
				if !scanner.Scan() {
					break
				}

				for j, s := range strings.Fields(scanner.Text()) {
					b[i][j] = Atoi(s)
					total += b[i][j]
				}
			}
			boards = append(boards, &Board{
				numbers:   b,
				remaining: total,
			})
		}
	}

	a1, a2 := PlayBingo(boards, drawed)
	return fmt.Sprintf("%d, %d", a1, a2)
}

func PlayBingo(boards []*Board, drawed []int) (bingo, lastBingo int) {
	for i, n := range drawed {
		for _, b := range boards {
			b.Mark(n)
			if i+1 < 4 {
				// No one can win in less than 4 draws
				continue
			}

			if b.won {
				continue
			}

			if b.Bingo() {
				//log.Printf("Board #%d got BINGO! remaining:%d", bn, b.remaining)
				lastBingo = b.remaining * n
				if bingo == 0 {
					bingo = lastBingo
				}
			}
		}
	}

	return
}

type Board struct {
	numbers   [5][5]int
	marked    [5][5]bool
	remaining int
	won       bool
}

func (b *Board) Mark(n int) {
	for i := range b.numbers {
		for j := range b.numbers[i] {
			if b.numbers[i][j] == n {
				b.marked[i][j] = true
				b.remaining -= b.numbers[i][j]
			}
		}
	}
}

func (b *Board) Bingo() bool {
	for i := range b.marked {
		row := true
		column := true
		for j := range b.marked[i] {
			if !b.marked[i][j] {
				row = false
			}
			if !b.marked[j][i] {
				column = false
			}
		}
		if row || column {
			b.won = true
			return true
		}
	}

	return false
}

func (b Board) String() string {
	var out strings.Builder
	bingo := ""
	if b.Bingo() {
		bingo = " BINGO!"
	}
	fmt.Fprintf(&out, "=Total: %d%s=\n", b.remaining, bingo)
	for i := range b.numbers {
		for j := range b.numbers[i] {
			marked := ' '
			if b.marked[i][j] {
				marked = '*'
			}

			fmt.Fprintf(&out, "%c%2d ", marked, b.numbers[i][j])
		}
		fmt.Fprintln(&out)
	}
	return out.String()
}
