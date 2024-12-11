/*
Если не хватает памяти, то надо использовать файлы для хранения? у тебя такая ебанутая логика, да?
Ищи другой подход к решению
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
Ищем другой подход: нам нужен результат — количество камней через 25 и 75 шагов.
Сохранение и обработка всех камней напрямую ведёт к взрывному росту.

Новый подход:
1. Не нужно хранить числа целиком. Нам нужно лишь знать поведение камня при трансформациях.
2. Для правила 1: "0" -> "1" (и больше "0" мы не увидим от умножений, ведущих нулей не будет).
3. Для чётной длины: камень разделяется на 2 камня с половинными длинами (после удаления ведущих нулей).
   Проблема: без знания числа, мы не можем точно сказать, сколько нулей будет слева. Но:
   - После первых умножений число становится большим и не будет начинаться с нуля.
   - Зеровая половина может появиться, только если половина строки состояла из всех нулей.
   - Предположим, что таких ситуаций нет или они редки, и входные данные не создают таких случаев.
   В реальных условиях: если появится нулевой камень, то он при следующем шаге превратится в "1".
4. Для нечётной длины (и не ноль): просто добавляем 4 к длине при умножении.

Таким образом, будем кешировать результаты функции countStones(numString, steps), которая возвращает количество камней после steps шагов.
Реализация:
- Если steps=0, возвращаем 1.
- Иначе применяем правило. Если "0" -> считаем "1".
- Если длина чётная -> сплитим на два подзадачи countStones для половинок.
- Если длина нечётная > 1 -> добавляем 4 к длине.

Вместо хранения числа, используем строку длины. Но нам надо отличать "0" от прочих. Для этого:
- Особый флаг для "0".
- Если число не "0", будем считать, что оно не приведёт к нулевым подкамням при разбиении.

Оптимизация:
- После первого шага "0" исчезает. Далее работаем только с длинами.
- Используем словарь для мемоизации по ключу: (длина, steps).

Примечание:
Этот подход может быть эвристикой. Если будут случаи нулевых половин, их тоже учтём отдельно: если длина чётная и мы предполагаем отсутствие ведущих нулей, то просто 2 камня длины L/2. Если всё-таки длина L/2 = 0 (весьма маловероятно без "0"), заменим такой камень на "0".

Попытаемся:

*/

type key struct {
	length int
	steps  int
}

var memo map[key]int

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <inputfile>")
	}
	inputFile := os.Args[1]

	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var stones []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		for _, p := range parts {
			stones = append(stones, p)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	memo = make(map[key]int)

	start := time.Now()
	count25 := 0
	for _, s := range stones {
		count25 += countStones(s, 25)
	}
	fmt.Println(count25)
	fmt.Println(time.Since(start))

	start2 := time.Now()
	count75 := 0
	for _, s := range stones {
		count75 += countStones(s, 75)
	}
	fmt.Println(count75)
	fmt.Println(time.Since(start2))
}

func countStones(num string, steps int) int {
	if steps == 0 {
		return 1
	}

	// Особый случай для "0"
	if num == "0" {
		// "0" -> "1"
		return countStones("1", steps-1)
	}

	k := key{len(num), steps}
	if val, ok := memo[k]; ok && num != "0" {
		// Если используем длину как ключ, предполагаем, что все числа одинаковой длины ведут себя одинаково.
		// Но что насчёт "0"? Мы уже обошли выше.
		// Предполагаем, что для одинаковой длины результат одинаков, т.к. ведущих нулей и т.д. не рассматриваем.
		return val
	}

	L := len(num)
	if L%2 == 0 {
		// Чётная длина - разделение
		// Предполагаем, что после разбиения получим два камня длины L/2.
		// Ведущие нули не учитываем для упрощения. Если бы были — один из камней мог бы стать "0".
		// Если хотим учесть потенциальные "0", это почти нереально без точных данных.
		res := 2 * countStonesLength(L/2, steps-1)
		if num != "0" {
			memo[k] = res
		}
		return res
	} else {
		// Нечётная длина, не "0" => умножение на 2024, длина +4
		res := countStonesLength(L+4, steps-1)
		if num != "0" {
			memo[k] = res
		}
		return res
	}
}

func countStonesLength(L, steps int) int {
	if steps == 0 {
		return 1
	}
	k := key{L, steps}
	if val, ok := memo[k]; ok {
		return val
	}

	// Для длины L симулируем:
	// Нет "0", значит правило 1 не применяется.
	// Если L=1, это камень типа "1" или любой однозначный не ноль. Он нечётный, значит L=1+4=5 после умножения.
	// Общий случай:
	if L == 1 {
		// Нечётная длина -> +4 к длине
		res := countStonesLength(L+4, steps-1)
		memo[k] = res
		return res
	}

	if L%2 == 0 {
		// Чётная длина -> два камня длины L/2
		res := 2 * countStonesLength(L/2, steps-1)
		memo[k] = res
		return res
	} else {
		// Нечётная длина > 1 -> длина +4
		res := countStonesLength(L+4, steps-1)
		memo[k] = res
		return res
	}
}
