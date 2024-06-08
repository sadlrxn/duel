package crash

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
/* @Internal
/* Lock and retrieve round record.
/* Should lock and retrieve where
/*   ID is matching provided roundID
*/
func lockAndRetrieveCrashRound(
	roundID uint,
	sessionId db_aggregator.UUID,
) (*models.CrashRound, error) {
	// 1. Validate parameter.
	if roundID == 0 {
		return nil, utils.MakeErrorWithCode(
			"crash_db",
			"lockAndRetrieveCrashRound",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided round id is 0"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"crash_db",
			"lockAndRetrieveCrashRound",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve crash round record.
	crashRound := models.CrashRound{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).First(
		&crashRound,
		roundID,
	); result.Error != nil {
		return nil, utils.MakeError(
			"crash_db",
			"lockAndRetrieveCrashRound",
			"failed to retrieve round record",
			result.Error,
		)
	}

	return &crashRound, nil
}

/*
/* @Internal
/* Lock and retrieve crash bet record.
/* Should lock and retrieve where
/*   ID is matching provided betID
*/
func lockAndRetrieveCrashBet(
	betID uint,
	sessionId db_aggregator.UUID,
) (*models.CrashBet, error) {
	// 1. Validate parameter.
	if betID == 0 {
		return nil, utils.MakeErrorWithCode(
			"crash_db",
			"lockAndRetrieveCrashBet",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided bet id is 0"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"crash_db",
			"lockAndRetrieveCrashBet",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve crash round record.
	crashBet := models.CrashBet{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).First(
		&crashBet,
		betID,
	); result.Error != nil {
		return nil, utils.MakeError(
			"crash_db",
			"lockAndRetrieveCrashBet",
			"failed to retrieve bet record",
			result.Error,
		)
	}

	return &crashBet, nil
}

/*
/* @Internal
/* Get bet counts made by user in specific round.
*/
func getBetCountMadeByUserForRound(
	userID uint,
	roundID uint,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter
	if userID == 0 ||
		roundID == 0 {
		return 0, utils.MakeErrorWithCode(
			"crash_db",
			"getBetCountMadeByUserForRound",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"userID: %d, roundID: %d",
				userID, roundID,
			),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"crash_db",
			"getBetCountMadeByUserForRound",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Get bet count made by user for round.
	var betCount int64
	if result := session.Model(
		&models.CrashBet{},
	).Where(
		"round_id = ? and user_id = ?",
		roundID, userID,
	).Count(&betCount); result.Error != nil {
		return 0, utils.MakeError(
			"crash_db",
			"getBetCountMadeByUserForRound",
			"failed to get bet count",
			result.Error,
		)
	}

	return uint(betCount), nil
}

/*
/* @Internal
/* Create crashBet record.
/* Context of this function is that the crashBet record is firstly made
/* when the bet is placed.
/* Doesn't check about cashOutAt minimum requirement here.
*/
func createCrashBet(
	crashBet *models.CrashBet,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if crashBet == nil ||
		crashBet.UserID == 0 ||
		crashBet.RoundID == 0 ||
		crashBet.BetAmount == 0 ||
		crashBet.ID != 0 {
		return utils.MakeErrorWithCode(
			"crash_db",
			"createCrashBet",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("crashBet: %v", crashBet),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"createCrashBet",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Create crash bet record.
	if result := session.Create(&crashBet); result.Error != nil {
		return utils.MakeError(
			"crash_db",
			"createCrashBet",
			"failed to create crash bet record",
			fmt.Errorf(
				"crashBet: %v, err: %v",
				crashBet, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Updates crashBet's profit and payoutMultiplier.
*/
func updateCrashBetPayoutFields(
	crashBet *models.CrashBet,
	profit int64,
	payoutMultiplier float64,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if crashBet == nil ||
		crashBet.ID == 0 ||
		profit == 0 ||
		payoutMultiplier == 0 ||
		crashBet.Profit != nil ||
		crashBet.PayoutMultiplier != nil {
		return utils.MakeErrorWithCode(
			"crash_db",
			"updateCrashBetPayoutFields",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"crashBet: %v, profit: %d, multiplier: %f",
				crashBet, profit, payoutMultiplier,
			),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"updateCrashBetPayoutFields",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update profit and payoutMultiplier.
	if result := session.Model(
		crashBet,
	).Clauses(
		clause.Returning{},
	).Updates(
		map[string]interface{}{
			"profit":            profit,
			"payout_multiplier": payoutMultiplier,
		},
	); result.Error != nil {
		return utils.MakeError(
			"crash_bet",
			"updateCrashBetPayoutFields",
			"failed to update profit and payoutMultiplier",
			fmt.Errorf(
				"record: %v, err: %v",
				crashBet, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Create crashRound record.
/* This function is created for the first time before crash.
/* This function is performed on main session.
*/
func createCrashRound(crashRound *models.CrashRound) error {
	// 1. Validate parameter.
	if crashRound == nil ||
		crashRound.ID == 0 ||
		crashRound.Seed == "" ||
		crashRound.Outcome != 0 ||
		crashRound.BetStartedAt != nil ||
		crashRound.RunStartedAt != nil ||
		crashRound.EndedAt != nil {
		return utils.MakeErrorWithCode(
			"crash_db",
			"createCrashRound",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("crashRound: %v", crashRound),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"createCrashRound",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Create crash bet record.
	if result := session.Create(&crashRound); result.Error != nil {
		return utils.MakeError(
			"crash_db",
			"createCrashRound",
			"failed to create crash round record",
			fmt.Errorf(
				"crashBet: %v, err: %v",
				crashRound, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Update outcome of crashRound.
/* This function is called for initial outcome calculation, and
/* case handling for 0 outcome retrieved on round retrieve.
/* This function is called on main session.
*/
func updateCrashRoundOutcome(
	crashRound *models.CrashRound,
	outcome float64,
) error {
	// 1. Validate parameter.
	if crashRound == nil ||
		crashRound.ID == 0 ||
		outcome == 0 ||
		crashRound.Outcome != 0 ||
		crashRound.BetStartedAt != nil ||
		crashRound.RunStartedAt != nil ||
		crashRound.EndedAt != nil {
		return utils.MakeErrorWithCode(
			"crash_db",
			"updateCrashRoundOutcome",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"crashRound: %v, outcome: %f",
				crashRound, outcome,
			),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"updateCrashRoundOutcome",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update outcome.
	if result := session.Model(
		crashRound,
	).Clauses(
		clause.Returning{},
	).Update(
		"outcome", outcome,
	); result.Error != nil {
		return utils.MakeError(
			"crash_bet",
			"updateCrashRoundOutcome",
			"failed to update outcome",
			fmt.Errorf(
				"record: %v, err: %v",
				crashRound, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Update betStartedAt of crashRound.
*/
func updateCrashRoundBetStartedAt(
	crashRound *models.CrashRound,
	betStartedAt time.Time,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if crashRound == nil ||
		crashRound.ID == 0 ||
		crashRound.BetStartedAt != nil ||
		crashRound.RunStartedAt != nil ||
		crashRound.EndedAt != nil ||
		crashRound.Outcome == 0 {
		return utils.MakeErrorWithCode(
			"crash_db",
			"updateCrashRoundBetStartedAt",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("crashRound: %v", crashRound),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"updateCrashRoundBetStartedAt",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update betStartedAt.
	if result := session.Model(
		crashRound,
	).Clauses(
		clause.Returning{},
	).Update(
		"bet_started_at", betStartedAt,
	); result.Error != nil {
		return utils.MakeError(
			"crash_bet",
			"updateCrashRoundBetStartedAt",
			"failed to update betStartedAt",
			fmt.Errorf(
				"record: %v, err: %v",
				crashRound, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Update runStartedAt of crashRound.
*/
func updateCrashRoundRunStartedAt(
	crashRound *models.CrashRound,
	runStartedAt time.Time,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if crashRound == nil ||
		crashRound.ID == 0 ||
		crashRound.BetStartedAt == nil ||
		crashRound.RunStartedAt != nil ||
		crashRound.EndedAt != nil ||
		crashRound.Outcome == 0 {
		return utils.MakeErrorWithCode(
			"crash_db",
			"updateCrashRoundRunStartedAt",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("crashRound: %v", crashRound),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"updateCrashRoundRunStartedAt",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update runStartedAt.
	if result := session.Model(
		crashRound,
	).Clauses(
		clause.Returning{},
	).Update(
		"run_started_at", runStartedAt,
	); result.Error != nil {
		return utils.MakeError(
			"crash_bet",
			"updateCrashRoundRunStartedAt",
			"failed to update runStartedAt",
			fmt.Errorf(
				"record: %v, err: %v",
				crashRound, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Update endedAt of crashRound.
*/
func updateCrashRoundEndedAt(
	crashRound *models.CrashRound,
	endedAt time.Time,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameter.
	if crashRound == nil ||
		crashRound.ID == 0 ||
		crashRound.BetStartedAt == nil ||
		crashRound.RunStartedAt == nil ||
		crashRound.EndedAt != nil ||
		crashRound.Outcome == 0 {
		return utils.MakeErrorWithCode(
			"crash_db",
			"updateCrashRoundEndedAt",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("crashRound: %v", crashRound),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"crash_db",
			"updateCrashRoundEndedAt",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update endedAt.
	if result := session.Model(
		crashRound,
	).Clauses(
		clause.Returning{},
	).Update(
		"ended_at", endedAt,
	); result.Error != nil {
		return utils.MakeError(
			"crash_bet",
			"updateCrashRoundEndedAt",
			"failed to update endedAt",
			fmt.Errorf(
				"record: %v, err: %v",
				crashRound, result.Error,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* Auto migrate crashRound table.
*/
func autoMigrateCrashRound() error {
	// 1. Retrieve db pointer.
	db, err := db_aggregator.GetSession()
	if err != nil {
		utils.MakeError(
			"crash_db",
			"autoMigrateCrashRound",
			"failed to retrieve db pointer",
			err,
		)
	}

	// 2. Migrate crash round table
	return db.AutoMigrate(
		&models.CrashRound{},
		&models.CrashBet{},
	)
}

/*
/* @Internal
/* Get total round count in crash round model.
*/
func getTotalCrashRoundCount() int64 {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return -1
	}

	// 2. Get total crash round count.
	count := int64(0)
	if result := session.Model(
		&models.CrashRound{},
	).Count(&count); result.Error != nil {
		return -1
	}

	return count
}

/*
/* @Internal
/* Get first unplayed round id.
*/
func getFirstUnplayedRoundID() uint {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return 0
	}

	// 2. Get first unplayed round id.
	roundID := uint(0)
	session.Model(
		&models.CrashRound{},
	).Select(
		"id as roundID",
	).Where(
		"bet_started_at is null",
	).Where(
		"run_started_at is null",
	).Where(
		"ended_at is null",
	).Order(
		"id",
	).Limit(1).Row().Scan(&roundID)

	return roundID
}

/*
* @Internal
* Get 10 round id, and multipliers.
 */
func getRoundHistory(currentRoundID uint) []RoundHistoryItem {
	history := []RoundHistoryItem{}

	session, err := db_aggregator.GetSession()
	if err != nil {
		return history
	}

	session.Model(
		&models.CrashRound{},
	).Select(
		"id", "outcome",
	).Where(
		"bet_started_at is not null",
	).Where(
		"run_started_at is not null",
	).Where(
		"ended_at is not null",
	).Where(
		"id < ?",
		currentRoundID,
	).Order(
		"id desc",
	).Limit(10).Scan(&history)

	return history
}

/**
* @External
* Get round detail for the provided round id.
* Returns error when:
* - roundID is zero or greater than max round count. `ErrCodeInvalidParameter`
* - Round not found with roundID and endedAt is not null. `ErrCodeNotFoundFinishedRound`
 */
func GetRoundHistoryDetail(roundID uint) (*RoundHistoryDetail, error) {
	// 1. Validate parameter.
	if roundID == 0 ||
		roundID > 2500000 {
		return nil, utils.MakeErrorWithCode(
			"crash_db",
			"GetRoundHistoryDetail",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("roundID: %d", roundID),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"crash_db",
			"GetRoundHistoryDetail",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Retrieve round record.
	roundInfo := models.CrashRound{}
	if result := session.Preload(
		"Bets.User",
	).Where(
		"id = ?",
		roundID,
	).Where(
		"bet_started_at is not null",
	).Where(
		"run_started_at is not null",
	).Where(
		"ended_at is not null",
	).First(&roundInfo); errors.Is(
		result.Error,
		gorm.ErrRecordNotFound,
	) {
		return nil, utils.MakeErrorWithCode(
			"crash_db",
			"GetRoundHistoryDetail",
			"round not found",
			ErrCodeNotFoundFinishedRound,
			fmt.Errorf(
				"roundID: %d, error: %v",
				roundID, err,
			),
		)
	} else if result.Error != nil {
		return nil, utils.MakeError(
			"crash_db",
			"GetRoundHistoryDetail",
			"failed to retrieve round info",
			fmt.Errorf(
				"roundID: %d, error: %v",
				roundID, err,
			),
		)
	}

	// 4. Build return meta.
	result := RoundHistoryDetail{
		ID:      roundInfo.ID,
		Seed:    roundInfo.Seed,
		Outcome: roundInfo.Outcome,
		Date:    *roundInfo.EndedAt,
		Bets:    []BetInRoundHistoryDetail{},
	}
	for _, bet := range roundInfo.Bets {
		user := user.GetUserInfoByID(bet.User.ID)
		if user == nil {
			log.LogMessage(
				"CrashRoundHistoryDetail",
				"Failed to get bet user info",
				"error",
				logrus.Fields{
					"UserID": bet.User.ID,
				},
			)
			continue
		}
		userInfo := utils.GetUserDataWithPermissions(*user, nil, 0)
		betInfo := BetInRoundHistoryDetail{
			User:            userInfo,
			BetAmount:       bet.BetAmount,
			PaidBalanceType: bet.PaidBalanceType,
		}
		if bet.Profit != nil &&
			bet.PayoutMultiplier != nil {
			betInfo.Profit = *bet.Profit
			betInfo.PayoutMultiplier = *bet.PayoutMultiplier
		}
		result.Bets = append(
			result.Bets,
			betInfo,
		)
	}

	return &result, nil
}
