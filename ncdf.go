package ncdf

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/parro-it/ncdf/types"
)

func readHeader(f *types.File) error {
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

	if f.Dimensions, f.DimensionsSeq, err = readDimensions(f); err != nil {
		return err
	}

	if f.Attrs, err = readAttributes(f); err != nil {
		return err
	}

	if f.Vars, err = readVars(f); err != nil {
		return err
	}

	return nil
}

func readDimensions(f *types.File) (map[string]types.Dimension, []*types.Dimension, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, nil, err
	}
	if t == types.ZeroTag {
		m, err := sectionNotPresent[types.Dimension](f)
		return m, []*types.Dimension{}, err
	}
	if t != types.DimensionTag {
		return nil, nil, fmt.Errorf("Expected DimensionTag, got %s", t.String())
	}
	lst, err := readListOf(f, func(f *types.File) (d types.Dimension, err error) {
		d = types.NewDimension(f)

		if d.Name, err = readString(f); err != nil {
			return d, err
		}

		if d.Len, err = readVal[int32](f); err != nil {
			return d, err
		}

		return
	})
	if err != nil {
		return nil, nil, err
	}
	res := map[string]types.Dimension{}
	seq := make([]*types.Dimension, len(lst))
	for i, d := range lst {
		res[d.Name] = d
		seq[i] = &lst[i]

	}
	return res, seq, nil
}

func sectionNotPresent[T any](f *types.File) (map[string]T, error) {
	t2, err := readTag(f)
	if err != nil {
		return nil, err
	}

	if t2 != types.ZeroTag {
		return nil, fmt.Errorf("Expected ZeroTag, got %s", t2.String())
	}

	return map[string]T{}, nil
}

func readAttributes(f *types.File) (map[string]types.Attr, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}

	if t == types.ZeroTag {
		return sectionNotPresent[types.Attr](f)
	}

	if t != types.AttributeTag {
		return nil, fmt.Errorf("Expected AttributeTag, got %s", t.String())
	}

	lst, err := readListOf(f, func(f *types.File) (a types.Attr, err error) {
		a = types.NewAttr(f)

		if a.Name, err = readString(f); err != nil {
			return
		}

		if a.Type, err = readVal[types.Type](f); err != nil {
			return
		}

		if a.Val, err = readValue(a, f); err != nil {
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

func readVars(f *types.File) (map[string]types.Var, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}
	if t == types.ZeroTag {
		return sectionNotPresent[types.Var](f)
	}

	if t != types.VariableTag {
		return nil, fmt.Errorf("Expected VariableTag, got %s", t.String())
	}

	lst, err := readListOf(f, func(f *types.File) (v types.Var, err error) {
		v = types.NewVar(f)

		if v.Name, err = readString(f); err != nil {
			return v, err
		}

		v.Dimensions, err = readListOf(f, func(f *types.File) (*types.Dimension, error) {
			id, err := readVal[int32](f)
			if err != nil {
				return nil, err
			}
			return f.DimensionsSeq[id], nil
		})

		if err != nil {
			return v, err
		}

		if v.Attrs, err = readAttributes(f); err != nil {
			return v, err
		}
		if v.Type, err = readVal[types.Type](f); err != nil {
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

	if err != nil {
		return nil, err
	}
	res := map[string]types.Var{}
	for _, d := range lst {
		res[d.Name] = d
	}
	return res, nil
}

func readAttrValue[T types.BaseType](f *types.File) ([]T, error) {
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

	restCount := 4 - (f.Count % 4)
	if restCount == 4 {
		return res, nil
	}

	_, err = f.ReadBytes(int(restCount))
	if err != nil {
		var empty []T
		return empty, err
	}
	f.Count += restCount

	return res, nil
}

func readValue(a types.Attr, f *types.File) (interface{}, error) {
	t := a.Type
	if t == types.Double {
		return readAttrValue[float64](f)
	}

	if t == types.Short {
		return readAttrValue[int16](f)
	}

	if t == types.Int {
		return readAttrValue[int32](f)
	}

	if t == types.Byte {
		return readAttrValue[byte](f)
	}

	if t == types.Float {
		return readAttrValue[float32](f)
	}

	if t == types.Char {
		return readString(f)
	}

	return nil, fmt.Errorf("Unsupported type <%s>", t)
}

func readListOf[T any](f *types.File, fn func(f *types.File) (T, error)) (list []T, err error) {
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

func readString(f *types.File) (string, error) {
	v, err := readAttrValue[byte](f)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func readVal[T any](f *types.File) (T, error) {
	var val T
	if err := f.Read(&val); err != nil {
		var empty T
		return empty, err
	}
	f.Count += uint64(unsafe.Sizeof(val))
	return val, nil
}

func readTag(f *types.File) (types.Tag, error) {

	buf, err := f.ReadBytes(4)
	if err != nil {
		return types.ZeroTag, err
	}
	f.Count += 4
	return types.Tag(buf[3]), err
}

// Open ...
func Open(file string) (*types.File, error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	f := types.NewFile(fd)
	if err := readHeader(f); err != nil {
		return nil, err
	}

	return f, nil
}

func VarData[T types.BaseType](v types.Var, f *types.File) ([]T, error) {
	if err := f.Seek(int64(v.Offset)); err != nil {
		return nil, err
	}
	data := make([]T, v.Size)
	if err := f.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
