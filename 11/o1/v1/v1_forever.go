/*
Напиши код на Go для решения следующей задачи.
Входные данные в файле указываемом аргументом в командной строке.
Выведи ответ и время решения после решения каждой части.

--- Day 11: Plutonian Pebbles ---
The ancient civilization on Pluto was known for its ability to manipulate spacetime, and while The Historians explore their infinite corridors, you've noticed a strange set of physics-defying stones.

At first glance, they seem like normal stones: they're arranged in a perfectly straight line, and each stone has a number engraved on it.

The strange part is that every time you blink, the stones change.

Sometimes, the number engraved on a stone changes. Other times, a stone might split in two, causing all the other stones to shift over a bit to make room in their perfectly straight line.

As you observe them for a while, you find that the stones have a consistent behavior. Every time you blink, the stones each simultaneously change according to the first applicable rule in this list:

If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.
No matter how the stones change, their order is preserved, and they stay on their perfectly straight line.

How will the stones evolve if you keep blinking at them? You take a note of the number engraved on each stone in the line (your puzzle input).

If you have an arrangement of five stones engraved with the numbers 0 1 10 99 999 and you blink once, the stones transform as follows:

The first stone, 0, becomes a stone marked 1.
The second stone, 1, is multiplied by 2024 to become 2024.
The third stone, 10, is split into a stone marked 1 followed by a stone marked 0.
The fourth stone, 99, is split into two stones marked 9.
The fifth stone, 999, is replaced by a stone marked 2021976.
So, after blinking once, your five stones would become an arrangement of seven stones engraved with the numbers 1 2024 1 0 9 9 2021976.

Here is a longer example:

Initial arrangement:
125 17

After 1 blink:
253000 1 7

After 2 blinks:
253 0 2024 14168

After 3 blinks:
512072 1 20 24 28676032

After 4 blinks:
512 72 2024 2 0 2 4 2867 6032

After 5 blinks:
1036288 7 2 20 24 4048 1 4048 8096 28 67 60 32

After 6 blinks:
2097446912 14168 4048 2 0 2 4 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
In this example, after blinking six times, you would have 22 stones. After blinking 25 times, you would have 55312 stones!

Consider the arrangement of stones in front of you. How many stones will you have after blinking 25 times?

--- Part Two ---
The Historians sure are taking a long time. To be fair, the infinite corridors are very large.

How many stones would you have after blinking a total of 75 times?


*/

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
