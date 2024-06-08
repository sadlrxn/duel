package types

import "github.com/gorilla/websocket"

type Room string

const (
	Jackpot      Room = "jackpot"
	Coinflip     Room = "coinflip"
	Payment      Room = "payment"
	GrandJackpot Room = "grandJackpot"
	Chat         Room = "chat"
	Crash        Room = "crash"
	None         Room = ""
)

type WSEvent struct {
	Conns   []*websocket.Conn
	Users   []uint
	Room    Room
	Message []byte
}

type WSMessage struct {
	Room      string      `json:"room"`
	EventType string      `json:"eventType"`
	Payload   interface{} `json:"payload"`
}
