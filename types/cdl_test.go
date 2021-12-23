package types

import (
	"testing"

	"github.com/parro-it/ncdf/ordmap"
	"github.com/stretchr/testify/assert"
)

func TestDimension(t *testing.T) {
	d := Dimension{
		Name: "test",
		Len:  (42),
	}
	assert.Equal(t, "test = 42;", d.CDL())
}

func TestDimensions(t *testing.T) {
	dd := []Dimension{
		{Name: "test", Len: (42)},
		{Name: "ciao-mondo", Len: (12)},
	}
	assert.Equal(t, `dimensions:
    test = 42;
    ciao-mondo = 12;
`, dimensionsCDL(dd))
}

func TestTypes(t *testing.T) {
	assert.Equal(t, "float", Float.CDL())
	assert.Equal(t, "double", Double.CDL())
	assert.Equal(t, "int", Int.CDL())
	assert.Equal(t, "short", Short.CDL())
	assert.Equal(t, "char", Char.CDL())
	assert.Equal(t, "byte", Byte.CDL())

}

func TestVar(t *testing.T) {
	v := Var{
		Dimensions: []*Dimension{{Name: "dim1"}, {Name: "dim2"}},
		Attrs: ordmap.From([]ordmap.Item[Attr, string]{
			{Attr{
				Name: "len",
				Val:  42,
				Type: Short,
			}, "len"},
			{Attr{
				Name: "alt",
				Val:  142,
				Type: Int,
			}, "alt"},
		}),
		Name: "test",
		Type: Double,
	}
	assert.Equal(t, `double test(dim1, dim2);
        test:len = 42;
        test:alt = 142;
`, v.CDL())
}
