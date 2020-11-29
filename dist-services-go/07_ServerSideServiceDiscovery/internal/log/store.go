package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	/* Encoding that is used when persisting the records sizes and index entries. */
	enc = binary.BigEndian
)

const (
	/* Number of bytes used to store the record's length. */
	lenWidth = 8
)

/* store is  */
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append the given bytes of a record to the store.
//
// It returns: `n`umber of bytes written, `pos`ition of the record in the store file,
// plus any `err`or in case of an issue.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size
	// First, writing the record's length.
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	// Second, writing the record.
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	// The total number of bytes written is the record's length
	// (stored as uint64, that's why lenWidth has the 8 (bytes) size).
	w += lenWidth
	s.size += uint64(w)
	// fmt.Printf("store.Append > len(p):%d, size:%d, ret n:%d pos:%d\n", len(p), s.size, uint64(w), pos)

	return uint64(w), pos, nil
}

// Read method reads and returns the record stored at the given `pos`ition.
func (s *store) Read(pos uint64) ([]byte, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	// Dirst, flush the buffer to disk.
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// ReadAt method reads `len(p)` bytes into p, beginning at the `off`set in the store file.
// It returns the number of bytes read, and an error in case of an issue.
func (s *store) ReadAt(p []byte, off int64) (int, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// Close method persists any buffered data and then closes the store's file.
func (s *store) Close() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
