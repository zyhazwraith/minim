package main

import (
	"github.com/smallnest/rpcx"
)

type Args struct {
	a int
}

type Reply struct {
	b int
}

type Arith int

func (t *Arith) Test(args *Args, reply *Reply) error {
	reply.b = args.a * 2
	return nil
}

func main() {
	server := rpcx.NewServer()
	server.RegisterName("Arith", new(Arith))
	server.Serve("tcp", "127.0.0.1:12345")
}
