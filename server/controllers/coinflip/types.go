package coinflip

import (
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"golang.org/x/sync/syncmap"
)

type Opponent string

const (
	Dueler Opponent = "dueler"
	Bot    Opponent = "bot"
)

type EventParam struct {
	EventType       string                    `json:"eventType"`
	RoundID         *uint                     `json:"roundId"`
	Side            string                    `json:"side"`
	Amount          int64                     `json:"amount"`
	Opponent        Opponent                  `json:"opponent"`
	PaidBalanceType models.PaidBalanceForGame `json:"paidBalanceType"`
}

type Controller struct {
	activeRounds   syncmap.Map
	round2Creator  syncmap.Map
	isRoundPending syncmap.Map
	roundLimit     uint
	minAmount      int64
	maxAmount      int64
	fee            int64
	EventEmitter   chan types.WSEvent
}
