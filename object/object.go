// Package object contains ipld definitions for multiverse objects.
package object

import (
	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-ipld-format"
)

func init() {
	format.Register(MCommit, DecodeCommitBlock)
}

const (
	// MCommit is the multicodec id for commits.
	MCommit = 0x300001
)

// DecodeCommitBlock decodes a commit from a block.
func DecodeCommitBlock(b blocks.Block) (format.Node, error) {
	return DecodeCommit(b.Cid(), b.RawData())
}
