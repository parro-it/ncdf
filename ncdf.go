package ncdf

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Version []byte

type File struct {
	fd         *os.File
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

func Open(file string) (*File, error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	f := &File{
		fd:      fd,
		Version: make([]byte, 4),
	}
	count := 0

	_, err = fd.Read(f.Version)
	if err != nil {
		fd.Close()
		return nil, err
	}
	count += len((f.Version))

	if err = f.Version.Check(); err != nil {
		fd.Close()
		return nil, err
	}

	if binary.Read(fd, binary.BigEndian, &f.NumRecs) != nil {
		fd.Close()
		return nil, err
	}
	count += 4

	buf, err := readTag(fd, &count)
	if err != nil {
		fd.Close()
		return nil, err
	}

	//NC_DIMENSION = \x00 \x00 \x00 \x0A

	var numdims int32
	if binary.Read(fd, binary.BigEndian, &numdims) != nil {
		fd.Close()
		return nil, err
	}
	count += 4
	//fmt.Println(numrecs, " dimensions")
	f.Dimensions = make([]Dimension, numdims)

	for i := 0; i < int(numdims); i++ {
		var nameLen int32
		if binary.Read(fd, binary.BigEndian, &nameLen) != nil {
			fd.Close()
			return nil, err
		}
		count += 4 + int(nameLen)

		buf = make([]byte, nameLen)
		_, err = fd.Read(buf)
		if err != nil {
			fd.Close()
			return nil, err
		}

		dimName := string(buf)
		//fmt.Println(dimName)
		restCount := 4 - (count % 4)
		if restCount < 4 {
			rest := make([]byte, restCount)
			_, err = fd.Read(rest)
			if err != nil {
				fd.Close()
				return nil, err
			}
			count += restCount
		}

		var dimLen int32
		if binary.Read(fd, binary.BigEndian, &dimLen) != nil {
			fd.Close()
			return nil, err
		}
		count += 4
		f.Dimensions[i] = Dimension{dimName, dimLen}
		//fmt.Printf("%s: %d\n", dimName, dimLen)
	}

	buf, err = readTag(fd, &count)
	if err != nil {
		fd.Close()
		return nil, err
	}
	// NC_ATTRIBUTE = \x00 \x00 \x00 \x0C

	var numgattrs int32
	if binary.Read(fd, binary.BigEndian, &numgattrs) != nil {
		fd.Close()
		return nil, err
	}
	count += 4

	f.Attrs = make([]Attr, numgattrs)

	for i := 0; i < int(numgattrs); i++ {
		var nameLen int32
		if binary.Read(fd, binary.BigEndian, &nameLen) != nil {
			fd.Close()
			return nil, err
		}
		count += 4 + int(nameLen)

		buf = make([]byte, nameLen)
		_, err = fd.Read(buf)
		if err != nil {
			fd.Close()
			return nil, err
		}

		attrName := string(buf)
		restCount := 4 - (count % 4)
		if restCount < 4 {
			rest := make([]byte, restCount)
			_, err = fd.Read(rest)
			if err != nil {
				fd.Close()
				return nil, err
			}
			count += restCount
		}

		var t Type
		if binary.Read(fd, binary.BigEndian, &t) != nil {
			fd.Close()
			return nil, err
		}
		count += 4

		var valLen int32
		if binary.Read(fd, binary.BigEndian, &valLen) != nil {
			fd.Close()
			return nil, err
		}
		count += 4 + int(valLen)

		buf = make([]byte, valLen)
		_, err = fd.Read(buf)
		if err != nil {
			fd.Close()
			return nil, err
		}

		valStr := string(buf)
		restCount = 4 - (count % 4)
		if restCount < 4 {
			rest := make([]byte, restCount)
			_, err = fd.Read(rest)
			if err != nil {
				fd.Close()
				return nil, err
			}
			count += restCount
		}

		f.Attrs[i] = Attr{
			Name: attrName,
			Val:  valStr,
		}

	}
	return f, nil
}

func readTag(fd *os.File, count *int) ([]byte, error) {
	buf := make([]byte, 4)
	_, err := fd.Read(buf)
	if err != nil {
		return nil, err
	}
	*count += 4
	return buf, err
}
