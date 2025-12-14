package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tinylib/msgp/gen"
)

// When stuff's going wrong, you'll be glad this is here!
const showGeneratedFile = false

func rename(filename, suffix, newSuffix string) string {
	return strings.TrimSuffix(filename, suffix) + newSuffix
}

// generate - returns filename, filenameGen, error
func generate(t *testing.T, content string) (string, string, error) {
	t.Helper()

	tempDir := t.TempDir()

	filename := filepath.Join(tempDir, "main.go")

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		return "", "", err
	}
	defer fd.Close()

	if _, err := fd.WriteString(content); err != nil {
		return "", "", err
	}

	mode := gen.Encode | gen.Decode | gen.Size | gen.Marshal | gen.Unmarshal | gen.Test
	if err := Run(filename, mode, false); err != nil {
		return "", "", err
	}

	filenameGen := rename(filename, ".go", "_gen.go")

	if showGeneratedFile {
		content, err := os.ReadFile(filenameGen)
		if err != nil {
			return "", "", err
		}
		t.Logf("generated %s content:\n%s", filenameGen, content)
	}

	return filename, filenameGen, nil
}

func goExec(t *testing.T, args ...string) {
	t.Helper()

	output, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v, output:\n%s", err, output)
	}
}

func runTest(t *testing.T, content string) {
	t.Helper()

	filename, filenameGen, err := generate(t, content)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	filenameGenTest := rename(filenameGen, "_gen.go", "_gen_test.go")

	goExec(t, "run", filename, filenameGen)
	goExec(t, "test", filename, filenameGen, filenameGenTest)
}
