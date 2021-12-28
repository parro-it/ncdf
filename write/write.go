package write

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/parro-it/ncdf/ordmap"
	"github.com/parro-it/ncdf/types"
)

// VarData ...
// TODO: use missing value for data
func VarData[T types.BaseType](v types.Var, data []T, fd io.WriterAt) error {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, data); err != nil {
		return err
	}
	if _, err := fd.WriteAt(buf.Bytes(), int64(v.Offset)); err != nil {
		return err
	}
	return nil
}

func Header(f *types.File, w io.Writer) error {
	// magic + version
	bytes := []byte{'C', 'D', 'F', 2}
	if _, err := w.Write(bytes); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, f.NumRecs); err != nil {
		return err
	}
	// dimensions

	if f.Dimensions == nil || len(f.Dimensions) == 0 {
		if err := writeTag(types.ZeroTag, w); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, int32(0)); err != nil {
			return err
		}
	} else {
		tag := types.DimensionTag

		if err := writeTag(tag, w); err != nil {
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
	}
	// attrs

	if err := writeAttrs(w, f.Attrs); err != nil {
		return err
	}

	// vars
	if f.Vars.Len() == 0 {
		if err := writeTag(types.ZeroTag, w); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, int32(0)); err != nil {
			return err
		}
	} else {
		if err := writeTag(types.VariableTag, w); err != nil {
			return err
		}

		if err := binary.Write(w, binary.BigEndian, int32(f.Vars.Len())); err != nil {
			return err
		}

		for _, v := range f.Vars.Values() {
			if err := writeVar(f, v, w); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeTag(tag types.Tag, w io.Writer) error {
	buf := []byte{0, 0, 0, byte(tag)}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

func writeAttrs(w io.Writer, attrs ordmap.OrderedMap[types.Attr, string]) error {
	if attrs.Len() == 0 {
		if err := writeTag(types.ZeroTag, w); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, int32(0)); err != nil {
			return err
		}
		return nil
	}
	if err := writeTag(types.AttributeTag, w); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, int32(attrs.Len())); err != nil {
		return err
	}

	for _, a := range attrs.Values() {
		if err := writeAttr(a, w); err != nil {
			return err
		}
	}
	return nil
}

func writeAttr(a types.Attr, w io.Writer) error {
	if err := writeSlice(w, []byte(a.Name)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, a.Type); err != nil {
		return err
	}

	err := writeAttrValue(a, w)
	if err != nil {
		return err
	}

	return nil
}

func writeSlice[T types.BaseType](w io.Writer, val []T) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(val))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, val); err != nil {
		return err
	}

	rest := types.FromValueType[T]().AlignForArrayOf(len(val))
	if rest > 0 {
		buf := make([]byte, rest)
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

func writeAttrValue(a types.Attr, w io.Writer) error {
	values := a.Val.([]int16)
	return writeSlice(w, values)
}

func writeVar(f *types.File, v types.Var, w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(v.Name))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(v.Name)); err != nil {
		return err
	}
	rest := 4 - (len(v.Name) % 4)
	if rest != 4 {
		buf := make([]byte, rest)
		if _, err := w.Write(buf); err != nil {
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
