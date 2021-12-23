package write

import (
	"bytes"
	"os"
	"testing"

	"github.com/parro-it/ncdf/read"
	"github.com/parro-it/ncdf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var d = types.Dimension{
	Name: "test",
	Len:  42,
}
var a = types.Attr{
	Name: "test",
	Val:  int16(42),
	Type: types.Short,
}
var v = types.Var{
	Name: "test",
	Attrs: map[string]types.Attr{
		"tst1": {
			Name: "tst1",
			Val:  []int16{42},
			Type: types.Short,
		},
		"tst2": {
			Name: "tst2",
			Val:  []int16{42},
			Type: types.Short,
		},
	},
	Dimensions: []*types.Dimension{&d, &d, &d},
	Type:       types.Short,
	Size:       42,
	Offset:     42,
}

var f = types.File{
	Version:    [4]byte{'C', 'D', 'F', 2},
	NumRecs:    0,
	Dimensions: []types.Dimension{d, d},
	Attrs: map[string]types.Attr{
		"a1": {
			Name: "a1",
			Val:  []int16{42},
			Type: types.Short,
		},
		"a2": {
			Name: "a2",
			Val:  []int16{42},
			Type: types.Short,
		},
	},
	Vars: map[string]types.Var{
		"v1": {
			Name:       "v1",
			Attrs:      map[string]types.Attr{},
			Dimensions: []*types.Dimension{&d, &d, &d},
			Type:       types.Short,
			Size:       42,
			Offset:     42,
		},
		"v2": {
			Name: "v2",
			Attrs: map[string]types.Attr{
				"a1": {
					Name: "a1",
					Val:  []int16{42},
					Type: types.Short,
				},
				"a2": {
					Name: "a2",
					Val:  []int16{42},
					Type: types.Short,
				},
			},
			Dimensions: []*types.Dimension{&d, &d, &d},
			Type:       types.Short,
			Size:       42,
			Offset:     42,
		},
	},
}

type BufCloser struct {
	bytes.Buffer
}

func (bc BufCloser) Close() error { return nil }
func (bc BufCloser) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func TestWriteHeader(t *testing.T) {
	fout, err := os.Create("/tmp/prova.nc")
	require.NoError(t, err)
	err = Header(&f, fout)
	require.NoError(t, err)
	require.NoError(t, fout.Close())

	expectedSize := int64(f.ByteSize())
	st, err := os.Stat("/tmp/prova.nc")
	require.NoError(t, err)
	assert.Equal(t, expectedSize, st.Size())

	f2, err := read.HeaderFromDisk(fout.Name())
	//f2.Count = 0
	require.NoError(t, err)
	assert.Equal(t, &f, f2)
	require.NoError(t, err)
}
