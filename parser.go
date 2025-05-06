package kiso

import (
	"encoding/hex"
	"fmt"
)

func asciiHexNibblesToBytes(b []byte) []byte {
	if len(b)%2 != 0 {
		panic("invalid length")
	}
	res, err := hex.DecodeString(string(b))
	if err != nil {
		panic(err)
	}
	return res
}

func bytesToAsciiHexNibbles(b []byte) []byte {
	return []byte(hex.EncodeToString(b))
}

func (isod *IsoDescription) formatMTI(mti int) []byte {
	d := isod.fieldDescriptions[0]
	if d.format == FieldFormat_ASCII_N {
		return []byte(fmt.Sprintf("%04d", mti))
	} else {
		bcdMti, _ := hex.DecodeString(fmt.Sprintf("%04d", mti))
		return bcdMti
	}
}
func (isod *IsoDescription) formatBitmap(b []byte) []byte {
	d := isod.fieldDescriptions[1]
	if d.format == FieldFormat_ASCII_Bitmap {
		return bytesToAsciiHexNibbles(b)
	} else {
		return b
	}
}
func (isod *IsoDescription) formatField(field *Field) []byte {
	d := isod.fieldDescriptions[field.n]

	var v = field.value

	if d.size < 0 {
		if isod.lengthFormat == LengthFormat_ASCII {
			if d.size >= -99 {
				v = append([]byte(fmt.Sprintf("%02d", len(v))), v...)
			} else if d.size >= -999 {
				v = append([]byte(fmt.Sprintf("%03d", len(v))), v...)
			} else if d.size >= -9999 {
				v = append([]byte(fmt.Sprintf("%04d", len(v))), v...)
			}
		} else if isod.lengthFormat == LengthFormat_BCD {
			if d.size >= -99 {
				v = append([]byte{byte((len(v)/10)<<4 | len(v)%10)}, v...)
			} else if d.size >= -999 {
				v = append([]byte{byte((len(v) / 100) & 0xF), byte((len(v)%100/10)<<4 | len(v)%10)}, v...)
			} else if d.size >= -9999 {
				v = append([]byte{byte((len(v)/1000)<<4 | len(v)%1000/100), byte((len(v)%100/10)<<4 | len(v)%10)}, v...)
			}
		} else if isod.lengthFormat == LengthFormat_HEX {
			if d.size >= -99 {
				v = append([]byte{byte(len(v) & 0xFF)}, v...)
			} else if d.size >= -999 {
				v = append([]byte{byte((len(v) >> 8) & 0xFF), byte(len(v) & 0xFF)}, v...)
			}
		}
	}
	return v
}

func (isod *IsoDescription) getField(n int, message []byte) ([]byte, int) {
	d := isod.fieldDescriptions[n]
	if d == nil {
		panic(fmt.Sprintf("field [%d] description not found", n))
	}

	if d.size > 0 {
		if d.format == FieldFormat_N {
			return message[:(d.size+1)/2], int((d.size + 1) / 2)
		}
		if d.format == FieldFormat_ASCII_B {
			return message[:d.size*2], int(d.size * 2)
		}
		if d.format == FieldFormat_ASCII_Bitmap {
			b := asciiHexNibblesToBytes(message[:d.size*2])
			if b[0]&0x80 == 0x80 {
				b = asciiHexNibblesToBytes(message[:d.size*4])
				if b[0+8]&0x80 == 0x80 {
					b = asciiHexNibblesToBytes(message[:d.size*6])
				}
			}
			return b, len(b) * 2
		}
		if d.format == FieldFormat_Bitmap {
			b := message[:d.size]
			if b[0]&0x80 == 0x80 {
				b = message[:d.size*2]
				if b[0+8]&0x80 == 0x80 {
					b = message[:d.size*3]
				}
			}
			return b, len(b)
		}
		return message[:d.size], int(d.size)
	}
	if d.size < 0 {
		size := 0
		lengthSize := 0
		if isod.lengthFormat == LengthFormat_ASCII {
			if d.size >= -99 {
				lengthSize = 2
				size = (int(message[0])-0x30)*10 + int(message[1]) - 0x30
			} else if d.size >= -999 {
				lengthSize = 3
				size = (int(message[0])-0x30)*100 + (int(message[1])-0x30)*10 + int(message[2]) - 0x30
			} else if d.size >= -9999 {
				lengthSize = 4
				size = (int(message[0])-0x30)*1000 + (int(message[1])-0x30)*100 + (int(message[2])-0x30)*10 + int(message[3]) - 0x30
			}
		} else if isod.lengthFormat == LengthFormat_BCD {
			if d.size >= -99 {
				lengthSize = 1
				size = int(message[0]>>4)*10 + int(message[0]&0x0F)
			} else if d.size >= -999 {
				lengthSize = 2
				size = int(message[0]&0x0F)*100 + int(message[1]>>4)*10 + int(message[1]&0x0F)
			} else if d.size >= -9999 {
				lengthSize = 2
				size = int(message[0]>>4)*1000 + int(message[0]&0x0F)*100 + int(message[1]>>4)*10 + int(message[1]&0x0F)
			}
		} else if isod.lengthFormat == LengthFormat_HEX {
			if d.size >= -99 {
				lengthSize = 1
				size = int(message[0])
			} else if d.size >= -9999 {
				lengthSize = 2
				size = (int(message[0]) << 8) | int(message[1])
			}
		}
		if size > 0 {
			return message[lengthSize : lengthSize+size], int(lengthSize + size)
		} else {
			panic(fmt.Sprintf("unsupported field(%d) size(%d)", n, size))
		}
	}
	panic(fmt.Sprintf("unsupported field(%d) format(%s)", n, d.format))
}

func Compose(mti int, fields []Field, isod IsoDescription) []byte {
	var message []byte
	var bitmap [24]byte
	var maxFieldNo int = 0

	for _, f := range fields {
		if f.n == 0 || f.n == 1 {
			continue
		}
		message = append(message, isod.formatField(&f)...)

		shiftBy := uint8(8 - f.n%8)
		if shiftBy == 8 {
			shiftBy = 0
		}
		bitmap[(f.n-1)/8] |= (1 << shiftBy)

		maxFieldNo = max(maxFieldNo, f.n) - 1
	}

	if maxFieldNo > 64 {
		bitmap[0] |= 0x80
	}
	if maxFieldNo > 128 {
		bitmap[8] |= 0x80
	}

	message = append(isod.formatBitmap(bitmap[:(maxFieldNo/64+1)*8]), message...)
	message = append(isod.formatMTI(mti), message...)
	return message
}

func Parse(message []byte, isod IsoDescription) []Field {
	var fields = make([]Field, 0)
	fieldNo := 0
	c := 0
	fieldBytes, cc := isod.getField(fieldNo, message[c:])
	c += cc
	fields = append(fields, Field{fieldNo, fieldBytes})
	fieldNo++

	fieldBytes, cc = isod.getField(fieldNo, message[c:])
	fmt.Println(len(fieldBytes))
	fields = append(fields, Field{fieldNo, fieldBytes})
	bitmap := fieldBytes
	c += cc

	s := ""
	for i := 0; i < len(bitmap); i++ {
		for j := 0; j < 8; j++ {

			if i == 0 && j == 0 {
				continue
			}
			fieldNo++
			if (bitmap[i] & (1 << uint8(7-j))) != 0 {
				fieldBytes, cc = isod.getField(fieldNo, message[c:])
				fields = append(fields, Field{fieldNo, fieldBytes})
				c += cc
				s += "1"
			} else {
				s += "0"
			}
		}
		s += " "

	}
	fmt.Println(s)
	return fields
}

func ParsToMap(message []byte, isod IsoDescription) map[int]Field {
	var fields = make(map[int]Field, 0)

	for i, f := range Parse(message, isod) {
		fields[i] = f
	}

	return fields
}
