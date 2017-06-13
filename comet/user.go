package main

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
