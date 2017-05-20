package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/smallnest/rpcx"
	_ "github.com/smallnest/rpcx/clientselector"
	"github.com/zyhazwraith/minim/server"
)

func main() {
	/*
		sv1 := &clientselector.ServerPeer{Network: "tcp", Address: "localhost:8989"}
		sv2 := &clientselector.ServerPeer{Network: "tcp", Address: "localhost:18989"}

		servers := []*clientselector.ServerPeer{sv1, sv2}

		s := clientselector.NewMultiClientSelector(servers, rpcx.RandomSelect, 10*time.Second)
	*/
	s := &rpcx.DirectClientSelector{Network: "tcp", Address: "localhost:8989", DialTimeout: 10 * time.Second}

	for i := 0; i < 2; i++ {
		callServer(s)
		time.Sleep(2 * time.Second)
	}

	/*
		s := &rpcx.DirectClientSelector{Network: "tcp", Address: "localhost:8989", DialTimeout: 10 * time.Second}
		client := rpcx.NewClient(s)

		args := &server.Args{7, 8}
		var reply server.Reply
		divCall := client.Go(context.Background(), "Arith.Mul", args, &reply, nil)
		replayCall := <-divCall.Done
		if replayCall.Error != nil {
			log.Fatal("Arith error: ", replayCall.Error)
		} else {
			fmt.Println(reply.C)
		}
	*/
}

func callServer(s rpcx.ClientSelector) {
	log.Print("creating client")
	client := rpcx.NewClient(s)

	args := &server.Args{7, 9}
	var reply server.Reply
	err := client.Call(context.Background(), "Arith.Mul", args, &reply)
	if err != nil {
		log.Fatal("arith error", err)
	} else {

		fmt.Println("arith mul: ", reply.C)
	}
	client.Close()
}
