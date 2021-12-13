package ncdf

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type Version [4]byte

type File struct {
	fd         *os.File
	count      uint64
	Version    Version
	NumRecs    int32
	Dimensions []Dimension
	Attrs      []Attr
	Vars       []Var
}

type Var struct {
	Dimensions []*Dimension
	Attrs      []Attr
	Name       string
	Type       Type
	Size       int32
	Offset     uint64
}

type Attr struct {
	Name string
	Val  interface{}
	Type Type
}

type Dimension struct {
	Name string
	Len  int32
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

func (f *File) readHeader() error {
	var err error

	if f.Version, err = readVal[[4]byte](f); err != nil {
		return err
	}

	if err = f.Version.Check(); err != nil {
		return err
	}

	if f.NumRecs, err = readVal[int32](f); err != nil {
		return err
	}

	if f.Dimensions, err = f.readDimensions(); err != nil {
		return err
	}

	if f.Attrs, err = f.readAttributes(); err != nil {
		return err
	}

	if f.Vars, err = f.readVars(); err != nil {
		return err
	}

	return nil
}

type BaseType interface {
	byte | int16 | int32 | float32 | float64
}

func readAttrValue[T BaseType](f *File) ([]T, error) {
	var val T
	var res []T
	nelems, err := readVal[int32](f)
	if err != nil {
		var empty []T
		return empty, err
	}

	for i := int32(0); i < nelems; i++ {
		val, err = readVal[T](f)
		if err != nil {
			var empty []T
			return empty, err
		}
		res = append(res, val)
	}

	restCount := 4 - (f.count % 4)
	if restCount == 4 {
		return res, nil
	}

	rest := make([]byte, restCount)
	_, err = f.fd.Read(rest)
	if err != nil {
		var empty []T
		return empty, err
	}
	f.count += restCount

	return res, nil
}

func (f *File) readAttributes() ([]Attr, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}

	if t == ZeroTag {
		t2, err := readTag(f)
		if err != nil {
			return nil, err
		}

		if t2 != ZeroTag {
			return nil, fmt.Errorf("Expected ZeroTag, got %s", t2.String())
		}

		return []Attr{}, nil
	}

	if t != AttributeTag {
		return nil, fmt.Errorf("Expected AttributeTag, got %s", t.String())
	}

	return readListOf(f, func(f *File) (a Attr, err error) {
		if a.Name, err = readString(f); err != nil {
			return a, err
		}

		if a.Type, err = readVal[Type](f); err != nil {
			return a, err
		}

		if a.Type == Double {
			if a.Val, err = readAttrValue[float64](f); err != nil {
				return a, err
			}

			return
		}

		if a.Type == Short {
			if a.Val, err = readAttrValue[int16](f); err != nil {
				return a, err
			}

			return
		}

		if a.Type == Int {
			if a.Val, err = readAttrValue[int32](f); err != nil {
				return a, err
			}
			return
		}

		if a.Type == Byte {
			if a.Val, err = readAttrValue[byte](f); err != nil {
				return a, err
			}
			return
		}

		if a.Type == Float {
			if a.Val, err = readAttrValue[float32](f); err != nil {
				return a, err
			}

			return
		}

		if a.Type == Char {
			if a.Val, err = readString(f); err != nil {
				return a, err
			}
			return
		}
		log.Panicf("unsupported type %s", a.Type.String())

		return
	})
}

func (f *File) readVars() ([]Var, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}
	if t != VariableTag {
		return nil, fmt.Errorf("Expected VariableTag, got %s", t.String())
	}

	return readListOf(f, func(f *File) (v Var, err error) {
		if v.Name, err = readString(f); err != nil {
			return v, err
		}

		v.Dimensions, err = readListOf(f, func(f *File) (*Dimension, error) {
			id, err := readVal[int32](f)
			if err != nil {
				return nil, err
			}
			return &f.Dimensions[id], nil
		})

		if err != nil {
			return v, err
		}

		if v.Attrs, err = f.readAttributes(); err != nil {
			return v, err
		}
		if v.Type, err = readVal[Type](f); err != nil {
			return v, err
		}

		if v.Size, err = readVal[int32](f); err != nil {
			return v, err
		}

		if v.Offset, err = readVal[uint64](f); err != nil {
			return v, err
		}
		fmt.Println(v.Name, v.Type)
		return
	})
}

func readListOf[T any](f *File, fn func(f *File) (T, error)) (list []T, err error) {
	len, err := readVal[int32](f)
	if err != nil {
		return nil, err
	}

	list = make([]T, len)

	for i := int32(0); i < len; i++ {
		list[i], err = fn(f)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (f *File) readDimensions() ([]Dimension, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}
	if t != DimensionTag {
		return nil, fmt.Errorf("Expected DimensionTag, got %s", t.String())
	}
	return readListOf(f, func(f *File) (d Dimension, err error) {
		if d.Name, err = readString(f); err != nil {
			return d, err
		}

		if d.Len, err = readVal[int32](f); err != nil {
			return d, err
		}

		return
	})
}

func readString(f *File) (string, error) {
	nameLen, err := readVal[int32](f)
	if err != nil {
		return "", err
	}

	buf := make([]byte, nameLen)
	_, err = f.fd.Read(buf)
	if err != nil {
		return "", err
	}
	f.count += uint64(nameLen)

	restCount := 4 - (f.count % 4)
	if restCount == 4 {
		return string(buf), nil
	}

	rest := make([]byte, restCount)
	_, err = f.fd.Read(rest)
	if err != nil {
		return "", err
	}
	f.count += restCount

	return string(buf), nil
}

func readVal[T any](f *File) (T, error) {
	var val T
	if err := binary.Read(f.fd, binary.BigEndian, &val); err != nil {
		var empty T
		return empty, err
	}
	f.count += uint64(unsafe.Sizeof(val))
	return val, nil
}

func Open(file string) (*File, error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	f := &File{
		fd: fd,
	}
	if err := f.readHeader(); err != nil {
		return nil, err
	}

	return f, nil
}

func readTag(f *File) (Tag, error) {
	buf := make([]byte, 4)
	_, err := f.fd.Read(buf)
	if err != nil {
		return ZeroTag, err
	}
	f.count += 4
	return Tag(buf[3]), err
}
