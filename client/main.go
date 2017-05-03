package main

import (
	"fmt"
	"net"
	"strconv"
)

const (
	host = "localhost"
	port = 9999
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return
	}

	fmt.Println("connect success")
	//send out message
	conn.Write([]byte("Hello world\n"))
	conn.Close()
}
