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

	if f.Dimensions, err = readDimensions(f); err != nil {
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

func readDimensions(f *types.File) ([]types.Dimension, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}
	if t != types.DimensionTag {
		return nil, fmt.Errorf("Expected DimensionTag, got %s", t.String())
	}
	return readListOf(f, func(f *types.File) (d types.Dimension, err error) {
		d = types.NewDimension(f)

		if d.Name, err = readString(f); err != nil {
			return d, err
		}

		if d.Len, err = readVal[int32](f); err != nil {
			return d, err
		}

		return
	})
}

func readAttributes(f *types.File) ([]types.Attr, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}

	if t == types.ZeroTag {
		t2, err := readTag(f)
		if err != nil {
			return nil, err
		}

		if t2 != types.ZeroTag {
			return nil, fmt.Errorf("Expected ZeroTag, got %s", t2.String())
		}

		return []types.Attr{}, nil
	}

	if t != types.AttributeTag {
		return nil, fmt.Errorf("Expected AttributeTag, got %s", t.String())
	}

	return readListOf(f, func(f *types.File) (a types.Attr, err error) {
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
}

func readVars(f *types.File) ([]types.Var, error) {
	t, err := readTag(f)
	if err != nil {
		return nil, err
	}
	if t != types.VariableTag {
		return nil, fmt.Errorf("Expected VariableTag, got %s", t.String())
	}

	return readListOf(f, func(f *types.File) (v types.Var, err error) {
		v = types.NewVar(f)

		if v.Name, err = readString(f); err != nil {
			return v, err
		}

		v.Dimensions, err = readListOf(f, func(f *types.File) (*types.Dimension, error) {
			id, err := readVal[int32](f)
			if err != nil {
				return nil, err
			}
			return &f.Dimensions[id], nil
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
