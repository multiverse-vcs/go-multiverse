package data

import (
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/libp2p/go-libp2p-core/crypto"
)

// Record contains named record info.
type Record struct {
	// Payload contains the record data.
	Payload []byte
	// Sequence is a version identifier.
	Sequence uint64
	// Signature is a signature of the payload.
	Signature []byte
}

// RecordFromCBOR decodes a record from an ipld node.
func RecordFromCBOR(data []byte) (*Record, error) {
	var rec Record
	if err := cbornode.DecodeInto(data, &rec); err != nil {
		return nil, err
	}

	return &rec, nil
}

// NewRecord returns a signed record containing the given payload.
func NewRecord(payload []byte, sequence uint64, key crypto.PrivKey) (*Record, error) {
	signature, err := key.Sign(payload)
	if err != nil {
		return nil, err
	}

	return &Record{
		Payload:   payload,
		Sequence:  sequence,
		Signature: signature,
	}, nil
}

// Encode returns the record raw bytes.
func (r *Record) Encode() ([]byte, error) {
	return cbornode.DumpObject(r)
}

// Verify returns true if the payload signature matches the given public key.
func (r *Record) Verify(key crypto.PubKey) (bool, error) {
	return key.Verify(r.Payload, r.Signature)
}
