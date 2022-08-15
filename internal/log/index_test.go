package log

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	idx, err := newIndex(f, c)
	require.NoError(t, err)
	_, _, err = idx.Read(-1)
	require.Equal(t, io.EOF, err)
	require.Equal(t, f.Name(), idx.Name())

	entries := []struct {
		Offset uint32
		Pos    uint64
	}{{
		Offset: 0,
		Pos:    0,
	},
		{
			Offset: 1,
			Pos:    10,
		},
	}
	for _, want := range entries {
		err = idx.Write(want.Offset, want.Pos)
		require.NoError(t, err)

		_, pos, err := idx.Read(int64(want.Offset))
		require.NoError(t, err)
		require.Equal(t, want.Pos, pos)
	}
	_, _, err = idx.Read(-1)
	require.NoError(t, err)
	_ = idx.Close()

	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	idx, err = newIndex(f, c)
	require.NoError(t, err)
	off, pos, err := idx.Read(-1)
	require.NoError(t, err)
	require.Equal(t, uint32(1), off)
	require.Equal(t, entries[1].Pos, pos)
}
