package main

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
