/*
Напиши код на Go для решения следующей задачи. Входные данные в файле указываемом аргументом в командной строке.
Выведи время решения каждой части.

https://adventofcode.com/2024/day/6

--- Day 6: Guard Gallivant ---
The Historians use their fancy device again, this time to whisk you all away to the North Pole prototype suit manufacturing lab... in the year 1518! It turns out that having direct access to history is very convenient for a group of historians.

You still have to be careful of time paradoxes, and so it will be important to avoid anyone from 1518 while The Historians search for the Chief. Unfortunately, a single guard is patrolling this part of the lab.

Maybe you can work out where the guard will go ahead of time so that The Historians can search safely?

You start by making a map (your puzzle input) of the situation. For example:

....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...
The map shows the current position of the guard with ^ (to indicate the guard is currently facing up from the perspective of the map). Any obstructions - crates, desks, alchemical reactors, etc. - are shown as #.

Lab guards in 1518 follow a very strict patrol protocol which involves repeatedly following these steps:

If there is something directly in front of you, turn right 90 degrees.
Otherwise, take a step forward.
Following the above protocol, the guard moves up several times until she reaches an obstacle (in this case, a pile of failed suit prototypes):

....#.....
....^....#
..........
..#.......
.......#..
..........
.#........
........#.
#.........
......#...
Because there is now an obstacle in front of the guard, she turns right before continuing straight in her new facing direction:

....#.....
........>#
..........
..#.......
.......#..
..........
.#........
........#.
#.........
......#...
Reaching another obstacle (a spool of several very long polymers), she turns right again and continues downward:

....#.....
.........#
..........
..#.......
.......#..
..........
.#......v.
........#.
#.........
......#...
This process continues for a while, but the guard eventually leaves the mapped area (after walking past a tank of universal solvent):

....#.....
.........#
..........
..#.......
.......#..
..........
.#........
........#.
#.........
......#v..
By predicting the guard's route, you can determine which specific positions in the lab will be in the patrol path. Including the guard's starting position, the positions visited by the guard before leaving the area are marked with an X:

....#.....
....XXXXX#
....X...X.
..#.X...X.
..XXXXX#X.
..X.X.X.X.
.#XXXXXXX.
.XXXXXXX#.
#XXXXXXX..
......#X..
In this example, the guard will visit 41 distinct positions on your map.

Predict the path of the guard. How many distinct positions will the guard visit before leaving the mapped area?


--- Part Two ---
While The Historians begin working around the guard's patrol route, you borrow their fancy device and step outside the lab. From the safety of a supply closet, you time travel through the last few months and record the nightly status of the lab's guard post on the walls of the closet.

Returning after what seems like only a few seconds to The Historians, they explain that the guard's patrol area is simply too large for them to safely search the lab without getting caught.

Fortunately, they are pretty sure that adding a single new obstruction won't cause a time paradox. They'd like to place the new obstruction in such a way that the guard will get stuck in a loop, making the rest of the lab safe to search.

To have the lowest chance of creating a time paradox, The Historians would like to know all of the possible positions for such an obstruction. The new obstruction can't be placed at the guard's starting position - the guard is there right now and would notice.

In the above example, there are only 6 different positions where a new obstruction would cause the guard to get stuck in a loop. The diagrams of these six situations use O to mark the new obstruction, | to show a position where the guard moves up/down, - to show a position where the guard moves left/right, and + to show a position where the guard moves both up/down and left/right.

Option one, put a printing press next to the guard's starting position:

....#.....
....+---+#
....|...|.
..#.|...|.
....|..#|.
....|...|.
.#.O^---+.
........#.
#.........
......#...
Option two, put a stack of failed suit prototypes in the bottom right quadrant of the mapped area:


....#.....
....+---+#
....|...|.
..#.|...|.
..+-+-+#|.
..|.|.|.|.
.#+-^-+-+.
......O.#.
#.........
......#...
Option three, put a crate of chimney-squeeze prototype fabric next to the standing desk in the bottom right quadrant:

....#.....
....+---+#
....|...|.
..#.|...|.
..+-+-+#|.
..|.|.|.|.
.#+-^-+-+.
.+----+O#.
#+----+...
......#...
Option four, put an alchemical retroencabulator near the bottom left corner:

....#.....
....+---+#
....|...|.
..#.|...|.
..+-+-+#|.
..|.|.|.|.
.#+-^-+-+.
..|...|.#.
#O+---+...
......#...
Option five, put the alchemical retroencabulator a bit to the right instead:

....#.....
....+---+#
....|...|.
..#.|...|.
..+-+-+#|.
..|.|.|.|.
.#+-^-+-+.
....|.|.#.
#..O+-+...
......#...
Option six, put a tank of sovereign glue right next to the tank of universal solvent:

....#.....
....+---+#
....|...|.
..#.|...|.
..+-+-+#|.
..|.|.|.|.
.#+-^-+-+.
.+----++#.
#+----++..
......#O..
It doesn't really matter what you choose to use as an obstacle so long as you and The Historians can put it into position without the guard noticing. The important thing is having enough options that you can find one that minimizes time paradoxes, and in this example, there are 6 different positions you could choose.

You need to get the guard stuck in a loop by adding a single new obstruction. How many different positions could you choose for this obstruction?

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type posDir struct {
	r, c, d int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program input.txt")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var grid []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		grid = append(grid, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	rows := len(grid)
	cols := len(grid[0])
	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	startR, startC, startD := -1, -1, -1
FindStart:
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			switch grid[i][j] {
			case '^':
				startR, startC, startD = i, j, 0
				break FindStart
			case '>':
				startR, startC, startD = i, j, 1
				break FindStart
			case 'v':
				startR, startC, startD = i, j, 2
				break FindStart
			case '<':
				startR, startC, startD = i, j, 3
				break FindStart
			}
		}
	}

	original := make([][]byte, rows)
	for i := range grid {
		original[i] = []byte(grid[i])
	}

	// Часть 1
	part1Start := time.Now()
	{
		r, c, d := startR, startC, startD
		visited := make(map[[2]int]bool)
		visited[[2]int{r, c}] = true

		for {
			nr, nc := r+dirs[d][0], c+dirs[d][1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				break
			}
			if original[nr][nc] == '#' {
				d = (d + 1) % 4
			} else {
				r, c = nr, nc
				visited[[2]int{r, c}] = true
			}
		}

		fmt.Println(len(visited)) // Результат части 1
	}
	part1Time := time.Since(part1Start)
	fmt.Println("Part1 time:", part1Time)

	// Часть 2
	part2Start := time.Now()
	simWithObstacle := func(or, oc int) bool {
		original[or][oc] = '#'
		r, c, d := startR, startC, startD
		visitedStates := make(map[posDir]bool)
		visitedStates[posDir{r, c, d}] = true

		for {
			nr, nc := r+dirs[d][0], c+dirs[d][1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				original[or][oc] = '.'
				return false
			}
			if original[nr][nc] == '#' {
				d = (d + 1) % 4
			} else {
				r, c = nr, nc
				st := posDir{r, c, d}
				if visitedStates[st] {
					original[or][oc] = '.'
					return true
				}
				visitedStates[st] = true
			}
		}
	}

	count := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if (i == startR && j == startC) || original[i][j] != '.' {
				continue
			}
			if simWithObstacle(i, j) {
				count++
			}
		}
	}

	fmt.Println(count) // Результат части 2
	part2Time := time.Since(part2Start)
	fmt.Println("Part2 time:", part2Time)
}
