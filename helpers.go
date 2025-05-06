package kiso

import (
	"fmt"
	"strconv"
)

func MakeFieldDescription(n int, description string, size int, format string) FieldDescription {
	return FieldDescription{
		n:           n,
		description: description,
		size:        size,
		format:      format,
	}
}

func (isod *IsoDescription) ToString(field Field) string {
	d := isod.fieldDescriptions[field.n]

	var value string
	if d.format == FieldFormat_ASCII_B {
		value = string(field.value)
	} else if d.format == FieldFormat_ASCII_Bitmap {
		value = string(bytesToAsciiHexNibbles(field.value))
	} else if d.format == FieldFormat_N || d.format == FieldFormat_B || d.format == FieldFormat_Z {
		value = fmt.Sprintf("%X", field.value)
	} else {
		value = string(field.value)
	}
	return fmt.Sprintf("DF.%03d : (%03d) %-20s %s", field.n, len(field.value), value, d.description)
}

func (isod *IsoDescription) MakeFieldAscii(n int, v string) Field {
	d := isod.fieldDescriptions[n]
	if d == nil {
		panic(fmt.Sprintf("field (%d) description not found", n))
	}
	if d.format == FieldFormat_ASCII_Bitmap || d.format == FieldFormat_Bitmap {
		panic("bitmaps are calculated inplace")
	}
	if d.format == FieldFormat_N ||
		d.format == FieldFormat_B ||
		d.format == FieldFormat_Z {
		panic("non-ascii field")
	}

	if d.size > 0 {
		if d.format == FieldFormat_ASCII_N {
			if len(v) > d.size {
				panic("too long")
			}
			v = fmt.Sprintf("%0*s", d.size, v)
		} else if d.format == FieldFormat_ASCII_B {
			if len(v) > d.size*2 {
				fmt.Printf("len %d, size %d", len(v), d.size*2)
				panic("too long")
			}
			v = fmt.Sprintf("%s%0*s", v, d.size*2-len(v), "")
		} else {
			if len(v) > d.size {
				panic("too long")
			}
			v = fmt.Sprintf("%-*s", d.size, v)
		}
	}

	return Field{
		n:     n,
		value: []byte(v)}
}

func (isod *IsoDescription) MakeFieldBinary(n int, b []byte) Field {
	d := isod.fieldDescriptions[n]
	if d == nil {
		panic(fmt.Sprintf("field [%d] description not found", n))
	}
	if d.format == FieldFormat_ASCII_Bitmap || d.format == FieldFormat_Bitmap {
		panic("bitmaps are calculated inplace")
	}
	if d.format != FieldFormat_N &&
		d.format != FieldFormat_B &&
		d.format != FieldFormat_Z {
		panic("non-byte field")
	}

	if d.size > 0 {
		if d.format == FieldFormat_N {
			if len(b) > (d.size+1)/2 {
				panic("too long")
			}
			diff := (d.size+1)/2 - len(b)
			b = append(make([]byte, diff), b...)
		} else if d.format == FieldFormat_B {
			if len(b) > d.size {
				panic("too long")
			}
			diff := d.size - len(b)
			b = append(b, make([]byte, diff)...)
		} else {
			if len(b) > d.size {
				panic("too long")
			}
			b = []byte(fmt.Sprintf("%-*s", d.size, string(b)))
		}
	}

	return Field{
		n:     n,
		value: b}
}

func (isod *IsoDescription) GetFieldValueAscii(field Field) string {
	d := isod.fieldDescriptions[field.n]
	if d.format == FieldFormat_ASCII_Bitmap ||
		d.format == FieldFormat_ASCII_B ||
		d.format == FieldFormat_ASCII_N ||
		d.format == FieldFormat_ASCII_Z ||
		d.format == FieldFormat_ANS ||
		d.format == FieldFormat_AN ||
		d.format == FieldFormat_NS ||
		d.format == FieldFormat_ANP {
		return string(field.value)

	} else {
		panic("non-ascii field")
	}
}

func (isod *IsoDescription) GetFieldValueBinary(field Field) []byte {
	d := isod.fieldDescriptions[field.n]
	if d.format != FieldFormat_N &&
		d.format != FieldFormat_B &&
		d.format != FieldFormat_Z {
		panic("non-byte field")
	} else {
		return field.value
	}
}

func MakeIsoDescription(lengthFormat string, fieldDescriptions []FieldDescription) IsoDescription {
	var isod = IsoDescription{
		lengthFormat: lengthFormat,
	}
	for i := 0; i < len(fieldDescriptions); i++ {
		isod.fieldDescriptions[fieldDescriptions[i].n] = &fieldDescriptions[i]
	}
	return isod
}

func (isod *IsoDescription) MakeSampleField(n int) *Field {
	d := isod.fieldDescriptions[n]
	if d == nil {
		return nil
	}
	if d.format != FieldFormat_N &&
		d.format != FieldFormat_B &&
		d.format != FieldFormat_Z {
		var f = isod.MakeFieldAscii(n, strconv.Itoa(n))
		return &f
	} else {
		var f = isod.MakeFieldBinary(n, []byte{byte(n)})
		return &f
	}

}
