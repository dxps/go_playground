package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytesSize = 32

// RememberToken is a helper function used for generating values of remember tokens.
func RememberToken() (string, error) {
	return genString(rememberTokenBytesSize)
}

// Bytes genrate n random bytes.
// It uses "crypto/rand" package, so it's save to use for strong uniqueness constraints.
func genBytes(n int) ([]byte, error) {

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String generate a byte slice of nBytes size and returns a Base64 encoded version of it.
func genString(nBytes int) (string, error) {
	b, err := genBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
