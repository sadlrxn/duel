package models

import (
	"time"

	"gorm.io/gorm"
)

type JackpotType string

const (
	Low    JackpotType = "jackpotLow"
	Medium JackpotType = "jackpotMedium"
	Wild   JackpotType = "jackpotWild"
	Grand  JackpotType = "grand"
)

type NftGameStatus string

const (
	ChargedAsFee NftGameStatus = "charged_as_fee"
)

type NftInGame struct {
	gorm.Model
	Name            string        `gorm:"not null" json:"name"`
	MintAddress     string        `gorm:"not null;varchar(50)" json:"mintAddress"`
	Image           string        `gorm:"not null" json:"image"`
	CollectionName  string        `gorm:"not null" json:"collectionName"`
	CollectionImage string        `gorm:"not null" json:"collectionImage"`
	Price           int64         `gorm:"not null" json:"price"`
	Status          NftGameStatus `json:"status"`
	BetID           uint          `gorm:"not null" json:"betId"`
}

type JackpotBet struct {
	gorm.Model
	UsdAmount int64       `json:"usdAmount"`
	Nfts      []NftInGame `gorm:"foreignKey:BetID" json:"nfts"`
	PlayerID  uint        `gorm:"not null" json:"playerId"`
}

type JackpotPlayer struct {
	gorm.Model
	UserID  uint         `gorm:"not null;index" json:"userId"`
	Bets    []JackpotBet `gorm:"foreignKey:PlayerID" json:"bets"`
	RoundID uint         `gorm:"not null" json:"roundId"`
}

type JackpotRound struct {
	gorm.Model
	StartedAt         time.Time       `json:"startedAt"`
	CountingStartedAt time.Time       `json:"countingStartedAt"`
	EndedAt           time.Time       `json:"endedAt"`
	Players           []JackpotPlayer `gorm:"foreignKey:RoundID" json:"players"`
	WinnerID          uint            `json:"winnerId"`
	ChargedFee        int64           `gorm:"not null" json:"chargedFee"`
	TicketID          string          `gorm:"not null" json:"ticketId"`
	Type              JackpotType     `gorm:"not null;default:normal;index:type" json:"type"`
	SignedString      *string         `gorm:"index:signed_string" json:"signedString"`

	RefTransactions []Transaction `gorm:"polymorphic:Owner;polymorphicValue:tx_jackpot_referenced" json:"refTransactions"`
}
