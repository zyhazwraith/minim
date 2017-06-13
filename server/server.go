package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@/meeidol")
	if err != nil {
		return
	}
	server := rpcx.NewServer()
	server.RegisterName("User", new(User))
	fmt.Println("start listen port 8972")
	server.Serve("tcp", "127.0.0.1:8972")
}
