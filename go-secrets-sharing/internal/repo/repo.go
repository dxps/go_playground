package repo

import "github.com/dxps/go_playground/go-secrets-sharing/internal/errors"

type Repo struct {
	memstore map[string]string
}

func NewRepo() *Repo {
	return &Repo{
		memstore: make(map[string]string),
	}
}

func (r *Repo) Add(hash, secret string) {
	r.memstore[hash] = secret
}

func (r *Repo) GetAndRemove(hash string) (secret string, err errors.AppError) {

	if val, exists := r.memstore[hash]; exists {
		delete(r.memstore, hash)
		return val, nil
	} else {
		return "", errors.EntryNotFound
	}
}
