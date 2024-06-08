package coinflip

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
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
)

func (c *Controller) ServeGameData(conn *websocket.Conn) {
	activeRoundPayloads := types.CoinflipRoundDataPayloads{}
	db := db.GetDB()
	c.activeRounds.Range(func(key, value interface{}) bool {
		var round models.CoinflipRound = value.(models.CoinflipRound)
		var tailsUser, headsUser, creator models.User
		db.First(&tailsUser, round.TailsUserID)
		db.First(&headsUser, round.HeadsUserID)
		var creatorID, ok = c.round2Creator.Load(round.ID)
		if !ok {
			return true
		}
		db.First(&creator, creatorID)

		activeRoundPayloads = append(activeRoundPayloads, types.CoinflipRoundDataPayload{
			RoundID:   round.ID,
			HeadsUser: utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser: utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:    round.Amount,
			Prize:     round.Prize,
			TicketID:  round.TicketID,
			CreatorID: creator.ID,
		})
		return true
	})

	sort.Sort(activeRoundPayloads)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "game_data",
		Payload:   activeRoundPayloads})
	c.EventEmitter <- types.WSEvent{Conns: []*websocket.Conn{conn}, Message: b}
}

func (c *Controller) Create(userID uint, eventParam EventParam) {
	if !c.validateAmount(userID, eventParam) {
		return
	}
	if !c.validateCount(userID, eventParam) {
		return
	}

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &eventParam.Amount,
		},
		Type:          models.TxCoinflipBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to cash in.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	var tailsUserID, headsUserID *uint
	var tailsUser, headsUser, userInfo models.User

	db := db.GetDB()
	if result := db.First(&userInfo, userID); result.Error != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to get user data.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}
	if eventParam.Side == string(models.Heads) {
		headsUserID = &userID
		if result := db.First(&headsUser, headsUserID); result.Error != nil {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Coinflip),
				EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to get user data.", RoundID: 0}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			b, _ = json.Marshal(types.WSMessage{
				EventType: "balance_update",
				Payload: types.BalanceUpdatePayload{
					UpdateType:  types.Increase,
					Balance:     eventParam.Amount,
					BalanceType: models.ChipBalanceForGame,
					Delay:       0,
				}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     userID,
				OwnerType:   models.TransactionUserReferenced,
			})
			return
		}
	} else if eventParam.Side == string(models.Tails) {
		tailsUserID = &userID
		if result := db.First(&tailsUser, tailsUserID); result.Error != nil {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Coinflip),
				EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to get user data.", RoundID: 0}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			b, _ = json.Marshal(types.WSMessage{
				EventType: "balance_update",
				Payload: types.BalanceUpdatePayload{
					UpdateType:  types.Increase,
					Balance:     eventParam.Amount,
					BalanceType: models.ChipBalanceForGame,
					Delay:       0,
				}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     userID,
				OwnerType:   models.TransactionUserReferenced,
			})
			return
		}
	} else {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Unavailable coin side.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}
	ticketID, err := c.validateTicketID(userID, eventParam.Amount, tx)
	if err != nil {
		return
	}

	round := models.CoinflipRound{
		TailsUserID: tailsUserID,
		HeadsUserID: headsUserID,
		Amount:      eventParam.Amount,
		Prize:       eventParam.Amount * 2 * (100 - c.fee) / 100,
		TicketID:    ticketID,
	}
	if result := db.Create(&round); result.Error != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to create a new game.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}

	c.activeRounds.Store(round.ID, round)
	c.round2Creator.Store(round.ID, userID)

	if err := transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionCoinflipReferenced,
	}); err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to cash in.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "created",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:   round.ID,
			HeadsUser: utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser: utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:    round.Amount,
			Prize:     round.Prize,
			TicketID:  round.TicketID,
			CreatorID: userInfo.ID,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}
	log.LogMessage("coinflip controller", "new round created", "success", logrus.Fields{"round": round.ID, "user": userID, "side": eventParam.Side, "amount": eventParam.Amount})
}

func (c *Controller) Join(userID uint, roundID uint) {
	db := db.GetDB()
	activeRound, prs := c.activeRounds.Load(roundID)
	c.isRoundPending.Store(roundID, true)
	defer c.isRoundPending.Delete(roundID)

	if !prs {
		var round models.CoinflipRound
		db.First(&round, roundID)
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Already ended round.", RoundID: roundID}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}
	round := activeRound.(models.CoinflipRound)

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &round.Amount,
		},
		Type:          models.TxCoinflipBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to cash in.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	if round.HeadsUserID == nil {
		if *round.TailsUserID == userID {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Coinflip),
				EventType: "message",
				Payload:   types.ErrorMessagePayload{Message: "Already betted", RoundID: roundID}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			b, _ = json.Marshal(types.WSMessage{
				EventType: "balance_update",
				Payload: types.BalanceUpdatePayload{
					UpdateType:  types.Increase,
					Balance:     round.Amount,
					BalanceType: models.ChipBalanceForGame,
					Delay:       0,
				}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     round.ID,
				OwnerType:   models.TransactionCoinflipReferenced,
			})
			return
		}
		round.HeadsUserID = &userID
	} else if round.TailsUserID == nil {
		if *round.HeadsUserID == userID {
			b, _ := json.Marshal(types.WSMessage{
				Room:      string(types.Coinflip),
				EventType: "message",
				Payload:   types.ErrorMessagePayload{Message: "Already betted", RoundID: roundID}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			b, _ = json.Marshal(types.WSMessage{
				EventType: "balance_update",
				Payload: types.BalanceUpdatePayload{
					UpdateType:  types.Increase,
					Balance:     round.Amount,
					BalanceType: models.ChipBalanceForGame,
					Delay:       0,
				}})
			c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     round.ID,
				OwnerType:   models.TransactionCoinflipReferenced,
			})
			return
		}
		round.TailsUserID = &userID
	} else {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "The round is full", RoundID: roundID}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     round.ID,
			OwnerType:   models.TransactionCoinflipReferenced,
		})
		return
	}

	signedString, err := c.validateRandomString(userID, round, tx)
	if err != nil {
		return
	}
	c.activeRounds.Store(roundID, round)
	candidates := utils.WinnerCandidates[uint]{{Entity: *round.HeadsUserID, Weight: uint64(round.Amount)}, {Entity: *round.TailsUserID, Weight: uint64(round.Amount)}}
	winnedCandidate := utils.GenerateWinnerWithArray(signedString, candidates, 2)
	winnerId := winnedCandidate.Winner

	round.SignedString = &signedString
	round.WinnerID = &winnerId
	round.EndedAt = time.Now()
	if result := db.Save(&round); result.Error != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to save game.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}
	transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionCoinflipReferenced,
	})

	_, err = transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		ToUser:   (*db_aggregator.User)(round.WinnerID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &round.Prize,
		},
		Type:          models.TxCoinflipProfit,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionCoinflipReferenced,
	})
	if err != nil {
		log.LogMessage("coinflip controller", "failed to transfer profit to winner", "error", logrus.Fields{"round": roundID, "error": err.Error()})
		return
	}

	fee := round.Amount*2 - round.Prize
	_, err = transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_FEE_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &fee,
		},
		Type:          models.TxCoinflipFee,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionCoinflipReferenced,
	})
	if err != nil {
		log.LogMessage("coinflip controller", "failed to transfer house fee", "error", logrus.Fields{"round": roundID, "error": err.Error()})
		return
	}

	afterWagerParams := wager.PerformAfterWagerParams{
		Players: []wager.PlayerInPerformAfterWagerParams{
			{
				UserID: *round.HeadsUserID,
				Bet:    round.Amount,
			},
			{
				UserID: *round.TailsUserID,
				Bet:    round.Amount,
			},
		},
		Type: models.Coinflip,
	}
	if winnerId == *round.HeadsUserID {
		afterWagerParams.Players[0].Profit = round.Prize - round.Amount
	} else {
		afterWagerParams.Players[1].Profit = round.Prize - round.Amount
	}
	if err := wager.AfterWager(afterWagerParams); err != nil {
		log.LogMessage(
			"coinflip_join",
			"failed to perform after wager",
			"error",
			logrus.Fields{
				"error":    err.Error(),
				"headUser": *round.HeadsUserID,
				"tailUser": *round.TailsUserID,
				"amount":   round.Amount,
				"winner":   winnerId,
				"profit":   round.Prize - round.Amount,
			},
		)
	}

	var tailsUser, headsUser, creator, winner models.User
	db.First(&tailsUser, round.TailsUserID)
	db.First(&headsUser, round.HeadsUserID)
	creatorID, _ := c.round2Creator.Load(roundID)
	db.First(&creator, creatorID)
	db.First(&winner, winnerId)

	c.activeRounds.Delete(roundID)
	c.round2Creator.Delete(roundID)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "joined",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:      round.ID,
			HeadsUser:    utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser:    utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:       round.Amount,
			Prize:        round.Prize,
			TicketID:     round.TicketID,
			SignedString: *round.SignedString,
			WinnerID:     winner.ID,
			CreatorID:    creator.ID,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}

	b, _ = json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     round.Prize,
			BalanceType: models.ChipBalanceForGame,
			Delay:       5.5,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{winnerId}, Message: b}
	log.LogMessage("coinflip controller", "joined", "success", logrus.Fields{"round": round.ID, "user": userID, "winner": map[string]any{"id": winner.ID, "name": winner.Name}})
}

func (c *Controller) Cancel(userID uint, roundID uint) {
	isPending, ok := c.isRoundPending.Load(roundID)
	if _, prs := c.activeRounds.Load(roundID); !prs || (ok && isPending.(bool)) {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Invalid Round", RoundID: roundID}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	creatorID, _ := c.round2Creator.Load(roundID)
	if creatorID.(uint) != userID {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Permission Denied", RoundID: roundID}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	c.activeRounds.Delete(roundID)
	c.round2Creator.Delete(roundID)

	var round models.CoinflipRound
	db := db.GetDB()
	db.First(&round, roundID)
	round.EndedAt = time.Now()
	db.Save(&round)

	_, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&userID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &round.Amount,
		},
		Type:          models.TxCoinflipCancel,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionCoinflipReferenced,
	})
	if err != nil {
		log.LogMessage("coinflip controller", "cancel failed", "error", logrus.Fields{"error": err.Error()})
		return
	}

	var creator models.User
	db.First(&creator, creatorID)

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "cancelled",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:   roundID,
			CreatorID: creator.ID,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}
	b, _ = json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     round.Amount,
			BalanceType: models.ChipBalanceForGame,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
	log.LogMessage("coinflip controller", "cancelled", "success", logrus.Fields{"round": round.ID, "user": userID})
}

func (c *Controller) BetAgainstBot(userID uint, eventParam EventParam) {
	if !c.validateAmount(userID, eventParam) {
		return
	}
	if !c.validateCount(userID, eventParam) {
		return
	}

	result, tx, _ := coupon.TryBet(coupon.TryBetWithCouponRequest{
		UserID:  userID,
		Balance: eventParam.Amount,
		Type:    models.CpTxCoinflipBet,
	})
	if result == coupon.CouponBetUnavailable && eventParam.PaidBalanceType == models.ChipBalanceForGame {
		c.betAgainstBotWithChips(userID, eventParam)
	} else if result == coupon.CouponBetSucceed && eventParam.PaidBalanceType == models.CouponBalanceForGame {
		c.betAgainstBotWithCoupon(userID, eventParam, tx)
	} else {
		if tx != 0 {
			coupon.Decline(tx)
		}
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				Wagered:     eventParam.Amount,
				BalanceType: eventParam.PaidBalanceType,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
	}
}

func (c *Controller) betAgainstBotWithChips(userID uint, eventParam EventParam) {
	tx1, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &eventParam.Amount,
		},
		Type:          models.TxCoinflipBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		log.LogMessage(
			"coinflip_main",
			"failed to transfer user's fund to coinflip temp",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to cash in.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		return
	}

	tx2, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_BOT_ID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &eventParam.Amount,
		},
		Type:          models.TxCoinflipBet,
		ToBeConfirmed: false,
	})
	if err != nil {
		log.LogMessage(
			"coinflip_main",
			"failed to transfer bot's fund to coinflip temp",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to cash in.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx1,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}

	var tailsUserID, headsUserID *uint
	var tailsUser, headsUser, userInfo models.User

	db := db.GetDB()
	db.First(&userInfo, userID)
	if eventParam.Side == string(models.Heads) {
		headsUserID = &userID
		tailsUserID = &config.COINFLIP_BOT_ID
		db.First(&headsUser, headsUserID)
		db.First(&tailsUser, tailsUserID)
	} else if eventParam.Side == string(models.Tails) {
		headsUserID = &config.COINFLIP_BOT_ID
		tailsUserID = &userID
		db.First(&headsUser, headsUserID)
		db.First(&tailsUser, tailsUserID)
	} else {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Invalid coin side.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx1,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx2,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}
	ticketID, err := c.validateTicketID(userID, eventParam.Amount, tx1)
	if err != nil {
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx2,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}

	round := models.CoinflipRound{
		TailsUserID: tailsUserID,
		HeadsUserID: headsUserID,
		Amount:      eventParam.Amount,
		Prize:       eventParam.Amount * 2 * (100 - c.fee) / 100,
		TicketID:    ticketID,
	}

	signedString, err := c.validateRandomString(userID, round, tx2)
	if err != nil {
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx1,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}

	candidates := utils.WinnerCandidates[uint]{{Entity: *round.HeadsUserID, Weight: uint64(round.Amount)}, {Entity: *round.TailsUserID, Weight: uint64(round.Amount)}}
	winnedCandidate := utils.GenerateWinnerWithArray(signedString, candidates, 2)
	winnerId := winnedCandidate.Winner

	round.SignedString = &signedString
	round.WinnerID = &winnerId
	round.EndedAt = time.Now()

	if result := db.Create(&round); result.Error != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to create a new round.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     eventParam.Amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx1,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx2,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return
	}

	transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx1,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionCoinflipReferenced,
	})
	transaction.Confirm(transaction.ConfirmRequest{
		Transaction: *tx2,
		OwnerID:     round.ID,
		OwnerType:   models.TransactionCoinflipReferenced,
	})

	_, err = transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		ToUser:   (*db_aggregator.User)(round.WinnerID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &round.Prize,
		},
		Type:          models.TxCoinflipProfit,
		ToBeConfirmed: true,
		OwnerID:       round.ID,
		OwnerType:     models.TransactionCoinflipReferenced,
	})
	if err != nil {
		log.LogMessage("coinflip controller", "failed to transfer profit to winner", "error", logrus.Fields{"round": round.ID, "error": err.Error()})
		return
	}

	refillBot := (*round.WinnerID != config.COINFLIP_BOT_ID)
	feeBackup := round.Amount*2 - round.Prize
	fee := feeBackup
	_, err = transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
		ToUser:   (*db_aggregator.User)(&config.COINFLIP_FEE_ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &fee,
		},
		Type:          models.TxCoinflipFee,
		ToBeConfirmed: true,
		HouseFeeMeta: &transaction.HouseFeeMeta{
			User:        db_aggregator.User(userID),
			WagerAmount: round.Amount,
		},
		OwnerID:   round.ID,
		OwnerType: models.TransactionCoinflipReferenced,
		EdgeDetails: transaction.EdgeDetailsInTransactionRequest{
			ShouldNotDistributeRevShare: refillBot,
		},
	})
	if err != nil {
		log.LogMessage("coinflip controller", "failed to transfer house fee", "error", logrus.Fields{"round": round.ID, "error": err.Error()})
		return
	}

	// Refill CF_TEMP wallet incase bot lost the game instead of giving
	// 80% of revshare to bot holders.
	if refillBot {
		refillFee := utils.RevShareFromFee(feeBackup)

		_, err = transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&config.COINFLIP_TEMP_ID),
			ToUser:   (*db_aggregator.User)(&config.COINFLIP_BOT_ID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &refillFee,
			},
			Type:          models.TxCoinflipRefill,
			ToBeConfirmed: true,
			OwnerID:       round.ID,
			OwnerType:     models.TransactionCoinflipReferenced,
		})
		if err != nil {
			log.LogMessage(
				"coinflip controller",
				"failed to perform refill transaction",
				"error",
				logrus.Fields{
					"round": round.ID,
					"error": err.Error(),
				})
			return
		}
	}

	c.round2Creator.Store(round.ID, userID)
	timer := time.NewTimer(time.Duration(5) * time.Second)
	go func() {
		<-timer.C
		c.round2Creator.Delete(round.ID)
	}()

	afterWagerParams := wager.PerformAfterWagerParams{
		Players: []wager.PlayerInPerformAfterWagerParams{
			{
				UserID: userID,
				Bet:    round.Amount,
			},
		},
		Type:        models.Coinflip,
		IsHouseGame: true,
	}
	if winnerId == userID {
		afterWagerParams.Players[0].Profit = round.Prize - round.Amount
	}
	if err := wager.AfterWager(afterWagerParams); err != nil {
		log.LogMessage(
			"coinflip_join_bet_against_bot_with_chips",
			"failed to perform after wager",
			"error",
			logrus.Fields{
				"error":    err.Error(),
				"headUser": *round.HeadsUserID,
				"tailUser": *round.TailsUserID,
				"amount":   round.Amount,
				"winner":   winnerId,
				"profit":   round.Prize - round.Amount,
			},
		)
	}

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "created",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:   round.ID,
			HeadsUser: utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser: utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:    round.Amount,
			Prize:     round.Prize,
			TicketID:  round.TicketID,
			CreatorID: userInfo.ID,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}

	b, _ = json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "joined",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:      round.ID,
			HeadsUser:    utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser:    utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:       round.Amount,
			Prize:        round.Prize,
			TicketID:     round.TicketID,
			SignedString: *round.SignedString,
			WinnerID:     *round.WinnerID,
			CreatorID:    userInfo.ID,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}

	if userInfo.ID == winnerId {
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Prize,
				BalanceType: models.ChipBalanceForGame,
				Delay:       5.5,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{winnerId}, Message: b}
	}

	log.LogMessage("coinflip controller", "play against bot", "success", logrus.Fields{"round": round.ID, "user": userID, "winner": winnerId})
}

func (c *Controller) betAgainstBotWithCoupon(userID uint, eventParam EventParam, txID uint) {
	var tailsUserID, headsUserID *uint
	var tailsUser, headsUser, userInfo models.User

	var couponAmount = eventParam.Amount

	db := db.GetDB()
	db.First(&userInfo, userID)
	if eventParam.Side == string(models.Heads) {
		headsUserID = &userID
		tailsUserID = &config.COINFLIP_BOT_ID
	} else if eventParam.Side == string(models.Tails) {
		headsUserID = &config.COINFLIP_BOT_ID
		tailsUserID = &userID
	} else {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Invalid coin side.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     couponAmount,
				Wagered:     couponAmount,
				BalanceType: models.CouponBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		if err := coupon.Decline(txID); err != nil {
			log.LogMessage("betAgainstBotWithCoupon", "failed to decline coupon bet", "error", logrus.Fields{"error": err.Error()})
		}
		return
	}
	db.First(&headsUser, headsUserID)
	db.First(&tailsUser, tailsUserID)

	ticketID, err := c.validateTicketID(userID, eventParam.Amount, nil)
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     couponAmount,
				Wagered:     couponAmount,
				BalanceType: models.CouponBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		if err := coupon.Decline(txID); err != nil {
			log.LogMessage("betAgainstBotWithCoupon", "failed to decline coupon bet", "error", logrus.Fields{"error": err.Error()})
		}
		return
	}

	round := models.CoinflipRound{
		TailsUserID:     tailsUserID,
		HeadsUserID:     headsUserID,
		Amount:          eventParam.Amount,
		Prize:           eventParam.Amount * 2 * (100 - c.fee) / 100,
		TicketID:        ticketID,
		PaidBalanceType: models.CouponBalanceForGame,
	}

	signedString, err := c.validateRandomString(userID, round, nil)
	if err != nil {
		b, _ := json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     couponAmount,
				Wagered:     couponAmount,
				BalanceType: models.CouponBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		if err := coupon.Decline(txID); err != nil {
			log.LogMessage("betAgainstBotWithCoupon", "failed to decline coupon bet", "error", logrus.Fields{"error": err.Error()})
		}
		return
	}

	candidates := utils.WinnerCandidates[uint]{{Entity: *round.HeadsUserID, Weight: uint64(round.Amount)}, {Entity: *round.TailsUserID, Weight: uint64(round.Amount)}}
	winnedCandidate := utils.GenerateWinnerWithArray(signedString, candidates, 2)
	winnerId := winnedCandidate.Winner

	round.SignedString = &signedString
	round.WinnerID = &winnerId
	round.EndedAt = time.Now()

	if result := db.Create(&round); result.Error != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message", Payload: types.ErrorMessagePayload{Message: "Failed to create a new round.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     couponAmount,
				Wagered:     couponAmount,
				BalanceType: models.CouponBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		if err := coupon.Decline(txID); err != nil {
			log.LogMessage("betAgainstBotWithCoupon", "failed to decline coupon bet", "error", logrus.Fields{"error": err.Error()})
		}
		return
	}

	coupon.Confirm(txID)

	if winnerId == userID {
		if _, err := coupon.Perform(coupon.CouponTransactionRequest{
			UserID:        userID,
			Balance:       round.Prize,
			Type:          models.CpTxCoinflipProfit,
			ToBeConfirmed: true,
		}); err != nil {
			log.LogMessage("coinflip controller", "failed to transfer profit to winner", "error", logrus.Fields{"round": round.ID, "error": err.Error()})
		}
	}

	c.round2Creator.Store(round.ID, userID)
	timer := time.NewTimer(time.Duration(5) * time.Second)
	go func() {
		<-timer.C
		c.round2Creator.Delete(round.ID)
	}()

	b, _ := json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "created",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:         round.ID,
			HeadsUser:       utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser:       utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:          round.Amount,
			Prize:           round.Prize,
			TicketID:        round.TicketID,
			CreatorID:       userInfo.ID,
			PaidBalanceType: models.CouponBalanceForGame,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}

	b, _ = json.Marshal(types.WSMessage{
		Room:      string(types.Coinflip),
		EventType: "joined",
		Payload: types.CoinflipRoundDataPayload{
			RoundID:         round.ID,
			HeadsUser:       utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser:       utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:          round.Amount,
			Prize:           round.Prize,
			TicketID:        round.TicketID,
			SignedString:    *round.SignedString,
			WinnerID:        *round.WinnerID,
			CreatorID:       userInfo.ID,
			PaidBalanceType: models.CouponBalanceForGame,
		}})
	c.EventEmitter <- types.WSEvent{Room: types.Coinflip, Message: b}

	if userInfo.ID == winnerId {
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     round.Prize,
				BalanceType: models.CouponBalanceForGame,
				Delay:       5.5,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{winnerId}, Message: b}
	}

	log.LogMessage("coinflip controller", "play against bot", "success", logrus.Fields{"round": round.ID, "user": userID, "winner": winnerId})
}
