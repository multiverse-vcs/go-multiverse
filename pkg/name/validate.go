package name

import (
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"
	record "github.com/libp2p/go-libp2p-record"
)

var _ record.Validator = (*Validator)(nil)

// Validator ensures records are valid.
type Validator struct{}

// Validate ensures that the signature matches the topic id.
func (v Validator) Validate(key string, value []byte) error {
	ns, k, err := record.SplitKey(key)
	if err != nil {
		return err
	}

	if ns != Namespace {
		return record.ErrInvalidRecordType
	}

	peerID, err := peer.Decode(k)
	if err != nil {
		return err
	}

	pub, err := peerID.ExtractPublicKey()
	if err != nil {
		return err
	}

	rec, err := RecordFromCBOR(value)
	if err != nil {
		return err
	}

	match, err := rec.Verify(pub)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("signature does not match")
	}

	return nil
}

// Select finds the best record by comparing sequence numbers.
func (v Validator) Select(key string, vals [][]byte) (int, error) {
	var ind int
	var max uint64

	for i, v := range vals {
		rec, err := RecordFromCBOR(v)
		if err != nil {
			return -1, err
		}

		if rec.Sequence >= max {
			ind, max = i, rec.Sequence
		}
	}

	return ind, nil
}
