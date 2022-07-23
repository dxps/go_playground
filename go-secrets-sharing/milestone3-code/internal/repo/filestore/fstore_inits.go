package filestore

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"golang.org/x/crypto/scrypt"
)

func initFileRepo(filePath string) (file *os.File, fileExists bool, err error) {

	if _, err := os.Stat(filePath); err != nil {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, false, fmt.Errorf("Create file error: %v", err)
		}
		return file, false, nil
	}
	file, err = os.OpenFile(filePath, os.O_RDWR, 0600)
	fileExists = true
	return
}

func initCrypto(password, salt string) (cipher.AEAD, []byte, error) {

	key, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	return gcm, nonce, err
}
