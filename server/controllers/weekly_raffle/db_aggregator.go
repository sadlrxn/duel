package weekly_raffle

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
* @Internal
* Creates weekly raffle record.
* Doesn't check about current round presence.
 */
func createWeeklyRaffleUnchecked(weeklyRaffle *models.WeeklyRaffle) error {
	// 1. Validate parameter.
	if weeklyRaffle == nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"createWeeklyRaffleUnchecked",
			"invalid parameter",
			errors.New("provided weeklyRaffle is nil pointer"),
		)
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"createWeeklyRaffleUnchecked",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Create weeklyRaffle record.
	if err := session.Create(weeklyRaffle).Error; err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"createWeeklyRaffleUnchecked",
			"failed to create weeklyRaffle record",
			fmt.Errorf(
				"record: %v, error: %v",
				*weeklyRaffle, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Try to retrieve current weekly raffle record.
* If not found gorm error, returns nil error with nil record.
 */
func retrieveCurWeeklyRaffle() (*models.WeeklyRaffle, error) {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"retrieveCurWeeklyRaffle",
			"failed to retrieve main session",
			err,
		)
	}

	// 2. Retrieve cur session.
	weeklyRaffle := models.WeeklyRaffle{}
	if err := addCurrentWeeklyRaffleQuery(
		session,
	).Last(&weeklyRaffle).Error; errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil, nil
	} else if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"retrieveCurWeeklyRaffle",
			"failed to retrieve current weekly raffle",
			err,
		)
	}

	return &weeklyRaffle, nil
}

/**
* @Internal
* Adds current weekly raffle query context to tx.
 */
func addCurrentWeeklyRaffleQuery(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return nil
	}
	return tx.Where(
		"end_at > ?",
		time.Now(),
	).Where(
		"ended = ?",
		false,
	)
}

/**
* @Internal
* Returns not performed but ended weekly raffle record started at
* the `startedAt`.
 */
func retrieveNotPerformedWeeklyRaffle(
	startedAt datatypes.Date,
) (*models.WeeklyRaffle, error) {
	// 1. Retrieve a main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"retrieveNotPerformedWeeklyRaffle",
			"failed to retrieve main session",
			err,
		)
	}

	// 2. Retrieve not performed weekly raffle.
	weeklyRaffle := models.WeeklyRaffle{}
	if err := session.Where(
		"started_at = ?",
		startedAt,
	).Where(
		"ended = ?",
		false,
	).Where(
		"end_at < ?",
		time.Now(),
	).Last(&weeklyRaffle).Error; err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"retrieveNotPerformedWeeklyRaffle",
			"failed to retrieve not performed weeklyRaffle",
			err,
		)
	}

	return &weeklyRaffle, nil
}

/**
* @Internal
* Updates ranks of winning tickets.
* Takes list of winning tickets' numbers as argument,
* first ticket is rank #1 and saved as 0 in the db.
 */
func updateRanksOfWinningTickets(
	startedAt datatypes.Date,
	winningTickets []uint,
	sessionId db_aggregator.UUID,
) ([]uint, error) {
	// 1. Validate parameters.
	if len(winningTickets) == 0 {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"invalid parameter",
			errors.New("provided winningTickets is empty slice"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Build query to set ranks.
	var updateRanksQuery = `UPDATE weekly_raffle_tickets
SET rank = CASE
%s
END
WHERE round_started_at = @started_at_value
AND ticket_id in @ticket_id_values;`
	var caseRelationQuery = ""
	var namedArgs = map[string]interface{}{
		"started_at_value": startedAt,
		"ticket_id_values": winningTickets,
	}
	for i, ticket := range winningTickets {
		caseRelationQuery = fmt.Sprintf(
			"%s WHEN ticket_id = @ticket_id_value_%d THEN %d",
			caseRelationQuery, i, i,
		)
		namedArgs[fmt.Sprintf(
			"ticket_id_value_%d",
			i,
		)] = ticket
	}

	if result := session.Exec(
		fmt.Sprintf(
			updateRanksQuery,
			caseRelationQuery,
		),
		namedArgs,
	); result.Error != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"failed to execute update query",
			result.Error,
		)
	} else if int(result.RowsAffected) != len(winningTickets) {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"mismatching updated rows with number of winning tickets",
			fmt.Errorf(
				"updatedRows: %d, winningTickets: %d",
				result.RowsAffected,
				len(winningTickets),
			),
		)
	}

	// 4. Retrieve winned users' IDs.
	if winnerIDs, err := reviewWinnerIDsFromTickets(
		startedAt,
		winningTickets,
		sessionId,
	); err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"failed to get winning userIDs",
			fmt.Errorf(
				"startedAt: %v, winningTickets: %v, error: %v",
				startedAt, winningTickets, err,
			),
		)
	} else {
		return winnerIDs, nil
	}
}

/**
* @Internal
* Review winnerIDs from tickets and startedAt.
 */
func reviewWinnerIDsFromTickets(
	startedAt datatypes.Date,
	winningTickets []uint,
	sessionId db_aggregator.UUID,
) ([]uint, error) {
	// 1. Validate parameters.
	if len(winningTickets) == 0 {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"reviewWinnerIDsFromTickets",
			"invalid parameter",
			errors.New("provided winningTickets is empty slice"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"reviewWinnerIDsFromTickets",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve winned users' IDs.
	var selectWinnerIDsQuery = `SELECT user_id 
FROM weekly_raffle_tickets
WHERE ticket_id in @ticket_id_values
AND round_started_at = @started_at_value
ORDER BY CASE
%s
END ASC;`
	var caseRelationQuery = ""
	var namedArgs = map[string]interface{}{
		"started_at_value": startedAt,
		"ticket_id_values": winningTickets,
	}
	for i, ticket := range winningTickets {
		caseRelationQuery = fmt.Sprintf(
			"%s WHEN ticket_id = @ticket_id_value_%d THEN %d",
			caseRelationQuery, i, i,
		)
		namedArgs[fmt.Sprintf(
			"ticket_id_value_%d",
			i,
		)] = ticket
	}

	if rows, err := session.Raw(
		fmt.Sprintf(
			selectWinnerIDsQuery,
			caseRelationQuery,
		),
		namedArgs,
	).Rows(); err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"reviewWinnerIDsFromTickets",
			"failed to select winner ids",
			err,
		)
	} else {
		defer rows.Close()
		winnerIDs := []uint{}
		for rows.Next() {
			winnerID := uint(0)
			rows.Scan(&winnerID)
			winnerIDs = append(winnerIDs, winnerID)
		}
		return winnerIDs, nil
	}
}

/**
* @Internal
* Update weekly raffle's ended to be true.
* Doesn't check about previous status.
 */
func setWeeklyRaffleEnded(
	weeklyRaffle *models.WeeklyRaffle,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameters.
	if weeklyRaffle == nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"setWeeklyRaffleEnded",
			"invalid parameter",
			errors.New("provided weeklyRaffle is nil pointer"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"setWeeklyRaffleEnded",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update weekly raffle's ended.
	if err := session.Model(
		weeklyRaffle,
	).Clauses(
		clause.Returning{},
	).Update(
		"ended",
		true,
	).Error; err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"setWeeklyRaffleEnded",
			"failed to update weeklyRaffle's ended as true",
			err,
		)
	}

	return nil
}

/**
* @Internal
* Returns last weekly raffle's prizes.
* If not exists returns nil.
 */
func getLastPrizes() []int64 {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	// 2. Retrieve last prizes.
	lastWeeklyRaffle := models.WeeklyRaffle{}
	if err := session.Where(
		"end_at < ?",
		time.Now(),
	).Last(
		&lastWeeklyRaffle,
	).Error; err != nil {
		return nil
	}

	return lastWeeklyRaffle.Prizes
}

/**
* @Internal
* Lock and retrieve unclaimed user tickets.
* Preloads weeklyRaffle.
* Where
* - matching userID, ticketIDs(primary key)
* - rank is not null
* - claimed is null
 */
func lockAndRetrieveUnclaimedTickets(
	userID uint,
	ticketIDs []uint,
	sessionId db_aggregator.UUID,
) ([]models.WeeklyRaffleTicket, error) {
	// 1. Validate parameters.
	if userID == 0 ||
		len(ticketIDs) == 0 {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"lockAndRetireveUnclaimedTickets",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, ticketIDs: %v",
				userID, ticketIDs,
			),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"lockAndRetireveUnclaimedTickets",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve user tickets.
	tickets := []models.WeeklyRaffleTicket{}
	if err := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"Round",
	).Where(
		"user_id = ?",
		userID,
	).Where(
		"ticket_id in ?",
		ticketIDs,
	).Where(
		"rank is not null",
	).Where(
		"claimed is null",
	).Find(&tickets).Error; err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"lockAndRetireveUnclaimedTickets",
			"failed to retrieve user tickets",
			fmt.Errorf(
				"userID: %d, ticketIDs: %v, error: %v",
				userID, ticketIDs, err,
			),
		)
	} else if len(tickets) != len(ticketIDs) {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"lockAndRetireveUnclaimedTickets",
			"mismatching retrieved tickets than requested ticketIDs",
			fmt.Errorf(
				"tickets: %d, ticketIDs: %d",
				len(tickets), len(ticketIDs),
			),
		)
	}

	return tickets, nil
}

/**
* @Internal
* Update weekly raffle's claimed.
 */
func updateRaffleTicketsClaimed(
	ticket *models.WeeklyRaffleTicket,
	sessionId db_aggregator.UUID,
) error {
	// 1. Validate parameters.
	if ticket == nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"updateRaffleTicketsClaimed",
			"invalid parameter",
			errors.New("provided ticket is nil pointer"),
		)
	}
	if ticket.ID == 0 {
		return utils.MakeError(
			"weekly_raffle_db",
			"updateRaffleTicketsClaimed",
			"invalid parameter",
			errors.New("provided ticket is with zero id"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"updateRaffleTicketsClaimed",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Update claimed.
	if err := session.Model(
		ticket,
	).Clauses(
		clause.Returning{},
	).Update(
		"claimed",
		ticket.Claimed,
	).Error; err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"updateRaffleTicketsClaimed",
			"failed to update weekly raffle ticket's claimed",
			fmt.Errorf(
				"ticketID: %d, claimed: %d, error: %s",
				ticket.ID, *ticket.Claimed, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Returns the last ticket id for the startedAt.
* If not found, returns 0.
 */
func retrieveMaxTicketID(
	startedAt datatypes.Date,
) (uint, error) {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_db",
			"retrieveMaxTicketID",
			"failed to retrieve main db session",
			err,
		)
	}

	// 2. Get max ticket ID.
	var maxTicketID sql.NullInt64
	if err := session.Model(
		&models.WeeklyRaffleTicket{},
	).Select(
		"max(ticket_id)",
	).Where(
		"round_started_at = ?",
		startedAt,
	).Row().Scan(&maxTicketID); err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_db",
			"retrieveMaxTicketID",
			"failed to get max ticket ID",
			fmt.Errorf(
				"startedAt: %v, error: %v",
				startedAt,
				err,
			),
		)
	}

	return uint(maxTicketID.Int64), nil
}

/**
* @Internal
* Issue weekly raffle tickets.
 */
func createWeeklyRaffleTickets(
	tickets *[]models.WeeklyRaffleTicket,
) error {
	// 1. Validate parameters.
	if tickets == nil ||
		len(*tickets) == 0 {
		return utils.MakeError(
			"weekly_raffle_db",
			"createWeeklyRaffleTickets",
			"invalid parameter",
			errors.New("provided tickets argument is empty"),
		)
	}

	// 2. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"weekly_raffel_db",
			"createWeeklyRaffleTickets",
			"failed to retrieve main session",
			err,
		)
	}

	// 3. Create new raffle tickets.
	if err := session.Create(
		tickets,
	).Error; err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"createWeeklylRaffleTickets",
			"failed to create new ticket records",
			fmt.Errorf(
				"tickets: %v, error: %v",
				tickets, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Updates prizes for specific weekly raffle prizes.
 */
func setWeeklyRafflePrizes(
	prizes []int64,
	startedAt datatypes.Date,
) error {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"setWeeklyRafflePrizes",
			"failed to retrieve main session",
			err,
		)
	}

	// 2. Update prizes of weekly raffle.
	if err := session.Model(
		&models.WeeklyRaffle{},
	).Where(
		"started_at = ?",
		startedAt,
	).Update(
		"prizes",
		prizes,
	).Error; err != nil {
		return utils.MakeError(
			"weekly_raffle_db",
			"setWeeklyRafflePrizes",
			"failed to update prizes of weekly raffle",
			fmt.Errorf(
				"startedAt: %v, prizes: %v, error: %v",
				startedAt, prizes, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Get user's weekly raffle tickets.
* Return tickets for the current weekly raffle round.
 */
func getUserWeeklyRaffleTickets(userID uint) []WeeklyRaffleTicketInStatus {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return []WeeklyRaffleTicketInStatus{}
	}

	// 2. Get current weekly raffle.
	raffle, err := retrieveCurWeeklyRaffle()
	if err != nil || raffle == nil {
		return []WeeklyRaffleTicketInStatus{}
	}

	// 3. Get users' tickets.
	tickets := []models.WeeklyRaffleTicket{}
	if err := session.Where(
		"round_started_at = ?",
		raffle.StartedAt,
	).Where(
		"user_id = ?",
		userID,
	).Find(&tickets).Error; err != nil {
		return []WeeklyRaffleTicketInStatus{}
	}

	// 4. Build result.
	result := []WeeklyRaffleTicketInStatus{}
	for _, ticket := range tickets {
		result = append(
			result,
			WeeklyRaffleTicketInStatus{
				CreatedAt: ticket.CreatedAt,
				TicketID:  ticket.TicketID,
			},
		)
	}

	return result
}

/**
* @Internal
* Get user's weekly raffle tickets.
* Return tickets for the given weekly raffle round.
 */
func getUserWeeklyRaffleTicketsForRound(userID uint, startedAt datatypes.Date) []WeeklyRaffleTicketInStatus {
	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return []WeeklyRaffleTicketInStatus{}
	}

	// 2. Get users' tickets.
	tickets := []models.WeeklyRaffleTicket{}
	if err := session.Where(
		"round_started_at = ?",
		startedAt,
	).Where(
		"user_id = ?",
		userID,
	).Find(&tickets).Error; err != nil {
		return []WeeklyRaffleTicketInStatus{}
	}

	// 3. Build result.
	result := []WeeklyRaffleTicketInStatus{}
	for _, ticket := range tickets {
		result = append(
			result,
			WeeklyRaffleTicketInStatus{
				CreatedAt: ticket.CreatedAt,
				TicketID:  ticket.TicketID,
			},
		)
	}

	return result
}

/**
* @Internal
* Get weekly raffle remaining time in seconds.
 */
func getWeeklyRaffleRemainingTime() uint {
	// 1. Retrieve current weekly raffle.
	raffle, err := retrieveCurWeeklyRaffle()
	if err != nil || raffle == nil {
		return 0
	}

	return uint(time.Until(raffle.EndAt).Seconds())
}

/**
* @Internal
* Get weekly raffle total prize.
 */
func getWeeklyRaffleTotalPrize() int64 {
	// 1. Retrieve current weekly raffle.
	raffle, err := retrieveCurWeeklyRaffle()
	if err != nil || raffle == nil {
		return 0
	}

	var totalPrize int64
	for _, prize := range raffle.Prizes {
		totalPrize += prize
	}
	return totalPrize
}

/*
* @Internal
* Get weekly raffle total ticket count.
 */
func getWeeklyRaffleTotalTicketsIssued() uint {
	// 1. Retrieve curretn weekly raffle.
	raffle, err := retrieveCurWeeklyRaffle()
	if err != nil || raffle == nil {
		return 0
	}

	// 2. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return 0
	}

	// 3. Get total tickets count.
	count := int64(0)
	if err := session.Model(
		&models.WeeklyRaffleTicket{},
	).Where(
		"round_started_at = ?",
		raffle.StartedAt,
	).Count(&count).Error; err != nil {
		return 0
	}

	return uint(count)
}

/*
* @Internal
* Get weekly raffle total ticket count for given round.
 */
func getWeeklyRaffleTotalTicketsIssuedForRound(startedAt datatypes.Date) uint {

	// 1. Retrieve main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return 0
	}

	// 2. Get total tickets count.
	count := int64(0)
	if err := session.Model(
		&models.WeeklyRaffleTicket{},
	).Where(
		"round_started_at = ?",
		startedAt,
	).Count(&count).Error; err != nil {
		return 0
	}

	return uint(count)
}

/*
* @Internal
* Get user's weekly raffle rewards which unclaimed.
 */
func getWeeklyRaffleRewardsForUser(userID uint) []WeeklyRaffleRewardsStatus {
	// 1. Retrieve a main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		log.LogMessage(
			"getWeeklyRaffleRewardsForUser",
			"failed to retrieve main session",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return nil
	}

	// 2. Retrieve user tickets.
	var tickets = []models.WeeklyRaffleTicket{}
	if err := session.Preload(
		"Round",
	).Where(
		"user_id = ?",
		userID,
	).Where(
		"rank IS NOT NULL",
	).Where(
		"claimed IS NULL",
	).Find(&tickets).Error; err != nil {
		log.LogMessage(
			"getWeeklyRaffleRewardsForUser",
			"failed to retrieve user tickets",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return nil
	}

	// 3. Build result from tickets.
	var result = []WeeklyRaffleRewardsStatus{}
	for _, ticket := range tickets {
		result = append(
			result,
			WeeklyRaffleRewardsStatus{
				ID:    ticket.ID,
				Date:  time.Time(ticket.Round.StartedAt),
				Prize: ticket.Round.Prizes[*ticket.Rank],
				Rank:  *ticket.Rank,
			},
		)
	}

	return result
}

/*
* @Internal
* Get recently ended but still not performed weekly raffles order by latest.
 */
func getUnperformedWeeklyRaffles() []models.WeeklyRaffle {
	// 1. Retrieve a main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		log.LogMessage(
			"getUnperformedWeeklyRaffles",
			"failed to retrieve main session",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return []models.WeeklyRaffle{}
	}

	// 2. Retrieve unperformed weekly raffles.
	var weeklyRaffles = []models.WeeklyRaffle{}
	if err := session.Where(
		"ended = ?",
		false,
	).Find(&weeklyRaffles).Error; err != nil {
		log.LogMessage(
			"getUnperformedWeeklyRaffles",
			"failed to retrieve unperformed weekly raffles",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return []models.WeeklyRaffle{}
	}
	return weeklyRaffles
}

/*
* @Internal
* Get weekly raffle prizing preview result.
 */
func getPrizingPreviewResult(
	winningTickets []uint,
	startedAt time.Time,
) (*WeeklyRafflePrizingResult, error) {
	// 1. Validate parameters.
	if len(winningTickets) == 0 {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"getPrizingPreviewResult",
			"invalid parameter",
			fmt.Errorf(
				"tickets: %v, startedAt: %v",
				winningTickets, startedAt,
			),
		)
	}

	// 2. Retrieve weekly raffle from the date.
	startedDate := datatypes.Date(time.Date(
		startedAt.Year(),
		startedAt.Month(),
		startedAt.Day(),
		0, 0, 0, 0,
		time.Local,
	))
	weeklyRaffle, err := retrieveNotPerformedWeeklyRaffle(startedDate)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"getPrizingPreviewResult",
			"failed to retrieve not performed weekly raffle",
			fmt.Errorf(
				"startedAt: %v, err: %v",
				startedAt, err,
			),
		)
	}

	// 3. Checks whether length of winningTickets is equal to length of prizes.
	if len(weeklyRaffle.Prizes) != len(winningTickets) {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"mismatching number of tickets",
			fmt.Errorf(
				"prizes: %v, tickets: %v",
				weeklyRaffle.Prizes, winningTickets,
			),
		)
	}

	// 4. Retrieve user ranks from winning tickets.
	if winnerIDs, err := reviewWinnerIDsFromTickets(
		startedDate,
		winningTickets,
		db_aggregator.MainSessionId(),
	); err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"updateRanksOfWinningTickets",
			"failed to get winning userIDs",
			fmt.Errorf(
				"startedAt: %v, winningTickets: %v, error: %v",
				startedAt, winningTickets, err,
			),
		)
	} else {
		var result = WeeklyRafflePrizingResult{
			StartedAt: time.Time(startedDate),
		}
		var winners = []WinnerInWeeklyRafflePrizingResult{}
		for i, winnerID := range winnerIDs {
			winners = append(
				winners,
				WinnerInWeeklyRafflePrizingResult{
					UserID:   winnerID,
					TicketID: winningTickets[i],
					Rank:     uint(i + 1),
					Prize:    utils.ConvertBalanceToChip(weeklyRaffle.Prizes[i]),
				},
			)
		}
		result.Winners = winners
		return &result, nil
	}
}

/*
* @Internal
* Get the last weekly raffle round.
 */
func getLastWeeklyRaffle() (*models.WeeklyRaffle, error) {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffel_db",
			"getLastWeeklyRaffle",
			"failed to retrieve main session",
			err,
		)
	}

	var weeklyRaffle = models.WeeklyRaffle{}
	if err := session.Last(
		&weeklyRaffle,
	).Error; errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil, nil
	} else if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_db",
			"getLastWeeklyRaffle",
			"failed to retrieve last weekly raffle",
			err,
		)
	}

	return &weeklyRaffle, nil
}
