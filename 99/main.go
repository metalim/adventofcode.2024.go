package main

import (
	"fmt"
	"os"
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

func parseInput(input string) []string {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func part1(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Printf("Part 1: \t\tin %v\n", time.Since(timeStart))
}

func part2(lines []string) {
	timeStart := time.Now()
	for _, line := range lines {
		_ = line
	}

	fmt.Printf("Part 2: \t\tin %v\n", time.Since(timeStart))
}
