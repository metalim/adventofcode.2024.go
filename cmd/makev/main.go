package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	version := 1
	name := "v1"
	for {
		if _, err := os.Stat(filepath.Join("o1", name)); os.IsNotExist(err) {
			break
		}
		version++
		name = "v" + strconv.Itoa(version)
	}
	folder := filepath.Join("o1", name)
	err := os.MkdirAll(folder, 0755)
	catch(err)
	err = os.WriteFile(filepath.Join(folder, name+".go"), nil, 0644)
	catch(err)
	if version > 1 {
		err = os.WriteFile(filepath.Join(folder, name+".txt"), nil, 0644)
		catch(err)
		return
	}
	prompt, err := os.ReadFile("../o1/prompt.txt")
	catch(err)
	prompt = append(prompt, "\n"...)
	task, err := os.ReadFile("task.txt")
	catch(err)
	prompt = append(prompt, task...)
	err = os.WriteFile(filepath.Join(folder, name+".txt"), prompt, 0644)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
