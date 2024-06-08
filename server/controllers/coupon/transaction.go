package coupon

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// @External
// Wraps performTransactionInSession function within main session.
// Returns error in case of transaction type is claim or exchange.
func Perform(transactionRequest CouponTransactionRequest) (uint, error) {
	if !isSupportedByPerform(transactionRequest.Type) {
		return 0, utils.MakeError(
			"coupon_transaction",
			"Perform",
			"this transaction is not supported by external perform",
			fmt.Errorf("transaction type: %s", transactionRequest.Type),
		)
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"coupon_transaction",
			"Perform",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	tx, err := performTransactionInSession(transactionRequest, sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_transaction",
			"Perform",
			"failed to perform transaction in session",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"coupon_transaction",
			"Perform",
			"failed to commit session",
			err,
		)
	}

	return tx, nil
}

// @Internal
// Checks whether transaction type is supported by performTransaction function.
// For now, all transactions except claim and exchange are supported.
func isSupportedByPerform(transactionType models.CouponTransactionType) bool {
	return transactionType == models.CpTxCoinflipBet ||
		transactionType == models.CpTxCoinflipProfit ||
		transactionType == models.CpTxDreamtowerBet ||
		transactionType == models.CpTxDreamtowerProfit ||
		transactionType == models.CpTxCrashBet ||
		transactionType == models.CpTxCrashProfit
}

// To Do
// @Internal
// Performs transaction changing user's coupon balance and leaving coupon transaction history.
// Returns recorded coupon transaction id and error object.
//   - Balance is zero.
//   - Active coupon is not found.
//   - For adding balance transactions like claim, and profit, we assume that the ToBeConfirmed
//     is always `true`.
func performTransactionInSession(
	transactionRequest CouponTransactionRequest,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter.
	if transactionRequest.Balance <= 0 {
		return 0, utils.MakeErrorWithCode(
			"coupon_transaction",
			"performTransactionInSession",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided perform balance should not be absent"),
		)
	}
	// !Archived for pending confirmed
	// if isAddingBalanceTransaction(transactionRequest.Type) &&  !transactionRequest.ToBeConfirmed {
	// 	return 0, utils.MakeError(
	// 		"coupon_transaction",
	// 		"performTransactionInSession",
	// 		"invalid parameter",
	// 		errors.New("for adding balance tx, toBeConfirmed should always be `true`"),
	// 	)
	// }

	// 3. Lock and retrieve user's active coupon record.
	activeCoupon, err := lockAndRetrieveActiveCoupon(transactionRequest.UserID, sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_transaction",
			"performTransactionInSession",
			"failed to lock and retrieve active coupon",
			err,
		)
	}

	// 4. Save current coupon record.
	backupBalance := activeCoupon.Balance

	// 5. Perform balance change. Mint or burn.
	if isAddingBalanceTransaction(transactionRequest.Type) {
		if err := addBalanceToClaimedCouponUnchecked(
			transactionRequest.Balance,
			activeCoupon,
			sessionId,
		); err != nil {
			return 0, utils.MakeError(
				"coupon_transaction",
				"performTransactionInSession",
				"failed to add balance to active coupon",
				err,
			)
		}
	} else {
		if err := removeBalanceFromClaimedCouponUnchecked(
			transactionRequest.Balance,
			activeCoupon,
			sessionId,
		); err != nil {
			return 0, utils.MakeError(
				"coupon_transaction",
				"performTransactionInSession",
				"failed to remove balance to active coupon",
				err,
			)
		}
	}

	// 6. Perform wager amount change.
	// ============= Migrated wager amount increase to `Confirm` method =============
	// if isWagerTransaction(transactionRequest.Type) {
	// 	if err := addWagerAmountUnchecked(
	// 		transactionRequest.Balance,
	// 		activeCoupon,
	// 		sessionId,
	// 	); err != nil {
	// 		return 0, utils.MakeError(
	// 			"coupon_transaction",
	// 			"performTransactionInSession",
	// 			"failed to add wager amount",
	// 			err,
	// 		)
	// 	}
	// }

	// 7. Leave transaction.
	status := models.CouponTransactionSucceed
	if !transactionRequest.ToBeConfirmed {
		status = models.CouponTransactionPending
	}
	tx, err := leaveTransactionHistory(
		&models.CouponTransaction{
			CouponID:      activeCoupon.CouponID,
			ClaimedUserID: activeCoupon.ClaimedUserID,
			PrevBalance:   backupBalance,
			TxBalance:     transactionRequest.Balance,
			NextBalance:   activeCoupon.Balance,
			Status:        status,
			Type:          transactionRequest.Type,
		},
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_transaction",
			"performTransactionInSession",
			"failed to leave transaction history",
			err,
		)
	}

	return tx, nil
}

// @Internal
// Checks whether transaction type is to add balance.
func isAddingBalanceTransaction(transactionType models.CouponTransactionType) bool {
	return transactionType == models.CpTxClaimCode ||
		transactionType == models.CpTxCoinflipProfit ||
		transactionType == models.CpTxDreamtowerProfit ||
		transactionType == models.CpTxCrashProfit
}

// @Internal
// Checks whether transaction type is wager one.
func isWagerTransaction(transactionType models.CouponTransactionType) bool {
	return transactionType == models.CpTxCoinflipBet ||
		transactionType == models.CpTxDreamtowerBet ||
		transactionType == models.CpTxCrashBet
}

// To Do
// @External
// Confirm pending transaction.
// Confirm is done with in main session without locking any record.
func Confirm(tx uint) error {
	// 1. Validate parameter.
	if tx == 0 {
		return utils.MakeErrorWithCode(
			"coupon_transaction",
			"Confirm",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided tx id is zero"),
		)
	}

	// 2. Get main session id.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Confirm",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Confirm transaction.
	cpTx, err := confirmTransactionStatus(
		tx, sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Confirm",
			"failed to confirm transaction",
			fmt.Errorf("tx: %d, err: %v", tx, err),
		)
	}

	// 4. Lock and retrieve active coupon.
	activeCoupon, err := lockAndRetrieveActiveCoupon(
		cpTx.ClaimedUserID,
		sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Confirm",
			"failed to lock andretrieve active coupon",
			err,
		)
	}
	if activeCoupon.CouponID != cpTx.CouponID {
		return utils.MakeError(
			"coupon_transaction",
			"performTransactionInSession",
			"retrieved active coupon mismatching with tx one",
			fmt.Errorf(
				"active: %v, transaction: %v",
				activeCoupon,
				cpTx,
			),
		)
	}

	// 5. Perform wager amount change.
	// ============= Migrated from `performTransactionInSession` method =============
	if isWagerTransaction(cpTx.Type) {
		if err := addWagerAmountUnchecked(
			cpTx.TxBalance,
			activeCoupon,
			sessionId,
		); err != nil {
			return utils.MakeError(
				"coupon_transaction",
				"performTransactionInSession",
				"failed to add wager amount",
				err,
			)
		}
	}

	// 6. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"performTransactionInSession",
			"failed to commit session",
			err,
		)
	}

	return nil
}

// To Do
// @External
// Decline pending transaction.
// Decline is done after locking activeCoupon record.
func Decline(tx uint) error {
	// 1. Validate parameter.
	if tx == 0 {
		return utils.MakeErrorWithCode(
			"coupon_transaction",
			"Decline",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided tx id is zero"),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve transaction record.
	transaction, err := lockAndRetrievePendingTransaction(tx, sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"failed to lock and retrieve transaction",
			fmt.Errorf("tx: %d, err: %v", tx, err),
		)
	}

	// 4. Lock and retrieve user's active coupon.
	activeCoupon, err := lockAndRetrieveActiveCoupon(
		transaction.ClaimedUserID,
		sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"failed to lock and retrieve active coupon",
			fmt.Errorf(
				"userID: %d, err: %v",
				transaction.ClaimedUserID,
				err,
			),
		)
	}
	if activeCoupon.CouponID != transaction.CouponID {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"transaction coupon mismatching with active one",
			fmt.Errorf(
				"transaction coupon: %s, active coupon: %s",
				transaction.CouponID,
				activeCoupon.CouponID,
			),
		)
	}

	// 5. Refund bonus balance.
	if isAddingBalanceTransaction(transaction.Type) {
		if err := removeBalanceFromClaimedCouponUnchecked(
			transaction.TxBalance,
			activeCoupon,
			sessionId,
		); err != nil {
			return utils.MakeError(
				"coupon_transaction",
				"Decline",
				"failed to remove balance from user",
				fmt.Errorf(
					"refund balance: %d, activeCoupon: %v",
					transaction.TxBalance,
					activeCoupon,
				),
			)
		}
	} else {
		if err := addBalanceToClaimedCouponUnchecked(
			transaction.TxBalance,
			activeCoupon,
			sessionId,
		); err != nil {
			return utils.MakeError(
				"coupon_transaction",
				"Decline",
				"failed to add balance to user",
				fmt.Errorf(
					"refund balance: %d, activeCoupon: %v",
					transaction.TxBalance,
					activeCoupon,
				),
			)
		}
	}

	// 6. Decline transaction history.
	if err := declineTransactionStatus(
		tx, activeCoupon.Balance, sessionId,
	); err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"failed to decline transaction history",
			fmt.Errorf(
				"txID: %d, err: %v",
				tx, err,
			),
		)
	}

	// 7. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"coupon_transaction",
			"Decline",
			"failed to commit session",
			err,
		)
	}

	return nil
}
