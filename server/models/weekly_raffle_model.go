package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type WeeklyRaffle struct {
	StartedAt datatypes.Date `gorm:"primaryKey;type:date" json:"startedAt"`
	EndAt     time.Time      `gorm:"index" json:"endAt"`
	Prizes    pq.Int64Array  `gorm:"type:bigint[]" json:"prizes"`
	Ended     bool           `gorm:"index" json:"ended"`
}

type WeeklyRaffleTicket struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	RoundStartedAt   datatypes.Date `gorm:"type:date;uniqueIndex:raffle_date_ticket;uniqueIndex:raffle_date_rank;index:raffle_date_user_claimed" json:"startedAt"`
	Round            WeeklyRaffle   `gorm:"foreignKey:RoundStartedAt" json:"round"`
	TicketID         uint           `gorm:"uniqueIndex:raffle_date_ticket" json:"ticketId"`
	UserID           uint           `gorm:"index:raffle_date_user_claimed" json:"userId"`
	User             User           `json:"user"`
	Rank             *uint          `gorm:"uniqueIndex:raffle_date_rank" json:"rank"`
	Claimed          *int64         `gorm:"index:raffle_date_user_claimed" json:"claimed"`
	ClaimTransaction *Transaction   `gorm:"polymorphic:Owner;polymorphicValue:tx_weekly_raffle_reward_referenced" json:"claimTransaction"`
}
