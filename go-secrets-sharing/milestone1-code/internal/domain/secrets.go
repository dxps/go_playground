package domain

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/apperrs"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo"
)

type Secrets struct {
	repo *repo.Repo
}

// `NewSecrets` creates a new instance of `Secrets`.
func NewSecrets(repo *repo.Repo) *Secrets {
	return &Secrets{repo}
}

// `Store` persists the plaintext secret and return the MD5 hash of it.
func (s *Secrets) Store(plainSecret string) (string, apperrs.AppError) {
	sh := s.getMD5Hash(plainSecret)
	if err := s.repo.Add(sh, plainSecret); err != nil {
		return "", err
	}
	return sh, nil
}

// `Retrieve` retrieves the associated secret, if exists.
// If it exists, the (hash, secret) pair is removed, since that hash can be used only once.
func (s *Secrets) Retrieve(hash string) (string, apperrs.AppError) {
	return s.repo.GetAndRemove(hash)
}

func (s *Secrets) getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
