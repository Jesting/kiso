package kiso

import (
	"encoding/hex"
	"fmt"
)

func asciiHexNibblesToBytes(b []byte) ([]byte, error) {
	if len(b)%2 != 0 {
		return []byte{}, fmt.Errorf("invalid length")
	}
	res, err := hex.DecodeString(string(b))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func bytesToAsciiHexNibbles(b []byte) []byte {
	return []byte(hex.EncodeToString(b))
}

func (isod *IsoDefinition) formatMTI(mti int) []byte {
	d := isod.fieldDescriptions[0]
	if d.format == FieldFormat_ASCII_N {
		return []byte(fmt.Sprintf("%04d", mti))
	} else {
		bcdMti, _ := hex.DecodeString(fmt.Sprintf("%04d", mti))
		return bcdMti
	}
}
func (isod *IsoDefinition) formatBitmap(b []byte) []byte {
	d := isod.fieldDescriptions[1]
	if d.format == FieldFormat_ASCII_Bitmap {
		return bytesToAsciiHexNibbles(b)
	} else {
		return b
	}
}
func (isod *IsoDefinition) formatField(field *Field) []byte {
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
			} else if d.size >= -9999 {
				v = append([]byte{byte((len(v) >> 8) & 0xFF), byte(len(v) & 0xFF)}, v...)
			}
		}
	}
	return v
}

func (isod *IsoDefinition) getField(n int, message []byte) ([]byte, int, error) {
	d := isod.fieldDescriptions[n]
	if d == nil {
		return nil, 0, fmt.Errorf("field [%d] description not found", n)
	}

	if d.size > 0 {
		if d.format == FieldFormat_N {
			return message[:(d.size+1)/2], int((d.size + 1) / 2), nil
		}
		if d.format == FieldFormat_ASCII_B {
			return message[:d.size*2], int(d.size * 2), nil
		}
		if d.format == FieldFormat_ASCII_Bitmap {
			b, err := asciiHexNibblesToBytes(message[:d.size*2])
			if err != nil {
				return nil, 0, err
			}
			if b[0]&0x80 == 0x80 {
				b, err = asciiHexNibblesToBytes(message[:d.size*4])
				if err != nil {
					return nil, 0, err
				}
				if b[0+8]&0x80 == 0x80 {
					b, err = asciiHexNibblesToBytes(message[:d.size*6])
					if err != nil {
						return nil, 0, err
					}
				}
			}
			return b, len(b) * 2, nil
		}
		if d.format == FieldFormat_Bitmap {
			b := message[:d.size]
			if b[0]&0x80 == 0x80 {
				b = message[:d.size*2]
				if b[0+8]&0x80 == 0x80 {
					b = message[:d.size*3]
				}
			}
			return b, len(b), nil
		}
		return message[:d.size], int(d.size), nil
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
			return message[lengthSize : lengthSize+size], int(lengthSize + size), nil
		} else {
			return nil, 0, fmt.Errorf("unsupported field(%d) size(%d)", n, size)
		}
	}
	return nil, 0, fmt.Errorf("unsupported field(%d) format(%s)", n, d.format)
}

func (isod *IsoDefinition) Compose(mti int, fields []*Field) ([]byte, error) {
	var message []byte
	var bitmap [24]byte
	var maxFieldNo int = -1

	for _, f := range fields {
		if f == nil {
			continue
		}
		if f.n <= maxFieldNo {
			return nil, fmt.Errorf("duplicate field %d or wrong field order", f.n)
		}
		maxFieldNo = max(maxFieldNo, f.n) - 1

		if f.n == 0 || f.n == 1 {
			continue
		}
		message = append(message, isod.formatField(f)...)

		shiftBy := uint8(8 - f.n%8)
		if shiftBy == 8 {
			shiftBy = 0
		}
		bitmap[(f.n-1)/8] |= (1 << shiftBy)

	}

	if maxFieldNo > 64 {
		bitmap[0] |= 0x80
	}
	if maxFieldNo > 128 {
		bitmap[8] |= 0x80
	}

	message = append(isod.formatBitmap(bitmap[:(maxFieldNo/64+1)*8]), message...)
	message = append(isod.formatMTI(mti), message...)
	return message, nil
}

func (isod *IsoDefinition) Parse(message []byte) ([]*Field, error) {
	var fields = make([]*Field, 0)
	fieldNo := 0
	c := 0
	fieldBytes, cc, err := isod.getField(fieldNo, message[c:])
	if err != nil {
		return nil, err
	}

	c += cc
	fields = append(fields, &Field{fieldNo, fieldBytes})
	fieldNo++

	fieldBytes, cc, err = isod.getField(fieldNo, message[c:])
	if err != nil {
		return nil, err
	}

	fields = append(fields, &Field{fieldNo, fieldBytes})
	bitmap := fieldBytes
	c += cc
	for i := 0; i < len(bitmap); i++ {
		for j := 0; j < 8; j++ {

			if i == 0 && j == 0 {
				continue
			}
			fieldNo++
			if (bitmap[i] & (1 << uint8(7-j))) != 0 {
				fieldBytes, cc, err = isod.getField(fieldNo, message[c:])
				if err != nil {
					return nil, err
				}
				fields = append(fields, &Field{fieldNo, fieldBytes})
				c += cc
			}
		}
	}
	return fields, nil
}

func (isod *IsoDefinition) ParseToMessage(messageBytes []byte) (*Message, error) {
	var message = Message{}
	parsed, err := isod.Parse(messageBytes)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(parsed); i++ {
		message.fields[parsed[i].n] = parsed[i]
	}
	message.mti, err = isod.GetMti(parsed[0])

	if err != nil {
		return nil, err
	}

	return &message, err
}
func (isod *IsoDefinition) ComposeFromMessage(message *Message) ([]byte, error) {
	var res, err = isod.Compose(message.mti, message.fields[:])
	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
