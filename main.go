package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func _filterDirs(files *[]os.FileInfo) {
	idx := 0
	for _, file := range *files {
		if file.IsDir() == true {
			(*files)[idx] = file
			idx++
		}
	}
	*files = (*files)[:idx]
}

func _dirTree(out io.Writer, path string, includeFiles bool, ident string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// Leave directories or all
	if includeFiles == false {
		_filterDirs(&files)
	}

	for idx, file := range files {
		lastFile := idx == len(files)-1

		writeText := func(text string) {
			out.Write([]byte(text))
		}

		writeText(ident)
		if lastFile == false {
			writeText("├───")
		} else {
			writeText("└───")
		}
		writeText(file.Name())

		writeSize := func(file os.FileInfo) {
			if file.IsDir() == true {
				return
			}

			writeText(" (")
			size := file.Size()
			if size > 0 {
				writeText(strconv.FormatInt(size, 10) + "b")
			} else {
				writeText("empty")
			}
			writeText(")")
		}

		if includeFiles == true {
			writeSize(file)
		}
		writeText("\n")

		//prepare for the next step
		if file.IsDir() == true {
			nextIdent := ident
			if lastFile == false {
				nextIdent += "│"
			}
			nextIdent += "\t"
			_dirTree(out, path+"/"+file.Name(), includeFiles, nextIdent)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	if err := _dirTree(out, path, printFiles, ""); err != nil {
		return err
	}
	return nil
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
