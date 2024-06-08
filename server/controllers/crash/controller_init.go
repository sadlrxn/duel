package crash

import (
	"sync"
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
)

type GameController struct {
	// Time duration between real time event emition
	eventInterval time.Duration
	// Event ticker which emits every `eventInterval` time.
	eventTicker *time.Ticker
	// Time duration to stay at `Betting` game status.
	bettingDuration time.Duration
	// Time duration to stay at `Pending` game status.
	pendingDuration time.Duration
	// Time duration to stay at `Preparing` game status.
	preparingDuration time.Duration
	// Bet count limit per user per round.
	betCountLimit uint
	// Minimum required bet amount per user per round.
	minBetAmount int64
	// Maximum limit bet amount per user per round.
	maxBetAmount int64
	// Multiplier increase rate per step.
	multiplierIncreaseRate float64
	// House edge when 10000 is 100 percentage.
	houseEdge int64
	// Crash temp user id.
	tempUserID uint
	// Crash fee user id.
	feeUserID uint
	// Maximum player limit per round.
	maxPlayerLimit uint
	// Minimum `cashOutAt`.
	minCashOutAt float64
	// The whole round data.
	round *models.CrashRound
	// Event emitter.
	EventEmitter chan types.WSEvent
	// Event queue for `CashInEvent`.
	cashInEvents chan CashInEvent
	// Event queue for `CashOutEvent`.
	cashOutEvents chan CashOutEvent
	// Current game status.
	roundStatus GameStatus
	// Last status updated time.
	lastStatusUpdated time.Time
	// Current multiplier.
	currentMultiplier float64
	// Next multiplier.
	nextMultiplier float64
	// Current step since starting.
	currentStep int64
	// Performed cash outs slice.
	performedCashOuts []CashOutForRealTimeEvent
	// Performed cash ins slice.
	performedCashIns []CashInForRealTimeEvent
	// Mutex for cash slice thread safe.
	cashMut sync.Mutex
	// Block crash flag.
	isBlockCrash bool
	// Mutex to block data race of status managing functions.
	statusMut sync.Mutex
	// Map to save cashIn events made by user. Cannot exceed `betCountLimit`.
	cashInCountsPerUser sync.Map
	// Map to save whether cashOut is requested for specific bet.
	cashOutFlagPerBet sync.Map
	// Max winning chips.
	maxCashOut int64
}

/*
/* @External
/* Initializes and allocates game controller's attributes.
*/
func (c *GameController) Init(initParams GameControllerInitParams) error {
	// Initialize static fields with `initParams`.
	c.eventInterval = time.Millisecond * time.Duration(initParams.EventIntervalMilli)
	c.bettingDuration = time.Millisecond * time.Duration(initParams.BettingDurationMilli)
	c.pendingDuration = time.Millisecond * time.Duration(initParams.PendingDurationMilli)
	c.preparingDuration = time.Millisecond * time.Duration(initParams.PreparingDurationMilli)
	c.betCountLimit = initParams.BetCountLimit
	c.minBetAmount = initParams.MinBetAmount
	c.maxBetAmount = initParams.MaxBetAmount
	c.multiplierIncreaseRate = initParams.MultiplierIncreaseRate
	c.houseEdge = initParams.HouseEdge
	c.tempUserID = initParams.TempUserID
	c.feeUserID = initParams.FeeUserID
	c.maxPlayerLimit = initParams.MaxPlayerLimit
	c.minCashOutAt = initParams.MinCashOutAt
	c.maxCashOut = initParams.MaxCashOut

	// Initialize other feilds with default values.
	c.eventTicker = time.NewTicker(c.eventInterval)
	c.roundStatus = Preparing
	c.lastStatusUpdated = time.Now()
	c.isBlockCrash = false

	// Initialize first round.
	if err := c.loadRoundOnInit(); err != nil {
		return utils.MakeError(
			"crash_controller_init",
			"Init",
			"failed to initialize first round",
			err,
		)
	}
	if c.cashInEvents == nil || c.cashOutEvents == nil {
		c.cashInEvents = make(chan CashInEvent, 1024)
		c.cashOutEvents = make(chan CashOutEvent, 1024)
		go c.eventListener()
	}

	// Start running ticker function.
	go c.runTicker()

	return nil
}
