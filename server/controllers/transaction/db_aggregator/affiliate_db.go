package db_aggregator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @External
// Creates affiliate codes owned by user.
// If provided duplicated code, returns error.
func createAffiliateCode(user User, codes []string, sessionId ...UUID) error {
	// 1. Validate parameter.
	if user == 0 {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	if len(codes) == 0 {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"invalid parameter",
			errors.New("no code provided"),
		)
	}
	lowerCodes := []string{}
	for _, code := range codes {
		lowerCode := strings.ToLower(code)
		if strings.Contains(lowerCode, " ") {
			return utils.MakeError(
				"affiliate_db",
				"createAffiliateCode",
				"invalid parameter",
				fmt.Errorf(
					"code(%s) is containing space",
					code,
				),
			)
		}

		for _, reserved := range config.AFFILIATE_RESERVED_WORDS {
			if strings.Contains(lowerCode, strings.ToLower(reserved)) {
				return utils.MakeError(
					"affiliate_db",
					"createAffiliateCode",
					"invalid parameter",
					fmt.Errorf(
						"code(%s) is containing reserved word(%s)",
						code,
						reserved,
					),
				)
			}
		}
		lowerCodes = append(lowerCodes, lowerCode)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve user info.
	userInfo := models.User{}
	if result := session.Preload("Statistics").First(&userInfo, user); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"failed to retrieve user info",
			result.Error,
		)
	}
	if userInfo.Statistics.TotalWagered < int64(config.AFFILIATE_WAGER_LIMIT_FOR_CREATION) {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"not enough wager amount",
			fmt.Errorf(
				"limit: %d, wagered: %d",
				config.AFFILIATE_WAGER_LIMIT_FOR_CREATION,
				userInfo.Statistics.TotalWagered,
			),
		)
	}

	// 4. Lock affiliates by creator id.
	userAffiliates := []models.Affiliate{}
	if result := session.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where(
		"creator_id = ?",
		user,
	).Find(&userAffiliates); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"failed to lock user affiliates",
			result.Error,
		)
	}
	if len(userAffiliates) >= 5 {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"too many affiliate codes",
			fmt.Errorf("limit: %d, actual: %d", 5, len(userAffiliates)),
		)
	}

	// 5. Check for code duplication.
	duplicatedAffiliates := []models.Affiliate{}
	if result := session.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Where(
		"lower(code) in ?",
		lowerCodes,
	).Find(&duplicatedAffiliates); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"failed to retrieve duplicated code count",
			result.Error,
		)
	}
	if len(duplicatedAffiliates) > 0 {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"containing duplicated code",
			fmt.Errorf("duplicating count is: %v", len(duplicatedAffiliates)),
		)
	}

	// 6. Create affiliate codes.
	newAffiliates := []models.Affiliate{}
	for _, code := range codes {
		newAffiliates = append(
			newAffiliates,
			models.Affiliate{
				Code:      code,
				CreatorID: uint(user),
			},
		)
	}
	if result := session.Create(&newAffiliates); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"createAffiliateCode",
			"failed to create affiliates",
			result.Error,
		)
	}

	return nil
}

// @External
// Deletes affiliate codes from owner.
// Auto claim rewards on delete.
// Returns claimed rewards.
func deleteAffiliateCode(user User, codes []string, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if user == 0 {
		return 0, utils.MakeError(
			"affiliate_db",
			"deleteAffiliateCode",
			"invalid parameter",
			errors.New("provided user is invalid"),
		)
	}
	if len(codes) == 0 {
		return 0, nil
	}

	// 2. Claim affiliate codes.
	claimed, err := claimAffiliateRewards(user, codes, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"deleteAffiliateCode",
			"failed to claim affiliate rewards",
			err,
		)
	}

	// 3. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"deleteAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 4. Delete affiliates.
	if result := session.Where(
		"creator_id = ?",
		user,
	).Where(
		"code in ?",
		codes,
	).Delete(&models.Affiliate{}); result.Error != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"deleteAffiliateCode",
			"failed to delete affiliates",
			result.Error,
		)
	}

	return claimed, nil
}

// @External
// Claims rewards accumulating on affiliate codes.
// Returns claimed rewards.
func claimAffiliateRewards(user User, codes []string, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if user == 0 {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"invalid parameter",
			errors.New("provided user is invalid"),
		)
	}
	if len(codes) == 0 {
		return 0, nil
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to get session",
			err,
		)
	}

	// 3. Retrieve affiliates.
	affiliates := []models.Affiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"creator_id = ?",
		user,
	).Where(
		"code in ?",
		codes,
	).Find(&affiliates); result.Error != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to retrieve affiliates",
			result.Error,
		)
	}
	if len(affiliates) != len(codes) {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"mismatching retrieved affiliate count",
			fmt.Errorf(
				"affiliates: %d, codes: %d",
				len(affiliates),
				len(codes),
			),
		)
	}

	// 4. Reset affiliate rewards.
	claimed := int64(0)
	for i, affiliate := range affiliates {
		claimed += affiliate.Reward
		affiliates[i].Reward = 0
	}
	if claimed == 0 {
		return 0, nil
	}
	if result := session.Save(&affiliates); result.Error != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to reset affiliate rewards",
			result.Error,
		)
	}

	// 5. Transfer claimed.
	balanceLoad := BalanceLoad{
		ChipBalance: &claimed,
	}
	transferResult, err := transfer(
		nil,
		&user,
		&balanceLoad,
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to transfer affiliate rewards to user",
			err,
		)
	}

	// 6. Record transaction as confirmed one.
	toWallet, err := GetUserWallet(&user, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to retrieve user's wallet",
			err,
		)
	}

	transaction, err := RecordTransaction(
		&TransactionLoad{
			FromWallet: nil,
			ToWallet:   toWallet,
			Balance:    balanceLoad,
			Type:       models.TxClaimAffiliateReward,
		},
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"claimAffiliateRewards",
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
			"affiliate_db",
			"claimAffiliateRewards",
			"failed to confirm transaction",
			err,
		)
	}

	return claimed, nil
}

// @External
// Activate a new affiliate code.
// Return a boolean whether this is the first activation or not.
func activateAffiliateCode(user User, code string, sessionId ...UUID) (bool, error) {
	// 1. Validate parameter.
	if user == 0 {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}
	if code == "" {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"invalid parameter",
			errors.New("provided code is empty string"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve affiliate.
	affiliate := models.Affiliate{}
	if result := session.Where(
		"lower(code) = ?",
		strings.ToLower(code),
	).First(&affiliate); result.Error != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"failed to retrieve affiliate",
			result.Error,
		)
	}
	if affiliate.CreatorID == uint(user) {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"cannot activate your own code",
			fmt.Errorf("code: %v, user: %d", code, user),
		)
	}

	// 4. Retrieve userInfo and check whether the creation time
	// is in allowed period to activate affiliate code.
	userExistence := int64(0)
	if err := session.Model(
		&models.User{},
	).Where(
		"created_at > ?",
		time.Now().Add(
			-time.Hour*time.Duration(config.AFFILIATE_ACTIVATION_TIMELINE_IN_HOURS),
		),
	).Where(
		"id = ?",
		user,
	).Count(&userExistence).Error; err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"activeAffiliateCode",
			"failed to retrieve user info",
			fmt.Errorf(
				"userID: %d, err: %v",
				user, err,
			),
		)
	} else if userExistence != 1 {
		return false, utils.MakeError(
			"affiliate_db",
			"activeAffiliateCode",
			"failed to retrieve user info",
			fmt.Errorf(
				"userID: %d, possibly expired affiliate code activation timeline",
				user,
			),
		)
	}

	// 5. Activate affiliate.
	activeAffiliate := models.ActiveAffiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"user_id = ?",
		user,
	).First(&activeAffiliate); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"failed to retrieve active affiliate",
			result.Error,
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		activeAffiliate.UserID = uint(user)
		activeAffiliate.AffiliateID = affiliate.ID
		if result := session.Create(&activeAffiliate); result.Error != nil {
			return false, utils.MakeError(
				"affiliate_db",
				"activateAffiliateCode",
				"failed to create active affiliate record on not found",
				result.Error,
			)
		}
	}
	if activeAffiliate.AffiliateID != affiliate.ID {
		activeAffiliate.AffiliateID = affiliate.ID
		if result := session.Save(&activeAffiliate); result.Error != nil {
			return false, utils.MakeError(
				"affiliate_db",
				"activateAffiliateCode",
				"failed to update active affiliate id",
				result.Error,
			)
		}
	}

	// 6. Activate additional rakeback conditionally.
	activated, err := setActivateAffiliateOnceForRakeback(user, sessionId...)
	if err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"activateAffiliateCode",
			"failed to active additional rakeback",
			err,
		)
	}

	// 7. Update affiliate life time. On failure doesn't return error,
	// but prints the error, since it is not critical issue.
	if err := activateAffiliateLifetimeUnchecked(
		activeAffiliate.UserID,
		activeAffiliate.AffiliateID,
		sessionId...,
	); err != nil {
		log.LogMessage(
			"affiliate_db_activateAffiliateCode",
			"failed to update affiliate lifetime",
			"error",
			logrus.Fields{
				"error":       err.Error(),
				"userID":      activeAffiliate.UserID,
				"affiliateID": activeAffiliate.AffiliateID,
			},
		)
	}

	return activated, nil
}

// @External
// Get activated affiliate code.
// If user doesn't have activated any affiliate code, return nil.
func getActiveAffiliateCode(user User, sessionId ...UUID) (*ActiveAffiliateMeta, error) {
	// 1. Validate parameters.
	if user == 0 {
		return nil, utils.MakeError(
			"affiliate_db",
			"getActiveAffiliateCode",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"affiliate_db",
			"getActiveAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve active affiliate code.
	activeAffiliate := models.ActiveAffiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"Affiliate.Creator",
	).Where(
		"user_id = ?",
		user,
	).First(&activeAffiliate); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utils.MakeError(
			"affiliate_db",
			"getActiveAffiliateCode",
			"failed to retrieve active affiliate",
			result.Error,
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &ActiveAffiliateMeta{
		ID:                  activeAffiliate.ID,
		Code:                activeAffiliate.Affiliate.Code,
		Rate:                activeAffiliate.Affiliate.CustomAffiliateRate,
		OwnerID:             activeAffiliate.Affiliate.Creator.ID,
		OwnerName:           activeAffiliate.Affiliate.Creator.Name,
		OwnerAvatar:         activeAffiliate.Affiliate.Creator.Avatar,
		IsFirstDepositBonus: activeAffiliate.Affiliate.IsFirstDepositBonus,
		FirstDepositDone:    activeAffiliate.FirstDepositDone,
	}, nil
}

// @External
// Get owned affilaite code meta.
func getOwnedAffiliateCode(user User, sessionId ...UUID) ([]AffiliateMeta, error) {
	// 1. Validate parameters.
	if user == 0 {
		return nil, utils.MakeError(
			"affiliate_db",
			"getOwnedAffiliateCode",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"affiliate_db",
			"getOwnedAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve owned affiliates.
	affiliates := []models.Affiliate{}
	if result := session.Preload(
		"ActiveAffiliates",
	).Where(
		"creator_id = ?",
		user,
	).Find(&affiliates); result.Error != nil {
		return nil, utils.MakeError(
			"affiliate_db",
			"getOwnedAffiliateCode",
			"failed to retrieve affiliates",
			result.Error,
		)
	}

	// 4. Build and return affiliates meta.
	affiliateMetas := []AffiliateMeta{}
	if len(affiliates) == 0 {
		return affiliateMetas, nil
	}
	for _, affiliate := range affiliates {
		affiliateMetas = append(
			affiliateMetas,
			AffiliateMeta{
				Code:         affiliate.Code,
				UserCnt:      uint(len(affiliate.ActiveAffiliates)),
				TotalEarned:  affiliate.TotalEarned,
				Reward:       affiliate.Reward,
				TotalWagered: affiliate.TotalWagered,
				Rate:         getAffiliateRate(affiliate.CustomAffiliateRate),
			},
		)
	}
	return affiliateMetas, nil
}

// @Internal
// Get appropriate affiliate rate from custom rate, min rate, and max rate.
// - If custom rate is less than min rate, return min rate.
// - If custom rate is bigger than max rate, return max rate.
// - If custom rate is in range of [min, max], return custom rate.
func getAffiliateRate(customRate uint) uint {
	if customRate < config.AFFILIATE_RATE_MIN {
		return config.AFFILIATE_RATE_MIN
	}
	if customRate > config.AFFILIATE_RATE_MAX {
		return config.AFFILIATE_RATE_MAX
	}
	return customRate
}

// @External
// Distribute affiliate rewards to activated code owner.
func distributeAffiliate(from User, feeAmount int64, wagerAmount int64, sessionId ...UUID) (int64, error) {
	// 1. Validate parameter.
	if from == 0 {
		return 0, utils.MakeError(
			"affiliate_db",
			"distributeAffiliate",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}
	if feeAmount == 0 {
		return 0, nil
	}

	// 2. Retrieve active affiliate.
	activeAffiliate, err := getActiveAffiliateCode(from, sessionId...)
	if err != nil ||
		activeAffiliate == nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"distributeAffiliate",
			"failed to retrieve active affiliate",
			err,
		)
	}
	activeCode := activeAffiliate.Code
	customRate := activeAffiliate.Rate
	if activeCode == "" {
		return 0, nil
	}

	// 3. Calculate distribute amount.
	distributed := feeAmount * int64(getAffiliateRate(customRate)) / 100 / 2

	if distributed == 0 {
		return 0, nil
	}

	// 4. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"distributeAffiliate",
			"failed to retrieve session",
			err,
		)
	}

	// 5. Retrieve affiliate.
	affiliate := models.Affiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"code = ?",
		activeCode,
	).First(&affiliate); result.Error != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"distributeAffiliate",
			"failed to retrieve affilate",
			result.Error,
		)
	}

	// 6. Update affiliate reward.
	affiliate.Reward += distributed
	affiliate.TotalEarned += distributed
	affiliate.TotalWagered += wagerAmount
	if result := session.Save(&affiliate); result.Error != nil {
		return 0, utils.MakeError(
			"affiliate_db",
			"distributeAffiliate",
			"failed to update affilate rewards",
			result.Error,
		)
	}

	// 7. Update affiliate lifetime statistics. Doesn't return error object but
	// only prints since it is not a critical issue.
	if err := distributeAffiliateLifetimeUnchecked(
		uint(from),
		affiliate.ID,
		wagerAmount,
		distributed,
		sessionId...,
	); err != nil {
		log.LogMessage(
			"affiliate_db_distributeAffiliateRewards",
			"failed to update affiliate lifetime stats",
			"error",
			logrus.Fields{
				"error":       err.Error(),
				"userID":      from,
				"affiliateID": affiliate.ID,
				"wagered":     wagerAmount,
				"reward":      distributed,
			},
		)
	}

	return distributed, nil
}

// @External
// Deactivate affiliate code.
func deactivateAffiliateCode(user User, code string, sessionId ...UUID) error {
	// 1. Validate parameter.
	if user == 0 {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}
	if code == "" {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"invalid parameter",
			errors.New("provided code is empty string"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve affiliate.
	affiliate := models.Affiliate{}
	if result := session.Where(
		"code = ?",
		code,
	).First(&affiliate); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"failed to retrieve affiliate",
			result.Error,
		)
	}

	// 4. Deactivate affiliate.
	activeAffiliate := models.ActiveAffiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"user_id = ?",
		user,
	).First(&activeAffiliate); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"failed to retrieve active affiliate",
			result.Error,
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"no activated code",
			fmt.Errorf("user: %d, code: %s", user, code),
		)
	}
	if activeAffiliate.AffiliateID != affiliate.ID {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"mismatching active affiliate id",
			fmt.Errorf("user: %d, active: %d, prequestid: %d, code: %s",
				user,
				activeAffiliate.AffiliateID,
				affiliate.ID,
				code,
			),
		)
	}

	if result := session.Delete(
		&models.ActiveAffiliate{},
		activeAffiliate.ID,
	); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateCode",
			"failed to update delete affiliate id",
			result.Error,
		)
	}

	// 5. Update affiliate lifetime. Doesn't return error object but only prints
	// since it is not critical issue.
	if err := deactivateAffiliateLifetimeUnchecked(
		activeAffiliate.UserID,
		activeAffiliate.AffiliateID,
		sessionId...,
	); err != nil {
		log.LogMessage(
			"affiliate_db_deactivateAffiliateCode",
			"failed to update affiliate lifetime",
			"error",
			logrus.Fields{
				"error":       err.Error(),
				"userID":      activeAffiliate.UserID,
				"affiliateID": activeAffiliate.AffiliateID,
			},
		)
	}

	return nil
}

// @External
// Sets custom rate of a affiliate code.
func setAffiliateCustomRate(code string, customRate uint, sessionId ...UUID) error {
	// 1. Validate parameters.
	if code == "" {
		return utils.MakeError(
			"affiliate_db",
			"setAffiliateCustomRate",
			"invalid parameter",
			errors.New("provided code is empty string"),
		)
	}

	// ================== Archived ==================
	// if customRate < config.AFFILIATE_RATE_MIN ||
	// 	customRate > config.AFFILIATE_RATE_MAX {
	// 	return utils.MakeError(
	// 		"affiliate_db",
	// 		"setAffiliateCustomRate",
	// 		"invalid parameter",
	// 		errors.New("custom rate is out of range"),
	// 	)
	// }
	// ================== Archived ==================

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"setAffiliateCustomRate",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve the session
	affiliateInfo := models.Affiliate{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"code = ?",
		code,
	).First(&affiliateInfo); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"setAffiliateCustomRate",
			"failed to retrieve affiliate",
			result.Error,
		)
	}

	// 4. Update custom rate.
	if affiliateInfo.CustomAffiliateRate == customRate {
		return nil
	}
	affiliateInfo.CustomAffiliateRate = customRate
	if result := session.Save(&affiliateInfo); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"setAffiliateCustomRate",
			"failed to update custom affiliate rate",
			result.Error,
		)
	}

	return nil
}

/**
* @Internal
* Sets `LastActivated` field as current timestamp.
* This function doesn't check about already having activated code or other cases.
* Should be called after successful activation of affiliate code.
 */
func activateAffiliateLifetimeUnchecked(
	userID uint,
	affiliateID uint,
	sessionId ...UUID,
) error {
	// 1. Validate parameter
	if userID == 0 ||
		affiliateID == 0 {
		return utils.MakeError(
			"affiliate_db",
			"acitveAffiliateLifetimeUnchecked",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, affiliateID: %d",
				userID, affiliateID,
			),
		)
	}

	// 2. Retrieve session
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"activateAffiliateLifetimeUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update `LastActivated` as current timestamp.
	// Create a new lifetime record if not exists for userID and affiliateID.
	affiliateLifetime := models.AffiliateLifetime{}
	if result := session.Where(
		"user_id = ?",
		userID,
	).Where(
		"affiliate_id = ?",
		affiliateID,
	).First(
		&affiliateLifetime,
	); result.Error != nil &&
		!errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.MakeError(
			"affiliate_db",
			"activeAffiliateLifetimeUnchecked",
			"failed to retrieve affiliate lifetime record",
			fmt.Errorf(
				"userID: %d, affiliateID: %d, err: %v",
				userID, affiliateID, result.Error,
			),
		)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		affiliateLifetime.UserID = userID
		affiliateLifetime.AffiliateID = affiliateID
		affiliateLifetime.LastActivated = time.Now()
		affiliateLifetime.IsActive = true
		if result := session.Create(&affiliateLifetime); result.Error != nil {
			return utils.MakeError(
				"affiliate_db",
				"activeAffiliateLifetimeUnchecked",
				"failed to create a new lifetime record on not found",
				fmt.Errorf(
					"userID: %d, affiliateID: %d, err: %v",
					userID, affiliateID, result.Error,
				),
			)
		}
	} else if result.Error == nil {
		affiliateLifetime.LastActivated = time.Now()
		affiliateLifetime.IsActive = true
		if result := session.Save(&affiliateLifetime); result.Error != nil {
			return utils.MakeError(
				"affiliate_db",
				"activeAffiliateLifetimeUnchecked",
				"failed to update last activate time for existing record",
				fmt.Errorf(
					"userID: %d, affiliateID: %d, err: %v",
					userID, affiliateID, result.Error,
				),
			)
		}
	} else {
		log.LogMessage(
			"affiliate_db",
			"what's this? your if statement is ridiculous, lol",
			"error",
			logrus.Fields{},
		)
	}

	return nil
}

/**
* @Internal
* Sets `LastDeactivated` and `Lifetime` fields on deactivation of affiliate code.
* Doesn't check about context of activated code or other cases.
* Should be called after successful deactivation of affiliate code.
 */
func deactivateAffiliateLifetimeUnchecked(
	userID uint,
	affiliateID uint,
	sessionId ...UUID,
) error {
	// 1. Validate parameter
	if userID == 0 ||
		affiliateID == 0 {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateLifetimeUnchecked",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, affiliateID: %d",
				userID, affiliateID,
			),
		)
	}

	// 2. Retrieve session
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateLifetimeUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve affiliate life time.
	affiliateLifetime := models.AffiliateLifetime{}
	if result := session.Where(
		"user_id = ?",
		userID,
	).Where(
		"affiliate_id = ?",
		affiliateID,
	).First(
		&affiliateLifetime,
	); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateLifetimeUnchecked",
			"failed to retrieve affiliate lifetime",
			fmt.Errorf(
				"userID: %d, affiliateID: %d, err: %v",
				userID, affiliateID, result.Error,
			),
		)
	}

	// 4. Update `LastDeactivated` and `Lifetime` field.
	now := time.Now()
	affiliateLifetime.LastDeactivated = &now
	lifeTimeForLastActivation := now.Sub(affiliateLifetime.LastActivated).Seconds()
	if lifeTimeForLastActivation < 0 {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateLifetimeUnchecked",
			"negative lifetime is detected",
			fmt.Errorf(
				"userID: %d, affiliateID: %d, lastActivated: %v, plusLifetime: %f",
				userID,
				affiliateID,
				affiliateLifetime.LastActivated,
				lifeTimeForLastActivation,
			),
		)
	}
	affiliateLifetime.Lifetime += uint(lifeTimeForLastActivation)
	affiliateLifetime.IsActive = false
	if result := session.Save(&affiliateLifetime); result.Error != nil {
		return utils.MakeError(
			"affiliate_db",
			"deactivateAffiliateLifetimeUnchecked",
			"failed to update affilaite lifetime",
			fmt.Errorf(
				"userID: %d, affiliateID: %d, err: %v",
				userID, affiliateID, result.Error,
			),
		)
	}

	return nil
}

/**
* @Internal
* Update wagered and reward for affiliate lifetime record.
* Should be called on successful distsribution of affiliate code.
 */
func distributeAffiliateLifetimeUnchecked(
	userID uint,
	affiliateID uint,
	wagered int64,
	reward int64,
	sessionId ...UUID,
) error {
	// 1. Validate parameter
	if userID == 0 ||
		affiliateID == 0 {
		return utils.MakeError(
			"affiliate_db",
			"distributeAffiliateLifetimeUnchecked",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, affiliateID: %d",
				userID, affiliateID,
			),
		)
	}

	// 2. Retrieve session
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"distributeAffiliateLifetimeUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update affiliate lifetime's stats.
	if result := session.Model(
		&models.AffiliateLifetime{},
	).Where(
		"user_id = ?",
		userID,
	).Where(
		"affiliate_id = ?",
		affiliateID,
	).Updates(
		map[string]interface{}{
			"total_wagered": gorm.Expr("total_wagered + ?", wagered),
			"total_reward":  gorm.Expr("total_reward + ?", reward),
		},
	); result.Error != nil ||
		result.RowsAffected != 1 {
		return utils.MakeError(
			"affiliate_db",
			"distributeAffiliateLifetimeUnchecked",
			"failed to update affiliate lifetime",
			fmt.Errorf(
				"userID: %d, affiliateID: %d, err: %v, rowsAffected: %d",
				userID, affiliateID, result.Error, result.RowsAffected,
			),
		)
	}

	return nil
}

/**
* @External
* Retrieves detail of affiliate code with the users' lifetime and rewards.
 */
func getAffiliateDetail(
	code string,
	sessionId ...UUID,
) (*AffiliateDetail, error) {
	// 1. Validate parameters.
	if code == "" {
		return nil, utils.MakeError(
			"affiliate_db",
			"getAffiliateDetail",
			"invalid parameter",
			errors.New("provided code argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"affiliate_db",
			"getAffiliateDetail",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve affiliate detail.
	affiliate := models.Affiliate{}
	if result := session.Preload(
		"AffiliateLifetimes.User",
	).Where(
		"code = ?",
		code,
	).First(&affiliate); result.Error != nil {
		return nil, utils.MakeError(
			"affiliate_db",
			"getAffiliateDetail",
			"failed to retrieve affiliate",
			fmt.Errorf(
				"code: %s, result: %v",
				code, result.Error,
			),
		)
	}

	// 4. Build affiliate detail.
	affiliateDetail := AffiliateDetail{
		Code:  affiliate.Code,
		Users: []UserInAffiliateDetail{},
	}
	for _, item := range affiliate.AffiliateLifetimes {
		if !item.IsActive {
			continue
		}
		affiliateDetail.Users = append(
			affiliateDetail.Users,
			UserInAffiliateDetail{
				ID:     item.User.ID,
				Name:   item.User.Name,
				Avatar: item.User.Avatar,
				Lifetime: item.Lifetime +
					uint(time.Since(item.LastActivated).Seconds()),
				Wagered: item.TotalWagered,
				Reward:  item.TotalReward,
			},
		)
	}

	return &affiliateDetail, nil
}

/**
* @External
* Updates affiliate's isFirstDepositBonus property.
* Performs in main session.
* Return previous type and error object on failure.
 */
func UpdateAffiliateFirstDepositBonus(
	affiliateCode string,
	isFirstDepositBonus bool,
) (bool, error) {
	// 1. Validate parameter.
	if len(affiliateCode) == 0 {
		return false, utils.MakeError(
			"affiliate_db",
			"UpdateAffiliateFirstDepositBonus",
			"invalid parameter",
			errors.New("provided affiliate code is empty string"),
		)
	}

	// 2. Get main session.
	session, err := getSession()
	if err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"UpdateAffiliateFirstDepositBonus",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Get affiliate type.
	affiliatInfo := models.Affiliate{}
	if err := session.Where(
		"code = ?",
		affiliateCode,
	).First(&affiliatInfo).Error; err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"UpdateAffiliateFirstDepositBonus",
			"failed to retrieve affiliate info",
			fmt.Errorf(
				"affiliateCode: %s, err: %v",
				affiliateCode, err,
			),
		)
	}

	// 4. Update affiliate type.
	if affiliatInfo.IsFirstDepositBonus == isFirstDepositBonus {
		return isFirstDepositBonus, nil
	}
	previous := affiliatInfo.IsFirstDepositBonus
	affiliatInfo.IsFirstDepositBonus = isFirstDepositBonus
	if err := session.Save(&affiliatInfo).Error; err != nil {
		return false, utils.MakeError(
			"affiliate_db",
			"UpdateAffiliateFirstDepositBonus",
			"failed to update affiliate type",
			fmt.Errorf(
				"affiliateCode: %s, affiliate1stDBonus: %v, err: %v",
				affiliateCode, isFirstDepositBonus, err,
			),
		)
	}

	return previous, nil
}

/**
* @External
* Update active affiliates first deposit done as true.
 */
func SetActiveAffiliateFirstDepositDone(
	affiliateID uint,
) error {
	// 1. Validate parameters.
	if affiliateID == 0 {
		return utils.MakeError(
			"affiliate_db",
			"SetActiveAffiliateFirstDepositDone",
			"invalid parameter",
			errors.New("provided affiliate id is zero"),
		)
	}

	// 2. Retrieve main session.
	session, err := getSession()
	if err != nil {
		return utils.MakeError(
			"affiliate_db",
			"SetActiveAffiliateFirstDepositDone",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Update first deposit done as true.
	if result := session.Model(
		&models.ActiveAffiliate{},
	).Where(
		"id = ?",
		affiliateID,
	).Update(
		"first_deposit_done",
		true,
	); result.Error != nil ||
		result.RowsAffected != 1 {
		return utils.MakeError(
			"affiliate_db",
			"SetActiveAffiliateFirstDepositDone",
			"failed to update first deposit done as true",
			fmt.Errorf(
				"id: %d, err: %v, rowsAffected: %d",
				affiliateID, result.Error, result.RowsAffected,
			),
		)
	}

	return nil
}
