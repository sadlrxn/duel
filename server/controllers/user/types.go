package user

import (
	"github.com/Duelana-Team/duelana-v1/controllers/chat"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/google/uuid"
)

type Controller struct {
	EventEmitter chan types.WSEvent
	Chat         *chat.Controller
}

type ServerStatisticsResult struct {
	TotalBets    int64 `json:"totalBets"`
	TotalWagered int64 `json:"totalWagered"`
	TotalProfit  int64 `json:"totalProfit"`
}

type CreateUserRequest models.User
type SaveUserRequest models.User

type UserLoadResponse struct {
	ID            uint         `json:"id"`
	Name          string       `json:"name"`
	WalletAddress string       `json:"walletAddress"`
	Role          models.Role  `json:"role"`
	Avatar        string       `json:"avatar"`
	Balances      userBalances `json:"balances"`
	Nfts          userNfts     `json:"nfts"`
}

type userBalances struct {
	Chip   userBalance `json:"chip"`
	Coupon userBalance `json:"coupon,omitempty"`
}

type userBalance struct {
	Code          uuid.UUID `json:"code,omitempty"`
	Balance       int64     `json:"balance"`
	Claimed       int64     `json:"claimed,omitempty"`
	Wagered       int64     `json:"wagered,omitempty"`
	WagerLimit    int64     `json:"wagerLimit,omitempty"`
	RemainingTime int64     `json:"remainingTime,omitempty"`
}

type userNfts struct {
	Deposited []types.NftDetails `json:"deposited"`
}

type userInfoRequest struct {
	UserID   *uint   `form:"userId"`
	UserName *string `form:"userName"`
}

type userInfoResponse struct {
	Info       types.User      `json:"info"`
	Statistics *userStatistics `json:"statistics,omitempty"`
}

type userStatistics struct {
	TotalRounds    uint  `json:"total_rounds"`
	WinnedRounds   uint  `json:"winned_rounds"`
	LostRounds     uint  `json:"lost_rounds"`
	BestStreaks    uint  `json:"best_streaks"`
	WorstStreaks   uint  `json:"worst_streaks"`
	TotalWagered   int64 `json:"total_wagered"`
	PrivateProfile bool  `json:"private_profile"`
}
