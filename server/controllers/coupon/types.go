package coupon

import (
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/google/uuid"
)

type CreateCouponRequest struct {
	Type                  models.CouponType `json:"type"`
	AccessUserNames       *[]string         `json:"accessUserNames"`
	AccessUserIDs         *[]uint           `json:"accessUserIDs"`
	AccessUserLimit       *int              `json:"accessUserLimit"`
	Balance               int64             `json:"balance"`
	Shortcut              *string           `json:"shortcut"`
	RequiredAffiliateCode *string           `json:"requiredAffiliateCode"`
}

type CouponTransactionRequest struct {
	Type          models.CouponTransactionType `json:"type"`
	UserID        uint                         `json:"userId"`
	Balance       int64                        `json:"balance"`
	ToBeConfirmed bool                         `json:"toBeConfirmed"`
}

type ActiveUserCouponMeta struct {
	Code          uuid.UUID `json:"code"`
	Balance       int64     `json:"balance"`
	Claimed       int64     `json:"claimed"`
	Wagered       int64     `json:"wagered"`
	WagerLimit    int64     `json:"wagerLimit"`
	RemainingTime int64     `json:"remainingTime"`
}

type TryBetWithCouponResult string

const (
	CouponBetUnavailable       TryBetWithCouponResult = "coupon-bet-unavailable"
	CouponBetSucceed           TryBetWithCouponResult = "coupon-bet-succeed"
	CouponBetInsufficientFunds TryBetWithCouponResult = "coupon-bet-insufficient-funds"
	CouponBetFailed            TryBetWithCouponResult = "coupon-bet-failed"
)

type TryBetWithCouponRequest struct {
	UserID  uint                         `json:"userId"`
	Balance int64                        `json:"balance"`
	Type    models.CouponTransactionType `json:"type"`
}
