package socket

import (
	"encoding/json"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitConnection(ctx *gin.Context, hub *Hub, userID *uint) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.LogMessage("websocket", "failed to connect websocket", "error", logrus.Fields{})
		return
	}

	client := &Client{hub: hub, conn: ws, send: make(chan []byte, 256), room: types.None, userID: userID}
	client.hub.register <- client

	var keys []uint
	hub.users.Range(func(key, value any) bool {
		keys = append(keys, key.(uint))
		return true
	})
	controllers.Chat.ServeChatContents(ws, keys, userID)

	b, _ := json.Marshal(types.WSMessage{
		EventType: "connect",
	})
	client.hub.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{ws}, Message: b}
	go client.Reader()
	go client.Writer()

	if userID != nil {
		log.LogMessage("websocket", "websocket started", "info", logrus.Fields{"user": *userID})
	} else {
		log.LogMessage("websocket", "websocket started", "info", logrus.Fields{})
	}
}
