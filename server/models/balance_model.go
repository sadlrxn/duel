package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BalanceOwnerType string

const (
	InWallet      BalanceOwnerType = "in-wallet"
	InHistory     BalanceOwnerType = "in-history"
	InTransaction BalanceOwnerType = "in-transaction"
)

type ChipBalance struct {
	gorm.Model
	Balance int64 `gorm:"not null;default:0" json:"balance"`
}

type NftBalance struct {
	gorm.Model
	Balance pq.StringArray `gorm:"type:text[]" json:"balance"`
}

type RakeBack struct { // To Be Removed: Rakeback Migration
	gorm.Model
	Balance int64 `gorm:"not null;default:0" json:"balance"`
}

type Balance struct {
	gorm.Model
	ChipBalanceID *uint            `json:"chipBalanceId"`
	ChipBalance   *ChipBalance     `gorm:"foreignKey:ChipBalanceID" json:"chipBalance"`
	NftBalanceID  *uint            `json:"nftBalanceId"`
	NftBalance    *NftBalance      `gorm:"foreignKey:NftBalanceID" json:"nftBalance"`
	OwnerID       uint             `gorm:"not null;index:owner_key" json:"ownerId"`
	OwnerType     BalanceOwnerType `gorm:"not null;index:owner_key" json:"ownerType"`
}
