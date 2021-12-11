package ncdf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	f, err := Open("fixtures/example.nc")
	assert.NoError(t, err)
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
		f, err := Open("fixtures/example.nc")
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, int32(0), f.NumRecs)
		f.Close()

	})

	t.Run("Dimensions", func(t *testing.T) {
		f, err := Open("fixtures/example.nc")
		require.NotNil(t, f)
		defer f.Close()
		assert.NoError(t, err)
		assert.Equal(t, []Dimension{
			{"longitude", 3600},
			{"latitude", 1801},
			{"time", 24},
		}, f.Dimensions)

	})

	t.Run("Global attributes", func(t *testing.T) {
		f, err := Open("fixtures/example.nc")
		require.NotNil(t, f)
		defer f.Close()
		assert.NoError(t, err)
		assert.Equal(t, []Attr{
			{"Conventions", "CF-1.6"},
			{"history", "2020-03-20 11:41:00 GMT by grib_to_netcdf-2.16.0: /opt/ecmwf/eccodes/bin/grib_to_netcdf -S param -o /cache/data4/adaptor.mars.internal-1584704431.3029954-14866-24-a3017812-b06b-4c9d-aee6-e5a74bbbfbc9.nc /cache/tmp/a3017812-b06b-4c9d-aee6-e5a74bbbfbc9-adaptor.mars.internal-1584704431.3035216-14866-8-tmp.grib"},
		}, f.Attrs)

	})

}
