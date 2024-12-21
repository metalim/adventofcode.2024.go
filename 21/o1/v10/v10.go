// Advent of Code 2024, Day 21: "Keypad Conundrum"
// -----------------------------------------------
// Реальное решение с учётом, что в задаче у нас три уровня "роботных" клавиатур (оба верхних —
// "directional keypad", нижний — тоже "directional keypad", но управляющий цифровой клавиатурой),
// а в самом низу — сама цифровая клавиатура (с цифрами и кнопкой 'A').
//
// Главный нюанс: в условии сказано, что каждая "directional keypad" имеет такую раскладку:
//
//     +---+---+
//     | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+
//
// При этом "рука" робота изначально указывает на 'A' (в правом верхнем углу),
// и мы можем нажимать всего 5 клавиш: '^','v','<','>','A'.
//
// Внизу у нас цифровая клавиатура:
//
//   +---+---+---+
//   | 7 | 8 | 9 |
//   +---+---+---+
//   | 4 | 5 | 6 |
//   +---+---+---+
//   | 1 | 2 | 3 |
//   +---+---+---+
//       | 0 | A |
//       +---+---+
//
// Где "рука" нижнего робота тоже изначально на 'A' (внизу справа).
// Один "нажим" пользователя на верхней клавиатуре ('^','v','<','>','A')
// может пробить несколько уровней вниз, если везде подряд 'A'.
//
// Нужно для каждого кода (например, "029A") найти кратчайшую
// последовательность нажатий на верхней клавиатуре,
// чтобы внизу был набран нужный код.
// Complexity кода = (длина такой последовательности) * (числовая часть кода).
// Суммируем по всем пяти кодам.
//
// Эта версия кода учитывает реальную схему переходов ("adjacency")
// для каждой из трёх роботных клавиатур (они одинаковые) —
// сделан так, чтобы действительно можно было сдвигаться
// (например, из верхнего правого 'A' по нажатию '<' перейти в '^', и т.д.).
//
// Также применяется "префиксная" отсечка: если уже набранная строка
// не совпадает с префиксом искомого кода, путь отбрасывается.
//
// В результате для sample.txt (5 кодов из условия) будет получаться 126384,
// а не 0.
//
// ------------------------------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------------------
// 1) Раскладка "directional keypad" (для верхнего и среднего робота).
//    Она выглядит как:
//
//       (row=0 col=1) '^'   (row=0 col=2) 'A'
//   (row=1 col=0) '<'   (row=1 col=1) 'v'   (row=1 col=2) '>'
//
//    Причём (row=0 col=0) — gap.
//
//    С учётом стрелок (up, down, left, right), задаём переходы так, чтобы:
//      - Из 'A' при нажатии '<' => '^', при нажатии 'v' => '>', и т.д.
//      - Там, где "дырка", будет пустая строка, означающая gap.
//
//    userInput: "^","v","<",">". Если нажали 'A' — это уже "активация" (обрабатывается отдельно).

var robotKeypadAdj = map[string]map[string]string{
	"^": {
		"^": "",  // вверх из '^' нет
		"v": "v", // вниз => 'v'
		"<": "",  // влево => gap
		">": "A", // вправо => 'A'
	},
	"A": {
		"^": "",  // вверх => gap
		"v": ">", // вниз => '>'
		"<": "^", // влево => '^'
		">": "",  // вправо => gap
	},
	"<": {
		"^": "",  // вверх => gap
		"v": "",  // вниз => gap (стрелка вниз из '<' идёт за поле)
		"<": "",  // влево => край
		">": "v", // вправо => 'v'
	},
	"v": {
		"^": "^", // вверх => '^'
		"v": "",  // вниз => gap
		"<": "<", // влево => '<'
		">": ">", // вправо => '>'
	},
	">": {
		"^": "A", // вверх => 'A'
		"v": "",  // вниз => gap
		"<": "v", // влево => 'v'
		">": "",  // вправо => gap
	},
}

// ----------------------------------------------------------
// 2) Раскладка "цифровой" клавиатуры, где рука может стоять
//    на 7,8,9,4,5,6,1,2,3,0,A.
//    Пример:
//        7 8 9
//        4 5 6
//        1 2 3
//          0 A
//    Аналогично задаём, куда попадём по нажатию '^','v','<','>'.

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

// ----------------------------------------------------------
// 3) Описание состояния BFS:
//    r3: позиция руки верхнего робота (строки "^","v","<",">","A")
//    r2: позиция руки среднего робота (то же)
//    r1: позиция руки нижнего робота (клавиши "7","8","9",...,"A")
//    typed: что уже набрали (цифры + 'A')

type State struct {
	r3    string
	r2    string
	r1    string
	typed string
}

// ----------------------------------------------------------
// 4) BFS для одного кода: найдём длину кратчайшей последовательности нажатий
//    на верхней клавиатуре, чтобы в итоге в typed был code.

func solveCode(code string) int {
	start := State{"A", "A", "A", ""}
	if code == "" {
		return 0
	}

	dist := map[State]int{start: 0}
	queue := []State{start}
	userButtons := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		d := dist[st]

		// Если достигли нужного набранного кода — вернём длину
		if st.typed == code {
			return d
		}
		// Префиксная отсечка: если текущий typed не начало искомого кода — отбрасываем
		if !strings.HasPrefix(code, st.typed) {
			continue
		}

		// Пытаемся нажать каждую из 5 кнопок
		for _, inp := range userButtons {
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

// nextState обрабатывает одно нажатие пользователя inp.
// Если inp — это стрелка, двигаем руку верхнего робота (r3).
// Если inp — 'A', то делаем многоступенчатую активацию r3 -> r2 -> r1.
// Возвращаем (State{}, false), если упираемся в "дырку".
func nextState(st State, inp string) (State, bool) {
	r3, r2, r1, typed := st.r3, st.r2, st.r1, st.typed

	if inp != "A" {
		// Движение верхнего робота
		r3next := robotKeypadAdj[r3][inp]
		if r3next == "" {
			return State{}, false
		}
		return State{r3: r3next, r2: r2, r1: r1, typed: typed}, true
	}

	// inp == "A" => активация верхнего робота
	return activate(r3, r2, r1, typed)
}

// activate моделирует "нажатие A" на r3.
// Если r3 — стрелка, она передаётся роботу r2.
// Если r3 — 'A', идём смотреть r2 и т.д.
func activate(r3, r2, r1, typed string) (State, bool) {
	switch r3 {
	case "^", "v", "<", ">":
		// Передаём эту стрелку роботу r2
		r2next := robotKeypadAdj[r2][r3]
		if r2next == "" {
			return State{}, false
		}
		// Остальное без изменений
		return State{r3: r3, r2: r2next, r1: r1, typed: typed}, true

	case "A":
		// Тогда активируем r2
		return activateR2(r2, r1, typed, r3)
	default:
		return State{}, false
	}
}

// activateR2: если r2 — стрелка, она передаётся r1 (цифровая клавиатура);
// если r2 — 'A', то нажимаем r1 (если это цифра/А) или пытаемся двигать (если это стрелка).
func activateR2(r2, r1, typed, r3 string) (State, bool) {
	switch r2 {
	case "^", "v", "<", ">":
		// Движение r1 на цифровой
		r1next := digitKeypadAdj[r1][r2]
		if r1next == "" {
			return State{}, false
		}
		return State{r3: r3, r2: r2, r1: r1next, typed: typed}, true

	case "A":
		// Тогда активируем r1
		return activateR1(r1, typed, r3, r2)

	default:
		return State{}, false
	}
}

// activateR1: нажимаем кнопку r1, если это '0'..'9' или 'A'.
// Если это '^','v','<','>', то пытались бы двигать внутри digitKeypad (но там их нет => gap).
func activateR1(r1, typed, r3, r2 string) (State, bool) {
	// Если r1 — одна из {'^','v','<','>'}, в digitKeypadAdj такого нет => gap
	// Но если r1 == "A" или цифра => добавляем к typed.
	if _, isDirKey := robotKeypadAdj[r1]; isDirKey && r1 != "A" {
		// То это "^","v","<",">" => gap в цифровой
		return State{}, false
	}
	// Иначе r1 — "0","1","2","...","9","A" => нажимаем => typed + r1
	return State{r3: r3, r2: r2, r1: r1, typed: typed + r1}, true
}

// parseNumericPart выдёргивает integer из кода "XYZ...A". "029A" => 29, "456A" => 456, "000A" => 0.
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

	// Part 1: вычисляем сумму complexities
	startPart1 := time.Now()
	sumComplexity := 0
	for _, code := range codes {
		seqLen := solveCode(code)
		numVal := parseNumericPart(code)
		sumComplexity += seqLen * numVal
	}
	t1 := time.Since(startPart1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", sumComplexity, t1)

	// Part 2 — не реализовано
	startPart2 := time.Now()
	fmt.Printf("Part 2: not implemented (time %.4fs)\n", time.Since(startPart2).Seconds())

	_ = time.Since(startAll)
}
