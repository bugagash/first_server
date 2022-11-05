package main

import (
	"net"
	"strings"
)

const (
	END_BYTES = "\000\001\002\003\004\005"
	PORT = ":8080"
)

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		panic("Server error!")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil { break }
		go handleConnect(conn)
	}
}

func handleConnect(conn net.Conn) {
	defer conn.Close()
	var (
		buffer = make([]byte, 512)
		message string
	)
	for {
		length, err := conn.Read(buffer)
		if (length == 0 || err != nil) { break }
		message += string(buffer[:length])
		if (strings.HasSuffix(message, END_BYTES)) {
			message = strings.TrimSuffix(message, END_BYTES)
			break
		}
	}
	conn.Write([]byte(strings.ToUpper(message)))
}