package cdl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: add stringifications to token to improve
// test names
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
	t.Run("uncomplete comment", func(t *testing.T) {
		tks, err := Tokenize(strings.NewReader(`/`))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: unexpected char `/`")

	})

	assertTokenizeTo(t, "integers", "42", Token{
		Type:   TkInt,
		Text:   "42",
		NumVal: 42,
	})

	assertTokenizeTo(t, "comments", "// this is a comment\n42 // this is a comment", Token{
		Type:   TkInt,
		Text:   "42",
		NumVal: 42,
	})

	assertTokenizeTo(t, "decimals", "42.15", Token{
		Type:   TkDec,
		Text:   "42.15",
		NumVal: 42.15,
	})

	assertTokenizeTo(t, "name", "ciao", Token{
		Type: TkName,
		Text: "ciao",
	})

	assertTokenizeTo(t, "string", `"42"`, Token{
		Type: TkStr,
		Text: "42",
	})

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

	assertTokenizeTo(t, "dimensions variables data", `dimensions variables data`, Token{
		Type: TkDimensions,
		Text: "dimensions",
	}, Token{
		Type: TkVariables,
		Text: "variables",
	}, Token{
		Type: TkData,
		Text: "data",
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

	assertTokenizeTo(t, "empty source", "", Token{
		Type: TkEmpty,
	})

	t.Run("failing reader", func(t *testing.T) {
		tks, err := Tokenize(failingReader("TEST"))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: TEST")

	})

	t.Run("wrong float", func(t *testing.T) {
		tks, err := Tokenize(strings.NewReader("12.12.12"))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: unexpected dot")

	})

	t.Run("unclosed string", func(t *testing.T) {
		tks, err := Tokenize(strings.NewReader(`"ciao`))
		assert.Empty(t, <-tks)
		e := <-err
		assert.EqualError(t, e, "Tokenization failed: unclosed string")

	})

}
