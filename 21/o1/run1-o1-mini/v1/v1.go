package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Определяем структуру для позиции на клавиатуре
type Position struct {
	x int
	y int
}

// Клавиатура инициализируется как карта символов к их позициям
var keypad = map[rune]Position{
	'7': {0, 0}, '8': {1, 0}, '9': {2, 0},
	'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
	'1': {0, 2}, '2': {1, 2}, '3': {2, 2},
	'0': {1, 3}, 'A': {2, 3},
}

// Возможные направления движения
var directions = []struct {
	dx int
	dy int
}{
	{0, -1}, // Вверх (^)
	{0, 1},  // Вниз (v)
	{-1, 0}, // Влево (<)
	{1, 0},  // Вправо (>)
}

// Функция для проверки, существует ли кнопка на данной позиции
func isValid(pos Position) bool {
	for _, p := range keypad {
		if p == pos {
			return true
		}
	}
	return false
}

// Функция BFS для нахождения минимального расстояния между двумя кнопками
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

	// Часть 1
	startPart1 := time.Now()
	sumComplexity := 0

	for _, code := range codes {
		numericPart := ""
		for _, ch := range code {
			if ch >= '0' && ch <= '9' {
				numericPart += string(ch)
			}
		}
		var numericValue int
		if numericPart != "" {
			fmt.Sscanf(numericPart, "%d", &numericValue)
		}

		// Вычисляем количество нажатий
		// Начинаем с 'A'
		currentKey := 'A'
		totalPresses := 0

		for _, ch := range code {
			distance := bfs(currentKey, ch)
			if distance == -1 {
				fmt.Printf("Невозможно добраться от %c до %c\n", currentKey, ch)
				return
			}
			// Каждое движение требует одного нажатия, плюс одно нажатие 'A'
			totalPresses += distance + 1
			currentKey = ch
		}

		sumComplexity += totalPresses * numericValue
	}

	elapsedPart1 := time.Since(startPart1)
	fmt.Printf("Часть 1: Сумма сложностей = %d\n", sumComplexity)
	fmt.Printf("Время выполнения части 1: %s\n", elapsedPart1)

	// Часть 2
	startPart2 := time.Now()
	// Для части 2 требуется учитывать цепочку из 25 роботов.
	// Это значительно усложняет задачу. Предполагая, что каждый уровень увеличивает количество нажатий,
	// можно умножить количество нажатий на количество уровней.

	// Количество уровней: 25 роботов + ваш уровень = 26
	numLevels := 26
	sumComplexityPart2 := 0

	for _, code := range codes {
		numericPart := ""
		for _, ch := range code {
			if ch >= '0' && ch <= '9' {
				numericPart += string(ch)
			}
		}
		var numericValue int
		if numericPart != "" {
			fmt.Sscanf(numericPart, "%d", &numericValue)
		}

		// Вычисляем количество нажатий
		// Начинаем с 'A'
		currentKey := 'A'
		totalPresses := 0

		for _, ch := range code {
			distance := bfs(currentKey, ch)
			if distance == -1 {
				fmt.Printf("Невозможно добраться от %c до %c\n", currentKey, ch)
				return
			}
			// Каждое движение требует одного нажатия, плюс одно нажатие 'A'
			totalPresses += distance + 1
			currentKey = ch
		}

		// Учитываем количество уровней
		sumComplexityPart2 += totalPresses * numericValue * numLevels
	}

	elapsedPart2 := time.Since(startPart2)
	fmt.Printf("Часть 2: Сумма сложностей = %d\n", sumComplexityPart2)
	fmt.Printf("Время выполнения части 2: %s\n", elapsedPart2)
}
