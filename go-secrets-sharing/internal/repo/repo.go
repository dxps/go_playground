package repo

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/errors"
)

type Repo struct {
	memstore  map[string]string
	filestore *fileStore
}

// `NewRepo` creates a `Repo` instance.
// `dataFilePath` is the (relative or absolute) path to the file to persist the data on disk.
func NewRepo(dataFilePath string) (*Repo, error) {

	file, fileExists, err := initFilestore(dataFilePath)
	if err != nil {
		return nil, err
	}
	filestore := &fileStore{
		filePath: dataFilePath,
		file:     file,
		mu:       &sync.Mutex{},
	}
	var memstore map[string]string
	if fileExists {
		data, err := loadFromFile(dataFilePath)
		if err != nil {
			return nil, err
		}
		memstore = data
	} else {
		memstore = make(map[string]string, 0)
	}
	return &Repo{memstore, filestore}, nil
}

func initFilestore(filePath string) (file *os.File, fileExists bool, err error) {

	if _, err := os.Stat(filePath); err != nil {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, false, err
		}
		return file, false, nil
	}
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	fileExists = true
	return
}

func loadFromFile(filePath string) (map[string]string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res := make(map[string]string, 0)
	var entry dataEntry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		err := json.Unmarshal(sc.Bytes(), &entry)
		if err != nil {
			return nil, err
		}
		res[entry.ID] = entry.Secret
	}
	return res, nil
}

func (r *Repo) Add(hash, secret string) error {

	if _, exists := r.memstore[hash]; exists {
		return nil
	}
	r.memstore[hash] = secret
	return r.filestore.Append(&dataEntry{
		ID:     hash,
		Secret: secret,
	})
}

func (r *Repo) GetAndRemove(hash string) (secret string, err errors.AppError) {

	if val, exists := r.memstore[hash]; exists {
		delete(r.memstore, hash)
		return val, nil
	} else {
		return "", errors.EntryNotFound
	}
	// TODO: Persist the change on disk by overwriting the file.
}

// --------------------------------------
//             File Store
// --------------------------------------

type fileStore struct {
	filePath string      // Path to the file.
	file     *os.File    // Opened (append only, for minimal efficiency on adding secrets) file handle.
	mu       *sync.Mutex // Synchronization primitive.
}

type dataEntry struct {
	ID     string `json:"id"`
	Secret string `json:"data"`
}

func (fs *fileStore) Append(entry *dataEntry) error {

	fs.mu.Lock()
	defer fs.mu.Unlock()

	bs, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	bs = append(bs, 10) // Adding new line character.
	_, err = fs.file.Write(bs)
	if err != nil {
		return err
	}
	return nil
}
