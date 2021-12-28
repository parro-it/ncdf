package types

import (
	"fmt"

	"github.com/parro-it/ncdf/ordmap"
)

// Version contains magic number for netcdf
// in the first 3 bytes, and version info in
// the 4 byte.
type Version [4]byte

// File represent an open netcdf file
// It has an os.File field containing
// the fd of file being read.
type File struct {
	//fd      io.ReadSeekCloser
	//Count   uint64
	Version Version
	NumRecs int32
	//Dimensions    map[string]Dimension
	Dimensions []Dimension
	Attrs      ordmap.OrderedMap[Attr, string]
	Vars       ordmap.OrderedMap[Var, string]
}

// Tag ...
type Tag byte

const (
	// ZeroTag is ZERO tag = \x00 \x00 \x00 \x00 // 32-bit zero
	ZeroTag Tag = 0x00
	// DimensionTag is NC_DIMENSION tag = \x00 \x00 \x00 \x0A // tag for list of dimensions
	DimensionTag Tag = 0x0A
	// VariableTag is NC_VARIABLE tag = \x00 \x00 \x00 \x0B // tag for list of variables
	VariableTag Tag = 0x0B
	// AttributeTag is NC_ATTRIBUTE tag = \x00 \x00 \x00 \x0C // tag for list of attributes
	AttributeTag Tag = 0x0C
)

// Var represents a netcdf variable
type Var struct {
	Dimensions []*Dimension
	Attrs      ordmap.OrderedMap[Attr, string]
	Name       string
	Type       Type
	Size       int32
	Offset     uint64
}

// Attr ...
type Attr struct {
	Name string
	Val  interface{}
	Type Type
	//file *File
}

// Dimension ...
type Dimension struct {
	Name string
	Len  int32
	//file *File
}

// Check ...
func (v Version) Check() error {
	if v[0] != 'C' ||
		v[1] != 'D' ||
		v[2] != 'F' {
		return fmt.Errorf("Invalid magic string %v", v[0:3])
	}
	if v[3] != 1 && v[3] != 2 {
		return fmt.Errorf("Invalid version %d", v[3])
	}
	return nil
}

func (t Tag) String() string {
	switch t {
	case ZeroTag:
		return "ZERO"
	case DimensionTag:
		return "NC_DIMENSION"
	case VariableTag:
		return "NC_VARIABLE"
	case AttributeTag:
		return "NC_ATTRIBUTE"
	}

	return "[UNKNOWN]"
}

type Attrs []Attr

func (attrs Attrs) Map() ordmap.OrderedMap[Attr, string] {
	var res ordmap.OrderedMap[Attr, string]
	for _, a := range attrs {
		res.Set(a.Name, a)
	}
	return res
}

type Vars []Var

func (vars Vars) Map() ordmap.OrderedMap[Var, string] {
	var res ordmap.OrderedMap[Var, string]
	for _, a := range vars {
		res.Set(a.Name, a)
	}
	return res
}
