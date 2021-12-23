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
	Attrs: Attrs{{
		Name: "a1",
		Val:  int16(42),
		Type: Short,
	}, {
		Name: "a2",
		Val:  int16(42),
		Type: Short,
	}}.Map(),
	Vars: Vars{{
		Name: "red",
		Attrs: Attrs{{
			Name: "t1",
			Val:  int16(42),
			Type: Short,
		}, {
			Name: "t2",
			Val:  int16(42),
			Type: Short,
		}}.Map(),

		Dimensions: []*Dimension{&dims[0], &dims[1]},
		Type:       Short,
	}, {
		Name: "blu",
		Attrs: Attrs{{
			Name: "t1",
			Val:  int16(42),
			Type: Short,
		}, {
			Name: "t2",
			Val:  int16(42),
			Type: Short,
		}}.Map(),
		Dimensions: []*Dimension{&dims[0], &dims[1]},
		Type:       Short,
	}}.Map(),
}

func TestComputeSizes(t *testing.T) {
	var headSz = 264
	assert.Equal(t, int32(headSz), file.ByteSize())
	file.ComputeSizes()
	assert.Equal(t, uint64(264), file.Vars.Get("red").Offset)
	assert.Equal(t, int32(12), file.Vars.Get("red").Size)
	assert.Equal(t, uint64(276), file.Vars.Get("blu").Offset)
	assert.Equal(t, int32(12), file.Vars.Get("blu").Size)
}
