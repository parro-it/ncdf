package ncdf

import (
	"os"
	"testing"

	"github.com/parro-it/ncdf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	f, err := Open("fixtures/exampl2.nc")
	require.NoError(t, err)
	assert.NotNil(t, f)
	//f.Close()
}

func TestVarData(t *testing.T) {
	f, err := Open("fixtures/exampl2.nc")
	require.NoError(t, err)
	require.NotNil(t, f)
	//defer f.Close()

	t2 := f.Vars["T2"]
	dim := 1
	for _, d := range t2.Dimensions {
		if d.Len == 0 {
			continue
		}
		dim *= int(d.Len)
	}

	fd, err := os.Open("fixtures/exampl2.nc")
	require.NoError(t, err)
	defer fd.Close()

	values, err := VarData[float32](t2, fd)
	require.NoError(t, err)
	require.Equal(t, dim*4, len(values))

	require.Equal(t, []float32{275.82614, 275.82596, 275.98535, 275.60385, 275.5405, 275.5339, 275.5177, 275.3091, 275.31165, 275.27158}, values[0:10])

}

func TestCheck(t *testing.T) {
	t.Run("bad magic string", func(t *testing.T) {
		f := &types.File{
			Version: [4]byte{1, 2, 3, 4},
		}
		assert.EqualError(t, f.Version.Check(), "Invalid magic string [1 2 3]")
	})

	t.Run("bad version", func(t *testing.T) {
		f := &types.File{
			Version: [4]byte{'C', 'D', 'F', 4},
		}
		assert.EqualError(t, f.Version.Check(), "Invalid version 4")
	})

	t.Run("Correct", func(t *testing.T) {
		f := &types.File{
			Version: [4]byte{'C', 'D', 'F', 1},
		}
		assert.NoError(t, f.Version.Check())
	})
	t.Run("NumRecs", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		assert.NoError(t, err)
		require.NotNil(t, f)
		assert.Equal(t, int32(1), f.NumRecs)
		//f.Close()

	})

	t.Run("Dimensions", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		//f.Close()
		assert.NoError(t, err)
		assert.Equal(t, []types.Dimension{
			{Name: "Time", Len: 0},
			{Name: "DateStrLen", Len: 19},
			{Name: "west_east", Len: 429},
			{Name: "south_north", Len: 468},
			{Name: "num_press_levels_stag", Len: 11},
		}, f.Dimensions)

	})

	t.Run("Variables", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		assert.NoError(t, err)
		f.Unlink()
		time := f.Vars["Times"]
		assert.Equal(t, "Time", time.Dimensions[0].Name)
		time.Dimensions = make([]*types.Dimension, 0)

		assert.Equal(t, types.Var{
			Dimensions: []*types.Dimension{},
			Attrs:      map[string]types.Attr{},
			Name:       "Times",
			Type:       2,
			Size:       20,
			Offset:     0x2f5c,
		}, time)

	})

	t.Run("Global attributes", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		assert.NoError(t, err)
		f.Unlink()
		assert.Equal(t, types.Attr{
			Name: "TITLE",
			Val:  " OUTPUT FROM WRF V3.8.1 MODEL",
			Type: types.Char,
		}, f.Attrs["TITLE"])

	})

}
