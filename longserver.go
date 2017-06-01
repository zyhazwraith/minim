package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr := net.ResolveTCPAddr("tcp", "localhost:8989")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Print(err)
		}
		tcpConn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	}
}

func handleConn(conn net.TCPConn, timeout int) {

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			log.Fatal(err)
		}

	}
}
