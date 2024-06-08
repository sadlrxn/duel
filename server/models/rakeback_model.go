package models

import (
	"time"

	"gorm.io/gorm"
)

type Rakeback struct {
	gorm.Model
	UserID                    uint      `gorm:"index:user_id;not null" json:"userId"`
	User                      User      `gorm:"foreignKey:UserID" json:"user"`
	TotalEarned               int64     `json:"totalEarned"`
	Reward                    int64     `json:"reward"`
	AdditionalRakebackRate    uint      `json:"additionalRakeBackRate"`
	AdditionalRakebackExpired time.Time `json:"additionalRakeBackExpired"`
	ActivatedAffiliateOnce    bool      `json:"activatedAffiliateOnce"`
}
