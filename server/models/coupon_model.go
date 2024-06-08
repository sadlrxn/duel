package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CouponType string

const (
	CouponForSpecUsers  CouponType = "coupon-for-spec-users"
	CouponForLimitUsers CouponType = "coupon-for-limit-users"
)

type Coupon struct {
	Code                  *uuid.UUID      `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"code"`
	CreatedAt             time.Time       `gorm:"index"`
	Type                  CouponType      `gorm:"not null" json:"type"`
	ClaimedCoupons        []ClaimedCoupon `gorm:"foreignKey:CouponID" json:"claimedCoupons"`
	AccessUserIDs         pq.Int64Array   `gorm:"type:bigint[]" json:"accessUserIds"`
	AccessUserLimit       uint            `json:"claimLimit"`
	BonusBalance          int64           `gorm:"not null" json:"bonusBalance"`
	RequiredAffiliateCode *string         `json:"requiredAffiliateCode"`
	RequiredAffiliate     *Affiliate      `gorm:"foreignKey:RequiredAffiliateCode;references:Code" json:"requiredAffiliate"`
}

type ClaimedCoupon struct {
	CreatedAt     time.Time `gorm:"index"`
	CouponID      uuid.UUID `gorm:"not null;primaryKey;autoIncrement:false" json:"couponId"`
	Coupon        Coupon    `gorm:"foreignKey:CouponID" json:"coupon"`
	ClaimedUserID uint      `gorm:"not null;index;primaryKey;autoIncrement:false" json:"claimedUserId"`
	ClaimedUser   User      `gorm:"foreignKey:ClaimedUserID" json:"claimedUser"`
	Wagered       int64     `gorm:"default:0" json:"wagered"`
	Balance       int64     `gorm:"not null" json:"balance"`
	Exchanged     int64     `gorm:"default:0;index" json:"exchanged"`
}

type CouponShortcut struct {
	CouponID uuid.UUID `gorm:"not null;primaryKey;autoIncrement:false" json:"couponId"`
	Coupon   Coupon    `gorm:"foreignKey:CouponID" json:"coupon"`
	Shortcut string    `gorm:"unique;index" json:"shortcut"`
}
