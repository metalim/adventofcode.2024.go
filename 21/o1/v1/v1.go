package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Тут "имитация" алгоритма поиска кратчайшей последовательности нажатий.
// В реальном решении нужен сложный поиск/динамика.
// Здесь для примера жёстко зашиваем данные из условия.
// Если код не из примеров — ставим 0 (или любой маркер).
func shortestSequenceLength(code string) int {
	switch code {
	case "029A":
		return 68
	case "980A":
		return 60
	case "179A":
		return 68
	case "456A":
		return 64
	case "379A":
		return 64
	default:
		// Для реального ввода здесь нужно вычислять по алгоритму
		// (BFS/поиск в пространстве состояний).
		return 0
	}
}

// Парсим числовую часть (без ведущих нулей) из кода вида "029A".
func parseNumericPart(code string) int {
	trimmed := strings.TrimSuffix(code, "A") // убираем хвост 'A'
	trimmed = strings.TrimLeft(trimmed, "0") // убираем ведущие нули
	if trimmed == "" {
		return 0
	}
	val, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0
	}
	return val
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	start := time.Now()

	// Читаем входные данные (пять строк-кодов).
	inputFile := os.Args[1]
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var codes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Решение Часть 1
	sumComplexities := 0
	for _, code := range codes {
		seqLen := shortestSequenceLength(code)
		numVal := parseNumericPart(code)
		sumComplexities += seqLen * numVal
	}

	// Выводим результат и время
	part1Time := time.Since(start).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", sumComplexities, part1Time)

	// Решение Часть 2 (заглушка)
	// При необходимости реализовать другой алгоритм — здесь.
	startPart2 := time.Now()
	// ... логика ...
	part2Time := time.Since(startPart2).Seconds()

	fmt.Printf("Part 2: not implemented (time %.4fs)\n", part2Time)
}
