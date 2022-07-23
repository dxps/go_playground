package repo

import (
	"log"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/apperrs"
	"github.com/dxps/go_playground/go-secrets-sharing/internal/repo/filestore"
)

type Repo struct {
	memStore  map[string]string
	fileStore *filestore.FileStore
}

// `NewRepo` creates a `Repo` instance.
// `dataFilePath` is the (relative or absolute) path to the file to persist the data on disk.
func NewRepo(dataFilePath, storePass, storeSalt string) (*Repo, error) {

	fileRepo, fileExists, err := filestore.NewFileStore(dataFilePath, storePass, storeSalt)
	if err != nil {
		return nil, err
	}
	var memRepo map[string]string
	if fileExists {
		data, err := fileRepo.LoadFile()
		if err != nil {
			return nil, err
		}
		memRepo = data
		log.Printf("%d entries loaded from file.", len(data))
	} else {
		memRepo = make(map[string]string, 0)
	}
	return &Repo{memRepo, fileRepo}, nil
}

func (r *Repo) Add(hash, secret string) error {

	if _, exists := r.memStore[hash]; exists {
		return nil
	}
	r.memStore[hash] = secret
	return r.fileStore.Append(hash, secret)
}

func (r *Repo) GetAndRemove(hash string) (secret string, err apperrs.AppError) {

	if val, exists := r.memStore[hash]; exists {
		delete(r.memStore, hash)
		r.fileStore.Flush(r.memStore)
		return val, nil
	} else {
		return "", apperrs.ErrEntryNotFound
	}
}
