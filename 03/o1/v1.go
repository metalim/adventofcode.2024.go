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
