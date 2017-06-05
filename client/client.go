package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func sender(conn net.Conn) {
	name := os.Args[1]
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(":")
		input, _ := reader.ReadBytes('\n')
		//fmt.Print(input)
		conn.Write([]byte(name + ": " + string(input)))
		//	fmt.Println("send over")
	}
}

func reader(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		fmt.Println()
		fmt.Print(string(buffer[:n]))
		fmt.Print(":")
	}
}
func main() {
	server := "127.0.0.1:8989"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	if os.Args[2] != "123456" {
		fmt.Println("log fail")
		os.Exit(1)
	}
	if len(os.Args[1]) < 4 {
		fmt.Println("log fail")
		os.Exit(1)
	}
	if os.Args[1][:4] != "test" {
		fmt.Println("log fail")
		os.Exit(1)
	}
	fmt.Println("logging success")
	go sender(conn)
	go reader(conn)
	for {

	}
}
