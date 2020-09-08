package files

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-block-format"
	"github.com/multiformats/go-multihash"
)

// BlockReader implements io.Reader.
type BlockReader struct {
	reader io.Reader
}

// NewBlockReader returns an io.Reader that reads blocks.
func NewBlockReader(reader io.Reader) *BlockReader {
	return &BlockReader{reader}
}

// Read reads bytes up to the size of the buffer.
func (br *BlockReader) Read(buffer []byte) (n int, err error) {
	return io.ReadFull(br.reader, buffer)
}

// ReadBlock returns a new block of the given size.
func (br *BlockReader) ReadBlock(size uint32) (blocks.Block, error) {
	data := make([]byte, size)
	if _, err := br.Read(data); err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	cid := cid.NewCidV1(cid.Raw, hash)
	return blocks.NewBlockWithCid(data, cid)
}