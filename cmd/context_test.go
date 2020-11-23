package cmd

import (
	"context"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
)

func TestInitContext(t *testing.T) {
	fs := memfs.New()
	if err := InitContext(fs, context.TODO()); err != nil {
		t.Fatalf("failed to init context")
	}

	if _, err := fs.Lstat(DotDir); err != nil {
		t.Fatalf("expected dot dir to exist")
	}

	path := fs.Join(DotDir, ConfigFile)
	if _, err := fs.Lstat(path); err != nil {
		t.Errorf("expected config file to exist")
	}
}
