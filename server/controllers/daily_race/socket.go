package daily_race

import (
	"encoding/json"
	"time"

	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
)

var EventEmitter chan types.WSEvent

/*
* @Internal
* Initialize socket event emitter.
 */
func initSocket(eventEmitter chan types.WSEvent) {
	EventEmitter = eventEmitter
}

/**
* @Internal
* Sends prizing websocket events to winners.
 */
func sendPrizingEvents(result *DailyRacePrizingResult) error {
	var resErr error
	for _, winner := range result.Winners {
		if wsEvent, err := buildPrizingEvent(
			result.Date,
			winner,
		); err == nil && wsEvent != nil {
			EventEmitter <- *wsEvent
		} else {
			resErr = utils.MakeError(
				"daily_race_socket",
				"sendPrizingEvents",
				"failed",
				err,
			)
		}
	}
	return resErr
}

/*
* @Internal
* Builds websocket event for winners.
 */
func buildPrizingEvent(date time.Time, winner WinnerInDailyRacePrizingResult) (*types.WSEvent, error) {
	if b, err := json.Marshal(types.WSMessage{
		EventType: "notification",
		Payload: gin.H{
			"date":  date,
			"rank":  winner.Rank,
			"prize": winner.Prize,
		},
	}); err == nil {
		return &types.WSEvent{
			Users:   []uint{winner.UserID},
			Message: b,
		}, nil
	} else {
		return nil, utils.MakeError(
			"daily_race_socket",
			"buildPrizingEvent",
			"failed to marshal json",
			err,
		)
	}
}
