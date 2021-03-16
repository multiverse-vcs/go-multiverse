package key

import (
	"github.com/libp2p/go-libp2p-core/crypto"
)

// DefaultType is the default private key type.
const DefaultType = crypto.Ed25519

// Generate returns a new private key.
func Generate() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateKeyPair(DefaultType, -1)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// Encode returns a base64 encoded version of the key.
func Encode(priv crypto.PrivKey) (string, error) {
	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}

	return crypto.ConfigEncodeKey(data), nil
}

// Decode returns a key from a base64 encoded version.
func Decode(encoded string) (crypto.PrivKey, error) {
	data, err := crypto.ConfigDecodeKey(encoded)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(data)
}
