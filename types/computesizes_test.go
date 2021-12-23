package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dims = []Dimension{{
	Name: "x",
	Len:  2,
}, {
	Name: "y",
	Len:  3,
}}

var file = File{
	Version:    [4]byte{},
	NumRecs:    0,
	Dimensions: dims,
	Attrs: map[string]Attr{
		"a1": {
			Name: "a1",
			Val:  int16(42),
			Type: Short,
		},
		"a2": {
			Name: "a2",
			Val:  int16(42),
			Type: Short,
		},
	},
	Vars: map[string]Var{
		"red": {
			Name: "red",
			Attrs: map[string]Attr{
				"t1": {
					Name: "t1",
					Val:  int16(42),
					Type: Short,
				},
				"t2": {
					Name: "t2",
					Val:  int16(42),
					Type: Short,
				},
			},
			Dimensions: []*Dimension{&dims[0], &dims[1]},
			Type:       Short,
			Size:       42,
			Offset:     42,
		},
		"blu": {
			Name: "blu",
			Attrs: map[string]Attr{
				"t1": {
					Name: "t1",
					Val:  int16(42),
					Type: Short,
				},
				"t2": {
					Name: "t2",
					Val:  int16(42),
					Type: Short,
				},
			},
			Dimensions: []*Dimension{&dims[0], &dims[1]},
			Type:       Short,
			Size:       42,
			Offset:     42,
		},
	},
}

func TestComputeSizes(t *testing.T) {
	var expected = 12 + 12 + 4 /*len*/ + 4 /*tag*/ + //dims
		20 + 20 + 4 + 4 + //attrs
		84 + 84 + 4 + 4 + //vars
		4 + 4 // magic + recs
	assert.Equal(t, int32(expected), file.ByteSize())

}
