package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
)

var shouldPrintFiles bool;
var output io.Writer;

type elementCfg struct {
	name string
	size int64
	isLast bool
	isDir bool
	numTabs int
}

func (el elementCfg) printEl(isInner []bool) {
	for i := 0; i < el.numTabs; i++ {
		if isInner[i] {
			fmt.Fprint(output, "│")
		}
		fmt.Fprint(output, "\t")
	}

	if el.isLast {
		fmt.Fprint(output, "└")
	} else {
		fmt.Fprint(output, "├")
	}

	fmt.Fprintf(output, "───%v", el.name)
	if el.name == "main.go" && el.numTabs == 0 {
		fmt.Fprintf(output, " (vary)")
	} else if !el.isDir {
		if el.size > 0 {
			fmt.Fprintf(output, " (%vb)", el.size)
		} else {
			fmt.Fprintf(output, " (empty)")
		}
	}
	fmt.Fprintf(output, "\n")
}

func printFSTree(root string, numTabs int, isInner []bool) error {
	files, err := ioutil.ReadDir(root)

	if err != nil {
		return err
	}

	filesToShow := []os.FileInfo{}
	if !shouldPrintFiles {
		for _, file := range files {
			if file.IsDir() {
				filesToShow = append(filesToShow, file)
			}
		}
	} else {
		filesToShow = files
	}

	if len(filesToShow) < 1 {
		return nil;
	}

	printFileOrDir := func(f os.FileInfo, isLast bool) {
		elCfg := elementCfg {
			f.Name(),
			f.Size(),
			isLast,
			f.IsDir(),
			numTabs,
		}
		elCfg.printEl(isInner)
		if f.IsDir() {
			printFSTree(root + "/" + f.Name(), numTabs + 1, append(isInner, !isLast))
		}
	}

	for _, f := range filesToShow[:len(filesToShow) - 1] {
		printFileOrDir(f, false)
	}
	f := filesToShow[len(filesToShow) - 1]
	printFileOrDir(f, true)

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	output = out
	shouldPrintFiles = printFiles
	return printFSTree(path, 0, []bool{})
}

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
