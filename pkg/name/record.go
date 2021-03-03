package name

import (
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func init() {
	cbornode.RegisterCborType(Payload{})
	cbornode.RegisterCborType(Record{})
}

// Payload contains record data.
type Payload struct {
	// Value contains the payload data.
	Value []byte
	// Sequence is a version and nonce.
	Sequence uint64
}

// Record is used for record signatures.
type Record struct {
	// Signature is a signature of the payload.
	Signature []byte

	Payload
}

// RecordFromCBOR decodes an envelope from cbor.
func RecordFromCBOR(data []byte) (*Record, error) {
	var rec Record
	if err := cbornode.DecodeInto(data, &rec); err != nil {
		return nil, err
	}

	return &rec, nil
}

// NewRecord returns a new record.
func NewRecord(value []byte) *Record {
	return &Record{
		Payload: Payload{Value: value},
	}
}

// Bytes returns the raw bytes of the record.
func (r *Record) Bytes() ([]byte, error) {
	return cbornode.DumpObject(r)
}

// Sign creates a signature for the record payload.
func (r *Record) Sign(key crypto.PrivKey) error {
	payload, err := cbornode.DumpObject(r.Payload)
	if err != nil {
		return err
	}

	signature, err := key.Sign(payload)
	if err != nil {
		return err
	}

	r.Signature = signature
	return nil
}

// Verify checks if the signature of the payload is valid.
func (r *Record) Verify(key crypto.PubKey) (bool, error) {
	payload, err := cbornode.DumpObject(r.Payload)
	if err != nil {
		return false, err
	}

	return key.Verify(payload, r.Signature)
}
