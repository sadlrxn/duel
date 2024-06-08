package models

import (
	"gorm.io/gorm"
)

type GameType string

const (
	Jackpot    GameType = "jackpot"
	Coinflip   GameType = "coinflip"
	Dreamtower GameType = "dreamtower"
	Crash      GameType = "crash"
)

type Game struct {
	gorm.Model
	Type        GameType `gorm:"not null;default:coinflip" json:"type"`
	IsActive    bool     `gorm:"not null;default:true" json:"isActive"`
	PlayerLimit uint8    `gorm:"not null;default:50" json:"playerLimit"`
	GameDetails string   `gorm:"type:text" json:"gameDetails"`
	OwnerID     uint     `json:"ownerId"`
	Description *string  `json:"description"`
	Avatar      *string  `json:"avatar"`
}
