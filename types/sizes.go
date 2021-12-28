package types

func (f File) ByteSize() int32 {
	var szAttrs int
	szAttrs += 8 // len+tag
	for _, it := range f.Attrs.Values() {
		szAttrs += int(it.ByteSize())
	}

	szAttrs += 8 // len+tag
	for _, it := range f.Dimensions {
		szAttrs += int(it.ByteSize())
	}

	szAttrs += 8 // len+tag
	for _, it := range f.Vars.Values() {
		szAttrs += int(it.ByteSize())
	}

	return int32(
		szAttrs +
			4 + // numrecs
			4 + // magic & Version
			0)

}
func (v Var) ByteSize() int32 {
	var szAttrs int32
	szAttrs += 4 + 4 // len+attr tag
	for _, a := range v.Attrs.Values() {
		szAttrs += a.ByteSize()
	}

	return 4 + int32(len(v.Dimensions))*4 + // Dimensions
		szAttrs +
		stringByteSize(v.Name) + // Name string
		4 + //Size
		8 + // Offset
		4 // Type

}

// TODO: add support for array values
func (a Attr) ByteSize() int32 {
	// pad value
	sz := a.ValueByteSize()

	return stringByteSize(a.Name) + // Name string
		4 + // Type
		4 + // len
		sz
}

func (a Attr) ValueByteSize() int32 {
	return int32(a.Type.ArraySize(1))
}

func (v Var) ValueByteSize() int32 {
	var len = 1

	for _, d := range v.Dimensions {
		len *= int(d.Len)
	}

	return int32(v.Type.ArraySize(len))

}

func stringByteSize(val string) int32 {
	return int32(4 + Byte.ArraySize(len(val)))
}

func (d Dimension) ByteSize() int32 {
	return stringByteSize(d.Name) + // Name string
		4 // Len

}
