/*
Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
Выведи время решения каждой части.

--- Day 3: Mull It Over ---
"Our computers are having issues, so I have no idea if we have any Chief Historians in stock! You're welcome to check the warehouse, though," says the mildly flustered shopkeeper at the North Pole Toboggan Rental Shop. The Historians head out to take a look.

The shopkeeper turns to you. "Any chance you can see why our computers are having issues again?"

The computer appears to be trying to run a program, but its memory (your puzzle input) is corrupted. All of the instructions have been jumbled up!

It seems like the goal of the program is just to multiply some numbers. It does that with instructions like mul(X,Y), where X and Y are each 1-3 digit numbers. For instance, mul(44,46) multiplies 44 by 46 to get a result of 2024. Similarly, mul(123,4) would multiply 123 by 4.

However, because the program's memory has been corrupted, there are also many invalid characters that should be ignored, even if they look like part of a mul instruction. Sequences like mul(4*, mul(6,9!, ?(12,34), or mul ( 2 , 4 ) do nothing.

For example, consider the following section of corrupted memory:

xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))
Only the four highlighted sections are real mul instructions. Adding up the result of each instruction produces 161 (2*4 + 5*5 + 11*8 + 8*5).

Scan the corrupted memory for uncorrupted mul instructions. What do you get if you add up all of the results of the multiplications?


--- Part Two ---
As you scan through the corrupted memory, you notice that some of the conditional statements are also still intact. If you handle some of the uncorrupted conditional statements in the program, you might be able to get an even more accurate result.

There are two new instructions you'll need to handle:

The do() instruction enables future mul instructions.
The don't() instruction disables future mul instructions.
Only the most recent do() or don't() instruction applies. At the beginning of the program, mul instructions are enabled.

For example:

xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))
This corrupted memory is similar to the example from before, but this time the mul(5,5) and mul(11,8) instructions are disabled because there is a don't() instruction before them. The other mul instructions function normally, including the one at the end that gets re-enabled by a do() instruction.

This time, the sum of the results is 48 (2*4 + 8*5).

Handle the new instructions; what do you get if you add up all of the results of just the enabled multiplications?

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <input_file>")
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := string(data)

	// Часть 1
	start1 := time.Now()
	reMul := regexp.MustCompile(`mul\(\s*(\d+)\s*,\s*(\d+)\s*\)`)
	sumPart1 := 0
	for _, match := range reMul.FindAllStringSubmatch(input, -1) {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		sumPart1 += x * y
	}
	duration1 := time.Since(start1)

	// Часть 2
	start2 := time.Now()
	reDo := regexp.MustCompile(`do\(\)`)
	reDont := regexp.MustCompile(`don't\(\)`)
	enabled := true
	sumPart2 := 0

	// Разобьем все на токены
	tokens := reMul.FindAllStringIndex(input, -1)
	doTokens := reDo.FindAllStringIndex(input, -1)
	dontTokens := reDont.FindAllStringIndex(input, -1)

	// Соберем все интервалы
	type item struct {
		pos int
		typ int // 0 - mul, 1 - do, 2 - don't
		sub []string
	}
	var items []item
	for _, t := range tokens {
		m := reMul.FindStringSubmatch(input[t[0]:t[1]])
		items = append(items, item{t[0], 0, m})
	}
	for _, t := range doTokens {
		items = append(items, item{t[0], 1, nil})
	}
	for _, t := range dontTokens {
		items = append(items, item{t[0], 2, nil})
	}

	// Отсортируем по позиции
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].pos < items[i].pos {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	for _, it := range items {
		switch it.typ {
		case 1:
			enabled = true
		case 2:
			enabled = false
		case 0:
			if enabled {
				x, _ := strconv.Atoi(it.sub[1])
				y, _ := strconv.Atoi(it.sub[2])
				sumPart2 += x * y
			}
		}
	}

	duration2 := time.Since(start2)

	fmt.Println("Part 1:", sumPart1, "Time:", duration1)
	fmt.Println("Part 2:", sumPart2, "Time:", duration2)
}
