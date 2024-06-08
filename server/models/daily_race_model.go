package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DailyRaceRewards struct {
	gorm.Model
	StartedAt        datatypes.Date `gorm:"type:date;uniqueIndex:date_rank;uniqueIndex:date_user" json:"startedAt"`
	UserID           uint           `gorm:"uniqueIndex:date_user;index" json:"userId"`
	Rank             uint           `gorm:"uniqueIndex:date_rank" json:"rank"`
	Prize            int64          `json:"prize"`
	Claimed          int64          `gorm:"index" json:"claimed"`
	ClaimTransaction *Transaction   `gorm:"polymorphic:Owner;polymorphicValue:tx_daily_race_rewards_referenced" json:"claimTransaction"`
	Approved         bool           `gorm:"index" json:"approved"`
}
