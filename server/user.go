package main

import (
	"log"
)

type User int
type Message int

type Args struct {
	Username string
	Password string
}

type Reply struct {
	Code   int
	Status bool
	Msg    string
}

type Msg struct {
	Username string
	Content  string
	RoomId   int
}

func (u *User) Login(args *Args, reply *Reply) error {
	var id int32
	var password, salt string

	row := db.QueryRow("select id,  password, salt from user where username = ?", args.Username)
	err := row.Scan(&id, &password, &salt)
	if password == args.Password {
		reply.Status = true
	} else {
		reply.Status = false
	}
	log.Print(args.Username, " login ", reply.Status)
	return err
}

func (u *User) Broad(args *Msg, reply *Reply) error {
	numLock.Lock()
	numPush++
	numLock.Unlock()
	wch <- *args
	return nil
}
