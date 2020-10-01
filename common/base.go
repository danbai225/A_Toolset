package common

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"time"
)

type Msg struct {
	Time int64       `json:"time"`
	Type string      `json:"type"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetMsg(Type string, msg string, data interface{}) Msg {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic info is: ", err)
		}
	}()
	return Msg{Time: time.Now().UnixNano() / 1e6, Type: string(Type), Msg: msg, Data: data}
}
func SendStructType(conn net.Conn, Type string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic info is: ", err)
		}
	}()
	SendStruct(conn, Type, "", nil)
}
func SendStructTypeAndData(conn net.Conn, Type string, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic info is: ", err)
		}
	}()
	SendStruct(conn, Type, "", data)
}
func SendStruct(conn net.Conn, Type string, msg string, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic info is: ", err)
		}
	}()
	_, err := conn.Write(BytesCombine(ToJsonBytes(GetMsg(Type, msg, ToJsonStr(data))), []byte("\n")))
	if err != nil {
		println(err)
	}
}
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
func ToJsonBytes(s interface{}) []byte {
	//Person 结构体转换为对应的 Json
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	return jsonBytes
}
func ToJsonStr(s interface{}) string {
	return string(ToJsonBytes(s))
}
