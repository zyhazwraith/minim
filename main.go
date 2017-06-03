package main

import (
	myserv "github.com/zyhazwraith/minim/server"

	"github.com/smallnest/rpcx"
)

func main() {
	server := rpcx.NewServer()
	server.RegisterName("Arith", new(myserv.Arith))
	server.Start("tcp", "localhost:8989")

	server2 := rpcx.NewServer()
	server2.RegisterName("Arith", new(myserv.Arith2))
	server2.Serve("tcp", "localhost:18989")
}
