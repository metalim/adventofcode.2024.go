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

	rules, packets := parseInput(string(bs))
	part1(rules, packets)
	part2(rules, packets)
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

var reInts = regexp.MustCompile(`\d+`)

func ints(s string) []int {
	m := reInts.FindAllString(s, -1)
	ns := make([]int, len(m))
	for i, s := range m {
		ns[i] = atoi(s)
	}
	return ns
}

func parseInput(input string) (rules [][]int, packets [][]int) {
	lines := strings.Split(input, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	var next int
	for i, line := range lines {
		if line == "" {
			next = i + 1
			break
		}
		rules = append(rules, ints(line))
	}
	for _, line := range lines[next:] {
		packets = append(packets, ints(line))
	}
	return
}

func getPositions(packet []int) map[int]int {
	pos := map[int]int{}
	for i, n := range packet {
		pos[n] = i
	}
	return pos
}

func isValid(packet []int, rules [][]int) bool {
	pos := getPositions(packet)
	for _, rule := range rules {
		if i, ok := pos[rule[0]]; ok {
			if j, ok := pos[rule[1]]; ok {
				if i > j {
					return false
				}
			}
		}
	}
	return true
}

func part1(rules [][]int, packets [][]int) {
	timeStart := time.Now()
	var sum int
	for _, packet := range packets {
		if isValid(packet, rules) {
			sum += packet[(len(packet)-1)/2]
		}
	}

	fmt.Printf("Part 1: \t\t%d\tin %v\n", sum, time.Since(timeStart))
}

func fix(packet []int, rules [][]int) []int {
	pos := getPositions(packet)
	for _, rule := range rules {
		if i, ok := pos[rule[0]]; ok {
			if j, ok := pos[rule[1]]; ok {
				if i > j {
					packet[i], packet[j] = packet[j], packet[i]
					pos[rule[0]] = j
					pos[rule[1]] = i
					return fix(packet, rules)
				}
			}
		}
	}
	return packet
}

func part2(rules [][]int, packets [][]int) {
	timeStart := time.Now()
	var sum int
	for _, packet := range packets {
		if isValid(packet, rules) {
			continue
		}
		fixed := fix(packet, rules)
		sum += fixed[(len(fixed)-1)/2]
	}

	fmt.Printf("Part 2: \t\t%d\tin %v\n", sum, time.Since(timeStart))
}
