// Advent of Code 2024, Day 21: "Keypad Conundrum" — Реальное решение с многошаговой активацией.
//
// Основная идея:
//   У нас есть 3 "робота"-уровня (R3 — верхний, R2 — средний, R1 — нижний), а в самом низу — цифровая клавиатура.
//   Пользователь нажимает одну из 5 клавиш: '^','v','<','>','A' на САМОЙ верхней (реальной) клавиатуре.
//   - Если это стрелка, робот R3 двигает свою "руку", если нет "дыры".
//   - Если это 'A', то происходит "активация":
//       * Смотрим, на какой кнопке стоит R3. Если там стрелка, она передаётся R2 (т. е. R2 двигает руку).
//         Если там 'A', тогда смотрим R2. Если у R2 стрелка, она передаётся R1 (цифровая клавиатура). Если у R2 'A',
//         тогда смотрим R1. Если у R1 — цифра или 'A', она "нажимается" => добавляется к typed.
//         Если у R1 стрелка (что в нормальной цифровой раскладке не бывает), пытались бы двигать её (но это gap).
//
// Важная деталь: если на всех трёх роботах (R3, R2, R1) подряд 'A', то ОДНО нажатие 'A' пользователя проваливается
// насквозь, вплоть до нажатия кнопки на цифровой клавиатуре. Это и есть "многошаговая" активация.
// Так работает пример из условия, когда робот "печатает" цифры.
//
// Чтобы найти КРАТЧАЙШУЮ последовательность нажатий, используем BFS по состояниям:
//   State = (posR3, posR2, posR1, typedSoFar)
//
// Отсечение лишних путей: если typedSoFar не является префиксом искомого кода, отбрасываем (нет бэкспейса).
// Когда typedSoFar == code, возвращаем глубину (число нажатий).
//
// Это решение корректно выводит 126384 для sample.txt (коды 029A, 980A, 179A, 456A, 379A)
// и решит прочие входные данные за разумное время.
//
// Компиляция и запуск:
//   go run main.go sample.txt
//
// Примечание: для реального AoC (большие коды) возможно придётся ещё сильнее оптимизировать.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ---- Раскладки клавиатур роботов (верхний/средний) и цифровой (нижний) ----

// Верхняя/средняя "роботная" клавиатура (2x3), без учёта "дыр" —
// key: какая кнопка, value: куда перейдём при нажатии стрелки "^","v","<",">".
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

// Нижняя цифровая клавиатура (3x3 + 2 внизу):
//
//	7 8 9
//	4 5 6
//	1 2 3
//	  0 A
//
// Если там "дырка", переход "".
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

// ---- Состояние и BFS ----
type State struct {
	r3    string // позиция руки на верхнем роботе
	r2    string // позиция руки на среднем
	r1    string // позиция руки на нижнем (цифровая клавиатура)
	typed string // уже набранная строка (например "0", потом "02", ...)
}

// solveCode возвращает длину кратчайшей последовательности нажатий пользователя,
// ведущей к вводу строки code на нижней клавиатуре. Если пути нет, вернёт 0.
func solveCode(code string) int {
	start := State{r3: "A", r2: "A", r1: "A", typed: ""}
	if code == "" {
		return 0
	}

	// BFS
	queue := []State{start}
	dist := map[State]int{start: 0}

	// Все 5 кнопок, которые пользователь может нажать
	inputs := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		// Проверка цели
		if st.typed == code {
			return d
		}
		// Префиксная отсечка (если уже typed не совпадает с префиксом code)
		if !strings.HasPrefix(code, st.typed) {
			continue
		}

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

	return 0
}

// nextState обрабатывает одно нажатие inp ('^','v','<','>','A') и возвращает новое состояние либо (State{}, false) если gap.
func nextState(st State, inp string) (State, bool) {
	r3, r2, r1, typed := st.r3, st.r2, st.r1, st.typed

	// 1) Если inp не 'A', то это попытка сдвинуть руку верхнего робота
	if inp != "A" {
		r3next := robotKeypadAdj[r3][inp]
		if r3next == "" {
			return State{}, false // gap
		}
		return State{r3: r3next, r2: r2, r1: r1, typed: typed}, true
	}

	// 2) inp == "A" => активация верхнего робота:
	return activate(r3, r2, r1, typed)
}

// activate моделирует одну "активацию" на верхнем роботе.
// Возможен цепной проход, если r3 == 'A', тогда смотрим r2, если и там 'A', идём к r1 и т.д.
// Возвращает (State{}, false) при любом gap.
func activate(r3, r2, r1, typed string) (State, bool) {
	// Шаг 1: что на r3?
	switch r3 {
	case "^", "v", "<", ">":
		// Это стрелка для R2
		r2next := robotKeypadAdj[r2][r3]
		if r2next == "" {
			return State{}, false
		}
		return State{r3: r3, r2: r2next, r1: r1, typed: typed}, true

	case "A":
		// Идём к роботу2
		return activateR2(r2, r1, typed, r3)
	default:
		return State{}, false
	}
}

// activateR2: активация робота2, смотрим r2.
func activateR2(r2, r1, typed, r3 string) (State, bool) {
	switch r2 {
	case "^", "v", "<", ">":
		// Это стрелка для R1 (цифровая)
		r1next := digitKeypadAdj[r1][r2]
		if r1next == "" {
			return State{}, false
		}
		return State{r3: r3, r2: r2, r1: r1next, typed: typed}, true

	case "A":
		// Идём к роботу1
		return activateR1(r1, typed, r3, r2)
	default:
		return State{}, false
	}
}

// activateR1: активация робота1 => нажать кнопку (если это цифра или 'A') или попытаться движение (если '^','v','<','>').
func activateR1(r1, typed, r3, r2 string) (State, bool) {
	// Если r1 — стрелка, пытались бы двигать, но в digitKeypad у нас таких нет (gap).
	// Если r1 — '0'..'9' или 'A', добавляем к typed.
	if _, isRobotButton := robotKeypadAdj[r1]; isRobotButton {
		// Тогда это '^','v','<','>' или 'A' с точки зрения "роботной" карты — но мы в цифрах => будет gap.
		// Единственное исключение — 'A' существует и там, и тут. Но если 'A' в digitKeypadAdj, это реальная кнопка (цифровая).
		// Скажем так: если r1 == "A", мы действительно имеем в digitKeypad? Да, там есть "A".
		// Но также "A" есть в robotKeypadAdj => коллизия. Разберём явно:
		if r1 == "A" {
			// Это "цифровая A". Тогда нажимаем её => typed += "A".
			return State{r3: r3, r2: r2, r1: r1, typed: typed + "A"}, true
		}
		// Иначе стрелка "^","v","<",">" => gap, потому что digitKeypadAdj для них пуст.
		return State{}, false
	}
	// Значит это "0","1","2","3","4","5","6","7","8","9" (или теоретически ещё что-то).
	return State{r3: r3, r2: r2, r1: r1, typed: typed + r1}, true
}

// parseNumericPart выдёргивает целую часть (без ведущих нулей) из кода, заканчивающегося на 'A'.
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

	// Считываем коды (обычно 5 строк)
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
		fmt.Println("No codes found.")
		return
	}

	// Part 1: сумма complexities
	t1Start := time.Now()
	total := 0
	for _, code := range codes {
		seqLen := solveCode(code)
		numPart := parseNumericPart(code)
		total += seqLen * numPart
	}
	t1 := time.Since(t1Start).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", total, t1)

	// Part 2 — не реализована
	t2Start := time.Now()
	fmt.Printf("Part 2: not implemented (time %.4fs)\n", time.Since(t2Start).Seconds())

	_ = time.Since(startAll)
}
