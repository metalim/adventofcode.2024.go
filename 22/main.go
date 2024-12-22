package main

import (
	"flag"
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

type Parsed []int

func parseInput(input string) Parsed {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	ints := make([]int, len(lines))
	for i, line := range lines {
		ints[i] = atoi(line)
	}
	return ints
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

const Mod = 16777216
const Repeat = 2000

func hash(n int) int {
	n = (n ^ (n * 64)) % Mod
	n = (n ^ (n / 32)) % Mod
	n = (n ^ (n * 2048)) % Mod
	return n
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	var sum int
	for _, n := range parsed {
		for i := 0; i < Repeat; i++ {
			n = hash(n)
		}
		sum += n
	}

	fmt.Printf("Part 1: %d\t\tin %v\n", sum, time.Since(timeStart))
}

type Seq [4]int

func part2(parsed Parsed) {
	timeStart := time.Now()
	seqsProfit := make(map[Seq]int)
	saw := make(map[Seq]bool)
	for _, n := range parsed {
		var seq [4]int
		clear(saw)

		for i := 0; i < Repeat; i++ {
			prevPrice := n % 10
			n = hash(n)
			price := n % 10
			delta := price - prevPrice
			seq[0], seq[1], seq[2], seq[3] = seq[1], seq[2], seq[3], delta

			// 0, 1, 2, 3
			// at 3 we have the sequence of 4 price changes
			if i < 3 {
				// we have no sequence yet
				continue
			}
			// we need only the first seq for this n
			if !saw[seq] {
				seqsProfit[seq] += price
				saw[seq] = true
			}
		}
	}
	fmt.Printf("Seqs: %d\n", len(seqsProfit))

	var maxBananas int
	var maxSeq Seq
	for seq, price := range seqsProfit {
		if price > maxBananas {
			maxBananas = price
			maxSeq = seq
		}
	}

	fmt.Printf("Part 2: %v: %d\t\tin %v\n", maxSeq, maxBananas, time.Since(timeStart))
}
