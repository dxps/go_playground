package log

import (
	"io/ioutil"
	"os"
	"testing"

	api "devisions.org/go-dist-svcs/log/api/v1"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {

	for scenario, fn := range map[string]func(t *testing.T, log *Log){
		"append and read a record succeeds": testAppendRead,
		"offset out of range error":         testOutOfRangeErr,
		"init with existing segments":       testInitExisting,
		"reader":                            testReader,
		"truncate":                          testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "log-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := NewLog(dir, c)
			require.NoError(t, err)
			fn(t, log)
		})
	}
}

func testAppendRead(t *testing.T, log *Log) {

	rec := &api.Record{Value: []byte("hello world")}
	off, err := log.Append(rec)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, rec, read)
}

func testOutOfRangeErr(t *testing.T, log *Log) {

	read, err := log.Read(1)
	require.Nil(t, read)
	require.Error(t, err)

}

func testInitExisting(t *testing.T, log *Log) {

	// Preparing the log for testing.
	rec := &api.Record{Value: []byte("hello world")}
	for i := 0; i < 3; i++ {
		_, err := log.Append(rec)
		require.NoError(t, err)
	}
	require.NoError(t, log.Close())

	off, err := log.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	off, err = log.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

	// Initing a log based on existing data.
	newLog, err := NewLog(log.Dir, log.Config)
	require.NoError(t, err)

	off, err = newLog.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	off, err = newLog.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

}

// testReader test that we can read the full, raw log as being stored on disk.
// This reading mechanism is used when taking snapshots.
func testReader(t *testing.T, log *Log) {

	rec := &api.Record{Value: []byte("hello world")}
	off, err := log.Append(rec)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	reader := log.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, rec, read)
}

func testTruncate(t *testing.T, log *Log) {

	// Preparing the log for testing.
	rec := &api.Record{Value: []byte("hello world")}
	for i := 0; i < 3; i++ {
		_, err := log.Append(rec)
		require.NoError(t, err)
	}

	err := log.Truncate(1)
	require.NoError(t, err)

	_, err = log.Read(0)
	require.Error(t, err)
}
