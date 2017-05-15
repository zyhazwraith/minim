package server

import (
	"errors"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Quotient(args *Args, quo *Args) error {
	if args.B == 0 {
		return errors.New("divided by zero")
	}
	quo.A = args.A / args.B
	quo.B = args.A % args.B
	return nil
}
