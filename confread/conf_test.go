package confread

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type readTestConfig struct {
	Value string `json:"value" toml:"value" yaml:"value"`
}

func TestReadUsesFirstExistingCandidateDirectory(t *testing.T) {
	fileName := "app.yaml"
	root := t.TempDir()
	firstDir := filepath.Join(root, "first")
	secondDir := filepath.Join(root, "second")
	writeReadTestConfig(t, firstDir, fileName, "value: first\n")
	writeReadTestConfig(t, secondDir, fileName, "value: second\n")

	var got readTestConfig
	if err := Read(fileName, &got, firstDir, secondDir); err != nil {
		t.Fatalf("Read(%q, positions %q) error = %v, want nil", fileName, []string{firstDir, secondDir}, err)
	}

	if got.Value != "first" {
		t.Errorf("Read(%q, positions %q) value = %q, want %q", fileName, []string{firstDir, secondDir}, got.Value, "first")
	}
}

func TestReadUsesLaterCandidateDirectoryWhenEarlierIsMissing(t *testing.T) {
	fileName := "app.yaml"
	root := t.TempDir()
	missingDir := filepath.Join(root, "missing")
	secondDir := filepath.Join(root, "second")
	writeReadTestConfig(t, secondDir, fileName, "value: second\n")

	var got readTestConfig
	if err := Read(fileName, &got, missingDir, secondDir); err != nil {
		t.Fatalf("Read(%q, positions %q) error = %v, want nil", fileName, []string{missingDir, secondDir}, err)
	}

	if got.Value != "second" {
		t.Errorf("Read(%q, positions %q) value = %q, want %q", fileName, []string{missingDir, secondDir}, got.Value, "second")
	}
}

func TestReadReturnsErrNotFoundWhenAllCandidatesAreMissing(t *testing.T) {
	fileName := "app.yaml"
	root := t.TempDir()
	firstDir := filepath.Join(root, "first")
	secondDir := filepath.Join(root, "second")

	var got readTestConfig
	err := Read(fileName, &got, firstDir, secondDir)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Read(%q, positions %q) error = %v, want errors.Is ErrNotFound", fileName, []string{firstDir, secondDir}, err)
	}
}

func writeReadTestConfig(t *testing.T, dir, fileName, contents string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v, want nil", dir, err)
	}
	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v, want nil", path, err)
	}
}
