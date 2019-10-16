package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//FindGoFilesToProcess - starting from root, look for .go files
func FindGoFilesToProcess(root string, verbose bool) []string {
	var files []string
	var folders []string
	folders = append(folders, root)
	for len(folders) > 0 {
		folder := folders[0]
		folders = folders[1:]

		if verbose {
			fmt.Printf("processing folder %s\n", folder)
		}

		fileTraverseErr := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if path == folder {
				if verbose {
					fmt.Printf("ignoring current folder %s\n", path)
				}
			} else if strings.HasSuffix(path, "_test") {
				if verbose {
					fmt.Printf("ignoring test folder %s\n", path)
				}
			} else if strings.Contains(path, ".git") {
				if verbose {
					fmt.Printf("ignoring git folder %s\n", path)
				}
			} else if strings.Contains(path, ".vscode") {
				if verbose {
					fmt.Printf("ignoring vscode folder %s\n", path)
				}
			} else if info.IsDir() {
				if verbose {
					fmt.Printf("ignoring folder %s\n", path)
				}
				// folders = append(folders, path)
				// if verbose {
				// 	fmt.Printf("adding folder %s to the folder list\n", path)
				// }
			} else if !strings.HasSuffix(path, "go") {
				if verbose {
					fmt.Printf("ignoring non-go file %s\n", path)
				}
			} else if strings.HasSuffix(path, "_test.go") {
				if verbose {
					fmt.Printf("ignoring test file %s\n", path)
				}
			} else {
				files = append(files, path)
				if verbose {
					fmt.Printf("adding file %s to the file list\n", path)
				}
			}
			return nil
		})

		if fileTraverseErr != nil {
			panic(fileTraverseErr)
		}
	}

	return files
}
