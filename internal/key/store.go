package key

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// DirName is the keystore directory name.
const DirName = "keystore"

// Store is used to read and write keys.
type Store struct {
	path string
}

// NewStore returns a store with the given root directory.
func NewStore(dir string) (*Store, error) {
	path := filepath.Join(dir, DirName)
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return &Store{path}, nil
}

// Delete removes the key with the given name from the store.
func (s *Store) Delete(name string) error {
	path := filepath.Join(s.path, name)

	return os.Remove(path)
}

// Get returns the key with the given name from the store.
func (s *Store) Get(name string) (crypto.PrivKey, error) {
	path := filepath.Join(s.path, name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Decode(string(data))
}

// Has returns true if the store contains the key with the given name.
func (s *Store) Has(name string) (bool, error) {
	path := filepath.Join(s.path, name)

	_, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Put stores the given key in the store under the given name.
func (s *Store) Put(name string, key crypto.PrivKey) error {
	has, err := s.Has(name)
	if err != nil {
		return err
	}

	if has {
		return errors.New("key already exists")
	}

	data, err := Encode(key)
	if err != nil {
		return err
	}

	path := filepath.Join(s.path, name)
	return os.WriteFile(path, []byte(data), 0644)
}
