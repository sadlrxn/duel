package crash

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
)

type GameStatus string

/*
/* `GameStatus` represents round status.
/* - `crash-status-betting`:
/*    In this status, users are able to make bets.
/* - `crash-status-pending`:
/*    In this status, users are not able to make bets, round is being preparing
/*    to start running on graph.
/* - `crash-status-running`:
/*    In this status, users are able to cash out, graph is running every 0.2s.
/* - `crash-status-preparing`:
/*    Game status becomes this case by crashing from `crash-status-running`.
/*    During this status, just shows the last round's multiplier to users and
/*    wait for the next round start - `crash-status-betting`.
*/
const (
	Betting   GameStatus = "crash-status-betting"
	Pending   GameStatus = "crash-status-pending"
	Running   GameStatus = "crash-status-running"
	Preparing GameStatus = "crash-status-preparing"
)

/*
/* Client -> Server event
/* `CashInEvent` is content type of event which is emitted from client to
/* server. This event is emitted during `Betting` game status.
*/
type CashInEvent struct {
	UserID      uint                      `json:"userId"`
	Amount      int64                     `json:"amount"`
	BalanceType models.PaidBalanceForGame `json:"balanceType"`
	RoundID     uint                      `json:"roundId"`
	CashOutAt   float64                   `json:"cashOutAt"`
}

/*
/* Client -> Server event
/* `CashOutEvent` is content type of event which is emitted from client to
/* server. This event is emitted during `Betting` game status.
*/
type CashOutEvent struct {
	UserID           uint `json:"userId"`
	RoundID          uint `json:"roundId"`
	BetID            uint `json:"betId"`
	PayoutMultiplier float64
	Type             CashOutType
}

type CashInForRealTimeEvent struct {
	User        types.User                `json:"user"`
	Amount      int64                     `json:"amount"`
	BalanceType models.PaidBalanceForGame `json:"balanceType"`
	BetID       uint                      `json:"betId"`
	CashOutAt   float64                   `json:"cashOutAt"`
}

type CashOutType string

const (
	AutoCashOut CashOutType = "crash-cashout-auto"
	StopCashOut CashOutType = "crash-cashout-stop"
)

type CashOutForRealTimeEvent struct {
	User        types.User                `json:"user"`
	Amount      int64                     `json:"amount"`
	BalanceType models.PaidBalanceForGame `json:"balanceType"`
	Multiplier  float64                   `json:"multiplier"`
	BetID       uint                      `json:"betId"`
}

/*
/* Server -> Client event
/* `RealTimeEvent` is content type of event which is emitted from server to
/* client. This event is emitted during `Betting` and `Running` game status.
/* - `Elapsed` is time passed after `Running` status.
/* - `Multiplier` is curreint multiplier for `Running` status.
/*   It represents crash point for `Preparing` status.
/* - `Timestamp` is server timestamp at the point of emitting event.
/* - `BetRemaining` is valid at status of `Betting` status.
/*   It represents time left until `Pending` status.
*/
type RealTimeEvent struct {
	RoundID     uint                      `json:"roundId"`
	CashIns     []CashInForRealTimeEvent  `json:"cashIns"`
	CashOuts    []CashOutForRealTimeEvent `json:"cashOuts"`
	Elapsed     int64                     `json:"elapsed"`
	Multiplier  float64                   `json:"multiplier"`
	RoundStatus GameStatus                `json:"roundStatus"`
}

type GameControllerInitParams struct {
	EventIntervalMilli     int64
	BettingDurationMilli   int64
	PendingDurationMilli   int64
	PreparingDurationMilli int64
	BetCountLimit          uint
	MinBetAmount           int64
	MaxBetAmount           int64
	MultiplierIncreaseRate float64
	HouseEdge              int64
	TempUserID             uint
	FeeUserID              uint
	MaxPlayerLimit         uint
	MinCashOutAt           float64
	MaxCashOut             int64
}

type CashInRequestParams CashInEvent

type CashOutRequestParams CashOutEvent

type CashOutRequestResult struct {
	Amount      int64
	BalanceType models.PaidBalanceForGame
}

type BetPayload struct {
	BetID            uint                      `json:"betId"`
	User             types.User                `json:"user"`
	BetAmount        int64                     `json:"betAmount"`
	PaidBalanceType  models.PaidBalanceForGame `json:"paidBalanceType"`
	Profit           *int64                    `json:"profit"`
	PayoutMultiplier *float64                  `json:"payoutMultiplier"`
	CashOutAt        *float64                  `json:"cashOutAt"`
}

type RoundPayload struct {
	RoundID      uint         `json:"roundId"`
	RoundStatus  GameStatus   `json:"roundStatus"`
	Elapsed      int64        `json:"elapsed"`
	Multiplier   float64      `json:"multiplier"`
	RunStartedAt *time.Time   `json:"runStartedAt"`
	Bets         []BetPayload `json:"bets"`
}

type RoundHistoryItem struct {
	ID      uint    `gorm:"column:id" json:"id"`
	Outcome float64 `gorm:"column:outcome" json:"outcome"`
}

type BetInRoundHistoryDetail struct {
	User             types.User                `json:"user"`
	BetAmount        int64                     `json:"betAmount"`
	Profit           int64                     `json:"profit"`
	PayoutMultiplier float64                   `json:"payoutMultiplier"`
	PaidBalanceType  models.PaidBalanceForGame `json:"paidBalanceType"`
}

type RoundHistoryDetail struct {
	ID      uint                      `json:"id"`
	Seed    string                    `json:"seed"`
	Outcome float64                   `json:"outcome"`
	Date    time.Time                 `json:"date"`
	Bets    []BetInRoundHistoryDetail `json:"bets"`
}
