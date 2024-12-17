package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	Custom bool
	Print  bool
	Brute  bool
	From   int
)

func main() {
	flag.BoolVar(&Custom, "custom", false, "run custom part 2")
	flag.BoolVar(&Brute, "brute", false, "run brute part 2")
	flag.IntVar(&From, "from", 0, "continue part 2 from a")
	flag.BoolVar(&Print, "print", false, "print debug info")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run . input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	parsed := parseInput(string(bs))
	part1(parsed)
	if Custom {
		part2_custom(parsed)
	} else if Brute {
		part2_brute(parsed)
	} else {
		part2(parsed)
	}
}

type Parsed struct {
	program   []int
	registers [3]int
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
		program:   make([]int, 0),
		registers: [3]int{},
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

func run(program []int, reg [3]int) ([3]int, []int) {
	var output []int
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
			reg[0] >>= combo
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
			reg[1] = reg[0] >> combo
		case 7: // cdv
			reg[2] = reg[0] >> combo
		}
	}
	return reg, output
}

func part1(parsed Parsed) {
	timeStart := time.Now()
	_, output := run(parsed.program, parsed.registers)
	var s strings.Builder

	for i, o := range output {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(strconv.Itoa(o))
	}
	fmt.Printf("Part 1: %s\t\tin %v\n", s.String(), time.Since(timeStart))
}

var jnz0 = []int{3, 0}

func part2(parsed Parsed) {
	timeStart := time.Now()
	i := len(parsed.program) - len(jnz0)
	if slices.Compare(parsed.program[i:], jnz0) != 0 {
		fmt.Println("Part 2: input is not supported")
		return
	}
	cycle := parsed.program[:i]
	fn := func(a int) int {
		_, out := run(cycle, [3]int{a, 0, 0})
		return out[0]
	}
	if a, ok := findA(parsed.program, 0, fn); ok {
		confirmed := run2(parsed.program, a)
		fmt.Printf("Part 2: %d, confirmed: %t\t\tin %v\n", a, confirmed, time.Since(timeStart))
	} else {
		fmt.Println("Part 2: solution not found")
	}
}

type Fn func(a int) int

func findA(out []int, a int, fn Fn) (int, bool) {
	if len(out) == 0 {
		return a, true
	}
	i := len(out) - 1
	o := out[i]
	a <<= 3
	for bits := 0; bits < 8; bits++ {
		na := a | bits
		v := fn(na)
		if v == o {
			if Print {
				fmt.Printf("i: %d, o: %d, a: %b\n", i, o, na)
			}
			if found, ok := findA(out[:i], na, fn); ok {
				return found, true
			}
		}
	}
	return 0, false
}
