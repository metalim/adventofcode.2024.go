/*
➜ go run ./o1/v1 sample.txt
# github.com/metalim/adventofcode.2024.go/17/o1/v1
o1/v1/v1.go:243:2: declared and not used: partOneStart
o1/v1/v1.go:245:2: declared and not used: partOneEnd
o1/v1/v1.go:253:2: declared and not used: partTwoStart
o1/v1/v1.go:255:2: declared and not used: partTwoEnd
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Registers структура для хранения регистров A, B, C
type Registers struct {
	A int
	B int
	C int
}

// Computer структура для хранения программы и регистров
type Computer struct {
	program []int
	reg     Registers
	ip      int
	output  []int
}

// NewComputer инициализирует новый Computer с заданной программой и регистрами
func NewComputer(program []int, reg Registers) *Computer {
	return &Computer{
		program: program,
		reg:     reg,
		ip:      0,
		output:  []int{},
	}
}

// getOperandValue возвращает значение операнда на основе его типа
func (comp *Computer) getOperandValue(operand int) int {
	switch {
	case operand >= 0 && operand <= 3:
		return operand
	case operand == 4:
		return comp.reg.A
	case operand == 5:
		return comp.reg.B
	case operand == 6:
		return comp.reg.C
	default:
		// Операнд 7 зарезервирован и не должен появляться
		return 0
	}
}

// Execute выполняет программу до остановки
func (comp *Computer) Execute() {
	for comp.ip < len(comp.program) {
		if comp.ip+1 >= len(comp.program) {
			break // Остановиться, если opcode или operand выходят за границы
		}
		opcode := comp.program[comp.ip]
		operand := comp.program[comp.ip+1]

		switch opcode {
		case 0: // adv
			denominator := 1 << comp.getOperandValue(operand) // 2^operand
			if denominator == 0 {
				comp.reg.A = 0
			} else {
				comp.reg.A = comp.reg.A / denominator
			}
			comp.ip += 2
		case 1: // bxl
			comp.reg.B = comp.reg.B ^ operand
			comp.ip += 2
		case 2: // bst
			comp.reg.B = operand % 8
			comp.ip += 2
		case 3: // jnz
			if comp.reg.A != 0 {
				comp.ip = operand
			} else {
				comp.ip += 2
			}
		case 4: // bxc
			comp.reg.B = comp.reg.B ^ comp.reg.C
			comp.ip += 2
		case 5: // out
			value := comp.getOperandValue(operand) % 8
			comp.output = append(comp.output, value)
			comp.ip += 2
		case 6: // bdv
			denominator := 1 << comp.getOperandValue(operand) // 2^operand
			if denominator == 0 {
				comp.reg.B = 0
			} else {
				comp.reg.B = comp.reg.A / denominator
			}
			comp.ip += 2
		case 7: // cdv
			denominator := 1 << comp.getOperandValue(operand) // 2^operand
			if denominator == 0 {
				comp.reg.C = 0
			} else {
				comp.reg.C = comp.reg.A / denominator
			}
			comp.ip += 2
		default:
			// Недопустимый opcode, остановить программу
			return
		}
	}
}

// parseInput читает входной файл и возвращает начальные регистры и программу
func parseInput(filename string) (Registers, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Registers{}, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reg := Registers{}
	var program []int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register A:") {
			valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Register A:"))
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return Registers{}, nil, err
			}
			reg.A = value
		} else if strings.HasPrefix(line, "Register B:") {
			valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Register B:"))
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return Registers{}, nil, err
			}
			reg.B = value
		} else if strings.HasPrefix(line, "Register C:") {
			valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Register C:"))
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return Registers{}, nil, err
			}
			reg.C = value
		} else if strings.HasPrefix(line, "Program:") {
			programStr := strings.TrimSpace(strings.TrimPrefix(line, "Program:"))
			programParts := strings.Split(programStr, ",")
			for _, part := range programParts {
				num, err := strconv.Atoi(strings.TrimSpace(part))
				if err != nil {
					return Registers{}, nil, err
				}
				program = append(program, num)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Registers{}, nil, err
	}

	return reg, program, nil
}

// solvePartOne запускает программу с начальными регистрами и возвращает строку вывода
func solvePartOne(reg Registers, program []int) (string, time.Duration) {
	comp := NewComputer(program, reg)
	start := time.Now()
	comp.Execute()
	duration := time.Since(start)

	// Преобразовать вывод в строку, разделённую запятыми
	outputStrs := []string{}
	for i, val := range comp.output {
		if i > 0 {
			outputStrs = append(outputStrs, ",")
		}
		outputStrs = append(outputStrs, strconv.Itoa(val))
	}
	output := strings.Join(outputStrs, "")

	return output, duration
}

// solvePartTwo находит наименьшее положительное начальное значение A, которое заставляет программу вывести саму себя
func solvePartTwo(initialReg Registers, program []int, programStr string) (int, time.Duration) {
	startTime := time.Now()
	a := 1
	for {
		// Инициализировать регистры с текущим A
		reg := Registers{
			A: a,
			B: initialReg.B,
			C: initialReg.C,
		}
		comp := NewComputer(program, reg)
		comp.Execute()

		// Преобразовать вывод в строку, разделённую запятыми
		outputStrs := []string{}
		for i, val := range comp.output {
			if i > 0 {
				outputStrs = append(outputStrs, ",")
			}
			outputStrs = append(outputStrs, strconv.Itoa(val))
		}
		output := strings.Join(outputStrs, "")

		// Проверить, совпадает ли вывод со строкой программы
		if output == programStr {
			duration := time.Since(startTime)
			return a, duration
		}

		a++

		// Чтобы предотвратить бесконечный цикл, можно установить разумный верхний предел
		if a > 1000000 {
			// Вернуть -1, чтобы указать, что не найдено в пределах лимита
			return -1, time.Since(startTime)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Пожалуйста, укажите входной файл в качестве аргумента командной строки.")
		return
	}

	filename := os.Args[1]
	reg, program, err := parseInput(filename)
	if err != nil {
		fmt.Printf("Ошибка при чтении входного файла: %v\n", err)
		return
	}

	// Часть Первая
	outputPartOne, durationPartOne := solvePartOne(reg, program)
	fmt.Printf("Часть Первая Ответ: %s\n", outputPartOne)
	fmt.Printf("Часть Первая Время: %v\n", durationPartOne)

	// Подготовить строку программы для сравнения во Второй части
	programStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(program)), ","), "[]")

	// Часть Вторая
	initialA, durationPartTwo := solvePartTwo(reg, program, programStr)
	if initialA != -1 {
		fmt.Printf("Часть Вторая Ответ: %d\n", initialA)
	} else {
		fmt.Println("Часть Вторая Ответ: Не найдено в пределах установленного лимита поиска.")
	}
	fmt.Printf("Часть Вторая Время: %v\n", durationPartTwo)
}