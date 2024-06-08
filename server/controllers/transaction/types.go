package transaction

import (
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
)

type TransactionRequest struct {
	FromUser          *db_aggregator.User
	ToUser            *db_aggregator.User
	Balance           db_aggregator.BalanceLoad
	Type              models.TransactionType
	ToBeConfirmed     bool
	HouseFeeMeta      *HouseFeeMeta
	OwnerType         models.TransactionOwnerType
	OwnerID           uint
	BatchHouseFeeMeta []HouseFeeMeta
	EdgeDetails       EdgeDetailsInTransactionRequest
}

type EdgeDetailsInTransactionRequest struct {
	ShouldNotDistributeRevShare bool
}

type RainRequest struct {
	FromUser *db_aggregator.User
	ToUsers  *[]db_aggregator.User
	Balance  db_aggregator.BalanceLoad
	Type     models.TransactionType
}

type HouseFeeMeta struct {
	User        db_aggregator.User
	WagerAmount int64
	FeeAmount   int64
}

type DuelBotsRequest struct {
	FromUser db_aggregator.User
	DuelBots []db_aggregator.Nft
}

type ConfirmRequest struct {
	Transaction db_aggregator.Transaction
	OwnerType   models.TransactionOwnerType
	OwnerID     uint
}

type DeclineRequest struct {
	Transaction db_aggregator.Transaction
	OwnerType   models.TransactionOwnerType
	OwnerID     uint
}
