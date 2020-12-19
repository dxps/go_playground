package data

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
)

type SomeData struct {
	Value uint32
}

func (d *SomeData) Encode(to []byte) {
	buf := bytes.Buffer{}
	_ = gob.NewEncoder(&buf).Encode(d)
	_, err := buf.Write(to)
	if err != nil {
		log.Printf("Failed to encode the data. Reason: %s", err)
	}
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
