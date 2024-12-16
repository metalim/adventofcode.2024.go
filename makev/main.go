package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Usage: makev <name>")
		os.Exit(1)
	}

	name := flag.Arg(0)
	folder := filepath.Join("o1", name)
	err := os.MkdirAll(folder, 0755)
	catch(err)
	err = os.WriteFile(filepath.Join(folder, name+".go"), []byte("/*\n\n*/\n\n"), 0644)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
