package ignore

import (
	"testing"
)

func TestMatch(t *testing.T) {
	ignore := New("test", "*.exe")

	if !ignore.Match("test/foo/bar.exe") {
		t.Error("expected ignore to match")
	}

	if !ignore.Match("foo.exe") {
		t.Error("expected ignore to match")
	}
}

func TestLoad(t *testing.T) {
	ignore, err := Load("testdata")
	if err != nil {
		t.Fatal("failed to load ignore file")
	}

	if !ignore.Match("foo/bar") {
		t.Error("expected ignore to match")
	}

	if !ignore.Match("bar/foo") {
		t.Error("expected ignore to match")
	}

	if ignore.Match("foo.exe") {
		t.Error("expected ignore not to match")
	}
}

func TestMerge(t *testing.T) {
	ignore, err := Load("testdata")
	if err != nil {
		t.Fatal("failed to load ignore file")
	}

	other := New("test", "*.exe")
	merge := ignore.Merge(other)

	if !merge.Match("foo/bar") {
		t.Error("expected ignore to match")
	}

	if !merge.Match("bar/foo") {
		t.Error("expected ignore to match")
	}

	if !merge.Match("foo.exe") {
		t.Error("expected ignore to match")
	}
}
