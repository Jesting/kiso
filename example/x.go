package main

import (
	"fmt"
	"os"

	iso "github.com/Jesting/kiso/parser"
)

var isod = iso.MakeIsoDescription(iso.LengthFormat_ASCII, iso.Fields93)

func main() {
	println(len(os.Args))
	println(os.Args[1])
	if len(os.Args) == 4 && os.Args[1] == "client" {
		client(os.Args[2], os.Args[3])
	} else if len(os.Args) == 4 && os.Args[1] == "server" {
		server(os.Args[2], os.Args[3])
	} else {
		fmt.Println("No/wrong arguments provided")
		fmt.Println("Usage: ./x client ip:port mti:xxxx,fieldNo:value,....")
		fmt.Println("Usage: ./x server ip:port echfieldNo,generatedfieldNo:,valuesfieldNo:value,....")
		return
	}
}
