package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrew-d/isbinary"
	"github.com/fatih/color"
	"github.com/y-yagi/configure"
)

type config struct {
	Home string `toml:"home"`
}

const cmd = "blogrep"

var (
	warningColor  = color.New(color.FgYellow).SprintFunc()
	filePathColor = color.New(color.FgGreen, color.Bold).SprintFunc()
)

func errorline(s string) {
	os.Stderr.WriteString(s + "\n")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERNS\n", os.Args[0])
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func cmdEdit() error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}

	return configure.Edit(cmd, editor)
}

func containsAllAndColorized(article *string, patterns []string) bool {
	for _, pattern := range patterns {
		if !strings.Contains(strings.ToLower(*article), strings.ToLower(pattern)) {
			return false
		}
		*article = strings.Replace(*article, pattern, warningColor(pattern), -1)
	}
	return true
}

func readAndGrep(patterns []string, writer io.Writer) filepath.WalkFunc {
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

		fileIsBinary, err := isbinary.TestReader(in)
		if err != nil {
			return err
		}
		if fileIsBinary {
			return nil
		}

		in.Seek(0, 0)
		data, err := ioutil.ReadAll(in)
		if err != nil {
			return err
		}

		articles = strings.Split(string(data), "\n***\n\n")

		for _, article := range articles {
			if containsAllAndColorized(&article, patterns) {
				fmt.Fprintln(writer, filePathColor(path))
				fmt.Fprintln(writer, article)
			}
		}

		return nil
	}
}
func init() {
	if !configure.Exist(cmd) {
		var cfg config
		cfg.Home = ""
		configure.Save(cmd, cfg)
	}
}

func main() {
	var edit bool

	flag.BoolVar(&edit, "c", false, "edit config")
	flag.Parse()

	if edit {
		os.Exit(msg(cmdEdit()))
	}

	args := os.Args[1:]

	if len(args) < 1 {
		usage()
		os.Exit(2)
	}

	var cfg config
	err := configure.Load(cmd, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Home) == 0 {
		fmt.Fprintf(os.Stderr, "Please specify home to config file.\n")
		os.Exit(1)
	}

	err = os.Chdir(cfg.Home)
	if err != nil {
		errorline(err.Error())
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
	err = filepath.Walk(cwd, readAndGrep(args, os.Stdout))
	if err != nil {
		errorline(err.Error())
		os.Exit(1)
	}
}
