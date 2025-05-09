package kiso

import (
	"reflect"
	"testing"

	iso "github.com/Jesting/kiso"
)

var allFields = []iso.FieldDescription{
	iso.MakeFieldDescription(0, "MTI_ASII_N", 4, iso.FieldFormat_ASCII_N),
	iso.MakeFieldDescription(1, "Bitmap_ASCII", 8, iso.FieldFormat_ASCII_Bitmap),
	iso.MakeFieldDescription(2, "F_ASCII_B", 8, iso.FieldFormat_ASCII_B),
	iso.MakeFieldDescription(3, "F_ASCII_N", -1200, iso.FieldFormat_ASCII_N),
	iso.MakeFieldDescription(4, "F_ASCII_Z", 12, iso.FieldFormat_ASCII_Z),
	iso.MakeFieldDescription(6, "F_ANS", -100, iso.FieldFormat_ANS),
	iso.MakeFieldDescription(7, "F_AN", 10, iso.FieldFormat_AN),
	iso.MakeFieldDescription(8, "F_NS", 10, iso.FieldFormat_NS),
	iso.MakeFieldDescription(9, "F_ANP", -10, iso.FieldFormat_ANP),
	iso.MakeFieldDescription(10, "F_N", 10, iso.FieldFormat_N),
	iso.MakeFieldDescription(11, "F_B", 10, iso.FieldFormat_B),
	iso.MakeFieldDescription(12, "F_Z", 10, iso.FieldFormat_Z),
	iso.MakeFieldDescription(13, "F_ASCII_B", -999, iso.FieldFormat_ASCII_B),
	iso.MakeFieldDescription(14, "F_ASCII_N", -99, iso.FieldFormat_ASCII_N),
	iso.MakeFieldDescription(15, "F_ASCII_Z", -9999, iso.FieldFormat_ASCII_Z),
	iso.MakeFieldDescription(16, "F_ANS", 55, iso.FieldFormat_ANS),
	iso.MakeFieldDescription(17, "F_AN", -100, iso.FieldFormat_AN),
	iso.MakeFieldDescription(18, "F_NS", 10, iso.FieldFormat_NS),
	iso.MakeFieldDescription(19, "F_ANP", -10, iso.FieldFormat_ANP),
	iso.MakeFieldDescription(20, "F_N", -9999, iso.FieldFormat_N),
	iso.MakeFieldDescription(21, "F_B", -999, iso.FieldFormat_B),
	iso.MakeFieldDescription(22, "F_Z", 10, iso.FieldFormat_Z),
}

func TestLenHex(t *testing.T) {

	var isod = iso.MakeIsoDescription(iso.LengthFormat_ASCII, allFields)

	var fields = []iso.Field{}

	for i := 2; i < 3*64; i++ {
		var field = isod.MakeSampleField(i)
		if field != nil {
			fields = append(fields, *field)
		}
	}

	var message = isod.Compose(1234, fields)
	var res = isod.Parse(message)

	for i := 2; i < len(fields); i++ {
		if !reflect.DeepEqual(fields[i-2], res[i]) {
			t.Fail()
		}
	}
	if isod.GetMti(res[0]) != 1234 {
		t.Fail()
	}
}

func TestLenAscii(t *testing.T) {
	var isod = iso.MakeIsoDescription(iso.LengthFormat_ASCII, allFields)

	var fields = []iso.Field{}

	for i := 2; i < 3*64; i++ {
		var field = isod.MakeSampleField(i)
		if field != nil {
			fields = append(fields, *field)
		}
	}

	var message = isod.Compose(1234, fields)

	var res = isod.Parse(message)

	for i := 2; i < len(fields); i++ {
		if !reflect.DeepEqual(fields[i-2], res[i]) {
			t.Fail()
		}
	}

	if isod.GetMti(res[0]) != 1234 {
		t.Fail()
	}
}
func TestLenBCD(t *testing.T) {
	var isod = iso.MakeIsoDescription(iso.LengthFormat_BCD, allFields)

	var fields = []iso.Field{}

	for i := 2; i < 3*64; i++ {
		var field = isod.MakeSampleField(i)
		if field != nil {
			fields = append(fields, *field)
		}
	}

	var message = isod.Compose(1210, fields)

	var res = isod.Parse(message)

	for i := 2; i < len(fields); i++ {
		if !reflect.DeepEqual(fields[i-2], res[i]) {
			t.Fail()
		}
	}

	if isod.GetMti(res[0]) != 1210 {
		t.Fail()
	}

}
