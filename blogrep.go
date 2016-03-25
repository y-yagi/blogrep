package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func errorline(s string) {
	os.Stderr.WriteString(s + "\n")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

func containsAll(article string, patterns []string) bool {
	for _, pattern := range patterns {
		if !strings.Contains(article, pattern) {
			return false
		}
	}
	return true
}

func readAndGrep(patterns []string) filepath.WalkFunc {
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
			if containsAll(article, patterns) {
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

	cwd, _ := os.Getwd()
	err := filepath.Walk(cwd, readAndGrep(args))
	if err != nil {
		errorline(err.Error())
		os.Exit(1)
	}
}
