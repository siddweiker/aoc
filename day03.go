package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
)

func init() {
	Register(Day03)
}

func Day03(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	nums := Report{}
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

	a1, a2 := nums.CountBits(numBits), nums.CountBitsFilter(numBits)
	return fmt.Sprintf("%d, %d", a1, a2)
}

type Report []uint

func (r Report) CountBits(length int) uint {
	common, least := uint(0), uint(0)
	for i := length - 1; i >= 0; i-- {
		common |= r.CommonBit(i)
	}
	least = (^uint(0) >> (bits.UintSize - length)) ^ common
	return common * least
}

func (r Report) CommonBit(pos int) uint {
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

func (r Report) LeastBit(pos int) uint {
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

func (r Report) Filter(f func(uint) bool) []uint {
	filtered := make([]uint, 0)
	for _, v := range r {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (r Report) CountBitsFilter(length int) uint {
	common, least := uint(0), uint(0)
	com := r
	lea := r

	for i := length - 1; i >= 0; i-- {
		com = com.Filter(func(n uint) bool {
			return n>>(i+1) == common>>(i+1)
		})
		common |= com.CommonBit(i)

		lea = lea.Filter(func(n uint) bool {
			return n>>(i+1) == least>>(i+1)
		})
		least |= lea.LeastBit(i)
	}

	return common * least
}
