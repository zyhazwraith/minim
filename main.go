package main

import (
	"./server"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	lis, err := net.Listen("tcp", ":9998")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	go http.Serve(lis, nil)
	for {

	}
}
