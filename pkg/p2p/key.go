package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
)

// DefaultKeyType is the default private key type.
const DefaultKeyType = crypto.Ed25519

// GenerateKey returns a new private key.
func GenerateKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(DefaultKeyType, -1)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// EncodeKey returns a base64 encoded version of the key.
func EncodeKey(priv crypto.PrivKey) (string, error) {
	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}

	return crypto.ConfigEncodeKey(data), nil
}

// DecodeKey returns a key from a base64 encoded version.
func DecodeKey(encoded string) (crypto.PrivKey, error) {
	data, err := crypto.ConfigDecodeKey(encoded)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(data)
}
