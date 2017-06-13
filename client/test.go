package main

import (
	"fmt"
)

func main() {
	data := []byte{1, '#'}
	data = append(data, []byte("123"))
	fmt.Print(data)
}
