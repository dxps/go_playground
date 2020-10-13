package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC is a wrapper around "crypto/hmac" package, exposing just the features we use.
type HMAC struct {
	hmac hash.Hash
}

// Hash returnes the hash of the provided input string using HMAC
// with the secret key provided when the HMAC object was created.
func (h HMAC) Hash(input string) string {

	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}

// NewHMAC creates and returns a new HMAC object.
func NewHMAC(key string) HMAC {

	h := hmac.New(sha256.New, []byte(key))
	return HMAC{h}
}
