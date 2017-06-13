package main

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
