package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CouponTransactionType string

const (
	CpTxClaimCode        CouponTransactionType = "cp-tx-claim-code"
	CpTxCoinflipBet      CouponTransactionType = "cp-tx-coinflip-bet"
	CpTxCoinflipProfit   CouponTransactionType = "cp-tx-coinflip-profit"
	CpTxDreamtowerBet    CouponTransactionType = "cp-tx-dreamtower-bet"
	CpTxDreamtowerProfit CouponTransactionType = "cp-tx-dreamtower-profit"
	CpTxCrashBet         CouponTransactionType = "cp-tx-crash-bet"
	CpTxCrashProfit      CouponTransactionType = "cp-tx-crash-profit"
	CpTxExchangeToChip   CouponTransactionType = "cp-tx-exchange-to-chip"
)

type CouponTransactionStatus string

const (
	CouponTransactionSucceed CouponTransactionStatus = "succeed"
	CouponTransactionFailed  CouponTransactionStatus = "failed"
	CouponTransactionPending CouponTransactionStatus = "pending"
)

type CouponTransaction struct {
	gorm.Model
	CouponID        uuid.UUID               `gorm:"not null" json:"couponId"`
	Coupon          Coupon                  `gorm:"foreignKey:CouponID" json:"coupon"`
	Type            CouponTransactionType   `json:"type"`
	ClaimedUserID   uint                    `gorm:"not null" json:"claimedUserId"`
	ClaimedUser     User                    `gorm:"foreignKey:ClaimedUserID" json:"claimedUser"`
	ActiveCoupon    ClaimedCoupon           `gorm:"foreignKey:CouponID,ClaimedUserID" json:"activeCoupon"`
	PrevBalance     int64                   `gorm:"default:0" json:"prevBalance"`
	TxBalance       int64                   `gorm:"default:0" json:"TxBalance"`
	NextBalance     int64                   `gorm:"default:0" json:"nextBalance"`
	AfterRefund     int64                   `json:"afterRefund"`
	Status          CouponTransactionStatus `json:"status"`
	RealTransaction *Transaction            `gorm:"polymorphic:Owner;polymorphicValue:tx_coupon_transaction_referenced" json:"realTransaction"`
}
