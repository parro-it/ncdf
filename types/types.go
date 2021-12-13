package types

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Version contains magic number for netcdf
// in the first 3 bytes, and version info in
// the 4 byte.
type Version [4]byte

// File represent an open netcdf file
// It has an os.File field containing
// the fd of file being read.
type File struct {
	fd         *os.File
	Count      uint64
	Version    Version
	NumRecs    int32
	Dimensions []Dimension
	Attrs      []Attr
	Vars       []Var
}

func (f *File) Read(data interface{}) error {
	return binary.Read(f.fd, binary.BigEndian, data)
}

func (f *File) ReadBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := f.fd.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Var represents a netcdf variable
type Var struct {
	Dimensions []*Dimension
	Attrs      []Attr
	Name       string
	Type       Type
	Size       int32
	Offset     uint64
	file       *File
}

func (f *File) Unlink() {
	for i, d := range f.Dimensions {
		d.file = nil
		f.Dimensions[i] = d
	}

	for i, v := range f.Vars {
		v.file = nil
		f.Vars[i] = v
	}

	for i, a := range f.Attrs {
		a.file = nil
		f.Attrs[i] = a
	}
}

func (a *Attr) UnlinkFile() {
	a.file = nil
}

func (d *Dimension) UnlinkFile() {
	d.file = nil
}

type Attr struct {
	Name string
	Val  interface{}
	Type Type
	file *File
}

type Dimension struct {
	Name string
	Len  int32
	file *File
}

func NewFile(fd *os.File) *File {
	return &File{fd: fd}
}

func NewAttr(f *File) Attr {
	return Attr{file: f}
}

func NewDimension(f *File) Dimension {
	return Dimension{file: f}
}

func NewVar(f *File) Var {
	return Var{file: f}
}

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

func (f *File) Close() error {
	return f.fd.Close()
}

type Type int32

const (
	Byte   Type = 1 // NC_BYTE = \x00 \x00 \x00 \x01 // 8-bit signed integers
	Char   Type = 2 // NC_CHAR = \x00 \x00 \x00 \x02 // text characters
	Short  Type = 3 // NC_SHORT = \x00 \x00 \x00 \x03 // 16-bit signed integers
	Int    Type = 4 // NC_INT = \x00 \x00 \x00 \x04 // 32-bit signed integers
	Float  Type = 5 // NC_FLOAT = \x00 \x00 \x00 \x05 // IEEE single precision floats
	Double Type = 6 // NC_DOUBLE = \x00 \x00 \x00 \x06 // IEEE double precision floats

)

func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

func (t Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, t[3])), nil
}

func (t Type) String() string {
	switch t {
	case Byte:
		return "NC_BYTE"
	case Char:
		return "NC_CHAR"
	case Short:
		return "NC_SHORT"
	case Int:
		return "NC_INT"
	case Float:
		return "NC_FLOAT"
	case Double:
		return "NC_DOUBLE"
	}

	return fmt.Sprintf("[UNKNOWN:%d]", t)
}

type Tag int32

const (
	ZeroTag      Tag = 0x00 // ZERO = \x00 \x00 \x00 \x00 // 32-bit zero
	DimensionTag Tag = 0x0A // NC_DIMENSION = \x00 \x00 \x00 \x0A // tag for list of dimensions
	VariableTag  Tag = 0x0B // NC_VARIABLE = \x00 \x00 \x00 \x0B // tag for list of variables
	AttributeTag Tag = 0x0C // NC_ATTRIBUTE = \x00 \x00 \x00 \x0C // tag for list of attributes
)

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
