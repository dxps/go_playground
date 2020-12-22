package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/devisions/go-playground/go-directio/internal/data"
	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

const STATE_FILE = "consumer.state"

type ConsumerState struct {
	ReadFilepath  string
	ReadBlocks    int
	maxBlocks     int
	saveFilepath  string
	saveBlocksize int
}

func (s *ConsumerState) encode(to []byte) error {
	buf := bytes.Buffer{}
	_ = gob.NewEncoder(&buf).Encode(*s)
	if len(to) < buf.Len() {
		return errors.New("cannot copy encoded into a smaller buffer")
	}
	copy(to, buf.Bytes())
	return nil
}

func decode(from []byte) (*ConsumerState, error) {
	s := &ConsumerState{}
	dec := gob.NewDecoder(bytes.NewReader(from))
	err := dec.Decode(s)
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding data")
	}
	return s, nil
}

func (s *ConsumerState) SaveToFile() error {
	f, err := data.OpenFileForWriting(s.saveFilepath, false)
	if err != nil {
		return errors.Wrap(err, "opening file for writing the state")
	}
	defer func() { _ = f.Close() }()
	block := directio.AlignedBlock(s.saveBlocksize)
	err = s.encode(block)
	if err != nil {
		return errors.Wrap(err, "encoding state")
	}
	_, err = f.Write(block)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("writing the state to file %s", s.saveFilepath))
	}
	return nil
}

func InitConsumerState(path string, saveBlocksize int) (*ConsumerState, error) {
	filepath := path + string(os.PathSeparator) + STATE_FILE
	f, err := data.OpenFileForReading(filepath)
	if err != nil {
		if os.IsNotExist(errors.Cause(err)) {
			// There is no `consumer.state` file. We'll return an empty object.
			// The file will eventually be created first time the state is saved.
			return &ConsumerState{
				saveFilepath:  filepath,
				saveBlocksize: saveBlocksize,
			}, nil
		}
		return nil, err
	}
	block := directio.AlignedBlock(saveBlocksize)
	_, err = f.Read(block)
	if err != nil {
		return nil, err
	}
	s, derr := decode(block)
	if derr != nil {
		return nil, derr
	}
	s.saveFilepath = filepath
	s.saveBlocksize = saveBlocksize
	return s, nil
}

func (s *ConsumerState) UseNew(filepath string) {
	if filepath != "" && s.ReadFilepath != filepath {
		s.ReadFilepath = filepath
		s.ReadBlocks = 1
	}
}

func (s *ConsumerState) IsEmpty() bool {
	return s.ReadFilepath == "" && s.ReadBlocks == 0
}

func (s *ConsumerState) SeekOffset() int64 {
	return int64(s.ReadBlocks) * int64(s.saveBlocksize)
}
