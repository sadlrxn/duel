package self_exclusion

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
* @External
* Update `Until` field of `SelfExclusion`.
 */
func exclude(userID uint, days uint) error {
	// 1. Validate parameter.
	if userID == 0 ||
		days == 0 {
		return utils.MakeError(
			"self_exclusion_db",
			"Exclude",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, days: %d",
				userID, days,
			),
		)
	}

	// 2. Retrieve `SelfExclusion`.
	selfExclusion, err := retrieveSelfExclusion(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if selfExclusion, err = createUnchecked(userID); err != nil {
			return utils.MakeError(
				"self_exclusion_db",
				"Exclude",
				"failed to create new selfExclusion",
				fmt.Errorf(
					"userID: %d, err: %v",
					userID, err,
				),
			)
		}
	} else if err != nil {
		return utils.MakeError(
			"self_exclusion_db",
			"Exclude",
			"failed to retrieve selfExclusion",
			fmt.Errorf(
				"userID: %d, err: %v",
				userID, err,
			),
		)
	}

	// 3. Update self excluded.
	selfExclusion.Until = time.Now().Add(time.Hour * 24 * time.Duration(days))
	if err := updateUntilUnchecked(selfExclusion); err != nil {
		return utils.MakeError(
			"self_exclusion_db",
			"Exclude",
			"failed to update self exclusion",
			fmt.Errorf(
				"selfExclusion: %v, err: %v",
				*selfExclusion, err,
			),
		)
	}

	return nil
}

/**
* @External
* Check whether specific user is in self-excluded status.
* If error happens while retriving user's self exclusion, return false.
 */
func ExclusionRemaining(userID uint) uint {
	selfExclusion, err := retrieveSelfExclusion(userID)
	if err != nil {
		return 0
	}
	if time.Now().Before(selfExclusion.Until) {
		return uint(time.Until(selfExclusion.Until) / time.Second)
	}
	return 0
}

/**
* @Internal
* Retrieve `SelfExclusion` record for specific user.
* This action isn't performed in session since doesn't take any security problem.
 */
func retrieveSelfExclusion(userID uint) (*models.SelfExclusion, error) {
	// 1. Validate parameter.
	if userID == 0 {
		return nil, utils.MakeError(
			"self_exclusion_db",
			"retrieveSelfExclusion",
			"invalid parameter",
			fmt.Errorf("userID: %d", userID),
		)
	}

	// 2. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"self_exclusion_db",
			"retrieveSelfExclusion",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Get `SelfExclusion` record.
	selfExclusion := models.SelfExclusion{}
	if result := session.First(
		&selfExclusion,
		userID,
	); result.Error != nil {
		//! Returns raw `result.Error` because caller of this fuction need to
		//! check whether the error is gorm.ErrRecordNotFound or not.
		return nil, result.Error
	}

	return &selfExclusion, nil
}

/**
* @Internal
* Create a new `SelfExclusion` record, without checkoug about duplication or
* presence of userID. Though since `UserID` is primaryKey of model,
* will through error case of userID duplication.
 */
func createUnchecked(userID uint) (*models.SelfExclusion, error) {
	// 1. Validate parameter.
	if userID == 0 {
		return nil, utils.MakeError(
			"self_exclusion_db",
			"createSelfExclusionUnchecked",
			"invalid parameter",
			fmt.Errorf("userID: %d", userID),
		)
	}

	// 2. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"self_exclusion_db",
			"createSelfExclusionUnchecked",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Create SelfExclusion record.
	selfExclusion := models.SelfExclusion{
		UserID: userID,
	}
	if result := session.Create(&selfExclusion); result.Error != nil {
		return nil, utils.MakeError(
			"self_exclusion_db",
			"createSelfExclusionUnchecked",
			"failed to create selfExclusion record",
			fmt.Errorf(
				"selfExclusion: %v, err: %v",
				selfExclusion, result.Error,
			),
		)
	}

	return &selfExclusion, nil
}

/**
* @Internal
* Updates selfExclusion record. Used to update `Until`.
* Doesn't check any comparison of time or presence of user.
 */
func updateUntilUnchecked(selfExclusion *models.SelfExclusion) error {
	// 1. Validate parameter.
	if selfExclusion == nil ||
		selfExclusion.UserID == 0 {
		return utils.MakeError(
			"self_exclusion_db",
			"updateSelfExclusionUnchecked",
			"invalid parameter",
			fmt.Errorf("selfExclusion: %v", selfExclusion),
		)
	}

	// 2. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"self_exclusion_db",
			"updateSelfExclusionUnchecked",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Update selfExclusion.
	if result := session.Model(
		selfExclusion,
	).Clauses(
		clause.Returning{},
	).Update(
		"until", selfExclusion.Until,
	); result.Error != nil || result.RowsAffected != 1 {
		return utils.MakeError(
			"self_exclusion_db",
			"updateSelfExclusionUnchecked",
			"failed ot update until of selfExclusion",
			fmt.Errorf(
				"selfExclusion: %v, err: %v",
				*selfExclusion, result.Error,
			),
		)
	}

	return nil
}

func getUserInfoByName(userName string) *models.User {
	if userName == "" {
		return nil
	}

	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	userInfo := models.User{}
	if result := session.Select(
		"id",
		"name",
		"wallet_address",
		"role",
		"avatar",
		"banned",
		"private_profile",
	).Where(
		"name = ?",
		userName,
	).First(
		&userInfo,
	); result.Error != nil {
		return nil
	}

	return &userInfo
}

/**
* @External
* Set self exclusion's until as current timestamp.
 */
func remove(userName string) error {
	// 1. Validate parameter.
	if userName == "" {
		return utils.MakeError(
			"self_exclusion_db",
			"remove",
			"invalid parameter",
			errors.New("provided userName is empty string"),
		)
	}

	// 2. Retrieve user info.
	userID := uint(0)
	if userInfo := getUserInfoByName(
		userName,
	); userInfo == nil {
		return utils.MakeError(
			"self_exclusion_db",
			"remove",
			"failed to retrieve user info by name",
			fmt.Errorf("name: %s", userName),
		)
	} else {
		userID = userInfo.ID
	}

	// 2. Retrieve `SelfExclusion`.
	selfExclusion, err := retrieveSelfExclusion(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		return utils.MakeError(
			"self_exclusion_db",
			"remove",
			"failed to retrieve selfExclusion",
			fmt.Errorf(
				"userID: %d, err: %v",
				userID, err,
			),
		)
	}

	// 3. Update self until as current time.
	selfExclusion.Until = time.Now()
	if err := updateUntilUnchecked(selfExclusion); err != nil {
		return utils.MakeError(
			"self_exclusion_db",
			"remove",
			"failed to remove self exclusion",
			fmt.Errorf(
				"selfExclusion: %v, err: %v",
				*selfExclusion, err,
			),
		)
	}

	return nil
}
