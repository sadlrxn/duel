package redis

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

// Set of user IDs which are cached to be not performing
// first deposit bonus system.
// User IDs are added to this set in case of first deposit bonus is done,
// or they activated other affiliate code than first deposit bonus one.
const FIRST_DEPOSIT_AFTER_AFFILIATE_ACTIVATION_SET_KEY_NAME = "set-first-deposit-after-affiliate-activation"

/**
* Add userIDs to set.
 */
func AddFirstDepositBonusDone(userIDs []uint) error {
	// 1. Validate parameter.
	if userIDs == nil || len(userIDs) == 0 {
		return utils.MakeError(
			"redis_first_deposit_after_affiliate",
			"AddFirstDepositBonusDone",
			"invalid parameter",
			errors.New(
				fmt.Sprintf("UserIDs: %v", userIDs),
			),
		)
	}

	// 2. Add userIDs to set.
	if _, err := rdb.SAdd(
		redis_ctx,
		FIRST_DEPOSIT_AFTER_AFFILIATE_ACTIVATION_SET_KEY_NAME,
		userIDs,
	).Result(); err != nil {
		return utils.MakeError(
			"redis_first_deposit_after_affiliate",
			"AddFirstDepositBonusDone",
			"failed to add userIDs to set",
			err,
		)
	}
	return nil
}

/**
* Remove userIDs from set.
 */
func RemoveFirstDepositBonusDone(userIDs []uint) error {
	// 1. Validate parameter.
	if userIDs == nil || len(userIDs) == 0 {
		return utils.MakeError(
			"redis_first_deposit_after_affiliate",
			"RemoveFirstDepositBonusDone",
			"invalid parameter",
			errors.New(
				fmt.Sprintf("UserIDs: %v", userIDs),
			),
		)
	}

	// 2. Remove userIDs from set.
	if _, err := rdb.SRem(
		redis_ctx,
		FIRST_DEPOSIT_AFTER_AFFILIATE_ACTIVATION_SET_KEY_NAME,
		userIDs,
	).Result(); err != nil {
		return utils.MakeError(
			"redis_first_deposit_after_affiliate",
			"RemoveFirstDepositBonusDone",
			"failed to remove userIDs from set",
			err,
		)
	}
	return nil
}

/**
* Check whether userID is in the set.
 */
func IsFirstDepositBonusDone(userID uint) bool {
	// 1. Validate parameter.
	if userID == 0 {
		return false
	}

	// 2. Check wheter userID exists in the set.
	if exists, err := rdb.SIsMember(
		redis_ctx,
		FIRST_DEPOSIT_AFTER_AFFILIATE_ACTIVATION_SET_KEY_NAME,
		userID,
	).Result(); err == nil {
		return exists
	} else if err != nil {
		log.LogMessage(
			"IsFirstDepositBonusDone",
			"failed to run SISMEMBER command",
			"error",
			logrus.Fields{
				"userID": userID,
				"error":  err.Error(),
			},
		)
		return false
	}
	return false
}
