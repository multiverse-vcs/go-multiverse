package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestInit(t *testing.T) {
	fs = afero.NewMemMapFs()

	args := []string{"multi", "init"}
	if err := app.Run(args); err != nil {
		t.Fatalf("failed to run init command")
	}

	store, err := openStore()
	if err != nil {
		t.Fatalf("failed to open store")
	}

	if _, err := store.ReadConfig(); err != nil {
		t.Errorf("failed to read config %s", err)
	}

	if _, err := store.ReadKey(); err != nil {
		t.Error("failed to read key")
	}
}
