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

func readHeader(f *types.File, fd io.ReadSeeker) error {
	var err error

	if f.Version, err = readVal[[4]byte](f, fd); err != nil {
		return err
	}

	if err = f.Version.Check(); err != nil {
		return err
	}

	if f.NumRecs, err = readVal[int32](f, fd); err != nil {
		return err
	}

	if f.Dimensions, err = readDimensions(f, fd); err != nil {
		return err
	}

	if f.Attrs, err = readAttributes(f, fd); err != nil {
		return err
	}

	if f.Vars, err = readVars(f, fd); err != nil {
		return err
	}

	return nil
}

func readDimensions(f *types.File, fd io.ReadSeeker) ([]types.Dimension, error) {
	t, err := readTag(f, fd)
	if err != nil {
		return nil, err
	}
	if t == types.ZeroTag {
		_, err := sectionNotPresent[types.Dimension](f, fd)
		return nil, err
	}
	if t != types.DimensionTag {
		return nil, fmt.Errorf("Expected DimensionTag, got %s", t.String())
	}
	lst, err := readListOf(f, fd, func(f *types.File) (d types.Dimension, err error) {
		d = types.NewDimension(f)

		if d.Name, err = readString(f, fd); err != nil {
			return d, err
		}

		if d.Len, err = readVal[int32](f, fd); err != nil {
			return d, err
		}

		return
	})
	if err != nil {
		return nil, err
	}

	return lst, nil
}

func sectionNotPresent[T any](f *types.File, fd io.ReadSeeker) (map[string]T, error) {
	t2, err := readTag(f, fd)
	if err != nil {
		return nil, err
	}

	if t2 != types.ZeroTag {
		return nil, fmt.Errorf("Expected ZeroTag, got %s", t2.String())
	}

	return map[string]T{}, nil
}

func readAttributes(f *types.File, fd io.ReadSeeker) (map[string]types.Attr, error) {
	t, err := readTag(f, fd)
	if err != nil {
		return nil, err
	}

	if t == types.ZeroTag {
		return sectionNotPresent[types.Attr](f, fd)
	}

	if t != types.AttributeTag {
		return nil, fmt.Errorf("Expected AttributeTag, got %s", t.String())
	}

	lst, err := readListOf(f, fd, func(f *types.File) (a types.Attr, err error) {
		a = types.NewAttr(f)

		if a.Name, err = readString(f, fd); err != nil {
			return
		}

		if a.Type, err = readVal[types.Type](f, fd); err != nil {
			return
		}

		if a.Val, err = readValue(a, f, fd); err != nil {
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

func readVars(f *types.File, fd io.ReadSeeker) (map[string]types.Var, error) {
	t, err := readTag(f, fd)
	if err != nil {
		return nil, err
	}
	if t == types.ZeroTag {
		return sectionNotPresent[types.Var](f, fd)
	}

	if t != types.VariableTag {
		return nil, fmt.Errorf("Expected VariableTag, got %s", t.String())
	}

	lst, err := readListOf(f, fd, func(f *types.File) (v types.Var, err error) {
		v = types.NewVar(f)

		if v.Name, err = readString(f, fd); err != nil {
			return v, err
		}

		v.Dimensions, err = readListOf(f, fd, func(f *types.File) (*types.Dimension, error) {
			id, err := readVal[int32](f, fd)
			if err != nil {
				return nil, err
			}
			return &f.Dimensions[id], nil
		})

		if err != nil {
			return v, err
		}

		if v.Attrs, err = readAttributes(f, fd); err != nil {
			return v, err
		}
		if v.Type, err = readVal[types.Type](f, fd); err != nil {
			return v, err
		}

		if v.Size, err = readVal[int32](f, fd); err != nil {
			return v, err
		}

		if v.Offset, err = readVal[uint64](f, fd); err != nil {
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

func readAttrValue[T types.BaseType](f *types.File, fd io.ReadSeeker) ([]T, error) {
	var val T
	var res []T
	nelems, err := readVal[int32](f, fd)
	if err != nil {
		var empty []T
		return empty, err
	}

	for i := int32(0); i < nelems; i++ {
		val, err = readVal[T](f, fd)
		if err != nil {
			var empty []T
			return empty, err
		}
		res = append(res, val)
	}

	restCount := 4 - (f.Count % 4)
	if restCount == 4 {
		return res, nil
	}

	f.Count += restCount
	_, err = fd.Seek(int64(f.Count), io.SeekStart)
	if err != nil {
		f.Count -= restCount
		var empty []T
		return empty, err
	}

	return res, nil
}

func readValue(a types.Attr, f *types.File, fd io.ReadSeeker) (interface{}, error) {
	t := a.Type
	if t == types.Double {
		return readAttrValue[float64](f, fd)
	}

	if t == types.Short {
		return readAttrValue[int16](f, fd)
	}

	if t == types.Int {
		return readAttrValue[int32](f, fd)
	}

	if t == types.Byte {
		return readAttrValue[byte](f, fd)
	}

	if t == types.Float {
		return readAttrValue[float32](f, fd)
	}

	if t == types.Char {
		return readString(f, fd)
	}

	return nil, fmt.Errorf("Unsupported type <%s>", t)
}

func readListOf[T any](f *types.File, fd io.ReadSeeker, fn func(f *types.File) (T, error)) (list []T, err error) {
	len, err := readVal[int32](f, fd)
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

func readString(f *types.File, fd io.ReadSeeker) (string, error) {
	v, err := readAttrValue[byte](f, fd)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func readVal[T any](f *types.File, fd io.ReadSeeker) (T, error) {
	var val T
	if err := binary.Read(fd, binary.BigEndian, &val); err != nil {
		var empty T
		return empty, err
	}
	f.Count += uint64(unsafe.Sizeof(val))
	return val, nil
}

func readTag(f *types.File, fd io.ReadSeeker) (types.Tag, error) {
	var buf [4]byte
	if err := binary.Read(fd, binary.BigEndian, &buf); err != nil {
		return types.ZeroTag, err
	}
	f.Count += 4
	return types.Tag(buf[3]), nil
}

// Open ...
func Open(file string) (*types.File, error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	f := types.NewFile()
	if err := readHeader(f, fd); err != nil {
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
