package main

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
)

const (
	port = 9999
	host = "localhost"
)

func main() {
	//config
	runtime.GOMAXPROCS(4 - 1)
	//start listener
	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		//		fmt.Println(err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	fmt.Println("start a goroutine to  handle: ", conn.RemoteAddr().String())

	buff := make([]byte, 2048)
	_, err := conn.Read(buff)

	if err != nil {
		return
	}

	fmt.Print("received: ", string(buff))
	fmt.Println("exit")
}
