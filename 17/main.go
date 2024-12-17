package main

import (
	"flag"
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
	registers [3]int
	program   []int
}

var reRegister = regexp.MustCompile(`Register (\w): (\d+)`)
var reProgram = regexp.MustCompile(`Program: (.*)`)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	catch(err)
	return i
}

func parseInput(input string) Parsed {
	parts := strings.Split(input, "\n\n")
	mRegisters := reRegister.FindAllStringSubmatch(parts[0], -1)
	if len(mRegisters) == 0 {
		panic("no registers")
	}
	mProgram := reProgram.FindAllStringSubmatch(parts[1], -1)
	if len(mProgram) == 0 {
		panic("no program")
	}
	parsed := Parsed{
		registers: [3]int{},
		program:   make([]int, 0),
	}
	for i, register := range mRegisters {
		parsed.registers[i] = toInt(register[2])
	}
	sProgram := strings.Split(mProgram[0][1], ",")
	for _, s := range sProgram {
		parsed.program = append(parsed.program, toInt(s))
	}
	return parsed
}

func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func run(parsed Parsed) ([3]int, []int) {
	reg := parsed.registers
	var output []int
	for i := 0; i < len(parsed.program); i += 2 {
		opcode := parsed.program[i]
		literal := parsed.program[i+1]

		var combo int
		switch literal {
		case 0, 1, 2, 3:
			combo = literal
		case 4:
			combo = reg[0]
		case 5:
			combo = reg[1]
		case 6:
			combo = reg[2]
		case 7:
			// ignore
		}
		// fmt.Printf("op: %d, lit: %d, combo: %d\n", opcode, literal, combo)
		// fmt.Printf("reg: %v\n", reg)

		switch opcode {
		case 0: // adv
			reg[0] /= pow(2, combo)
		case 1: // bxl
			reg[1] ^= literal
		case 2: // bst
			reg[1] = combo % 8
		case 3: // jnz
			if reg[0] != 0 {
				i = literal - 2
			}
		case 4: // bxc
			reg[1] ^= reg[2]
		case 5: // out
			output = append(output, combo%8)
		case 6: // bdv
			reg[1] = reg[0] / pow(2, combo)
		case 7: // cdv
			reg[2] = reg[0] / pow(2, combo)
		}
	}
	return reg, output
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	_, output := run(parsed)
	var s strings.Builder

	for i, o := range output {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(strconv.Itoa(o))
	}
	fmt.Printf("Part 1: %s\t\tin %v\n", s.String(), time.Since(timeStart))
}

func run2(program []int, a int) bool {
	reg := [3]int{a, 0, 0}
	var oPos int
	for i := 0; i < len(program); i += 2 {
		opcode := program[i]
		literal := program[i+1]

		var combo int
		switch literal {
		case 0, 1, 2, 3:
			combo = literal
		case 4:
			combo = reg[0]
		case 5:
			combo = reg[1]
		case 6:
			combo = reg[2]
		case 7:
			// ignore
		}

		switch opcode {
		case 0: // adv
			reg[0] /= pow(2, combo)
		case 1: // bxl
			reg[1] ^= literal
		case 2: // bst
			reg[1] = combo % 8
		case 3: // jnz
			if reg[0] != 0 {
				i = literal - 2
			}
		case 4: // bxc
			reg[1] ^= reg[2]
		case 5: // out
			out := combo % 8
			if oPos > len(program) {
				return false
			}
			if out != program[oPos] {
				return false
			}
			oPos++
		case 6: // bdv
			reg[1] = reg[0] / pow(2, combo)
		case 7: // cdv
			reg[2] = reg[0] / pow(2, combo)
		}
	}
	return oPos == len(program)
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	for a := 0; ; a++ {
		if a%1e7 == 0 {
			fmt.Printf("a: %d\n", a)
		}
		if run2(parsed.program, a) {
			fmt.Printf("Part 2: %d\t\tin %v\n", a, time.Since(timeStart))
			break
		}
	}
}
