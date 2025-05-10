package kiso

type IsoDefinition struct {
	lengthFormat      string
	fieldDescriptions [64 * 3]*FieldDescription
}

type FieldDescription struct {
	n           int
	description string
	size        int
	format      string
}

type Field struct {
	n     int
	value []byte
}

type Message struct {
	mti    int
	fields [64 * 3]*Field
}

const (
	LengthFormat_ASCII = "ascii"
	LengthFormat_BCD   = "bcd"
	LengthFormat_HEX   = "hex"

	FieldFormat_ASCII_Bitmap = "ascii_bitmap"
	FieldFormat_ASCII_B      = "ascii_b"
	FieldFormat_ASCII_N      = "ascii_n"
	FieldFormat_ASCII_Z      = "ascii_z"
	FieldFormat_Bitmap       = "bitmap"
	FieldFormat_ANS          = "ans"
	FieldFormat_AN           = "an"
	FieldFormat_NS           = "ns"
	FieldFormat_ANP          = "anp"
	FieldFormat_N            = "n"
	FieldFormat_B            = "b"
	FieldFormat_Z            = "z"
)
