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
	"time"
	"unicode"
)

type Runner func(io.Reader) string

var drivers = struct {
	sync.RWMutex
	solutions []Runner
}{}

func main() {
	var (
		test bool
		day  int
	)
	flag.BoolVar(&test, "test", false, "Use test file")
	flag.IntVar(&day, "day", 0, "Specify Day to run")
	flag.Parse()

	filestr := "data/%s.txt"
	if test {
		filestr = "data/%s.test.txt"
	}

	names, days := Solutions()
	var totalRuntime time.Duration
	for i, r := range days {
		if day != 0 && day != i+1 {
			continue
		}
		f, err := os.Open(fmt.Sprintf(filestr, names[i]))
		if err != nil {
			defer f.Close()
		}
		start := time.Now()
		answer := r(f)
		ran := time.Since(start)
		totalRuntime += ran
		log.Printf("%-5s Answer: %-35s [%s]", names[i], answer, ran)
	}
	log.Printf("Total Runtime %s [%s]", strings.Repeat("=", 35), totalRuntime)
}

func Register(r Runner) {
	drivers.Lock()
	defer drivers.Unlock()
	if r == nil {
		panic("Register runner is nil")
	}
	drivers.solutions = append(drivers.solutions, r)
}

func Solutions() ([]string, []Runner) {
	named := map[string]Runner{}
	namesOrder := []string{}
	for _, r := range drivers.solutions {
		fn := getFunctionName(r)
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

func getFunctionName(i interface{}) string {
	return strings.TrimPrefix(
		runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(),
		"main.",
	)
}

// Simple helper functions

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

func ContainsAll(s, chars string) bool {
	for _, c := range chars {
		if !strings.ContainsRune(s, c) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
