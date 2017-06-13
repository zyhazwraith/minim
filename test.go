package main

import (
	"fmt"
	"strings"
)

func main() {
	data := string("127.0.0.1:8989")
	fmt.Println(data[len(data)-4:])
}
