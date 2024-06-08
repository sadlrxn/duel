package crash

import (
	"fmt"
	"math"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/*
	@Internal

/* Handle cash in from user to round.
/* Returns betID and error object.
/* Possible errors:
  - Game status is not betting, and pending: `ErrCodeNotStatusForCashIn`
  - Invalid parameter or mismatching round id: `ErrCodeInvlidParameter`
  - Max wager amount exceed: `ErrCodeMoreThanMaxBetAmount`
  - Less than min wager amount: `ErrCodeLessThanMinBetAmount`
  - Max bet count limit exceed: `ErrCodeMaxBetCountExceed`
  - Insufficient user balance: `ErrCodeInsufficientUserBalance`
  - Bet balance type mismatching: `ErrCodeBalanceTypeMismatching`
*/
func (c *GameController) cashIn(params CashInRequestParams) (uint, error) {
	// 1. Check game status.
	if !c.isStatusForCashIn() ||
		!c.isBettingRound() {
		return 0, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashIn",
			"not status for cash in",
			ErrCodeNotStatusForCashIn,
			fmt.Errorf(
				"request: %v, status: %v, round: %v",
				params, c.roundStatus, c.round,
			),
		)
	}

	// 2. Validate parameter.
	if params.UserID == 0 ||
		params.Amount <= 0 ||
		params.RoundID != c.round.ID ||
		(params.CashOutAt != 0 &&
			params.CashOutAt < c.minCashOutAt) {
		return 0, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashIn",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("params: %v", params),
		)
	}

	// 3. Check wager amount range.
	if c.greaterThanMaxWager(params.Amount) {
		return 0, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashIn",
			"exceeds max bet amount",
			ErrCodeMaxBetCountExceed,
			fmt.Errorf("params: %v", params),
		)
	}
	if c.lessThanMinWager(params.Amount) {
		return 0, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashIn",
			"less than min bet amount",
			ErrCodeLessThanMinBetAmount,
			fmt.Errorf("params: %v", params),
		)
	}

	// 4. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 5. Lock and retrieve crash round record.
	round, err := lockAndRetrieveCrashRound(
		c.round.ID,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"failed to lock and retrieve crash round",
			fmt.Errorf(
				"roundId: %d, err: %v",
				c.round.ID, err,
			),
		)
	}

	// 6. Check status with retrieved round again.
	if round.BetStartedAt == nil ||
		round.RunStartedAt != nil {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"invalid retrieved round status for cash in",
			fmt.Errorf("round: %v", round),
		)
	}

	// 7. Check bet count made by this user to round.
	betCount, err := getBetCountMadeByUserForRound(
		params.UserID,
		c.round.ID,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"failed to get bet count made by user for round",
			fmt.Errorf(
				"userID: %d, roundID: %d, err: %v",
				params.UserID, c.round.ID, err,
			),
		)
	}
	if betCount >= c.betCountLimit {
		return 0, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashIn",
			"bet count limit exceeds",
			ErrCodeMaxBetCountExceed,
			fmt.Errorf(
				"params: %v, betCount: %d",
				params, betCount,
			),
		)
	}

	// 8. Try bet with coupon.
	couponBetResult, cpTxID, err := coupon.TryBet(coupon.TryBetWithCouponRequest{
		UserID:  params.UserID,
		Balance: params.Amount,
		Type:    models.CpTxCrashBet,
	})

	// 9. Check cash in balance type with coupon bet result.
	if params.BalanceType == models.ChipBalanceForGame {
		if couponBetResult != coupon.CouponBetUnavailable {
			if cpTxID > 0 {
				coupon.Decline(cpTxID)
			}
			return 0, utils.MakeErrorWithCode(
				"crash_controller_cash",
				"cashIn",
				"balance type mismatching, expected: chip",
				ErrCodeBalanceTypeMismatching,
				fmt.Errorf(
					"params: %v, cpRes: %v, cpTxID: %d, cpTxErr: %v",
					params, couponBetResult, cpTxID, err,
				),
			)
		}
	} else if params.BalanceType == models.CouponBalanceForGame {
		if couponBetResult == coupon.CouponBetUnavailable {
			return 0, utils.MakeErrorWithCode(
				"crash_controller_cash",
				"cashIn",
				"balance type mismatching, expected: coupon",
				ErrCodeBalanceTypeMismatching,
				fmt.Errorf(
					"params: %v, cpRes: %v, cpTxID: %d, cpTxErr: %v",
					params, couponBetResult, cpTxID, err,
				),
			)
		} else if couponBetResult == coupon.CouponBetInsufficientFunds {
			return 0, utils.MakeErrorWithCode(
				"crash_controller_cash",
				"cashIn",
				"insufficient coupon balance",
				ErrCodeInsufficientUserBalance,
				fmt.Errorf(
					"params: %v, cpTxErr: %v",
					params, err,
				),
			)
		} else if couponBetResult == coupon.CouponBetFailed {
			return 0, utils.MakeError(
				"crash_controller_cash",
				"cashIn",
				"failed to bet coupon balance",
				fmt.Errorf(
					"params: %v, cpTxErr: %v",
					params, err,
				),
			)
		}
	}

	// 10. Bet real chips in case of real chip bet.
	betAmount := params.Amount
	var txId *db_aggregator.Transaction
	if params.BalanceType == models.ChipBalanceForGame {
		txId, err = transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&params.UserID),
			ToUser:   (*db_aggregator.User)(&config.CRASH_TEMP_ID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &betAmount,
			},
			Type:          models.TxCrashBet,
			ToBeConfirmed: false,
		})
		if err != nil && strings.Contains(err.Error(), "insufficient funds") {
			return 0, utils.MakeErrorWithCode(
				"crash_controller_cash",
				"cashIn",
				"insufficient real chip balance",
				ErrCodeInsufficientUserBalance,
				fmt.Errorf(
					"params: %v, txErr: %v",
					params, err,
				),
			)
		} else if err != nil {
			return 0, utils.MakeError(
				"crash_controller_cash",
				"cashIn",
				"failed to bet real chip balance",
				fmt.Errorf(
					"params: %v, txErr: %v",
					params, err,
				),
			)
		}
	}

	// 11. Create bet record.
	crashBet := models.CrashBet{
		UserID:          params.UserID,
		RoundID:         c.round.ID,
		BetAmount:       params.Amount,
		PaidBalanceType: params.BalanceType,
	}
	if params.CashOutAt >= c.minCashOutAt {
		crashBet.CashOutAt = &params.CashOutAt
	}
	if err := createCrashBet(
		&crashBet,
		sessionId,
	); err != nil {
		if params.BalanceType == models.ChipBalanceForGame {
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *txId,
				OwnerType:   models.TransactionCrashRoundReferenced,
				OwnerID:     c.round.ID,
			})
		} else if params.BalanceType == models.CouponBalanceForGame {
			coupon.Decline(cpTxID)
		}
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"failed to create crash bet record",
			fmt.Errorf(
				"params: %v, crashBetRecord: %v, err: %v",
				params, crashBet, err,
			),
		)
	}

	// 12. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		var declineError error
		if params.BalanceType == models.ChipBalanceForGame {
			declineError = transaction.Decline(transaction.DeclineRequest{
				Transaction: *txId,
				OwnerType:   models.TransactionCrashRoundReferenced,
				OwnerID:     c.round.ID,
			})
		} else if params.BalanceType == models.CouponBalanceForGame {
			declineError = coupon.Decline(cpTxID)
		}
		return 0, utils.MakeError(
			"crash_controller_cash",
			"cashIn",
			"failed to commit session",
			fmt.Errorf(
				"params: %v, err: %v, declineErr: %v",
				params, err, declineError,
			),
		)
	}

	// 13. Add bet record to c.round.
	c.round.Bets = append(c.round.Bets, crashBet)

	// 14. Confirm transaction.
	// Should not be decline on failure because bet is already committed to round.
	if params.BalanceType == models.ChipBalanceForGame {
		if err := transaction.Confirm(transaction.ConfirmRequest{
			Transaction: *txId,
			OwnerType:   models.TransactionCrashBetReferencedForCashIn,
			OwnerID:     crashBet.ID,
		}); err != nil {
			return 0, utils.MakeError(
				"crash_controller_cash",
				"cashIn",
				"failed to confirm real chip bet",
				fmt.Errorf(
					"params: %v, txId: %d, err: %v",
					params, *txId, err,
				),
			)
		}
	} else if params.BalanceType == models.CouponBalanceForGame {
		if err := coupon.Confirm(cpTxID); err != nil {
			return 0, utils.MakeError(
				"crash_controller_cash",
				"cashIn",
				"failed to confirm coupon bet",
				fmt.Errorf(
					"params: %v, cpTxID: %d, err: %v",
					params, cpTxID, err,
				),
			)
		}
	}

	return crashBet.ID, nil
}

/*
	@Internal

/* Handle cash out from round to user.
/* Returns cashed out amount, balance type and error object.
/* Possible errors:
  - Invalid parameter or mismatching round id: `ErrCodeInvalidParameter`
  - Game status is not running, and preparing: `ErrCodeNotStatusForCashOut`
  - Invalid bet requirement with user and round: `ErrCodeInvalidBetForCashout`
  - Insufficient pool balance: `ErrCodeInsufficientPoolBalance`
*/
func (c *GameController) cashOut(params CashOutRequestParams) (*CashOutRequestResult, error) {
	// 0. Refine payoutMultiplier decimals to be 2.
	params.PayoutMultiplier = math.Floor(params.PayoutMultiplier*100) / 100

	// 1. Check game status.
	if !c.isStatusForCashOut() ||
		!c.isRunningRound() {
		return nil, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashOut",
			"not status for cash out",
			ErrCodeNotStatusForCashOut,
			fmt.Errorf(
				"request: %v, status: %v, round: %v",
				params, c.roundStatus, c.round,
			),
		)
	}

	// 2. Validate parameter.
	if params.UserID == 0 ||
		params.RoundID != c.round.ID ||
		params.BetID == 0 ||
		params.PayoutMultiplier > c.nextMultiplier ||
		params.PayoutMultiplier < c.minCashOutAt {
		return nil, utils.MakeErrorWithCode(
			"crash_controller_cash",
			"cashOut",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"request: %v, roundID: %d, curMul: %f, minCashoutAt: %f",
				params, c.round.ID, c.nextMultiplier, c.minCashOutAt,
			),
		)
	}

	// 3. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 4. Lock and retrieve crash round record.
	round, err := lockAndRetrieveCrashRound(
		c.round.ID,
		sessionId,
	)
	if err != nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"failed to lock and retrieve crash round",
			fmt.Errorf(
				"roundId: %d, err: %v",
				c.round.ID, err,
			),
		)
	}

	// 5. Check status with retrieved round again.
	if round.RunStartedAt == nil ||
		round.EndedAt != nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"invalid retrieved round status for cash out",
			fmt.Errorf("round: %v", round),
		)
	}

	// 6. Check whether multiplier is greater than outcome.
	if round.Outcome < params.PayoutMultiplier {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"greater payout multiplier than round outcome",
			fmt.Errorf(
				"round: %v, params: %v",
				round, params,
			),
		)
	}

	// 7. Get bet record for BetID.
	crashBet := (*models.CrashBet)(nil)
	for i, bet := range c.round.Bets {
		if bet.ID == params.BetID {
			crashBet = &c.round.Bets[i]
		}
	}
	if crashBet == nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"failed to find bet record in controller",
			fmt.Errorf(
				"c.round: %v, betID: %d",
				c.round, params.BetID,
			),
		)
	}

	// 8. Check whether crash bet is already paid out.
	if crashBet.Profit != nil ||
		crashBet.PayoutMultiplier != nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"already paid out",
			fmt.Errorf(
				"params: %v, crashBet: %v",
				params, crashBet,
			),
		)
	}

	// 9. Check crashBet's userID is matching with event's userID.
	if crashBet.UserID != params.UserID {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"not bet placed user",
			fmt.Errorf(
				"params: %v, crashBet: %v",
				params, crashBet,
			),
		)
	}

	// 10. Update crash bet's payout fields.
	profit := int64(float64(crashBet.BetAmount) * params.PayoutMultiplier)
	if profit > c.maxCashOut {
		profit = c.maxCashOut
	}
	if err := updateCrashBetPayoutFields(
		crashBet,
		profit,
		params.PayoutMultiplier,
		sessionId,
	); err != nil {
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"failed to update crash bet payout fields",
			fmt.Errorf(
				"crashBet: %v, profit: %d, params: %v",
				c.round, profit, params,
			),
		)
	}

	// 11. Payout profit.
	cpTxId := uint(0)
	txId := (*db_aggregator.Transaction)(nil)

	if crashBet.PaidBalanceType == models.CouponBalanceForGame {
		cpTxId, err = coupon.Perform(coupon.CouponTransactionRequest{
			Type:          models.CpTxCrashProfit,
			UserID:        params.UserID,
			Balance:       profit,
			ToBeConfirmed: false,
		})
		if err != nil {
			return nil, utils.MakeError(
				"crash_controller_cash",
				"cashOut",
				"failed to payout coupon profit",
				fmt.Errorf(
					"params: %v, profit: %d, err: %v",
					params, profit, err,
				),
			)
		}
	} else if crashBet.PaidBalanceType == models.ChipBalanceForGame {
		backupProfit := profit
		txId, err = transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&config.CRASH_TEMP_ID),
			ToUser:   (*db_aggregator.User)(&params.UserID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &backupProfit,
			},
			ToBeConfirmed: false,
			Type:          models.TxCrashProfit,
		})
		if err != nil && strings.Contains(err.Error(), "insufficient funds") {
			return nil, utils.MakeErrorWithCode(
				"crash_controller_cash",
				"cashOut",
				"insufficient pool balance to payout real chips profit",
				ErrCodeInsufficientPoolBalance,
				fmt.Errorf(
					"params: %v, profit: %d, err: %v",
					params, profit, err,
				),
			)
		} else if err != nil {
			return nil, utils.MakeError(
				"crash_controller_cash",
				"cashOut",
				"failed to payout real chips profit",
				fmt.Errorf(
					"params: %v, profit: %d, err: %v",
					params, profit, err,
				),
			)
		}
	}

	// 12. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		var declineError error
		if crashBet.PaidBalanceType == models.CouponBalanceForGame {
			declineError = coupon.Decline(cpTxId)
		} else if crashBet.PaidBalanceType == models.ChipBalanceForGame {
			declineError = transaction.Decline(transaction.DeclineRequest{
				Transaction: *txId,
				OwnerType:   models.TransactionCrashBetReferenced,
				OwnerID:     crashBet.ID,
			})
		}
		return nil, utils.MakeError(
			"crash_controller_cash",
			"cashOut",
			"failed to commit session",
			fmt.Errorf(
				"params: %v, commitError: %v, declineError: %v",
				params, err, declineError,
			),
		)
	}

	// 13. Confirm payout transaction.
	if crashBet.PaidBalanceType == models.CouponBalanceForGame {
		if err := coupon.Confirm(cpTxId); err != nil {
			return nil, utils.MakeError(
				"crash_controller_cash",
				"cashOut",
				"failed confirm coupon payout transaction",
				fmt.Errorf(
					"params: %v, err: %v",
					params, err,
				),
			)
		}
	} else if crashBet.PaidBalanceType == models.ChipBalanceForGame {
		if err := transaction.Confirm(transaction.ConfirmRequest{
			Transaction: *txId,
			OwnerType:   models.TransactionCrashBetReferencedForCashOut,
			OwnerID:     crashBet.ID,
		}); err != nil {
			return nil, utils.MakeError(
				"crash_controller_cash",
				"cashOut",
				"failed confirm real chip payout transaction",
				fmt.Errorf(
					"params: %v, err: %v",
					params, err,
				),
			)
		}
	}

	return &CashOutRequestResult{
		Amount:      profit,
		BalanceType: crashBet.PaidBalanceType,
	}, nil
}

/*
* @Internal
* Handle fee transfer from temp to fee.
 */
func (c *GameController) chargeFee() (int64, error) {
	// 1. Check game status.
	if !c.isPreparingStatus() ||
		!c.isEndedRound() {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"chargeFee",
			"not status for charge fee",
			fmt.Errorf(
				"round: %v, roundStatu: %v",
				c.round, c.roundStatus,
			),
		)
	}

	// 2. Calculate total chip bet amount and build batch house fee meta.
	totalFeeAmount := int64(0)
	batchHouseFeeMeta := []transaction.HouseFeeMeta{}
	for _, bet := range c.round.Bets {
		if bet.PaidBalanceType == models.ChipBalanceForGame {
			feeAmount := bet.BetAmount * c.houseEdge / 10000
			if feeAmount > 0 {
				totalFeeAmount += feeAmount
				batchHouseFeeMeta = append(
					batchHouseFeeMeta,
					transaction.HouseFeeMeta{
						User:        db_aggregator.User(bet.UserID),
						WagerAmount: bet.BetAmount,
						FeeAmount:   feeAmount,
					},
				)
			}
		}
	}
	if totalFeeAmount <= 0 {
		return 0, nil
	}

	// 3. Transfer fee.
	backupFee := totalFeeAmount
	if _, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.CRASH_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&config.CRASH_FEE_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &backupFee,
		},
		Type:              models.TxCrashFee,
		ToBeConfirmed:     true,
		OwnerType:         models.TransactionCrashRoundReferencedForFee,
		OwnerID:           c.round.ID,
		BatchHouseFeeMeta: batchHouseFeeMeta,
	}); err != nil {
		return 0, utils.MakeError(
			"crash_controller_cash",
			"chargeFee",
			"failed to transfer fee",
			fmt.Errorf(
				"fee: %d, err: %v",
				totalFeeAmount, err,
			),
		)
	}

	return totalFeeAmount, nil
}

/*
	@Internal

/* Checks whether amount is greater than max wager amount
*/
func (c *GameController) greaterThanMaxWager(amount int64) bool {
	return amount > c.maxBetAmount
}

/*
	@Internal

/* Checks whether amount is less than min wager amount
*/
func (c *GameController) lessThanMinWager(amount int64) bool {
	return amount < c.minBetAmount
}

/*
	@Internal

/* Checks whether game status is for cash in.
*/
func (c *GameController) isStatusForCashIn() bool {
	return c.roundStatus == Betting ||
		c.roundStatus == Pending

}

/*
	@Internal

/* Checks whether game status is for cash out.
*/
func (c *GameController) isStatusForCashOut() bool {
	return c.roundStatus == Running ||
		c.roundStatus == Preparing
}
