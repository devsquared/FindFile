package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	fileName := flag.String("file", "", "name of file")
	pathName := flag.String("path", "", "path to find file")
	logLevel := flag.String("logLevel", "verbose", "amount of logs returned - verbose or simple")

	flag.Parse()
	log.SetFlags(0)

	path := resolveInitialPath(*pathName)

	info, err := os.Lstat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	verbose := true

	if *logLevel != "verbose" {
		verbose = false
	}

	fmt.Println("Attempting to find " + *fileName + " in the parent path " + path + "...")
	pathContaintingFile := findFile(path, info, *fileName, verbose)

	if pathContaintingFile != "" {
		fmt.Println("Found the file in " + pathContaintingFile)
	} else {
		fmt.Println("Did not find file.")
	}
}

func resolveInitialPath(givenPath string) string {
	if givenPath == "" {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		return path
	}

	return givenPath
}

func findFile(currentPath string, info os.FileInfo, fileName string, verbose bool) string {
	if !info.IsDir() {
		return ""
	}

	dir, err := os.Open(currentPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer dir.Close()

	if verbose {
		fmt.Println("Searching for " + fileName + " within the path: " + currentPath)
	}

	if _, err := os.Stat(currentPath + string(os.PathSeparator) + snapshotJARName); !os.IsNotExist(err) {
		// path/to/whatever does exist

		return currentPath
	}

	// continue reading
	files, err := dir.Readdir(-1)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var allPaths []string

	for _, file := range files {
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		foundFile := findFile(currentPath+string(os.PathSeparator)+file.Name(), file, fileName, verbose)
		allPaths = append(allPaths, foundFile)
	}

	pathWithFile := ""
	for _, path := range allPaths {
		if path != "" {
			pathWithFile = path
		}
	}

	return pathWithFile
}
