package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
)

func FileSizeFormat(size int) string {
	if size == 0 {
		return "empty"
	}
	return strconv.Itoa(size) + "b"
}

func MakeFirstPrefix(curIdx, size int, prefix string) string {
	if curIdx == size-1 {
		return prefix + "\t"
	}
	return prefix + "│\t"
}

func MakeSecondPrefix(curIdx, size int) string {
	if curIdx == size-1 {
		return "└───"
	}
	return "├───"
}

func dirTreeFull(output io.Writer, path string, printFiles bool, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if printFiles {
		filesCount := len(files)
		for idx, file := range files {
			if file.IsDir() {
				_, err := fmt.Fprintf(output, "%s%s%s\n", prefix, MakeSecondPrefix(idx, filesCount), file.Name())
				if err != nil {
					return err
				}
				newPath := fmt.Sprintf("%s/%s", path, file.Name())
				err = dirTreeFull(output, newPath, printFiles, MakeFirstPrefix(idx, filesCount, prefix))
				if err != nil {
					return err
				}
			} else {
				fileInfo, _ := file.Info()
				fileSize := int(fileInfo.Size())
				_, err := fmt.Fprintf(output, "%s%s%s (%s)\n", prefix, MakeSecondPrefix(idx, filesCount), file.Name(), FileSizeFormat(fileSize))
				if err != nil {
					return err
				}
			}
		}
	} else {
		dirFiles := make([]fs.DirEntry, 0, len(files))
		for _, file := range files {
			if file.IsDir() {
				dirFiles = append(dirFiles, file)
			}
		}
		filesCount := len(dirFiles)
		for idx, file := range dirFiles {
			if file.IsDir() {
				_, err := fmt.Fprintf(output, "%s%s%s\n", prefix, MakeSecondPrefix(idx, filesCount), file.Name())
				if err != nil {
					return err
				}
				newPath := fmt.Sprintf("%s/%s", path, file.Name())
				err = dirTreeFull(output, newPath, printFiles, MakeFirstPrefix(idx, filesCount, prefix))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	return dirTreeFull(output, path, printFiles, "")
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
