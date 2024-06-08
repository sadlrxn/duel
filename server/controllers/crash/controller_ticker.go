package crash

import (
	"math"
	"time"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/sirupsen/logrus"
)

func (c *GameController) runTicker() {
	for {
		select {
		case <-c.eventTicker.C:
			c.ticker()
		}
	}
}

func (c *GameController) ticker() {
	switch c.roundStatus {
	case Betting:
		log.LogMessage(
			"crash ticker",
			"betting",
			"info",
			logrus.Fields{
				"isBettingDurationPassed": c.isBettingDurationPassed(),
			},
		)
		// 1. Perform event broadcast.
		if err := c.emitRealTimeEvent(); err != nil {
			log.LogMessage(
				"crash ticker",
				"failed to broadcast betting event",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
			return
		}

		// 2. Check status duration and call `startPending` accordingly.
		if c.isBettingDurationPassed() {
			if err := c.startPending(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to start pending period",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			}
		}
	case Pending:
		log.LogMessage(
			"crash ticker",
			"pending",
			"info",
			logrus.Fields{
				"isPerformedCashInsEmpty": c.isPerformedCashInsEmpty(),
				"isPendingDurationPassed": c.isPendingDurationPassed(),
				"isCashInEventsEmpty":     c.isCashInEventsEmpty(),
			},
		)
		// 1. Perform event broadcast if `cashInEvents` has element.
		if !c.isPerformedCashInsEmpty() {
			if err := c.emitRealTimeEvent(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to broadcast performed cash-ins",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
				return
			}
		}
		// 2. If `pendingDuration` passed, and no events in `cashInEvents`, call `startRunning`.
		if c.isPendingDurationPassed() && c.isCashInEventsEmpty() {
			if err := c.startRunning(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to start running period",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			}
		}
	case Running:
		log.LogMessage(
			"crash ticker",
			"running",
			"info",
			logrus.Fields{
				"nextMultiplier": c.nextMultiplier,
			},
		)
		// 1. `currentMultiplier` <- `nextMultiplier`, ++currentStep, Calculate `nextMultiplier`.
		c.toNextRunningStep()

		// 2. If there's any cash-ins which `autoCashOutAt` is between
		// `currentMultiplier` and `nextMultiplier`, insert cash-out
		// events into `cashOutEvents`.
		c.getShouldBePerformedCashOuts()

		if c.nextMultiplier < c.round.Outcome {
			// 3. If `nextMultiplier` is less than round's Outcome, broadcast running event.
			if err := c.emitRealTimeEvent(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to broadcast running event.",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
				return
			}
		} else {
			// 4. Else call `startPreparing`.
			if err := c.startPreparing(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to start preparing period.",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			}
		}
	case Preparing:
		log.LogMessage(
			"crash ticker",
			"preparing",
			"info",
			logrus.Fields{
				"isPerformedCashOutsEmpty":     c.isPerformedCashOutsEmpty(),
				"isCashOutEventsEmpty":         c.isCashOutEventsEmpty(),
				"isPreparingRound":             c.isPreparingRound(),
				"isPreparingDurationPassed":    c.isPreparingDurationPassed(),
				"isRoundInitializedForBetting": c.isRoundInitializedForBetting(),
			},
		)
		// 1. Perform event broadcast if `cashOutEvents` has element.
		if !c.isPerformedCashOutsEmpty() {
			if err := c.emitRealTimeEvent(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to broadcast performed cash-outss",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
				return
			}
		}

		// 2. If no events in `cashOutEvents`, call `prepareBetting`.
		if c.isCashOutEventsEmpty() &&
			!c.isPreparingRound() {
			if err := c.prepareBetting(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to prepare betting",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			}
		}

		// 3. If `preparingDuration` is passed and no events in `cashOutEvents`, call `startBetting`.
		if c.isCashOutEventsEmpty() &&
			c.isPreparingDurationPassed() &&
			c.isRoundInitializedForBetting() {
			if err := c.startBetting(); err != nil {
				log.LogMessage(
					"crash ticker",
					"failed to start betting",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			}
		}
	default:
	}
}

/*
/* @Internal
/* Checks whether betting duration passed after status updated as
/* 'betting'.
*/
func (c *GameController) isBettingDurationPassed() bool {
	return c.roundStatus == Betting &&
		time.Now().After(c.lastStatusUpdated.Add(c.bettingDuration))
}

/*
/* @Internal
/* Checks whether pending duration passed after status updated as
/* 'pending'.
*/
func (c *GameController) isPendingDurationPassed() bool {
	return c.roundStatus == Pending &&
		time.Now().After(c.lastStatusUpdated.Add(c.pendingDuration))
}

/*
/* @Internal
/* Checks whether preparing duration passed after status updated as
/* 'preparing'.
*/
func (c *GameController) isPreparingDurationPassed() bool {
	return c.roundStatus == Preparing &&
		time.Now().After(c.lastStatusUpdated.Add(c.preparingDuration))
}

/*
/* @Internal
/* Checks whether `performedCashIns` is empty.
*/
func (c *GameController) isPerformedCashInsEmpty() bool {
	return len(c.performedCashIns) == 0
}

/*
/* @Internal
/* Checks whether `performedCashOuts` is empty.
*/
func (c *GameController) isPerformedCashOutsEmpty() bool {
	return len(c.performedCashOuts) == 0
}

/*
/* @Internal
/* Checks whether `cashInEvents` is empty.
*/
func (c *GameController) isCashInEventsEmpty() bool {
	return len(c.cashInEvents) == 0
}

/*
/* @Internal
/* Checks whether `cashOutEvents` is empty.
*/
func (c *GameController) isCashOutEventsEmpty() bool {
	return len(c.cashOutEvents) == 0
}

/*
/* @Internal
/* Calculate `nextMultiplier` and update status as next running step.
/* - `currentMultiplier` <- `nextMultiplier`.
/* - calculate `nextMultiplier` on next step.
/* - increase `currentStep` by 1.const
/* - set `nextMultiplier` as minimum value of it's value & round's
/*   `Outcome`.
*/
func (c *GameController) toNextRunningStep() {
	c.currentMultiplier = c.nextMultiplier
	c.nextMultiplier = c.nextMultiplier * c.multiplierIncreaseRate
	c.currentStep++

	c.nextMultiplier = math.Min(c.nextMultiplier, c.round.Outcome)
}

/*
/* @internal
/* Returns bets which should be cashed out in current step.
*/
func (c *GameController) getShouldBePerformedCashOuts() {
	for _, bet := range c.round.Bets {
		log.LogMessage(
			"getShouldBePerformedCashOuts",
			"checking",
			"info",
			logrus.Fields{
				"cashoutAt":         bet.CashOutAt,
				"profit":            bet.Profit,
				"payoutMultiplier":  bet.PayoutMultiplier,
				"shouldCashedOut":   c.isShouldBeCashedOut(bet),
				"currentMultiplier": c.currentMultiplier,
				"nextMultiplier":    c.nextMultiplier,
			},
		)
		if c.isShouldBeCashedOut(bet) {
			event := CashOutEvent{
				UserID:           bet.UserID,
				RoundID:          c.round.ID,
				BetID:            bet.ID,
				PayoutMultiplier: *bet.CashOutAt,
				Type:             AutoCashOut,
			}
			c.CashOut(event)
		}
	}
}

/*
/* @Internal
/* Is should be cashed out bet in current step.
*/
func (c *GameController) isShouldBeCashedOut(bet models.CrashBet) bool {
	return bet.CashOutAt != nil &&
		bet.Profit == nil &&
		bet.PayoutMultiplier == nil &&
		// *bet.CashOutAt > c.currentMultiplier &&
		*bet.CashOutAt <= c.nextMultiplier
}

/*
* @Internal
* Checks whether round is okay to start betting.
 */
func (c *GameController) isRoundInitializedForBetting() bool {
	return c.isPreparingRound() &&
		c.round.Outcome > 0
}
