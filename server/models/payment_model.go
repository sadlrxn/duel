package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	Success PaymentStatus = "success"
	Failed  PaymentStatus = "failed"
	Pending PaymentStatus = "pending"
)

type SolDetail struct {
	SolAmount int64 `json:"solAmount"`
	UsdAmount int64 `json:"usdAmount"`
}

type NftDetail struct {
	Mints         pq.StringArray `gorm:"type:text[]" json:"mints"`
	MintAddresses string         `gorm:"type:text" json:"mintAddresses"`
}

type Payment struct {
	gorm.Model
	UserID        uint          `gorm:"not null" json:"userId"`
	Type          string        `gorm:"not null;default:deposit_sol;index:type" json:"type"`
	Status        PaymentStatus `gorm:"not null;default:pending;index:status" json:"status"`
	SolDetail     SolDetail     `gorm:"embedded"`
	NftDetail     NftDetail     `gorm:"embedded"`
	TxHash        string        `gorm:"type:varchar(100)" json:"txHash"`
	TransactionID *uint         `json:"transactionId"`
	Transaction   *Transaction  `gorm:"foreignKey:TransactionID" json:"transaction"`
}
