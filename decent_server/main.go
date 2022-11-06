package main

import (
	"os"
	"bufio"
	"fmt"
	"log"
	"strings"
	"net"
	"encoding/json"
)

/*
	N - N1, N2, N3
	
	N1 - N
	N2 - N
	N3 - N
*/

type Node struct {
	Connections map[Address]bool
	Address 	Address
}

type Address struct {
	IPv4 string
	Port string
}

func (addr *Address) to_str() string {
	st := addr.IPv4 + ":" + addr.Port
	return st
}

type Package struct {
	To 		Address
	From 	Address
	data 	string
}

func init() {
	if len(os.Args) != 2 {
		panic("Not correct initialization!")
	}
}

func main() {
	NewNode(os.Args[1]).Run(handleServer, handleClient)
}

func NewNode(address string) *Node {
	splited := strings.Split(address, ":")
	if (len(splited) != 2) { return nil }
	return &Node{
		Connections: make(map[Address]bool),
		Address: Address{
			IPv4: splited[0],
			Port: ":"+splited[1],
		},
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(node)
	handleClient(node)
}

func handleServer(node *Node) {
	listen, err := net.Listen("tcp", "0.0.0.0" + node.Address.Port)
	if err != nil {
		panic("Can't listen to server!")
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil { break }
		go handleConnection(node, conn)
	}
}

func handleConnection(node *Node, conn net.Conn) {
	defer conn.Close()
	var (
		buffer =  make([]byte, 512)
		message string
		pack Package
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil { break }
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil { return }
	node.ConnectTo([]string{pack.From.to_str()})
	fmt.Println(pack.From, "--", pack.data)
}

func handleClient(node *Node) {
	for {
		message := InputString()
		splited := strings.Split(message, " ")
		switch splited[0] {
			case "/exit": os.Exit(0)
			case "/connect": node.ConnectTo(splited[1:])
			case "/network": node.PrintNetwork()
			default: node.SendMessageToAll(message)
		}
	}
}

func (node* Node) PrintNetwork() {
	for addr, _ := range node.Connections {
		log.Println("|", addr)
	}
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		split := strings.Split(addr, ":")
		_addr = Address{
			IPv4 = split[0],
			Port = split[1],
		}
		node.Connections[addr] = true
	}
}

func (node *Node) SendMessageToAll(message string) {
	var new_package = Package{
		From: Address{
			IPv4: node.Address.IPv4,
			Port: node.Address.Port,
		},
		data: message,
	}
	for addr := range node.Connections {
		new_package.To = addr
		node.Send(&new_package)
	}
}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To.to_str())
	if err != nil {
		delete(node.Connections, pack.To)
		return
	}
	defer conn.Close()
	json_pack, err := json.Marshal(*pack)
	if err != nil { return }
	conn.Write(json_pack)
}

func InputString() string {
	msg,_ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
