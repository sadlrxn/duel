package coupon

import (
	"fmt"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
)

// To Do
// @External
// For user.
// Exchanges bonus balance to real chip after reaching requirements.
// Maximum exchangeable amount is COUPON_MAXIMUM_EXCHANGE.
// Returns exchanged amount and error object.
// Returns error on
//   - Provided code is not active for user. ErrCode: ErrCodeNotActiveCodeForExchange
//   - Wager limit is not reached to limit for exchanging. ErrCode: ErrCodeNotReachingExchangeWager
//   - On transfer failure, specifically returns error with code. ErrCode: ErrCodeExchangeTransactionFailure
func Exchange(userID uint, code uuid.UUID) (int64, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		code == uuid.Nil {
		return 0, utils.MakeErrorWithCode(
			"coupon_exchange",
			"Exchange",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"userID: %d, code: %v",
				userID, code,
			),
		)
	}

	// 2. Check for existence of playing rounds by coupon.
	if existingRoundWithCoupon(userID) {
		return 0, utils.MakeErrorWithCode(
			"coupon_exchange",
			"Exchange",
			"existing playing round",
			ErrCodeExistingPlayingRounds,
			fmt.Errorf(
				"userID: %d, code: %v",
				userID, code,
			),
		)
	}

	// 2. Start a session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"Exchange",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Retrieve active coupon record.
	activeCoupon, err := lockAndRetrieveActiveCoupon(userID, sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"Exchange",
			fmt.Sprintf(
				"failed to retrieve active coupon. userID: %d",
				userID,
			),
			err,
		)
	}

	// 4. Check whether provided code is equal to active code.
	if *activeCoupon.Coupon.Code != code {
		return 0, utils.MakeErrorWithCode(
			"coupon_exchange",
			"Exchange",
			"active code mismatching with request code",
			ErrCodeNotActiveCodeForExchange,
			fmt.Errorf(
				"userID: %d, codeParam: %v, activeCode: %v",
				userID, code, *activeCoupon.Coupon.Code,
			),
		)
	}

	// 5. Check whether wager limit is reached requirement.
	if activeCoupon.Wagered < int64(config.COUPON_REQUIRED_WAGER_TIMES)*activeCoupon.Coupon.BonusBalance {
		return 0, utils.MakeErrorWithCode(
			"coupon_exchange",
			"Exchange",
			"wager limit not reached to requirement",
			ErrCodeNotReachingExchangeWager,
			fmt.Errorf(
				"activeCode: %v, claimed: %v, wagered: %v",
				code, activeCoupon.Coupon.BonusBalance, activeCoupon.Wagered,
			),
		)
	}

	// 6. Perform coupon transaction. Remove the whole remaining coupon balance.
	if activeCoupon.Balance == 0 {
		return 0, nil
	}
	exchangeBalance := activeCoupon.Balance
	if exchangeBalance > config.COUPON_MAXIMUM_EXCHANGE {
		exchangeBalance = config.COUPON_MAXIMUM_EXCHANGE
	}
	couponTxID, err := performTransactionInSession(
		CouponTransactionRequest{
			Type:          models.CpTxExchangeToChip,
			UserID:        userID,
			Balance:       activeCoupon.Balance,
			ToBeConfirmed: true,
		},
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"Exchange",
			"failed to perform coupon transaction",
			err,
		)
	}

	// 7. Update exchanged field in active coupon field.
	if err := updateExchangedAmountUnchecked(
		exchangeBalance,
		activeCoupon,
		sessionId,
	); err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"Exchang",
			"failed to update exchanged amount",
			err,
		)
	}

	// 8. Perform real chips transaction.
	if _, err := giveRealChipForExchange(
		userID,
		exchangeBalance,
		couponTxID,
		sessionId,
	); err != nil {
		return 0, utils.MakeErrorWithCode(
			"coupon_exchange",
			"Exchange",
			"failed to give real chips for exchange",
			ErrCodeExchangeTransactionFailure,
			err,
		)
	}

	// 9. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"Exchang",
			"failed to commit session",
			err,
		)
	}

	return exchangeBalance, nil
}

// @Internal
// Perform real chips transfer for exchange transaction.
// Returns generated txID, and error object.
func giveRealChipForExchange(
	userID uint,
	balance int64,
	couponTxID uint,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		balance <= 0 ||
		couponTxID == 0 {
		return 0, utils.MakeError(
			"coupon_exchange",
			"giveRealChipForExchange",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, balance: %d, couponTxID: %d",
				userID, balance, couponTxID,
			),
		)
	}

	// 2. Give out real chips to user.
	realTxResult, err := db_aggregator.Transfer(
		(*db_aggregator.User)(&config.COUPON_TEMP_ID),
		(*db_aggregator.User)(&userID),
		&db_aggregator.BalanceLoad{
			ChipBalance: &balance,
		},
		sessionId,
	)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient funds") {
			return 0, utils.MakeErrorWithCode(
				"coupon_exchange",
				"giveRealChipForExchange",
				"insufficient admin temp wallet balance",
				ErrCodeInsufficientAdminBalance,
				err,
			)
		}
		return 0, utils.MakeError(
			"coupon_exchange",
			"giveRealChipForExchange",
			"failed to perform real chips transfer",
			err,
		)
	}

	// 3. Leave transaction.
	transactionHistory := models.Transaction{
		FromWallet: (*uint)(realTxResult.FromWallet),
		ToWallet:   (*uint)(realTxResult.ToWallet),
		Balance: models.Balance{
			ChipBalance: &models.ChipBalance{
				Balance: balance,
			},
		},
		Type:   models.TxExchangeCouponToChips,
		Status: models.TransactionSucceed,

		FromWalletPrevID: (*uint)(realTxResult.FromPrevBalance),
		FromWalletNextID: (*uint)(realTxResult.FromNextBalance),
		ToWalletPrevID:   (*uint)(realTxResult.ToPrevBalance),
		ToWalletNextID:   (*uint)(realTxResult.ToNextBalance),
		OwnerID:          couponTxID,
		OwnerType:        models.TransactionCouponTransactionReferenced,
	}
	if err := db_aggregator.LeaveRealTransaction(
		&transactionHistory,
		sessionId,
	); err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"giveRealChipForExchange",
			"failed to leave transaction",
			err,
		)
	}

	return transactionHistory.ID, nil
}

/**
* @Internal
* Checks whether the provided user has running game placed
* bet with bonus balance.
 */
func existingRoundWithCoupon(
	userID uint,
) bool {
	return existingDreamtowerRoundWithCoupon(userID) ||
		existingCrashRoundWithCoupon(userID)
}
