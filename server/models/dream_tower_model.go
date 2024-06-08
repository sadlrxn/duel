package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DreamTowerStatus string

const (
	DreamTowerPlaying DreamTowerStatus = "playing"
	DreamTowerWin     DreamTowerStatus = "win"
	DreamTowerLoss    DreamTowerStatus = "loss"
	DreamTowerCashout DreamTowerStatus = "cashout"
)

type DifficultyLevel string

const (
	LevelEasy   DifficultyLevel = "Easy"
	LevelMedium DifficultyLevel = "Medium"
	LevelHard   DifficultyLevel = "Hard"
	LevelExpert DifficultyLevel = "Expert"
	LevelMaster DifficultyLevel = "Master"
)

type DreamTowerDifficulty struct {
	Level       DifficultyLevel `json:"level"`
	BlocksInRow uint            `json:"blocksInRow"`
	StarsInRow  uint            `json:"starsInRow"`
}

type DreamTowerRound struct {
	gorm.Model
	UserID          uint                 `gorm:"not null;index:user_id" json:"userId"`
	BetAmount       int64                `gorm:"not null;index" json:"betAmount"`
	SeedPairID      uint                 `gorm:"not null" json:"seedPairId"`
	Nonce           uint                 `gorm:"not null" json:"nonce"`
	Bets            pq.Int32Array        `gorm:"type:integer[]" json:"bets"`
	Status          DreamTowerStatus     `gorm:"not null;index:status" json:"status"`
	Difficulty      DreamTowerDifficulty `gorm:"not null;embedded" json:"difficulty"`
	Profit          *int64               `json:"profit"`
	PaidBalanceType PaidBalanceForGame   `gorm:"not null;default:chip" json:"paidBalanceType"`

	RefTransactions []Transaction `gorm:"polymorphic:Owner;polymorphicValue:tx_dream_tower_referenced" json:"refTransactions"`
}
