package transaction

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func getRakebackRewards(user *db_aggregator.User) (int64, int64, error) {
	if user == nil || *user == 0 {
		return 0, 0, utils.MakeError(
			"transaction",
			"getRakebackRewards",
			"invalid parameter",
			fmt.Errorf("provided parameter user: %p", user),
		)
	}

	rakebackInfo, err := db_aggregator.GetRakebackInfo(*user)
	if err != nil {
		return 0, 0, utils.MakeError(
			"transaction",
			"getRakebackRewards",
			"failed to retrieve rakeback info",
			err,
		)
	}

	return rakebackInfo.TotalEarned, rakebackInfo.Reward, nil
}

// @External
// Claim rakeback rewards.
func claimRakeback(user *db_aggregator.User) (int64, error) {
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimRakeback",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	rewards, err := db_aggregator.ClaimRakeback(*user, sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimRakeback",
			"failed to claim rake back",
			err,
		)
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"transaction",
			"claimRakeback",
			"failed to commit session",
			err,
		)
	}

	return rewards, nil
}
