package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
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
		fmt.Println("Usage: go run main.go input.txt")
		os.Exit(1)
	}

	bs, err := os.ReadFile(flag.Arg(0))
	catch(err)

	input := parseInput(string(bs))
	part1(input)
	part2(input)
}

type Input string

func parseInput(input string) Input {
	lines := strings.Split(input, "\n")
	return Input(lines[0])
}

const FREE = -1

func buildDisk(input Input) (disk []int) {
	var id int
	var free bool
	for _, r := range input {
		val := id
		if free {
			val = FREE
		} else {
			id++
		}
		disk = append(disk, slices.Repeat([]int{val}, int(r-'0'))...)
		free = !free
	}
	return
}

func diskChecksum(disk []int) int {
	var checksum int
	for i, v := range disk {
		if v != FREE {
			checksum += i * v
		}
	}
	return checksum
}

func part1(input Input) {
	timeStart := time.Now()
	disk := buildDisk(input)
	j := len(disk) - 1
	for i, v := range disk {
		if v != FREE {
			continue
		}
		for disk[j] == FREE {
			j--
		}
		if i >= j {
			break
		}
		disk[i] = disk[j]
		disk[j] = FREE
		j--
	}
	checksum := diskChecksum(disk)
	fmt.Printf("Part 1: %d\t\tin %v\n", checksum, time.Since(timeStart))
}

const NOT_FOUND = -1

func part2(input Input) {
	timeStart := time.Now()
	disk := buildDisk(input)
	for j := len(disk) - 1; j >= 0; j-- {
		if disk[j] == FREE {
			continue
		}

		// get length of file
		fileLength := 1
		id := disk[j]
		for j > 0 && disk[j-1] == id {
			fileLength++
			j--
		}

		// find free space to fit file
		fitPos := NOT_FOUND
		for i := 0; i < j; i++ {
			if disk[i] != FREE {
				continue
			}
			freeLength := 0
			for i < j && disk[i] == FREE {
				freeLength++
				i++
			}
			if freeLength >= fileLength {
				fitPos = i - freeLength
				break
			}
		}
		if fitPos == NOT_FOUND {
			continue
		}

		// move file to fitPos
		copy(disk[fitPos:], disk[j:j+fileLength])
		for k := j; k < j+fileLength; k++ {
			disk[k] = FREE
		}
	}
	fmt.Printf("Part 2: %d\t\tin %v\n", diskChecksum(disk), time.Since(timeStart))
}
