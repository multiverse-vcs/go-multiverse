package core

import (
	"testing"

	fsutil "github.com/go-git/go-billy/v5/util"
)

func TestIgnoreDefault(t *testing.T) {
	mock := NewMockContext()

	IgnoreRules = []string{"foo"}

	rules, err := mock.Ignore()
	if err != nil {
		t.Fatalf("failed to load ignore rules")
	}

	if len(rules) != 1 {
		t.Fatalf("unexpected ignore rules")
	}

	if rules[0] != "foo" {
		t.Errorf("unexpected ignore rules")
	}
}

func TestIgnoreFile(t *testing.T) {
	mock := NewMockContext()

	IgnoreRules = []string{"foo"}

	file := mock.Fs.Join(mock.Fs.Root(), IgnoreFile)
	if err := fsutil.WriteFile(mock.Fs, file, []byte("bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	rules, err := mock.Ignore()
	if err != nil {
		t.Fatalf("failed to load ignore rules")
	}

	if len(rules) != 2 {
		t.Fatalf("unexpected ignore rules")
	}

	if rules[0] != "foo" {
		t.Errorf("unexpected ignore rules")
	}

	if rules[1] != "bar" {
		t.Errorf("unexpected ignore rules")
	}
}
