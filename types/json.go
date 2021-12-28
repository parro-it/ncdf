package types

import "fmt"

// MarshalJSON ...
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

// MarshalJSON ...
func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, v[3])), nil
}
