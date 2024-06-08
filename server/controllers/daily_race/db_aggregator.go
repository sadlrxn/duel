package daily_race

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
* @Internal
* Returns approved daily race rewards which belong to the user.
 */
func getDailyRaceRewardsForUser(
	userID uint,
) []DailyRaceRewardsStatus {
	if userID == 0 {
		return nil
	}

	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	dailyRaceRewards := []models.DailyRaceRewards{}
	if err := session.Where(
		"user_id = ?",
		userID,
	).Where(
		"approved = ?",
		true,
	).Where(
		"claimed = ?",
		0,
	).Find(&dailyRaceRewards).Error; err != nil {
		return nil
	}

	result := []DailyRaceRewardsStatus{}
	for _, reward := range dailyRaceRewards {
		result = append(
			result,
			DailyRaceRewardsStatus{
				ID:    reward.ID,
				Date:  time.Time(reward.StartedAt),
				Prize: reward.Prize,
				Rank:  reward.Rank + 1,
			},
		)
	}

	return result
}

/**
* @Internal
* Lock and retrieve unclaimed(`claimed == 0`) and approved(`approved == true`)
* daily reward records by userID and rewardIDs.
* If not matching retrieved record count with length of rewardIDs array,
* returns error.
 */
func lockAndRetrieveUnclaimedDailyRewards(
	userID uint,
	rewardIDs []uint,
	sessionId db_aggregator.UUID,
) ([]models.DailyRaceRewards, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		len(rewardIDs) == 0 {
		return nil, utils.MakeError(
			"daily_race_db",
			"lockAndRetrieveUnclaimedDailyRewards",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, rewardIDs: %v",
				userID, rewardIDs,
			),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"coupon_db",
			"lockAndRetrieveUnclaimedDailyRewards",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve daily race rewards.
	dailyRaceRewards := []models.DailyRaceRewards{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where(
		"user_id = ?",
		userID,
	).Where(
		"id in ?",
		rewardIDs,
	).Where(
		"claimed = ?",
		0,
	).Where(
		"approved = ?",
		true,
	).Find(&dailyRaceRewards); result.Error != nil {
		return nil, utils.MakeError(
			"daily_race_db",
			"lockAndRetrieveUnclaimedDailyRewards",
			"failed to retrieve daily race rewards",
			fmt.Errorf(
				"userID: %d, rewardIDs: %v, err: %v",
				userID, rewardIDs, result.Error,
			),
		)
	} else if len(dailyRaceRewards) != len(rewardIDs) {
		return nil, utils.MakeError(
			"daily_race_db",
			"lockAndRetrieveUnclaimedDailyRewards",
			"mistmatching record count with rewardIDs arg",
			fmt.Errorf(
				"rewardIDs count: %d, retrieved count: %d",
				len(rewardIDs), len(dailyRaceRewards),
			),
		)
	}

	return dailyRaceRewards, nil
}

/**
* @Internal
* Set claimed of daily rewards as prize.
* If claimed is not zero returns error.
 */
func updateDailyRewardsClaimed(
	dailyRewards []models.DailyRaceRewards,
	sessionId db_aggregator.UUID,
) (int64, error) {
	// 1. Validate parameter.
	if len(dailyRewards) == 0 {
		return 0, utils.MakeError(
			"daily_race_db",
			"updateDailyRewardsClaimed",
			"invalid parameter",
			errors.New("provided dailyRewards is empty slice"),
		)
	}

	// 2. Retrieve the session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"daily_race_db",
			"updateDailyRewardsClaimed",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Prepare daily rewards claimed.
	totalPrize := int64(0)
	for i, reward := range dailyRewards {
		if reward.ID == 0 ||
			reward.Claimed != 0 {
			return 0, utils.MakeError(
				"daily_race_db",
				"updateDailyRewardsClaimed",
				"invalid prameter",
				fmt.Errorf(
					"rewardID: %d, reward.Calimed: %d",
					reward.ID, reward.Claimed,
				),
			)
		}
		dailyRewards[i].Claimed = reward.Prize
		totalPrize += reward.Prize
	}

	// 4. Update records.
	if result := session.Model(
		&dailyRewards,
	).Clauses(
		clause.Returning{},
	).Update(
		"claimed", gorm.Expr("prize"),
	); result.Error != nil || result.RowsAffected != int64(len(dailyRewards)) {
		return 0, utils.MakeError(
			"daily_race_db",
			"updateDailyRewardsClaimed",
			"failed to update claimed of daily rewards",
			fmt.Errorf(
				"err: %v, rowsAffected: %d, dailyRewards count: %d",
				result.Error, result.RowsAffected, len(dailyRewards),
			),
		)
	}

	return totalPrize, nil
}

/**
* @Internal
* Create daily race rewards all at once in one trasaction.
 */
func createDailyRaceRewardsInSession(
	rewards []models.DailyRaceRewards,
) error {
	// 1. Validate parameter.
	if len(rewards) == 0 {
		return utils.MakeError(
			"daily_race_db",
			"createDailyRaceRewardsInSession",
			"invalid parameter",
			errors.New("provided rewards is empty slice"),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"daily_race_db",
			"createDailyRaceRewardsInSession",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Create records.
	if err := session.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&rewards).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return utils.MakeError(
			"daily_race_db",
			"createDailyRaceRewardsInSession",
			"failed to create daily race rewards",
			fmt.Errorf(
				"rewards: %v, err: %v",
				rewards, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Sets daily rewards as approved one.
* If rewardIDs is nil or zero length,
* approve all unapproved ones since daysAgo from now.
 */
func approveDailyRaceReward(
	rewardIDs []uint,
	daysAgo uint,
) ([]models.DailyRaceRewards, error) {
	// 1. Validate parameter.
	if daysAgo == 0 {
		return nil, nil
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"daily_race_db",
			"approveDailyRaceReward",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Perform updating approved field in on session.
	rewards := []models.DailyRaceRewards{}
	since := time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo))
	if err := session.Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(
			&rewards,
		).Clauses(
			clause.Returning{},
		)
		if len(rewardIDs) > 0 {
			tx = tx.Where(
				"id in ?",
				rewardIDs,
			)
		}
		if result := tx.Where(
			"approved = ?",
			false,
		).Where(
			"started_at > ?",
			time.Date(
				since.Year(),
				since.Month(),
				since.Day(),
				0, 0, 0, 0,
				time.Local,
			),
		).Where(
			"claimed = ?",
			0,
		).Update(
			"approved",
			true,
		); result.Error != nil {
			return result.Error
		} else if len(rewardIDs) > 0 &&
			result.RowsAffected != int64(len(rewardIDs)) {
			return fmt.Errorf(
				"mismatching updated count. expected: %d, actual: %d",
				len(rewardIDs), result.RowsAffected,
			)
		}

		return nil
	}); err != nil {
		return nil, utils.MakeError(
			"daily_race_db",
			"approveDailyRaceReward",
			"failed to perform approving",
			fmt.Errorf(
				"rewardIDs: %v, daysAgo: %d, err: %v",
				rewardIDs, daysAgo, err,
			),
		)
	}

	return rewards, nil
}

/**
* @Internal
* Gets unapproved daily race rewards since days ago.
 */
func getUnapprovedDailyRaceRewards(
	daysAgo uint,
) []models.DailyRaceRewards {
	if daysAgo == 0 {
		return nil
	}

	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	since := time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo))
	rewards := []models.DailyRaceRewards{}
	if err := session.Where(
		"approved = ?",
		false,
	).Where(
		"started_at > ?",
		time.Date(
			since.Year(),
			since.Month(),
			since.Day(),
			0, 0, 0, 0,
			time.Local,
		),
	).Where(
		"claimed = ?",
		0,
	).Order(
		"started_at desc",
	).Find(&rewards).Error; err != nil {
		log.LogMessage(
			"daily_race_db_getUnapprovedDailyRaceRewards",
			"failed to retrieve daily race rewards",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return nil
	}

	return rewards
}
