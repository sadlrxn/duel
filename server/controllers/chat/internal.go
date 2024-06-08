package chat

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) detectAndHandleCommand(user models.User, content string) (string, bool) {
	var command, param1, param2 string
	n, _ := fmt.Sscanf(content, `/%s %s %s`, &command, &param1, &param2)
	if n < 2 {
		return content, true
	}

	switch command {
	case "mute":
		if ok, err := regexp.MatchString(`^/mute\s\w+$`, content); err == nil && ok {
			return content, c.mute(user, param1)
		}
		return content, true
	case "unmute":
		if ok, err := regexp.MatchString(`^/unmute\s\w+$`, content); err == nil && ok {
			return content, c.unmute(user, param1)
		}
		return content, true

	case "ban":
		if ok, err := regexp.MatchString(`^/ban\s\w+$`, content); err == nil && ok {
			return content, c.ban(user, param1)
		}
		return content, true
	case "unban":
		if ok, err := regexp.MatchString(`^/unban\s\w+$`, content); err == nil && ok {
			return content, c.unban(user, param1)
		}
		return content, true
	case "setMaxLength":
		if ok, err := regexp.MatchString(`^/setMaxLength\s\d+$`, content); err == nil && ok {
			return content, c.setMaxLength(user, param1)
		}
		return content, true
	case "setWagerLimit":
		if ok, err := regexp.MatchString(`^/setWagerLimit\s\d+$`, content); err == nil && ok {
			return content, c.setWagerLimit(user, param1)
		}
		return content, true
	case "setChatCooldown":
		if ok, err := regexp.MatchString(`^/setChatCooldown\s\d+$`, content); err == nil && ok {
			return content, c.setChatCooldown(user, param1)
		}
		return content, true
	case "rain":
		if ok, err := regexp.MatchString(`^/rain\s\d+\s\d+$`, content); err == nil && ok {
			return c.rain(user, param1, param2, content)
		}
		return content, true
	default:
		return content, false
	}
}

func (c *Controller) mute(user models.User, param string) bool {
	var target models.User
	db := db.GetDB()
	if err := db.Where("name = ?", param).First(&target).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Invalid userName : " + param}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}

	if !(user.Role == models.AdminRole && target.Role != models.AdminRole ||
		user.Role == models.ModeratorRole && (target.Role != models.AdminRole && target.Role != models.ModeratorRole)) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	c.isMuted.Store(target.ID, true)
	log.LogMessage(user.Name+" has muted", target.Name, "info", logrus.Fields{"handler": user.Role, "target": target.Role})

	timer := time.NewTimer(config.MUTE_DURATION)
	go func() {
		<-timer.C
		c.unmute(user, param)
	}()
	return true
}

func (c *Controller) unmute(user models.User, param string) bool {
	var target models.User
	db := db.GetDB()
	if err := db.Where("name = ?", param).First(&target).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Invalid userName : " + param}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}

	if !(user.Role == models.AdminRole && target.Role != models.AdminRole ||
		user.Role == models.ModeratorRole && (target.Role != models.AdminRole && target.Role != models.ModeratorRole)) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	c.isMuted.Store(target.ID, false)
	log.LogMessage(user.Name+" has unmuted", target.Name, "info", logrus.Fields{"handler": user.Role, "target": target.Role})
	return true
}

func (c *Controller) ban(user models.User, param string) bool {
	var target models.User
	db := db.GetDB()
	if err := db.Where("name = ?", param).First(&target).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Invalid userName : " + param}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}

	if !(user.Role == models.AdminRole && target.Role != models.AdminRole ||
		user.Role == models.ModeratorRole && (target.Role != models.AdminRole && target.Role != models.ModeratorRole)) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	target.Banned = true
	db.Save(&target)
	log.LogMessage(user.Name+" has banned", target.Name, "info", logrus.Fields{"handler": user.Role, "target": target.Role})
	return true
}

func (c *Controller) unban(user models.User, param string) bool {
	var target models.User
	db := db.GetDB()
	if err := db.Where("name = ?", param).First(&target).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Invalid userName : " + param}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}

	if !(user.Role == models.AdminRole && target.Role != models.AdminRole ||
		user.Role == models.ModeratorRole && (target.Role != models.AdminRole && target.Role != models.ModeratorRole)) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	target.Banned = false
	db.Save(&target)
	log.LogMessage(user.Name+" has unbanned", target.Name, "info", logrus.Fields{"handler": user.Role, "target": target.Role})
	return true
}

func (c *Controller) setMaxLength(user models.User, param string) bool {
	if user.Role != models.AdminRole {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	var limit uint
	n, err := fmt.Sscanf(param, "%d", &limit)
	if n == 0 || err != nil {
		return false
	}
	c.maxLength = limit

	msg, _ := json.Marshal(gin.H{"maxLength": limit})
	cmdContent := types.ChatContent{Author: utils.GetUserDataWithPermissions(user, nil, 0), Message: string(msg), Time: uint64(time.Now().UnixMilli())}
	c.ChatBroadcastMessage("command", cmdContent)
	return true
}

func (c *Controller) setWagerLimit(user models.User, param string) bool {
	if user.Role != models.AdminRole {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	var limit int64
	n, err := fmt.Sscanf(param, "%d", &limit)
	if n == 0 || err != nil {
		return false
	}
	c.wagerLimit = utils.ConvertChipToBalance(limit)

	msg, _ := json.Marshal(gin.H{"wagerLimit": limit})
	cmdContent := types.ChatContent{Author: utils.GetUserDataWithPermissions(user, nil, 0), Message: string(msg), Time: uint64(time.Now().UnixMilli())}
	c.ChatBroadcastMessage("command", cmdContent)
	return true
}

func (c *Controller) setChatCooldown(user models.User, param string) bool {
	if user.Role != models.AdminRole {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{Message: "Permission denied"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return false
	}
	var chatCooldown int
	n, err := fmt.Sscanf(param, "%d", &chatCooldown)
	if n == 0 || err != nil {
		return false
	}
	c.chatCooldown = chatCooldown
	msg, _ := json.Marshal(gin.H{"chatCooldown": chatCooldown})
	cmdContent := types.ChatContent{Author: utils.GetUserDataWithPermissions(user, nil, 0), Message: string(msg), Time: uint64(time.Now().UnixMilli())}
	c.ChatBroadcastMessage("command", cmdContent)
	return false
}

func (c *Controller) rain(user models.User, param1 string, param2 string, content string) (string, bool) {
	var split int
	var amount int64
	n, err := fmt.Sscanf(param1, "%d", &split)
	if n != 1 || err != nil {
		return content, true
	}
	n, err = fmt.Sscanf(param2, "%d", &amount)
	if n != 1 || err != nil {
		return content, true
	}

	if split <= 0 || amount <= 0 {
		return content, false
	}

	balance := amount
	if user.Wallet.Balance.ChipBalance.Balance < balance {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrNotEnoughChip}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return content, false
	}

	selectedUserIDs, selectedUsers := c.selectUsersForRain(user.ID, split)
	if len(selectedUserIDs) != split {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Chat),
			EventType: "error",
			Payload:   types.ErrorMessagePayload{ErrorCode: types.ErrNotEnoughPeople}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
		return content, false
	}

	var receipients []db_aggregator.User
	for _, receipient := range selectedUserIDs {
		receipients = append(receipients, db_aggregator.User(receipient))
	}

	balanceForEach := balance / int64(split)
	_, err = transaction.Rain(&transaction.RainRequest{
		FromUser: (*db_aggregator.User)(&user.ID),
		ToUsers:  &receipients,
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &balanceForEach,
		},
		Type: models.TxRain,
	})
	if err != nil {
		log.LogMessage("rain handler", "failed to distribute funds", "error", logrus.Fields{"error": err.Error()})
		return content, false
	}

	b, _ := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Decrease,
			Balance:     balanceForEach * int64(len(receipients)),
			BalanceType: models.ChipBalanceForGame,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}

	b, _ = json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     balanceForEach,
			BalanceType: models.ChipBalanceForGame,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{Users: selectedUserIDs, Message: b}

	var newContent = content + " "
	bytes, _ := json.Marshal(selectedUsers)
	newContent += string(bytes)
	return newContent, true
}

func (c *Controller) selectUsersForRain(userID uint, count int) ([]uint, []types.User) {
	db := db.GetDB()
	var selectedUserIDs []uint
	var selectedUsers []types.User

	recentlyWageredUserIDs := redis.ZRevRangeRecentlyWagered()
	utils.ShuffleSlice(recentlyWageredUserIDs)
	for _, recentUserID := range recentlyWageredUserIDs {
		if len(selectedUserIDs) >= count {
			break
		}
		if _, ok := c.activeUsers.Load(recentUserID); !ok {
			continue
		}
		onlineUser := models.User{}
		if result := db.Preload(
			"Statistics",
		).First(
			&onlineUser,
			recentUserID,
		); result.Error != nil {
			continue
		}
		if onlineUser.Statistics.TotalWagered >= c.rainMinWager &&
			recentUserID != userID {
			selectedUserIDs = append(
				selectedUserIDs,
				recentUserID,
			)
			selectedUsers = append(
				selectedUsers,
				utils.GetUserDataWithPermissions(
					onlineUser,
					nil,
					0,
				),
			)
		}

	}

	return selectedUserIDs, selectedUsers
}

func (c *Controller) filterMessage(content string) bool {
	b, _ := regexp.MatchString(`(\$ \w+ .*)`, content)
	return b
}
