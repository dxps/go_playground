package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	writeData  = []byte("hello world")
	writeWidth = uint64(len(writeData) + lenWidth)
)

func TestStoreAppendRead(t *testing.T) {

	f, err := ioutil.TempFile("", "store_append_read_store_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// testing a new store created from scratch
	s, err := newStore(f)
	require.NoError(t, err)
	testAppend(t, s)
	testRead(t, s)
	testReadAt(t, s)

	// testing an existing store, used for example after a restart
	s, err = newStore(f)
	require.NoError(t, err)
	testRead(t, s)
}

func testAppend(t *testing.T, s *store) {

	t.Helper()
	for i := uint64(1); i < 4; i++ {
		n, pos, err := s.Append(writeData)
		require.NoError(t, err)
		require.Equal(t, pos+n, writeWidth*i)
	}
}

func testRead(t *testing.T, s *store) {

	t.Helper()
	var pos uint64
	for i := uint64(1); i < 4; i++ {
		readData, err := s.Read(pos)
		require.NoError(t, err)
		require.Equal(t, readData, writeData)
		pos += writeWidth
	}
}

func testReadAt(t *testing.T, s *store) {

	t.Helper()
	for i, off := uint64(1), int64(0); i < 4; i++ {
		b := make([]byte, lenWidth)
		n, err := s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, lenWidth, n)
		off += int64(n)

		size := enc.Uint64(b)
		b = make([]byte, size)
		n, err = s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, writeData, b)
		require.Equal(t, int(size), n)
		off += int64(n)
	}
}

func TestStoreClose(t *testing.T) {

	f, err := ioutil.TempFile("", "close_store_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	s, err := newStore(f)
	require.NoError(t, err)

	_, _, err = s.Append(writeData)
	require.NoError(t, err)

	_, beforeSize, err := openFile(f.Name())
	require.NoError(t, err)

	require.NoError(t, s.Close())

	_, afterSize, _ := openFile(f.Name())
	require.True(t, afterSize > beforeSize)
}

// TODO: Is it really needed to return the file pointer?
func openFile(name string) (file *os.File, size int64, err error) {

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	return f, fi.Size(), nil
}
