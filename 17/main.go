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

var Custom bool
var From int
var Print bool

func main() {
	flag.BoolVar(&Custom, "custom", false, "run custom part 2")
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
	} else {
		part2(parsed)
	}
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
			reg[0] >>= combo
		case 1: // bxl
			reg[1] ^= literal
		case 2: // bst
			reg[1] = combo % 8
		case 3: // jnz
			if reg[0] != 0 {
				i = literal - 2
			}
		case 4: // bxc, ignore operand
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
			reg[1] = reg[0] >> combo
		case 7: // cdv
			reg[2] = reg[0] >> combo
		}
	}
	return oPos == len(program)
}

func part2(parsed Parsed) {
	timeStart := time.Now()
	if From == 0 {
		fmt.Println(`!!! This will take a "few" days !!!`)
		fmt.Println(`You might want the --custom or --from <val>`)
	}
	for a := From; ; a++ {
		if a%1e7 == 0 {
			fmt.Printf("a: %d\n", a)
		}
		if run2(parsed.program, a) {
			fmt.Printf("Part 2: %d %b\t\tin %v\n", a, a, time.Since(timeStart))
			break
		}
	}
}

// see [input_assembly.txt]
// 010 100 001 001 111 101 000 011 001 100 100 101 101 101 011 000
// 48 bits, start from a=0 and generate from the end
func findA(out []int, a int, fn func(a int) int) (int, bool) {
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

type Fn func(a int) int

func part2_custom(parsed Parsed) {
	start := time.Now()
	outs := [][]int{
		{2, 4, 1, 1, 7, 5, 0, 3, 1, 4, 4, 5, 5, 5, 3, 0},
		{2, 4, 1, 2, 7, 5, 1, 7, 4, 4, 0, 3, 5, 5, 3, 0},
	}

	fns := []Fn{
		func(A int) int {
			return A&7 ^ 5 ^ (A>>(A&7^1))&7
		},
		func(A int) int {
			return A&7 ^ 5 ^ (A>>(A&7^2))&7
		},
	}

	var fn Fn
	for i, o := range outs {
		if slices.Compare(parsed.program, o) == 0 {
			fn = fns[i]
			fmt.Println("found formula for", o)
		}
	}
	if fn == nil {
		fmt.Println("Unsupported input, sorry")
		return
	}

	if a, ok := findA(parsed.program, 0, fn); ok {
		confirmed := run2(parsed.program, a)
		fmt.Printf("Part 2 custom: %d, confirmed: %t\t\tin %v\n", a, confirmed, time.Since(start))
	} else {
		fmt.Println("Part 2 custom: solution not found")
	}
}
