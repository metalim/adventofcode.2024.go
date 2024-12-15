package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

var noColor = color.New(color.FgWhite)
var box = color.New(color.FgYellow)
var colors = map[rune]*color.Color{
	'#': color.New(color.FgRed),
	'@': color.New(color.FgHiGreen),
	'O': box,
	'[': box,
	']': box,
	'.': noColor,
}

func init() {
	noColor.DisableColor()
}

const ClearScreen = "\033[H\033[2J"
const RedrawScreen = "\033[H"
const HideCursor = "\033[?25l"
const ShowCursor = "\033[?25h"

func initPrint() {
	if Print {
		fmt.Print(ClearScreen)
	}
}

func printGrid(room map[Point]rune, W, H, i int, instructions string) {
	if !Print {
		return
	}
	var buf bytes.Buffer
	buf.WriteString(RedrawScreen)
	buf.WriteString(HideCursor)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			c := room[Point{X: x, Y: y}]
			col, ok := colors[c]
			if !ok {
				col = noColor
			}
			buf.WriteString(col.Sprintf("%c", c))
		}
		buf.WriteString("\n")
	}
	next := '*'
	if i < len(instructions)-1 {
		next = rune(instructions[i+1])
	}
	fmt.Fprintf(&buf, "%d/%d, instruction: %c, next: %c", i+1, len(instructions), instructions[i], next)
	buf.WriteString("\n")
	buf.WriteString(ShowCursor)
	os.Stdout.Write(buf.Bytes())
	os.Stdout.Sync()
	time.Sleep(Delay)
}
