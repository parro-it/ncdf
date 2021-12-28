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

func TestAttr(t *testing.T) {
	aa := map[string]types.Attr{
		"1-a = 42.42;":  {Name: "a", Type: types.Float, Val: 42.42},
		"2-a = 42.42;":  {Name: "a", Type: types.Double, Val: 42.42},
		"3-a = 42;":     {Name: "a", Type: types.Int, Val: 42},
		"4-a = 42;":     {Name: "a", Type: types.Short, Val: 42},
		`5-a = "ciao";`: {Name: "a", Type: types.Char, Val: "ciao"},
		"6-a = 42;":     {Name: "a", Type: types.Byte, Val: 42},
	}

	for expected, actual := range aa {
		assert.Equal(t, expected[2:], CDLAttr(&actual))
	}
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
