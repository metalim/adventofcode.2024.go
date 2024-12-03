package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
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

func parseInput(lines []string) ([]int, []int) {
	var list1, list2 []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		list1 = append(list1, atoi(fields[0]))
		list2 = append(list2, atoi(fields[1]))
	}
	slices.Sort(list1)
	slices.Sort(list2)
	return list1, list2
}

func part1(lines []string) {
	timeStart := time.Now()
	list1, list2 := parseInput(lines)

	var sum int
	for i, v := range list1 {
		sum += abs(v - list2[i])
	}
	fmt.Printf("Part 1: %d\t\tin %v\n", sum, time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	list1, list2 := parseInput(lines)
	freqs := make(map[int]int)
	for _, v := range list2 {
		freqs[v]++
	}

	var sum int
	for _, v := range list1 {
		sum += v * freqs[v]
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", sum, time.Since(timeStart))
}
