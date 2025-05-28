package kiso

import (
	"fmt"
	"strconv"
	"strings"
)

func MakeFieldDescription(n int, description string, size int, format string) FieldDescription {
	return FieldDescription{
		n:           n,
		description: description,
		size:        size,
		format:      format,
	}
}

func (isod *IsoDefinition) MakeFieldAscii(n int, v string) *Field {
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

	return &Field{
		n:     n,
		value: []byte(v)}
}

func (isod *IsoDefinition) MakeFieldBinary(n int, b []byte) *Field {
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

	return &Field{
		n:     n,
		value: b}
}

func (isod *IsoDefinition) GetFieldValueAscii(field *Field) (string, error) {
	d := isod.fieldDescriptions[field.n]
	if d.format == FieldFormat_ASCII_Bitmap ||
		d.format == FieldFormat_ASCII_B ||
		d.format == FieldFormat_ASCII_N ||
		d.format == FieldFormat_ASCII_Z ||
		d.format == FieldFormat_ANS ||
		d.format == FieldFormat_AN ||
		d.format == FieldFormat_NS ||
		d.format == FieldFormat_ANP {
		return string(field.value), nil

	} else {
		return "", fmt.Errorf("non-ascii field")
	}
}

func (isod *IsoDefinition) GetFieldValueBinary(field *Field) ([]byte, error) {
	d := isod.fieldDescriptions[field.n]
	if d.format != FieldFormat_N &&
		d.format != FieldFormat_B &&
		d.format != FieldFormat_Z {
		return nil, fmt.Errorf("non-byte field")
	} else {
		return field.value, nil
	}
}

func (isod *IsoDefinition) GetMti(field *Field) (int, error) {
	if field.n != 0 {
		return 0, fmt.Errorf("not mti field")
	}
	d := isod.fieldDescriptions[field.n]

	if d.format == FieldFormat_ASCII_N {
		return int(field.value[0]-'0')*1000 +
			int(field.value[1]-'0')*100 +
			int(field.value[2]-'0')*10 +
			int(field.value[3]-'0'), nil
	}
	if d.format == FieldFormat_N {
		return int(field.value[0]>>4)*1000 +
			int(field.value[0]&0x0F)*100 +
			int(field.value[1]>>4)*10 +
			int(field.value[1]&0x0F), nil
	}
	return 0, fmt.Errorf("non mti format field")

}

func MakeIsoDescription(lengthFormat string, fieldDescriptions []FieldDescription) IsoDefinition {
	var isod = IsoDefinition{
		lengthFormat: lengthFormat,
	}
	for i := 0; i < len(fieldDescriptions); i++ {
		isod.fieldDescriptions[fieldDescriptions[i].n] = &fieldDescriptions[i]
	}
	return isod
}

func (isod *IsoDefinition) MakeSampleField(n int) *Field {
	d := isod.fieldDescriptions[n]
	if d == nil {
		return nil
	}
	if d.format != FieldFormat_N &&
		d.format != FieldFormat_B &&
		d.format != FieldFormat_Z {
		var f = isod.MakeFieldAscii(n, strconv.Itoa(n))
		return f
	} else {
		var f = isod.MakeFieldBinary(n, []byte{byte(n)})
		return f
	}
}

func (m *Message) GetMti() int {
	return m.mti
}
func (m *Message) SetMti(mti int) {
	m.mti = mti
}
func (m *Message) GetField(n int) *Field {
	return m.fields[n]
}
func (m *Message) SetField(f *Field) {
	m.fields[f.n] = f
}

func (isod *IsoDefinition) MessageToString(m *Message) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("MTI %04d, ", m.mti))
	maxDescrLen := 0
	maxFieldNo := 0
	var bitmap [24]byte

	for i := 2; i < len(m.fields); i++ {
		if m.fields[i] != nil {
			maxFieldNo = i - 1
			maxDescrLen = max(len(isod.fieldDescriptions[i].description), maxDescrLen)
			shiftBy := uint8(8 - i%8)
			if shiftBy == 8 {
				shiftBy = 0
			}
			bitmap[(i-1)/8] |= (1 << shiftBy)
		}
	}
	if maxFieldNo > 64 {
		bitmap[0] |= 0x80
	}
	if maxFieldNo > 128 {
		bitmap[8] |= 0x80
	}
	bitmapToPrint := bitmap[:(maxFieldNo/64+1)*8]

	sb.WriteString(fmt.Sprintf("bitmap: %X", bitmapToPrint))
	nums := ""
	for i := 0; i < len(bitmapToPrint); i++ {
		if i%8 == 0 {
			sb.WriteString("\n")
			sb.WriteString(nums)
			sb.WriteString("\n")
			nums = ""
		}
		sb.WriteString(fmt.Sprintf("%08b ", bitmap[i]))
		nums += fmt.Sprintf("%-8s|", fmt.Sprintf("%d", i*8+1))
		nums = strings.ReplaceAll(nums, " ", ".")
	}
	sb.WriteString("\n")
	sb.WriteString(nums)
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("%-7s|%-*s|%-6s|%-s\n", strings.Repeat("-", 7), maxDescrLen, strings.Repeat("-", maxDescrLen+2), strings.Repeat("-", 6), strings.Repeat("-", 15)))
	sb.WriteString(fmt.Sprintf("%-7s|%-*s|%-6s|%-s\n", " FIELD", maxDescrLen+2, " NAME", " LEN", " VALUE"))
	sb.WriteString(fmt.Sprintf("%-7s|%-*s|%-6s|%-s\n", strings.Repeat("-", 7), maxDescrLen, strings.Repeat("-", maxDescrLen+2), strings.Repeat("-", 6), strings.Repeat("-", 15)))
	for i := 2; i < len(m.fields); i++ {
		if m.fields[i] != nil {
			d := isod.fieldDescriptions[m.fields[i].n]

			var value string
			if d.format == FieldFormat_ASCII_B {
				value = string(m.fields[i].value)
			} else if d.format == FieldFormat_ASCII_Bitmap {
				value = string(bytesToAsciiHexNibbles(m.fields[i].value))
			} else if d.format == FieldFormat_N || d.format == FieldFormat_B || d.format == FieldFormat_Z {
				value = fmt.Sprintf("%X", m.fields[i].value)
			} else {
				value = string(m.fields[i].value)
			}

			sb.WriteString(fmt.Sprintf("DE.%03d : %-*s : %04d : %s\n", m.fields[i].n, maxDescrLen, d.description, len(m.fields[i].value), value))
		}
	}
	return sb.String()
}
