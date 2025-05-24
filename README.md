# KISO
## _for Payments101_

KISO is simple [ISO8583] parser + some standard's primitives

**Code supports:**
* LLVar, LLLVar, LLLLVar field length
* ASCII, BCD & HEX length prefixes
* ASCII & Binary fields
* MTI composition
* 3 bitmaps
* Padding
* Message output

## To build the app you need:
- [Go compiler]

## Import

```
go get github.com/Jesting/kiso
```

## Use
Example ISO8583 client and server code (powered by KISO) is located under the ./example dierectory

To start client:
```
cd example
go build
./client.sh 
```
To start server:
```
cd example
go build
./server.sh 
```
## Notes

## License
...

**Free for any use**

[//]: #
   [Go compiler]: <https://go.dev>
   [ISO8583]: <https://en.wikipedia.org/wiki/ISO_8583>
 
