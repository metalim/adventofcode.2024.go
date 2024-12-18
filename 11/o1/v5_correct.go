package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

// Требование: не использовать файлы для временного хранения.
// Новый подход: использовать мемоизацию результатов для (число, шаги).
// Правила трансформации:
// 1. "0" -> "1"
// 2. Чётная длина: разделяем пополам, обрезаем ведущие нули, считаем сумму для каждой половинки
// 3. Нечётная длина (кроме "0"): умножаем на 2024
//
// Главная проблема в предположениях о ведущих нулях. Сейчас не делаем никаких упрощающих предположений —
// просто выполняем точные операции. Надеемся на разумный входной размер.
//
// Мемоизация будет по ключу: число+"|"+strconv.Itoa(steps).

var memo = make(map[string]int)
var m2024 = big.NewInt(2024)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile>")
	}
	inputFile := os.Args[1]

	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var stones []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		for _, p := range strings.Fields(line) {
			stones = append(stones, p)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	count25 := 0
	for _, s := range stones {
		count25 += countStones(s, 25)
	}
	fmt.Println(count25)
	fmt.Println(time.Since(start))

	start2 := time.Now()
	count75 := 0
	for _, s := range stones {
		count75 += countStones(s, 75)
	}
	fmt.Println(count75)
	fmt.Println(time.Since(start2))
}

func countStones(num string, steps int) int {
	if steps == 0 {
		return 1
	}

	key := num + "|" + fmt.Sprint(steps)
	if v, ok := memo[key]; ok {
		return v
	}

	var res int
	if num == "0" {
		// Правило 1
		res = countStones("1", steps-1)
	} else {
		L := len(num)
		if L%2 == 0 {
			// Чётная длина: правило 2
			half := L / 2
			leftStr := trimLeadingZeros(num[:half])
			rightStr := trimLeadingZeros(num[half:])
			res = countStones(leftStr, steps-1) + countStones(rightStr, steps-1)
		} else {
			// Нечётная длина: правило 3
			// умножаем на 2024
			// Используем big.Int
			bi := new(big.Int)
			bi.SetString(num, 10)
			bi.Mul(bi, m2024)
			res = countStones(bi.String(), steps-1)
		}
	}

	memo[key] = res
	return res
}

func trimLeadingZeros(s string) string {
	s = strings.TrimLeft(s, "0")
	if s == "" {
		s = "0"
	}
	return s
}
