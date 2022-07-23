package filestore

import (
	"bufio"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileStore struct {
	filePath string      // Path to the file.
	file     *os.File    // Open file handle.
	mu       *sync.Mutex // File write synchronization primitive.
	cipher   cipher.AEAD
	nonce    []byte
}

type dataEntry struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

func NewFileStore(dataFilePath, storePass, storeSalt string) (r *FileStore, fileExists bool, err error) {

	var file *os.File
	file, fileExists, err = initFileRepo(dataFilePath)
	if err != nil {
		return
	}
	cipher, nonce, err := initCrypto(storePass, storeSalt)
	r = &FileStore{
		filePath: dataFilePath,
		file:     file,
		mu:       &sync.Mutex{},
		cipher:   cipher,
		nonce:    nonce,
	}
	return
}

func (s *FileStore) LoadFile() (map[string]string, error) {

	res := make(map[string]string, 1)
	var entry dataEntry
	var err error
	sc := bufio.NewScanner(s.file)
	for sc.Scan() {
		// Decrypt before unmarshalling.
		bs := []byte(sc.Text())
		bs, err = s.decrypt(bs)
		if err != nil {
			return nil, fmt.Errorf("Entry decryption error: %v", err)
		}
		err = json.Unmarshal(bs, &entry)
		if err != nil {
			return nil, fmt.Errorf("Entry unmarshal error: %v", err)
		}
		res[entry.Key] = entry.Value
	}
	return res, nil
}

func (s *FileStore) Append(key, val string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.writeEntry(key, val); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) Flush(data map[string]string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.file.Truncate(0); err != nil {
		return fmt.Errorf("File truncate error: %v", err)
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("File seek error: %v", err)
	}
	for k, v := range data {
		if err := s.writeEntry(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStore) writeEntry(key, val string) error {

	bs, err := json.Marshal(dataEntry{Key: key, Value: val})
	if err != nil {
		return fmt.Errorf("Entry marshal error: %v", err)
	}
	bs = append(bs, 10) // Adding new line character.
	// Encrypting before storing.
	bs, err = s.encrypt(bs)
	if err != nil {
		return fmt.Errorf("Entry encryption error: %v", err)
	}
	_, err = s.file.Write(bs)
	if err != nil {
		return fmt.Errorf("Write entry error: %v", err)
	}
	return nil
}
