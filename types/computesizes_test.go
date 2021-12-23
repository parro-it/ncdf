package types

import (
	"testing"

	"github.com/parro-it/ncdf/ordmap"
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
	Attrs: ordmap.From([]ordmap.Item[Attr, string]{
		{Attr{
			Name: "a1",
			Val:  int16(42),
			Type: Short,
		}, "a1"},
		{Attr{
			Name: "a2",
			Val:  int16(42),
			Type: Short,
		}, "a2"},
	}),
	Vars: ordmap.From([]ordmap.Item[Var, string]{
		{Var{
			Name: "red",
			Attrs: ordmap.From([]ordmap.Item[Attr, string]{
				{Attr{
					Name: "t1",
					Val:  int16(42),
					Type: Short,
				}, "t1"},
				{Attr{
					Name: "t2",
					Val:  int16(42),
					Type: Short,
				}, "t2"},
			}),
			Dimensions: []*Dimension{&dims[0], &dims[1]},
			Type:       Short,
		}, "red"},
		{Var{
			Name: "blu",
			Attrs: ordmap.From([]ordmap.Item[Attr, string]{
				{Attr{
					Name: "t1",
					Val:  int16(42),
					Type: Short,
				}, "t1"},
				{Attr{
					Name: "t2",
					Val:  int16(42),
					Type: Short,
				}, "t2"},
			}),
			Dimensions: []*Dimension{&dims[0], &dims[1]},
			Type:       Short,
		}, "blu"},
	}),
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
