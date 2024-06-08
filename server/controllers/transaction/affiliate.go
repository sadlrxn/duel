package transaction

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

func createAffiliateCode(user db_aggregator.User, codes []string) error {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"transaction",
			"createAffiliateCode",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	if err := db_aggregator.CreateAffiliateCode(
		user,
		codes,
		sessionId,
	); err != nil {
		return utils.MakeError(
			"transaction",
			"createAffiliateCode",
			"failed to create affiliate codes",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"transaction",
			"createAffiliateCode",
			"failed to commit session",
			err,
		)
	}

	return nil
}

func deleteAffiliateCode(user db_aggregator.User, codes []string) (int64, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"deleteAffiliateCode",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	claimed, err := db_aggregator.DeleteAffiliateCode(
		user,
		codes,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"deleteAffiliateCode",
			"failed to delete affiliate codes",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"transaction",
			"deleteAffiliateCode",
			"failed to commit session",
			err,
		)
	}

	return claimed, nil
}

func claimAffiliateRewards(user db_aggregator.User, codes []string) (int64, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimAffiliateRewards",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	claimed, err := db_aggregator.ClaimAffiliateRewards(
		user,
		codes,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimAffiliateRewards",
			"failed to claim affiliate rewards",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimAffiliateRewards",
			"failed to commit session",
			err,
		)
	}

	return claimed, nil
}

func activateAffiliateCode(user db_aggregator.User, code string) (bool, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return false, utils.MakeError(
			"transaction",
			"activateAffiliateCode",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	isFirst, err := db_aggregator.ActivateAffiliateCode(
		user,
		code,
		sessionId,
	)
	if err != nil {
		return false, utils.MakeError(
			"transaction",
			"activateAffiliateCode",
			"failed to activate affiliate code",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return false, utils.MakeError(
			"transaction",
			"activateAffiliateCode",
			"failed to commit session",
			err,
		)
	}

	return isFirst, nil
}

func deactivateAffiliateCode(user db_aggregator.User, code string) error {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"transaction",
			"deactivateAffiliateCode",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	if err := db_aggregator.DeactivateAffiliateCode(
		user,
		code,
		sessionId,
	); err != nil {
		return utils.MakeError(
			"transaction",
			"deactivateAffiliateCode",
			"failed to activate affiliate code",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"transaction",
			"deactivateAffiliateCode",
			"failed to commit session",
			err,
		)
	}

	return nil
}

/**
* @External
* Try apply for first deposit bonus for every chip deposit.
 */
func TryApplyForFirstDepositBonus(
	userID uint,
	depositAmount int64,
) (int64, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		depositAmount <= 0 {
		return 0, nil
	}

	// 2. Checks whether the userID is cached to not perform
	// first deposit bonus from redis.
	if redis.IsFirstDepositBonusDone(userID) {
		return 0, nil
	}

	// 3. Retrieve active affiliate coupon.
	activeAffiliate, err := db_aggregator.GetActiveAffiliateCode(
		db_aggregator.User(userID),
	)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate",
			"TryApplyForFirstDepositBonus",
			"failed to retrieve active affiliate code",
			fmt.Errorf(
				"userID: %d, err: %v",
				userID, err,
			),
		)
	}
	if activeAffiliate == nil {
		return 0, nil
	}

	// 4. Checks whether should perform first deposit bonus,
	// if not, adds to cache to not perform for next time as well.
	// Adds to cache on successful perform as well.
	if shouldPerformFirstDepositBonus(activeAffiliate) {
		if bonus, err := coupon.PerformFirstDepositBonus(
			userID, depositAmount,
		); err != nil {
			return 0, utils.MakeError(
				"affiliate",
				"TryApplyForFirstDepositBonus",
				"failed to perform first deposit bonus",
				fmt.Errorf(
					"userID: %d, depositAmount: %d, err: %v",
					userID, depositAmount, err,
				),
			)
		} else {
			if err := db_aggregator.SetActiveAffiliateFirstDepositDone(
				activeAffiliate.ID,
			); err != nil {
				log.LogMessage(
					"affiliate_TryApplyForFirstDepositBonus",
					"failed to set active affiliate first deposit done",
					"error",
					logrus.Fields{
						"userID":        userID,
						"affiliateCode": activeAffiliate.Code,
						"error":         err.Error(),
					},
				)
			}
			redis.AddFirstDepositBonusDone([]uint{userID})
			return bonus, nil
		}
	} else {
		redis.AddFirstDepositBonusDone([]uint{userID})
		return 0, nil
	}
}

/**
* @Internal
* Checks whether should perform first deposit bonus from the
* retrieved active affiliate meta.
 */
func shouldPerformFirstDepositBonus(
	activeAffiliate *db_aggregator.ActiveAffiliateMeta,
) bool {
	return activeAffiliate != nil &&
		activeAffiliate.IsFirstDepositBonus &&
		!activeAffiliate.FirstDepositDone
}

/**
* @External
* Updates affiliate's isFirstDepositBonus property.
* If previous type was not supporting first deposit bonus, and setting
* as another one which supports, should remove the activated userIDs from
* cache not to perform deposit bonus.
 */
func UpdateAffiliateFirstDepositBonus(
	affiliateCode string,
	isFirstDepositBonus bool,
) error {
	// 1. Validate parameter.
	if len(affiliateCode) == 0 {
		return utils.MakeError(
			"affiliate",
			"UpdateAffiliateFirstDepositBonus",
			"invalid parameter",
			errors.New("provided affiliate code is empty string"),
		)
	}

	// 2. Update affiliate type.
	previous, err := db_aggregator.UpdateAffiliateFirstDepositBonus(
		affiliateCode,
		isFirstDepositBonus,
	)
	if err != nil {
		return utils.MakeError(
			"affiliate",
			"UpdateAffiliateFirstDepositBonus",
			"failed to update affiliate type",
			fmt.Errorf(
				"affiliateCode: %s, isFirstDepositBonus: %v, err: %v",
				affiliateCode, isFirstDepositBonus, err,
			),
		)
	}

	// 3. If previous type is not supporting first deposit bonus,
	// while the new one is, remove the using IDs from cache.
	if isFirstDepositBonus && !previous {
		detail, err := db_aggregator.GetAffiliateDetail(
			affiliateCode,
		)
		if err != nil {
			log.LogMessage(
				"affiliate_UpdateAffiliateFirstDepositBonus",
				"failed to get affiliate detail",
				"error",
				logrus.Fields{
					"affiliateCode": affiliateCode,
					"error":         err.Error(),
				},
			)
		} else if detail != nil {
			userIDs := []uint{}
			for _, user := range detail.Users {
				userIDs = append(
					userIDs,
					user.ID,
				)
			}
			// redis.RemoveFirstDepositBonusDone(userIDs)
		}
	}

	return nil
}
