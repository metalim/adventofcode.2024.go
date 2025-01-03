package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// OperandType определяет тип операнда: Литеральный или Комбо
type OperandType int

const (
	LITERAL OperandType = iota
	COMBO
)

// Registers структура для хранения регистров A, B, C
type Registers struct {
	A int64
	B int64
	C int64
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

// getOperandType возвращает тип операнда (Литеральный или Комбо) в зависимости от opcode
func getOperandType(opcode int) OperandType {
	switch opcode {
	case 1, 3: // bxl и jnz используют Литеральные операнды
		return LITERAL
	default: // Все остальные инструкции используют Комбо операнды
		return COMBO
	}
}

// getOperandValue возвращает значение операнда на основе его типа
func (comp *Computer) getOperandValue(operand int, opType OperandType) int64 {
	if opType == LITERAL {
		return int64(operand)
	}
	// Комбо операнд
	switch operand {
	case 0, 1, 2, 3:
		return int64(operand)
	case 4:
		return comp.reg.A
	case 5:
		return comp.reg.B
	case 6:
		return comp.reg.C
	case 7:
		// Операнд 7 зарезервирован и не должен появляться, но если появляется, возвращаем 7
		return 7
	default:
		// Недопустимый операнд, возвращаем 0
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
		opType := getOperandType(opcode)

		switch opcode {
		case 0: // adv
			denominator := int64(1) << comp.getOperandValue(operand, opType)
			if denominator == 0 {
				comp.reg.A = 0
			} else {
				comp.reg.A = comp.reg.A / denominator
			}
			comp.ip += 2
		case 1: // bxl
			comp.reg.B = comp.reg.B ^ comp.getOperandValue(operand, opType)
			comp.ip += 2
		case 2: // bst
			comp.reg.B = comp.getOperandValue(operand, opType) % 8
			comp.ip += 2
		case 3: // jnz
			if comp.reg.A != 0 {
				comp.ip = int(comp.getOperandValue(operand, opType))
			} else {
				comp.ip += 2
			}
		case 4: // bxc
			comp.reg.B = comp.reg.B ^ comp.reg.C
			comp.ip += 2
		case 5: // out
			value := comp.getOperandValue(operand, opType) % 8
			comp.output = append(comp.output, int(value))
			comp.ip += 2
		case 6: // bdv
			denominator := int64(1) << comp.getOperandValue(operand, opType)
			if denominator == 0 {
				comp.reg.B = 0
			} else {
				comp.reg.B = comp.reg.A / denominator
			}
			comp.ip += 2
		case 7: // cdv
			denominator := int64(1) << comp.getOperandValue(operand, opType)
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
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return Registers{}, nil, err
			}
			reg.A = value
		} else if strings.HasPrefix(line, "Register B:") {
			valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Register B:"))
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return Registers{}, nil, err
			}
			reg.B = value
		} else if strings.HasPrefix(line, "Register C:") {
			valueStr := strings.TrimSpace(strings.TrimPrefix(line, "Register C:"))
			value, err := strconv.ParseInt(valueStr, 10, 64)
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

// solvePartTwoParallel находит наименьшее положительное начальное значение A, которое заставляет программу вывести саму себя, используя параллелизм
func solvePartTwoParallel(initialReg Registers, program []int, programStr string) (int64, time.Duration) {
	startTime := time.Now()
	found := make(chan int64)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Количество параллельных горутин, можно настроить в зависимости от CPU
	numWorkers := 8

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			a := int64(workerID + 1)
			for {
				select {
				case <-done:
					return
				default:
					// Продолжаем поиск
				}

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
				for j, val := range comp.output {
					if j > 0 {
						outputStrs = append(outputStrs, ",")
					}
					outputStrs = append(outputStrs, strconv.Itoa(val))
				}
				output := strings.Join(outputStrs, "")

				// Проверить, совпадает ли вывод со строкой программы
				if output == programStr {
					select {
					case found <- a:
						// Найдено, сигнализируем другим горутинам остановиться
						close(done)
						return
					default:
						return
					}
				}

				a += int64(numWorkers)
			}
		}(i)
	}

	// Ожидание результата или завершение всех горутин
	var result int64 = -1
	select {
	case a := <-found:
		result = a
		close(done) // Закрываем канал done, чтобы остановить остальные горутины
	case <-time.After(10 * time.Minute): // Установите таймаут по необходимости
		close(done) // Таймаут достигнут, прекращаем поиск
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	duration := time.Since(startTime)

	return result, duration
}

// solvePartTwo находит наименьшее положительное начальное значение A, которое заставляет программу вывести саму себя
func solvePartTwo(initialReg Registers, program []int, programStr string) (int64, time.Duration) {
	a, duration := solvePartTwoParallel(initialReg, program, programStr)
	return a, duration
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
