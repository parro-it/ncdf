package ncdf

import (
	"encoding/binary"
	"fmt"
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
}

type Attr struct {
	Name string
	Val  string
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
	Byte  Type = 1 // NC_BYTE = \x00 \x00 \x00 \x01 // 8-bit signed integers
	Char  Type = 2 // NC_CHAR = \x00 \x00 \x00 \x02 // text characters
	Short Type = 3 // NC_SHORT = \x00 \x00 \x00 \x03 // 16-bit signed integers
	Int   Type = 4 // NC_INT = \x00 \x00 \x00 \x04 // 32-bit signed integers
	Float Type = 5 // NC_FLOAT = \x00 \x00 \x00 \x05 // IEEE single precision floats
)

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

	if _, err = readTag(f.fd, &f.count); err != nil {
		return err
	}

	//NC_DIMENSION = \x00 \x00 \x00 \x0A

	//fmt.Printf("%s: %d\n", dimName, dimLen)
	f.Dimensions, err = f.readDimensions()
	if err != nil {
		return err
	}

	_, err = readTag(f.fd, &f.count)
	if err != nil {
		return err
	}
	// NC_ATTRIBUTE = \x00 \x00 \x00 \x0C

	numgattrs, err := readVal[int32](f)
	if err != nil {
		return err
	}

	f.Attrs = make([]Attr, numgattrs)

	for i := 0; i < int(numgattrs); i++ {
		var attrName string

		if attrName, err = readString(f); err != nil {
			return err
		}

		if _, err := readVal[Type](f); err != nil {
			return err
		}

		var valStr string

		if valStr, err = readString(f); err != nil {
			return err
		}

		f.Attrs[i] = Attr{
			Name: attrName,
			Val:  valStr,
		}

	}
	return nil
}

func readListOf[T any](f *File, fn func(f *File) T) ([]T, error) {
	len, err := readVal[int32](f)
	if err != nil {
		return nil, err
	}

	list := make([]T, len)

	for i := int32(0); i < len; i++ {
		list[i] = fn(f)
	}

	return list, nil
}

func (f *File) readDimensions() ([]Dimension, error) {
	numdims, err := readVal[int32](f)
	if err != nil {
		return nil, err
	}

	list := make([]Dimension, numdims)

	for i := 0; i < int(numdims); i++ {

		dimName, err := readString(f)
		if err != nil {
			return nil, err
		}

		dimLen, err := readVal[int32](f)
		if err != nil {
			return nil, err
		}
		list[i] = Dimension{dimName, dimLen}

	}
	return list, nil
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

func readTag(fd *os.File, count *uint64) ([]byte, error) {
	buf := make([]byte, 4)
	_, err := fd.Read(buf)
	if err != nil {
		return nil, err
	}
	*count += 4
	return buf, err
}
