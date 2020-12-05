package core

import (
	"testing"

	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/spf13/afero"
)

func TestIgnoreDefault(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	IgnoreRules = []string{"foo"}
	rules, err := Ignore(store)
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
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("failed to create storage")
	}

	if err := afero.WriteFile(store.Cwd, IgnoreFile, []byte("bar"), 0644); err != nil {
		t.Fatalf("failed to write file")
	}

	IgnoreRules = []string{"foo"}
	rules, err := Ignore(store)
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
