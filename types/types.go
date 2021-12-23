package types

import (
	"fmt"
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
	Attrs      map[string]Attr
	Vars       map[string]Var
}

func (f File) ByteSize() int32 {
	var szAttrs int
	szAttrs += 8 // len+tag
	for _, it := range f.Attrs {
		szAttrs += int(it.ByteSize())
	}

	szAttrs += 8 // len+tag
	for _, it := range f.Dimensions {
		szAttrs += int(it.ByteSize())
	}

	szAttrs += 8 // len+tag
	for _, it := range f.Vars {
		szAttrs += int(it.ByteSize())
	}

	return int32(
		szAttrs +
			4 + // numrecs
			4 + // magic & Version
			0)

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
	Attrs      map[string]Attr
	Name       string
	Type       Type
	Size       int32
	Offset     uint64
	//file       *File
}

func (v Var) ByteSize() int32 {
	var szAttrs int32
	szAttrs += 4 + 4 // len+attr tag
	for _, a := range v.Attrs {
		szAttrs += a.ByteSize()
	}

	return 4 + int32(len(v.Dimensions))*4 + // Dimensions
		szAttrs +
		stringByteSize(v.Name) + // Name string
		4 + //Size
		8 + // Offset
		4 // Type

}

// BaseType ...
type BaseType interface {
	byte | int16 | int32 | float32 | float64
}

// Attr ...
type Attr struct {
	Name string
	Val  interface{}
	Type Type
	//file *File
}

// TODO: add support for array values
// TODO: add padding for 32bit alignment
func (a Attr) ByteSize() int32 {
	var sz int32
	switch a.Type {
	case Double:
		sz = 8
	case Short:
		sz = 2
	case Int:
		sz = 4
	case Byte:
		sz = 1
	case Float:
		sz = 4
	case Char:
		sz = 1
	}
	// pad value
	if sz < 4 {
		sz = 4
	}

	return stringByteSize(a.Name) + // Name string
		4 + // Type
		4 + // len
		sz
}

// Dimension ...
type Dimension struct {
	Name string
	Len  int32
	//file *File
}

func stringByteSize(val string) int32 {
	len := len(val)
	rest := 4 - (len % 4)
	if rest != 4 {
		len += rest
	}
	return int32(4 + len)
}

func (d Dimension) ByteSize() int32 {
	return stringByteSize(d.Name) + // Name string
		4 // Len

}

// Type ...
type Type int32

const (
	// Byte is type NC_BYTE = \x00 \x00 \x00 \x01 // 8-bit signed integers
	Byte Type = 1
	// Char is type NC_CHAR = \x00 \x00 \x00 \x02 // text characters
	Char Type = 2
	// Short is type NC_SHORT = \x00 \x00 \x00 \x03 // 16-bit signed integers
	Short Type = 3
	// Int is type NC_INT = \x00 \x00 \x00 \x04 // 32-bit signed integers
	Int Type = 4
	// Float is type NC_FLOAT = \x00 \x00 \x00 \x05 // IEEE single precision floats
	Float Type = 5
	// Double is type NC_DOUBLE = \x00 \x00 \x00 \x06 // IEEE double precision floats
	Double Type = 6
)

/*
func (f *File) Read(data interface{}) error {
	return binary.Read(f.fd, binary.BigEndian, data)
}

// SeekTo ...
func (f *File) SeekTo(n int64) error {
	_, err := f.fd.Seek(n, 0)
	return err
}

// ReadBytes ...
func (f *File) ReadBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := f.fd.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Unlink ...
func (f *File) Unlink() {
	for i, d := range f.Dimensions {
		d.file = nil
		f.Dimensions[i] = d
	}

	for i, v := range f.Vars {
		v.file = nil
		for i, a := range v.Attrs {
			a.file = nil
			v.Attrs[i] = a
		}

		f.Vars[i] = v
	}

	for i, a := range f.Attrs {
		a.file = nil
		f.Attrs[i] = a
	}
}
*/

// NewFile ...
func NewFile() *File {
	return &File{}
}

// NewAttr ...
func NewAttr() Attr {
	return Attr{}
}

// NewDimension ...
func NewDimension() Dimension {
	return Dimension{}
}

// NewVar ...
func NewVar() Var {
	return Var{}
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

/*
// Close ...
func (f *File) Close() error {
	return f.fd.Close()
}
*/
// MarshalJSON ...
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

// MarshalJSON ...
func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, v[3])), nil
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
