package filestore

import (
	"crypto/rand"
	"io"
)

func (s *FileStore) encrypt(plaintext []byte) ([]byte, error) {

	if _, err := io.ReadFull(rand.Reader, s.nonce); err != nil {
		return nil, err
	}

	return s.cipher.Seal(s.nonce, s.nonce, plaintext, nil), nil
}

func (s *FileStore) decrypt(encData []byte) ([]byte, error) {

	nonce := encData[:s.cipher.NonceSize()]
	encData = encData[s.cipher.NonceSize():]
	data, err := s.cipher.Open(nil, nonce, encData, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
