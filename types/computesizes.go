package types

func (f *File) ComputeSizes() *File {
	offset := uint64(f.ByteSize())
	for _, it := range f.Vars.Items() {
		name := it.K
		v := it.V
		v.Offset = offset
		v.Size = v.ValueByteSize()
		f.Vars.Set(name, v)
		offset += uint64(v.Size)
	}
	return f
}
