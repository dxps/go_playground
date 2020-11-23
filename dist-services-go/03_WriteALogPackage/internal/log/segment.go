package log

import (
	"fmt"
	"os"
	"path"

	"google.golang.org/protobuf/proto"
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

// Append adds the given `rec`ord to the segment and returns its `off`set.
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
		// index's offsets are relative to the base offset
		uint32(s.nextOffset-uint64(s.baseOffset)),
		pos,
	); err != nil {
		return 0, err
	}
	s.nextOffset++
	return cur, nil
}
