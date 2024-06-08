package db_aggregator

import "github.com/Duelana-Team/duelana-v1/models"

type User uint
type Wallet uint
type Nft string
type Transaction uint
type Balance uint

type TransactionLoad struct {
	FromWallet *Wallet
	ToWallet   *Wallet
	Balance    BalanceLoad
	Type       models.TransactionType
	Status     models.TransactionStatus

	FromWalletPrevID *Balance
	FromWalletNextID *Balance
	RefundPrevID     *Balance
	RefundNextID     *Balance
	ToWalletPrevID   *Balance
	ToWalletNextID   *Balance

	OwnerID   uint
	OwnerType models.TransactionOwnerType

	Receipients *[]uint
}

type BalanceLoad struct {
	Balance     *Balance
	ChipBalance *int64
	NftBalance  *[]Nft
}

type ChangedBalance struct {
	ChipBalance bool
	NftBalance  bool
}

type TransferResult struct {
	FromWallet      *Wallet
	ToWallet        *Wallet
	FromPrevBalance *Balance
	FromNextBalance *Balance
	ToPrevBalance   *Balance
	ToNextBalance   *Balance
}

type RainResult struct {
	FromPrevBalance *Balance
	FromNextBalance *Balance
}

type BalanceHistoryChain struct {
	Prev *Balance
	Next *Balance
}

type QueryArgType string

const (
	Preload QueryArgType = "preload"
	Where   QueryArgType = "where"
	Update  QueryArgType = "update"
	Select  QueryArgType = "select"
	Lower   QueryArgType = "lower"
	Upper   QueryArgType = "upper"
)

type QueryArgConventionType string

const (
	LowerConvention QueryArgConventionType = "lower-convention"
	UpperConvention QueryArgConventionType = "upper-convention"
)

type AffiliateMeta struct {
	Code         string `json:"code"`
	UserCnt      uint   `json:"userCnt"`
	TotalEarned  int64  `json:"totalEarned"`
	Reward       int64  `json:"reward"`
	TotalWagered int64  `json:"totalWagered"`
	Rate         uint   `json:"rate"`
}

type UserInAffiliateDetail struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Lifetime uint   `json:"lifetime"`
	Wagered  int64  `json:"wagered"`
	Reward   int64  `json:"reward"`
}

type AffiliateDetail struct {
	Code  string                  `json:"code"`
	Users []UserInAffiliateDetail `json:"users"`
}

type ActiveAffiliateMeta struct {
	ID                  uint
	Code                string `json:"code"`
	Rate                uint   `json:"rate"`
	OwnerID             uint   `json:"ownerId"`
	OwnerName           string `json:"ownerName"`
	OwnerAvatar         string `json:"ownerAvatar"`
	IsFirstDepositBonus bool
	FirstDepositDone    bool
}
