# KISO
## _for Payments101_

KISO is simple ISO8583 parser + some standard primitives
**Code supports:**
    * LLVar, LLLVar, LLLLVar fields length
    * ASCII, BCD & HEX length prefixes
    * ASCII & Binary fields
    * MTI composition
    * 3 bitmaps
    * Padding

## To use the app you need:
- [Go compiler]

## Import

```
go get github.com/Jesting/kiso
```

## Use
```
	package main
	import (
	"fmt"
	iso "github.com/Jesting/kiso"
	)
	
	func main() {
	
		var isod = iso.MakeIsoDescription(iso.LengthFormat_ASCII, iso.Fields93)
		var message = iso.Compose(1100, []iso.Field{
			isod.MakeFieldAscii(2, "1234567891234567"),
			isod.MakeFieldAscii(3, "000000"),
			isod.MakeFieldAscii(4, "53100"),
			isod.MakeFieldAscii(12, "250428001122"),
			isod.MakeFieldAscii(14, "2505"),
			isod.MakeFieldAscii(22, "010203010203"),
			isod.MakeFieldAscii(23, "003"),
			isod.MakeFieldAscii(24, "100"),
			isod.MakeFieldAscii(25, "1506"),
			isod.MakeFieldAscii(28, "250428"),
			isod.MakeFieldAscii(29, "000"),
			isod.MakeFieldAscii(35, "1234567891234567=25052011000000000000"),
			isod.MakeFieldAscii(41, "12345678"),
			isod.MakeFieldAscii(42, "1234567890"),
			isod.MakeFieldAscii(49, "978"),
			isod.MakeFieldAscii(52, "AB67CDEF1A34A890"),
		}, isod)

		fmt.Println("Message:", string(message))

		var res = iso.Parse(message, isod)

		for i := 0; i < len(res); i++ {
			fmt.Println(isod.ToString(res[i]))
		}
	}
```

## Notes


## License
...

**Free for any use**

[//]: #
   [Go compiler]: <https://go.dev>
 
