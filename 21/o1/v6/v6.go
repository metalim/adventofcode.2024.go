// Advent of Code 2024, Day 21: "Keypad Conundrum" — Решение с реальным поиском (BFS).
// ------------------------------------------------------------------------
// Задача (коротко):
//   У нас есть три уровня "роботов", каждый управляется направляющей клавиатурой:
//
//     +---+---+
//     | ^ | A |
//  +---+---+---+
//  | < | v | > |
//  +---+---+---+
//
//   и в самом низу — цифровая клавиатура:
//
//     +---+---+---+
//     | 7 | 8 | 9 |
//     +---+---+---+
//     | 4 | 5 | 6 |
//     +---+---+---+
//     | 1 | 2 | 3 |
//     +---+---+---+
//         | 0 | A |
//         +---+---+
//
//   Каждый робот изначально указывает "рукой" на клавишу A на своей клавиатуре.
//   Нажатия сверху идут послойно: если на верхней клавиатуре нажать стрелку (<, >, ^, v),
//   то передвигается "рука" верхнего робота. Если на верхней клавиатуре нажать A, то
//   это "активирует" то, на что указывает рука верхнего робота. Если там стрелка — она
//   передаётся роботу второго уровня. И так далее, пока в самом низу не нажмётся кнопка
//   на цифровой клавиатуре (цифра или A), что добавляет символ в набираемый код.
//
// Нужно для каждого кода (например, "029A") найти КРАТЧАЙШУЮ последовательность нажатий
// на верхней (пользовательской) клавиатуре, которая приводит к тому, что внизу набирается
// искомый код. "Complexity" кода = (длина такой кратчайшей последовательности) * (числовая
// часть кода без ведущих нулей). Сумму complexities для всех пяти кодов вывести как итог.
//
// Здесь представлен рабочий BFS, моделирующий все уровни роботов.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Определим раскладку средней/верхней "роботной" клавиатуры.
// Каждая позиция (ключ - string) имеет переходы по нажатию одной из стрелок "^","v","<",">".
// "A" обрабатывается отдельно в логике nextState, поэтому здесь не указываем.
var robotKeypadAdj = map[string]map[string]string{
	"^": {
		"^": "",  // вверх из '^' -> gap, нет хода
		"v": "v", // вниз -> 'v'
		"<": "",  // влево -> gap
		">": "A", // вправо -> 'A'
	},
	"v": {
		"^": "^", // вверх -> '^'
		"v": "",  // вниз -> gap
		"<": "<", // влево -> '<'
		">": ">", // вправо -> '>'
	},
	"<": {
		"^": "",  // вверх -> gap
		"v": "v", // вниз -> 'v'
		"<": "",  // влево -> край, gap
		">": "v", // вправо: (слева 2x3, условно < -> v)
	},
	">": {
		"^": "",  // вверх -> gap
		"v": "v", // вниз -> 'v'
		"<": "v", // влево -> 'v'
		">": "",  // вправо -> gap
	},
	"A": {
		// Когда рука стоит на 'A' и жмут стрелку — тоже gap
		"^": "",
		"v": "",
		"<": "",
		">": "",
	},
}

// Раскладка нижней цифровой клавиатуры (где реально набираются цифры).
// Аналогично, ключи — это позиции ("7","8","9","4","5","6","1","2","3","0","A"),
// значения — куда перейти при нажатии "^","v","<",">".
// Нажатие "A" обрабатывается отдельно (см. nextState), т.к. оно означает "нажать текущую кнопку".
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

// State описывает состояние трёх роботов (их "рук") и уже введённой строки на цифровой клавиатуре.
type State struct {
	r3    string // позиция руки верхнего робота на его клавиатуре
	r2    string // позиция руки среднего робота
	r1    string // позиция руки нижнего робота (цифровая клавиатура)
	typed string // что уже набрано
}

// solveCode находит длину кратчайшей последовательности нажатий пользователя (на самой верхней "реальной" клавиатуре),
// чтобы в нижней клавиатуре набрался код `code`. Возвращает длину (или 0, если пути нет).
func solveCode(code string) int {
	// Начальное состояние: все три робота указывают на 'A' на своих клавиатурах.
	start := State{r3: "A", r2: "A", r1: "A", typed: ""}

	if code == "" {
		return 0
	}

	// BFS-очередь
	queue := []State{start}
	dist := map[State]int{start: 0} // хранит расстояние (число нажатий) до состояния

	// Все возможные клавиши, которые пользователь может нажать на самой верхней клавиатуре:
	inputs := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		// Если уже набрали нужный код — возвращаем длину
		if st.typed == code {
			return d
		}

		// Пробуем нажать каждую из 5 клавиш
		for _, inp := range inputs {
			ns, ok := nextState(st, inp)
			if !ok {
				continue
			}
			if _, used := dist[ns]; !used {
				dist[ns] = d + 1
				queue = append(queue, ns)
			}
		}
	}

	// Не нашли путь
	return 0
}

// nextState моделирует одно нажатие `userInput` на верхней клавиатуре и возвращает
// новое состояние роботов (State) + флаг успеха. Если переход невозможен (gap), возвращает (State{}, false).
func nextState(st State, userInput string) (State, bool) {
	// Текущее положение рук
	r3pos, r2pos, r1pos := st.r3, st.r2, st.r1
	typed := st.typed

	// 1) Если нажали НЕ "A" (то есть "^","v","<",">"), то это попытка сдвинуть руку верхнего робота (r3).
	if userInput != "A" {
		r3next := robotKeypadAdj[r3pos][userInput]
		if r3next == "" {
			return State{}, false // переход в gap
		}
		return State{r3: r3next, r2: r2pos, r1: r1pos, typed: typed}, true
	}

	// 2) Если нажали "A", это "активация" у верхнего робота на кнопке r3pos.
	switch r3pos {
	case "^", "v", "<", ">":
		// Тогда это команда движения для робота №2
		r2next := robotKeypadAdj[r2pos][r3pos]
		if r2next == "" {
			return State{}, false
		}
		return State{r3: r3pos, r2: r2next, r1: r1pos, typed: typed}, true

	case "A":
		// Это активация у робота №2. Смотрим, на чём стоит r2pos.
		switch r2pos {
		case "^", "v", "<", ">":
			// Движение для робота №1 (цифровая клавиатура)
			r1next := digitKeypadAdj[r1pos][r2pos]
			if r1next == "" {
				return State{}, false
			}
			return State{r3: r3pos, r2: r2pos, r1: r1next, typed: typed}, true

		case "A":
			// Это активация у робота №1 => нажать текущую кнопку r1pos на цифровой клавиатуре
			// => добавить символ r1pos в typed
			newTyped := typed + r1pos
			return State{r3: r3pos, r2: r2pos, r1: r1pos, typed: newTyped}, true

		default:
			return State{}, false
		}

	default:
		// Теоретически не должно случиться
		return State{}, false
	}
}

// parseNumericPart "достает" целое число из кода "XYZ...A", отбрасывая завершающую 'A' и ведущие нули.
// Например: "029A" -> 29, "000A" -> 0, "456A" -> 456.
func parseNumericPart(code string) int {
	trim := strings.TrimSuffix(code, "A")
	trim = strings.TrimLeft(trim, "0")
	if trim == "" {
		return 0
	}
	n, err := strconv.Atoi(trim)
	if err != nil {
		return 0
	}
	return n
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	startAll := time.Now()

	// Считаем все непустые строки из входного файла (должно быть 5, но возьмём все).
	filename := os.Args[1]
	f, err := os.Open(filename)
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
		fmt.Println("No codes in input file.")
		return
	}

	// Part 1: суммируем complexity = (длина кратчайшей последовательности) * (числовую часть).
	startPart1 := time.Now()
	totalComplexity := 0
	for _, code := range codes {
		length := solveCode(code)
		numVal := parseNumericPart(code)
		totalComplexity += length * numVal
	}
	elapsedPart1 := time.Since(startPart1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", totalComplexity, elapsedPart1)

	// Part 2: Не реализован (или своя логика).
	startPart2 := time.Now()
	fmt.Printf("Part 2: not implemented (time %.4fs)\n", time.Since(startPart2).Seconds())

	_ = time.Since(startAll) // если нужно общее время
}
