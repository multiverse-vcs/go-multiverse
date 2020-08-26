package files

import "io"

// Splitter reads bytes from a file and splits them into smaller blocks.
type Splitter interface {
	Reader() io.Reader
	NextBytes() ([]byte, error)
}

type sizeSplitter struct {
	reader io.Reader
	buffer []byte
}

// NewSizeSplitter returns a splitter that creates blocks of the given size.
func NewSizeSplitter(reader io.Reader, size uint32) Splitter {
	buffer := make([]byte, size)
	return &sizeSplitter{reader, buffer}
}

// NextBytes reads bytes up to the block size of the splitter.
func (s *sizeSplitter) NextBytes() ([]byte, error) {
	size, err := io.ReadFull(s.reader, s.buffer)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	return s.buffer[:size], nil
}

// Reader returns the io.Reader from the splitter.
func (s *sizeSplitter) Reader() io.Reader {
	return s.reader
}