package crash

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func (c *GameController) Pause() {
	c.isBlockCrash = true
}

func (c *GameController) Start() error {
	if c.round != nil {
		return utils.MakeError(
			"crash_controller_maintain",
			"Start",
			"currently playing round",
			errors.New("round pointer is not nil"),
		)
	}
	c.isBlockCrash = false
	return c.Init(GameControllerInitParams{
		EventIntervalMilli:     config.CRASH_EVENT_INTERVAL_MILLI,
		BettingDurationMilli:   config.CRASH_BETTING_DURATION_MILLI,
		PendingDurationMilli:   config.CRASH_PENDING_DURATION_MILLI,
		PreparingDurationMilli: config.CRASH_PREPARING_DURATION_MILLI,
		BetCountLimit:          config.CRASH_BET_COUNT_LIMIT,
		MinBetAmount:           config.CRASH_MIN_BET_AMOUNT,
		MaxBetAmount:           config.CRASH_MAX_BET_AMOUNT,
		MultiplierIncreaseRate: config.CRASH_MULTIPLIER_INCREASE_RATE,
		HouseEdge:              config.CRASH_HOUSE_EDGE,
		MaxPlayerLimit:         config.CRASH_MAX_PLAYER_LIMIT,
		MinCashOutAt:           config.CRASH_MIN_CASH_OUT_AT,
		TempUserID:             config.CRASH_TEMP_ID,
		FeeUserID:              config.CRASH_FEE_ID,
		MaxCashOut:             config.CRASH_MAX_CASH_OUT,
	})
}
