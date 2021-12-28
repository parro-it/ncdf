package cdl

import (
	"testing"

	"github.com/parro-it/ncdf/types"
	"github.com/stretchr/testify/assert"
)

func TestDimension(t *testing.T) {
	d := types.Dimension{
		Name: "test",
		Len:  (42),
	}
	assert.Equal(t, "test = 42;", CDLDimension(&d))
}

func TestDimensions(t *testing.T) {
	dd := []types.Dimension{
		{Name: "test", Len: (42)},
		{Name: "ciao-mondo", Len: (12)},
	}
	assert.Equal(t, `dimensions:
    test = 42;
    ciao-mondo = 12;
`, dimensionsCDL(dd))
}

func TestTypes(t *testing.T) {
	assert.Equal(t, "float", CDLType(types.Float))
	assert.Equal(t, "double", CDLType(types.Double))
	assert.Equal(t, "int", CDLType(types.Int))
	assert.Equal(t, "short", CDLType(types.Short))
	assert.Equal(t, "char", CDLType(types.Char))
	assert.Equal(t, "byte", CDLType(types.Byte))

}

func TestVar(t *testing.T) {
	v := types.Var{
		Dimensions: []*types.Dimension{{Name: "dim1"}, {Name: "dim2"}},
		Attrs: types.Attrs{{
			Name: "len",
			Val:  42,
			Type: types.Short,
		}, {
			Name: "alt",
			Val:  142,
			Type: types.Int,
		}}.Map(),
		Name: "test",
		Type: types.Double,
	}
	assert.Equal(t, `double test(dim1, dim2);
        test:len = 42;
        test:alt = 142;
`, CDLVar(&v))
}
