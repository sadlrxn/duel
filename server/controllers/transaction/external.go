package transaction

import (
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"

	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) error {
	return initialize(db)
}

func Transfer(transactionRequest *TransactionRequest) (*db_aggregator.Transaction, error) {
	return transfer(transactionRequest)
}

func Confirm(confirmRequest ConfirmRequest) error {
	return confirm(confirmRequest)
}

func Decline(declineRequest DeclineRequest) error {
	return decline(declineRequest)
}

func ClaimRakeback(user *db_aggregator.User) (int64, error) {
	return claimRakeback(user)
}

func StakeDuelBots(request DuelBotsRequest) error {
	return stakeDuelBots(request)
}

func UnstakeDuelBots(request DuelBotsRequest) (int64, error) {
	return unstakeDuelBots(request)
}

func ClaimDuelBotsRewards(request DuelBotsRequest) (int64, error) {
	return claimDuelBotsRewards(request)
}

func GetRakebackRewards(user *db_aggregator.User) (int64, int64, error) {
	return getRakebackRewards(user)
}

func CreateAffiliateCode(user db_aggregator.User, codes []string) error {
	return createAffiliateCode(user, codes)
}

func DeleteAffiliateCode(user db_aggregator.User, codes []string) (int64, error) {
	return deleteAffiliateCode(user, codes)
}

func ClaimAffiliateRewards(user db_aggregator.User, codes []string) (int64, error) {
	return claimAffiliateRewards(user, codes)
}

func ActivateAffiliateCode(user db_aggregator.User, code string) (bool, error) {
	return activateAffiliateCode(user, code)
}

func DeactivateAffiliateCode(user db_aggregator.User, code string) error {
	return deactivateAffiliateCode(user, code)
}

func Rain(request *RainRequest) (*db_aggregator.Transaction, error) {
	return rain(request)
}
