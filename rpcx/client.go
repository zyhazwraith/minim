package main

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx"
)

type Args struct {
	a int
}

type Reply struct {
	b int
}

type Arith int

func main() {
	server := &rpcx.DirectClientSelector{
		Network: "tcp",
		Address: "127.0.0.1:12345",
	}
	client := rpcx.NewClient(server)
	args := &Args{1}
	var reply Reply
	reply.b = 4
	err := client.Call(context.Background(), "Arith.Test", args, &reply)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(reply.b)
	}
}
