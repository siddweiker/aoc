package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
)

func init() {
	Register(Day16)
}

func Day16(r io.Reader) string {
	scanner := bufio.NewScanner(r)

	var byteStr strings.Builder
	if scanner.Scan() {
		line := scanner.Text()

		dec, err := hex.DecodeString(line)
		if err != nil {
			log.Fatal(err)
		}
		for _, b := range dec {
			fmt.Fprintf(&byteStr, "%08b", b)
		}
	}

	a1, a2, _ := SubPackets(byteStr.String())
	return fmt.Sprintf("%d, %d", a1, a2)
}

func SubPackets(str string) (int, int, string) {
	if str == "" || strings.Count(str, "0") == len(str) {
		return 0, 0, ""
	}

	version := Binary(str[0:3])
	typeID := Binary(str[3:6])
	str = str[6:]

	if typeID == 4 {
		v := uint(0)
		for {
			done := str[0]
			var n uint
			Sscanf(str[1:5], "%b", &n)
			v |= n // Set the group bits in v
			str = str[5:]

			if done == '0' {
				break
			}

			v <<= 4 // Move the bits over for the next group
		}

		return version, int(v), str
	}

	operator := str[0]
	str = str[1:]
	if operator == '0' {
		subBits := Binary(str[:15])
		str = str[15:]

		ver, val, subStr := SubPackets(str[:subBits])
		str = str[subBits:]
		// Read sub length until all bits are exhausted
		for subStr != "" {
			vr, vl, s := SubPackets(subStr)
			ver += vr
			val = Combine(typeID, val, vl)
			subStr = s
		}

		return version + ver, val, str
	}

	subPacks := Binary(str[:11])
	str = str[11:]

	ver, val, s := SubPackets(str)
	str = s
	// Read sub number of packets
	for i := 1; i < subPacks; i++ {
		vr, vl, s := SubPackets(str)
		ver += vr
		val = Combine(typeID, val, vl)
		str = s
	}

	return version + ver, val, str
}

func Binary(s string) int {
	var n int
	Sscanf(s, "%b", &n)
	return n
}

func Combine(typeID, a, b int) (result int) {
	switch typeID {
	case 0:
		result = a + b
	case 1:
		result = a * b
	case 2:
		result = Min(a, b)
	case 3:
		result = Max(a, b)
	case 5:
		if a > b {
			result = 1
		}
	case 6:
		if a < b {
			result = 1
		}
	case 7:
		if a == b {
			result = 1
		}
	}
	return
}
