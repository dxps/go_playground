package internal

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/devisions/go-playground/go-directio/internal/data"
	"github.com/pkg/errors"
)

type ConsumerData struct {
	Data         *data.SomeData
	FromFilepath string
}

// GetFileForReading looks for the next available file to be used for reading.
// If `curr`ent (existing) file reached the max size, it looks for the next available one, if any.
// If `curr`ent (existing) file did not reached the max size, nil is returned, meaning that it can be reused.
func GetFileForReading(curr *os.File, path string, maxsize int64) (*os.File, error) {
	// First, let's check the current file size, if provided.
	if curr != nil {
		fi, err := curr.Stat()
		if err != nil {
			return nil, errors.Wrap(err, "trying to get the current file info")
		}
		if fi.Size() >= maxsize {
			file, err := getNextFileNameForReading(path, curr.Name())
			if err != nil {
				return nil, errors.Wrap(err, "trying to get next file for reading")
			}
			filepath := path + string(os.PathSeparator) + file
			return data.OpenFileForReading(filepath)
		}
		return nil, nil // Existing is not returned again.
	}
	// Nothing was read before, (no `curr`ent exists).
	file, err := getFirstFileNameForReading(path)
	if err != nil {
		return nil, errors.Wrap(err, "trying to get first file for reading")
	}
	filepath := path + string(os.PathSeparator) + file
	return data.OpenFileForReading(filepath)
}

func getFirstFileNameForReading(iopath string) (string, error) {
	fs, err := ioutil.ReadDir(iopath)
	if err != nil {
		return "", err
	}
	for _, f := range fs {
		if ".dat" == path.Ext(f.Name()) {
			return f.Name(), nil
		}
	}
	return "", os.ErrNotExist
}

func getNextFileNameForReading(iopath string, lastFilePath string) (string, error) {
	fs, err := ioutil.ReadDir(iopath)
	if err != nil {
		return "", err
	}
	var fi, f1 os.FileInfo
	var lastFilePathFound bool
	for _, f := range fs {
		fname := f.Name()
		if ".dat" == path.Ext(fname) {
			if f1 == nil {
				f1 = f
			}
			if fname == lastFilePath {
				lastFilePathFound = true
				continue
			}
			if lastFilePathFound {
				return fi.Name(), nil
			}
			fi = f
		}
	}
	if lastFilePathFound /* lastFilePath was found, but there was no next to return immediately. */ ||
		f1 == nil /* Or lastFilePath was not found, but there is no first (any) file. */ {
		return "", os.ErrNotExist
	}
	// Returning the 1st file found.
	return f1.Name(), nil
}
