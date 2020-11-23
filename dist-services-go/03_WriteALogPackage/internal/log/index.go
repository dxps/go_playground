package log

import (
	"io"
	"os"

	"github.com/tysontate/gommap"
)

var (
	// number of bytes used for storing the record's offset
	offWidth uint64 = 4
	// number of bytes used for storing the record's position
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth
)

// An index entry contains two fields:
// 1. The record's offset (stored as uint32)
// 2. Its position (stored as uint64) in the store file.

// index represents the index file where the `off`set and `pos`ition of each record
// is stored for quicker search.
type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {

	idx := &index{file: f}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	// Grow the file to the maximum index size before memory-mapping it.
	// After being memory mapped, the file cannot be resized.
	if err = os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

// Close method synchronizes the memory-mapped file with the persisted file,
// flushes any file system cached data, and closes the file.
func (i *index) Close() error {

	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	// Truncate the file according to the data that's actually in it.
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}

// Read returns the associated record's `off`set and `pos`ition in the store
// based on the given segment offset. This `setOff` is relative to the segment's
// base offset: 0 is the offset of the 1st entry, 1 is the 2nd entry, and so on.
func (i *index) Read(segOff int64) (off uint32, pos uint64, err error) {

	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if segOff == -1 {
		off = uint32((i.size / entWidth) - 1)
	} else {
		off = uint32(segOff)
	}
	pos = uint64(off) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}
	off = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])
	return off, pos, nil
}

// Write appends the given `off`set and `pos`ition to the index.
func (i *index) Write(off uint32, pos uint64) error {

	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += uint64(entWidth)
	return nil
}

// Name returns the index's file path.
func (i *index) Name() string {
	return i.file.Name()
}
