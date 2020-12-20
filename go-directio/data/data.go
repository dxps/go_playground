package data

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/pkg/errors"
)

// Using a blocksize smaller than directio's 4K one.
const BlockSize = 64

type SomeData struct {
	Value uint32
}

func (d *SomeData) Encode(to []byte) error {
	buf := bytes.Buffer{}
	_ = gob.NewEncoder(&buf).Encode(*d)
	if len(to) < buf.Len() {
		return errors.New("cannot copy encoded into a smaller buffer")
	}
	copy(to, buf.Bytes())
	return nil
}

func Decode(from []byte) (*SomeData, error) {
	d := &SomeData{}
	dec := gob.NewDecoder(bytes.NewReader(from))
	err := dec.Decode(d)
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding data")
	}
	return d, nil
}
