package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var d = Dimension{
	Name: "test",
	Len:  42,
}
var a = Attr{
	Name: "test",
	Val:  int16(42),
	Type: Short,
}
var v = Var{
	Name: "test",
	Attrs: map[string]Attr{
		"tst1": {
			Name: "tst1",
			Val:  int16(42),
			Type: Short,
		},
		"tst2": {
			Name: "tst2",
			Val:  int16(42),
			Type: Short,
		},
	},
	Dimensions: []*Dimension{nil, nil, nil},
	Type:       Short,
	Size:       42,
	Offset:     42,
}

var f = File{
	Version:    [4]byte{},
	NumRecs:    0,
	Dimensions: []Dimension{d, d},
	Attrs: map[string]Attr{
		"a1": a,
		"a2": a,
	},
	Vars: map[string]Var{
		"v1": v,
		"v2": v,
	},
}

func TestFileSize(t *testing.T) {
	var expected = 12 + 12 + 4 + 4 + //dims
		20 + 20 + 4 + 4 + //attrs
		88 + 88 + 4 + 4 + //vars
		4 + 4 // magic + recs
	assert.Equal(t, int32(expected), f.ByteSize())
}

func TestDimSize(t *testing.T) {
	assert.Equal(t, int32(12), d.ByteSize())
}

func TestAttrSize(t *testing.T) {

	assert.Equal(t, int32(20), a.ByteSize())
}

func TestVarSize(t *testing.T) {

	assert.Equal(t, int32(88), v.ByteSize())
}
