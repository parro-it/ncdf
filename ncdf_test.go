package ncdf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	i, err := os.Stat("fixtures/exampl2.nc")
	require.NoError(t, err)
	fmt.Println(i)
	f, err := Open("fixtures/exampl2.nc")
	require.NoError(t, err)
	assert.NotNil(t, f)
	f.Close()
}

func TestCheck(t *testing.T) {
	t.Run("bad magic string", func(t *testing.T) {
		f := &File{
			Version: [4]byte{1, 2, 3, 4},
		}
		assert.EqualError(t, f.Version.Check(), "Invalid magic string [1 2 3]")
	})

	t.Run("bad version", func(t *testing.T) {
		f := &File{
			Version: [4]byte{'C', 'D', 'F', 4},
		}
		assert.EqualError(t, f.Version.Check(), "Invalid version 4")
	})

	t.Run("Correct", func(t *testing.T) {
		f := &File{
			Version: [4]byte{'C', 'D', 'F', 1},
		}
		assert.NoError(t, f.Version.Check())
	})

	t.Run("NumRecs", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		assert.NoError(t, err)
		require.NotNil(t, f)
		assert.Equal(t, int32(1), f.NumRecs)
		f.Close()

	})

	t.Run("Dimensions", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		defer f.Close()
		assert.NoError(t, err)
		assert.Equal(t, []Dimension{
			{"Time", 0},
			{"DateStrLen", 19},
			{"west_east", 429},
			{"south_north", 468},
			{"num_press_levels_stag", 11},
		}, f.Dimensions)

	})

	t.Run("Variables", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		defer f.Close()
		assert.NoError(t, err)
		assert.Equal(t, "Time", f.Vars[0].Dimensions[0].Name)
		f.Vars[0].Dimensions = make([]*Dimension, 0)
		assert.Equal(t, Var{
			Dimensions: []*Dimension{},
			Attrs:      []Attr{},
			Name:       "Times",
			Type:       2,
			Size:       20,
			Offset:     0x2f5c,
		}, f.Vars[0])

	})

	t.Run("Global attributes", func(t *testing.T) {
		f, err := Open("fixtures/exampl2.nc")
		require.NotNil(t, f)
		defer f.Close()
		assert.NoError(t, err)
		assert.Equal(t, Attr{"TITLE", " OUTPUT FROM WRF V3.8.1 MODEL", Char}, f.Attrs[0])

	})

}
