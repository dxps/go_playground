package server

import (
	"fmt"
	"sync"
)

// Log is a commit log structure.
type Log struct {
	mu      sync.Mutex
	records []Record
}

// NewLog creates a new instance of a Log.
func NewLog() *Log {
	return &Log{}
}

// Append adds the provided record to the log.
func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)
	return record.Offset, nil
}

// Read fetches the record that (might) exists at the provided offset.
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}

// Record is an entry of a commit log.
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// ErrOffsetNotFound is an error returned when the offset is not found.
var ErrOffsetNotFound = fmt.Errorf("offset not found")
