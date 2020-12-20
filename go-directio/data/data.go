package data

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"log"
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
		return errors.New("Unable to store encoded data since received buffer is smaller.")
	}
	// n, err := buf.Write(to)
	copy(to, buf.Bytes())
	return nil
}

func Decode(from []byte) *SomeData {
	d := &SomeData{}
	dec := gob.NewDecoder(bytes.NewReader(from))
	err := dec.Decode(d)
	if err != nil && err != io.EOF {
		log.Println("Failed to decode bytes. Reason:", err)
		return nil
	}
	return d
}
