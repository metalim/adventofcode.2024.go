/*
Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
Выведи время решения каждой части.

https://adventofcode.com/2024/day/7

--- Day 7: Bridge Repair ---
The Historians take you to a familiar rope bridge over a river in the middle of a jungle. The Chief isn't on this side of the bridge, though; maybe he's on the other side?

When you go to cross the bridge, you notice a group of engineers trying to repair it. (Apparently, it breaks pretty frequently.) You won't be able to cross until it's fixed.

You ask how long it'll take; the engineers tell you that it only needs final calibrations, but some young elephants were playing nearby and stole all the operators from their calibration equations! They could finish the calibrations if only someone could determine which test values could possibly be produced by placing any combination of operators into their calibration equations (your puzzle input).

For example:

190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20
Each line represents a single equation. The test value appears before the colon on each line; it is your job to determine whether the remaining numbers can be combined with operators to produce the test value.

Operators are always evaluated left-to-right, not according to precedence rules. Furthermore, numbers in the equations cannot be rearranged. Glancing into the jungle, you can see elephants holding two different types of operators: add (+) and multiply (*).

Only three of the above equations can be made true by inserting operators:

190: 10 19 has only one position that accepts an operator: between 10 and 19. Choosing + would give 29, but choosing * would give the test value (10 * 19 = 190).
3267: 81 40 27 has two positions for operators. Of the four possible configurations of the operators, two cause the right side to match the test value: 81 + 40 * 27 and 81 * 40 + 27 both equal 3267 (when evaluated left-to-right)!
292: 11 6 16 20 can be solved in exactly one way: 11 + 6 * 16 + 20.
The engineers just need the total calibration result, which is the sum of the test values from just the equations that could possibly be true. In the above example, the sum of the test values for the three equations listed above is 3749.

Determine which equations could possibly be true. What is their total calibration result?


--- Part Two ---
The engineers seem concerned; the total calibration result you gave them is nowhere close to being within safety tolerances. Just then, you spot your mistake: some well-hidden elephants are holding a third type of operator.

The concatenation operator (||) combines the digits from its left and right inputs into a single number. For example, 12 || 345 would become 12345. All operators are still evaluated left-to-right.

Now, apart from the three equations that could be made true using only addition and multiplication, the above example has three more equations that can be made true by inserting operators:

156: 15 6 can be made true through a single concatenation: 15 || 6 = 156.
7290: 6 8 6 15 can be made true using 6 * 8 || 6 * 15.
192: 17 8 14 can be made true using 17 || 8 + 14.
Adding up all six test values (the three that could be made before using only + and * plus the new three that can now be made by also using ||) produces the new total calibration result of 11387.

Using your new knowledge of elephant hiding spots, determine which equations could possibly be true. What is their total calibration result?

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var opsPart1 = []string{"+", "*"}
var opsPart2 = []string{"+", "*", "||"}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <inputfile>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Part 1
	startPart1 := time.Now()
	sumPart1 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValue(nums, testVal, opsPart1) {
			sumPart1 += testVal
		}
	}
	fmt.Println("Part 1:", sumPart1, "Time:", time.Since(startPart1))

	// Part 2
	startPart2 := time.Now()
	sumPart2 := 0
	for _, line := range lines {
		testVal, nums := parseLine(line)
		if canMakeValue(nums, testVal, opsPart2) {
			sumPart2 += testVal
		}
	}
	fmt.Println("Part 2:", sumPart2, "Time:", time.Since(startPart2))
}

func parseLine(line string) (int, []int) {
	parts := strings.Split(line, ":")
	testVal, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	right := strings.Fields(strings.TrimSpace(parts[1]))
	var nums []int
	for _, r := range right {
		n, _ := strconv.Atoi(r)
		nums = append(nums, n)
	}
	return testVal, nums
}

func canMakeValue(nums []int, target int, ops []string) bool {
	if len(nums) == 1 {
		return nums[0] == target
	}
	return tryOps(nums, target, ops, 0, nums[0])
}

func tryOps(nums []int, target int, ops []string, idx int, current int) bool {
	if idx == len(nums)-1 {
		return current == target
	}
	nextVal := nums[idx+1]
	for _, op := range ops {
		switch op {
		case "+":
			if tryOps(nums, target, ops, idx+1, current+nextVal) {
				return true
			}
		case "*":
			if tryOps(nums, target, ops, idx+1, current*nextVal) {
				return true
			}
		case "||":
			concatVal := concatNumbers(current, nextVal)
			// При конкатенации нужно пересобрать цепочку вычислений,
			// т.к. конкатенация происходит "на месте" текущего и следующего элемента.
			// Поэтому мы создадим новый срез чисел, в котором текущий элемент заменим на конкатенированный,
			// а следующий пропустим.
			newNums := make([]int, 0, len(nums)-1)
			newNums = append(newNums, nums[:idx]...)
			newNums = append(newNums, concatVal)
			newNums = append(newNums, nums[idx+2:]...)
			if canMakeValue(newNums, target, ops) {
				return true
			}
		}
	}
	return false
}

func concatNumbers(a, b int) int {
	as := strconv.Itoa(a)
	bs := strconv.Itoa(b)
	res, _ := strconv.Atoi(as + bs)
	return res
}
