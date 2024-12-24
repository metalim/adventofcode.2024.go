package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zeebo/xxh3"
)

var Verbose bool

func Verbosef(format string, a ...any) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

// Sanitize the repository by removing all the tasks and input files to make little lizard happy
func main() {
	flag.BoolVar(&Verbose, "v", false, "Verbose output")
	flag.Parse()
	sanitizeInputsAndTasks()
}

var actions = []string{
	"Angry Alligator annihilated",
	"Bold Bear buried",
	"Clever Chameleon cleared",
	"Daring Dragonfly deleted",
	"Eager Eagle erased",
	"Fierce Falcon finished",
	"Gentle Giraffe gutted",
	"Happy Hedgehog hid",
	"Inquisitive Iguana invalidated",
	"Jumpy Jackal junked",
	"Keen Kangaroo kicked out",
	"Lazy Lizard lost",
	"Mellow Monkey moved away",
	"Naughty Newt nixed",
	"Obnoxious Owl obliterated",
	"Proud Panther purged",
	"Quick Quokka quashed",
	"Rowdy Raccoon removed",
	"Sneaky Snake shredded",
	"Tiny Turtle trashed",
	"Unique Unicorn undid",
	"Valiant Vulture vanished",
	"Witty Wolf wiped",
	"Xenophobic Xerus x-ed out",
	"Young Yak yielded nothing",
	"Zany Zebra zeroed out",
}

func sanitizeInputsAndTasks() {
	var inputsSanitized int
	var tasksSanitized int
	// Walk all txt files and remove inputs and task descriptions
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if strings.HasSuffix(path, ".git") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".txt") {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), "input") {
			if sanitizeInput(path) {
				inputsSanitized++
			}
		} else {
			if sanitizeTask(path) {
				tasksSanitized++
			}
		}
		// if inputsSanitized > 1 || tasksSanitized > 1 {
		// 	return filepath.SkipAll
		// }
		return nil
	})
	fmt.Printf("Sanitized %d inputs and %d tasks\n", inputsSanitized, tasksSanitized)
}

const fmtReplace = "<%s the content, and left this: %x>\n"

var reSanitized = regexp.MustCompile(`^<[\w\s\-]+ the content, and left this: \w+>`)
var reInputAOC = regexp.MustCompile(`^input\d?\.txt$`)

func sanitizeInput(path string) bool {
	Verbosef("Checking input %s\n", path)
	// ignore custom inputs, like input_123.txt
	name := filepath.Base(path)
	if !reInputAOC.MatchString(name) {
		Verbosef("Ignoring custom input %s\n", path)
		return false
	}

	content, err := os.ReadFile(path)
	catch(err)

	if reSanitized.Match(content) {
		Verbosef("Input already sanitized in %s\n", path)
		return false
	}

	fmt.Printf("Sanitizing input %s\n", path)
	action, hash := getActionHash(content)
	replace := fmt.Sprintf(fmtReplace, action, hash)
	err = os.WriteFile(path, []byte(replace), 0644)
	catch(err)
	return true
}

var reTask = regexp.MustCompile(`(?s)(--- Day \d+: [^\n]*---\s*)\n([^\n]*)\n(.*)`)

const fmtDayReplace = "$1\n" + fmtReplace

func sanitizeTask(path string) bool {
	Verbosef("Checking %s\n", path)
	content, err := os.ReadFile(path)
	catch(err)

	m := reTask.FindSubmatch(content)
	if m == nil {
		Verbosef("No task found in %s\n", path)
		return false
	}

	task := m[0]
	// day := m[1]
	trigger := m[2]
	// theRest := m[3]

	if reSanitized.Match(trigger) {
		Verbosef("Task already sanitized in %s\n", path)
		return false
	}

	fmt.Printf("Sanitizing task %s\n", path)
	action, hash := getActionHash(task)
	replace := fmt.Sprintf(fmtDayReplace, action, hash)
	content = reTask.ReplaceAll(content, []byte(replace))
	catch(os.WriteFile(path, content, 0644))
	return true
}

func getActionHash(content []byte) (string, uint64) {
	content = bytes.TrimSpace(content)
	hash := xxh3.Hash(content)
	i := hash % uint64(len(actions))
	return actions[i], hash
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func Fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
