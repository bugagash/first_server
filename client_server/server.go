package main

import (
	"net"
	"log"
	"strings"
)

const (
	END_BYTES = "\000\001\002\003\004\005"
	PORT = ":8080"
)

var (
	Connections = make(map[net.Conn]bool)
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
	//defer conn.Close()
	Connections[conn] = true
	var (
		buffer = make([]byte, 512)
		message string
	)
	close: for {
		message = ""
		for {
			length, err := conn.Read(buffer)
			if err != nil { break close }
			message += string(buffer[:length])
			if (strings.HasSuffix(message, END_BYTES)) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		log.Println(message)
		for addr := range Connections {
			if addr == conn { continue }
			addr.Write([]byte(message + END_BYTES))
		}
	}
	delete(Connections, conn)
}