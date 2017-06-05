package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var connPool = make([]*net.Conn, 20)
var connCnt int

func main() {

	//建立socket，监听端口
	connCnt = 0
	netListen, err := net.Listen("tcp", "localhost:8989")
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
		connPool[connCnt] = &conn
		connCnt++
	}
	for {

	}
}

func transMessage(conn net.Conn, ch chan string) {
	for {
		msg := <-ch
		for i := 0; i < connCnt; i = i + 1 {
			s1 := conn.RemoteAddr().String()
			s2 := (*connPool[i]).RemoteAddr().String()
			if s1 != s2 {
				(*connPool[i]).Write([]byte(msg))
			}
		}
	}
}

//处理连接
func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)
	ch := make(chan string, 10)
	go transMessage(conn, ch)
	for {

		n, err := conn.Read(buffer)

		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		if n == 0 {
			return
		}
		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		fmt.Print(string(buffer[:n]))
		ch <- string(buffer[:n])
	}

}
func Log(v ...interface{}) {
	return
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
