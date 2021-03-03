package name

import (
	"testing"

	"github.com/multiverse-vcs/go-multiverse/internal/p2p"
)

func TestRecordSignature(t *testing.T) {
	key, err := p2p.GenerateKey()
	if err != nil {
		t.Fatal("failed to generate key")
	}

	rec := &Record{}
	rec.Sequence = 1
	rec.Value = []byte("bafyreiakcek7msekxf67tdvmsjbkyxus4iy6j3ed5n3hgfelly26hlm2lu")

	if err := rec.Sign(key); err != nil {
		t.Fatal("failed to sign envelope")
	}

	match, err := rec.Verify(key.GetPublic())
	if err != nil {
		t.Fatal("failed to verify envelope")
	}

	if !match {
		t.Error("envelope signature does not match")
	}
}
