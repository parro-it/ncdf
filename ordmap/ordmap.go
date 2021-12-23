package ordmap

type OrderedMap[TVal any, TKey comparable] struct {
	seq []TKey
	idx map[TKey]TVal
}

type Item[TVal any, TKey comparable] struct {
	V TVal
	K TKey
}

// From return a new OrderedMap filled with
// the key/values specified in args
func From[TVal any, TKey comparable](args []Item[TVal, TKey]) OrderedMap[TVal, TKey] {
	var res OrderedMap[TVal, TKey]
	for _, a := range args {
		res.Set(a.K, a.V)
	}
	return res
}

// Len returns number of elements contained in the map.
// It returns 0 for an unitialized OrderedMap
func (m *OrderedMap[TVal, TKey]) Len() int {
	return len(m.seq)
}

// Values returns all items contained in the map,
// in the same order as they was inserted in the map.
// It returns an empty slice for an unitialized OrderedMap.
func (m *OrderedMap[TVal, TKey]) Values() []TVal {
	res := make([]TVal, len(m.idx))
	for i, key := range m.seq {
		res[i] = m.idx[key]
	}
	return res
}

// Items returns all key and items contained in the map,
// in the same order as they was inserted.
// It returns an empty slice for an unitialized OrderedMap.
func (m *OrderedMap[TVal, TKey]) Items() []Item[TVal, TKey] {
	res := make([]Item[TVal, TKey], len(m.idx))
	for i, key := range m.seq {
		res[i] = Item[TVal, TKey]{m.idx[key], key}
	}
	return res
}

// Keys returns all keys contained in the map,
// in the same order as they was inserted in the map.
// It returns an empty slice for an unitialized OrderedMap.
func (m *OrderedMap[TVal, TKey]) Keys() []TKey {
	return m.seq
}

// Keys returns all keys contained in the map,
// in the same order as they was inserted in the map.
// It returns an empty slice for an unitialized OrderedMap.
// TODO: no append to m.seq if key is existing
func (m *OrderedMap[TVal, TKey]) Set(key TKey, value TVal) {
	m.seq = append(m.seq, key)
	if m.idx == nil {
		m.idx = make(map[TKey]TVal)
	}
	m.idx[key] = value
}

// Get returns the value contained in the map indexed by key,
// It returns an empty value of TVal if the key doesn't exist.
func (m *OrderedMap[TVal, TKey]) Get(key TKey) TVal {
	if v, ok := m.idx[key]; ok {
		return v
	}
	var empty TVal
	return empty
}

// Get returns if the key exists in the map.
func (m *OrderedMap[TVal, TKey]) Has(key TKey) bool {
	_, ok := m.idx[key]
	return ok
}

// Get returns the sequence position for the key,
// or -1 if it doesn't exist in the map.
func (m *OrderedMap[TVal, TKey]) Find(key TKey) int {

	var seqIdx int = -1
	for i := 0; i < len(m.seq); i++ {
		if m.seq[i] == key {
			seqIdx = i
			break
		}
	}
	return seqIdx
}

// Del remove the item indexed by the key,
// if it exists in the map.
func (m *OrderedMap[TVal, TKey]) Del(key TKey) {
	seqIdx := m.Find(key)
	if seqIdx == -1 {
		return

	}

	delete(m.idx, key)
	m.seq = append(m.seq[0:seqIdx], m.seq[seqIdx+1:]...)
}
