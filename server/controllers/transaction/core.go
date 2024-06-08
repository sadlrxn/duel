package transaction

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// ========================== Session management is reserved ==========================
// // Once a transaction is started, a channel is created with the key as
// // corresponding user ID.
// // Once the transaction is finished, a `true` value is pushed to the channel.
// var isIdle map[User](chan bool)

// // The size of isIdle map.
// const MAXIMUM_TRANSACTION_ROOM = 1024

// // @External
// // Initializes isIdle and db_aggregator.
// func initialize(db *gorm.DB) error {
// 	if err := db_aggregator.Initialize(db); err != nil {
// 		return err
// 	}

// 	isIdle = make(map[User](chan bool), MAXIMUM_TRANSACTION_ROOM)
// 	return nil
// }

// // @External
// // Returns whether a user session created or not.
// func isUserSessionStarted(user *User) bool {
// 	_, prs := isIdle[*user]
// 	return prs
// }

// // @External
// // Starts user session
// func startUserSession(user *User) error {
// 	if isUserSessionStarted(user) {
// 		return errors.New("already started session for the user")
// 	}

// 	isIdle[*user] = make(chan bool)
// 	return nil
// }

// // @External
// // Ends user session
// func endUserSession(user *User) error {
// 	if !isUserSessionStarted(user) {
// 		return errors.New("no session for the user")
// 	}

// 	close(isIdle[*user])
// 	return nil
// }
// ========================== Session management is reserved ==========================

// @External
// Initializes isIdle and db_aggregator.
func initialize(db *gorm.DB) error {
	if err := db_aggregator.Initialize(db); err != nil {
		return err
	}
	return nil
}

// @External
// Checks whether transaction is fee type.
func isFeeTransaction(txType models.TransactionType) bool {
	if txType == models.TxJackpotFee ||
		txType == models.TxCoinflipFee ||
		txType == models.TxGrandJackpotFee ||
		txType == models.TxDreamtowerFee ||
		txType == models.TxCrashFee {
		return true
	}
	return false
}

// @Internal
// Checks whether transaction is withdraw type.
func isWithdrawTransaction(txType models.TransactionType) bool {
	if txType == models.TxWithdrawSol ||
		txType == models.TxWithdrawNft ||
		txType == models.TxWithdrawSpl {
		return true
	}
	return false
}

// @Internal
// Get withdraw fee internal.
func addWithdrawFeeToRequest(transactionRequest *TransactionRequest) (int64, error) {
	if !isWithdrawTransaction(transactionRequest.Type) {
		return 0, nil
	}
	if transactionRequest.Type == models.TxWithdrawSol {
		return 0, nil
	}

	chipTransferWithFee := int64(0)
	if transactionRequest.Balance.ChipBalance != nil {
		chipTransferWithFee = *transactionRequest.Balance.ChipBalance
	}
	withdrawFee := int64(0)
	if transactionRequest.Type == models.TxWithdrawSpl {
		withdrawFee = config.WITHDRAW_FEE_PER_SPL
	} else if transactionRequest.Type == models.TxWithdrawNft {
		if transactionRequest.Balance.NftBalance == nil {
			return 0, utils.MakeError(
				"transaction",
				"addWithdrawFeeToRequest",
				"failed to get withdraw fee",
				errors.New("nft balance should not be null for nft withdraw tx"),
			)
		}
		withdrawFee = config.WITHDRAW_FEE_PER_SPL * int64(len(*transactionRequest.Balance.NftBalance))
	}
	chipTransferWithFee += withdrawFee
	transactionRequest.Balance.ChipBalance = &chipTransferWithFee

	return withdrawFee, nil
}

func shouldDistributeRevShare(txRequest *TransactionRequest) bool {
	if txRequest == nil {
		return false
	}

	return !txRequest.EdgeDetails.ShouldNotDistributeRevShare
}

// @Internal
// Handles transfer request.
func transfer(transactionRequest *TransactionRequest) (*db_aggregator.Transaction, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	if isFeeTransaction(transactionRequest.Type) &&
		transactionRequest.Balance.ChipBalance != nil {
		var totalDistributed int64
		var burnAmount int64

		// Distribute staking rewards to duelbots.
		if shouldDistributeRevShare(transactionRequest) {
			if distributed, err := db_aggregator.DistributeFeeToDuelBots(
				*transactionRequest.Balance.ChipBalance,
				sessionId,
			); err != nil {
				log.LogMessage(
					"transaction",
					"failed to distribute staking rewards",
					"error",
					logrus.Fields{
						"error": err.Error(),
					},
				)
			} else {
				totalDistributed += distributed
				burnAmount += distributed
			}
		} else {
			// This instruction assumes that distributin rev share is at the first part of
			// distribution actions for the fee transaction.
			totalDistributed += utils.RevShareFromFee(*transactionRequest.Balance.ChipBalance)
		}

		// Prepare batch house fee meta slice.
		batchHouseFeeMeta := []HouseFeeMeta{}
		if transactionRequest.BatchHouseFeeMeta != nil {
			batchHouseFeeMeta = append(
				batchHouseFeeMeta,
				transactionRequest.BatchHouseFeeMeta...,
			)
		}
		if transactionRequest.HouseFeeMeta != nil {
			batchHouseFeeMeta = append(
				batchHouseFeeMeta,
				HouseFeeMeta{
					User:        transactionRequest.HouseFeeMeta.User,
					WagerAmount: transactionRequest.HouseFeeMeta.WagerAmount,
					FeeAmount:   *transactionRequest.Balance.ChipBalance,
				},
			)
		}

		// Distribute rakeback.
		{
			for _, houseFeeMeta := range batchHouseFeeMeta {
				distributed, err := db_aggregator.DistributeRakeback(
					houseFeeMeta.User,
					houseFeeMeta.FeeAmount,
					sessionId,
				)
				if err != nil {
					log.LogMessage(
						"transaction",
						"failed to distribute rakeback",
						"error",
						logrus.Fields{
							"error":        err.Error(),
							"houseFeeMeta": houseFeeMeta,
						},
					)
				} else {
					totalDistributed += distributed
					burnAmount += distributed
				}
			}
		}

		// Distribute affiliate. Assuming that distribution for rakeback also means affiliate.
		{
			for _, houseFeeMeta := range batchHouseFeeMeta {
				distributed, err := db_aggregator.DistributeAffiliate(
					houseFeeMeta.User,
					houseFeeMeta.FeeAmount,
					houseFeeMeta.WagerAmount,
					sessionId,
				)
				if err != nil {
					log.LogMessage(
						"transaction",
						"failed to distribute affiliate",
						"error",
						logrus.Fields{
							"error": err.Error(),
						},
					)
				} else {
					totalDistributed += distributed
					burnAmount += distributed
				}
			}
		}

		if err := db_aggregator.Burn(
			transactionRequest.FromUser,
			burnAmount,
			sessionId,
		); err != nil {
			return nil, err
		}

		*transactionRequest.Balance.ChipBalance -= totalDistributed
	}

	// Include fee for withdraw transactions
	withdrawFee, err := addWithdrawFeeToRequest(transactionRequest)
	if err != nil {
		log.LogMessage(
			"transaction",
			"transfer",
			"error",
			logrus.Fields{
				"message": "failed to add withdraw fee to transaction request",
				"error":   err.Error(),
			},
		)
	} else if withdrawFee > 0 {
		log.LogMessage(
			"transaction",
			"transfer",
			"info",
			logrus.Fields{
				"message": "withdraw fee is added to transaction",
				"type":    transactionRequest.Type,
				"chips":   transactionRequest.Balance.ChipBalance,
				"nfts":    transactionRequest.Balance.NftBalance,
				"fee":     withdrawFee,
			},
		)
	}

	var toUser *db_aggregator.User = nil
	if transactionRequest.ToBeConfirmed {
		toUser = transactionRequest.ToUser
	}

	transferResult, err := db_aggregator.Transfer(
		transactionRequest.FromUser,
		toUser,
		&transactionRequest.Balance,
		sessionId,
	)
	if err != nil {
		return nil, err
	}

	fromWallet, err := db_aggregator.GetUserWallet(transactionRequest.FromUser, sessionId)
	if err != nil {
		return nil, err
	}
	toWallet, err := db_aggregator.GetUserWallet(transactionRequest.ToUser, sessionId)
	if err != nil {
		return nil, err
	}

	transaction, err := db_aggregator.RecordTransaction(&db_aggregator.TransactionLoad{
		FromWallet:       fromWallet,
		ToWallet:         toWallet,
		Balance:          transactionRequest.Balance,
		Type:             transactionRequest.Type,
		FromWalletPrevID: transferResult.FromPrevBalance,
		FromWalletNextID: transferResult.FromNextBalance,
	}, sessionId)
	if err != nil {
		return nil, err
	}
	if transactionRequest.ToBeConfirmed {
		if err := db_aggregator.ConfirmTransaction(&db_aggregator.TransactionLoad{
			ToWalletPrevID: transferResult.ToPrevBalance,
			ToWalletNextID: transferResult.ToNextBalance,
			OwnerID:        transactionRequest.OwnerID,
			OwnerType:      transactionRequest.OwnerType,
		}, transaction, sessionId); err != nil {
			return nil, err
		}
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, err
	}

	return transaction, nil
}

// @Internal
// Confirms a transaction. Currently only supports for `MainWallet` type.
func confirm(confirmRequest ConfirmRequest) error {
	if confirmRequest.Transaction == 0 {
		return errors.New("not invalid transaction id")
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	transactionLoad, err := db_aggregator.GetTransactionLoad(&confirmRequest.Transaction, sessionId)
	if err != nil {
		return err
	}

	if transactionLoad.Status != models.TransactionStatus(models.TransactionPending) {
		return errors.New("not a pending transaction")
	}

	if transactionLoad.ToWallet == nil {
		if err := db_aggregator.ConfirmTransaction(
			&db_aggregator.TransactionLoad{
				OwnerID:   confirmRequest.OwnerID,
				OwnerType: confirmRequest.OwnerType,
			},
			&confirmRequest.Transaction,
			sessionId); err != nil {
			return err
		}

		if err := db_aggregator.CommitSession(sessionId); err != nil {
			return err
		}

		return nil
	}

	toUser, err := db_aggregator.GetWalletUser(transactionLoad.ToWallet, sessionId)
	if err != nil {
		return err
	}

	transferResult, err := db_aggregator.Transfer(nil, toUser, &transactionLoad.Balance, sessionId)
	if err != nil {
		return err
	}

	if err := db_aggregator.ConfirmTransaction(&db_aggregator.TransactionLoad{
		ToWalletPrevID: transferResult.ToPrevBalance,
		ToWalletNextID: transferResult.ToNextBalance,
		OwnerID:        confirmRequest.OwnerID,
		OwnerType:      confirmRequest.OwnerType,
	}, &confirmRequest.Transaction, sessionId); err != nil {
		return err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return err
	}

	return nil
}

// @Internal
// Cancels a transaction.
func decline(declineRequest DeclineRequest) error {
	if declineRequest.Transaction == 0 {
		return errors.New("not invalid transaction id")
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	transactionLoad, err := db_aggregator.GetTransactionLoad(&declineRequest.Transaction, sessionId)
	if err != nil {
		return err
	}

	if transactionLoad.Status != models.TransactionStatus(models.TransactionPending) {
		return errors.New("not a pending transaction")
	}

	fromUser, err := db_aggregator.GetWalletUser(transactionLoad.FromWallet, sessionId)
	if err != nil {
		return err
	}

	transferResult, err := db_aggregator.Transfer(nil, fromUser, &transactionLoad.Balance, sessionId)
	if err != nil {
		return err
	}

	if err := db_aggregator.DeclineTransaction(&db_aggregator.TransactionLoad{
		RefundPrevID: transferResult.ToPrevBalance,
		RefundNextID: transferResult.ToNextBalance,
		OwnerID:      declineRequest.OwnerID,
		OwnerType:    declineRequest.OwnerType,
	}, &declineRequest.Transaction, sessionId); err != nil {
		return err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return err
	}

	return nil
}

// @External
// Stake duel bots.
func stakeDuelBots(request DuelBotsRequest) error {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	if err := db_aggregator.StakeDuelBots(
		request.FromUser,
		request.DuelBots,
		sessionId,
	); err != nil {
		return err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return err
	}

	return nil
}

// @External
// Unstake duel bots.
func unstakeDuelBots(request DuelBotsRequest) (int64, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	rewards, err := db_aggregator.UnstakeDuelBots(
		request.FromUser,
		request.DuelBots,
		sessionId,
	)
	if err != nil {
		return 0, err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, err
	}

	return rewards, nil
}

// @External
// Claim duel bot rewards.
func claimDuelBotsRewards(request DuelBotsRequest) (int64, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	rewards, err := db_aggregator.ClaimDuelBotsRewards(
		request.FromUser,
		request.DuelBots,
		sessionId,
	)
	if err != nil {
		return 0, err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, err
	}

	return rewards, nil
}

// @External
// Handles rain request
func rain(rainRequest *RainRequest) (*db_aggregator.Transaction, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	rainResult, err := db_aggregator.Rain(
		rainRequest.FromUser,
		rainRequest.ToUsers,
		&rainRequest.Balance,
		sessionId,
	)
	if err != nil {
		return nil, err
	}

	fromWallet, err := db_aggregator.GetUserWallet(rainRequest.FromUser, sessionId)
	if err != nil {
		return nil, err
	}

	var receipients []uint
	for _, toUser := range *rainRequest.ToUsers {
		toWallet, err := db_aggregator.GetUserWallet(&toUser, sessionId)
		if err != nil {
			return nil, err
		}
		receipients = append(receipients, uint(*toWallet))
	}

	transaction, err := db_aggregator.RecordTransaction(&db_aggregator.TransactionLoad{
		FromWallet:       fromWallet,
		ToWallet:         nil,
		Balance:          rainRequest.Balance,
		Type:             rainRequest.Type,
		FromWalletPrevID: rainResult.FromPrevBalance,
		FromWalletNextID: rainResult.FromNextBalance,
		Receipients:      &receipients,
	}, sessionId)
	if err != nil {
		return nil, err
	}

	if err := db_aggregator.ConfirmTransaction(&db_aggregator.TransactionLoad{
		ToWalletPrevID: nil,
		ToWalletNextID: nil,
		OwnerID:        uint(*fromWallet),
		OwnerType:      models.TransactionWalletReferenced,
	}, transaction, sessionId); err != nil {
		return nil, err
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, err
	}

	return transaction, nil
}
