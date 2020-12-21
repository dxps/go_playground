package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

// GetFileForIO looks for the next available file to be used for writing.
// If existing file reached the max size, it initializes a new file (with a new name and opens it)
// and returns it.
func GetFileForWriting(curr *os.File, path string, maxsize int64) (*os.File, error) {
	// First, let's check the current file size, if it exists.
	if curr != nil {
		fi, err := curr.Stat()
		if err != nil {
			return nil, errors.Wrap(err, "trying to get the current file info")
		}
		if fi.Size() >= maxsize {
			return openNewFileForWriting(path)
		}
		return curr, nil
	}
	// This should be the startup phase, since there's no `curr` provided.
	file, size, err := getFileNameForIO(path, "", true)
	if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return nil, errors.Wrap(err, "trying to get file for writing")
	}
	if os.IsNotExist(errors.Cause(err)) || file == "" || size >= maxsize {
		return openNewFileForWriting(path)
	}
	filepath := path + string(os.PathSeparator) + file
	return OpenFileForWriting(filepath, true)
}

// GetFileForReading looks for the next available file to be used for reading.
// If existing file reached the max size, it initializes a new file (with a new name and opens it)
// and returns it.
func GetFileForReading(curr *os.File, path string, lastReadFilePath string, maxsize int64) (*os.File, error) {
	// First, let's check the current file size, if it exists.
	if curr != nil {
		fi, err := curr.Stat()
		if err != nil {
			return nil, errors.Wrap(err, "trying to get the current file info")
		}
		if fi.Size() >= maxsize {
			file, _, err := getFileNameForIO(path, curr.Name(), false)
			if err != nil {
				return nil, err
			}
			filepath := path + string(os.PathSeparator) + file
			return OpenFileForReading(filepath)
		}
		return curr, nil
	}
	// Here, nothing was read before, so there is no `curr` provided.
	file, _, err := getFileNameForIO(path, lastReadFilePath, false)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrap(err, "trying to get file for reading")
		}
		return nil, err // Send back the 'non existence' case as is.
	}
	filepath := path + string(os.PathSeparator) + file
	return OpenFileForReading(filepath)
}

// It returns the (latest, first, or next after last) file's name and size, if no `lastFile` (aka "") was provided.
// If `lastFilePath` is provided, it returns the next file found in the `iopath`.
// It returns `os.ErrNotExist` error if no file was found in the provided `iopath`.
func getFileNameForIO(iopath string, lastFilePath string, latest bool) (string, int64, error) {
	fs, err := ioutil.ReadDir(iopath)
	if err != nil {
		return "", 0, errors.Wrap(err, fmt.Sprintf("looking for files on path '%s'", iopath))
	}
	var fi os.FileInfo
	var lastFilePathFound bool
	for _, f := range fs {
		fname := f.Name()
		if ".dat" == path.Ext(fname) {
			if fname == lastFilePath {
				lastFilePathFound = true
				continue
			}
			if lastFilePathFound {
				return fi.Name(), fi.Size(), nil
			}
			fi = f
		}
	}
	if fi == nil {
		return "", 0, os.ErrNotExist
	}
	return fi.Name(), fi.Size(), nil
}

func openNewFileForWriting(path string) (*os.File, error) {
	filepath := fmt.Sprintf("%s%s%d.dat", path, string(os.PathSeparator), time.Now().UTC().Unix())
	return OpenFileForWriting(filepath, false)
}

func OpenFileForWriting(filepath string, append bool) (*os.File, error) {
	mode := os.O_CREATE | os.O_WRONLY
	if append {
		mode = mode | os.O_APPEND
	}
	f, err := directio.OpenFile(filepath, mode, 0665)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("while opening file %s for writing", filepath))
	}
	if append {
		_, err = f.Seek(0, 2) // Go to the end of the file.
		if err != nil {
			return nil, errors.Wrap(err, "seeking to the file end")
		}
	}
	return f, nil
}

func OpenFileForReading(filepath string) (*os.File, error) {
	f, err := directio.OpenFile(filepath, os.O_RDONLY, 0665)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("while opening file %s for reading", filepath))
	}
	return f, nil
}

func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
