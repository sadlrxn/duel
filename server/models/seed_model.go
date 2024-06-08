package models

import "gorm.io/gorm"

type ClientSeed struct {
	gorm.Model
	Seed     string `gorm:"not null" json:"seed"`
	SeedPair *SeedPair
}

type ServerSeed struct {
	gorm.Model
	Seed     string `gorm:"not null" json:"seed"`
	Hash     string `gorm:"not null" json:"hash"`
	SeedPair *SeedPair
}

type SeedPair struct {
	gorm.Model
	UserID           uint       `gorm:"not null;index:user_id" json:"userId"`
	User             User       `gorm:"not null" json:"user"`
	ClientSeedID     uint       `gorm:"not null" json:"clientSeedId"`
	ClientSeed       ClientSeed `gorm:"not null" json:"clientSeed"`
	ServerSeedID     uint       `gorm:"not null" json:"serverSeedId"`
	ServerSeed       ServerSeed `gorm:"not null" json:"serverSeed"`
	NextServerSeedID uint       `gorm:"not null" json:"nextServerSeedId"`
	NextServerSeed   ServerSeed `gorm:"not null" json:"nextServerSeed"`
	Nonce            uint       `gorm:"not null;default:0" json:"nonce"`
	UsingCount       uint       `json:"usingCount"`
	IsExpired        bool       `gorm:"not null;index:is_expired" json:"isExpired"`
}
