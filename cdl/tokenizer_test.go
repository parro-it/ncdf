package cdl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertTokenizeTo(t *testing.T, name, code string, expected ...interface{}) {
	t.Run(fmt.Sprintf("%s -> %v", code, expected), func(t *testing.T) {
		r := strings.NewReader(code)

		tks, err := Tokenize(r)

		for _, extk := range expected {
			select {
			case tk := <-tks:
				assert.Equal(t, extk, tk)
			case <-time.After(100 * time.Millisecond):
				assert.FailNowf(t, "FAILURE", "Expecting %v but found closed chan.", extk)
			}
		}

		assert.Empty(t, <-tks)
		require.NoError(t, <-err)

	})
}

type failingReader string

func (r failingReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("%s", r)
}

func TestTokenize(t *testing.T) {
	assertTokenizeTo(t, "integers", "42", Token{
		Type:   TkInt,
		Text:   "42",
		NumVal: 42,
	})

	assertTokenizeTo(t, "name", "ciao", Token{
		Type: TkName,
		Text: "ciao",
	})

	assertTokenizeTo(t, "string", `"42"`, Token{
		Type: TkStr,
		Text: "42",
	})

	// TODO: reserved word for netcdf
	assertTokenizeTo(t, "reserved words", `byte short int float double char`, Token{
		Type: TkVarType,
		Text: "byte",
	}, Token{
		Type: TkVarType,
		Text: "short",
	}, Token{
		Type: TkVarType,
		Text: "int",
	}, Token{
		Type: TkVarType,
		Text: "float",
	}, Token{
		Type: TkVarType,
		Text: "double",
	}, Token{
		Type: TkVarType,
		Text: "char",
	})
	assertTokenizeTo(t, "netcdf", `netcdf`, Token{
		Type: TkNetCdf,
		Text: "netcdf",
	})
	assertTokenizeTo(t, "single chars", "{}()=:;,", Token{
		Type: TkCurOpen,
		Text: "{",
	}, Token{
		Type: TkCurClose,
		Text: "}",
	}, Token{
		Type: TkParOpen,
		Text: "(",
	}, Token{
		Type: TkParClose,
		Text: ")",
	}, Token{
		Type: TkEqual,
		Text: "=",
	}, Token{
		Type: TkColon,
		Text: ":",
	}, Token{
		Type: TkSemicolon,
		Text: ";",
	}, Token{
		Type: TkComma,
		Text: ",",
	})
	return
	assertTokenizeTo(t, "empty source", "", Token{
		Type: TkEmpty,
	})

	t.Run("failing reader", func(t *testing.T) {
		tks, err := Tokenize(failingReader("TEST"))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: TEST")

	})

	/*
		t.Run("wrong type", func(t *testing.T) {
			assert.Equal(t, "<Unknown type 99>", TkType(99).String())
			assert.Equal(t, "TkStr", TkStr.String())
		})
	*/

	t.Run("wrong float", func(t *testing.T) {
		tks, err := Tokenize(strings.NewReader("12.12.12"))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: strconv.ParseFloat: parsing \"12.12.12\": invalid syntax")

	})

}
