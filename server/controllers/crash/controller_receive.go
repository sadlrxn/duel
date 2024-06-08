package crash

import (
	"runtime/debug"

	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

// Check current status `crash-status-betting`
func (c *GameController) CashIn(event CashInEvent) {
	if c.isBettingStatus() &&
		c.isValidCashInEvent(event) &&
		c.checkAndIncreaseCashInCountsPerUser(event.UserID) {
		c.cashInEvents <- event
	} else {
		log.LogMessage(
			"crash CashIn",
			"cash-in blocked",
			"info",
			logrus.Fields{
				"userID":                              event.UserID,
				"isBettingStatus":                     c.isBettingStatus(),
				"isValidCashInEvent":                  c.isValidCashInEvent(event),
				"checkAndIncreaseCashInCountsPerUser": c.checkAndIncreaseCashInCountsPerUser(event.UserID),
			},
		)
		if err := c.emitRefundEvent(event); err != nil {
			log.LogMessage(
				"crash CashIn",
				"an error occured during refund failed cash-in.",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
		}
	}
}

// Check current status `crash-status-running`
func (c *GameController) CashOut(event CashOutEvent) {
	if event.PayoutMultiplier == 0 {
		event.PayoutMultiplier = c.currentMultiplier
	}
	if c.isRunningStatus() &&
		c.isValidCashOutEvent(event) &&
		c.checkAndSetCashOutFlagPerBet(event.BetID) {
		c.cashOutEvents <- event
	}
}

/*
/* @internal
/* Pop cash event from channels and call handler.
*/
func (c *GameController) eventListener() {
	defer func() {
		if r := recover(); r != nil {
			log.LogMessage(
				"crash",
				"event listener recovered",
				"info",
				logrus.Fields{
					"recover": r,
					"stack":   string(debug.Stack()),
				},
			)
		}
	}()

	for {
		select {
		case cashInEvent := <-c.cashInEvents:
			log.LogMessage(
				"crash_event_listener",
				"poped cash-in event",
				"info",
				logrus.Fields{
					"event":   cashInEvent,
					"isEmpty": c.isEmptyCashIn(cashInEvent),
				},
			)
			if !c.isEmptyCashIn(cashInEvent) {
				c.cashInHandler(cashInEvent)
			}
		case cashOutEvent := <-c.cashOutEvents:
			log.LogMessage(
				"crash_event_listener",
				"poped cash-out event",
				"info",
				logrus.Fields{"event": cashOutEvent},
			)
			if !c.isEmptyCashOut(cashOutEvent) {
				c.cashOutHandler(cashOutEvent)
			}
		}
	}
}

/*
	@Internal

/* Check whether `CashInEvent` is valid.
*/
func (c *GameController) isValidCashInEvent(event CashInEvent) bool {
	return (event.CashOutAt == 0 || event.CashOutAt >= c.minCashOutAt) &&
		event.Amount >= c.minBetAmount &&
		event.Amount <= c.maxBetAmount &&
		c.round != nil &&
		event.RoundID == c.round.ID
}

/*
	@Internal

/* Check whether `CashOutEvent` is valid.
*/
func (c *GameController) isValidCashOutEvent(event CashOutEvent) bool {
	return event.PayoutMultiplier <= c.nextMultiplier &&
		event.PayoutMultiplier >= c.minCashOutAt &&
		c.round != nil &&
		event.RoundID == c.round.ID
}

// Check current status `crash-status-betting`, and `crash-status-pending`
// CashMut
func (c *GameController) cashInHandler(event CashInEvent) {
	betID, err := c.cashIn(CashInRequestParams(event))
	if err != nil {
		c.decreaseCashInCountsPerUser(event.UserID)
		log.LogMessage(
			"crash cashInHandler",
			"an error occured during cash-in.",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		if err := c.emitRefundEvent(event); err != nil {
			log.LogMessage(
				"crash cashInHandler",
				"an error occured during refund failed cash-in.",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
		}
		return
	} else {
		log.LogMessage(
			"crash_cashInHandler",
			"successfully cashed in",
			"success",
			logrus.Fields{
				"event": event,
			},
		)
	}

	log.LogMessage(
		"crash_cash_in_handler",
		"cash-in performed",
		"success",
		logrus.Fields{"betID": betID},
	)

	user := user.GetUserInfoByID(event.UserID)
	if user == nil {
		return
	}

	c.cashMut.Lock()
	defer c.cashMut.Unlock()
	c.performedCashIns = append(
		c.performedCashIns,
		CashInForRealTimeEvent{
			User:        utils.GetUserDataWithPermissions(*user, nil, 0),
			Amount:      event.Amount,
			BalanceType: event.BalanceType,
			BetID:       betID,
			CashOutAt:   event.CashOutAt,
		})
	log.LogMessage(
		"crash_cash_in_handler",
		"insert performed result",
		"success",
		logrus.Fields{"betID": betID},
	)
}

// Check current status `crash-status-running`, and `crash-status-preparing`
// CashMut
func (c *GameController) cashOutHandler(event CashOutEvent) {
	result, err := c.cashOut(CashOutRequestParams(event))
	if err != nil || result == nil {
		c.resetCashOutFlagPerBet(event.BetID)
		log.LogMessage(
			"crash cashOutHandler",
			"an error occured during cash-out.",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return
	} else {
		log.LogMessage(
			"crash_cashOutHandler",
			"successfully cashed out",
			"success",
			logrus.Fields{
				"event": event,
			},
		)
	}

	user := user.GetUserInfoByID(event.UserID)
	if user == nil {
		return
	}

	c.cashMut.Lock()
	defer c.cashMut.Unlock()
	c.performedCashOuts = append(
		c.performedCashOuts,
		CashOutForRealTimeEvent{
			User:        utils.GetUserDataWithPermissions(*user, nil, 0),
			Amount:      result.Amount,
			BalanceType: result.BalanceType,
			Multiplier:  event.PayoutMultiplier,
			BetID:       event.BetID,
		})
}

/*
/* @Internal
/* Returns whether the cash-in event is empty.
*/
func (c *GameController) isEmptyCashIn(event CashInEvent) bool {
	return event.UserID == 0
}

/*
/* @Internal
/* Returns whether the cash-out event is empty.
*/
func (c *GameController) isEmptyCashOut(event CashOutEvent) bool {
	return event.UserID == 0
}

/*
/* @Internal
/* Increase `cashInCountsPerUser` by 1.
*/
func (c *GameController) checkAndIncreaseCashInCountsPerUser(userID uint) bool {
	value, ok := c.cashInCountsPerUser.Load(userID)
	if ok && value.(uint) >= c.betCountLimit {
		return false
	} else if ok && value.(uint) < c.betCountLimit {
		c.cashInCountsPerUser.Store(userID, value.(uint)+1)
		return true
	} else if !ok {
		c.cashInCountsPerUser.Store(userID, uint(1))
		return true
	}
	return false
}

/*
/* @Internal
/* Decrease `cashInCoutsPerUser` by 1.
*/
func (c *GameController) decreaseCashInCountsPerUser(userID uint) {
	value, ok := c.cashInCountsPerUser.Load(userID)
	if ok && value.(uint) > 0 {
		c.cashInCountsPerUser.Store(userID, value.(uint)-1)
	}
}

/*
/* @Internal
/* Set `cashOutFlagPerBet`.
*/
func (c *GameController) checkAndSetCashOutFlagPerBet(betID uint) bool {
	value, ok := c.cashOutFlagPerBet.Load(betID)
	if ok && value.(bool) == true {
		return false
	} else if ok && value.(bool) == false {
		c.cashOutFlagPerBet.Store(betID, true)
		return true
	} else if !ok {
		c.cashOutFlagPerBet.Store(betID, true)
		return true
	}
	return false
}

/*
/* @Internal
/* Reset `cashOutFlagPerBet`.
*/
func (c *GameController) resetCashOutFlagPerBet(betID uint) {
	value, ok := c.cashOutFlagPerBet.Load(betID)
	if ok && value.(bool) == true {
		c.cashOutFlagPerBet.Store(betID, false)
	}
}

/*
/* @Internal
/* Clear `cashInCountsPerUser` and `cashOutFlagPerBet`.
*/
func (c *GameController) clearFlagMaps() {
	c.cashInCountsPerUser.Range(func(key, value any) bool {
		c.cashInCountsPerUser.Delete(key)
		return true
	})
	c.cashOutFlagPerBet.Range(func(key, value any) bool {
		c.cashOutFlagPerBet.Delete(key)
		return true
	})
}
