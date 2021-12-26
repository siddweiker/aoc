package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Runner func(io.Reader) string

var (
	test      = false
	driversMu sync.RWMutex
	drivers   = []Runner{}
)

func main() {
	flag.BoolVar(&test, "test", false, "Use test file")
	flag.Parse()
	filestr := "data/%s.txt"
	if test {
		filestr = "data/%s.test.txt"
	}

	names, drivers := Drivers()
	for i, r := range drivers {
		f, err := os.Open(fmt.Sprintf(filestr, names[i]))
		if err != nil {
			defer f.Close()
		}
		log.Printf("#%d:%s Answer: %s", i, names[i], r(f))
	}
}

func Register(r Runner) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if r == nil {
		panic("Register runner is nil")
	}
	drivers = append(drivers, r)
}

func Drivers() ([]string, []Runner) {
	named := map[string]Runner{}
	namesOrder := []string{}
	for _, r := range drivers {
		fn := GetFunctionName(r)
		namesOrder = append(namesOrder, fn)
		named[fn] = r
	}

	sort.Slice(namesOrder, func(i, j int) bool {
		if len(namesOrder[i]) == len(namesOrder[j]) {
			return namesOrder[i] < namesOrder[j]
		}
		return len(namesOrder[i]) < len(namesOrder[j])
	})
	runners := make([]Runner, len(namesOrder))
	for i, fn := range namesOrder {
		runners[i] = named[fn]
	}

	return namesOrder, runners
}

func GetFunctionName(i interface{}) string {
	return strings.TrimPrefix(
		runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(),
		"main.",
	)
}

func Sscanf(str, format string, a ...interface{}) {
	_, err := fmt.Sscanf(str, format, a...)
	if err != nil {
		log.Printf("error parsing line '%s': %v", str, err)
	}
}

func Atoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("error parsing number '%s': %v", str, err)
	}
	return i
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Coming in 1.18
func Cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func ContainsAll(s, chars string) bool {
	for _, c := range chars {
		if !strings.ContainsRune(s, c) {
			return false
		}
	}
	return true
}
