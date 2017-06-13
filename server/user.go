package main

import ()

type User int

type Args struct {
	username string
	password string
}

type Reply struct {
	code   int
	status bool
	msg    string
}

func (u *User) Login(args *Args, reply *Reply) error {
	reply.status = false
	return nil
	/*
		var id int32
		var password, salt string
		//	db.Query("select * from user")
		reply.status = true
		reply.msg = args.username + args.password
		return nil

		row := db.QueryRow("select id,  password, salt from user where username = ?", args.username)
		err := row.Scan(&id, &password, &salt)
		if password == args.password {
			reply.status = true
		} else {
			reply.status = false
		}
	*/
	//	return err
}
