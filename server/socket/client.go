package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/coinflip"
	"github.com/Duelana-Team/duelana-v1/controllers/crash"
	"github.com/Duelana-Team/duelana-v1/controllers/jackpot"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	room types.Room

	userID *uint
}

type message struct {
	MsgType string `json:"type"`
	Room    string `json:"room"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

func (c *Client) listenCoinflip(content string) error {
	var params []coinflip.EventParam
	err := json.Unmarshal([]byte(content), &params)
	if err != nil {
		return err
	}

	for _, eventParam := range params {
		if admin.GetGameBlocked(admin.GAME_CONTROLLER_COINFLIP) {
			return fmt.Errorf("coinflip blocked by admin")
		}
		if eventParam.EventType == "bet" {
			if eventParam.Opponent == coinflip.Bot {
				controllers.Coinflip.BetAgainstBot(*c.userID, eventParam)
			} else {
				if eventParam.RoundID == nil {
					controllers.Coinflip.Create(*c.userID, eventParam)
				} else {
					controllers.Coinflip.Join(*c.userID, *eventParam.RoundID)
				}
			}
		} else if eventParam.EventType == "cancel" {
			controllers.Coinflip.Cancel(*c.userID, *eventParam.RoundID)
		}
	}
	return nil
}

func (c *Client) listenJackpotLow(content string) {
	if admin.GetGameBlocked(admin.GAME_CONTROLLER_JACKPOT) {
		return
	}
	var betParam struct {
		Amount int      `json:"amount"`
		Nfts   []string `json:"nfts"`
	}
	json.Unmarshal([]byte(content), &betParam)

	nftAmount, nfts := user.GetNftDetailsFromMintAddresses(betParam.Nfts)

	betData := jackpot.BetData{
		Amount:    int64(betParam.Amount),
		NftAmount: nftAmount,
		Nfts:      nfts,
		Time:      time.Now(),
	}

	if betData.Amount+betData.NftAmount == 0 {
		log.LogMessage(string(c.room)+" websocket reader", "invalid bet data", "error", logrus.Fields{"user": *c.userID})
		return
	}

	controllers.JackpotLow.Bet(*c.userID, jackpot.BetData{Amount: betData.Amount, NftAmount: betData.NftAmount, Nfts: betData.Nfts, Time: time.Now()})
}

func (c *Client) listenJackpotMedium(content string) {
	if admin.GetGameBlocked(admin.GAME_CONTROLLER_JACKPOT) {
		return
	}
	var betParam struct {
		Amount int      `json:"amount"`
		Nfts   []string `json:"nfts"`
	}
	json.Unmarshal([]byte(content), &betParam)

	nftAmount, nfts := user.GetNftDetailsFromMintAddresses(betParam.Nfts)

	betData := jackpot.BetData{
		Amount:    int64(betParam.Amount),
		NftAmount: nftAmount,
		Nfts:      nfts,
		Time:      time.Now(),
	}

	if betData.Amount+betData.NftAmount == 0 {
		log.LogMessage(string(c.room)+" websocket reader", "invalid bet data", "error", logrus.Fields{"user": *c.userID})
		return
	}

	controllers.JackpotMedium.Bet(*c.userID, jackpot.BetData{Amount: betData.Amount, NftAmount: betData.NftAmount, Nfts: betData.Nfts, Time: time.Now()})
}

func (c *Client) listenJackpotWild(content string) {
	if admin.GetGameBlocked(admin.GAME_CONTROLLER_JACKPOT) {
		return
	}
	var betParam struct {
		Amount int      `json:"amount"`
		Nfts   []string `json:"nfts"`
	}
	json.Unmarshal([]byte(content), &betParam)

	nftAmount, nfts := user.GetNftDetailsFromMintAddresses(betParam.Nfts)

	betData := jackpot.BetData{
		Amount:    int64(betParam.Amount),
		NftAmount: nftAmount,
		Nfts:      nfts,
		Time:      time.Now(),
	}

	if betData.Amount+betData.NftAmount == 0 {
		log.LogMessage(string(c.room)+" websocket reader", "invalid bet data", "error", logrus.Fields{"user": *c.userID})
		return
	}

	controllers.JackpotWild.Bet(*c.userID, jackpot.BetData{Amount: betData.Amount, NftAmount: betData.NftAmount, Nfts: betData.Nfts, Time: time.Now()})
}

func (c *Client) listenGrandJackpot(content string) {
	if admin.GetGameBlocked(admin.GAME_CONTROLLER_GRAND_JACKPOT) {
		return
	}
	var betParam struct {
		Amount int      `json:"amount"`
		Nfts   []string `json:"nfts"`
	}
	json.Unmarshal([]byte(content), &betParam)

	nftAmount, nfts := user.GetNftDetailsFromMintAddresses(betParam.Nfts)

	betData := jackpot.BetData{
		Amount:    int64(betParam.Amount),
		NftAmount: nftAmount,
		Nfts:      nfts,
		Time:      time.Now(),
	}

	if betData.Amount+betData.NftAmount == 0 {
		log.LogMessage(string(c.room)+" websocket reader", "invalid bet amount", "error", logrus.Fields{"user": *c.userID})
		return
	}

	controllers.GrandJackpot.Bet(*c.userID, jackpot.BetData{Amount: betData.Amount, NftAmount: betData.NftAmount, Nfts: betData.Nfts, Time: time.Now()})
}

func (c *Client) listenCrash(content string) error {
	if admin.GetGameBlocked(admin.GAME_CONTROLLER_CRASH) {
		return errors.New("crash game blocked by admin.")
	}

	var event struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	err := json.Unmarshal([]byte(content), &event)
	if err != nil {
		return utils.MakeError(
			"websocket reader",
			"listenCrash",
			"failed to unmarshal event param.",
			err,
		)
	}

	switch event.Type {
	case "cash-in":
		var cashInEvent crash.CashInEvent
		err := json.Unmarshal([]byte(event.Content), &cashInEvent)
		if err != nil {
			return utils.MakeError(
				"websocket reader",
				"listenCrash",
				"failed to unmarshal cash in event.",
				err,
			)
		}
		cashInEvent.UserID = *c.userID
		controllers.Crash.CashIn(cashInEvent)
	case "cash-out":
		var cashOutEvent crash.CashOutEvent
		err := json.Unmarshal([]byte(event.Content), &cashOutEvent)
		if err != nil {
			return utils.MakeError(
				"websocket reader",
				"listenCrash",
				"failed to unmarshal cash out event.",
				err,
			)
		}
		cashOutEvent.UserID = *c.userID
		controllers.Crash.CashOut(cashOutEvent)
	default:
		return utils.MakeError(
			"websocket reader",
			"listenCrash",
			"invalid event type for crash listener.",
			err,
		)
	}
	return nil
}

func (c *Client) Reader() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.LogMessage(string(c.room)+" websocket reader", "unexpectedly closed", "error", logrus.Fields{"error": err.Error()})
			}
			break
		}

		var message message
		json.Unmarshal(msg, &message)

		ok, err := middlewares.WebsocketRateLimiter(
			message.MsgType+"/"+message.Room,
			c.userID,
			c.conn,
		)
		if err != nil {
			log.LogMessage(
				"websocket reader",
				"occured an error during rate limit",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
		}

		if ok {
			switch message.MsgType + message.Room {
			case "visit" + string(types.Coinflip):
				c.room = types.Coinflip
				go controllers.Coinflip.ServeGameData(c.conn)
			case "visit" + string(types.Jackpot):
				c.room = types.Jackpot
				go controllers.JackpotLow.ServeRoundData(c.conn)
				go controllers.JackpotMedium.ServeRoundData(c.conn)
				go controllers.JackpotWild.ServeRoundData(c.conn)
			case "visit" + string(types.GrandJackpot):
				c.room = types.GrandJackpot
				go controllers.GrandJackpot.ServeRoundData(c.conn)
			case "visit" + string(types.None):
				c.room = types.None
			case "visit" + string(types.Crash):
				c.room = types.Crash
				go controllers.Crash.EmitRoundData(c.conn)
			case "event" + string(types.Coinflip):
				if c.userID != nil {
					go c.listenCoinflip(message.Content)
				}
			case "event" + string(types.Jackpot):
				if c.userID != nil {
					if message.Level == "low" {
						go c.listenJackpotLow(message.Content)
					} else if message.Level == "medium" {
						go c.listenJackpotMedium(message.Content)
					} else if message.Level == "wild" {
						go c.listenJackpotWild(message.Content)
					}
				}
			case "event" + string(types.GrandJackpot):
				if c.userID != nil {
					go c.listenGrandJackpot(message.Content)
				}
			case "event" + string(types.Crash):
				if c.userID != nil {
					go c.listenCrash(message.Content)
				}
			case "message" + string(types.Chat):
				if c.userID != nil {
					go controllers.Chat.RecieveMessage(*c.userID, strings.TrimSpace(message.Content))
				}
			case "reply" + string(types.Chat):
				if c.userID != nil {
					go controllers.Chat.RecieveMessage(*c.userID, strings.TrimSpace(message.Content))
				}
			case "delete" + string(types.Chat):
				if c.userID != nil {
					go controllers.Chat.DeleteMessage(*c.userID, strings.TrimSpace(message.Content))
				}
			case "sponsor" + string(types.Chat):
				if c.userID != nil {
					go controllers.Chat.SponsorMessage(*c.userID, strings.TrimSpace(message.Content))
				}
			}
		}
	}
}

func (c *Client) Writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.LogMessage(string(c.room)+" websocket writer", "failed to get websocket writter", "error", logrus.Fields{"error": err.Error})
				return
			}
			mes, err := buildMessage(message)
			if err != nil {
				log.LogMessage(string(c.room)+" websocket writer", "failed to build message", "error", logrus.Fields{"error": err.Error, "message": message})
				continue
			}
			_, err = w.Write(mes)
			if err != nil {
				log.LogMessage(string(c.room)+" websocket writer", "failed to write message", "error", logrus.Fields{"error": err.Error()})
			}

			// Add queued events to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	event := <-c.send
			// 	mes, err = buildMessage(event)
			// 	if err != nil {
			// 		log.LogMessage(string(c.room)+" websocket writer", "failed to build message", "error", logrus.Fields{"error": err.Error, "message": message})
			// 		continue
			// 	}
			// 	w.Write(mes)
			// }
			if err := w.Close(); err != nil {
				log.LogMessage(string(c.room)+" websocket writer", "failed to close websocket writter", "error", logrus.Fields{"error": err.Error})
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if c.userID != nil {
					log.LogMessage(string(c.room)+" websocket writer", "websocket closed for user", "info", logrus.Fields{"user": *c.userID, "error": err.Error})
				} else {
					log.LogMessage(string(c.room)+" websocket writer", "websocket closed", "info", logrus.Fields{"error": err.Error})
				}
				return
			}
		}
	}
}

// @Interanl
// Build websocket message from event.
func buildMessage(message []byte) (mes []byte, err error) {
	var event interface{}
	if err = json.Unmarshal(message, &event); err != nil {
		return
	}
	var data struct {
		Event interface{} `json:"event"`
		Time  time.Time   `json:"time"`
	}
	data.Event = event
	data.Time = time.Now()
	mes, err = json.Marshal(data)
	return
}
