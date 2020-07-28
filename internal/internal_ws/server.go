package internal_ws

import (
	"fmt"
	"golang_im/pkg/log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	// CheckOrigin 检查源，如果 CheckOrigin 函数返回 false，则 Upgrade 方法会使 WebSocket 握手失败，HTTP 状态为403
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级协议
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("连接错误:", err)
		return
	}

	ctx := NewWSConn(Conn, 0, 0, 0)
	ctx.DoConn()
}

// websocket服务路由
func StartWSServer(address string) {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("websocket server start")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Panic(err)
		return
	}
}
