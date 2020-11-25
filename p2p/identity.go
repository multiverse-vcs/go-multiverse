package p2p

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// DefaultKeyType is the type of key to use.
const DefaultKeyType = crypto.Ed25519

// Identity is a wrapper for keys.
type Identity struct {
	// Key is the private key of the peer.
	Key crypto.PrivKey
}

// MarshalJSON encodes the identity into json.
func (i *Identity) MarshalJSON() ([]byte, error) {
	data, err := crypto.MarshalPrivateKey(i.Key)
	if err != nil {
		return nil, err
	}

	base64 := crypto.ConfigEncodeKey(data)
	return json.Marshal(base64)
}

// UnmarshalJSON decodes the identity from json.
func (i *Identity) UnmarshalJSON(data []byte) error {
	var base64 string
	if err := json.Unmarshal(data, &base64); err != nil {
		return err
	}

	enc, err := crypto.ConfigDecodeKey(base64)
	if err != nil {
		return err
	}

	i.Key, err = crypto.UnmarshalPrivateKey(enc)
	return err
}
