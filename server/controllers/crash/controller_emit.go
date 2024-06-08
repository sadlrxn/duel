package crash

import (
	"encoding/json"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
/* @Internal
/* Broadcasts real time events.
/*
/* The event includes:
/*		- 'RoundID`: current round's ID.
/*		- `Status`: current round status.
/* 		- `CashIns`: performed cash-ins during previous period.
/* 		- `CashOuts`: performed cash-outs during previous period.
/*		- `Elapsed`: elapsed time from the last status updated.
/*		- `Multiplier`: multiplier which should be reached during
/*			the next period.
/*		- `Timestamp`: current time on server side.
*/
func (c *GameController) emitRealTimeEvent() error {
	c.cashMut.Lock()

	event := c.buildRealTimeEvent()
	c.clearPerformedBets()

	c.cashMut.Unlock()

	b, err := json.Marshal(types.WSMessage{
		Room:      string(types.Crash),
		EventType: "game_status",
		Payload:   event,
	})
	if err != nil {
		return utils.MakeError(
			"CrashEmit",
			"emitRealTimeEvent",
			"failed to build event message.",
			err,
		)
	}
	c.EventEmitter <- types.WSEvent{Room: types.Crash, Message: b}
	return nil
}

/*
/* @Internal
/* Emits `balance_update` event when any cash-in fails.
/* Check `balanceType` and sets related feild with converted as chip
/* amount.
*/
func (c *GameController) emitRefundEvent(event CashInEvent) error {
	b, err := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     event.Amount,
			BalanceType: event.BalanceType,
			Delay:       0,
		}})
	if err != nil {
		return utils.MakeError(
			"CrashEmit",
			"emitRefundEvent",
			"failed to build event message.",
			err,
		)
	}
	c.EventEmitter <- types.WSEvent{Users: []uint{event.UserID}, Message: b}
	return nil
}

/*
/* @External
/* Emits current round data via websocket connection.
*/
func (c *GameController) EmitRoundData(conn *websocket.Conn) error {
	var roundPayload RoundPayload
	if c.round != nil {
		roundPayload = RoundPayload{
			RoundID:      c.round.ID,
			RoundStatus:  c.roundStatus,
			Elapsed:      time.Now().UnixMilli() - c.lastStatusUpdated.UnixMilli(),
			Multiplier:   c.currentMultiplier,
			RunStartedAt: c.round.RunStartedAt,
		}
		var betPayloads []BetPayload
		for _, bet := range c.round.Bets {
			betUser := user.GetUserInfoByID(bet.UserID)
			if betUser == nil {
				return utils.MakeError(
					"crash_controller_emit",
					"EmitRoundData",
					"failed to get bet user model",
					nil,
				)
			}
			var profit int64
			if bet.Profit != nil {
				profit = *bet.Profit
			}
			betPayloads = append(betPayloads, BetPayload{
				BetID:            bet.ID,
				BetAmount:        bet.BetAmount,
				PaidBalanceType:  bet.PaidBalanceType,
				Profit:           &profit,
				PayoutMultiplier: bet.PayoutMultiplier,
				CashOutAt:        bet.CashOutAt,
				User:             utils.GetUserDataWithPermissions(*betUser, nil, 0),
			})
		}
		roundPayload.Bets = betPayloads
	}

	currentRoundID := config.CRASH_SEED_CHAIN_LENGTH
	if c.round != nil {
		currentRoundID = c.round.ID
	}

	b, err := json.Marshal(types.WSMessage{
		Room:      string(types.Crash),
		EventType: "game_data",
		Payload: gin.H{
			"round":   roundPayload,
			"history": getRoundHistory(currentRoundID),
		},
	})
	if err != nil {
		return utils.MakeError(
			"CrashEmit",
			"EmitRoundData",
			"failed to build event message.",
			err,
		)
	}
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}
	return nil
}

/*
/* @Internal
/* Returns real time event to be broadcasted from current round.
*/
func (c *GameController) buildRealTimeEvent() RealTimeEvent {
	return RealTimeEvent{
		RoundID:     c.round.ID,
		RoundStatus: c.roundStatus,
		Elapsed:     time.Now().UnixMilli() - c.lastStatusUpdated.UnixMilli(),
		Multiplier:  c.nextMultiplier,
		CashIns:     c.performedCashIns,
		CashOuts:    c.performedCashOuts,
	}
}

/*
/* @Internal
/* Empty `performedCashIns` & `performedCashOuts` temp arrays.
*/
func (c *GameController) clearPerformedBets() {
	c.performedCashIns = []CashInForRealTimeEvent{}
	c.performedCashOuts = []CashOutForRealTimeEvent{}
}
