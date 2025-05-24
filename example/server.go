package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	iso "github.com/Jesting/kiso"
)

func readMessage(con *net.Conn) ([]byte, error) {
	lenBytes := make([]byte, 4)
	cnt, err := (*con).Read(lenBytes)
	if err != nil {
		return lenBytes, err
	}
	if cnt != 4 {
		return lenBytes, fmt.Errorf("length bytes are less than 4 %d", cnt)
	}

	fmt.Printf("<<<%X\n", lenBytes)

	length := int(lenBytes[0])<<24 | int(lenBytes[1])<<16 | int(lenBytes[2])<<8 | int(lenBytes[3])
	println("Expected message length:", length)
	bytes := make([]byte, length)

	cnt, err = (*con).Read(bytes)
	if err != nil {
		return lenBytes, err
	}
	if cnt != length {
		return lenBytes, fmt.Errorf("message size is not equal to %d", length)
	}

	fmt.Printf("<<<%X\n", bytes)
	fmt.Printf("ASCII<<<%s\n", string(bytes))

	return bytes, nil
}
func writeMessage(con *net.Conn, bytes []byte) {
	println("Calculated message length:", len(bytes))
	var lenBytes = []byte{byte((len(bytes) >> 24) & 0xFF), byte((len(bytes) >> 16) & 0xFF), byte((len(bytes) >> 8) & 0xFF), byte(len(bytes) & 0xFF)}

	fmt.Printf(">>>%X\n", lenBytes)

	(*con).Write(lenBytes)

	fmt.Printf(">>>%X\n", bytes)
	fmt.Printf("ASCII>>>%s\n", string(bytes))
	(*con).Write(bytes)
}

func genRRN() string {
	rrn := make([]byte, 12)

	for i := 0; i < 12; i++ {
		rrn[i] = byte('0' + (rand.Int() % 10))
	}
	return string(rrn)
}

func composeResponse(req *iso.Message, feldsEcho []int, feldsGen []int, feldsVal map[int]string) *iso.Message {
	var resp = &iso.Message{}
	resp.SetMti(req.GetMti() + 10)

	for _, fieldNo := range feldsEcho {
		if req.GetField(fieldNo) != nil {
			resp.SetField(req.GetField(fieldNo))
		}
	}
	for _, fieldNo := range feldsGen {
		if fieldNo == 37 {
			resp.SetField(isod.MakeFieldAscii(fieldNo, genRRN()))
		} else {
			resp.SetField(isod.MakeSampleField(fieldNo))
		}
	}

	for no, v := range feldsVal {
		resp.SetField(isod.MakeFieldAscii(no, v))
	}
	return resp
}

func exchange(con *net.Conn, feldsEcho []int, feldsGen []int, feldsVal map[int]string) {
	defer func() {
		(*con).Close()
		println("Connection closed.")
	}()

	bytes, err := readMessage(con)
	if err != nil {
		println("Error receiving message:", err)
		return
	}

	message, err := isod.ParseToMessage(bytes)
	if err != nil {
		println("Error parsing message:", err)
		return
	}

	println("Request message:")
	println(isod.MessageToString(message))

	var responseMessage = composeResponse(message, feldsEcho, feldsGen, feldsVal)

	println("Response message:")
	println(isod.MessageToString(responseMessage))

	responseMessageBytes, err := isod.ComposeFromMessage(responseMessage)
	if err != nil {
		println("Error composing message:", err)
		return
	}

	writeMessage(con, responseMessageBytes)
	println("Message sent.")

}

func server(addr string, input string) {
	var fields = strings.Split(input, ",")

	feldsEcho := []int{}
	feldsGen := []int{}
	feldsVal := make(map[int]string)
	for i := 0; i < len(fields); i++ {
		noAndValue := strings.Split(fields[i], ":")
		var fieldNo, err = strconv.Atoi(noAndValue[0])
		if err != nil {
			println("Wrong field number:", fields[i])
			return
		}
		if len(noAndValue) == 1 {
			feldsEcho = append(feldsEcho, fieldNo)
		} else if len(noAndValue) == 2 {
			if len(noAndValue[1]) == 0 {
				feldsGen = append(feldsGen, fieldNo)
			} else {
				feldsVal[fieldNo] = noAndValue[1]
			}
		} else {
			println("Wrong field format:", fields[i])
			return
		}
	}

	println("Starting server on: ", addr)
	tcp, err := net.Listen("tcp", addr)

	if err != nil {
		println("Listen error:", err)
		return
	}
	println("Listening...")
	for {
		con, err := tcp.Accept()
		if err != nil {
			println("Connection error:", err)
		} else {
			println("Accepted connection:", con.RemoteAddr().String())
			go exchange(&con, feldsEcho, feldsGen, feldsVal)
		}
	}
}
