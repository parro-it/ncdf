package cdl

import (
	"strings"
	"testing"

	"github.com/parro-it/ncdf/types"
	"github.com/stretchr/testify/require"
)

func assertParseTo(t *testing.T, code string, expectedFile *types.File, expectedErr string) {
	r := strings.NewReader(code)

	tks, tknErr := Tokenize(r)

	p := Parser{Tokens: tks}
	f, err := p.Parse()
	if expectedErr == "" {
		require.NoError(t, <-tknErr)
		require.NoError(t, err)
	} else {
		require.EqualError(t, err, expectedErr)
	}
	if expectedFile == nil {
		require.Nil(t, f)
	} else {
		require.Equal(t, *f, *expectedFile)
	}

}

func TestFailures(t *testing.T) {
	assertParseTo(t, "netcdf fname {}", &types.File{}, "")

	assertParseTo(t, "ciao", nil, "Parse failed: expected netcdf word")
	assertParseTo(t, "netcdf {", nil, "Parse failed: expected file name")
	assertParseTo(t, "netcdf fname", nil, "Parse failed: expected {")
	assertParseTo(t, "netcdf fname {", nil, "Parse failed: expected }")
}
