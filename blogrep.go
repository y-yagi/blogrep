package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

func readAndGrep(pattern string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		var articles []string

		if info.IsDir() {
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}

		defer in.Close()

		data, err := ioutil.ReadAll(in)
		if err != nil {
			return err
		}
		articles = strings.Split(string(data), "***")

		for _, article := range articles {
			if strings.Contains(article, pattern) {
				fmt.Fprintln(os.Stdout, article)
			}
		}

		return nil
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	filepath.Walk("testdata", readAndGrep(args[0]))
}
