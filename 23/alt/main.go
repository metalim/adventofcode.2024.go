package main

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"sort"
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

type Parsed map[string]map[string]bool

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	parsed := make(Parsed)
	for _, line := range lines {
		parts := strings.Split(line, "-")
		if parsed[parts[0]] == nil {
			parsed[parts[0]] = make(map[string]bool)
		}
		if parsed[parts[1]] == nil {
			parsed[parts[1]] = make(map[string]bool)
		}
		parsed[parts[0]][parts[1]] = true
		parsed[parts[1]][parts[0]] = true
	}
	return parsed
}

type Triplet [3]string

func part1(parsed Parsed) {
	timeStart := time.Now()
	// find triplets, where each computes is connected to the other two

	triplets := make(map[Triplet]bool)
	for c1, cs := range parsed {
		for c2 := range cs {
			for c3 := range parsed[c2] {
				if c3 == c1 {
					continue
				}
				if !parsed[c3][c1] {
					continue
				}
				t := Triplet{c1, c2, c3}

				if c1[0] == 't' || c2[0] == 't' || c3[0] == 't' {
					sort.Strings(t[:])
					triplets[t] = true
				}
			}
		}
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", len(triplets), time.Since(timeStart))
}

type Party map[string]bool

var calls int

func dfs(connected Parsed, visited Party) Party {
	calls++
	var largestParty Party
NextC:
	for c := range connected {
		if visited[c] {
			continue
		}
		// check if c is connected to all visited
		for prev := range visited {
			if !connected[c][prev] {
				continue NextC
			}
		}
		// it is connected to all visited, so we can add it to the party

		visited[c] = true
		party := dfs(connected, visited)
		if len(party) > len(largestParty) {
			largestParty = party
		}
		delete(visited, c)
	}
	if len(largestParty) == 0 {
		largestParty = maps.Clone(visited)
	}
	return largestParty
}

// does it need memo? :)
func part2(parsed Parsed) {
	timeStart := time.Now()

	largestParty := dfs(parsed, make(Party))
	fmt.Printf("calls: %d\n", calls)
	sorted := make([]string, 0, len(largestParty))
	for c := range largestParty {
		sorted = append(sorted, c)
	}
	sort.Strings(sorted)
	fmt.Printf("Part 2: %d, %v\t\tin %v\n", len(largestParty), strings.Join(sorted, ","), time.Since(timeStart))
}
