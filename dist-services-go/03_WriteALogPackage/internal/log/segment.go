package log

import (
	api "devisions.org/go-dist-svcs/log/api/v1"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"os"
	"path"
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint64
	config                 Config
}

func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {

	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}
	var err error
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}
	if off, _, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
		fmt.Printf("[tmp] newSegment > Starting with baseOffset:%d, got index's off:%d and set nextOffset:%d\n",
			baseOffset, off, s.nextOffset)
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1
		fmt.Printf("[tmp] newSegment > Starting with baseOffset:%d, got index's off:%d and set nextOffset:%d\n",
			baseOffset, off, s.nextOffset)
	}
	return s, nil
}

// Append adds the given `rec`ord to the segment and returns its `off`set
// that is relative to the segment's base offset.
func (s *segment) Append(rec *api.Record) (off uint64, err error) {

	cur := s.nextOffset
	rec.Offset = cur
	p, err := proto.Marshal(rec)
	if err != nil {
		return 0, err
	}
	_, pos, err := s.store.Append(p)
	if err != nil {
		return 0, err
	}
	if err = s.index.Write(
		// index's offsets are relative to the segment's base offset.
		uint32(s.nextOffset-s.baseOffset),
		pos,
	); err != nil {
		return 0, err
	}
	s.nextOffset++
	return cur, nil
}

func (s *segment) Read(off uint64) (*api.Record, error) {

	_, pos, err := s.index.Read(int64(off - s.baseOffset))
	if err != nil {
		return nil, err
	}
	p, err := s.store.Read(pos)
	if err != nil {
		return nil, err
	}
	rec := &api.Record{}
	err = proto.Unmarshal(p, rec)
	return rec, err
}

// IsMaxed tells whether the segment has reached its max size.
func (s *segment) IsMaxed() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes || s.index.size >= s.config.Segment.MaxIndexBytes
}

func (s *segment) Remove() error {

	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.index.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		return err
	}
	return nil
}

func (s *segment) Close() error {

	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	return nil
}

// nearestMultiple returns the multiple of `k` that is nearest and lesser than j.
func nearestMultiple(j, k uint64) uint64 {

	if j >= 0 {
		return (j / k) * k
	}
	return ((j - k + 1) / k) * k
}
