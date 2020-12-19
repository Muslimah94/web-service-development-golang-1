package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {
	return myDirTree(out, path, printFiles, "", false)
}

func myDirTree(out *os.File, path string, printFiles bool, lastPrefix string, isLast bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	filesInfo, err := file.Readdir(-1)
	if err != nil {
		return err
	}
	sort.SliceStable(filesInfo, func(i, j int) bool { return filesInfo[i].Name() < filesInfo[j].Name() })
	var startSymbol, nextStartSymbol, nextPath, size string
	for i, file := range filesInfo {
		if i+1 == len(filesInfo) {
			startSymbol = `└`
			nextStartSymbol = ``
			isLast = true
		} else {
			startSymbol = `├`
			nextStartSymbol = `│`
		}

		if file.IsDir() {
			fmt.Fprintf(out, "%s%s───%s\n", lastPrefix, startSymbol, file.Name())
			nextPath = filepath.Join(path, file.Name())
			nextPrefix := fmt.Sprintf("%s%s\t", lastPrefix, nextStartSymbol)
			myDirTree(out, nextPath, printFiles, nextPrefix, isLast)
		} else if printFiles == true {
			if file.Size() > 0 {
				size = fmt.Sprintf("%db", file.Size())
			} else {
				size = "empty"
			}
			fmt.Fprintf(out, "%s%s───%s (%s)\n", lastPrefix, startSymbol, file.Name(), size)
		}
	}
	return nil

}
