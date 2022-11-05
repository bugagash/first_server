package main

import (
	"fmt";
	"os";
	"net";
	"bufio";
	"strings"
)

const (
	ADDR_SERVER = ":8080"
	END_BYTES = "\000\001\002\003\004\005"
)

func main() {
	conn, err := net.Dial("tcp", ADDR_SERVER)
	if err != nil {
		panic("Can't connect to server")
	}
	defer conn.Close()
	conn.Write([]byte(InputString() + END_BYTES))
	var (
		buffer = make([]byte, 512)
		message string
	)
	for {
		length, err := conn.Read(buffer)
		if (length == 0 || err != nil) { break }
		message = string(buffer[:length])
	}
	fmt.Println(message)
}

func InputString() {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}