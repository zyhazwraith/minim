package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		handleConn(tcpConn)
	}
}

func handleConn(conn *net.TCPConn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Print(err)
			return
		}
		data := buffer[:n]
		messager := make(chan byte)

		go heartBeat(conn, messager)

		go gravelChan(data, messager)
	}
}

func heartBeat(conn *net.Conn, readchan chan byte) {
	select {
	case fk := <-readchan:
		log.Print("received ", fk)
		conn.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
		log.Print("xu zhu le")
	case <-time.After(time.Second * 5):
		log.Print("nothing receied in 5 secs")
		conn.Close()
	}
}

func gravelChan(data []byte, messager chan byte) {
	for _, v := range data {
		messager <- v
	}
	close(messager)
}
