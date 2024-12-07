package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(os.Args[1])
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

var reInts = regexp.MustCompile(`\d+`)

func ints(s string) []int {
	ss := reInts.FindAllString(s, -1)
	is := make([]int, len(ss))
	for i, s := range ss {
		n, err := strconv.Atoi(s)
		catch(err)
		is[i] = n
	}
	return is
}

func parseInput(input string) [][]int {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	linesInts := make([][]int, len(lines))
	for i, line := range lines {
		linesInts[i] = ints(line)
	}
	return linesInts
}

func isValid(result int, ns []int) bool {
	if len(ns) == 1 {
		return result == ns[0]
	}
	if isValid(result-ns[len(ns)-1], ns[:len(ns)-1]) {
		return true
	}
	if result%ns[len(ns)-1] == 0 && isValid(result/ns[len(ns)-1], ns[:len(ns)-1]) {
		return true
	}
	return false
}

func part1(lines [][]int) {
	timeStart := time.Now()
	var sum int
	for _, line := range lines {
		if isValid(line[0], line[1:]) {
			sum += line[0]
		}
	}

	fmt.Printf("Part 1: \t\t%d\tin %v\n", sum, time.Since(timeStart))
}

func concat(ns []int) int {
	var s strings.Builder
	for _, n := range ns {
		s.WriteString(strconv.Itoa(n))
	}
	n, err := strconv.Atoi(s.String())
	catch(err)
	return n
}

// concat|| should also be evaluated left to right !!! No priority
// so we can't calculate the result of last operation first
// and have to go left to right instead
func isValid3(result int, ns []int) bool {
	if len(ns) == 1 {
		return result == ns[0]
	}
	if isValid3(result, append([]int{ns[0] + ns[1]}, ns[2:]...)) {
		return true
	}
	if isValid3(result, append([]int{ns[0] * ns[1]}, ns[2:]...)) {
		return true
	}
	if isValid3(result, append([]int{concat(ns[:2])}, ns[2:]...)) {
		return true
	}

	return false
}

func part2(lines [][]int) {
	timeStart := time.Now()
	var sum int
	for _, line := range lines {
		if isValid3(line[0], line[1:]) {
			sum += line[0]
		}
	}

	fmt.Printf("Part 2: \t\t%d\tin %v\n", sum, time.Since(timeStart))
}
