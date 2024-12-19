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

type Parsed struct {
	Patterns []string
	Designs  []string
}

func parseInput(input string) Parsed {
	parts := strings.Split(input, "\n\n")
	patterns := strings.Split(parts[0], ", ")
	designs := strings.Split(parts[1], "\n")
	if len(designs[len(designs)-1]) == 0 {
		designs = designs[:len(designs)-1]
	}
	return Parsed{patterns, designs}
}

var mem = make(map[string]int)

func getPossible(patterns []string, design string) int {
	if v, ok := mem[design]; ok {
		return v
	}
	var count int
	for _, pattern := range patterns {
		nd := strings.TrimPrefix(design, pattern)
		if len(nd) == len(design) {
			continue
		}
		if nd == "" {
			count++
			continue
		}
		count += getPossible(patterns, nd)
	}
	mem[design] = count
	return count
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	var count int
	for _, design := range parsed.Designs {
		if getPossible(parsed.Patterns, design) > 0 {
			count++
		}
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", count, time.Since(timeStart))
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	var count int
	for _, design := range parsed.Designs {
		count += getPossible(parsed.Patterns, design)
	}

	fmt.Printf("Part 2: %d\t\tin %v\n", count, time.Since(timeStart))
}
