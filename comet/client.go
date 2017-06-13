package main

import (
	"bufio"
	"fmt"
	json "github.com/bitly/go-simplejson"
	"github.com/zyhazwraith/minim/proto"
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
var username, password string
var room int

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

	var res bool
	fmt.Print("username: ")
	fmt.Scanf("%s", &username)
	fmt.Print("password: ")
	fmt.Scanf("%s", &password)
	fmt.Print("room: ")
	fmt.Scanf("%d", &room)
	for {
		myinfo := proto.RegInfo{username, password, room}
		message := proto.Message{proto.REQ_REG, myinfo}
		data, _ := proto.PackTcp(message)
		//		fmt.Println(string(data))
		conn.Write(data)
		conn.Read(data)
		body, _ := proto.UnpackTcp(data)
		//		conn.Write([]byte{REQ_REG, '#', '2'})
		//		conn.Read(data)
		js, err := json.NewJson(body)
		if err != nil {
			return
		}
		res, err = js.Get("Body").Get("Status").Bool()
		break
	}
	if res == false {
		fmt.Println("disconnect")
		conn.Close()
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
		danmu := proto.Danmu{username, string(data)}
		message := proto.Message{proto.REQ, danmu}
		data, _ = proto.PackTcp(message)
		wch <- data
	}
}

func readHandle(conn *net.TCPConn) {
	for {
		data := make([]byte, 1024)
		n, _ := conn.Read(data)
		if n == 0 {
			dch <- true
			return
		}
		if data[0] == '#' {
			rch <- data
		}
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
			recvMsg(data)
		case <-ticker.C:
			//send heartbeat
			//			wch <- []byte{REQ_HB, '#', 'h', 'b'}
			message := proto.Message{proto.REQ_HB, ""}
			heartbeat, _ := proto.PackTcp(message)
			wch <- heartbeat
		}
	}
}

func recvMsg(data []byte) {
	jsBody, _ := proto.UnpackTcp(data)
	js, _ := json.NewJson(jsBody)
	op, _ := js.Get("Op").Int()
	if op != proto.REQ {
		return
	}
	danmu := js.Get("Danmu")
	uname, _ := danmu.Get("username").String()
	content, _ := danmu.Get("content").String()
	if uname != username {
		fmt.Print(username, ":", content)
	}
}
