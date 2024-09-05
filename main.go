package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级到WebSocket连接时出错:", err)
		return
	}
	defer conn.Close()
	log.Println("客户端已连接")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("从客户端读取时出错:", err)
			return
		}
		log.Printf("这是客户端发来的信息: %s ------------- 消息类型: %v", string(p), messageType)

		err = conn.WriteMessage(messageType, []byte("你好！我是server端。。。"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/log", handleConnection)
	log.Println("WebSocket服务器正在侦听：1024")
	err := http.ListenAndServe(":1024", nil)
	if err != nil {
		log.Println("WebSocket服务器启动出错:", err)
	}
}
