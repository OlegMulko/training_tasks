package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
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

func dirTree(out io.Writer, path string, printFiles bool) error {

	err := printFilesLevel(out, path, "", printFiles)
	if err != nil {
		return err
	}
	return nil
}

func printFilesLevel(out io.Writer, path string, prefixString string, printFiles bool) error {

	var filesSize int
	var subPrefixString string

	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	filesInfo, err := file.Readdir(0)
	if err != nil {
		return err
	}

	sort.SliceStable(filesInfo, func(i, j int) bool {
		return filesInfo[i].Name() < filesInfo[j].Name()
	})

	if printFiles {
		filesSize = len(filesInfo)
	} else {
		filesSize = getCountDirs(filesInfo)
	}

	i := 1

	for _, fileInfo := range filesInfo {
		if !fileInfo.IsDir() && !printFiles {
			continue
		}
		if i == filesSize {
			fmt.Fprintln(out, getLine(prefixString, true, fileInfo))
		} else {
			fmt.Fprintln(out, getLine(prefixString, false, fileInfo))
		}

		if fileInfo.IsDir() {

			path := path + string(os.PathSeparator) + fileInfo.Name()

			if i != filesSize {
				subPrefixString = prefixString + "│	"
			} else {
				subPrefixString = prefixString + "	"
			}

			printFilesLevel(out, path, subPrefixString, printFiles)
		}
		i++
	}
	return nil
}

func getCountDirs(filesInfo []os.FileInfo) int {

	var size int = 0

	if len(filesInfo) == 0 {
		return size
	}

	for _, fileInfo := range filesInfo {
		if !fileInfo.IsDir() {
			continue
		}
		size++
	}
	return size
}

func getLine(prefixString string, last bool, fileInfo os.FileInfo) string {

	printString := prefixString

	if last {
		printString += "└"
	} else {
		printString += "├"
	}

	printString += "───" + fileInfo.Name()

	if !fileInfo.IsDir() {
		sizeString := ""
		if fileInfo.Size() == 0 {
			sizeString = "empty"
		} else {
			sizeString = strconv.Itoa((int)(fileInfo.Size())) + "b"
		}
		printString += " (" + sizeString + ")"
	}
	return printString
}
