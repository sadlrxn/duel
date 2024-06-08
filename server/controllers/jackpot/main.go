package jackpot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

func (c *Controller) initLastRound() (bool, error) {
	db := db.GetDB()
	var round models.JackpotRound
	if result := db.Preload("Players.Bets.Nfts").Where("type = ?", c.Type).Last(&round); result.Error != nil {
		log.LogMessage(
			"Jackpot Controller",
			"Failed to get last round while initializing",
			"error",
			logrus.Fields{"error": result.Error.Error()},
		)
		return false, result.Error
	}

	if !round.EndedAt.IsZero() {
		return false, nil
	}

	log.LogMessage("Jackpot Controller", "Found an old round", "info", logrus.Fields{"roundId": round.ID})

	// Information needs to be loaded
	// - roundID
	// - playerToBets
	// - betPlayers
	// - nftsInGame
	// - totalPlayers
	// - ticketID

	for _, player := range round.Players {
		c.betPlayers = append(c.betPlayers, player.UserID)

		user := models.User{}
		if result := db.First(&user, player.UserID); result.Error != nil {
			log.LogMessage(
				"Jackpot Controller",
				"failed to get user's info",
				"error",
				logrus.Fields{
					"userId": player.UserID,
				},
			)
			return false, errors.New("failed to initialize the last round")
		}

		userInfo := utils.GetUserDataWithPermissions(user, nil, 0)
		bet4Player := Bet4Player{
			UserID:         player.UserID,
			UserName:       userInfo.Name,
			UserAvatar:     userInfo.Avatar,
			TotalUsdAmount: 0,
			TotalNftAmount: 0,
			Bets:           []BetData{},
		}

		for _, bet := range player.Bets {
			betData := BetData{
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

	c.totalPlayers = uint(len(round.Players))

	c.create(round.ID, round.TicketID)

	if len(c.betPlayers) > 1 {
		c.start(c.betPlayers[1])
	}

	return true, nil
}

func (c *Controller) Init(minBetAmount int64, maxBetAmount int64, betCountLimit uint, playerLimit uint, countingTime uint, rollingTime uint, fee int64) {
	c.minBetAmount = minBetAmount
	c.maxBetAmount = maxBetAmount
	c.betCountLimit = betCountLimit
	c.playerLimit = playerLimit
	c.countingTime = countingTime
	c.rollingTime = rollingTime
	c.fee = fee
	c.status = Available
	c.playerToBets = syncmap.Map{}
	c.betPlayers = []uint{}
	c.nftsInGame = []models.NftInGame{}
	c.lockUser = syncmap.Map{}

	if ok, err := c.initLastRound(); !(err == nil && ok) {
		c.setAvailable()
	}
}

func (c *Controller) ServeRoundData(conn *websocket.Conn) {
	players := []types.PlayerInJackpotRound{}
	var winnerInfo models.User
	var winner types.User
	db := db.GetDB()

	if c.winnerID != nil {
		db.First(&winnerInfo, c.winnerID)
		winner = utils.GetUserDataWithPermissions(winnerInfo, nil, 0)
	}
	for i := 0; i < int(c.totalPlayers); i++ {
		betUserId := c.betPlayers[i]
		nfts := []types.NftDetails{}
		betPerUser, _ := c.playerToBets.Load(betUserId)
		for j := 0; j < len(betPerUser.(Bet4Player).Bets); j++ {
			nfts = append(nfts, betPerUser.(Bet4Player).Bets[j].Nfts...)
		}
		players = append(players, types.PlayerInJackpotRound{
			ID:        betPerUser.(Bet4Player).UserID,
			Name:      betPerUser.(Bet4Player).UserName,
			Avatar:    betPerUser.(Bet4Player).UserAvatar,
			UsdAmount: betPerUser.(Bet4Player).TotalUsdAmount,
			NftAmount: betPerUser.(Bet4Player).TotalNftAmount,
			Nfts:      nfts,
			BetCount:  uint(len(betPerUser.(Bet4Player).Bets)),
		})
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
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
			CountingTime:    c.countingTime,
		}})
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}
}

func (c *Controller) Bet(userID uint, betData BetData) {
	if _, prs := c.lockUser.Load(userID); prs {
		c.emitErrMessageWithBalanceUpdate("Please wait for a moment.", userID, betData)
		return
	}

	c.lockUser.Store(userID, true)
	defer c.lockUser.Delete(userID)

	if c.status != Created && c.status != Started {
		c.emitErrMessageWithBalanceUpdate("Please wait for the next round.", userID, betData)
		return
	}
	playerBets, prs := c.playerToBets.Load(userID)
	finalAmount := betData.Amount + betData.NftAmount
	if prs {
		finalAmount += playerBets.(Bet4Player).TotalUsdAmount + playerBets.(Bet4Player).TotalNftAmount
	}

	if finalAmount > c.maxBetAmount || finalAmount < c.minBetAmount {
		c.emitErrMessageWithBalanceUpdate("Exceed bet amount limit.", userID, betData)
		return
	}
	if utils.IsDuplicateInArray(betData.Nfts) {
		c.emitErrMessageWithBalanceUpdate("Duplicated NFTs.", userID, betData)
		return
	}

	if !prs && c.totalPlayers == c.playerLimit {
		c.emitErrMessageWithBalanceUpdate("Exceed total player limit.", userID, betData)
		return
	} else if prs && len(playerBets.(Bet4Player).Bets) == int(c.betCountLimit) {
		c.emitErrMessageWithBalanceUpdate("Exceed Player Bet Count Limit.", userID, betData)
		return
	}

	var user models.User
	db := db.GetDB()
	db.First(&user, userID)
	userInfo := utils.GetUserDataWithPermissions(user, nil, 0)

	nftMintAddresses := []string{}
	for i := 0; i < len(betData.Nfts); i++ {
		nftMintAddresses = append(nftMintAddresses, betData.Nfts[i].MintAddress)
	}

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userID),
		ToUser:   (*db_aggregator.User)(&config.JACKPOT_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &betData.Amount,
			NftBalance:  db_aggregator.ConvertStringArrayToNftArray(&nftMintAddresses)},
		Type:          models.TxJackpotBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		c.emitErrMessageWithBalanceUpdate("", userID, betData)
		return
	}

	c.checkRemainingTime()

	nfts := []models.DepositedNft{}
	db.Where("mint_address IN ?", nftMintAddresses).Find(&nfts)

	var round models.JackpotRound
	var player models.JackpotPlayer

	if c.status == Created {
		if err := db.First(&round, c.roundID).Error; err != nil {
			c.emitErrMessageWithBalanceUpdate("Failed to get round data from DB", userID, betData)
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     round.ID,
				OwnerType:   models.TransactionJackpotReferenced,
			})
			log.LogMessage("jackpot controller", "failed to get round data from db", "error", logrus.Fields{"round": c.roundID, "ticket": *c.ticketID})
			return
		}
		if !prs {
			player = models.JackpotPlayer{
				UserID:  userID,
				RoundID: round.ID,
			}
			db.Create(&player)
			if c.totalPlayers > 0 {
				round.StartedAt = time.Now()
				db.Save(&round)
				c.start(userID)
			}
			c.totalPlayers++
		} else {
			if err := db.Model(&models.JackpotPlayer{}).Where("round_id = ?", round.ID).Where("user_id = ?", userID).First(&player).Error; err != nil {
				c.emitErrMessageWithBalanceUpdate("Failed to get player data from DB", userID, betData)
				transaction.Decline(transaction.DeclineRequest{
					Transaction: *tx,
					OwnerID:     round.ID,
					OwnerType:   models.TransactionJackpotReferenced,
				})
				log.LogMessage("jackpot controller", "failed to get Player data from db", "error", logrus.Fields{"round": c.roundID, "user": userID})
				return
			}
		}
	} else if c.status == Started {
		if err := db.First(&round, c.roundID).Error; err != nil {
			c.emitErrMessageWithBalanceUpdate("Can not find created round in DB", userID, betData)
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     round.ID,
				OwnerType:   models.TransactionJackpotReferenced,
			})
			return
		}
		if !prs {
			player = models.JackpotPlayer{
				UserID:  userID,
				RoundID: round.ID,
			}
			db.Create(&player)
			c.totalPlayers++
		} else {
			if err := db.Where("user_id = ?", userID).Where("round_id = ?", c.roundID).First(&player).Error; err != nil {
				c.emitErrMessageWithBalanceUpdate("Can not find player in DB", userID, betData)
				transaction.Decline(transaction.DeclineRequest{
					Transaction: *tx,
					OwnerID:     round.ID,
					OwnerType:   models.TransactionJackpotReferenced,
				})
				return
			}
		}
	} else {
		if err := db.Where("user_id = ?", userID).Where("round_id = ?", c.roundID).First(&player).Error; err != nil {
			c.emitErrMessageWithBalanceUpdate("Round already ended", userID, betData)
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     round.ID,
				OwnerType:   models.TransactionJackpotReferenced,
			})
			return
		}
	}

	var bet = models.JackpotBet{
		PlayerID:  player.ID,
		UsdAmount: betData.Amount,
	}
	db.Create(&bet)

	collectionID2Info := make(map[uint]*models.NftCollection)
	totalNftPrice := int64(0)
	var nftsInGame []models.NftInGame
	for _, nft := range nfts {
		if _, prs = collectionID2Info[nft.CollectionID]; !prs {
			var collectionInfo models.NftCollection
			if err := db.First(&collectionInfo, nft.CollectionID).Error; err != nil {
				c.emitErrMessageWithBalanceUpdate("Can not find NFT", userID, betData)
				transaction.Decline(transaction.DeclineRequest{
					Transaction: *tx,
					OwnerID:     round.ID,
					OwnerType:   models.TransactionJackpotReferenced,
				})
				return
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

	p, prs := c.playerToBets.Load(userID)
	if !prs {
		var betPerUser = Bet4Player{
			UserID:         userInfo.ID,
			UserName:       userInfo.Name,
			UserAvatar:     userInfo.Avatar,
			TotalUsdAmount: betData.Amount,
			TotalNftAmount: totalNftPrice,
			Bets:           []BetData{betData},
		}
		c.playerToBets.Store(userID, betPerUser)
		c.betPlayers = append(c.betPlayers, userID)
	} else {
		playerBet := p.(Bet4Player)
		playerBet.TotalUsdAmount += betData.Amount
		playerBet.TotalNftAmount += totalNftPrice
		playerBet.Bets = append(playerBet.Bets, betData)
		c.playerToBets.Store(userID, playerBet)
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(c.Type),
		EventType: "bet",
		Payload: types.JackpotBetPayload{
			RoundID:   c.roundID,
			ID:        userInfo.ID,
			Name:      userInfo.Name,
			Avatar:    userInfo.Avatar,
			UsdAmount: betData.Amount,
			NftAmount: totalNftPrice,
			Nfts:      betData.Nfts}})
	c.EventEmitter <- types.WSEvent{Room: c.Room, Message: b}
	transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionJackpotReferenced,
	})

	log.LogMessage("jackpot controller", "betted", "success", logrus.Fields{"round": c.roundID, "user": userID, "bet": betData})
}
