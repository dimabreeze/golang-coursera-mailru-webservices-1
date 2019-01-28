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

func _readFiles(path string, printFiles bool) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Leave directories or all
	if printFiles == false {
		_filterDirs(&files)
	}

	return files, nil
}

func _dirTree(out io.Writer, path string, printFiles bool, ident string) error {

	files, err := _readFiles(path, printFiles)
	if err != nil {
		return err
	}

	for idx, file := range files {
		lastEntry := idx == len(files)-1

		write := func(text string) {
			out.Write([]byte(text))
		}

		write(ident)
		if lastEntry == false {
			write("├───")
		} else {
			write("└───")
		}
		write(file.Name())

		if printFiles == true && file.IsDir() == false {
			write(" (")
			size := file.Size()
			if size > 0 {
				write(strconv.FormatInt(size, 10) + "b")
			} else {
				write("empty")
			}
			write(")")
		}
		write("\n")

		//prepare for the next step
		if file.IsDir() == true {
			nextIdent := ident
			if lastEntry == false {
				nextIdent += "│"
			}
			nextIdent += "\t"
			_dirTree(out, path+"/"+file.Name(), printFiles, nextIdent)
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
