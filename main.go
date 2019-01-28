package main

import (
	"io"
	"os"
	"sort"
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
	f, err := os.Open(path)
	defer func() {
		f.Close()
	}()
	if err != nil {
		return nil, err
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Leave directories or all
	if printFiles == false {
		_filterDirs(&files)
	}

	//Sort slice by file name
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	return files, nil
}

func _dirTree(out io.Writer, path string, printFiles bool, ident string) error {

	files, err := _readFiles(path, printFiles)
	if err != nil {
		return err
	}

	for idx, file := range files {
		out.Write([]byte(ident))

		lastEntry := idx == len(files)-1
		if lastEntry == false {
			out.Write([]byte("├"))
		} else {
			out.Write([]byte("└"))
		}
		out.Write([]byte("───" + file.Name()))

		if printFiles == true {
			if file.IsDir() == false {
				out.Write([]byte(" ("))
				size := file.Size()
				if size > 0 {
					out.Write([]byte(strconv.FormatInt(size, 10) + "b"))
				} else {
					out.Write([]byte("empty"))
				}
				out.Write([]byte(")"))
			}
		}
		out.Write([]byte("\n"))

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
