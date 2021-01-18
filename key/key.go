// Package key implements a keystore.
package key

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	// DefaultType is the default private key type.
	DefaultType = crypto.Ed25519
	// DefaultKeyName is the name of the default key.
	DefaultKeyName = "default"
)

// Keystore stores keys in directory.
type Keystore struct {
	root string
}

// NewKeystore returns a  new keystore in the given root directory.
func NewKeystore(root string) (*Keystore, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &Keystore{root}, nil
}

// DefaultKey returns the default key from the store.
func (ks *Keystore) DefaultKey() (crypto.PrivKey, error) {
	key, err := ks.GetKey(DefaultKeyName)
	if os.IsNotExist(err) {
		return ks.NewKey(DefaultKeyName)
	}

	return key, nil
}

// NewKey generates a new key with the given name.
func (ks *Keystore) NewKey(name string) (crypto.PrivKey, error) {
	kpath := path.Join(ks.root, name)
	if _, err := os.Stat(kpath); err == nil {
		return nil, errors.New("key already exists")
	}

	priv, _, err := crypto.GenerateKeyPair(DefaultType, -1)
	if err != nil {
		return nil, err
	}

	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	encoded := crypto.ConfigEncodeKey(data)
	if err := ioutil.WriteFile(kpath, []byte(encoded), 0644); err != nil {
		return nil, err
	}

	return priv, nil
}

// GetKey returns the private key from keystore.
func (ks *Keystore) GetKey(name string) (crypto.PrivKey, error) {
	kpath := path.Join(ks.root, name)
	if _, err := os.Stat(kpath); err != nil {
		return nil, err
	}

	encoded, err := ioutil.ReadFile(kpath)
	if err != nil {
		return nil, err
	}

	data, err := crypto.ConfigDecodeKey(string(encoded))
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(data)
}
