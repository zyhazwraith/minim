package main

import (
	//	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		log.Print("handle ", tcpConn.RemoteAddr().String())
		go handleConn(tcpConn)
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
		log.Print("received ", data)
		log.Print("length: ", n)
	}
}

func heartBeat(conn *net.TCPConn, readchan chan byte) {
	for ok := 1; ok < 20; ok++ {
		select {
		case fk := <-readchan:
			log.Print("received ", fk)
			conn.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
			log.Print("xu zhu le")
		case <-time.After(time.Second * 5):
			log.Print("nothing receied in 5 secs")
			conn.Close()
			//			ok = false
		}
	}
}

func gravelChan(data []byte, messager chan byte) {
	//	log.Print("data len ", len(data))
	for _, v := range data {
		messager <- v
		log.Print("push")
	}
	// if close channel in advance, fk will read `` data
	//	close(messager)
}
