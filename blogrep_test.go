package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const TESTDATA_DIR = "/testdata/"

func TestSearchText(t *testing.T) {
	buf := &bytes.Buffer{}
	patterns := []string{"i18n", "revert"}

	cwd, _ := os.Getwd()
	if err := filepath.Walk(cwd+TESTDATA_DIR, readAndGrep(patterns, buf)); err != nil {
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