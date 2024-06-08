package models

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID    uint          `gorm:"not null" json:"userId"`
	Balance   Balance       `gorm:"not null;polymorphic:Owner;polymorphicValue:in-wallet" json:"balance"`
	History   []Balance     `gorm:"polymorphic:Owner;polymorphicValue:in-history" json:"history"`
	IncomeTx  []Transaction `gorm:"foreignKey:ToWallet" json:"incomeTx"`
	OutcomeTx []Transaction `gorm:"foreignKey:FromWallet" json:"outcomeTx"`

	RefTransactions []Transaction `gorm:"polymorphic:Owner;polymorphicValue:tx_wallet_referenced" json:"refTransactions"`
}
