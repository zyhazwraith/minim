package main

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx"
	"time"
)

func main() {
	server := &rpcx.DirectClientSelector{
		Network:     "tcp",
		Address:     "127.0.0.1:8972",
		DialTimeout: 10 * time.Second,
	}
	client := rpcx.NewClient(server)
	args := &Args{"test1", "123456"}
	var reply Reply
	reply.status = true
	err := client.Call(context.Background(), "User.Login", args, &reply)
	if err != nil {
		fmt.Println("User ", err)
	} else {
		//		fmt.Println(reply.msg)
		fmt.Println(reply.status)
	}
	client.Close()
}
