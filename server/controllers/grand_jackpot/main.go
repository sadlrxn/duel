package grand_jackpot

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/jackpot"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/controllers/wager"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

type DesiredTimes struct {
	start  time.Time
	finish time.Time
	count  time.Time
	end    time.Time
}

// Grand JP status
type Status string

const (
	Started  Status = "started"
	Finished Status = "finished"
	Counting Status = "counting"
	Rolling  Status = "rolling"
	Ended    Status = "ended"
)

type Controller struct {
	status          Status
	minBetAmount    int64
	countingTime    uint
	rollingTime     uint
	fee             int64
	roundID         uint
	playerToBets    syncmap.Map
	betPlayers      []uint
	nftsInGame      []models.NftInGame
	winnerID        *uint
	ticketID        *string
	signedString    *string
	lastUpdated     time.Time
	EventEmitter    chan types.WSEvent
	usdProfit       int64
	nfts4Profit     []types.NftDetails
	usdFee          int64
	nfts4Fee        []types.NftDetails
	totalProfit     int64
	totalAmount     int64
	totalFee        int64
	rollingDuration uint64
	candidates      []types.User
	desiredTimes    DesiredTimes
	lockUser        syncmap.Map
}

func (c *Controller) initLastRound() (bool, error) {
	db := db.GetDB()
	var round models.JackpotRound
	if result := db.Preload("Players.Bets.Nfts").Where("type = ?", models.Grand).Last(&round); result.Error != nil {
		log.LogMessage(
			"Grand Jackpot Controller",
			"Failed to get last round while initializing",
			"error",
			logrus.Fields{"error": result.Error.Error()},
		)
		return false, nil
	}

	if !round.EndedAt.IsZero() {
		c.lastUpdated = round.EndedAt
		return false, nil
	}
	log.LogMessage("Grand Jackpot Controller", "Found an old round", "info", logrus.Fields{"roundId": round.ID})

	for _, player := range round.Players {
		c.betPlayers = append(c.betPlayers, player.UserID)

		userInfo := models.User{}
		if result := db.First(&userInfo, player.UserID); result.Error != nil {
			log.LogMessage(
				"Grand Jackpot Controller",
				"failed to get user's info",
				"error",
				logrus.Fields{
					"userId": player.UserID,
				},
			)
			return false, errors.New("failed to initialize the last round")
		}
		userData := utils.GetUserDataWithPermissions(userInfo, nil, 0)

		bet4Player := jackpot.Bet4Player{
			UserID:         userData.ID,
			UserName:       userData.Name,
			UserAvatar:     userData.Avatar,
			UserRole:       userData.Role,
			TotalUsdAmount: 0,
			TotalNftAmount: 0,
			Bets:           []jackpot.BetData{},
		}

		for _, bet := range player.Bets {
			betData := jackpot.BetData{
				Amount:    bet.UsdAmount,
				Nfts:      []types.NftDetails{},
				NftAmount: 0,
				Time:      bet.CreatedAt,
			}

			for _, nft := range bet.Nfts {
				betData.NftAmount += nft.Price
				betData.Nfts = append(
					betData.Nfts,
					types.NftDetails{
						Name:            nft.Name,
						MintAddress:     nft.MintAddress,
						Image:           nft.Image,
						CollectionName:  nft.CollectionName,
						CollectionImage: nft.CollectionImage,
						Price:           nft.Price,
					},
				)
			}

			c.nftsInGame = append(c.nftsInGame, bet.Nfts...)

			bet4Player.TotalUsdAmount += betData.Amount
			bet4Player.TotalNftAmount += betData.NftAmount
			bet4Player.Bets = append(bet4Player.Bets, betData)
		}
		c.playerToBets.Store(player.UserID, bet4Player)
	}

	c.ticketID = &round.TicketID
	c.roundID = round.ID
	c.winnerID = nil
	c.signedString = nil
	c.usdProfit = 0
	c.nfts4Profit = []types.NftDetails{}
	c.usdFee = 0
	c.nfts4Fee = []types.NftDetails{}
	c.totalProfit = 0
	c.totalAmount = 0
	c.totalFee = 0
	c.rollingDuration = 0
	c.candidates = []types.User{}

	if !round.CountingStartedAt.IsZero() {
		c.status = Counting
		c.lastUpdated = round.CountingStartedAt
	} else {
		c.status = Started
		c.lastUpdated = round.StartedAt
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID, TicketID: c.ticketID},
	})
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
	log.LogMessage("grand jackpot controller", "round "+string(c.status), "info", logrus.Fields{"round": c.roundID, "ticketId": c.ticketID})
	return true, nil
}

// @External
// Initializer
func (c *Controller) Init(_minBetAmount int64, _countingTime uint, _rollingTime uint, _fee int64) {
	c.status = Ended
	c.lastUpdated = time.Now()
	c.minBetAmount = _minBetAmount
	c.countingTime = _countingTime
	c.rollingTime = _rollingTime
	c.fee = _fee
	c.playerToBets = syncmap.Map{}
	c.lockUser = syncmap.Map{}

	time.Local = time.UTC

	if ok, err := c.initLastRound(); !(err == nil && ok) {
		c.desiredTimes.start = config.GetServerConfig().NextGrandJackpotStartAt
		// c.desiredTimes.start = c.lastUpdated.Add(time.Duration(c.rollingTime) * time.Second)

		// if time.Now().Before(c.desiredTimes.start) {
		// 	timer := time.NewTimer(time.Until(c.desiredTimes.start))
		// 	go func() {
		// 		<-timer.C
		// 		c.start()
		// 	}()
		// } else {
		// 	for time.Now().After(c.desiredTimes.start) {
		// 		c.desiredTimes.start = c.desiredTimes.start.Add(time.Duration(c.rollingTime+config.GRAND_JACKPOT_BETTING_TIME) * time.Second)
		// 	}
		// 	timer := time.NewTimer(time.Until(c.desiredTimes.start))
		// 	go func() {
		// 		<-timer.C
		// 		c.start()
		// 	}()
		// }

		// for time.Now().After(c.desiredTimes.start) {
		// 	c.desiredTimes.start = c.desiredTimes.start.Add(time.Duration(c.rollingTime+config.GRAND_JACKPOT_BETTING_TIME) * time.Second)
		// }
		fmt.Println("START AT", c.desiredTimes.start)
		if time.Now().Before(c.desiredTimes.start) {
			timer := time.NewTimer(time.Until(c.desiredTimes.start))
			go func() {
				<-timer.C
				c.start()
			}()
		}
	} else {
		c.desiredTimes.end = c.lastUpdated.Add(time.Duration(config.GRAND_JACKPOT_BETTING_TIME) * time.Second)
		if time.Now().Before(c.desiredTimes.end) {
			timer := time.NewTimer(time.Until(c.desiredTimes.end))
			go func() {
				<-timer.C
				if err := c.end(); err != nil {
					log.LogMessage("grand jackpot controller", "error occured on round end", "error", logrus.Fields{"error": err.Error})
				}
			}()
		}
	}
}

// @External
// Send current round data via websocket connection
func (c *Controller) ServeRoundData(conn *websocket.Conn) {
	players := []types.PlayerInJackpotRound{}
	db := db.GetDB()
	var winnerInfo models.User
	var winner types.User

	if c.winnerID != nil {
		db.First(&winnerInfo, c.winnerID)
		winner = utils.GetUserDataWithPermissions(winnerInfo, nil, 0)
	}

	for _, userID := range c.betPlayers {
		var nfts []types.NftDetails
		b, _ := c.playerToBets.Load(userID)
		betPerUser := b.(jackpot.Bet4Player)
		for _, bet := range betPerUser.Bets {
			nfts = append(nfts, bet.Nfts...)
		}
		players = append(players, types.PlayerInJackpotRound{
			ID:        betPerUser.UserID,
			Role:      betPerUser.UserRole,
			Name:      betPerUser.UserName,
			Avatar:    betPerUser.UserAvatar,
			UsdAmount: betPerUser.TotalUsdAmount,
			NftAmount: betPerUser.TotalNftAmount,
			Nfts:      nfts,
			BetCount:  uint(len(betPerUser.Bets)),
		})
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: "game_data",
		Payload: types.JackpotPayload{
			Status:          string(c.status),
			RoundID:         c.roundID,
			TicketID:        c.ticketID,
			SignedString:    c.signedString,
			Players:         players,
			Winner:          winner,
			Offset:          uint64(time.Now().UnixMilli()) - uint64(c.lastUpdated.UnixMilli()),
			Prize:           c.totalAmount - c.totalFee,
			UsdProfit:       c.usdProfit,
			NftProfit:       c.nfts4Profit,
			UsdFee:          c.usdFee,
			NftFee:          c.nfts4Fee,
			Candidates:      c.candidates,
			RollingDuration: c.rollingDuration,
		},
	})
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}
}

// @External
// Grand JP bet handler
func (c *Controller) Bet(userID uint, betData jackpot.BetData) error {
	if _, prs := c.lockUser.Load(userID); prs {
		c.emitErrMessageWithBalanceUpdate("Please wait for a moment.", userID, betData)
		return errors.New("scam bet request")
	}

	c.lockUser.Store(userID, true)
	defer c.lockUser.Delete(userID)

	if err := c.validateStatusOnBet(userID, betData); err != nil {
		return err
	}
	if err := c.validateBetAmount(userID, betData); err != nil {
		return err
	}
	if err := c.validateBetNFTsDuplication(userID, betData); err != nil {
		return err
	}

	var mintAddresses []string
	for _, nft := range betData.Nfts {
		mintAddresses = append(mintAddresses, nft.MintAddress)
	}

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userID),
		ToUser:   (*db_aggregator.User)(&config.GRAND_JACKPOT_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &betData.Amount,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&mintAddresses),
		},
		Type:          models.TxGrandJackpotBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		c.emitErrMessageWithBalanceUpdate("", userID, betData)
		return err
	}

	db := db.GetDB()
	var nfts []models.DepositedNft
	if err := db.Where("mint_address IN ?", mintAddresses).Find(&nfts).Error; err != nil {
		c.emitErrMessageWithBalanceUpdate("Failed to get nfts from db", userID, betData)
		return err
	}

	var round models.JackpotRound
	var player models.JackpotPlayer
	if err := db.First(&round, c.roundID).Error; err != nil {
		c.emitErrMessageWithBalanceUpdate("Failed to get round from db", userID, betData)
		return err
	}

	{
		_, prs := c.playerToBets.Load(userID)
		if !prs {
			player = models.JackpotPlayer{
				UserID:  userID,
				RoundID: round.ID,
			}
			db.Create(&player)
			var user models.User
			db.First(&user, userID)
			c.candidates = append(c.candidates, utils.GetUserDataWithPermissions(user, nil, 0))
		} else {
			if err := db.Where("user_id = ?", userID).Where("round_id = ?", round.ID).First(&player).Error; err != nil {
				c.emitErrMessageWithBalanceUpdate("Failed to get player from db", userID, betData)
				transaction.Decline(transaction.DeclineRequest{
					Transaction: *tx,
					OwnerID:     round.ID,
					OwnerType:   models.TransactionJackpotReferenced,
				})
				return err
			}
		}
	}

	totalNftPrice, err := c.placeBet(player, userID, betData, nfts)
	if err != nil {
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     round.ID,
			OwnerType:   models.TransactionJackpotReferenced,
		})
		return err
	}

	c.finalizeAndEmitEvent(userID, betData, totalNftPrice)
	transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionJackpotReferenced,
	})

	log.LogMessage("grand jackpot controller", "betted", "success", logrus.Fields{"round": c.roundID, "user": userID, "bet": betData})
	return nil
}

// @Internal
// Place bet
func (c *Controller) placeBet(player models.JackpotPlayer, userID uint, betData jackpot.BetData, nfts []models.DepositedNft) (int64, error) {
	db := db.GetDB()
	var bet = models.JackpotBet{
		PlayerID:  player.ID,
		UsdAmount: betData.Amount,
	}
	db.Create(&bet)

	collectionID2Info := make(map[uint]*models.NftCollection)
	totalNftPrice := int64(0)
	var nftsInGame []models.NftInGame
	for _, nft := range nfts {
		if _, prs := collectionID2Info[nft.CollectionID]; !prs {
			var collectionInfo models.NftCollection
			if err := db.First(&collectionInfo, nft.CollectionID).Error; err != nil {
				c.emitErrMessageWithBalanceUpdate("Failed to ge NFT from db", userID, betData)
				return totalNftPrice, err
			}
			collectionID2Info[nft.CollectionID] = &collectionInfo
		}
		collectionInfo := collectionID2Info[nft.CollectionID]
		nftsInGame = append(nftsInGame, models.NftInGame{
			Name:            nft.Name,
			MintAddress:     nft.MintAddress,
			Image:           nft.Image,
			CollectionName:  collectionInfo.Name,
			CollectionImage: collectionInfo.Image,
			Price:           collectionInfo.FloorPrice,
			BetID:           bet.ID,
		})
		totalNftPrice += collectionInfo.FloorPrice
	}
	if len(nftsInGame) > 0 {
		db.Create(&nftsInGame)
		c.nftsInGame = append(c.nftsInGame, nftsInGame...)
	}
	return totalNftPrice, nil
}

// @Internal
// Finalize bet and emit event
func (c *Controller) finalizeAndEmitEvent(userID uint, betData jackpot.BetData, totalNftPrice int64) {
	db := db.GetDB()
	p, prs := c.playerToBets.Load(userID)
	var user models.User
	db.First(&user, userID)
	userInfo := utils.GetUserDataWithPermissions(user, nil, 0)
	if !prs {
		var betPerUser = jackpot.Bet4Player{
			UserID:         userInfo.ID,
			UserName:       userInfo.Name,
			UserAvatar:     userInfo.Avatar,
			UserRole:       userInfo.Role,
			TotalUsdAmount: betData.Amount,
			TotalNftAmount: totalNftPrice,
			Bets:           []jackpot.BetData{betData},
		}
		c.playerToBets.Store(userID, betPerUser)
		c.betPlayers = append(c.betPlayers, userID)
	} else {
		playerBets := p.(jackpot.Bet4Player)
		playerBets.TotalUsdAmount += betData.Amount
		playerBets.TotalNftAmount += totalNftPrice
		playerBets.Bets = append(playerBets.Bets, betData)
		c.playerToBets.Store(userID, playerBets)
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: "bet",
		Payload: types.JackpotBetPayload{
			RoundID:   c.roundID,
			ID:        userInfo.ID,
			Name:      userInfo.Name,
			Avatar:    userInfo.Avatar,
			Role:      userInfo.Role,
			UsdAmount: betData.Amount,
			NftAmount: totalNftPrice,
			Nfts:      betData.Nfts,
		},
	})
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
}

// @Internal
// Status validator on bet
func (c *Controller) validateStatusOnBet(userID uint, betData jackpot.BetData) error {
	if c.status != Started {
		c.emitErrMessageWithBalanceUpdate("Please wait for the next round.", userID, betData)
		return errors.New("invalid status")
	}
	return nil
}

// @Interanl
// Bet amount validator
func (c *Controller) validateBetAmount(userID uint, betData jackpot.BetData) error {
	if betData.Amount+betData.NftAmount < c.minBetAmount {
		c.emitErrMessageWithBalanceUpdate("Exceed bet amount limit", userID, betData)
		return errors.New("invalid bet amount")
	}
	return nil
}

// @Interanl
// Bet nft validator
func (c *Controller) validateBetNFTsDuplication(userID uint, betData jackpot.BetData) error {
	if utils.IsDuplicateInArray(betData.Nfts) {
		c.emitErrMessageWithBalanceUpdate("Duplicated NFTs", userID, betData)
		return errors.New("duplicated nfts")
	}
	return nil
}

// @Internal
// Emit an error message with a balance refund event
func (c *Controller) emitErrMessageWithBalanceUpdate(message string, userID uint, betData jackpot.BetData) {
	if message != "" {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.GrandJackpot),
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

// @Interanl
// Start a new round
func (c *Controller) start() error {
	if c.status != Ended {
		return errors.New("New round is not ready.")
	}

	c.status = Started
	currentTime := time.Now()

	c.desiredTimes.end = currentTime.Add(time.Duration(config.GRAND_JACKPOT_BETTING_TIME) * time.Second)
	if time.Now().Before(c.desiredTimes.end) {
		timer := time.NewTimer(time.Until(c.desiredTimes.end))
		go func() {
			<-timer.C
			if err := c.end(); err != nil {
				log.LogMessage("grand jackpot controller", "error occured on round end", "error", logrus.Fields{"error": err.Error})
			}
		}()
	}

	ticketID, err := utils.RequestTicketID()
	if err != nil {
		return err
	}

	db := db.GetDB()
	var round = models.JackpotRound{
		Players:   []models.JackpotPlayer{},
		TicketID:  ticketID,
		Type:      models.Grand,
		StartedAt: currentTime,
	}
	db.Create(&round)

	c.lastUpdated = round.StartedAt
	c.roundID = round.ID
	c.ticketID = &ticketID
	c.winnerID = nil
	c.playerToBets = syncmap.Map{}
	c.betPlayers = []uint{}
	c.nftsInGame = []models.NftInGame{}
	c.usdProfit = 0
	c.nfts4Profit = []types.NftDetails{}
	c.usdFee = 0
	c.nfts4Fee = []types.NftDetails{}
	c.totalProfit = 0
	c.totalAmount = 0
	c.totalFee = 0
	c.rollingDuration = 0
	c.candidates = []types.User{}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID, TicketID: c.ticketID},
	})
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
	log.LogMessage("grand jackpot controller", "round started", "info", logrus.Fields{"round": c.roundID, "ticketId": c.ticketID})
	return nil
}

// @Internal
// Finish & block betting
func (c *Controller) finishBetting() error {
	if c.status != Started {
		return errors.New("Invalid Status")
	}

	c.lastUpdated = time.Now()
	c.status = Finished

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID, TicketID: c.ticketID},
	})
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
	log.LogMessage("grand jackpot controller", "betting finished", "info", logrus.Fields{"round": c.roundID, "ticketId": c.ticketID})
	return nil
}

// @Internal
// Start down counting before rolling
func (c *Controller) startCountDown() error {
	if c.status != Finished {
		return errors.New("Invalid Status")
	}

	c.lastUpdated = time.Now()
	c.status = Counting

	db := db.GetDB()
	var round models.JackpotRound
	db.First(&round, c.roundID)

	round.CountingStartedAt = c.lastUpdated
	db.Save(&round)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
		EventType: string(c.status),
		Payload:   types.JackpotPayload{RoundID: c.roundID, TicketID: c.ticketID},
	})
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
	log.LogMessage("grand jackpot controller", "counting", "info", logrus.Fields{"round": c.roundID, "ticketId": c.ticketID})
	return nil
}

// @Interanl
// End current round and emit events
func (c *Controller) end() error {
	if c.status != Started {
		return errors.New("Cannot end round.")
	}

	c.lastUpdated = time.Now()
	c.status = Rolling

	// c.desiredTimes.start = c.lastUpdated.Add(time.Duration(c.rollingTime) * time.Second)
	// if time.Now().Before(c.desiredTimes.start) {
	// 	timer := time.NewTimer(time.Until(c.desiredTimes.start))
	// 	go func() {
	// 		<-timer.C
	// 		log.LogMessage("grand jackpot controller", "ended", "info", logrus.Fields{})
	// 		c.status = Ended
	// 		c.start()
	// 	}()
	// }

	if len(c.betPlayers) == 0 {
		var round models.JackpotRound
		db := db.GetDB()
		if err := db.Preload("Players.Bets.Nfts").First(&round, c.roundID).Error; err != nil {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.GrandJackpot),
				EventType: "message",
				Payload:   types.ErrorMessagePayload{Message: "Failed to get started round from db"},
			})
			c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
			log.LogMessage("grand jackpot controller", "failed to get started round from db", "error", logrus.Fields{"round": c.roundID})
			return errors.New("failed to get started round from db")
		}

		round.EndedAt = c.lastUpdated
		db.Save(&round)
		return errors.New("Empty round")
	}
	err := c.calculateSignedString()
	if err != nil {
		return err
	}
	winnedCandidate, err := c.determineWinner()
	if err != nil {
		return err
	}
	c.usdProfit, c.nfts4Profit, c.totalProfit, c.usdFee, c.nfts4Fee, c.totalFee, c.totalAmount = c.calculateJackpots()
	err = c.updateDbAtRoundEnd()
	if err != nil {
		return err
	}

	c.chargeNFTsAsFee(c.nfts4Fee)
	c.prizeWinner(c.usdProfit, c.nfts4Profit)
	c.updateStatistics(c.totalProfit)
	c.emitRoundEndEvent(winnedCandidate)
	return nil
}

// @Internal
// Calculate random string with current ticket
func (c *Controller) calculateSignedString() error {
	signedString, err := utils.GenerateRandomString(*c.ticketID)
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.GrandJackpot),
			EventType: "messge",
			Payload:   types.ErrorMessagePayload{Message: "Failed to generate a random string"},
		})
		c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
		log.LogMessage("grand jackpot controller", "failed to generate a random string", "error", logrus.Fields{"round": c.roundID, "ticket": *c.ticketID})
		return errors.New("failed to generate a random string")
	}
	c.signedString = &signedString
	return nil
}

// @Internal
// Determine the winner from random string
func (c *Controller) determineWinner() (utils.PickWinnerResult[uint], error) {
	candidates := utils.WinnerCandidates[uint]{}
	c.playerToBets.Range(func(key, value any) bool {
		playerBets := value.(jackpot.Bet4Player)
		if playerBets.UserRole == models.AdminRole {
			return true
		}
		candidates = append(candidates, utils.WinnerCandidate[uint]{ID: playerBets.UserID, Entity: playerBets.UserID, Weight: uint64(playerBets.TotalUsdAmount + playerBets.TotalNftAmount)})
		return true
	})

	winnedCandidate := utils.GenerateWinnerWithArray(*c.signedString, candidates, 50)
	winnerID := winnedCandidate.Winner

	c.winnerID = &winnerID
	return winnedCandidate, nil
}

// @Internal
// Store updated data to db when a round ends
func (c *Controller) updateDbAtRoundEnd() error {
	var round models.JackpotRound
	db := db.GetDB()
	if err := db.Preload("Players.Bets.Nfts").First(&round, c.roundID).Error; err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.GrandJackpot),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Failed to get started round from db"},
		})
		c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
		log.LogMessage("grand jackpot controller", "failed to get started round from db", "error", logrus.Fields{"round": c.roundID})
		return errors.New("failed to get started round from db")
	}

	round.EndedAt = time.Now()
	round.SignedString = c.signedString
	round.WinnerID = *c.winnerID
	round.ChargedFee = c.usdFee
	db.Save(&round)

	var mintAddresses4Fee []string
	for _, nft := range c.nfts4Fee {
		mintAddresses4Fee = append(mintAddresses4Fee, nft.MintAddress)
	}

	fee := c.usdFee
	_, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.GRAND_JACKPOT_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&config.GRAND_JACKPOT_FEE_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &fee,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&mintAddresses4Fee),
		},
		Type:          models.TxGrandJackpotFee,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionJackpotReferenced,
	})
	if err != nil {
		log.LogMessage("grand jackpot controller", "failed to transfer fee", "error", logrus.Fields{"round": c.roundID, "error": err.Error()})
	}
	return nil
}

// @Interanl
// Calculate values related to the round
func (c *Controller) calculateJackpots() (usdProfit int64, nfts4Profit []types.NftDetails, totalProfit int64, usdFee int64, nfts4Fee []types.NftDetails, totalFee int64, totalAmount int64) {
	var prize int64
	var nfts []types.NftDetails
	c.playerToBets.Range(func(key, value any) bool {
		playerBets := value.(jackpot.Bet4Player)
		if playerBets.UserID != *c.winnerID {
			prize += playerBets.TotalUsdAmount + playerBets.TotalNftAmount
			usdProfit += playerBets.TotalUsdAmount
		}
		for _, bet := range playerBets.Bets {
			nfts = append(nfts, bet.Nfts...)
		}
		totalAmount += playerBets.TotalUsdAmount + playerBets.TotalNftAmount
		return true
	})
	totalFee = prize * c.fee / 100
	usdFee = totalFee
	w, _ := c.playerToBets.Load(*c.winnerID)
	winnerBets := w.(jackpot.Bet4Player)
	if totalFee > usdProfit+winnerBets.TotalUsdAmount {
		usdFee = usdProfit + winnerBets.TotalUsdAmount
		usdProfit = 0
	} else {
		usdProfit = usdProfit - totalFee + winnerBets.TotalUsdAmount
	}

	nfts4Fee, nfts4Profit = utils.DetermineNFTs4Fee(nfts, totalFee-usdFee)
	totalFee = usdFee
	for _, nft := range nfts4Fee {
		totalFee += nft.Price
	}
	totalProfit = prize - totalFee
	return
}

// @Internal
// Transfer NFTs for fee to mixpanel wallet
func (c *Controller) chargeNFTsAsFee(nfts []types.NftDetails) {
	db := db.GetDB()
	isFee := make(map[string]bool)
	for _, nft := range nfts {
		isFee[nft.MintAddress] = true
	}
	for i, nftInGame := range c.nftsInGame {
		if isFee[nftInGame.MintAddress] {
			c.nftsInGame[i].Status = models.ChargedAsFee
		}
	}
	db.Save(&c.nftsInGame)
}

// @Interanl
// Transfer profit to winner
func (c *Controller) prizeWinner(usdProfit int64, nfts4Profit []types.NftDetails) {
	b, _ := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     usdProfit,
			BalanceType: models.ChipBalanceForGame,
			Nfts:        nfts4Profit,
			Delay:       float32(config.GRAND_JACKPOT_ROLLING_DURATION),
		},
	})
	c.EventEmitter <- types.WSEvent{Users: []uint{*c.winnerID}, Message: b}

	var mintaddresses []string
	for _, nft := range nfts4Profit {
		mintaddresses = append(mintaddresses, nft.MintAddress)
	}

	_, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.GRAND_JACKPOT_TEMP_ID),
		ToUser:   (*db_aggregator.User)(c.winnerID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &usdProfit,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&mintaddresses),
		},
		Type:          models.TxGrandJackpotProfit,
		ToBeConfirmed: true,
		OwnerID:       c.roundID,
		OwnerType:     models.TransactionJackpotReferenced,
	})
	if err != nil {
		log.LogMessage("grand jackpot controller", "failed to transfer profit to winner", "error", logrus.Fields{"round": c.roundID, "error": err.Error()})
	}
}

// @Internal
// Update statistics of all palyers betted in the round
func (c *Controller) updateStatistics(profit int64) {
	params := wager.PerformAfterWagerParams{
		Players: []wager.PlayerInPerformAfterWagerParams{},
		Type:    models.Jackpot,
	}

	c.playerToBets.Range(func(key, value any) bool {
		userID := key.(uint)
		playerBets := value.(jackpot.Bet4Player)
		playerInPerformAfterWager := wager.PlayerInPerformAfterWagerParams{
			UserID: userID,
			Bet:    playerBets.TotalUsdAmount + playerBets.TotalNftAmount,
		}
		if userID == *c.winnerID {
			playerInPerformAfterWager.Profit = profit
		}

		params.Players = append(
			params.Players,
			playerInPerformAfterWager,
		)
		return true
	})

	if err := wager.AfterWager(params); err != nil {
		log.LogMessage(
			"grand_jackpot_internal_update_statistics",
			"failed to perform after wager",
			"error",
			logrus.Fields{
				"error":  err.Error(),
				"params": params,
			},
		)
	}
}

// @Internal
// Emit round end event
func (c *Controller) emitRoundEndEvent(winnedCandidate utils.PickWinnerResult[uint]) {
	db := db.GetDB()
	var resultCandidates []types.User
	var winner models.User
	db.First(&winner, *c.winnerID)
	for _, candidate := range winnedCandidate.CandidatesWithCount {
		var userData models.User
		db.First(&userData, candidate.Entity)
		resultCandidates = append(resultCandidates, utils.GetUserDataWithPermissions(userData, nil, uint(candidate.Count)))
	}

	c.rollingDuration = uint64(config.GRAND_JACKPOT_ROLLING_DURATION*1000) - uint64(time.Now().UnixMilli()-c.lastUpdated.UnixMilli())
	c.candidates = resultCandidates

	w, _ := c.playerToBets.Load(*c.winnerID)
	winnerBets := w.(jackpot.Bet4Player)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.GrandJackpot),
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
	c.EventEmitter <- types.WSEvent{Room: types.GrandJackpot, Message: b}
	log.LogMessage("grand jackpot controller", "round ended", "info", logrus.Fields{"round": c.roundID, "winner": map[string]any{"id": winner.ID, "name": winner.Name}})
}
