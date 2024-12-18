package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const SuffGo = ".go"
const SuffTxt = ".txt"

func main() {
	filepath.WalkDir(".", func(pathGo string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasPrefix(filepath.Ext(pathGo), SuffGo) {
			return nil
		}
		if split(pathGo) {
			fmt.Printf("split: %s\n", pathGo)
			// os.Exit(0)
		}
		return nil
	})
}

var reV = regexp.MustCompile(`^v\d+`)

func split(pathGo string) bool {
	dir := filepath.Dir(pathGo)
	nameFull := strings.TrimSuffix(filepath.Base(pathGo), filepath.Ext(pathGo))
	m := reV.FindStringSubmatchIndex(nameFull)
	if len(m) == 0 {
		return false
	}
	nameShort := nameFull[m[0]:m[1]]
	pathTxt := filepath.Join(dir, nameShort+SuffTxt)
	pathTxtFull := filepath.Join(dir, nameFull+SuffTxt)
	if pathTxtFull != pathTxt {
		_, err := os.Stat(pathTxtFull)
		if err == nil {
			err = os.Rename(pathTxtFull, pathTxt)
			catch(err)
			return true
		}
	}

	content, err := os.ReadFile(pathGo)
	catch(err)
	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || lines[0] != "/*" {
		return false
	}
	endComment := slices.Index(lines, "*/")
	if endComment == -1 {
		return false
	}
	pkgMain := slices.Index(lines, "package main")
	if pkgMain == -1 {
		return false
	}

	_, err = os.Stat(pathTxt)
	if err == nil {
		fmt.Printf("already split: %s\n", pathTxt)
		return false
	}
	lines[endComment] = ""
	for lines[endComment-1] == "" {
		endComment--
	}
	comments := lines[1 : endComment+1]
	err = os.WriteFile(pathTxt, []byte(strings.Join(comments, "\n")), 0644)
	catch(err)

	program := lines[pkgMain:]
	err = os.WriteFile(pathGo, []byte(strings.Join(program, "\n")), 0644)
	catch(err)
	return true
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
