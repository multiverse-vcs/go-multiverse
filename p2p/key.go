package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
)

// KeyType is the default private key type.
const KeyType = crypto.Ed25519

// GenerateKey creates a new private key.
func GenerateKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(KeyType, -1)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// DecodeKey decodes the private key from the encoded string.
func DecodeKey(encoded string) (crypto.PrivKey, error) {
	data, err := crypto.ConfigDecodeKey(encoded)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(data)
}

// EncodeKey encodes the private key into an encoded string.
func EncodeKey(key crypto.PrivKey) (string, error) {
	data, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return "", err
	}

	return crypto.ConfigEncodeKey(data), nil
}
