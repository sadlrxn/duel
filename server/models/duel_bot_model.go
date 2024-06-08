package models

import "gorm.io/gorm"

type DuelBotStatus string

const (
	DuelBotNormal DuelBotStatus = "normal"
	DuelBotStaked DuelBotStatus = "staked"
)

type DuelBot struct {
	gorm.Model

	DepositedNftID uint          `gorm:"index:deposited_nft_id" json:"depositedNftId"`
	DepositedNft   DepositedNft  `gorm:"foreignKey:DepositedNftID" json:"depositedNft"`
	Status         DuelBotStatus `gorm:"index:status;default:normal" json:"status"`
	TotalEarned    int64         `json:"totalEarned"`
	StakingReward  int64         `json:"stakingReward"`
	StakingUserID  *uint         `gorm:"index:staking_user_id" json:"stakingUserId"`
	StakingUser    *User         `gorm:"foreignKey:StakingUserID" json:"stakingUser"`
}
