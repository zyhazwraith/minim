package main

import (
	server "./server"
	"fmt"
	"log"
	"net/rpc"
)

const (
	servAddr = "localhost"
)

func main() {
	client, err := rpc.DialHTTP("tcp", servAddr+":9998")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	args := &server.Args{7, 0}
	/*
		var reply int
		err = client.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error: ", err)
		}
		fmt.Printf("Arith: %d * %d = %d\n", args.A, args.B, reply)
	*/
	quotient := server.Args{}
	divCall := client.Go("Arith.Quotient", args, &quotient, nil)
	replyCall := <-divCall.Done
	if replyCall.Error != nil {
		log.Fatal("arith error: ", replyCall.Error)
	}
	fmt.Println(replyCall.Reply)
}
