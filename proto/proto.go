package proto

import (
	"encoding/json"
	"fmt"
	"strconv"
	//	simplejson "github.com/bitly/go-simplejson"
)

const (
	REQ_REG int = 1
	RES_REG int = 2
	REQ_HB  int = 3
	RES_HB  int = 4
	REQ     int = 5
	RES     int = 6
	STAT    int = 7
)

type Message struct {
	Op   int
	Body interface{}
}

type FeedBack struct {
	Status bool
	Error  string
}

type RegInfo struct {
	Username string
	Password string
	RoomId   int
}

type Danmu struct {
	Username string
	Content  string
	RoomId   int
}

func PackTcp(data interface{}) ([]byte, error) {
	buff := make([]byte, 1024)
	body, err := json.Marshal(data)
	if err != nil {
	}
	buff = []byte(fmt.Sprintf("%04d", len(body)))
	buff[0] = '#'
	buff = append(buff, body...)
	tail := fmt.Sprintf("%01024d", 1)
	buff = append(buff, tail...)
	buff = buff[:1024]
	return buff, nil
	//	fmt.Print(string(buff), cap(buff), len(buff))
}

func UnpackTcp(data []byte) ([]byte, error) {
	data[0] = '0'
	bodyLen, _ := strconv.Atoi(string(data[:4]))
	body := data[4 : bodyLen+4]
	return body, nil
}
