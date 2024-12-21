// Advent of Code 2024, Day 21: "Keypad Conundrum" — Решение (BFS) с отсечением неподходящих префиксов.
//
// Описание задачи см. в предыдущих версиях.
// Проблема: полный BFS без оптимизаций может взрывообразно расти.
// Причина: при нажатии 'A' на нижнем роботе мы можем набрать любую цифру, тем самым быстро
// создавая множество лишних "typed"-последовательностей.
//
// Оптимизация: префиксная фильтрация (prefix pruning).
//   - Цель: набрать ровно code (например, "029A").
//   - Если текущая набранная строка typed не совпадает с префиксом code, дальнейший путь бессмыслен.
//     (т.к. нет "Backspace" и продолжать набор уже неверной строки нельзя).
//   - Если typed уже достигла длины code, но при этом typed != code, тоже отсекаем.
//   - Как только typed == code, мы нашли решение.
//
// Сокращает пространство поиска очень сильно, что спасает от зависания.
//
// Алгоритм:
//   1) BFS по состояниям (r3pos, r2pos, r1pos, typed).
//   2) При формировании nextState, если typed не префикс code, отбрасываем.
//   3) Как только typed == code, возвращаем глубину.
//
// Этого достаточно, чтобы выполнять быстро даже на полном переборе.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Раскладка средней/верхней "роботной" клавиатуры
// (ключи и значения — строки "^", "v", "<", ">", "A").
var robotKeypadAdj = map[string]map[string]string{
	"^": {
		"^": "",
		"v": "v",
		"<": "",
		">": "A",
	},
	"v": {
		"^": "^",
		"v": "",
		"<": "<",
		">": ">",
	},
	"<": {
		"^": "",
		"v": "v",
		"<": "",
		">": "v",
	},
	">": {
		"^": "",
		"v": "v",
		"<": "v",
		">": "",
	},
	"A": {
		"^": "",
		"v": "",
		"<": "",
		">": "",
	},
}

// Раскладка нижней цифровой клавиатуры.
var digitKeypadAdj = map[string]map[string]string{
	"7": {
		"^": "",
		"v": "4",
		"<": "",
		">": "8",
	},
	"8": {
		"^": "",
		"v": "5",
		"<": "7",
		">": "9",
	},
	"9": {
		"^": "",
		"v": "6",
		"<": "8",
		">": "",
	},
	"4": {
		"^": "7",
		"v": "1",
		"<": "",
		">": "5",
	},
	"5": {
		"^": "8",
		"v": "2",
		"<": "4",
		">": "6",
	},
	"6": {
		"^": "9",
		"v": "3",
		"<": "5",
		">": "",
	},
	"1": {
		"^": "4",
		"v": "",
		"<": "",
		">": "2",
	},
	"2": {
		"^": "5",
		"v": "0",
		"<": "1",
		">": "3",
	},
	"3": {
		"^": "6",
		"v": "A",
		"<": "2",
		">": "",
	},
	"0": {
		"^": "2",
		"v": "",
		"<": "",
		">": "A",
	},
	"A": {
		"^": "3",
		"v": "",
		"<": "0",
		">": "",
	},
}

// State описывает состояние трёх роботов (r3 — верхний, r2 — средний, r1 — нижний) и текущую набранную строку.
type State struct {
	r3    string
	r2    string
	r1    string
	typed string
}

// solveCode находит длину кратчайшей последовательности нажатий для набора code.
// Если пути нет, возвращает 0.
func solveCode(code string) int {
	start := State{r3: "A", r2: "A", r1: "A", typed: ""}

	// Быстрый случай: пустой код
	if code == "" {
		return 0
	}

	// BFS
	queue := []State{start}
	dist := map[State]int{start: 0}
	inputs := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		// Если уже набрано ровно нужное — готово
		if st.typed == code {
			return d
		}

		// Пробуем все входы
		for _, inp := range inputs {
			ns, ok := nextState(st, inp, code)
			if !ok {
				continue
			}
			// Если ns.typed длиннее code или не совпадает с префиксом, пропустим
			if len(ns.typed) > len(code) {
				continue
			}
			if !strings.HasPrefix(code, ns.typed) {
				continue
			}
			if _, used := dist[ns]; !used {
				dist[ns] = d + 1
				queue = append(queue, ns)
			}
		}
	}

	// Не нашли
	return 0
}

// nextState моделирует одно нажатие inp на верхней клавиатуре, возвращает новое состояние + флаг успеха (false, если gap).
// Параметр code нужен только чтобы при наборе "A" (внизу) проверять префикс или что-то ещё, но в данном случае не обязателен.
// Однако используем его в префиксной проверке выше, поэтому здесь он лишь для контекста (при желании можно убрать).
func nextState(st State, inp string, code string) (State, bool) {
	r3pos, r2pos, r1pos := st.r3, st.r2, st.r1
	typed := st.typed

	if inp != "A" {
		// Движение руки верхнего робота
		r3next := robotKeypadAdj[r3pos][inp]
		if r3next == "" {
			return State{}, false
		}
		return State{r3: r3next, r2: r2pos, r1: r1pos, typed: typed}, true
	}

	// inp == "A": активация на уровне r3pos
	switch r3pos {
	case "^", "v", "<", ">":
		// Это движение для робота №2
		r2next := robotKeypadAdj[r2pos][r3pos]
		if r2next == "" {
			return State{}, false
		}
		return State{r3: r3pos, r2: r2next, r1: r1pos, typed: typed}, true

	case "A":
		// Активация робота №2 -> смотрим r2pos
		switch r2pos {
		case "^", "v", "<", ">":
			// Движение для робота №1
			r1next := digitKeypadAdj[r1pos][r2pos]
			if r1next == "" {
				return State{}, false
			}
			return State{r3: r3pos, r2: r2pos, r1: r1next, typed: typed}, true
		case "A":
			// Активация робота №1 -> нажать текущую кнопку r1pos => добавить символ
			return State{r3: r3pos, r2: r2pos, r1: r1pos, typed: typed + r1pos}, true
		default:
			return State{}, false
		}

	default:
		return State{}, false
	}
}

// parseNumericPart достаёт целую часть из кода вида "029A" => 29, "000A" => 0, "456A" => 456.
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	startAll := time.Now()

	// Считаем все коды (строки) из входного файла
	inputFile := os.Args[1]
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var codes []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			codes = append(codes, line)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	if len(codes) == 0 {
		fmt.Println("No codes in input.")
		return
	}

	// Part 1: считаем сумму complexity для всех кодов
	startPart1 := time.Now()
	sumComplexity := 0
	for _, code := range codes {
		seqLen := solveCode(code) // длина кратчайшей последовательности
		numVal := parseNumericPart(code)
		sumComplexity += seqLen * numVal
	}
	timePart1 := time.Since(startPart1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", sumComplexity, timePart1)

	// Part 2: не реализовано
	startPart2 := time.Now()
	fmt.Printf("Part 2: not implemented (time %.4fs)\n", time.Since(startPart2).Seconds())

	_ = time.Since(startAll)
}
