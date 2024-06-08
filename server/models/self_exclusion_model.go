package models

import "time"

type SelfExclusion struct {
	UserID uint      `gorm:"primarykey;autoIncrement:false" json:"userId"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
	Until  time.Time `json:"until"`
}
