package jackpot

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"golang.org/x/sync/syncmap"
)

type Status string

const (
	Available Status = "available"
	Created   Status = "created"
	Started   Status = "started"
	Rolling   Status = "rolling"
)

type BetData struct {
	Amount    int64              `json:"amount"`
	NftAmount int64              `json:"nftAmount"`
	Nfts      []types.NftDetails `json:"nfts"`
	Time      time.Time          `json:"time"`
}

type Bet4Player struct {
	UserID         uint
	UserRole       models.Role
	UserName       string
	UserAvatar     string
	TotalUsdAmount int64
	TotalNftAmount int64
	Bets           []BetData
}

type Controller struct {
	Room types.Room
	status          Status
	minBetAmount    int64
	maxBetAmount    int64
	betCountLimit   uint
	playerLimit     uint
	countingTime    uint
	rollingTime     uint
	fee             int64
	Type            models.JackpotType
	roundID         uint
	playerToBets    syncmap.Map
	betPlayers      []uint
	nftsInGame      []models.NftInGame
	totalPlayers    uint
	timer           *time.Timer
	winnerID        *uint
	ticketID        *string
	signedString    *string
	lastUpdated     time.Time
	EventEmitter    chan types.WSEvent
	usdProfit       int64
	nfts4Profit     []types.NftDetails
	usdFee          int64
	nfts4Fee        []types.NftDetails
	totalAmount     int64
	totalFee        int64
	rollingDuration uint64
	candidates      []types.User
	lockUser syncmap.Map
}
