package types

func (f *File) ComputeSizes() *File {
	offset := uint64(f.ByteSize())
	for name, v := range f.Vars {
		v.Offset = offset
		v.Size = v.ValueByteSize()
		f.Vars[name] = v
		offset += uint64(v.Size)
	}
	return f
}
