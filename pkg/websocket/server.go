package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WsServer struct {
	upgrade websocket.Upgrader
}

func NewWsServer() *WsServer {
	ws := new(WsServer)
	ws.upgrade = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return ws
}

func (ws *WsServer) Run() {
	http.HandleFunc("/ws", ws.wsHandle)
	if err := http.ListenAndServe(":8688", nil); err != nil {
		panic(err)
	}
}

func (ws *WsServer) wsHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		log.Printf("msg: %s", string(msg))
	}
}
