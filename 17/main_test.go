package main

import (
	"testing"

	// "github.com/zeebo/assert"
	"github.com/stretchr/testify/assert"
)

/*
If register C contains 9, the program 2,6 would set register B to 1.
If register A contains 10, the program 5,0,5,1,5,4 would output 0,1,2.
If register A contains 2024, the program 0,1,5,4,3,0 would output 4,2,5,6,7,7,7,7,3,1,0 and leave 0 in register A.
If register B contains 29, the program 1,7 would set register B to 26.
If register B contains 2024 and register C contains 43690, the program 4,0 would set register B to 44354.
*/
func TestCPU(t *testing.T) {
	reg, output := run([]int{2, 6}, [3]int{0, 0, 9})
	assert.Equal(t, 1, reg[1])

	reg, output = run([]int{5, 0, 5, 1, 5, 4}, [3]int{10, 0, 0})
	assert.Equal(t, []int{0, 1, 2}, output)

	reg, output = run([]int{0, 1, 5, 4, 3, 0}, [3]int{2024, 0, 0})
	assert.Equal(t, []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}, output)
	assert.Equal(t, 0, reg[0])

	reg, output = run([]int{1, 7}, [3]int{0, 29, 0})
	assert.Equal(t, 26, reg[1])

	reg, _ = run([]int{4, 0}, [3]int{0, 2024, 43690})
	assert.Equal(t, 44354, reg[1])
}
