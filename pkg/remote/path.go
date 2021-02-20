package remote

import (
	"errors"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
)

// Path is a repository identifier.
type Path string

// NewPath returns a new path.
func NewPath(id peer.ID, name string) Path {
	return Path(strings.Join([]string{id.Pretty(), name}, "/"))
}

// String returns the path as a string.
func (p Path) String() string {
	return string(p)
}

// PeerID returns the peer ID part of the path.
func (p Path) PeerID() (peer.ID, error) {
	parts, err := p.Split()
	if err != nil {
		return "", err
	}

	return peer.Decode(parts[0])
}

// Name returns the name part of the path.
func (p Path) Name() (string, error) {
	parts, err := p.Split()
	if err != nil {
		return "", err
	}

	return parts[1], nil
}

// Split returns the path parts.
func (p Path) Split() ([]string, error) {
	parts := strings.SplitN(string(p), "/", 3)
	if len(parts) < 2 {
		return nil, errors.New("invalid path")
	}

	return parts, nil
}

// Verify checks if the path is valid.
func (p Path) Verify() error {
	_, err := p.PeerID()
	if err != nil {
		return err
	}

	return nil
}
