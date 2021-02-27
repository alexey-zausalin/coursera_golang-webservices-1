package main

import (
	"fmt"
	_ "fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	_ "strings"
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
	return dirPrintTree(out, path, "", printFiles)
}

func dirPrintTree(out io.Writer, path string, prefix string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		files = filterOnLyDirs(files)
	}

	for i := 0; i < len(files); i++ {
		isLastPosition := i == len(files)-1

		if isLastPosition {
			_, err = out.Write([]byte(prefix + getLastPositionPrefix() + getFileFormattedData(files[i]) + "\n"))
		} else {
			_, err = out.Write([]byte(prefix + getMiddlePositionPrefix() + getFileFormattedData(files[i]) + "\n"))
		}

		if err != nil {
			return err
		}

		if files[i].IsDir() {
			subLevelPrefix := prefix
			if isLastPosition {
				subLevelPrefix += getTopLevelLastPositionPrefix()
			} else {
				subLevelPrefix += getTopLevelMiddlePositionPrefix()
			}

			err = dirPrintTree(out, path+string(filepath.Separator)+files[i].Name(), subLevelPrefix, printFiles)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func filterOnLyDirs(files []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file)
		}
	}
	return result
}

func getMiddlePositionPrefix() string {
	return "├───"
}

func getLastPositionPrefix() string {
	return "└───"
}

func getTopLevelMiddlePositionPrefix() string {
	return "│\t"
}

func getTopLevelLastPositionPrefix() string {
	return "\t"
}

func getFileFormattedData(fileInfo os.FileInfo) string {
	if fileInfo.IsDir() {
		return fileInfo.Name()
	}

	size := fileInfo.Size()

	if size == 0 {
		return fileInfo.Name() + " (empty)"
	}

	return fmt.Sprintf("%s (%db)", fileInfo.Name(), size)
}
