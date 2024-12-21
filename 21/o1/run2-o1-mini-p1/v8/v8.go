package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Position представляет координаты на клавиатуре
type Position struct {
	x, y int
}

// Keypad представляет раскладку клавиатуры
type Keypad struct {
	layout map[string]Position
	keys   map[Position]string
}

// NewDirectionalKeypad создает направляющую клавиатуру
func NewDirectionalKeypad() *Keypad {
	layout := map[string]Position{
		"^": {1, 0},
		"A": {2, 0},
		"<": {0, 1},
		"v": {1, 1},
		">": {2, 1},
	}
	keys := make(map[Position]string)
	for k, v := range layout {
		keys[v] = k
	}
	return &Keypad{layout: layout, keys: keys}
}

// NewNumericKeypad создает цифровую клавиатуру
func NewNumericKeypad() *Keypad {
	layout := map[string]Position{
		"7": {0, 0}, "8": {1, 0}, "9": {2, 0},
		"4": {0, 1}, "5": {1, 1}, "6": {2, 1},
		"1": {0, 2}, "2": {1, 2}, "3": {2, 2},
		"0": {1, 3}, "A": {2, 3},
	}
	keys := make(map[Position]string)
	for k, v := range layout {
		keys[v] = k
	}
	return &Keypad{layout: layout, keys: keys}
}

// State представляет текущее состояние в BFS
type State struct {
	yourPos    Position // Позиция вашей клавиатуры
	robot1Pos  Position // Позиция Робота 1
	robot2Pos  Position // Позиция Робота 2
	codeIndex  int      // Текущий индекс в коде
	pressCount int      // Общее количество нажатий кнопок
}

// Queue представляет очередь для BFS
type Queue []State

// Enqueue добавляет состояние в очередь
func (q *Queue) Enqueue(s State) {
	*q = append(*q, s)
}

// Dequeue удаляет и возвращает первое состояние из очереди
func (q *Queue) Dequeue() State {
	s := (*q)[0]
	*q = (*q)[1:]
	return s
}

// isEmpty проверяет, пуста ли очередь
func (q *Queue) isEmpty() bool {
	return len(*q) == 0
}

func main() {
	startTotal := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Пожалуйста, укажите входной файл в качестве аргумента.")
		return
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Инициализация клавиатур
	yourKeypad := NewDirectionalKeypad()
	robot1Keypad := NewDirectionalKeypad()
	robot2Keypad := NewDirectionalKeypad()
	numericKeypad := NewNumericKeypad()

	// Чтение кодов из входного файла
	scanner := bufio.NewScanner(file)
	var codes []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	elapsedTotal := time.Since(startTotal)
	fmt.Printf("Входные данные обработаны за %v\n", elapsedTotal)

	totalComplexity := 0

	for _, code := range codes {
		start := time.Now()
		sequenceLength := bfs(yourKeypad, robot1Keypad, robot2Keypad, numericKeypad, code)
		if sequenceLength == -1 {
			fmt.Printf("Код: %s невозможно ввести на клавиатуре.\n", code)
			continue
		}
		numericValue := getNumericValue(code)
		complexity := sequenceLength * numericValue
		totalComplexity += complexity
		elapsed := time.Since(start)
		fmt.Printf("Код: %s, Длина последовательности: %d, Числовое значение: %d, Сложность: %d, Время: %v\n",
			code, sequenceLength, numericValue, complexity, elapsed)
	}

	fmt.Printf("Общая сложность: %d\n", totalComplexity)
}

// bfs выполняет поиск в ширину для нахождения минимальной длины последовательности нажатий
func bfs(yourKeypad, robot1Keypad, robot2Keypad, numericKeypad *Keypad, code string) int {
	target := strings.ToUpper(code)
	codeRunes := []rune(target)
	codeLength := len(codeRunes)

	// Инициализация начального состояния
	initialState := State{
		yourPos:    yourKeypad.layout["A"],
		robot1Pos:  robot1Keypad.layout["A"],
		robot2Pos:  robot2Keypad.layout["A"],
		codeIndex:  0,
		pressCount: 0,
	}

	queue := &Queue{}
	queue.Enqueue(initialState)

	visited := make(map[string]bool)
	visited[stateKey(initialState)] = true

	for !queue.isEmpty() {
		current := queue.Dequeue()

		// Если все символы введены
		if current.codeIndex == codeLength {
			return current.pressCount
		}

		// Попробовать все возможные направления: '^', 'v', '<', '>'
		for direction := range yourKeypad.layout {
			if direction == "A" {
				continue
			}

			// Попытаться нажать направление
			newYourPos, valid := moveDirection(current.yourPos, direction, yourKeypad)
			if !valid {
				continue
			}
			newRobot1Pos, valid1 := moveDirection(current.robot1Pos, direction, robot1Keypad)
			if !valid1 {
				continue
			}
			newRobot2Pos, valid2 := moveDirection(current.robot2Pos, direction, robot2Keypad)
			if !valid2 {
				continue
			}

			newState := State{
				yourPos:    newYourPos,
				robot1Pos:  newRobot1Pos,
				robot2Pos:  newRobot2Pos,
				codeIndex:  current.codeIndex,
				pressCount: current.pressCount + 3, // Нажатие влияет на все уровни
			}

			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue.Enqueue(newState)
			}
		}

		// Попытаться нажать 'A'
		currentDesiredChar := string(codeRunes[current.codeIndex])

		// Проверить, указывает ли Робот 2 на нужный символ на цифровой клавиатуре
		currentRobot2Button := getKeyAt(current.robot2Pos, numericKeypad)
		if currentRobot2Button == currentDesiredChar {
			newState := State{
				yourPos:    current.yourPos,
				robot1Pos:  current.robot1Pos,
				robot2Pos:  current.robot2Pos,
				codeIndex:  current.codeIndex + 1,
				pressCount: current.pressCount + 3, // Нажатие 'A' влияет на все уровни
			}
			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				queue.Enqueue(newState)
			}
		}
	}

	// Если код невозможно ввести
	return -1
}

// moveDirection пытается переместить руку в заданном направлении на указанной клавиатуре
func moveDirection(current Position, direction string, keypad *Keypad) (Position, bool) {
	newX, newY := current.x, current.y
	switch direction {
	case "^":
		newY--
	case "v":
		newY++
	case "<":
		newX--
	case ">":
		newX++
	default:
		return current, false
	}

	// Проверка, существует ли новая позиция на клавиатуре
	if _, exists := keypad.keys[Position{newX, newY}]; exists {
		return Position{newX, newY}, true
	}
	return current, false
}

// stateKey генерирует уникальный ключ для состояния
func stateKey(state State) string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", state.yourPos.x, state.yourPos.y,
		state.robot1Pos.x, state.robot1Pos.y,
		state.robot2Pos.x, state.robot2Pos.y, state.codeIndex)
}

// getKeyAt возвращает символ на заданной позиции клавиатуры
func getKeyAt(pos Position, keypad *Keypad) string {
	if key, exists := keypad.keys[pos]; exists {
		return key
	}
	return ""
}

// getNumericValue извлекает числовое значение из кода, игнорируя ведущие нули
func getNumericValue(code string) int {
	num := 0
	for _, char := range code {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		}
	}
	return num
}
