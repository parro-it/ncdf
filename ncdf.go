package ncdf

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/parro-it/ncdf/types"
)

func WriteHeader(f *types.File, w io.Writer) error {
	// magic + version
	bytes := []byte{'C', 'D', 'F', 2}
	if err := binary.Write(w, binary.BigEndian, bytes); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, f.NumRecs); err != nil {
		return err
	}

	// dimensions
	var tag [4]byte
	tag[3] = byte(types.DimensionTag)
	if err := binary.Write(w, binary.BigEndian, tag); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, int32(len(f.Dimensions))); err != nil {
		return err
	}

	for _, d := range f.Dimensions {
		if err := writeDimension(d, w); err != nil {
			return err
		}
	}

	// attrs

	if err := writeAttrs(w, f.Attrs); err != nil {
		return err
	}

	// vars
	tag[3] = byte(types.VariableTag)
	if err := binary.Write(w, binary.BigEndian, tag); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, int32(len(f.Vars))); err != nil {
		return err
	}

	for _, v := range f.Vars {
		if err := writeVar(f, v, w); err != nil {
			return err
		}
	}

	return nil
}

func writeAttrs(w io.Writer, attrs map[string]types.Attr) error {
	var tag [4]byte
	tag[3] = byte(types.AttributeTag)
	if err := binary.Write(w, binary.BigEndian, tag); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, int32(len(attrs))); err != nil {
		return err
	}

	for _, a := range attrs {
		if err := writeAttr(a, w); err != nil {
			return err
		}
	}
	return nil
}

func writeAttr(a types.Attr, w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(a.Name))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, []byte(a.Name)); err != nil {
		return err
	}
	rest := 4 - (len(a.Name) % 4)
	if rest != 4 {
		buf := make([]byte, rest)
		if err := binary.Write(w, binary.BigEndian, buf); err != nil {
			return err
		}
	}
	if err := binary.Write(w, binary.BigEndian, a.Type); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, int32(1)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, a.Val); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, int16(0)); err != nil {
		return err
	}

	return nil
}

func writeVar(f *types.File, v types.Var, w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(v.Name))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, []byte(v.Name)); err != nil {
		return err
	}
	rest := 4 - (len(v.Name) % 4)
	if rest != 4 {
		buf := make([]byte, rest)
		if err := binary.Write(w, binary.BigEndian, buf); err != nil {
			return err
		}
	}
	if err := binary.Write(w, binary.BigEndian, int32(len(v.Dimensions))); err != nil {
		return err
	}

	findDim := func(d *types.Dimension) int32 {
		for idx, dt := range f.Dimensions {
			if d.Name == dt.Name {
				return int32(idx)
			}

		}
		return -1
	}

	for _, d := range v.Dimensions {
		if err := binary.Write(w, binary.BigEndian, findDim(d)); err != nil {
			return err
		}
	}

	if err := writeAttrs(w, v.Attrs); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, v.Type); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, v.Size); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, v.Offset); err != nil {
		return err
	}

	return nil
}

func writeDimension(d types.Dimension, w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(d.Name))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, []byte(d.Name)); err != nil {
		return err
	}

	rest := 4 - (len(d.Name) % 4)
	if rest != 4 {
		buf := make([]byte, rest)
		if err := binary.Write(w, binary.BigEndian, buf); err != nil {
			return err
		}
	}

	if err := binary.Write(w, binary.BigEndian, d.Len); err != nil {
		return err
	}
	return nil
}

func ReadHeader(fd io.ReadSeeker) (*types.File, error) {
	f := &types.File{}
	var err error

	if f.Version, err = readSingleValue[[4]byte](fd); err != nil {
		return nil, err
	}

	if err = f.Version.Check(); err != nil {
		return nil, err
	}

	if f.NumRecs, err = readSingleValue[int32](fd); err != nil {
		return nil, err
	}

	if f.Dimensions, err = readDimensions(fd); err != nil {
		return nil, err
	}

	if f.Attrs, err = readAttributes(fd); err != nil {
		return nil, err
	}

	if f.Vars, err = readVars(f.Dimensions, fd); err != nil {
		return nil, err
	}

	return f, nil
}

func readDimensions(fd io.ReadSeeker) ([]types.Dimension, error) {
	t, err := readTag(fd)
	if err != nil {
		return nil, err
	}
	if t == types.ZeroTag {
		_, err := sectionNotPresent[types.Dimension](fd)
		return nil, err
	}
	if t != types.DimensionTag {
		return nil, fmt.Errorf("Expected DimensionTag, got %s", t.String())
	}
	lst, err := readListOfObjects(fd, func() (d types.Dimension, err error) {
		d = types.NewDimension()

		if d.Name, err = readString(fd); err != nil {
			return d, err
		}

		if d.Len, err = readSingleValue[int32](fd); err != nil {
			return d, err
		}

		return
	})
	if err != nil {
		return nil, err
	}

	return lst, nil
}

func sectionNotPresent[T any](fd io.ReadSeeker) (map[string]T, error) {
	t2, err := readTag(fd)
	if err != nil {
		return nil, err
	}

	if t2 != types.ZeroTag {
		return nil, fmt.Errorf("Expected ZeroTag, got %s", t2.String())
	}

	return map[string]T{}, nil
}

func readAttributes(fd io.ReadSeeker) (map[string]types.Attr, error) {
	t, err := readTag(fd)
	if err != nil {
		return nil, err
	}

	if t == types.ZeroTag {
		return sectionNotPresent[types.Attr](fd)
	}

	if t != types.AttributeTag {
		return nil, fmt.Errorf("Expected AttributeTag, got %s", t.String())
	}

	lst, err := readListOfObjects(fd, func() (a types.Attr, err error) {
		a = types.NewAttr()

		if a.Name, err = readString(fd); err != nil {
			return
		}

		if a.Type, err = readSingleValue[types.Type](fd); err != nil {
			return
		}

		if a.Val, err = readAttributeValue(a, fd); err != nil {
			return
		}

		return
	})
	if err != nil {
		return nil, err
	}
	res := map[string]types.Attr{}
	for _, d := range lst {
		res[d.Name] = d
	}
	return res, nil
}

func readVars(dims []types.Dimension, fd io.ReadSeeker) (map[string]types.Var, error) {
	t, err := readTag(fd)
	if err != nil {
		return nil, err
	}
	if t == types.ZeroTag {
		return sectionNotPresent[types.Var](fd)
	}

	if t != types.VariableTag {
		return nil, fmt.Errorf("Expected VariableTag, got %s", t.String())
	}

	lst, err := readListOfObjects(fd, func() (v types.Var, err error) {
		v = types.NewVar()

		if v.Name, err = readString(fd); err != nil {
			return v, err
		}

		v.Dimensions, err = readListOfObjects(fd, func() (*types.Dimension, error) {
			id, err := readSingleValue[int32](fd)
			if err != nil {
				return nil, err
			}
			return &dims[id], nil
		})

		if err != nil {
			return v, err
		}

		if v.Attrs, err = readAttributes(fd); err != nil {
			return v, err
		}
		if v.Type, err = readSingleValue[types.Type](fd); err != nil {
			return v, err
		}

		if v.Size, err = readSingleValue[int32](fd); err != nil {
			return v, err
		}

		if v.Offset, err = readSingleValue[uint64](fd); err != nil {
			return v, err
		}
		return
	})

	if err != nil {
		return nil, err
	}
	res := map[string]types.Var{}
	for _, d := range lst {
		res[d.Name] = d
	}
	return res, nil
}

func readListOfValues[T types.BaseType](fd io.ReadSeeker) ([]T, error) {
	nelems, err := readSingleValue[int32](fd)
	if err != nil {
		var empty []T
		return empty, err
	}

	var res []T
	var val T

	for i := int32(0); i < nelems; i++ {
		val, err = readSingleValue[T](fd)
		if err != nil {
			var empty []T
			return empty, err
		}
		res = append(res, val)
	}
	if unsafe.Sizeof(val) < 4 {
		sz := int64(unsafe.Sizeof(val) * uintptr(nelems))
		restCount := 4 - (sz % 4)
		if restCount == 4 {
			return res, nil
		}

		_, err = fd.Seek(restCount, io.SeekCurrent)
		if err != nil {
			var empty []T
			return empty, err
		}
	}

	return res, nil
}

// TODO: add support for multiple values
func readAttributeValue(a types.Attr, fd io.ReadSeeker) (interface{}, error) {
	t := a.Type
	if t == types.Double {
		return readListOfValues[float64](fd)
	}

	if t == types.Short {
		return readListOfValues[int16](fd)
	}

	if t == types.Int {
		return readListOfValues[int32](fd)
	}

	if t == types.Byte {
		return readListOfValues[byte](fd)
	}

	if t == types.Float {
		return readListOfValues[float32](fd)
	}

	if t == types.Char {
		return readString(fd)
	}

	return nil, fmt.Errorf("Unsupported type <%s>", t)
}

func readListOfObjects[T any](fd io.ReadSeeker, fn func() (T, error)) (list []T, err error) {
	len, err := readSingleValue[int32](fd)
	if err != nil {
		return nil, err
	}

	list = make([]T, len)

	for i := int32(0); i < len; i++ {
		list[i], err = fn()
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func readString(fd io.ReadSeeker) (string, error) {
	v, err := readListOfValues[byte](fd)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func readSingleValue[T any](fd io.ReadSeeker) (T, error) {
	var val T
	if err := binary.Read(fd, binary.BigEndian, &val); err != nil {
		var empty T
		return empty, err
	}
	return val, nil
}

func readTag(fd io.ReadSeeker) (types.Tag, error) {
	var buf [4]byte
	if err := binary.Read(fd, binary.BigEndian, &buf); err != nil {
		return types.ZeroTag, err
	}
	return types.Tag(buf[3]), nil
}

// Open ...
func Open(file string) (*types.File, error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	f, err := ReadHeader(fd)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// VarData ...
func VarData[T types.BaseType](v types.Var, fd io.ReadSeeker) ([]T, error) {
	if _, err := fd.Seek(int64(v.Offset), io.SeekStart); err != nil {
		return nil, err
	}
	data := make([]T, v.Size)
	if err := binary.Read(fd, binary.BigEndian, &data); err != nil {
		return nil, err
	}
	return data, nil
}
