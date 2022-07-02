package repo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/dxps/go_playground/go-secrets-sharing/internal/apperrs"
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
		data, err := loadFromFile(file)
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
			return nil, false, fmt.Errorf("Create file error: %v", err)
		}
		return file, false, nil
	}
	file, err = os.OpenFile(filePath, os.O_RDWR, 0600)
	fileExists = true
	return
}

func loadFromFile(f *os.File) (map[string]string, error) {

	res := make(map[string]string, 1)
	var entry dataEntry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		err := json.Unmarshal([]byte(sc.Text()), &entry)
		if err != nil {
			return nil, fmt.Errorf("Entry unmarshal error: %v", err)
		}
		res[entry.Key] = entry.Value
	}
	return res, nil
}

func (r *Repo) Add(hash, secret string) error {

	if _, exists := r.memstore[hash]; exists {
		return nil
	}
	r.memstore[hash] = secret
	return r.filestore.Append(hash, secret)
}

func (r *Repo) GetAndRemove(hash string) (secret string, err apperrs.AppError) {

	if val, exists := r.memstore[hash]; exists {
		delete(r.memstore, hash)
		r.filestore.Flush(r.memstore)
		return val, nil
	} else {
		return "", apperrs.ErrEntryNotFound
	}
}

// --------------------------------------
//             File Store
// --------------------------------------

type fileStore struct {
	filePath string      // Path to the file.
	file     *os.File    // Open file handle.
	mu       *sync.Mutex // File write synchronization primitive.
}

type dataEntry struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

func (fs *fileStore) Append(key, val string) error {

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := fs.writeEntry(key, val); err != nil {
		return err
	}
	return nil
}

func (fs *fileStore) Flush(data map[string]string) error {

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := fs.file.Truncate(0); err != nil {
		return fmt.Errorf("File truncate error: %v", err)
	}
	if _, err := fs.file.Seek(0, 0); err != nil {
		return fmt.Errorf("File seek error: %v", err)
	}
	for k, v := range data {
		if err := fs.writeEntry(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (fs *fileStore) writeEntry(key, val string) error {

	bs, err := json.Marshal(dataEntry{Key: key, Value: val})
	if err != nil {
		return fmt.Errorf("Entry marshal error: %v", err)
	}
	bs = append(bs, 10) // Adding new line character.
	_, err = fs.file.Write(bs)
	if err != nil {
		return fmt.Errorf("Write entry error: %v", err)
	}
	return nil
}
