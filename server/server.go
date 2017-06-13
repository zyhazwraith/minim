package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx"
	"log"
	"sync"
	"time"
)

var db *sql.DB
var wch chan Msg
var numPush int
var numLock sync.Mutex
var cometAddr []string
var drpcx *rpcx.DirectClientSelector

func main() {
	wch = make(chan Msg)
	numPush = 0
	var err error
	db, err = sql.Open("mysql", "root:root@/meeidol")
	if err != nil {
		return
	}
	log.Print("connect mysql database")
	server := rpcx.NewServer()
	server.RegisterName("User", new(User))
	server.RegisterName("Message", new(Message))
	log.Println("start listen port 8972")
	go server.Serve("tcp", "127.0.0.1:8972")
	go boardMsg()
	// serve as client
	cometAddr = make([]string, 1)
	cometAddr[0] = "127.0.0.1:9001"
	drpcx = &rpcx.DirectClientSelector{
		Network:     "tcp",
		Address:     cometAddr[0],
		DialTimeout: 10 * time.Second,
	}
	time.Sleep(1 * time.Hour)
}

func boardMsg() {
	for {
		select {
		case msg := <-wch:
			log.Print(string(msg.Username))
			numLock.Lock()
			numPush--
			log.Print("broad message: ", msg.Username, " ", msg.Content)
			client := rpcx.NewClient(drpcx)
			client.Call(context.Background(), "Push.Broad", &msg, nil)
			numLock.Unlock()
		}
	}
}
