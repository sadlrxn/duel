package jackpot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/controllers/wager"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

func (c *Controller) setAvailable() {
	c.status = Available
	c.roundID = 0
	c.lastUpdated = time.Now()
	c.playerToBets = syncmap.Map{}
	c.betPlayers = []uint{}
	c.nftsInGame = []models.NftInGame{}
	c.totalPlayers = 0
	c.winnerID = nil
	c.ticketID = nil
	c.signedString = nil
	c.usdProfit = 0
	c.nfts4Profit = []types.NftDetails{}
	c.usdFee = 0
	c.nfts4Fee = []types.NftDetails{}
	c.totalAmount = 0
	c.totalFee = 0
	c.rollingDuration = 0
	c.candidates = []types.User{}
	c.countingTime = config.JACKPOT_COUNTING_TIME

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
		EventType: string(c.status)})
	c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}

	ticketID, err := utils.RequestTicketID()
	if err != nil {
		log.LogMessage("jackpot controller", "Failed to generate ticket ID", "error", logrus.Fields{"error": err.Error()})
		return
	}

	var round = models.JackpotRound{
		Players:  []models.JackpotPlayer{},
		TicketID: ticketID,
		Type:     c.Type,
	}
	db := db.GetDB()
	db.Create(&round)

	c.create(round.ID, ticketID)
}

func (c *Controller) create(roundID uint, ticketID string) {
	if c.status != Available {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Status Not Available"}})
		c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
		return
	}

	c.status = Created
	c.lastUpdated = time.Now()
	c.roundID = roundID
	c.ticketID = &ticketID

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID, TicketID: c.ticketID}})
	c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
	log.LogMessage("jackpot controller", "new round created", "info", logrus.Fields{"round": c.roundID, "ticket": *c.ticketID})
}

func (c *Controller) start(userID uint) {
	if c.status != Created {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "No Round Created"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	c.status = Started
	c.lastUpdated = time.Now()
	c.timer = time.NewTimer(time.Duration(c.countingTime) * time.Second)

	go func() {
		<-c.timer.C
		c.end()
	}()

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID}})
	c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
	log.LogMessage("jackpot controller", "round started", "info", logrus.Fields{"round": c.roundID})
}

func (c *Controller) end() {
	if c.status != Started {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Cannot End Round Not Started"}})
		c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
		return
	}
	c.status = Rolling
	c.lastUpdated = time.Now()
	c.timer = time.NewTimer(time.Duration(c.rollingTime) * time.Second)
	go func() {
		<-c.timer.C
		c.setAvailable()
	}()

	winnedCandidate, err := c.determineWinner()
	if err != nil {
		return
	}
	usdProfit, nfts4Profit, totalProfit, usdFee, nfts4Fee, totalFee, totalAmount, err := c.updateDbAtRoundEnd()
	if err != nil {
		return
	}

	c.usdProfit = usdProfit
	c.nfts4Profit = nfts4Profit
	c.usdFee = usdFee
	c.nfts4Fee = nfts4Fee
	c.totalAmount = totalAmount
	c.totalFee = totalFee

	c.chargeNFTsAsFee(nfts4Fee)
	c.prizeWinner(usdProfit, nfts4Profit)
	c.updateStatistics(totalProfit)
	c.emitRoundEndEvent(winnedCandidate)
}

func (c *Controller) checkRemainingTime() {
	if c.status == Started {
		currentTime := time.Now()
		bettingDuration := currentTime.Sub(c.lastUpdated)
		if bettingDuration >= time.Duration((c.countingTime-config.JACKPOT_TAIL)*uint(time.Second)) {
			c.countingTime += config.JACKPOT_EXTRA_TIME
			c.timer.Stop()
			c.timer = time.NewTimer(time.Until(c.lastUpdated.Add(time.Duration(c.countingTime) * time.Second)))
			go func() {
				<-c.timer.C
				c.end()
			}()
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(c.Type),
				EventType: "resetTime",
				Payload: gin.H{
					"countingTime": c.countingTime,
				}})
			c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
			log.LogMessage("jackpot controller", "betting delayed", "info", logrus.Fields{"round": c.roundID, "countingTime": c.countingTime})
		}
	}
}

func (c *Controller) determineWinner() (utils.PickWinnerResult[uint], error) {
	signedString, err := utils.GenerateRandomString(*c.ticketID)
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Failed to generate a random string"}})
		c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
		log.LogMessage("jackpot controller", "failed to generate a random string", "error", logrus.Fields{"round": c.roundID, "ticket": *c.ticketID})
		return utils.PickWinnerResult[uint]{}, errors.New("failed to generate a random string")
	}

	candidates := utils.WinnerCandidates[uint]{}
	c.playerToBets.Range(func(key, value any) bool {
		playerBets := value.(Bet4Player)
		playerID := key.(uint)
		candidates = append(candidates, utils.WinnerCandidate[uint]{ID: playerID, Entity: playerID, Weight: uint64(playerBets.TotalUsdAmount + playerBets.TotalNftAmount)})
		return true
	})
	winnedCandidate := utils.GenerateWinnerWithArray(signedString, candidates, 50)
	winnerId := winnedCandidate.Winner

	c.signedString = &signedString
	c.winnerID = &winnerId
	return winnedCandidate, nil
}

func (c *Controller) updateDbAtRoundEnd() (int64, []types.NftDetails, int64, int64, []types.NftDetails, int64, int64, error) {
	var round models.JackpotRound
	db := db.GetDB()
	if err := db.Preload("Players.Bets.Nfts").First(&round, c.roundID).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "can not find started round in DB "}})
		c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
		log.LogMessage("jackpot controller", "failed to get started Round from DB", "error", logrus.Fields{"round": c.roundID})
		return 0, []types.NftDetails{}, 0, 0, []types.NftDetails{}, 0, 0, errors.New("failed to get started round from db")
	}
	usdProfit, nfts4Profit, totalProfit, usdFee, nfts4Fee, totalFee, totalAmount := c.calculateJackpots()
	round.EndedAt = time.Now()
	round.SignedString = c.signedString
	round.WinnerID = *c.winnerID
	round.ChargedFee = usdFee
	db.Save(&round)

	var mintAddresses4Fee = []string{}
	for index := range nfts4Fee {
		mintAddresses4Fee = append(mintAddresses4Fee, nfts4Fee[index].MintAddress)
	}

	fee := usdFee
	_, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.JACKPOT_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&config.JACKPOT_FEE_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &fee,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&mintAddresses4Fee),
		},
		Type:          models.TxJackpotFee,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionJackpotReferenced,
	})
	if err != nil {
		log.LogMessage("jackpot controller", "failed to transfer fee", "error", logrus.Fields{"round": c.roundID, "error": err.Error()})
	}

	return usdProfit, nfts4Profit, totalProfit, usdFee, nfts4Fee, totalFee, totalAmount, nil
}

func (c *Controller) calculateJackpots() (int64, []types.NftDetails, int64, int64, []types.NftDetails, int64, int64) {
	var prize, usdProfit, totalAmount int64
	nfts := []types.NftDetails{}
	c.playerToBets.Range(func(key, value any) bool {
		playerID := key.(uint)
		playerBets := value.(Bet4Player)
		if playerID != *c.winnerID {
			prize += playerBets.TotalUsdAmount + playerBets.TotalNftAmount
			usdProfit += playerBets.TotalUsdAmount
		}
		for i := 0; i < len(playerBets.Bets); i++ {
			nfts = append(nfts, playerBets.Bets[i].Nfts...)
		}
		totalAmount += playerBets.TotalUsdAmount + playerBets.TotalNftAmount
		return true
	})
	w, _ := c.playerToBets.Load(*c.winnerID)
	winnerBets := w.(Bet4Player)
	totalFee := prize * c.fee / 100
	usdFee := totalFee
	if totalFee > usdProfit+winnerBets.TotalUsdAmount {
		usdFee = usdProfit + winnerBets.TotalUsdAmount
		usdProfit = 0
	} else {
		usdProfit = usdProfit - totalFee + winnerBets.TotalUsdAmount
	}

	nfts4Fee, nfts4Profit := utils.DetermineNFTs4Fee(nfts, totalFee-usdFee)
	totalFee = usdFee
	for i := 0; i < len(nfts4Fee); i++ {
		totalFee += nfts4Fee[i].Price
	}
	totalProfit := prize - totalFee
	return usdProfit, nfts4Profit, totalProfit, usdFee, nfts4Fee, totalFee, totalAmount
}

func (c *Controller) chargeNFTsAsFee(nfts []types.NftDetails) {
	db := db.GetDB()
	isFee := make(map[string]bool)
	for i := 0; i < len(nfts); i++ {
		isFee[nfts[i].MintAddress] = true
	}
	for i := 0; i < len(c.nftsInGame); i++ {
		if isFee[c.nftsInGame[i].MintAddress] {
			c.nftsInGame[i].Status = models.ChargedAsFee
		}
	}
	db.Save(&c.nftsInGame)
}

func (c *Controller) prizeWinner(usdPrize int64, nfts4Profit []types.NftDetails) {
	b, _ := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     usdPrize,
			BalanceType: models.ChipBalanceForGame,
			Nfts:        nfts4Profit,
			Delay:       15,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{*c.winnerID}, Message: b}

	var mintAddresses = []string{}
	for index := range nfts4Profit {
		mintAddresses = append(mintAddresses, nfts4Profit[index].MintAddress)
	}
	_, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.JACKPOT_TEMP_ID),
		ToUser:   (*db_aggregator.User)(c.winnerID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &usdPrize,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&mintAddresses),
		},
		Type:          models.TxJackpotProfit,
		ToBeConfirmed: true,
		OwnerID:       c.roundID,
		OwnerType:     models.TransactionJackpotReferenced,
	})
	if err != nil {
		log.LogMessage("jackpot controller", "failed to transfer profit to winner", "error", logrus.Fields{"round": c.roundID, "error": err.Error()})
		return
	}
}

func (c *Controller) updateStatistics(profit int64) {
	params := wager.PerformAfterWagerParams{
		Players: []wager.PlayerInPerformAfterWagerParams{},
		Type:    models.Jackpot,
	}

	c.playerToBets.Range(func(key, value any) bool {
		userID := key.(uint)
		playerBets := value.(Bet4Player)
		playerInAfterWager := wager.PlayerInPerformAfterWagerParams{
			UserID: userID,
			Bet:    playerBets.TotalUsdAmount + playerBets.TotalNftAmount,
		}
		if userID == *c.winnerID {
			playerInAfterWager.Profit = profit
		}

		params.Players = append(
			params.Players,
			playerInAfterWager,
		)
		return true
	})

	if err := wager.AfterWager(params); err != nil {
		log.LogMessage(
			"jackpot_internal_update_statistics",
			"failed to perform after wager",
			"error",
			logrus.Fields{
				"error":  err.Error(),
				"params": params,
			},
		)
	}
}

func (c *Controller) emitRoundEndEvent(winnedCandidate utils.PickWinnerResult[uint]) {
	db := db.GetDB()
	var resultCandidates []types.User
	var winner models.User
	db.First(&winner, *c.winnerID)
	for i := 0; i < len(winnedCandidate.CandidatesWithCount); i++ {
		userID := winnedCandidate.CandidatesWithCount[i].Entity
		var userData models.User
		db.First(&userData, userID)
		resultCandidates = append(resultCandidates, utils.GetUserDataWithPermissions(userData, nil, uint(winnedCandidate.CandidatesWithCount[i].Count)))
	}

	c.rollingDuration = uint64(15*1000) - uint64(time.Now().UnixMilli()-c.lastUpdated.UnixMilli())
	c.candidates = resultCandidates

	w, _ := c.playerToBets.Load(*c.winnerID)
	winnerBets := w.(Bet4Player)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
		EventType: string(c.status),
		Payload: types.JackpotEndPayload{
			RoundID:         c.roundID,
			TicketID:        *c.ticketID,
			Winner:          utils.GetUserDataWithPermissions(winner, nil, 0),
			UsdProfit:       c.usdProfit,
			NftProfit:       c.nfts4Profit,
			UsdFee:          c.usdFee,
			NftFee:          c.nfts4Fee,
			Prize:           c.totalAmount - c.totalFee,
			Chance:          float32(float64(winnerBets.TotalUsdAmount+winnerBets.TotalNftAmount) / float64(c.totalAmount) * 100),
			Candidates:      resultCandidates,
			RollingDuration: c.rollingDuration,
		}})
	c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
	log.LogMessage("jackpot controller", "round ended", "info", logrus.Fields{"round": c.roundID, "winner": map[string]any{"id": winner.ID, "name": winner.Name}})
}

func (c *Controller) emitErrMessageWithBalanceUpdate(message string, userID uint, betData BetData) {
	if message != "" {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(c.Type),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: message},
		})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
	}
	{
		b, _ := json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     betData.Amount,
				BalanceType: models.ChipBalanceForGame,
				Nfts:        betData.Nfts,
				Delay:       0,
			},
		})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
	}
}
