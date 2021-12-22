package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
)

func init() {
	Register(Day3)
}

func Day3(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	nums := report{}
	var numBits int
	for scanner.Scan() {
		line := scanner.Text()
		if numBits == 0 {
			numBits = len(line)
		}
		var num uint
		Sscanf(line, "%b", &num)
		nums = append(nums, num)
	}

	common, least := nums.CountBits(numBits)
	oxygen, co2 := nums.CountBitsFilter(numBits)

	return fmt.Sprintf("%d, %d", common*least, oxygen*co2)
}

type report []uint

func (r report) CountBits(lenght int) (common uint, least uint) {
	for i := lenght - 1; i >= 0; i-- {
		common |= r.CommonBit(i)
	}
	least = (^uint(0) >> (bits.UintSize - lenght)) ^ common
	return
}

func (r report) CommonBit(pos int) uint {
	if len(r) == 1 {
		return r[0] & (1 << pos)
	}

	set := 0
	for _, n := range r {
		// If there is a 1, increment its bit position
		if n&(1<<pos) > 0 {
			set++
		}
	}

	if set > (len(r)/2) || len(r)-set == set {
		return 1 << pos
	}

	return 0
}

func (r report) LeastBit(pos int) uint {
	if len(r) == 1 {
		return r[0] & (1 << pos)
	}

	set := 0
	for _, n := range r {
		// If there is a 1, increment its bit position
		if n&(1<<pos) > 0 {
			set++
		}
	}

	if set < len(r)-set {
		return 1 << pos
	}

	return 0
}

func (r report) Filter(f func(uint) bool) []uint {
	filtered := make([]uint, 0)
	for _, v := range r {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (r report) CountBitsFilter(lenght int) (common uint, least uint) {
	com := r
	lea := r

	for i := lenght - 1; i >= 0; i-- {
		com = com.Filter(func(n uint) bool {
			return n>>(i+1) == common>>(i+1)
		})
		common |= com.CommonBit(i)

		lea = lea.Filter(func(n uint) bool {
			return n>>(i+1) == least>>(i+1)
		})
		least |= lea.LeastBit(i)
	}

	return
}
