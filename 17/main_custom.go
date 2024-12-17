package main

import (
	"fmt"
	"slices"
	"time"
)

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
