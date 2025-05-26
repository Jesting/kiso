package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	iso "github.com/Jesting/kiso"
)

func composeRequest(input string) (*iso.Message, error) {
	var fields = strings.Split(input, ",")
	var mti = 0
	noAndValue := strings.Split(fields[0], ":")
	if noAndValue[0] == "mti" {
		mti, _ = strconv.Atoi(noAndValue[1])
		if mti == 0 {
			return nil, fmt.Errorf("wrong mti provided")
		}
	} else {
		return nil, fmt.Errorf("mti not found")
	}
	var message = iso.Message{}

	for i := 1; i < len(fields); i++ {
		noAndValue = strings.Split(fields[i], ":")
		if len(noAndValue) != 2 {
			return nil, fmt.Errorf("wrong field format (%s)", fields[i])
		}
		var fieldNo, _ = strconv.Atoi(noAndValue[0])
		message.SetField(isod.MakeFieldAscii(fieldNo, noAndValue[1]))
	}
	message.SetMti(mti)
	return &message, nil
}

func client(addr string, input string) {
	print("Starting client, connecting to: ", addr)
	con, err := net.Dial("tcp", addr)

	if err != nil {
		println("Dial error:", err)
		return
	}
	defer func() {
		con.Close()
		println("Connection closed.")
	}()

	println("Connected.")

	request, err := composeRequest(input)
	if err != nil {
		println("Error building message:", err)
		return
	}
	println("Request message:")
	println(isod.MessageToString(request))
	requestBytes, err := isod.ComposeFromMessage(request)
	if err != nil {
		println("Error composing message:", err)
		return
	}
	writeMessage(&con, requestBytes)

	println("Message sent.")

	responseBytes, err := readMessage(&con)
	if err != nil {
		println("Error receiving message:", err)
		return
	}
	response, err := isod.ParseToMessage(responseBytes)
	if err != nil {
		println("Error parsing message:", err)
		return
	}
	println("Response message:")
	println(isod.MessageToString(response))

}
