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
		},
	},
}

func TestComputeSizes(t *testing.T) {
	var headSz = 264
	assert.Equal(t, int32(headSz), file.ByteSize())
	file.ComputeSizes()
	assert.Equal(t, uint64(264), file.Vars["red"].Offset)
	assert.Equal(t, int32(12), file.Vars["red"].Size)
	assert.Equal(t, uint64(276), file.Vars["blu"].Offset)
	assert.Equal(t, int32(12), file.Vars["blu"].Size)
}
