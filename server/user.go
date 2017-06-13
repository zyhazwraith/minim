package main

import ()

type User int

type Args struct {
	Username string
	Password string
}

type Reply struct {
	Code   int
	Status bool
	Msg    string
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
	return err
}
