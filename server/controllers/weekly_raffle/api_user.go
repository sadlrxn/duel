package weekly_raffle

import (
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func prepareWeeklyRaffleStatus(userID uint, count uint) *WeeklyRaffleStatus {
	if pendingUntil != nil {
		return prepareWeeklyRaffleStatusForPending(userID, count)
	} else {
		return prepareWeeklyRaffleStatusForRunning(userID, count)
	}
}

func prepareWeeklyRaffleStatusForRunning(
	userID uint,
	count uint,
) *WeeklyRaffleStatus {
	var myRank int = -1
	var myTickets []WeeklyRaffleTicketInStatus
	if userID != 0 {
		myRank = redis.GetUserWeeklyRaffleRank(userID)
		myTickets = getUserWeeklyRaffleTickets(userID)
	}

	if myRank >= 0 &&
		myRank < int(count) {
		count++
	}

	playerIDs, ticketCounts, err := redis.GetWeeklyRaffleTicketsPerUser(count)
	if err != nil {
		log.LogMessage(
			"prepareWeeklyRaffleStatus",
			"failed to get weekly raffle tickets per user",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return nil
	}

	var result = WeeklyRaffleStatus{
		Status:             WeeklyRaffleStatusRunning,
		Players:            []UserInWeeklyRaffleStatus{},
		Remaining:          getWeeklyRaffleRemainingTime(),
		TotalTicketsIssued: getWeeklyRaffleTotalTicketsIssued(),
		ChipsPerTicket:     config.WEEKLY_RAFFLE_CHIPS_WAGER_PER_TICKET,
		TotalPrize:         getWeeklyRaffleTotalPrize(),
	}

	if len(playerIDs) == 0 {
		return &result
	}

	if myRank != -1 {
		myInfo := user.GetUserInfoByID(userID)
		if myInfo != nil {
			result.Me = UserInWeeklyRaffleStatus{
				ID:          myInfo.ID,
				Name:        myInfo.Name,
				Avatar:      myInfo.Avatar,
				Rank:        myRank + 1,
				TicketCount: uint(len(myTickets)),
				Tickets:     myTickets,
			}
		}
	}

	for i := range playerIDs {
		playerInfo := user.GetUserInfoByID(playerIDs[i])
		if playerInfo != nil && playerIDs[i] != userID {
			playerData := utils.GetUserDataWithPermissions(
				*playerInfo,
				nil,
				0,
			)
			result.Players = append(
				result.Players,
				UserInWeeklyRaffleStatus{
					ID:          playerData.ID,
					Name:        playerData.Name,
					Avatar:      playerData.Avatar,
					Rank:        i + 1,
					TicketCount: ticketCounts[i],
				},
			)
		}
	}

	return &result
}

func prepareWeeklyRaffleStatusForPending(
	userID uint,
	count uint,
) *WeeklyRaffleStatus {
	lastWeeklyRaffle, err := getLastWeeklyRaffle()
	if err != nil {
		return nil
	}

	var result = WeeklyRaffleStatus{
		Status:         WeeklyRaffleStatusPending,
		Players:        []UserInWeeklyRaffleStatus{},
		Remaining:      uint(time.Until(*pendingUntil).Seconds()),
		ChipsPerTicket: config.WEEKLY_RAFFLE_CHIPS_WAGER_PER_TICKET,
		TotalPrize:     utils.ConvertChipToBalance(2000),
	}

	if lastWeeklyRaffle == nil {
		return &result
	}

	var totalPrize int64
	for _, prize := range lastWeeklyRaffle.Prizes {
		totalPrize += prize
	}

	result.TotalPrize = totalPrize
	result.TotalTicketsIssued = getWeeklyRaffleTotalTicketsIssuedForRound(lastWeeklyRaffle.StartedAt)

	index := getIndexFromStartedAt(lastWeeklyRaffle.StartedAt)
	var myRank int = -1
	var myTickets []WeeklyRaffleTicketInStatus
	if userID != 0 {
		myRank = redis.GetUserWeeklyRaffleRank(userID, index)
		myTickets = getUserWeeklyRaffleTicketsForRound(userID, lastWeeklyRaffle.StartedAt)
	}

	if myRank >= 0 &&
		myRank < int(count) {
		count++
	}

	playerIDs, ticketCounts, err := redis.GetWeeklyRaffleTicketsPerUser(count, index)
	if err != nil {
		log.LogMessage(
			"prepareWeeklyRaffleStatus",
			"failed to get weekly raffle tickets per user",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return nil
	}

	if len(playerIDs) == 0 {
		return &result
	}

	if myRank != -1 {
		myInfo := user.GetUserInfoByID(userID)
		if myInfo != nil {
			result.Me = UserInWeeklyRaffleStatus{
				ID:          myInfo.ID,
				Name:        myInfo.Name,
				Avatar:      myInfo.Avatar,
				Rank:        myRank + 1,
				TicketCount: uint(len(myTickets)),
				Tickets:     myTickets,
			}
		}
	}

	for i := range playerIDs {
		playerInfo := user.GetUserInfoByID(playerIDs[i])
		if playerInfo != nil && playerIDs[i] != userID {
			playerData := utils.GetUserDataWithPermissions(
				*playerInfo,
				nil,
				0,
			)
			result.Players = append(
				result.Players,
				UserInWeeklyRaffleStatus{
					ID:          playerData.ID,
					Name:        playerData.Name,
					Avatar:      playerData.Avatar,
					Rank:        i + 1,
					TicketCount: ticketCounts[i],
				},
			)
		}
	}

	return &result
}

func GetWeeklyRaffleStatusHandler(ctx *gin.Context) {
	var params struct {
		Count uint `form:"count"`
	}

	if err := ctx.Bind(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
			},
		)
		return
	}

	if result := prepareWeeklyRaffleStatus(
		middlewares.GetAuthUserID(ctx, false),
		params.Count,
	); result == nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to get weekly raffle status",
			},
		)
	} else {
		ctx.JSON(
			http.StatusOK,
			result,
		)
	}
}

func GetWeeklyRaffleRewardsHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	rewards := getWeeklyRaffleRewardsForUser(userID)
	if rewards == nil {
		rewards = []WeeklyRaffleRewardsStatus{}
	}

	tickets := getUserWeeklyRaffleTickets(userID)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"rewards": rewards,
			"tickets": len(tickets),
		},
	)
}

func ClaimWeeklyRaffleRewardsHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	var params struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
			},
		)
		return
	}

	if claimed, err := claimReward(
		userID,
		params.IDs,
	); err == nil {
		ctx.JSON(
			http.StatusOK,
			claimed,
		)
	} else {
		log.LogMessage(
			"weekly_raffle_ClaimWeeklyRaffleRewardsHandler",
			"failed to claim rewards",
			"error",
			logrus.Fields{
				"userID":    userID,
				"rewardIDs": params.IDs,
				"error":     err.Error(),
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to claim rewards.",
			},
		)
	}
}
