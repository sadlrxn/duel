package models

import (
	"time"

	"gorm.io/gorm"
)

type ActiveAffiliate struct {
	gorm.Model
	AffiliateID      uint      `json:"affiliateId"`
	Affiliate        Affiliate `gorm:"foreignKey:AffiliateID" json:"affiliate"`
	UserID           uint      `gorm:"unique" json:"userId"`
	User             User      `gorm:"foreignKey:UserID" json:"user"`
	FirstDepositDone bool      `json:"firstDepositDone"`
}

type Affiliate struct {
	gorm.Model
	Code                string              `gorm:"type:varchar(50);not null;index;unique" json:"code"`
	CreatorID           uint                `gorm:"index:creator_id" json:"creatorId"`
	Creator             User                `gorm:"foreignKey:CreatorID" json:"creator"`
	ActiveAffiliates    []ActiveAffiliate   `json:"activeAffiliates"`
	AffiliateLifetimes  []AffiliateLifetime `json:"affiliateLifetimes"`
	TotalWagered        int64               `json:"totalWagered"`
	TotalEarned         int64               `json:"totalEarned"`
	Reward              int64               `json:"reward"`
	CustomAffiliateRate uint                `json:"customAffiliateRate"`
	IsFirstDepositBonus bool                `json:"isFirstDepositBonus"`
}

type AffiliateLifetime struct {
	UserID          uint       `gorm:"primarykey;autoIncrement:false;" json:"userId"`
	AffiliateID     uint       `gorm:"primarykey;autoIncrement:false;" json:"affiliateId"`
	User            User       `gorm:"foreignkey:UserID" json:"user"`
	Affiliate       Affiliate  `gorm:"foreignkey:AffiliateID" json:"affiliate"`
	Lifetime        uint       `json:"lifetime"`
	LastActivated   time.Time  `json:"lastActivated"`
	LastDeactivated *time.Time `json:"lastDeactivated"`
	TotalWagered    int64      `json:"totalWagered"`
	TotalReward     int64      `json:"totalReward"`
	IsActive        bool       `json:"isActive"`
}
