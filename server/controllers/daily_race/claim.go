package daily_race

import (
	"fmt"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/**
* Claims user's daily race prize.
* Returns claimed amount and an error object.
 */
func claimReward(
	userID uint,
	rewardIDs []uint,
) (int64, error) {
	// 1. Validate parameters.
	if userID == 0 ||
		rewardIDs == nil ||
		len(rewardIDs) == 0 {
		return 0, utils.MakeErrorWithCode(
			"daily_race_claim",
			"claimReward",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"userID: %d, rewardIDs: %d",
				userID, rewardIDs,
			),
		)
	}

	// 2. Start a session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"daily_race_claim",
			"claimReward",
			"failed to start a new session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve daily race rewards record.
	dailyRaceRewards, err := lockAndRetrieveUnclaimedDailyRewards(
		userID,
		rewardIDs,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"daily_race_claim",
			"claimReward",
			"failed to lock and retrieve daily reward records",
			fmt.Errorf(
				"rewardIDs: %v, err: %v",
				rewardIDs, err,
			),
		)
	}

	// 4. Update claimed of the rewards.
	totalClaimed, err := updateDailyRewardsClaimed(
		dailyRaceRewards,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"daily_race_claim",
			"claimReward",
			"failed to update daily rewards claimed",
			fmt.Errorf(
				"dailyRewards: %v, err: %v",
				dailyRaceRewards, err,
			),
		)
	}

	// 5. Perform chip transaction and leave history.
	for i, dailyReward := range dailyRaceRewards {
		if _, err := giveChipsForClaim(
			userID,
			dailyReward.Prize,
			dailyReward.ID,
			sessionId,
		); err != nil {
			return 0, utils.MakeError(
				"daily_race_claim",
				"claimReward",
				"failed to give chips for claim",
				fmt.Errorf(
					"i: %d, record: %v, err: %v",
					i, dailyReward, err,
				),
			)
		}
	}

	// 6. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"daily_race_claim",
			"claimReward",
			"failed to commit session",
			err,
		)
	}

	return totalClaimed, nil
}

/**
* @Internal
* Gives chips for claim.
* Returns generated tx id, and error object.
* Utilize dailyRaceRewardID as ownerID of polymorphic association.
 */
func giveChipsForClaim(
	userID uint,
	amount int64,
	dailyRaceRewardID uint,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		amount <= 0 ||
		dailyRaceRewardID == 0 {
		return 0, utils.MakeError(
			"daily_race_claim",
			"giveChipsForClaim",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, amount: %d, dailyRaceRewardID: %d",
				userID, amount, dailyRaceRewardID,
			),
		)
	}

	// 2. Give chips for claiming to the user.
	txResult, err := db_aggregator.Transfer(
		(*db_aggregator.User)(&config.DAILY_RACE_TEMP_ID),
		(*db_aggregator.User)(&userID),
		&db_aggregator.BalanceLoad{
			ChipBalance: &amount,
		},
		sessionId,
	)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient funds") &&
			strings.Contains(err.Error(), "removeChipsFromUser") {
			return 0, utils.MakeErrorWithCode(
				"coupon_exchange",
				"giveChipsForClaim",
				"insufficient admin temp wallet balance",
				ErrCodeInsufficientAdminBalance,
				err,
			)
		}
		return 0, utils.MakeError(
			"coupon_exchange",
			"giveChipsForClaim",
			"failed to perform real chips transfer",
			err,
		)
	}

	// 3. Leave transaction.
	transactionHistory := models.Transaction{
		FromWallet: (*uint)(txResult.FromWallet),
		ToWallet:   (*uint)(txResult.ToWallet),
		Balance: models.Balance{
			ChipBalance: &models.ChipBalance{
				Balance: amount,
			},
		},
		Type:   models.TxClaimDailyRaceReward,
		Status: models.TransactionSucceed,

		FromWalletPrevID: (*uint)(txResult.FromPrevBalance),
		FromWalletNextID: (*uint)(txResult.FromNextBalance),
		ToWalletPrevID:   (*uint)(txResult.ToPrevBalance),
		ToWalletNextID:   (*uint)(txResult.ToNextBalance),
		OwnerID:          dailyRaceRewardID,
		OwnerType:        models.TransactionDailyRaceRewardsReferenced,
	}
	if err := db_aggregator.LeaveRealTransaction(
		&transactionHistory,
		sessionId,
	); err != nil {
		return 0, utils.MakeError(
			"coupon_exchange",
			"giveChipsForClaim",
			"failed to leave transaction",
			err,
		)
	}

	return transactionHistory.ID, nil
}
