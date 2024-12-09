package main

import (
	"flag"
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
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

type Input []string

func parseInput(input string) Input {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return Input(lines)
}

func part1(input Input) {
	timeStart := time.Now()
	for _, line := range input {
		fmt.Println(line)
	}

	fmt.Printf("Part 1: \t\tin %v\n", time.Since(timeStart))
}

func part2(input Input) {
	timeStart := time.Now()
	for _, line := range input {
		_ = line
	}

	fmt.Printf("Part 2: \t\tin %v\n", time.Since(timeStart))
}
