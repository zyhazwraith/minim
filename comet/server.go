package main

import (
	"context"
	"fmt"
	"net"
	//	"strconv"
	json "github.com/bitly/go-simplejson"
	"github.com/smallnest/rpcx"
	"github.com/zyhazwraith/minim/proto"
	"log"
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
	rch    chan []byte
	wch    chan []byte
	dch    chan bool
	u      string
	roomId int
	conn   *net.TCPConn
}

func newCS(uid string, roomId int, conn *net.TCPConn) *CS {
	return &CS{
		rch:    make(chan []byte),
		wch:    make(chan []byte),
		dch:    make(chan bool),
		u:      uid,
		roomId: roomId,
		conn:   conn,
	}
}

var cmap map[string]*CS
var rclient *rpcx.Client
var rserver *rpcx.DirectClientSelector
var roomNum map[int]int

func main() {
	// start rpcx service first
	rserver = &rpcx.DirectClientSelector{
		Network:     "tcp",
		Address:     "localhost:8972",
		DialTimeout: 10 * time.Second,
	}
	cmap = make(map[string]*CS)
	roomNum = make(map[int]int)
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8989")
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen: ", err)
	}
	//	go pushAll()
	go serve(listen)
	// also serve as rpcx server
	cometServer := rpcx.NewServer()
	cometServer.RegisterName("Push", new(Push))
	log.Println("start serve as rpcx server")
	go cometServer.Serve("tcp", "127.0.0.1:9001")
	time.Sleep(1 * time.Hour)
}

func serve(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("accept :", err)
		}
		log.Println("connected: ", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var client *CS
	// auth first
	var authSuc bool
	authSuc = false
	for {
		conn.Read(data)
		if data[0] != '#' {
			break
		}
		//		fmt.Println("client: ", string(data))
		jsonData, _ := proto.UnpackTcp(data)
		js, _ := json.NewJson(jsonData)
		op, _ := js.Get("Op").Int()
		body := js.Get("Body")
		username, _ := body.Get("Username").String()
		password, _ := body.Get("Password").String()
		roomId, _ := body.Get("RoomId").Int()
		if op == proto.REQ_REG {
			args := &Args{username, password}
			var reply Reply
			rclient = rpcx.NewClient(rserver)
			rclient.Call(context.Background(), "User.Login", args, &reply)
			feedback := proto.FeedBack{true, ""}
			if reply.Status == true {
				log.Print(args.Username, " auth success")
				uid = conn.RemoteAddr().String()
				// add client map
				client = newCS(uid, roomId, conn)
				cmap[client.u] = client
				// add roomNum
				if _, ok := roomNum[roomId]; ok {
					roomNum[roomId]++
				} else {
					roomNum[roomId] = 1
				}
				log.Print(username, " entered ", roomId, " now this comet has ", roomNum[roomId])
				authSuc = true
			} else {
				feedback = proto.FeedBack{false, ""}
				//				conn.Write([]byte{REQ, '#', 's', 'b'})
				authSuc = false
			}
			message := proto.Message{proto.STAT, feedback}
			data, _ := proto.PackTcp(message)
			conn.Write(data)
		}
		break
	}
	if authSuc == false {
		log.Println(conn.RemoteAddr().String(), " disconnect")
		conn.Close()
		return
	}
	log.Println(conn.RemoteAddr().String(), " auth success")
	go writeHandle(conn, client)
	go readHandle(conn, client)
	go work(client)
	select {
	case <-client.dch:
		fmt.Println("close handler goroutine")
		roomNum[client.roomId]--
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
		data := make([]byte, 1024)
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return
		}
		n, err := conn.Read(data)
		if err != nil {
		}
		if n != 0 && data[0] == '#' {
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
			log.Println(client.u, " heartbeat time out")
			client.dch <- true
			return
		}
	}
}

func handleMsg(client *CS, data []byte) {
	jsBody, _ := proto.UnpackTcp(data)
	js, _ := json.NewJson(jsBody)
	op, _ := js.Get("Op").Int()
	//	fmt.Println(op == proto.REQ)
	if op != proto.REQ {
		return
	}
	danmu := js.Get("Body")
	username, _ := danmu.Get("Username").String()
	content, _ := danmu.Get("Content").String()
	roomid, _ := danmu.Get("RoomId").Int()
	args := &Msg{username, content, roomid}
	var reply Reply
	rclient = rpcx.NewClient(rserver)
	rclient.Call(context.Background(), "User.Broad", args, &reply)
}

/*
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
*/
func getJson(data []byte) *json.Json {
	js, _ := json.NewJson(data)
	return js
}

func preMsg(username string, content string) []byte {
	body := proto.Danmu{username, content, 0}
	data, _ := proto.PackTcp(proto.Message{proto.REQ, body})
	return data
}
