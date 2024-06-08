package crash

import (
	"errors"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/wager"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

/*
/* @Internal
/* 1. Save round's `BetStartedAt` as current time.
/* 2. Set the game status as `crash-status-betting`.
/* 3. Send the first `crash-status-betting` event with remaining betting time as
/*    `bettingDuration`.
/*
/* This function is called at the first time of server starting, and
/* after `preparingDuration` since `crash-status-preparing` status updated while making sure
/* that all cashouts are performed and new round is initialized successfully.
/*
/* Returns error object in case of:
  - round is nil. `ErrCodeEmptyRound`
  - ((betStartedAt, runStartedAt, or endedAt) isn't nil),
    feeTx isn't nil or len(bets) > 0. `ErrCodeNotPreparingRound`
  - Current status is not `crash-status-preparing`. `ErrCodeNotPreparingStatus`
*/
func (c *GameController) startBetting() error {
	if c.round == nil {
		return utils.MakeErrorWithCode(
			"crash status",
			"startBetting",
			"round is nil",
			ErrCodeEmptyRound,
			nil,
		)
	}

	if !c.isPreparingRound() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startBetting",
			"round is not preparing",
			ErrCodeNotPreparingRound,
			nil,
		)
	}

	if !c.isPreparingStatus() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startBetting",
			"status is not preparing",
			ErrCodeNotPreparingStatus,
			nil,
		)
	}

	// 1. Save round's `BetStartedAt` as current time.
	if err := c.updateRoundBetStartedAt(); err != nil {
		return utils.MakeError(
			"crash status",
			"startBetting",
			"failed to update round bet started at.",
			err,
		)
	}

	// 2. Set the game status as `crash-status-betting`.
	c.updateRoundStatus(Betting)

	c.currentMultiplier = 0
	c.nextMultiplier = 1
	c.currentStep = 0

	// 3. Send the first `crash-status-betting` event with remaining betting time as
	//    `bettingDuration`.
	if err := c.emitRealTimeEvent(); err != nil {
		return utils.MakeError(
			"crash status",
			"startBetting",
			"failed to emit the first betting event.",
			err,
		)
	}

	log.LogMessage(
		"crash_controller_status",
		"started betting",
		"info",
		logrus.Fields{
			"time":        time.Now(),
			"round":       c.round,
			"roundStatus": c.roundStatus,
		},
	)

	return nil
}

/*
/* @Internal
/* 1. Insert last event with zero userID to `cashInEvents`.
/* 2. Set the game status as `crash-status-pending`.
/* 3. Send `crash-status-pending` event.
/*
/* This function is called after `bettingDuration` since `crash-status-betting` status updated.
/* After this function, while the status is `crash-status-pending`, no cashIns are accepted.
/*
/* Returns error object in case of:
  - round is nil. `ErrCodeEmptyRound`
  - `betStartedAt` is nil or `runStartedAt` is not nil. `ErrCodeNotBettingRound`
  - Current status is not `crash-status-betting`. `ErrCodeNotBettingStatus`
*/
func (c *GameController) startPending() error {
	if c.round == nil {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPending",
			"round is nil",
			ErrCodeEmptyRound,
			nil,
		)
	}

	if !c.isBettingRound() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPending",
			"round is not betting",
			ErrCodeNotBettingRound,
			nil,
		)
	}

	if !c.isBettingStatus() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPending",
			"status is not betting",
			ErrCodeNotBettingStatus,
			nil,
		)
	}

	// 1. Insert last event with zero userID to `cashInEvents`.
	c.insertEmptyCashInEvent()

	// 2. Set the game status as `crash-status-pending`.
	c.updateRoundStatus(Pending)

	// 3. Send `crash-status-pending` event.
	if err := c.emitRealTimeEvent(); err != nil {
		return utils.MakeError(
			"crash status",
			"startPending",
			"failed to emit the first pending event.",
			err,
		)
	}

	log.LogMessage(
		"crash_controller_status",
		"started pending",
		"info",
		logrus.Fields{
			"time":        time.Now(),
			"round":       c.round,
			"roundStatus": c.roundStatus,
		},
	)

	return nil
}

/*
/* @Internal
/* 1. Init `currentMultiplier`=0, `nextMultiplier`=1, `currentStep`=0.
/* 2. Save round's `runStartedAt`.
/* 3. Set the game status as `crash-status-running`.
/* 4. Send the first running event with 0 elapsed and `nextMultiplier`.
/*
/* This function is called after `pendingDuration` since `crash-status-pending` status updated while
/* making sure that no events in `cashInEvents`.
/* After this function, while the status is `crash-status-running`, cashOuts are accepted.
/*
/* Returns error object in case of:
  - round is nil. `ErrCodeEmptyRound`
  - `endedAt` is not nil. `ErrCodeNotPendingRound`
  - Current status is not `crash-status-pending`. `ErrCodeNotPendingStatus`
*/
func (c *GameController) startRunning() error {
	if c.round == nil {
		return utils.MakeErrorWithCode(
			"crash status",
			"startRunning",
			"round is nil",
			ErrCodeEmptyRound,
			nil,
		)
	}

	if !c.isPendingRound() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startRunning",
			"round is not pending",
			ErrCodeNotPendingRound,
			nil,
		)
	}

	if !c.isPendingStatus() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startRunning",
			"status is not pending",
			ErrCodeNotPendingStatus,
			nil,
		)
	}

	// 1. Init `currentMultiplier`=0, `nextMultiplier`=1, `currentStep`=0.
	c.currentMultiplier = 0
	c.nextMultiplier = 1
	c.currentStep = 0

	// 2. Save round's `runStartedAt`.
	if err := c.updateRoundRunStartedAt(); err != nil {
		return utils.MakeError(
			"crash status",
			"startRunning",
			"failed to update round run started at.",
			err,
		)
	}

	// 3. Set the game status as `crash-status-running`.
	c.updateRoundStatus(Running)

	// 4. Send the first running event with 0 elapsed and `nextMultiplier`.
	if err := c.emitRealTimeEvent(); err != nil {
		return utils.MakeError(
			"crash status",
			"startRunning",
			"failed to emit the first running event.",
			err,
		)
	}

	log.LogMessage(
		"crash_controller_status",
		"started running",
		"info",
		logrus.Fields{
			"time":        time.Now(),
			"round":       c.round,
			"roundStatus": c.roundStatus,
		},
	)

	return nil
}

/*
/* @Internal
/* 1. Insert last event with zero userID to `cashOutEvents`.
/* 2. Set the game status as `crash-status-preparing`.
/* 3. Send the crash event with remaining cashOuts if exist.
/* 4. Transfer house fee to feeWallet.
/*
/* This function is called if `nextMultiplier` >= `round.Outcome` and
/* current status is `crash-status-running`.
/*
/* Returns error object incase of:
  - round is nil. `ErrCodeEmptyRound`
  - `startedAt` is nil or `endedAt` is not nil. `ErrCodeNotRunningRound`
  - Current status is not `crash-status-running`. `ErrCodeNotRunningStatus`
*/
func (c *GameController) startPreparing() error {
	if c.round == nil {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPreparing",
			"round is nil",
			ErrCodeEmptyRound,
			nil,
		)
	}

	if !c.isRunningRound() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPreparing",
			"round is not running",
			ErrCodeNotRunningRound,
			nil,
		)
	}

	if !c.isRunningStatus() {
		return utils.MakeErrorWithCode(
			"crash status",
			"startPreparing",
			"status is not running",
			ErrCodeNotRunningStatus,
			nil,
		)
	}

	// 1. Insert last event with zero userID to `cashOutEvents`.
	c.insertEmptyCashOutEvent()

	// 2. Set the game status as `crash-status-preparing`.
	c.updateRoundStatus(Preparing)

	// 3. Send the crash event with remaining cashOuts if exist.
	if err := c.emitRealTimeEvent(); err != nil {
		return utils.MakeError(
			"crash status",
			"startPreparing",
			"failed to emit the first preparing event.",
			err,
		)
	}

	log.LogMessage(
		"crash_controller_status",
		"started preparing",
		"info",
		logrus.Fields{
			"time":        time.Now(),
			"round":       c.round,
			"roundStatus": c.roundStatus,
		},
	)

	// 4. Transfer house fee to feeWallet.
	return nil
}

/*
/* @Internal
/* 1. Save round's `endedAt` as current time.
/* 2. Load next round to `round`.
/* 3. Initializes `performedCashOuts` and `performedCashIns`.
/*
/* This function is called if there is not evet in `cashOutEvents` channel and,
/* current game status is `crash-status-preparing`.
/*
/* Returns error object in ccase of:
  - round is nil. `ErrCodeEmptyRound`
  - Current status is not `crash-status-preparing`. `ErrCodeNotPreparingStatus`
*/
func (c *GameController) prepareBetting() error {
	if c.round == nil {
		return utils.MakeErrorWithCode(
			"crash status",
			"prepareBetting",
			"round is nil",
			ErrCodeEmptyRound,
			nil,
		)
	}

	if !c.isPreparingStatus() {
		return utils.MakeErrorWithCode(
			"crash status",
			"prepareBetting",
			"status is not preparing",
			ErrCodeNotPreparingStatus,
			nil,
		)
	}

	// 1. Save round's `endedAt` as current time.
	if err := c.updateRoundEndedAt(); err != nil {
		return utils.MakeError(
			"crash status",
			"prepareBetting",
			"failed to update round ended at",
			err,
		)
	}

	c.updatePlayerStatistics()

	// 2. Transfer fee from temp to fee wallet.
	charged, err := c.chargeFee()
	if err != nil {
		log.LogMessage(
			"crash_controller_status",
			"failed to charge fee",
			"error",
			logrus.Fields{
				"round": c.round,
				"error": err.Error(),
			},
		)
	} else {
		log.LogMessage(
			"crash_controller_status",
			"successfully charged fee",
			"success",
			logrus.Fields{
				"round":   c.round,
				"charged": charged,
			},
		)
	}

	// 3. Initializes `performedCashOuts` and `performedCashIns`.
	c.cashMut.Lock()

	c.clearPerformedBets()

	c.cashMut.Unlock()

	// 4. Initializes `cashInCountsPerUser` and `cashOutFlagPerBet`.
	c.clearFlagMaps()

	// 5. Check is blocked.
	if c.isBlockCrash {
		c.eventTicker.Stop()
		c.round = nil
		return utils.MakeError(
			"crash status",
			"prepareBetting",
			"crash is blocked by admin",
			errors.New("isBlockCrash flag is set"),
		)
	}

	// 6. Load next round to `round`.
	if err := c.loadNextRound(); err != nil {
		return utils.MakeError(
			"crash status",
			"prepareBetting",
			"failed to load next round.",
			err,
		)
	}

	log.LogMessage(
		"crash_controller_status",
		"preparing betting",
		"info",
		logrus.Fields{
			"time":        time.Now(),
			"round":       c.round,
			"roundStatus": c.roundStatus,
		},
	)

	return nil
}

/*
/* @Internal
/* Checks whether the round is preparing.
*/
func (c *GameController) isPreparingRound() bool {
	return c.round != nil &&
		c.round.BetStartedAt == nil &&
		c.round.RunStartedAt == nil &&
		c.round.EndedAt == nil &&
		c.round.FeeTransaction == nil &&
		len(c.round.Bets) == 0
}

/*
/* @Internal
/* Checks whether the round is betting.
*/
func (c *GameController) isBettingRound() bool {
	return c.round != nil &&
		c.round.BetStartedAt != nil &&
		c.round.RunStartedAt == nil &&
		c.round.EndedAt == nil
}

/*
/* @Internal
/* Checks whether the round is pending.
*/
func (c *GameController) isPendingRound() bool {
	return c.round != nil &&
		c.round.BetStartedAt != nil &&
		c.round.BetStartedAt.Before(time.Now().Add(-c.bettingDuration)) &&
		c.round.RunStartedAt == nil &&
		c.round.EndedAt == nil
}

/*
/* @Internal
/* Checks whether the round is running.
*/
func (c *GameController) isRunningRound() bool {
	return c.round != nil &&
		c.round.BetStartedAt != nil &&
		c.round.RunStartedAt != nil &&
		c.round.EndedAt == nil
}

/*
* @Internal
* Checks whether the round is ended round.
 */
func (c *GameController) isEndedRound() bool {
	return c.round != nil &&
		c.round.BetStartedAt != nil &&
		c.round.RunStartedAt != nil &&
		c.round.EndedAt != nil
}

/*
/* @Internal
/* Checks whether game status is `crash-status-preparing`.
*/
func (c *GameController) isPreparingStatus() bool {
	return c.roundStatus == Preparing
}

/*
/* @Internal
/* Checks whether game status is `crash-status-betting`.
*/
func (c *GameController) isBettingStatus() bool {
	return c.roundStatus == Betting
}

/*
/* @Internal
/* Checks whether game status is `crash-status-pending`.
*/
func (c *GameController) isPendingStatus() bool {
	return c.roundStatus == Pending
}

/*
/* @Internal
/* Checks whether game status is `crash-status-running`.
*/
func (c *GameController) isRunningStatus() bool {
	return c.roundStatus == Running
}

/*
/* @Internal
/* Updates current round's status & set `lastStatusUpdated` as current time.
/* And should lock `statusMut`.
*/
func (c *GameController) updateRoundStatus(status GameStatus) {
	c.statusMut.Lock()

	c.roundStatus = status
	switch status {
	case Betting:
		c.lastStatusUpdated = *c.round.BetStartedAt
	case Pending:
		c.lastStatusUpdated = time.Now()
	case Running:
		c.lastStatusUpdated = *c.round.RunStartedAt
	case Preparing:
		c.lastStatusUpdated = time.Now()
	}

	c.statusMut.Unlock()
}

/*
/* @Internal
/* Insert an empty cash-in event with `UserID` as 0, to detect
/* whether all cash-ins performed.
*/
func (c *GameController) insertEmptyCashInEvent() {
	c.CashIn(CashInEvent{
		UserID:  0,
		RoundID: c.round.ID,
		Amount:  c.minBetAmount,
	})
}

/*
/* @Internal
/* Insert an empty cash-out event with `UserID` as 0, to detect
/* whether all cash-outs performed.
*/
func (c *GameController) insertEmptyCashOutEvent() {
	c.CashOut(CashOutEvent{
		UserID:  0,
		RoundID: c.round.ID,
	})
}

/*
/* @Internal
/* Update user statistics after round ended.
*/
func (c *GameController) updatePlayerStatistics() {
	params := wager.PerformAfterWagerParams{
		Players: []wager.PlayerInPerformAfterWagerParams{},
		Type:    models.Crash,
	}
	for _, bet := range c.round.Bets {
		if bet.PaidBalanceType != models.ChipBalanceForGame {
			continue
		}
		playerInPerformAfterWager := wager.PlayerInPerformAfterWagerParams{
			UserID: bet.UserID,
			Bet:    bet.BetAmount,
		}
		if bet.Profit != nil && bet.PayoutMultiplier != nil {
			playerInPerformAfterWager.Profit = *bet.Profit
		}
		params.Players = append(
			params.Players,
			playerInPerformAfterWager,
		)
	}
	if err := wager.AfterWager(params); err != nil {
		log.LogMessage(
			"crash_controller_status_updatePlayerStatistics",
			"failed to update wager status",
			"error",
			logrus.Fields{
				"error":  err.Error(),
				"params": params,
			},
		)
	}
}
