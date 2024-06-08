package weekly_raffle

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/**
* @Internal
* Claims weekly raffle reward.
* Returns claimed amount and error object.
 */
func claimReward(
	userID uint,
	ticketIDs []uint,
) (int64, error) {
	// 1. Validate parameters.
	if userID == 0 {
		return 0, utils.MakeError(
			"weekly_raffle_claim",
			"claimReward",
			"invalid parameter",
			errors.New("provided userID is zero"),
		)
	}
	if len(ticketIDs) == 0 {
		return 0, nil
	}

	// 3. Start a session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_claim",
			"claimReward",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve unclaimed user tickets.
	tickets, err := lockAndRetrieveUnclaimedTickets(
		userID,
		ticketIDs,
		sessionId,
	)
	if err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_claim",
			"claimReward",
			"failed to lock and retrieve weekly raffle tickets",
			fmt.Errorf(
				"userID: %d, tickets: %v, err: %v",
				userID, ticketIDs, err,
			),
		)
	}

	// 5. Get total amount.
	totalClaimed := int64(0)
	for i, ticket := range tickets {
		if len(ticket.Round.Prizes) <= int(*ticket.Rank) {
			return 0, utils.MakeError(
				"weekly_raffle_claim",
				"claimReward",
				"index out of prizes slice range",
				fmt.Errorf(
					"prizeLen: %d, ticketRank: %d, ticketID: %d,",
					len(ticket.Round.Prizes), *ticket.Rank, ticket.ID,
				),
			)
		}
		prize := ticket.Round.Prizes[*ticket.Rank]
		totalClaimed += prize
		tickets[i].Claimed = &prize
		if _, err := giveChipsForClaim(
			userID,
			prize,
			ticket.ID,
			sessionId,
		); err != nil {
			return 0, utils.MakeError(
				"weekly_raffle_claim",
				"claimReward",
				"failed to transfer chips for weekly raffle reward",
				fmt.Errorf(
					"userID: %d, prize: %d, ticketID: %d, index: %d, error: %v",
					userID, prize, ticket.ID, i, err,
				),
			)
		}
		if err := updateRaffleTicketsClaimed(
			&tickets[i],
			sessionId,
		); err != nil {
			return 0, utils.MakeError(
				"weekly_raffle_claim",
				"claimReward",
				"update raffle ticket claimed",
				fmt.Errorf(
					"ticket: %v, index: %d, error: %v",
					ticket, i, err,
				),
			)
		}
	}

	// 6. Commit session
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_claim",
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
	weeklyRaffleTicketID uint,
	sessionId db_aggregator.UUID,
) (uint, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		amount <= 0 ||
		weeklyRaffleTicketID == 0 {
		return 0, utils.MakeError(
			"daily_race_claim",
			"giveChipsForClaim",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, amount: %d, weeklyRaffleTicketID: %d",
				userID, amount, weeklyRaffleTicketID,
			),
		)
	}

	// 2. Give chips for claiming to the user.
	txResult, err := db_aggregator.Transfer(
		(*db_aggregator.User)(&config.WEEKLY_RAFFLE_TEMP_ID),
		(*db_aggregator.User)(&userID),
		&db_aggregator.BalanceLoad{
			ChipBalance: &amount,
		},
		sessionId,
	)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient funds") &&
			strings.Contains(err.Error(), "removeChipsFromUser") {
			return 0, utils.MakeError(
				"coupon_exchange",
				"giveChipsForClaim",
				"insufficient admin temp wallet balance",
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
		Type:   models.TxClaimWeeklyRaffleReward,
		Status: models.TransactionSucceed,

		FromWalletPrevID: (*uint)(txResult.FromPrevBalance),
		FromWalletNextID: (*uint)(txResult.FromNextBalance),
		ToWalletPrevID:   (*uint)(txResult.ToPrevBalance),
		ToWalletNextID:   (*uint)(txResult.ToNextBalance),
		OwnerID:          weeklyRaffleTicketID,
		OwnerType:        models.TransactionWeeklyRaffleRewardReferenced,
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
