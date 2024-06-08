package coupon

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @Internal
// Lock and retrieve unexpired coupon record by code.
// Should lock and retrieve where
//   - Code is equal to provided one
//   - CreatedAt is greater than 14 days before now
func lockAndRetrieveCoupon(
	code uuid.UUID,
	sessionId db_aggregator.UUID,
) (*models.Coupon, error) {
	// 1. Validate parameter.
	if code == uuid.Nil {
		return nil, utils.MakeErrorWithCode(
			"coupon_db",
			"lockAndRetrieveCoupon",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided code is nil uuid"),
		)
	}

	// 2. Retrieve session
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveCoupon",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve coupon record.
	coupon := models.Coupon{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"ClaimedCoupons",
	).Preload(
		"RequiredAffiliate.ActiveAffiliates",
	).Where(
		"code = ?",
		code,
	).Where(
		"created_at > ?",
		time.Unix(time.Now().Unix()-int64(time.Hour.Seconds())*24*int64(config.COUPON_CODE_LIFE_TIME_IN_DAYS), 0),
	).Order(
		"created_at",
	).First(&coupon); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utils.MakeErrorWithCode(
			"coupon_db",
			"lockAndRetrieveCoupon",
			"failed to retrieve coupon record",
			ErrCodeCouponCodeNotFound,
			result.Error,
		)
	} else if result.Error != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveCoupon",
			"failed to retrieve coupon record",
			result.Error,
		)
	}

	return &coupon, nil
}

// @Internal
// Lock and retrieve user's active claimed coupon record.
// Should lock and retrieve where
//   - ClaimedUserID is equal to userID
//   - CreatedAt is greater than 8 hours before than now
//   - Exchanged is 0
func lockAndRetrieveActiveCoupon(
	userID uint,
	sessionId db_aggregator.UUID,
) (*models.ClaimedCoupon, error) {
	// 1. Validate parameter.
	if userID == 0 {
		return nil, utils.MakeErrorWithCode(
			"coupon_db",
			"lockAndRetrieveActiveCoupon",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user id is invalid"),
		)
	}

	// 2. Retrieve session
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveActiveCoupon",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve claimedCoupon record.
	claimedCoupon := models.ClaimedCoupon{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"Coupon",
	).Where(
		"claimed_user_id = ?",
		userID,
	).Where(
		"created_at > ?",
		getActiveCouponLifeTimeBeforeNow(),
	).Where(
		"exchanged = 0",
	).Order(
		"created_at",
	).First(&claimedCoupon); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utils.MakeErrorWithCode(
			"coupon_db",
			"lockAndRetrieveActiveCoupon",
			"failed to retrieve active coupon record",
			ErrCodeCouponCodeNotFound,
			result.Error,
		)
	} else if result.Error != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveActiveCoupon",
			"failed to retrieve active coupon record",
			result.Error,
		)
	}

	return &claimedCoupon, nil
}

// @Internal
// Lock and retrieve coupon transaction.
// Should lock and retrieve where
//   - ID is equal to provided txID
func lockAndRetrievePendingTransaction(
	txID uint,
	sessionId db_aggregator.UUID,
) (*models.CouponTransaction, error) {
	// 1. Validate parameter.
	if txID == 0 {
		return nil, utils.MakeErrorWithCode(
			"coupon_db",
			"lockAndRetrieveTransaction",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided txID is zero"),
		)
	}

	// 2. Retrieve session
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveTransaction",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve transaction record.
	transaction := models.CouponTransaction{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"status = ?",
		models.CouponTransactionPending,
	).Where(
		"id = ?",
		txID,
	).First(&transaction); result.Error != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveTransaction",
			"failed to retrieve transaction record",
			result.Error,
		)
	}

	return &transaction, nil
}

// @Internal
// Add balance to coupon record.
// Doesn't check about active coupon expiration or other logic.
// Returns updated balance.
func addBalanceToClaimedCouponUnchecked(
	balance int64,
	claimedCoupon *models.ClaimedCoupon,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if balance <= 0 {
		return utils.MakeError(
			"coupon_db",
			"addBalanceToClaimedCouponUnchecked",
			"invalid parameter",
			fmt.Errorf("provided balance: %d", balance),
		)
	}
	if claimedCoupon == nil {
		return utils.MakeError(
			"coupon_db",
			"addBalanceToClaimedCouponUnchecked",
			"invalid parameter",
			errors.New("claimedCoupon is nil pointer"),
		)
	}
	if claimedCoupon.CouponID == uuid.Nil ||
		claimedCoupon.ClaimedUserID == 0 {
		return utils.MakeError(
			"coupon_db",
			"addBalanceToClaimedCouponUnchecked",
			"invalid parameter",
			fmt.Errorf("provided claimedCoupon: %v", *claimedCoupon),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"addBalanceToClaimedCouponUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Add balance and save.
	if result := session.Model(
		&claimedCoupon,
	).Clauses(
		clause.Returning{},
	).Update(
		"balance", gorm.Expr("balance + ?", balance),
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"coupon_db",
			"addBalanceToClaimedCouponUnchecked",
			"failed add balance properly",
			result.Error,
		)
	}

	return nil
}

// @Internal
// Remove balance from coupon record.
// Doesn't check about active coupon expiration or other logic.
// Returns updated balance and an error object.
// Throws error object on insufficient balance.
func removeBalanceFromClaimedCouponUnchecked(
	balance int64,
	claimedCoupon *models.ClaimedCoupon,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if balance <= 0 {
		return utils.MakeError(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"invalid parameter",
			fmt.Errorf("provided balance: %d", balance),
		)
	}
	if claimedCoupon == nil {
		return utils.MakeError(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"invalid parameter",
			errors.New("claimedCoupon is nil pointer"),
		)
	}
	if claimedCoupon.CouponID == uuid.Nil ||
		claimedCoupon.ClaimedUserID == 0 {
		return utils.MakeError(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"invalid parameter",
			fmt.Errorf("provided claimedCoupon: %v", *claimedCoupon),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Remove balance and save.
	if result := session.Model(
		&claimedCoupon,
	).Clauses(
		clause.Returning{},
	).Update(
		"balance", gorm.Expr("balance - ?", balance),
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"failed remove balance properly",
			result.Error,
		)
	}
	if claimedCoupon.Balance < utils.GetMinimumChipInBalance()-balance {
		return utils.MakeErrorWithCode(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"failed remove balance properly",
			ErrCodeZeroBonusBalance,
			fmt.Errorf(
				"zero bonus balance. claimedCoupon: %v, balance: %d",
				*claimedCoupon, balance,
			),
		)
	} else if claimedCoupon.Balance < 0 {
		return utils.MakeErrorWithCode(
			"coupon_db",
			"removeBalanceFromClaimedCouponUnchecked",
			"failed remove balance properly",
			ErrCodeInsufficientBonusBalance,
			fmt.Errorf(
				"insufficient funds. claimedCoupon: %v, balance: %d",
				*claimedCoupon, balance,
			),
		)
	}

	return nil
}

// @Internal
// Add wager amount.
// Doesn't check about active coupon expiration or other logic.
func addWagerAmountUnchecked(
	wager int64,
	claimedCoupon *models.ClaimedCoupon,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if wager <= 0 {
		return utils.MakeError(
			"coupon_db",
			"addWagerAmountUnchecked",
			"invalid parameter",
			fmt.Errorf("provided wager: %d", wager),
		)
	}
	if claimedCoupon == nil {
		return utils.MakeError(
			"coupon_db",
			"addWagerAmountUnchecked",
			"invalid parameter",
			errors.New("claimedCoupon is nil pointer"),
		)
	}
	if claimedCoupon.CouponID == uuid.Nil ||
		claimedCoupon.ClaimedUserID == 0 {
		return utils.MakeError(
			"coupon_db",
			"addWagerAmountUnchecked",
			"invalid parameter",
			fmt.Errorf("provided claimedCoupon: %v", *claimedCoupon),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"addWagerAmountUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update wagered and save.
	if result := session.Model(
		&claimedCoupon,
	).Clauses(
		clause.Returning{},
	).Update(
		"wagered", gorm.Expr("wagered + ?", wager),
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"coupon_db",
			"addWagerAmountUnchecked",
			"failed to add wager amount properly",
			result.Error,
		)
	}

	return nil
}

// @Internal
// Leave transaction history.
// Returns generated transaction id and error object.
func leaveTransactionHistory(
	transactionHistory *models.CouponTransaction,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter.
	if transactionHistory == nil {
		return 0, utils.MakeError(
			"coupon_db",
			"leaveTransactionHistory",
			"invalid parameter",
			errors.New("provided transactionHistory is nil pointer"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_db",
			"leaveTransactionHistory",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Create coupon transaction histsory record.
	if result := session.Create(&transactionHistory); result.Error != nil {
		return 0, utils.MakeError(
			"coupon_db",
			"leaveTransactionHistory",
			"failed to create new coupon tx history record",
			result.Error,
		)
	}

	return transactionHistory.ID, nil
}

// @Internal
// Update exchanged amount.
// Doesn't check about active coupon expiration or other logic.
func updateExchangedAmountUnchecked(
	exchanged int64,
	claimedCoupon *models.ClaimedCoupon,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if exchanged <= 0 {
		return utils.MakeError(
			"coupon_db",
			"updateExchangedAmountUnchecked",
			"invalid parameter",
			fmt.Errorf("provided exchanged: %d", exchanged),
		)
	}
	if claimedCoupon == nil {
		return utils.MakeError(
			"coupon_db",
			"updateExchangedAmountUnchecked",
			"invalid parameter",
			errors.New("claimedCoupon is nil pointer"),
		)
	}
	if claimedCoupon.CouponID == uuid.Nil ||
		claimedCoupon.ClaimedUserID == 0 {
		return utils.MakeError(
			"coupon_db",
			"updateExchangedAmountUnchecked",
			"invalid parameter",
			fmt.Errorf("provided claimedCoupon: %v", *claimedCoupon),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"updateExchangedAmountUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update exchanged and save.
	if result := session.Model(
		&claimedCoupon,
	).Clauses(
		clause.Returning{},
	).Update(
		"exchanged", exchanged,
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"coupon_db",
			"updateExchangedAmountUnchecked",
			"failed to update exchanged",
			result.Error,
		)
	}

	return nil
}

// @Internal
// Create a claimed coupon record.
// Doesn't check about existance of another active claimed coupon or other logic.
func createClaimedCouponUnchecked(
	newClaimedCoupon *models.ClaimedCoupon,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if newClaimedCoupon == nil {
		return utils.MakeError(
			"coupon_db",
			"createClaimedCouponUnchecked",
			"invalid parameter",
			errors.New("provided newClaimedCoupon record is nil pointer"),
		)
	}
	if newClaimedCoupon.CouponID == uuid.Nil ||
		newClaimedCoupon.ClaimedUserID == 0 {
		return utils.MakeError(
			"coupon_db",
			"createClaimedCouponUnchecked",
			"invalid parameter",
			fmt.Errorf("provided new claimed coupon; %v", *newClaimedCoupon),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"createClaimedCouponUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Create a new claimed coupon record.
	if result := session.Create(&newClaimedCoupon); result.Error != nil {
		return utils.MakeError(
			"coupon_db",
			"createClaimedCouponUnchecked",
			"failed to create claimed coupon record",
			result.Error,
		)
	}

	return nil
}

// @Internal
// Update transaction's status to be `Succeed`.
// Returns error in case of not found `Pending` tx with provided ID.
func confirmTransactionStatus(
	txID uint,
	sessionId db_aggregator.UUID,
) (*models.CouponTransaction, error) {
	// 1. Validate parameter.
	if txID == 0 {
		return nil, utils.MakeError(
			"coupon_db",
			"confirmTransactionStatus",
			"invalid parameter",
			errors.New("provided txID is zero"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"confirmTransactionStatus",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Add confirm status and save.
	updatedTx := models.CouponTransaction{}
	if result := session.Model(
		&updatedTx,
	).Clauses(
		clause.Returning{},
	).Where(
		"id = ?",
		txID,
	).Where(
		"status = ?",
		models.CouponTransactionPending,
	).Update(
		"status", models.CouponTransactionSucceed,
	); result.Error != nil || result.RowsAffected != 1 {
		return nil, utils.MakeError(
			"coupon_db",
			"confirmTransactionStatus",
			"failed to confirm transaction",
			fmt.Errorf("txId: %d, error: %v", txID, result.Error),
		)
	}

	return &updatedTx, nil
}

// @Internal
// Update transaction's status to be `Failed`.
// And update transaction's afterRefund field.
// Returns error in case of not found `Pending` tx with provided ID.
func declineTransactionStatus(
	txID uint,
	afterRefund int64,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if txID == 0 {
		return utils.MakeError(
			"coupon_db",
			"declineTransactionStatus",
			"invalid parameter",
			errors.New("provided txID is zero"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"declineTransactionStatus",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Decline status and save.
	if result := session.Model(
		&models.CouponTransaction{},
	).Where(
		"id = ?",
		txID,
	).Where(
		"status = ?",
		models.CouponTransactionPending,
	).Updates(
		map[string]interface{}{
			"status":       models.CouponTransactionFailed,
			"after_refund": afterRefund,
		},
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"coupon_db",
			"declineTransactionStatus",
			"failed to decline transaction",
			fmt.Errorf(
				"txId: %d, afterRefund: %d, error: %v, rowsAffected: %v",
				txID, afterRefund, result.Error, result.RowsAffected,
			),
		)
	}

	return nil
}

// @External
// Get currently active claimed coupon within main session.
func GetActiveUserCoupon(userID uint) *ActiveUserCouponMeta {
	// 2. Retrieve session
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	// 3. Lock and retrieve claimedCoupon record.
	activeCoupon := models.ClaimedCoupon{}
	if err := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"Coupon",
	).Where(
		"claimed_user_id = ?",
		userID,
	).Where(
		"created_at > ?",
		getActiveCouponLifeTimeBeforeNow(),
	).Order(
		"created_at",
	).First(&activeCoupon).Error; err != nil {
		if !errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			log.LogMessage(
				"coupon_db_GetActiveUserCoupon",
				"failed to get active user coupon",
				"error",
				logrus.Fields{
					"userID": userID,
					"error":  err.Error(),
				},
			)
		}
		return nil
	}

	return &ActiveUserCouponMeta{
		Code:       *activeCoupon.Coupon.Code,
		Balance:    activeCoupon.Balance,
		Claimed:    activeCoupon.Coupon.BonusBalance,
		Wagered:    activeCoupon.Wagered,
		WagerLimit: activeCoupon.Coupon.BonusBalance * int64(config.COUPON_REQUIRED_WAGER_TIMES),
		RemainingTime: time.Until(time.Unix(
			activeCoupon.CreatedAt.Unix()+int64(time.Hour.Seconds())*int64(config.COUPON_BALANCE_LIFE_TIME_IN_HOURS), 0,
		)).Milliseconds(),
	}
}

// @Internal
// A helper function to get timestamp `COUPON_BALANCE_LIFE_TIME_IN_HOUR` hours ago then now.
func getActiveCouponLifeTimeBeforeNow() time.Time {
	return time.Unix(time.Now().Unix()-int64(time.Hour.Seconds())*int64(config.COUPON_BALANCE_LIFE_TIME_IN_HOURS), 0)
}

/**
* @Internal
* Retrieves coupon shortcut record from shortcut.
* Returns nil in case of error.
* Out of this function, assumes that the issue is not found one.
 */
func retrieveCouponShortcut(
	shortcut string,
) (*models.CouponShortcut, error) {
	// 1. Validate parameter.
	if shortcut == "" {
		return nil, utils.MakeError(
			"coupon_db",
			"retrieveCouponShortcut",
			"invalid parameter",
			errors.New("provided shortcut is empty string"),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"retrieveCouponShortcut",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Retrieve coupon shortcut.
	couponShortcut := models.CouponShortcut{}
	if result := session.Where(
		"shortcut = ?",
		shortcut,
	).Last(
		&couponShortcut,
	); result.Error != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"retrieveCouponShortcut",
			"failed to retrieve coupon shortcut",
			fmt.Errorf(
				"shortcut: %s, err: %v",
				shortcut, result.Error,
			),
		)
	}

	return &couponShortcut, nil
}

/**
* @Internal
* Create coupon shortcut record.
* Doesn't check about couponID existence.
* Tries to catch duplicated error with the model constraints.
* Shortcut and CouponID should be both unique.
 */
func createCouponShortcutUnchecked(
	couponID uuid.UUID,
	shortcut string,
) error {
	// 1. Validate parameter.
	if shortcut == "" ||
		couponID == uuid.Nil {
		return utils.MakeError(
			"coupon_db",
			"createCouponShortcutUnchecked",
			"invalid parameter",
			fmt.Errorf(
				"couponID: %v, shortcut: %s",
				couponID, shortcut,
			),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"createCouponShortcutUnchecked",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Create coupon shortcut.
	if result := session.Create(
		&models.CouponShortcut{
			CouponID: couponID,
			Shortcut: shortcut,
		},
	); result.Error != nil {
		return utils.MakeError(
			"coupon_db",
			"createCouponShortcutUnchecked",
			"failed to create coupon shortcut",
			fmt.Errorf(
				"couponID: %v, shortcut: %s, err: %v",
				couponID,
				shortcut,
				result.Error,
			),
		)
	}

	return nil
}

/**
* @Internal
* Delete coupon shortcut record.
 */
func deleteCouponShortcutUnchecked(
	shortcut string,
) error {
	// 1. Validate parameter.
	if shortcut == "" {
		return utils.MakeError(
			"coupon_db",
			"deleteCouponShortcutUnchecked",
			"invalid parameter",
			errors.New("provided shortcut is empty string"),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"deleteCouponShortcutUnchecked",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Delete coupon shortcut.
	if result := session.Where(
		"shortcut = ?",
		shortcut,
	).Delete(
		&models.CouponShortcut{},
	); result.Error != nil {
		return utils.MakeError(
			"coupon_db",
			"deleteCouponShortcutUnchecked",
			"failed to delete coupon shortcut",
			fmt.Errorf(
				"shortcut: %s, err: %v",
				shortcut, result.Error,
			),
		)
	}

	return nil
}

/**
* @Internal
* Checks existence of dreamtower bets with coupon.
 */
func existingDreamtowerRoundWithCoupon(
	userID uint,
) bool {
	// 1. Validate parameters.
	if userID == 0 {
		return false
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return false
	}

	// 3. Count currently playing rounds with coupon balance.
	var playingRoundsWithCoupon int64
	if result := session.Model(
		&models.DreamTowerRound{},
	).Where(
		"user_id = ? AND status = ? AND paid_balance_type = ?",
		userID,
		models.DreamTowerPlaying,
		models.CouponBalanceForGame,
	).Count(&playingRoundsWithCoupon); result.Error != nil {
		log.LogMessage(
			"existingDreamtowerRoundWithCoupon",
			"failed to count playing rounds",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
		return false
	}

	// 4. Return result.
	if playingRoundsWithCoupon > 0 {
		return true
	}
	return false
}

/**
* @Internal
* Checks existence of crash bets with coupon.
 */
func existingCrashRoundWithCoupon(
	userID uint,
) bool {
	// 1. Validate parameters.
	if userID == 0 {
		return false
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return false
	}

	// 3. Get currently playing crash round.
	var playingRound models.CrashRound
	if result := session.Where(
		"bet_started_at IS NOT NULL AND ended_at IS NULL",
	).Last(&playingRound); result.Error != nil {
		log.LogMessage(
			"existingCrashRoundWithCoupon",
			"failed to get currently playing round",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
		return false
	}

	// 4. Count user's bets on playing crash round.
	var betCount int64
	if result := session.Model(
		&models.CrashBet{},
	).Where(
		"user_id = ? AND round_id = ? AND paid_balance_type = ? AND profit IS NULL AND payout_multiplier IS NULL",
		userID,
		playingRound.ID,
		models.CouponBalanceForGame,
	).Count(&betCount); result.Error != nil {
		log.LogMessage(
			"existingCrashRoundWithCoupon",
			"failed to count bets on currently running round.",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
		return false
	}

	// 5. Check count and return result.
	if betCount > 0 {
		return true
	}
	return false
}
