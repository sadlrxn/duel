package models

import (
	"time"

	"gorm.io/gorm"
)

type ServerConfig struct {
	gorm.Model
	ShouldNotReset                  bool      `gorm:"default:false" json:"shouldReset"`
	BaseRakeBackRate                uint      `gorm:"not null;default:5" json:"baseRakeBackRate"`
	AdditionalRakeBackRate          uint      `json:"additionalRakeBackRate"`
	NextGrandJackpotStartAt         time.Time `json:"nextGrandJackpotStartAt"`
	ApiRateLimitConfiguration       string    `json:"apiRateLimitConfiguration"`
	WebsocketRateLimitConfiguration string    `json:"websocketRateLimitConfiguration"`
	CrashClientSeed                 string    `json:"crashClientSeed"`
}
