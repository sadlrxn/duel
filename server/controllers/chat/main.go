package chat

import (
	"encoding/json"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

type Controller struct {
	activeUsers  syncmap.Map
	maxCount     int
	chatContents []types.ChatContent
	currentIndex int
	maxLength    uint
	wagerLimit   int64
	rainMinWager int64
	chatCooldown int
	isMuted      syncmap.Map
	Index        uint
	EventEmitter chan types.WSEvent
}

func (c *Controller) Init(maxCount int, maxLength uint, wagerLimit int64, chatCooldown int, rainMinWager int64) {
	c.maxCount = maxCount
	c.maxLength = maxLength
	c.wagerLimit = wagerLimit
	c.chatCooldown = chatCooldown
	c.isMuted = syncmap.Map{}
	c.activeUsers = syncmap.Map{}
	c.rainMinWager = rainMinWager
}

func (c *Controller) ServeChatContents(conn *websocket.Conn, users []uint, viewerID *uint) {
	length := len(c.chatContents)
	contents := []types.ChatContent{}
	for i := c.currentIndex; i < c.currentIndex+length; i++ {
		contents = append(contents, c.chatContents[i%length])
	}
	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Chat),
		EventType: "messages",
		Payload: gin.H{
			"contents":     contents,
			"maxLength":    c.maxLength,
			"wagerLimit":   c.wagerLimit,
			"chatCooldown": c.chatCooldown,
			"tipMaxAmount": config.TIP_MAX_AMOUNT,
			"tipMinAmount": config.TIP_MIN_AMOUNT,
			"commands":     config.CHAT_COMMANDS,
		},
	})
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}

	var activeUsers []types.User
	var viewer *models.User
	db := db.GetDB()
	if viewerID != nil {
		db.First(&viewer, viewerID)
	}
	for _, userID := range users {
		var user models.User
		db.First(&user, userID)
		activeUsers = append(activeUsers, utils.GetUserDataWithPermissions(user, viewer, 0))
		c.activeUsers.Store(userID, true)
	}
	b, _ = json.Marshal(types.WSMessage{Room: string(types.Chat), EventType: "active_users", Payload: activeUsers})
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}
}

func (c *Controller) RecieveMessage(userID uint, content string) {
	var messageParam struct {
		ReplyTo *uint  `json:"replyTo"`
		Message string `json:"message"`
	}
	json.Unmarshal([]byte(content), &messageParam)
	if c.filterMessage(messageParam.Message) {
		return
	}

	db := db.GetDB()
	var user models.User
	if err := db.Preload("Statistics").Preload("Wallet.Balance.ChipBalance").First(&user, userID).Error; err != nil {
		return
	}

	if user.Role != models.AdminRole && user.Role != models.ModeratorRole {
		if user.Banned {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrBanned}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
		if user.Statistics.TotalWagered < c.wagerLimit {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrLocked}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
		if muted, prs := c.isMuted.Load(userID); prs && muted.(bool) {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrMuted}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
	}

	var newMessage string
	var ok bool
	if newMessage, ok = c.detectAndHandleCommand(user, messageParam.Message); !ok {
		log.LogMessage("invalid command from "+user.Name, messageParam.Message, "error", logrus.Fields{})
		return
	}

	if len(messageParam.Message) > int(c.maxLength) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Exceeds maximum letters"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	c.Index++
	chatContent := types.ChatContent{ID: c.Index, Author: utils.GetUserDataWithPermissions(user, nil, 0), Message: newMessage, Time: uint64(time.Now().UnixMilli()), ReplyTo: messageParam.ReplyTo, Sponsors: []uint{}}
	log.LogMessage("chat from "+user.Name, newMessage, "info", logrus.Fields{"user": userID})
	c.ChatBroadcastMessage("message", chatContent)
}

func (c *Controller) ChatBroadcastMessage(eventType string, chatContent types.ChatContent) {
	if eventType == "message" {
		length := len(c.chatContents)
		if length == c.currentIndex {
			c.chatContents = append(c.chatContents, chatContent)
		} else {
			c.chatContents[c.currentIndex] = chatContent
		}
		c.currentIndex = (c.currentIndex + 1) % c.maxCount
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Chat),
		EventType: eventType,
		Payload:   chatContent,
	})
	c.EventEmitter <- types.WSEvent{Room: types.Chat, Message: b}
}

func (c *Controller) ActivateUser(userID uint) {
	db := db.GetDB()
	var user models.User
	db.First(&user, userID)

	c.activeUsers.Store(userID, true)
	c.activeUsers.Range(func(key, value any) bool {
		var viewer models.User
		db.First(&viewer, key.(uint))
		b, _ := json.Marshal(types.WSMessage{Room: string(types.Chat), EventType: "active_user", Payload: utils.GetUserDataWithPermissions(user, &viewer, 0)})
		c.EventEmitter <- types.WSEvent{Room: types.Chat, Users: []uint{viewer.ID}, Message: b}
		return true
	})
}

func (c *Controller) DeactivateUser(userID uint) {
	db := db.GetDB()
	var user models.User
	db.First(&user, userID)

	c.activeUsers.Delete(userID)
	c.activeUsers.Range(func(key, value any) bool {
		var viewer models.User
		db.First(&viewer, key.(uint))
		b, _ := json.Marshal(types.WSMessage{Room: string(types.Chat), EventType: "deactive_user", Payload: utils.GetUserDataWithPermissions(user, &viewer, 0)})
		c.EventEmitter <- types.WSEvent{Room: types.Chat, Users: []uint{viewer.ID}, Message: b}
		return true
	})
}

func (c *Controller) DeleteMessage(userID uint, content string) {
	db := db.GetDB()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return
	}

	if user.Role != models.AdminRole && user.Role != models.ModeratorRole {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return
	}

	var params struct {
		ID uint `json:"id"`
	}
	json.Unmarshal([]byte(content), &params)

	for index, content := range c.chatContents {
		if content.ID == params.ID {
			c.chatContents[index].Deleted = true
			b, _ := json.Marshal(types.WSMessage{Room: string(types.Chat), EventType: "delete_message", Payload: content.ID})
			c.EventEmitter <- types.WSEvent{Room: types.Chat, Message: b}
			return
		}
	}
}

func (c *Controller) SponsorMessage(userID uint, content string) {
	db := db.GetDB()
	var user models.User
	if err := db.Preload("Statistics").First(&user, userID).Error; err != nil {
		return
	}

	if user.Role != models.AdminRole && user.Role != models.ModeratorRole {
		if user.Banned {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrBanned}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
		if user.Statistics.TotalWagered < c.wagerLimit {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrLocked}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
		if muted, prs := c.isMuted.Load(userID); prs && muted.(bool) {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Chat),
				EventType: "error",
				Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrMuted}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			return
		}
	}

	var params struct {
		ID uint `json:"id"`
	}
	json.Unmarshal([]byte(content), &params)

	for index, content := range c.chatContents {
		if content.ID == params.ID {
			var newSponsors []uint
			var flag bool
			for _, sponsorID := range content.Sponsors {
				if sponsorID == userID {
					flag = true
					continue
				}
				newSponsors = append(newSponsors, sponsorID)
			}
			if !flag {
				newSponsors = append(newSponsors, userID)
			}
			c.chatContents[index].Sponsors = newSponsors
			b, _ := json.Marshal(types.WSMessage{Room: string(types.Chat), EventType: "sponsor_message", Payload: gin.H{
				"id":       content.ID,
				"sponsors": newSponsors,
			}})
			c.EventEmitter <- types.WSEvent{Room: types.Chat, Message: b}
			return
		}
	}
}

func (c *Controller) IsMuted(userID uint) bool {
	muted, prs := c.isMuted.Load(userID)
	if prs && muted.(bool) {
		return true
	}
	return false
}
