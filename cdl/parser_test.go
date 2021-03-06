package cdl

import (
	"strings"
	"testing"

	"github.com/parro-it/ncdf/ordmap"
	"github.com/parro-it/ncdf/types"
	"github.com/stretchr/testify/require"
)

func assertParseTo(t *testing.T, code string, expectedFile *types.File, expectedErr string) {
	r := strings.NewReader(code)

	tks, _ := Tokenize(r)

	p := Parser{Tokens: tks}
	f, err := p.Parse()
	if expectedErr == "" {
		// TODO: check for tokenizer errors
		//require.NoError(t, <-tknErr)
		require.NoError(t, err)
	} else {
		require.EqualError(t, err, expectedErr)
	}
	if expectedFile == nil {
		require.Nil(t, f)
	} else {
		require.Equal(t, *expectedFile, *f)
	}

}

func TestAttributes(t *testing.T) {
	assertParseTo(t, "netcdf fname {dimensions: a=1; variables:float pippo (a); pippo:len=15; :lon=45;}", &types.File{
		Dimensions: []types.Dimension{
			{Name: "a", Len: 1},
		},
		Vars: types.Vars{{
			Dimensions: []*types.Dimension{{Name: "a", Len: 1}},
			Name:       "pippo",
			Type:       types.Float,
			Size:       4,
			Attrs: types.Attrs{{
				Name: "len",
				Type: types.Short,
				Val:  []int16{15},
			}}.Map(),
		}}.Map(),
		Attrs: types.Attrs{{
			Name: "lon",
			Type: types.Short,
			Val:  []int16{45},
		}}.Map(),
	}, "")

}

func TestFailures(t *testing.T) {
	assertParseTo(t, "netcdf fname {dimensions: a=10; variables:float pippo (a);}", &types.File{
		Dimensions: []types.Dimension{
			{Name: "a", Len: 10},
		},
		Vars: types.Vars{{
			Dimensions: []*types.Dimension{{Name: "a", Len: 10}},
			Name:       "pippo",
			Type:       types.Float,
			Size:       40,
		}}.Map(),
	}, "")

	assertParseTo(t, "netcdf fname {variables:float pippo (a)}", nil, "Parse failed: unknown dimension name `a`")
	assertParseTo(t, "netcdf fname {variables:float pippo (}", nil, "Parse failed: dimension name expected")
	assertParseTo(t, "netcdf fname {variables:float pippo }", nil, "Parse failed: dimension list expected")
	assertParseTo(t, "netcdf fname {variables:float }", nil, "Parse failed: variable name expected")
	assertParseTo(t, "netcdf fname {variables:wrong}", nil, "Parse failed: : expected")
	assertParseTo(t, "netcdf fname {variables:}", &types.File{Vars: ordmap.OrderedMap[types.Var, string]{}}, "")

	assertParseTo(t, "netcdf fname {dimensions: a=10; b = 20;}", &types.File{Dimensions: []types.Dimension{
		{Name: "a", Len: 10},
		{Name: "b", Len: 20},
	}}, "")

	assertParseTo(t, "netcdf fname {}", &types.File{}, "")
	assertParseTo(t, "ciao", nil, "Parse failed: expected netcdf word")
	assertParseTo(t, "netcdf {", nil, "Parse failed: expected file name")
	assertParseTo(t, "netcdf fname", nil, "Parse failed: expected {")
	assertParseTo(t, "netcdf fname {", nil, "Parse failed: expected }")
	assertParseTo(t, "netcdf fname {dimensions}", nil, "Parse failed: `:` is required after a `dimensions` directive")
	assertParseTo(t, "netcdf fname {dimensions:}", &types.File{Dimensions: []types.Dimension{}}, "")
	assertParseTo(t, "netcdf fname {variables}", nil, "Parse failed: `:` is required after a `variables` directive")

}
