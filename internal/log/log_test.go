package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	log_v1 "github.com/wangmingjunabc/proglog/api/v1"
	"google.golang.org/protobuf/proto"
)

func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T, log *Log){
		//"append and read a record success": testAppendRead,
		//"offset out of range error": testOutOfRangeErr,
		"init the existing segments": testInitExisting,
		//"read": testReader,
		//"truncate": testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			c.Segment.MaxStoreBytes = 32
			c.Segment.MaxIndexBytes = 1024
			log, err := NewLog(dir, c)
			require.NoError(t, err)
			fn(t, log)
		})
	}
}

func testAppendRead(t *testing.T, log *Log) {
	record := &log_v1.Record{
		Value: []byte("hello world"),
	}
	off, err := log.Append(record)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, record.Value, read.Value)
}

func testOutOfRangeErr(t *testing.T, log *Log) {
	read, err := log.Read(1)
	require.Nil(t, read)
	require.Error(t, err)
}

func testInitExisting(t *testing.T, log *Log) {
	record := &log_v1.Record{Value: []byte("hello world")}
	for i := 0; i < 3; i++ {
		_, err := log.Append(record)
		require.NoError(t, err)
	}
	require.NoError(t, log.Close())

	offset, err := log.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), offset)

	highestOffset, err := log.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), highestOffset)

	n, err := NewLog(log.Dir, log.Config)
	require.NoError(t, err)

	offset, err = n.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), offset)

	highestOffset, err = n.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), highestOffset)
}

func testReader(t *testing.T, log *Log) {
	record := &log_v1.Record{Value: []byte("hello world")}
	u, err := log.Append(record)
	require.NoError(t, err)
	require.Equal(t, uint64(0), u)

	reader := log.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &log_v1.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, record.Value, read.Value)
}

func testTruncate(t *testing.T, log *Log) {
	record := &log_v1.Record{Value: []byte("hello world")}
	for i := 0; i < 3; i++ {
		_, err := log.Append(record)
		require.NoError(t, err)
	}
	err := log.Truncate(1)
	require.NoError(t, err)

	_, err = log.Read(0)
	require.Error(t, err)
}
