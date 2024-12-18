package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile>")
	}

	inputFile := os.Args[1]
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var stones []*big.Int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		for _, p := range parts {
			n := new(big.Int)
			n.SetString(p, 10)
			stones = append(stones, n)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	stones = evolve(stones, 25)
	fmt.Println(len(stones))
	fmt.Println(time.Since(start))

	startPart2 := time.Now()
	stones = evolve(stones, 50)
	fmt.Println(len(stones))
	fmt.Println(time.Since(startPart2))
}

func evolve(stones []*big.Int, times int) []*big.Int {
	for i := 0; i < times; i++ {
		stones = blink(stones)
	}
	return stones
}

func blink(stones []*big.Int) []*big.Int {
	var newStones []*big.Int
	m2024 := big.NewInt(2024)
	for _, s := range stones {
		if s.Sign() == 0 {
			// Rule 1
			newStones = append(newStones, big.NewInt(1))
		} else {
			str := s.String()
			if len(str)%2 == 0 {
				// Rule 2
				half := len(str) / 2
				leftStr := strings.TrimLeft(str[:half], "0")
				rightStr := strings.TrimLeft(str[half:], "0")
				if leftStr == "" {
					leftStr = "0"
				}
				if rightStr == "" {
					rightStr = "0"
				}
				left := new(big.Int)
				left.SetString(leftStr, 10)
				right := new(big.Int)
				right.SetString(rightStr, 10)
				newStones = append(newStones, left, right)
			} else {
				// Rule 3
				ns := new(big.Int)
				ns.Mul(s, m2024)
				newStones = append(newStones, ns)
			}
		}
	}
	return newStones
}
