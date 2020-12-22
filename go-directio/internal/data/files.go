package data

import (
	"fmt"
	"os"

	"github.com/ncw/directio"
	"github.com/pkg/errors"
)

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
		_ = f.Close()
		return nil, errors.Wrap(err, fmt.Sprintf("while opening file %s for reading", filepath))
	}
	return f, nil
}

func DeleteFileIfReachedMaxSize(filepath string, maxSize int64) (bool, error) {
	f, err := OpenFileForReading(filepath)
	defer func() { _ = f.Close() }()
	if err != nil {
		if !os.IsNotExist(errors.Cause(err)) {
			return false, err
		}
		return false, nil
	}
	fi, err := f.Stat()
	if err != nil {
		return false, errors.Wrap(err, "trying to get file stat")
	}
	if fi.Size() >= maxSize {
		err := DeleteFile(filepath)
		return true, errors.Wrap(err, "deleting file "+filepath)
	}
	return false, nil
}

func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
