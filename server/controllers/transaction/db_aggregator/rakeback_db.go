package db_aggregator

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func retrieveRakebackInfoAndRate(user User, sessionId ...UUID) (*models.Rakeback, uint, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, 0, utils.MakeError(
			"rakeback_db",
			"retrieveRakebackInfoAndRate",
			"invalid parameter",
			errors.New("provided to user argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, 0, utils.MakeError(
			"rakeback_db",
			"retrieveRakebackInfoAndRate",
			"failed to find session",
			err,
		)
	}

	// 3. Retrieve rakeback info.
	rakebackInfo := models.Rakeback{}
	if result := session.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where(
		"user_id = ?",
		user,
	).First(&rakebackInfo); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, 0, utils.MakeError(
			"rakeback_db",
			"retrieveRakebackInfoAndRate",
			"failed to retrieve rakeback info",
			result.Error,
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newRakebackInfo, err := generateRakebackInfo(user, sessionId...)
		if err != nil || newRakebackInfo == nil {
			return nil, 0, utils.MakeError(
				"rakeback_db",
				"retrieveRakebackInfoAndRate",
				"failed to generate new rakeback on not found",
				err,
			)
		}

		rakebackInfo = *newRakebackInfo
	}
	serverConfig := config.GetServerConfig()
	rate := serverConfig.BaseRakeBackRate + serverConfig.AdditionalRakeBackRate
	if time.Now().Before(
		rakebackInfo.AdditionalRakebackExpired,
	) {
		rate2 := serverConfig.BaseRakeBackRate + rakebackInfo.AdditionalRakebackRate
		if rate < rate2 {
			rate = rate2
		}
	}
	if rate > config.RAKEBACK_MAX {
		rate = config.RAKEBACK_MAX
	}

	// 3. Check whether rate is too high
	if rate > 20 {
		return &rakebackInfo, rate, utils.MakeError(
			"rakeback_db",
			"retrieveRakebackInfoAndRate",
			"too high rakeback rate",
			fmt.Errorf("rate: %d, limit: 20", rate),
		)
	}

	return &rakebackInfo, rate, nil
}

// @External
// Distribute rakeback rewards to user.
func distributeRakeback(to User, feeAmount int64, sessionId ...UUID) (int64, error) {
	if to == 0 {
		return 0, utils.MakeError(
			"rakeback_db",
			"distributeRakeback",
			"invalid parameter",
			errors.New("provided to user argument is invalid"),
		)
	}

	if feeAmount == 0 {
		return 0, nil
	}

	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"distributeRakeback",
			"failed to find session",
			err,
		)
	}

	rakebackInfo, rate, err := retrieveRakebackInfoAndRate(to, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"distributeRakeback",
			"failed to retrieve rakeback info and rate",
			err,
		)
	}
	distributed := feeAmount * int64(rate) / 100

	rakebackInfo.TotalEarned += distributed
	rakebackInfo.Reward += distributed
	if result := session.Save(&rakebackInfo); result.Error != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"distributeRakeback",
			"failed to update rakeback record",
			result.Error,
		)
	}

	return distributed, nil
}

// @External
// Claim user's rakeback and leave transaction record.
func claimRakeback(user User, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if user == 0 {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"invalid parameter",
			errors.New("provided user is invalid"),
		)
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to get session",
			err,
		)
	}

	// 3. Retrieve rake back record.
	rakebackInfo := models.Rakeback{}
	if result := session.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where(
		"user_id = ?",
		user,
	).First(&rakebackInfo); result.Error != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to retrieve rakeback info",
			result.Error,
		)
	}

	// 4. Reset rake back rewards.
	rakebackAmount := rakebackInfo.Reward
	balanceLoad := BalanceLoad{
		ChipBalance: &rakebackAmount,
	}
	rakebackInfo.Reward = 0
	if result := session.Save(&rakebackInfo); result.Error != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to update rakeback reward amount to 0",
			result.Error,
		)
	}

	// 5. Transfer rewards to user.
	/*transferResult*/
	transferResult, err := transfer(
		nil,
		&user,
		&balanceLoad,
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to transfer rakeback rewards to user",
			err,
		)
	}

	// ================== Archived for transaction ==================
	// 6. Record transaction as confirmed one.
	toWallet, err := GetUserWallet(&user, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to retrieve user's wallet",
			err,
		)
	}

	transaction, err := RecordTransaction(
		&TransactionLoad{
			FromWallet: nil,
			ToWallet:   toWallet,
			Balance:    balanceLoad,
			Type:       models.TxClaimRakebackReward,
		},
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to record transaction",
			err,
		)
	}

	if err := ConfirmTransaction(
		&TransactionLoad{
			ToWalletPrevID: transferResult.ToPrevBalance,
			ToWalletNextID: transferResult.ToNextBalance,
			OwnerID:        uint(*toWallet),
			OwnerType:      models.TransactionWalletReferenced,
		},
		transaction,
		sessionId...,
	); err != nil {
		return 0, utils.MakeError(
			"rakeback_db",
			"claimRakeback",
			"failed to confirm transaction",
			err,
		)
	}
	// ================== Archived for transaction ==================

	return *balanceLoad.ChipBalance, nil
}

// @External
// Retrieve rakeback total earning, and claimable reward.
func getRakebackInfo(user User, sessionId ...UUID) (models.Rakeback, error) {
	nothing := models.Rakeback{}

	// 1. Validate parameter.
	if user == 0 {
		return nothing, utils.MakeError(
			"rakeback_db",
			"getRakebackInfo",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return nothing, utils.MakeError(
			"rakeback_db",
			"getRakebackInfo",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Retrieve rakeback rewards.
	rakebackInfo := models.Rakeback{}
	if result := session.Where(
		"user_id = ?",
		user,
	).First(&rakebackInfo); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nothing, utils.MakeError(
			"rakeback_db",
			"getRakebackInfo",
			"failed to retrieve rakeback info",
			result.Error,
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newRakebackInfo, err := generateRakebackInfo(user, sessionId...)
		if err != nil {
			return nothing, utils.MakeError(
				"rakeback_db",
				"getRakebackInfo",
				"failed to generate new rakeback on not found",
				err,
			)
		}

		if newRakebackInfo != nil {
			return *newRakebackInfo, nil
		}
	}

	return rakebackInfo, nil
}

// @External
// Create a new rakeback record for a user if not exists.
func generateRakebackInfo(user User, sessionId ...UUID) (*models.Rakeback, error) {
	if user == 0 {
		return nil, utils.MakeError(
			"rakeback_db",
			"generateRakebackInfo",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"rakeback_db",
			"generateRakebackInfo",
			"failed to retrieve main session",
			err,
		)
	}

	rakebackInfo := models.Rakeback{}
	if result := session.Where(
		"user_id = ?",
		user,
	).First(&rakebackInfo); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utils.MakeError(
			"rakeback_db",
			"generateRakebackInfo",
			"failed to retrieve rakeback info",
			result.Error,
		)
	} else if result.Error == nil {
		return nil, nil
	}

	newRakeback := models.Rakeback{
		UserID: uint(user),
	}
	if result := session.Create(&newRakeback); result.Error != nil {
		return nil, utils.MakeError(
			"rakeback_db",
			"generateRakebackInfo",
			"failed to create rakeback info",
			result.Error,
		)
	}

	return &newRakeback, nil
}

// @Internal
// Update additional rakeback.
func updateUserAdditionalRakebackUnchecked(
	user User,
	duration time.Duration,
	additionalRate uint,
	sessionId ...UUID,
) (bool, error) {
	// 1. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return false, utils.MakeError(
			"rakeback_db",
			"updateAdditionalRakebackUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 2. Retrieve rakeback info.
	rakebackInfo, err := getRakebackInfo(user, sessionId...)
	if err != nil {
		return false, utils.MakeError(
			"rakeback_db",
			"updateAdditionalRakebackUnchecked",
			"failed to get rakeback info",
			err,
		)
	}
	if rakebackInfo.ActivatedAffiliateOnce {
		return false, nil
	}
	newRakebackTime := time.Now().Add(duration)
	if newRakebackTime.Before(rakebackInfo.AdditionalRakebackExpired) {
		return false, nil
	}

	// 3. Update rakeback info.
	rakebackInfo.AdditionalRakebackRate = additionalRate
	rakebackInfo.AdditionalRakebackExpired = newRakebackTime
	rakebackInfo.ActivatedAffiliateOnce = true

	if result := session.Save(&rakebackInfo); result.Error != nil {
		return false, utils.MakeError(
			"rakeback_db",
			"updateAdditionalRakebackUnchecked",
			"failed to update rakeback info",
			result.Error,
		)
	}

	return true, nil
}

// @External
// Set activate affiliate once flag.
func setActivateAffiliateOnceForRakeback(user User, sessionId ...UUID) (bool, error) {
	// 1. Validate parameter.
	if user == 0 {
		return false, utils.MakeError(
			"rakeback_db",
			"setActivateAffiliateOnce",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Set activate affiliate once.
	activated, err := updateUserAdditionalRakebackUnchecked(
		user,
		time.Hour*24,
		5,
		sessionId...,
	)
	if err != nil {
		return false, utils.MakeError(
			"rakeback_db",
			"setActivateAffiliateOnce",
			"failed to activate affiliate",
			err,
		)
	}

	return activated, nil
}
