package coinflip

import (
	"encoding/json"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func (c *Controller) refundPaidBalance(userID uint, amount int64, paidBalanceType models.PaidBalanceForGame) {
	b, _ := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     amount,
			Wagered:     amount,
			BalanceType: paidBalanceType,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
}

func (c *Controller) validateAmount(userID uint, eventParam EventParam) bool {
	if eventParam.Amount < c.minAmount || eventParam.Amount > c.maxAmount {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Invalid bet amount.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}

		c.refundPaidBalance(userID, eventParam.Amount, eventParam.PaidBalanceType)
		return false
	}
	return true
}

func (c *Controller) validateCount(userID uint, eventParam EventParam) bool {
	roundCount := uint(0)
	c.round2Creator.Range(func(key, value any) bool {
		if value.(uint) == userID {
			roundCount++
		}
		return true
	})
	if roundCount == c.roundLimit {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Exceed round count limit.", RoundID: 0}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}

		c.refundPaidBalance(userID, eventParam.Amount, eventParam.PaidBalanceType)
		return false
	}
	return true
}

func (c *Controller) validateTicketID(userID uint, amount int64, tx *db_aggregator.Transaction) (string, error) {
	ticketID, err := utils.RequestTicketID()

	if err != nil && tx != nil {
		b, _ := json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     amount,
				BalanceType: models.ChipBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userID}, Message: b}
		transaction.Decline(transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     userID,
			OwnerType:   models.TransactionUserReferenced,
		})
		return "", err
	}
	return ticketID, nil
}

func (c *Controller) validateRandomString(userID uint, round models.CoinflipRound, tx *db_aggregator.Transaction) (string, error) {
	signedString, err := utils.GenerateRandomString(round.TicketID)
	if err != nil && tx != nil {
		b, _ := json.Marshal(types.WSMessage{
			Room:      string(types.Coinflip),
			EventType: "message",
			Payload:   types.ErrorMessagePayload{Message: "Failed to generate random string.", RoundID: round.ID}})
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

		declineRequest := transaction.DeclineRequest{
			Transaction: *tx,
			OwnerID:     round.ID,
			OwnerType:   models.TransactionCoinflipReferenced,
		}
		if round.ID == 0 {
			declineRequest.OwnerID = userID
			declineRequest.OwnerType = models.TransactionUserReferenced
		}
		transaction.Decline(declineRequest)
		return "", err
	}
	return signedString, nil
}
