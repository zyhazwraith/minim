package main

import (
	"fmt"
	"net"
	//	"strconv"
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

type CS struct {
	rch  chan []byte
	wch  chan []byte
	dch  chan bool
	u    string
	conn *net.TCPConn
}

func newCS(uid string, conn *net.TCPConn) *CS {
	return &CS{
		rch:  make(chan []byte),
		wch:  make(chan []byte),
		dch:  make(chan bool),
		u:    uid,
		conn: conn,
	}
}

var cmap map[string]*CS

func main() {
	cmap = make(map[string]*CS)
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen: ", err)
	}
	//	go pushAll()
	serve(listen)
	for {
	}
}

func pushAll() {
	time.Sleep(15 * time.Second)
	for _, v := range cmap {
		fmt.Println("push to all")
		v.wch <- []byte{REQ, '#', 'p', 'u', 's', 'h'}
	}
}

func serve(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept :", err)
		}
		fmt.Println("connected: ", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var client *CS
	// auth first
	for {
		conn.Read(data)
		fmt.Println("client: ", string(data))
		if data[0] == REQ_REG {
			conn.Write([]byte{RES_REG, '#', 'o', 'k'})
			uid = conn.RemoteAddr().String()
			client = newCS(uid, conn)
			cmap[client.u] = client
			break
		} else {
			conn.Write([]byte{RES_REG, '#', '0'})
		}
	}
	fmt.Println("finish auth")
	go writeHandle(conn, client)
	go readHandle(conn, client)
	go work(client)
	select {
	case <-client.dch:
		fmt.Println("close handler goroutine")
	}
}

func writeHandle(conn *net.TCPConn, client *CS) {
	//	tick := time.NewTicker(20 * time.Second)
	for {
		select {
		case d := <-client.wch:
			conn.Write(d)
			//		case <-tick.C:
			//			if _, ok := cmap[client.u]; !ok {
			//				fmt.Println("conn die, close writehandle")
			//				return
			//			}
		}
	}
}

func readHandle(conn *net.TCPConn, client *CS) {
	for {
		data := make([]byte, 128)
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return
		}
		n, err := conn.Read(data)
		if err != nil {
		}
		if n != 0 {
			client.rch <- data
		}
	}
}

func work(client *CS) {
	//	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case data := <-client.rch:
			handleMsg(client, data)
		case <-time.After(time.Second * 60):
			delete(cmap, client.u)
			fmt.Println("heartbeat time out")
			client.dch <- true
			return
		}
	}
}

func handleMsg(client *CS, data []byte) {
	if data[0] == REQ {
		// msg recv log
		//		fmt.Println("recv msg, send ack")
		client.wch <- []byte{RES, '#'}
		// notice that slice is a reference
		username := []byte(client.conn.RemoteAddr().String())
		// wrong  use
		//			newdata := append(data[:2], username...)
		newdata := []byte{REQ, '#'}
		newdata = append(newdata, username...)
		newdata = append(newdata, []byte(" say: ")...)
		newdata = append(newdata, data[2:]...)
		fmt.Print(string(newdata[2:]))
		for _, v := range cmap {
			if v != client {
				v.wch <- newdata
			}
		}
	} else if data[0] == REQ_HB {
		//heart beat log
		//		fmt.Println("recv ht, send ack")
		client.wch <- []byte{RES_HB, '#', 'h', 'b'}
	}
	// res && res_hb do not need feedback
}
