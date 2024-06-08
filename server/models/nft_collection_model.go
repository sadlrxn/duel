package models

import (
	"gorm.io/gorm"
)

type NftCollection struct {
	gorm.Model
	Name       string         `gorm:"not null;varchar(50)" json:"name"`
	MoonRank   string         `gorm:"varchar(50)" json:"moonRank"`
	Solanart   string         `gorm:"varchar(50)" json:"solanart"`
	MagicEden  string         `gorm:"varchar(50)" json:"magicEden"`
	HyperSpace string         `gorm:"varchar(50)" json:"hyperSpace"`
	HowRare    string         `gorm:"varchar(50)" json:"howRare"`
	Image      string         `gorm:"not null" json:"image"`
	FloorPrice int64          `gorm:"not null;default:0" json:"floorPrice"`
	Nfts       []DepositedNft `gorm:"foreignKey:CollectionID" json:"nfts"`
}
