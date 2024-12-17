package main

import (
	"fmt"
	"time"
)

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

func part2_brute(parsed Parsed) {
	timeStart := time.Now()
	if From == 0 {
		fmt.Println(`!!! This will take a "few" days !!!`)
		fmt.Println(`You might want the --from <val>`)
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
