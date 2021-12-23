package ordmap_test

import (
	"testing"

	"github.com/parro-it/ncdf/ordmap"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	assert.Equal(t, 0, m.Len())
}

func TestSet(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	assert.False(t, m.Has("a"))
	m.Set("a", 12)
	assert.Equal(t, 1, m.Len())
	assert.True(t, m.Has("a"))
}

func TestGet(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	assert.False(t, m.Has("a"))
	assert.Equal(t, 0, m.Get("a"))
	m.Set("a", 12)
	assert.Equal(t, 1, m.Len())
	v := m.Get("a")
	assert.Equal(t, 12, v)
}

func TestDel(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	m.Del("a")
	m.Set("b", 13)
	m.Set("a", 12)
	assert.Equal(t, 2, m.Len())
	m.Del("a")
	assert.Equal(t, 1, m.Len())
	assert.Equal(t, 0, m.Get("a"))
}

func TestFind(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	assert.Equal(t, -1, m.Find("c"))

	m.Set("b", 13)
	m.Set("a", 12)
	assert.Equal(t, 1, m.Find("a"))
	assert.Equal(t, 0, m.Find("b"))
	assert.Equal(t, -1, m.Find("c"))
}

func TestValues(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	m.Set("b", 13)
	m.Set("a", 12)
	assert.Equal(t, []int{13, 12}, m.Values())
}

func TestItems(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	m.Set("b", 13)
	m.Set("a", 12)
	assert.Equal(t, []ordmap.Item[int, string]{
		{13, "b"},
		{12, "a"},
	}, m.Items())
}
func TestFrom(t *testing.T) {
	m := ordmap.From([]ordmap.Item[int, string]{
		{13, "b"},
		{12, "a"},
	})
	assert.Equal(t, []ordmap.Item[int, string]{
		{13, "b"},
		{12, "a"},
	}, m.Items())
}
func TestKeys(t *testing.T) {
	var m ordmap.OrderedMap[int, string]
	m.Set("b", 13)
	m.Set("a", 12)
	assert.Equal(t, []string{"b", "a"}, m.Keys())
}
