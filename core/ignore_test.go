package core

import (
	"testing"

	"github.com/spf13/afero"
)

func TestIgnoreDefault(t *testing.T) {
	fs = afero.NewMemMapFs()

	IgnoreRules = []string{"foo"}
	rules, err := Ignore()
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
	fs = afero.NewMemMapFs()

	if err := afero.WriteFile(fs, IgnoreFile, []byte("bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	IgnoreRules = []string{"foo"}
	rules, err := Ignore()
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
