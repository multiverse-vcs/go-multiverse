package files

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// Block implements the IPFS Block interface.
type Block struct {
	cid  cid.Cid
	data []byte
}

// NewBlock creates a Block from data.
func NewBlock(data []byte) (*Block, error) {
	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	return &Block{data: data, cid: cid.NewCidV1(cid.Raw, hash)}, nil
}

// Multihash returns the hash contained in the block CID.
func (b *Block) Multihash() multihash.Multihash {
	return b.cid.Hash()
}

// RawData returns the block raw contents as a byte slice.
func (b *Block) RawData() []byte {
	return b.data
}

// Cid returns the content identifier of the block.
func (b *Block) Cid() cid.Cid {
	return b.cid
}

// String provides a human-readable representation of the block CID.
func (b *Block) String() string {
	return fmt.Sprintf("[Block %s]", b.Cid())
}

// Loggable returns a go-log loggable item.
func (b *Block) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"block": b.Cid().String(),
	}
}