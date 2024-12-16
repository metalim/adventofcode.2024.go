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
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	parsed := parseInput(string(bs))
	part1(parsed)
	part2(parsed)
}

type Parsed []string

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return Parsed(lines)
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	for _, line := range parsed {
		fmt.Println(line)
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", 0, time.Since(timeStart))
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	for _, line := range parsed {
		_ = line
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", 0, time.Since(timeStart))
}
