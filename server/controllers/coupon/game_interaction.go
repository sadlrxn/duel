package coupon

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// To Do
// @External
// Try bet with coupon balance first.
// Returns
//   - `CouponBetUnavailable`: No active coupon code, should bet with real chip.
//   - `CouponBetSucceed`: Bet with coupon balance succeed.
//   - `CouponBetInsufficientFunds`: Active coupon, but insufficient balance.
//   - `CouponBetFailed`: There is active coupon, but failed with unknown error.
func TryBet(request TryBetWithCouponRequest) (TryBetWithCouponResult, uint, error) {
	// 1. Validate parameter.
	if request.UserID == 0 ||
		request.Balance <= 0 ||
		!isSupportedTxTypeByTryBet(request.Type) {
		return CouponBetFailed, 0, utils.MakeErrorWithCode(
			"coupon_game_interaction",
			"TryBet",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("request: %v", request),
		)
	}

	// 2. Try perform coupon transaction with the type.
	if txId, err := Perform(
		CouponTransactionRequest{
			Type:          request.Type,
			UserID:        request.UserID,
			Balance:       request.Balance,
			ToBeConfirmed: false,
		},
	); err == nil {
		return CouponBetSucceed, txId, nil
	} else if utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		return CouponBetUnavailable, 0, utils.MakeError(
			"coupon_game_interaction",
			"TryBet",
			"not found active coupon",
			fmt.Errorf(
				"request: %v, err: %v",
				request, err,
			),
		)
	} else if utils.IsErrorCode(err, ErrCodeInsufficientBonusBalance) {
		return CouponBetInsufficientFunds, 0, utils.MakeError(
			"coupon_game_interaction",
			"TryBet",
			"insufficient bonus balance",
			fmt.Errorf(
				"request: %v, err: %v",
				request, err,
			),
		)
	} else if utils.IsErrorCode(err, ErrCodeZeroBonusBalance) {
		return CouponBetUnavailable, 0, utils.MakeError(
			"coupon_game_interaction",
			"TryBet",
			"zero bonus balance",
			fmt.Errorf(
				"request: %v, err: %v",
				request, err,
			),
		)
	} else {
		return CouponBetFailed, 0, utils.MakeError(
			"coupon_game_interaction",
			"TryBet",
			"failed to try perform coupon transaction",
			fmt.Errorf(
				"request: %v, err: %v",
				request, err,
			),
		)
	}
}

// To Do
// @Internal
// Check whether transaction type is supported by try bet.
func isSupportedTxTypeByTryBet(txType models.CouponTransactionType) bool {
	return txType == models.CpTxCoinflipBet ||
		txType == models.CpTxDreamtowerBet ||
		txType == models.CpTxCrashBet
}
