package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrew-d/isbinary"
	"github.com/fatih/color"
)

func errorline(s string) {
	os.Stderr.WriteString(s + "\n")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

func containsAllAndColorized(article *string, patterns []string) bool {
	warningColor := color.New(color.FgYellow).SprintFunc()
	for _, pattern := range patterns {
		if !strings.Contains(strings.ToLower(*article), strings.ToLower(pattern)) {
			return false
		}
		*article = strings.Replace(*article, pattern, warningColor(pattern), -1)
	}
	return true
}

func readAndGrep(patterns []string) filepath.WalkFunc {
	filePathColor := color.New(color.FgGreen, color.Bold).SprintFunc()
	return func(path string, info os.FileInfo, err error) error {
		var articles []string

		if info.IsDir() {
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}

		fileIsBinary, err := isbinary.TestReader(in)
		if err != nil {
			return err
		}
		if fileIsBinary {
			return nil
		}

		defer in.Close()

		data, err := ioutil.ReadAll(in)
		if err != nil {
			return err
		}
		articles = strings.Split(string(data), "\n***\n\n")

		for _, article := range articles {
			if containsAllAndColorized(&article, patterns) {
				fmt.Fprintln(os.Stdout, filePathColor(path))
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
