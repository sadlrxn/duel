package models

import "gorm.io/gorm"

type Role string

const (
	AdminRole      Role = "admin"
	ModeratorRole  Role = "moderator"
	AmbassadorRole Role = "ambassador"
	UserRole       Role = "user"
)

type DepositedNft struct {
	gorm.Model
	Name         string  `gorm:"not null;varchar(50)" json:"name"`
	CollectionID uint    `gorm:"not null" json:"collectionId"`
	MintAddress  string  `gorm:"not null;varchar(50);index" json:"mintAddress"`
	Image        string  `gorm:"not null" json:"image"`
	WalletID     *uint   `gorm:"index" json:"walletId"`
	Wallet       *Wallet `gorm:"foreignKey:WalletID" json:"wallet"`
}

type GameStats struct {
	TotalRounds  uint  `gorm:"not null;default:0" json:"totalRounds"`
	WinnedRounds uint  `gorm:"not null;default:0" json:"winnedRounds"`
	LostRounds   uint  `gorm:"not null;default:0" json:"lostRounds"`
	Wagered      int64 `gorm:"not null;default:0" json:"wagered"`
	Profit       int64 `gorm:"not null;default:0" json:"profit"`
	Loss         int64 `gorm:"not null;default:0" json:"loss"`
}

type Statistics struct {
	gorm.Model
	UserID          uint      `gorm:"not null;index" json:"userId"`
	JackpotStats    GameStats `gorm:"not null;embedded;embeddedPrefix:jackpot_" json:"jackpotStats"`
	CoinflipStats   GameStats `gorm:"not null;embedded;embeddedPrefix:coinflip_" json:"coinflipStats"`
	DreamtowerStats GameStats `gorm:"not null;embedded;embeddedPrefix:dreamtower_" json:"dreamtowerStats"`
	CrashStats      GameStats `gorm:"not null;embedded;embeddedPrefix:crash_" json:"crashStats"`
	WinStreaks      uint      `gorm:"not null;default:0" json:"winStreaks"`
	LoseStreaks     uint      `gorm:"not null;default:0" json:"loseStreaks"`
	BestStreaks     uint      `gorm:"not null;default:0" json:"bestStreaks"`
	WorstStreaks    uint      `gorm:"not null;default:0" json:"worstStreaks"`
	TotalWagered    int64     `gorm:"not null;default:0" json:"totalWagered"`
	TotalWin        int64     `gorm:"not null;default:0" json:"totalWin"`
	TotalLoss       int64     `gorm:"not null;default:0" json:"totalLoss"`
	MaxProfit       int64     `gorm:"not null;default:0" json:"maxProfit"`
	TotalProfit     int64     `gorm:"not null;default:0" json:"totalProfit"`
}

type User struct {
	gorm.Model
	Name           string     `gorm:"type:varchar(16);not null;uniqueIndex" json:"name"`
	WalletAddress  string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"walletAddress"`
	Role           Role       `gorm:"not null;default:user" json:"role"`
	Nonce          string     `gorm:"type:varchar(64)" json:"nonce"`
	Avatar         string     `json:"avatar"`
	Wallet         Wallet     `gorm:"not null" json:"wallet"`
	Statistics     Statistics `gorm:"foreignKey:UserID" json:"statistics"`
	PrivateProfile bool       `gorm:"not null;default:false" json:"privateProfile"`
	Banned         bool       `gorm:"not null;default:false" json:"banned"`
	IpAddress      string     `json:"ipAddress"`
}
