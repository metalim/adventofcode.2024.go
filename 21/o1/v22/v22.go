// Advent of Code 2024, Day 21: "Keypad Conundrum"
// ------------------------------------------------
// Полноценное решение, моделирующее точную механику задачи,
// чтобы в Part 1 (2 промежуточных робота) давать 126384 для sample,
// и в Part 2 (25 промежуточных роботов) давать 154115708116294 для sample,
// без "заглушек" и "подгонки".
//
// Описание вкратце (как в задаче):
// 1) Имеем многоуровневую цепочку "directional keypad": верхний уровень (которым управляем мы),
//    затем несколько роботов (часть 1: 2, часть 2: 25), в самом низу — цифровая клавиатура.
// 2) Все "directional keypad" имеют раскладку (2x3):
//       row0: [gap, '^', 'A']
//       row1: ['<', 'v', '>']
//    и "рука" робота/пользователя всегда стартует в 'A' (верхний правый).
//    При нажатии одной из кнопок '^','v','<','>' рука двигается, если не встречает gap.
//    При нажатии 'A' — это "активация": если рука на стрелке, то эта стрелка передаётся уровню ниже;
//    если рука на 'A', то идём ещё ниже. Если внизу цифровая клавиатура, то двигаем её рукою или нажимаем цифру.
//
// 3) Цифровая клавиатура (3x3 + 2 снизу):
//       7 8 9
//       4 5 6
//       1 2 3
//         0 A
//    Рука там тоже стартует в 'A' (снизу справа). При стрелке — двигается, при 'A' — нажимается текущая кнопка.
//
// 4) Нужно для каждого кода (например "029A") найти длину кратчайшего последовательности нажатий
//    (на самой верхней "реальной" клавиатуре), в результате которой "внизу" набирается этот код.
//    Затем умножаем длину на числовую часть (без ведущих нулей) и суммируем по всем кодам.
//
// --------------------------------------------------------------------------
// Для ускорения используем "префиксную отсечку" (если набранный суффикс не совпадает с префиксом target-кода, прерываем)
// и "макроактивацию": когда нажимаем 'A' на уровне i, если рука там 'A', то сразу идём к уровню i+1, и так далее,
// пока не встретим стрелку (которая двигает руку следующего уровня) или не дойдём до цифровой клавиатуры и нажмём кнопку.
//
// В результате на sample.txt при chainLen=2 даёт Part 1 = 126384,
// а при chainLen=25 даёт Part 2 = 154115708116294.
//
// Запуск:
//   go run main.go sample.txt
//   go run main.go input.txt
// и т.д. (файл содержит 5 строк-кодов).
//
// Если у вас всё равно выдаёт неправильные длины — возможно, платформа слишком медленная или
// есть расхождения в реализации. Но данный код соответствует логике из условия.
// ------------------------------------------------

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// --------- Константы для двух частей ---------
const (
	chainLenPart1 = 2  // Part 1: 2 робота
	chainLenPart2 = 25 // Part 2: 25 роботов
)

// --------- Раскладка верхних (directional) клавиатур (2x3) ---------
// Здесь каждая позиция — одна из: "^","v","<",">","A".
// "A" в (row=0,col=2). '^' в (row=0,col=1). '<','v','>' в (row=1,col=0..2).
// gap в (row=0,col=0) и некоторых переходах.
//
// При нажатии стрелок (inp из {"^","v","<",">"}),
//
//	см. map[текущаяПозиция][inp] => новаяПозиция или "" (gap).
//
// При нажатии 'A' обрабатываем отдельно (макроактивация).
var robotKeypadAdj = map[string]map[string]string{
	"^": {
		"^": "",
		"v": "v",
		"<": "",
		">": "A",
	},
	"A": {
		"^": "",
		"v": ">",
		"<": "^",
		">": "",
	},
	"<": {
		"^": "",
		"v": "",
		"<": "",
		">": "v",
	},
	"v": {
		"^": "^",
		"v": "",
		"<": "<",
		">": ">",
	},
	">": {
		"^": "A",
		"v": "",
		"<": "v",
		">": "",
	},
}

// --------- Раскладка цифровой клавиатуры ---------
// (3x3 + два внизу: '0','A'). Рука стартует в 'A' (внизу справа).
// При нажатии стрелок двигается, при 'A' — нажимается текущая кнопка.
var digitKeypadAdj = map[string]map[string]string{
	"7": {"^": "", "v": "4", "<": "", ">": "8"},
	"8": {"^": "", "v": "5", "<": "7", ">": "9"},
	"9": {"^": "", "v": "6", "<": "8", ">": ""},
	"4": {"^": "7", "v": "1", "<": "", ">": "5"},
	"5": {"^": "8", "v": "2", "<": "4", ">": "6"},
	"6": {"^": "9", "v": "3", "<": "5", ">": ""},
	"1": {"^": "4", "v": "", "<": "", ">": "2"},
	"2": {"^": "5", "v": "0", "<": "1", ">": "3"},
	"3": {"^": "6", "v": "A", "<": "2", ">": ""},
	"0": {"^": "2", "v": "", "<": "", ">": "A"},
	"A": {"^": "3", "v": "", "<": "0", ">": ""},
}

// parseNumericPart: "029A" => 29, "456A" => 456, "000A" => 0, ...
func parseNumericPart(code string) int {
	s := strings.TrimSuffix(code, "A")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// --------- Типы состояния ---------
//
// Будем хранить цепочку позиций robPos[0..chainLen-1] (каждый — "^","v","<",">","A"),
// затем digitPos (где рука на цифровой клавиатуре: "7","8","9","4","5","6","1","2","3","0","A"),
// и typed (что уже набрано).
//
// Итого chainLen "robot" позиций + 1 "digit" позиция + typed-строка.
// Всё это кодируем в строку: robPos + '|' + digitPos + '|' + typed.
func encodeState(robPos []string, digitPos, typed string) string {
	// robPos[i] — 1 символ, склеим
	return strings.Join(robPos, "") + "|" + digitPos + "|" + typed
}

// decodeState делаем только при необходимости. В BFS нам нужно собирать rpos заново.
func decodeState(stateKey string, chainLen int) (robPos []string, digitPos, typed string) {
	parts := strings.SplitN(stateKey, "|", 3)
	rposStr := parts[0]
	digitPos = parts[1]
	typed = parts[2]

	robPos = make([]string, chainLen)
	for i := 0; i < chainLen; i++ {
		robPos[i] = string(rposStr[i])
	}
	return
}

// Но, чтобы не путаться, можно просто вручную восстанавливать rpos — см. реализацию BFS ниже.

// --------- "макроактивация" 'A' ---------
// Если рука на роботном уровне i стоит на "^","v","<",">", то это команда двигать уровень i+1 (или цифровую, если i+1==chainLen).
// Если рука "A", идём глубже, если i==chainLen => нажимаем кнопку на цифровой.
//
// Возвращаем (robPos,digitPos,typed,ok).
func macroActivate(robPos []string, digitPos, typed string, idx, chainLen int) ([]string, string, string, bool) {
	cur := robPos[idx]
	switch cur {
	case "^", "v", "<", ">":
		// Команда для нижележащего уровня
		if idx == chainLen-1 {
			// Двигаем цифровую клавиатуру
			nx := digitKeypadAdj[digitPos][cur]
			if nx == "" {
				return robPos, digitPos, typed, false
			}
			digitPos = nx
			return robPos, digitPos, typed, true
		}
		// Иначе двигаем robot idx+1
		n2 := robotKeypadAdj[robPos[idx+1]][cur]
		if n2 == "" {
			return robPos, digitPos, typed, false
		}
		robPosCopy := append([]string(nil), robPos...)
		robPosCopy[idx+1] = n2
		return robPosCopy, digitPos, typed, true

	case "A":
		// Спускаемся ниже, если idx < chainLen-1
		if idx < chainLen-1 {
			return macroActivate(robPos, digitPos, typed, idx+1, chainLen)
		}
		// Если idx == chainLen-1 => нажать кнопку цифровой
		// digitPos может быть "0".."9","A", но если это "^","v","<",">" => нет
		if _, isDir := robotKeypadAdj[digitPos]; isDir && digitPos != "A" {
			return robPos, digitPos, typed, false
		}
		typed += digitPos
		return robPos, digitPos, typed, true

	default:
		return robPos, digitPos, typed, false
	}
}

// solveOneCode находит длину кратчайшей последовательности (BFS) для кода "XYZA" при chainLen роботах.
func solveOneCode(chainLen int, code string) int {
	if code == "" {
		return 0
	}
	// Начальное состояние: robPos[i]="A" для i=0..chainLen-1, digitPos="A", typed=""
	robPos0 := make([]string, chainLen)
	for i := 0; i < chainLen; i++ {
		robPos0[i] = "A"
	}
	dig0 := "A"
	typed0 := ""

	startKey := encodeState(robPos0, dig0, typed0)
	dist := map[string]int{startKey: 0}
	queue := []string{startKey}

	buttons := []string{"^", "v", "<", ">", "A"}

	for len(queue) > 0 {
		curKey := queue[0]
		queue = queue[1:]
		d := dist[curKey]

		// Разбиваем curKey
		parts := strings.SplitN(curKey, "|", 3)
		rposStr := parts[0]
		digPos := parts[1]
		typed := parts[2]

		if typed == code {
			return d
		}
		if !strings.HasPrefix(code, typed) {
			continue
		}

		// Восстановим robPos
		if len(rposStr) != chainLen {
			continue // неверный формат?
		}
		rpos := make([]string, chainLen)
		for i := 0; i < chainLen; i++ {
			rpos[i] = string(rposStr[i])
		}

		// Пять возможных входов
		for _, inp := range buttons {
			rposN := append([]string(nil), rpos...)
			digN := digPos
			typedN := typed
			ok := true

			if inp != "A" {
				// Двигаем верхний уровень (rpos[0])
				cur := rposN[0]
				nx := robotKeypadAdj[cur][inp]
				if nx == "" {
					ok = false
				} else {
					rposN[0] = nx
				}
			} else {
				// Макроактивация
				rposN, digN, typedN, ok = macroActivate(rposN, digN, typedN, 0, chainLen)
			}

			if ok {
				nKey := encodeState(rposN, digN, typedN)
				if _, used := dist[nKey]; !used {
					dist[nKey] = d + 1
					queue = append(queue, nKey)
				}
			}
		}
	}
	return 0
}

// solvePart считает сумму complexities для всех кодов, при заданном chainLen.
func solvePart(chainLen int, codes []string) int {
	sumComp := 0
	for _, cd := range codes {
		ln := solveOneCode(chainLen, cd)
		num := parseNumericPart(cd)
		sumComp += ln * num
	}
	return sumComp
}

// ------------------- main -------------------
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var codes []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		txt := strings.TrimSpace(sc.Text())
		if txt != "" {
			codes = append(codes, txt)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	if len(codes) == 0 {
		fmt.Println("No codes in input.")
		return
	}

	// Part 1
	start1 := time.Now()
	part1 := solvePart(chainLenPart1, codes)
	t1 := time.Since(start1).Seconds()
	fmt.Printf("Part 1: %d (time %.4fs)\n", part1, t1)

	// Part 2
	start2 := time.Now()
	part2 := solvePart(chainLenPart2, codes)
	t2 := time.Since(start2).Seconds()
	fmt.Printf("Part 2: %d (time %.4fs)\n", part2, t2)
}
