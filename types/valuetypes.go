package types

import "fmt"

// Type ...
type Type int32

const (
	// Unknown is a placeholder for an empty type value
	Unknown Type = 0
	// Byte is type NC_BYTE = \x00 \x00 \x00 \x01 // 8-bit signed integers
	Byte Type = 1
	// Char is type NC_CHAR = \x00 \x00 \x00 \x02 // text characters
	Char Type = 2
	// Short is type NC_SHORT = \x00 \x00 \x00 \x03 // 16-bit signed integers
	Short Type = 3
	// Int is type NC_INT = \x00 \x00 \x00 \x04 // 32-bit signed integers
	Int Type = 4
	// Float is type NC_FLOAT = \x00 \x00 \x00 \x05 // IEEE single precision floats
	Float Type = 5
	// Double is type NC_DOUBLE = \x00 \x00 \x00 \x06 // IEEE double precision floats
	Double Type = 6
)

func (t Type) CDLName() string {
	switch t {
	case Byte:
		return "byte"
	case Char:
		return "char"
	case Short:
		return "short"
	case Int:
		return "int"
	case Float:
		return "float"
	case Double:
		return "double"
	}

	return fmt.Sprintf("[unknown type:%d]", t)
}

func FromValueType[T BaseType]() Type {
	var v T
	switch (interface{})(v).(type) {
	case float32:
		return Float
	case float64:
		return Double
	case int32:
		return Int
	case int16:
		return Short
	case byte:
		return Byte
	}
	return Unknown
}

func FromCDLName(typeName string) Type {
	switch typeName {
	case "float":
		return Float
	case "byte":
		return Byte
	case "char":
		return Char
	case "short":
		return Short
	case "int":
		return Int
	case "double":
		return Double
	}

	return Unknown
}

func (t Type) ValueToString(value interface{}) string {
	var format string
	switch t {
	case Byte, Short, Int:
		format = "%d"
	case Char:
		format = "%s"
	case Float, Double:
		format = "%g"
	default:
		format = fmt.Sprintf("[UNKNOWN TYPE:%d. VALUE: %%v]", t)
	}
	return fmt.Sprintf(format, value)

}

func (t Type) String() string {
	switch t {
	case Byte:
		return "NC_BYTE"
	case Char:
		return "NC_CHAR"
	case Short:
		return "NC_SHORT"
	case Int:
		return "NC_INT"
	case Float:
		return "NC_FLOAT"
	case Double:
		return "NC_DOUBLE"
	}

	return fmt.Sprintf("[UNKNOWN TYPE:%d]", t)
}

// BaseType ...
type BaseType interface {
	byte | int16 | int32 | float32 | float64
}

// AlignForArrayOf returns the size in bytes of
// padding bytes needed to align to
// 32 bits boundary an array of n elements
// of this type.
// Returns -1 if n <= 0
func (t Type) AlignForArrayOf(n int) int {
	if n <= 0 {
		return -1
	}
	len := t.ScalarSize()
	if len >= 4 {
		return 0
	}
	len *= n

	rest := 4 - (len % 4)
	if rest == 4 {
		rest = 0
	}
	return rest
}

// ArraySize returns the size in bytes of
// an array of n elements, aligned to 32 bits.
// Returns -1 if n <= 0
func (t Type) ArraySize(n int) int {
	if n <= 0 {
		return -1
	}
	v := n * t.ScalarSize()
	pd := t.AlignForArrayOf(n)
	return v + pd
}

// ScalarSize returns the size in bytes of
// a single scalar value of this type.
func (t Type) ScalarSize() int {
	var sz int
	switch t {
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
	return sz
}
