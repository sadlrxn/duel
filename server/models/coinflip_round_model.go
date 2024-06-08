package models

import (
	"time"

	"gorm.io/gorm"
)

type PaidBalanceForGame string

const (
	ChipBalanceForGame   PaidBalanceForGame = "chip"
	CouponBalanceForGame PaidBalanceForGame = "coupon"
)

type CoinflipSide string

const (
	Heads CoinflipSide = "heads"
	Tails CoinflipSide = "tails"
)

type CoinflipRound struct {
	gorm.Model
	TailsUserID     *uint              `gorm:"index" json:"tailsUserId"`
	HeadsUserID     *uint              `gorm:"index" json:"headsUserId"`
	Amount          int64              `gorm:"not null;default:0" json:"amount"`
	EndedAt         time.Time          `gorm:"index" json:"endedAt"`
	WinnerID        *uint              `json:"winnerId"`
	Prize           int64              `json:"prize"`
	TicketID        string             `gorm:"not null" json:"ticketId"`
	SignedString    *string            `gorm:"index" json:"signedString"`
	PaidBalanceType PaidBalanceForGame `gorm:"not null;default:chip" json:"paidBalanceType"`

	RefTransactions []Transaction `gorm:"polymorphic:Owner;polymorphicValue:tx_coinflip_referenced" json:"refTransactions"`
}
