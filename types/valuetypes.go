package types

// Type ...
type Type int32

const (
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

// BaseType ...
type BaseType interface {
	byte | int16 | int32 | float32 | float64
}
