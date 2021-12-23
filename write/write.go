package write

import (
	"encoding/binary"
	"io"

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

	// attrs

	if err := writeAttrs(w, f.Attrs); err != nil {
		return err
	}

	// vars
	if err := writeTag(types.VariableTag, w); err != nil {
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

func writeTag(tag types.Tag, w io.Writer) error {
	var buf [4]byte
	buf[3] = byte(tag)
	if err := binary.Write(w, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}

func writeAttrs(w io.Writer, attrs map[string]types.Attr) error {
	if err := writeTag(types.AttributeTag, w); err != nil {
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

	values := a.Val.([]int16)
	if err := binary.Write(w, binary.BigEndian, int32(len(values))); err != nil {
		return err
	}
	for _, v := range values {
		if err := binary.Write(w, binary.BigEndian, v); err != nil {
			return err
		}
	}
	// TODO: implements meaningful padding
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
