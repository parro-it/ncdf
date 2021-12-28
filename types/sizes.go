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
// TODO: add padding for 32bit alignment
func (a Attr) ByteSize() int32 {
	// pad value
	sz := a.ValueByteSize()

	return stringByteSize(a.Name) + // Name string
		4 + // Type
		4 + // len
		sz
}

func (a Attr) ValueByteSize() int32 {
	var sz int32
	switch a.Type {
	case Double:
		sz = 8
	case Short:
		sz = 2
	case Int:
		sz = 4
	case Byte:
		sz = 1
	case Float:
		sz = 4
	case Char:
		sz = 1
	}

	if sz < 4 {
		sz = 4
	}
	return sz
}

func (v Var) ValueByteSize() int32 {
	var sz int32
	switch v.Type {
	case Double:
		sz = 8
	case Short:
		sz = 2
	case Int:
		sz = 4
	case Byte:
		sz = 1
	case Float:
		sz = 4
	case Char:
		sz = 1
	}
	for _, d := range v.Dimensions {
		sz *= d.Len
	}

	rest := 4 - (sz % 4)
	if rest != 4 {
		sz += rest
	}

	return sz
}

func stringByteSize(val string) int32 {
	len := len(val)
	rest := 4 - (len % 4)
	if rest != 4 {
		len += rest
	}
	return int32(4 + len)
}

func (d Dimension) ByteSize() int32 {
	return stringByteSize(d.Name) + // Name string
		4 // Len

}
