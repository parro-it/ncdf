package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert.Equal(t, 1, Byte.ScalarSize())
	assert.Equal(t, 1, Char.ScalarSize())
	assert.Equal(t, 2, Short.ScalarSize())
	assert.Equal(t, 4, Int.ScalarSize())
	assert.Equal(t, 4, Float.ScalarSize())
	assert.Equal(t, 8, Double.ScalarSize())
}

func TestAlignForArrayOf(t *testing.T) {
	assert.Equal(t, 3, Byte.AlignForArrayOf(1))
	assert.Equal(t, 3, Char.AlignForArrayOf(1))
	assert.Equal(t, 2, Short.AlignForArrayOf(1))
	assert.Equal(t, 0, Int.AlignForArrayOf(1))
	assert.Equal(t, 0, Float.AlignForArrayOf(1))
	assert.Equal(t, 0, Double.AlignForArrayOf(1))

	assert.Equal(t, -1, Byte.AlignForArrayOf(-1))
	assert.Equal(t, -1, Char.AlignForArrayOf(-1))
	assert.Equal(t, -1, Short.AlignForArrayOf(-1))
	assert.Equal(t, -1, Int.AlignForArrayOf(-1))
	assert.Equal(t, -1, Float.AlignForArrayOf(-1))
	assert.Equal(t, -1, Double.AlignForArrayOf(-1))

	assert.Equal(t, -1, Byte.AlignForArrayOf(0))
	assert.Equal(t, -1, Char.AlignForArrayOf(0))
	assert.Equal(t, -1, Short.AlignForArrayOf(0))
	assert.Equal(t, -1, Int.AlignForArrayOf(0))
	assert.Equal(t, -1, Float.AlignForArrayOf(0))
	assert.Equal(t, -1, Double.AlignForArrayOf(0))

	assert.Equal(t, 2, Byte.AlignForArrayOf(2))
	assert.Equal(t, 2, Char.AlignForArrayOf(2))
	assert.Equal(t, 0, Short.AlignForArrayOf(2))
	assert.Equal(t, 0, Int.AlignForArrayOf(2))
	assert.Equal(t, 0, Float.AlignForArrayOf(2))
	assert.Equal(t, 0, Double.AlignForArrayOf(2))

	assert.Equal(t, 1, Byte.AlignForArrayOf(3))
	assert.Equal(t, 1, Char.AlignForArrayOf(3))
	assert.Equal(t, 2, Short.AlignForArrayOf(3))
	assert.Equal(t, 0, Int.AlignForArrayOf(3))
	assert.Equal(t, 0, Float.AlignForArrayOf(3))
	assert.Equal(t, 0, Double.AlignForArrayOf(3))

	assert.Equal(t, 0, Byte.AlignForArrayOf(4))
	assert.Equal(t, 0, Char.AlignForArrayOf(4))
	assert.Equal(t, 0, Short.AlignForArrayOf(4))
	assert.Equal(t, 0, Int.AlignForArrayOf(4))
	assert.Equal(t, 0, Float.AlignForArrayOf(4))
	assert.Equal(t, 0, Double.AlignForArrayOf(4))
}

func TestArraySize(t *testing.T) {
	assert.Equal(t, 4, Byte.ArraySize(1))
	assert.Equal(t, 4, Char.ArraySize(1))
	assert.Equal(t, 4, Short.ArraySize(1))
	assert.Equal(t, 4, Int.ArraySize(1))
	assert.Equal(t, 4, Float.ArraySize(1))
	assert.Equal(t, 8, Double.ArraySize(1))

	assert.Equal(t, -1, Byte.ArraySize(-1))
	assert.Equal(t, -1, Char.ArraySize(-1))
	assert.Equal(t, -1, Short.ArraySize(-1))
	assert.Equal(t, -1, Int.ArraySize(-1))
	assert.Equal(t, -1, Float.ArraySize(-1))
	assert.Equal(t, -1, Double.ArraySize(-1))

	assert.Equal(t, -1, Byte.ArraySize(0))
	assert.Equal(t, -1, Char.ArraySize(0))
	assert.Equal(t, -1, Short.ArraySize(0))
	assert.Equal(t, -1, Int.ArraySize(0))
	assert.Equal(t, -1, Float.ArraySize(0))
	assert.Equal(t, -1, Double.ArraySize(0))

	assert.Equal(t, 4, Byte.ArraySize(2))
	assert.Equal(t, 4, Char.ArraySize(2))
	assert.Equal(t, 4, Short.ArraySize(2))
	assert.Equal(t, 8, Int.ArraySize(2))
	assert.Equal(t, 8, Float.ArraySize(2))
	assert.Equal(t, 16, Double.ArraySize(2))

	assert.Equal(t, 4, Byte.ArraySize(3))
	assert.Equal(t, 4, Char.ArraySize(3))
	assert.Equal(t, 8, Short.ArraySize(3))
	assert.Equal(t, 12, Int.ArraySize(3))
	assert.Equal(t, 12, Float.ArraySize(3))
	assert.Equal(t, 24, Double.ArraySize(3))

	assert.Equal(t, 4, Byte.ArraySize(4))
	assert.Equal(t, 4, Char.ArraySize(4))
	assert.Equal(t, 8, Short.ArraySize(4))
	assert.Equal(t, 16, Int.ArraySize(4))
	assert.Equal(t, 16, Float.ArraySize(4))
	assert.Equal(t, 32, Double.ArraySize(4))

	assert.Equal(t, 8, Byte.ArraySize(5))
	assert.Equal(t, 8, Char.ArraySize(5))
	assert.Equal(t, 12, Short.ArraySize(5))
	assert.Equal(t, 20, Int.ArraySize(5))
	assert.Equal(t, 20, Float.ArraySize(5))
	assert.Equal(t, 40, Double.ArraySize(5))
}
