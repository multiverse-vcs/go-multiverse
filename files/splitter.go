package files

import "io"

// Splitter reads bytes from a file and splits them into blocks.
type Splitter interface {
	NextBlock() (*Block, error)
}

type sizeSplitter struct {
	reader io.Reader
	size   uint32
}

// NewSizeSplitter returns a splitter that creates blocks of the given size.
func NewSizeSplitter(reader io.Reader, size uint32) Splitter {
	return &sizeSplitter{reader, size}
}

// NextBlock reads bytes up to the block size of the splitter.
func (s *sizeSplitter) NextBlock() (*Block, error) {
	buffer := make([]byte, 0, s.size)

	size, err := io.ReadFull(s.reader, buffer)
	if err == nil {
		return NewBlock(buffer)
	}

	// finished reading blocks
	if err == io.ErrUnexpectedEOF {
		return NewBlock(buffer[:size])		
	}

	return nil, err
}
