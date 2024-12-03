package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func ints(s string) []int {
	fields := strings.Fields(s)
	ns := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		catch(err)
		ns[i] = n
	}
	return ns
}

func sign(n int) int {
	if n > 0 {
		return 1
	}
	return -1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	catch(err)
	lines := strings.Split(string(bs), "\n")

	part1(lines)
	part2(lines)
}

func isSafe(ns []int) bool {
	var dir int
	for i, v := range ns[1:] {
		diff := v - ns[i]
		if dir == 0 {
			dir = sign(diff)
		}
		diff *= dir
		if diff < 1 || 3 < diff {
			return false
		}
	}
	return true
}

func part1(lines []string) {
	timeStart := time.Now()

	var safe int
	for _, line := range lines {
		if line == "" {
			continue
		}
		ns := ints(line)
		if isSafe(ns) {
			safe++
		}
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", safe, time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	var safe int
	for _, line := range lines {
		if line == "" {
			continue
		}
		ns := ints(line)
		if isSafe(ns) {
			safe++
			continue
		}
		for i := 0; i < len(ns); i++ {
			fixed := append(append([]int{}, ns[:i]...), ns[i+1:]...)
			if isSafe(fixed) {
				safe++
				break
			}
		}
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", safe, time.Since(timeStart))
}
