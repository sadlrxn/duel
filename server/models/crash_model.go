package models

import (
	"time"

	"gorm.io/gorm"
)

type CrashRound struct {
	gorm.Model
	Seed           string       `gorm:"not null;unique" json:"seed"`
	Outcome        float64      `json:"outcome"`
	BetStartedAt   *time.Time   `json:"betStartedAt"`
	RunStartedAt   *time.Time   `json:"runStartedAt"`
	EndedAt        *time.Time   `json:"endedAt"`
	Bets           []CrashBet   `gorm:"foreignkey:RoundID" json:"bets"`
	FeeTransaction *Transaction `gorm:"polymorphic:Owner;polymorphicValue:tx_crash_round_referenced_for_fee" json:"feeTransaction"`
}

type CrashBet struct {
	gorm.Model
	UserID             uint               `gorm:"not null;index:user_round_bet" json:"userId"`
	User               User               `json:"user"`
	RoundID            uint               `gorm:"not null;index:user_round_bet" json:"crashRoundId"`
	Round              CrashRound         `json:"round"`
	CashInTransaction  Transaction        `gorm:"polymorphic:Owner;polymorphicValue:tx_crash_bet_referenced_for_cash_in" json:"cashInTransaction"`
	CashOutTransaction *Transaction       `gorm:"polymorphic:Owner;polymorphicValue:tx_crash_bet_referenced_for_cash_out" json:"cashOutTransaction"`
	BetAmount          int64              `json:"betAmount"`
	Profit             *int64             `json:"profit"`
	PayoutMultiplier   *float64           `json:"payoutMultiplier"`
	CashOutAt          *float64           `json:"cashOutAt"`
	PaidBalanceType    PaidBalanceForGame `gorm:"not null;default:chip" json:"paidBalanceType"`
}
