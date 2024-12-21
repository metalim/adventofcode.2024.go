package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Position представляет координаты на клавиатуре
type Position struct {
	x int
	y int
}

// Клавиатура с расположением кнопок
var keypad = map[rune]Position{
	'7': {0, 0}, '8': {1, 0}, '9': {2, 0},
	'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
	'1': {0, 2}, '2': {1, 2}, '3': {2, 2},
	'0': {1, 3}, 'A': {2, 3},
}

// Направления движения
var directions = []struct {
	dx int
	dy int
}{
	{0, -1}, // Вверх (^)
	{0, 1},  // Вниз (v)
	{-1, 0}, // Влево (<)
	{1, 0},  // Вправо (>)
}

// Проверка валидности позиции на клавиатуре
func isValid(pos Position) bool {
	for _, p := range keypad {
		if p == pos {
			return true
		}
	}
	return false
}

// BFS для нахождения минимального количества движений между двумя кнопками
func bfs(start, end rune) int {
	startPos, ok := keypad[start]
	if !ok {
		return -1
	}
	endPos, ok := keypad[end]
	if !ok {
		return -1
	}

	type Node struct {
		pos  Position
		dist int
	}

	visited := make(map[Position]bool)
	queue := []Node{{startPos, 0}}
	visited[startPos] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == endPos {
			return current.dist
		}

		for _, dir := range directions {
			newPos := Position{current.pos.x + dir.dx, current.pos.y + dir.dy}
			if isValid(newPos) && !visited[newPos] {
				visited[newPos] = true
				queue = append(queue, Node{newPos, current.dist + 1})
			}
		}
	}
	return -1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Пожалуйста, укажите имя входного файла в качестве аргумента.")
		return
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var codes []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Функция для расчета сложности
	calculateComplexity := func(codes []string, numLevels int) (int64, time.Duration) {
		start := time.Now()
		var sumComplexity int64 = 0

		for _, code := range codes {
			numericPart := ""
			for _, ch := range code {
				if ch >= '0' && ch <= '9' {
					numericPart += string(ch)
				}
			}
			var numericValue int
			if numericPart != "" {
				numericValue, err = strconv.Atoi(numericPart)
				if err != nil {
					fmt.Printf("Ошибка при парсинге числовой части кода %s: %v\n", code, err)
					os.Exit(1)
				}
			}

			// Вычисляем количество нажатий
			currentKey := 'A'
			totalPresses := 0

			for _, ch := range code {
				if ch != 'A' {
					distance := bfs(currentKey, ch)
					if distance == -1 {
						fmt.Printf("Невозможно добраться от %c до %c\n", currentKey, ch)
						os.Exit(1)
					}
					// Каждое движение требует одного нажатия, плюс одно нажатие 'A'
					totalPresses += distance + 1
				} else {
					// Нажатие 'A' требует только одного нажатия
					totalPresses += 1
				}
				currentKey = ch
			}

			// Учитываем количество уровней
			complexity := int64(totalPresses) * int64(numericValue) * int64(numLevels)
			sumComplexity += complexity
		}

		elapsed := time.Since(start)
		return sumComplexity, elapsed
	}

	// Часть 1
	numLevelsPart1 := 4 // Ваша клавиатура + 2 робота + цифровая клавиатура
	sumComplexityPart1, elapsedPart1 := calculateComplexity(codes, numLevelsPart1)
	fmt.Printf("Часть 1: Сумма сложностей = %d\n", sumComplexityPart1)
	fmt.Printf("Время выполнения части 1: %s\n", elapsedPart1)

	// Часть 2
	numLevelsPart2 := 27 // Ваша клавиатура + 25 роботов + цифровая клавиатура
	sumComplexityPart2, elapsedPart2 := calculateComplexity(codes, numLevelsPart2)
	fmt.Printf("Часть 2: Сумма сложностей = %d\n", sumComplexityPart2)
	fmt.Printf("Время выполнения части 2: %s\n", elapsedPart2)
}
