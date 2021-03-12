package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	tree := Tree{path: path}
	treeResult, err := tree.String(printFiles)
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(treeResult))

	return err
}

type Tree struct {
	path   string
	parent *Tree
	isLast bool
}

func (t *Tree) String(withFiles bool) (string, error) {
	files, err := ioutil.ReadDir(t.path)
	if err != nil {
		return "", err
	}

	if !withFiles {
		files = filterOnLyDirs(files)
	}

	prefix := t.getStringPrefix()

	stringBuilder := strings.Builder{}
	for i := 0; i < len(files); i++ {
		isLastPosition := i == len(files)-1

		node := Node{files[i], isLastPosition}

		stringBuilder.WriteString(prefix)
		stringBuilder.WriteString(node.String())
		stringBuilder.WriteString("\n")

		if files[i].IsDir() {
			tree := Tree{filepath.Join(t.path, files[i].Name()), t, isLastPosition}
			treeResult, err := tree.String(withFiles)
			if err != nil {
				return "", nil
			}

			stringBuilder.WriteString(treeResult)
		}
	}

	return stringBuilder.String(), nil
}

func (t *Tree) getStringPrefix() string {
	if t.parent == nil {
		return ""
	}

	result := t.parent.getStringPrefix()
	if t.isLast {
		return result + "\t"
	}

	return result + "│\t"
}

type Node struct {
	fileInfo os.FileInfo
	isLast   bool
}

func (n *Node) String() string {
	return n.getStringPrefix() + n.getFileInfo()
}

func (n *Node) getStringPrefix() string {
	if n.isLast {
		return "└───"
	}

	return "├───"
}

func (n *Node) getFileInfo() string {
	if n.fileInfo.IsDir() {
		return n.fileInfo.Name()
	}

	size := n.fileInfo.Size()

	if size == 0 {
		return n.fileInfo.Name() + " (empty)"
	}

	return fmt.Sprintf("%s (%db)", n.fileInfo.Name(), size)
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
