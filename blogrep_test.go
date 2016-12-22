package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const TestdataDir = "/testdata/"

func TestSearchText(t *testing.T) {
	buf := &bytes.Buffer{}
	patterns := []string{"i18n", "revert"}

	cwd, _ := os.Getwd()
	if err := filepath.Walk(cwd+TestdataDir, readAndGrep(patterns, buf)); err != nil {
		t.Fatalf("got error %s", err.Error())
	}

	outputString := buf.String()

	if !strings.Contains(outputString, "i18n") {
		t.Errorf("Should contain %s in stdout. stdout %s", "i18n", outputString)
	}

	if !strings.Contains(outputString, "revert") {
		t.Errorf("Should contain %s in stdout. stdout %s", "revert", outputString)
	}
}

func TestSearchTopText(t *testing.T) {
	buf := &bytes.Buffer{}
	patterns := []string{"流し読み"}

	cwd, _ := os.Getwd()
	if err := filepath.Walk(cwd+TestdataDir, readAndGrep(patterns, buf)); err != nil {
		t.Fatalf("got error %s", err.Error())
	}

	outputString := buf.String()

	if !strings.Contains(outputString, "流し読み") {
		t.Errorf("Should contain %s in stdout. stdout %s", "流し読み", outputString)
	}
}
