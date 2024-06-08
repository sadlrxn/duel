package crash

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/*
/* @Internal
/* Loads next round, if c.round is nil, it load first unplayed round.
/* This function *SHOULD NOT* be called by other part than this file.
*/
func (c *GameController) loadNextRoundUnchecked() error {
	// 0. Check whether crash is blocked by admin.
	if c.isBlockCrash {
		return utils.MakeError(
			"crash_controller_db",
			"loadNextRoundUnchecked",
			"crash is blocked by admin",
			errors.New("isBlockCrash flag is set"),
		)
	}

	// 1. Get next crash round id.
	nextRoundID := uint(0)
	if c.round != nil {
		nextRoundID = c.round.ID + 1
	} else {
		nextRoundID = getFirstUnplayedRoundID()
	}

	// 2. Retrieve crashRound for that id.
	round, err := lockAndRetrieveCrashRound(
		nextRoundID,
		db_aggregator.MainSessionId(),
	)
	if err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadNextRoundUnchecked",
			"failed to retrieve crash round",
			fmt.Errorf(
				"roundID: %d, err : %v",
				nextRoundID, err,
			),
		)
	}

	// 3. Update the controller's round.
	c.round = round

	// 4. Update outcome if missing.
	if err := c.updateRoundOutcomeIfMissing(); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadNextRoundUnchecked",
			"failed to update outcome of new round",
			fmt.Errorf(
				"round: %v, err: %v",
				c.round, err,
			),
		)
	}

	return nil
}

/*
/* @Internal
/* This function is called on initializing GameController.
*/
func (c *GameController) loadRoundOnInit() error {
	// 1. Check whether controller's round is nil.
	if c.round != nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadRoundOnInit",
			"already initialized round",
			fmt.Errorf("round: %v", c.round),
		)
	}

	// 2. Load first unplayed round.
	if err := c.loadNextRoundUnchecked(); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadRoundOnInit",
			"failed to load first unplayed round",
			fmt.Errorf("err: %v", err),
		)
	}

	return nil
}

/*
/* @Internal
/* This function is called on loading next round after previous round.
*/
func (c *GameController) loadNextRound() error {
	// 1. Check whether controller's round is nil.
	if c.round == nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadNextRound",
			"prev round is nil pointer",
			fmt.Errorf("round: %v", c.round),
		)
	}

	// 2. Load first unplayed round.
	if err := c.loadNextRoundUnchecked(); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"loadNextRound",
			"failed to load next round",
			fmt.Errorf("err: %v", err),
		)
	}

	return nil
}

/*
/* @Internal
/* This function is called when updating outcome.
*/
func (c *GameController) updateRoundOutcomeIfMissing() error {
	// 1. Validate status for updating outcome of round.
	if c.round == nil ||
		!c.isPreparingStatus() ||
		!c.isPreparingRound() {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundOutcomeIfMissing",
			"not valid case for updating outcome",
			fmt.Errorf(
				"round: %v, roundStatus: %v",
				c.round, c.roundStatus,
			),
		)
	}
	if c.round.Outcome > 0 {
		return nil
	}

	// 2. Calculate out come for the current round.
	serverConfig := config.GetServerConfig()
	outcome := calculateOutCome(
		c.round.Seed,
		serverConfig.CrashClientSeed,
		c.houseEdge,
	)

	// 3. Update outcome.
	if err := updateCrashRoundOutcome(
		c.round,
		outcome,
	); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundOutcomeIfMissing",
			"failed to update crash outcome",
			fmt.Errorf(
				"round: %v, err: %v",
				c.round, outcome,
			),
		)
	}

	return nil
}

/*
* @Internal
* This function updates round's betStartedAt in *MAIN* session.
* This function is called in `startBetting`
* before update the round status.
 */
func (c *GameController) updateRoundBetStartedAt() error {
	// 1. Validate status to update round's betStartedAt.
	if !c.isPreparingStatus() ||
		!c.isRoundInitializedForBetting() {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundBetStartedAt",
			"not valid case for updating betStartedAt",
			fmt.Errorf(
				"round: %v, roundStatus: %v",
				c.round, c.roundStatus,
			),
		)
	}

	// 2. Update betStartedAt in the main session.
	if err := updateCrashRoundBetStartedAt(
		c.round,
		time.Now(),
		db_aggregator.MainSessionId(),
	); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundBetStartedAt",
			"failed to update betStartedAt",
			err,
		)
	}

	return nil
}

/*
* @Internal
* This function updates round's runStartedAt in session after
* *LOCKING* crash round record.
* This function is called in `startRunning` before updating the
* round status.
 */
func (c *GameController) updateRoundRunStartedAt() error {
	// 1. Validate status to update round's runStartedAt.
	if !c.isPendingStatus() ||
		!c.isPendingRound() {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundRunStartedAt",
			"not valid case for updating runStartedAt",
			fmt.Errorf(
				"round: %v, roundStatus: %v",
				c.round, c.roundStatus,
			),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundRunStartedAt",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve crash round record.
	round, err := lockAndRetrieveCrashRound(
		c.round.ID,
		sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundRunStartedAt",
			"failed to lock and retrieve crash round",
			fmt.Errorf(
				"roundId: %d, error: %v",
				c.round.ID, err,
			),
		)
	}

	// 4. Update crash round's runStartedAt.
	if err := updateCrashRoundRunStartedAt(
		round,
		time.Now(),
		sessionId,
	); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundRunStartedAt",
			"failed to update run started at of crash round",
			fmt.Errorf(
				"round: %v, runStartedAt: %v, err: %v",
				round, time.Now(), err,
			),
		)
	}

	// 5. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundRunStartedAt",
			"failed to commit session",
			err,
		)
	}

	// 6. Update controller's round.
	round.Bets = c.round.Bets
	round.FeeTransaction = c.round.FeeTransaction
	c.round = round

	return nil
}

/*
* @Internal
* This function updates round's endedAt in session after
* *LOCKING* crash round record.
* This function is called in `prepareBetting` before fetching
* new round's data.
* This function should be called in preparing status and ended round.
* At this point, round.EndedAt is not saved in database.
 */
func (c *GameController) updateRoundEndedAt() error {
	// 1. Validate status to update round's runStartedAt.
	if !c.isPreparingStatus() ||
		!c.isRunningRound() {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundEndedAt",
			"not valid case for updating runStartedAt",
			fmt.Errorf(
				"round: %v, roundStatus: %v",
				c.round, c.roundStatus,
			),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundEndedAt",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve crash round record.
	round, err := lockAndRetrieveCrashRound(
		c.round.ID,
		sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundEndedAt",
			"failed to lock and retrieve crash round",
			fmt.Errorf(
				"roundId: %d, error: %v",
				c.round.ID, err,
			),
		)
	}

	// 4. Update crash round's endedAt.
	if err := updateCrashRoundEndedAt(
		round,
		c.lastStatusUpdated,
		sessionId,
	); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundEndedAt",
			"failed to update ended at of crash round",
			fmt.Errorf(
				"round: %v, endedAt: %v, err: %v",
				round, time.Now(), err,
			),
		)
	}

	// 5. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"crash_controller_db",
			"updateRoundEndedAt",
			"failed to commit session",
			err,
		)
	}

	// 6. Update controller's round.
	round.Bets = c.round.Bets
	round.FeeTransaction = c.round.FeeTransaction
	c.round = round

	return nil
}
