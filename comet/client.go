package main

import (
	"bufio"
	"fmt"
	_ "log"
	"net"
	"os"
	"time"
)

const (
	REQ_REG byte = 1
	RES_REG byte = 2
	REQ_HB  byte = 3
	RES_HB  byte = 4
	REQ     byte = 5
	RES     byte = 6
)

var dch chan bool
var rch chan []byte
var wch chan []byte

func main() {
	dch = make(chan bool)
	rch = make(chan []byte)
	wch = make(chan []byte)
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go handleConn(conn)
	select {
	case <-dch:
		fmt.Print("close")
	}
}

func handleConn(conn *net.TCPConn) {
	data := make([]byte, 128)
	for {
		conn.Write([]byte{REQ_REG, '#', '2'})
		conn.Read(data)
		if data[0] == RES_REG {
			break
		}

	}
	fmt.Println("auth finish")
	go readHandle(conn)
	go writeHandle(conn)
	go work()
	go input()
}

func input() {
	for {
		fmt.Print(": ")
		reader := bufio.NewReader(os.Stdin)
		data, _ := reader.ReadBytes('\n')
		data = append([]byte{REQ, '#'}, data...)
		wch <- data
	}
}

func readHandle(conn *net.TCPConn) {
	for {
		data := make([]byte, 128)
		n, _ := conn.Read(data)
		if n == 0 {
			dch <- true
			return
		}
		rch <- data
	}
}

func writeHandle(conn *net.TCPConn) {
	for {
		select {
		case msg := <-wch:
			//			fmt.Println("send data: ", string(msg))
			conn.Write(msg)
		}
	}
}

func work() {
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case data := <-rch:
			//			log.Print("work recv ", string(msg))
			//			wch <- []byte{RES, '#', 'x', 'x'}
			if data[0] == REQ {
				fmt.Println("")
				fmt.Print(string(data[1:]))
				fmt.Print(": ")
				// should send ack here
				// wch <- ack
			}
		case <-ticker.C:
			//send heartbeat
			wch <- []byte{REQ_HB, '#', 'h', 'b'}
		}
	}
}
