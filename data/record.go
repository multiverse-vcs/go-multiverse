package data

import (
	cbornode "github.com/ipfs/go-ipld-cbor"
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

// NewRecord returns a record containing the given payload.
func NewRecord(payload []byte, sequence uint64, signature []byte) *Record {
	return &Record{
		Payload:   payload,
		Sequence:  sequence,
		Signature: signature,
	}
}
