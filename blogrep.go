package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

func readAndGrep(path string, info os.FileInfo, err error) error {
	fmt.Println(path)
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	fmt.Println(args)

	filepath.Walk(".", readAndGrep)
}
