package dreamtower

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm/clause"
)

// @Internal
// Get User's currently playing round
func getUserPlayingRound(user *db_aggregator.User, lock bool, sessionId ...db_aggregator.UUID) (*models.DreamTowerRound, error) {
	if user == nil {
		return nil, utils.MakeError(
			"dreamtower", "getUserPlayingRound", "invalid user", nil,
		)
	}

	session, err := db_aggregator.GetSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"dreamtower", "getUserPlayingRound", "failed to get session", err,
		)
	}

	var playingRound models.DreamTowerRound
	if lock {
		if result := session.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND status = ?", user, models.DreamTowerPlaying).Last(&playingRound); result.Error != nil {
			return nil, utils.MakeError(
				"dreamtower", "getUserPlayingRound", "failed to get round", result.Error,
			)
		}
	} else {
		if result := session.Where("user_id = ? AND status = ?", user, models.DreamTowerPlaying).Last(&playingRound); result.Error != nil {
			return nil, utils.MakeError(
				"dreamtower", "getUserPlayingRound", "failed to get round", result.Error,
			)
		}
	}

	return &playingRound, nil
}

// @External
// Create new round
func createRound(round *models.DreamTowerRound) error {
	if round == nil {
		return errors.New("invalid round")
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return err
	}

	if result := session.Create(&round); result.Error != nil {
		return result.Error
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return err
	}

	return nil
}

// @External
// Save round with params below
// Bets, Status, Profit
func saveRound(user *db_aggregator.User, request *saveRoundRequest) (*models.DreamTowerRound, error) {
	if user == nil {
		return nil, errors.New("invalid round")
	}
	if request == nil || request.bets == nil && request.profit == nil && request.status == nil {
		return nil, errors.New("invalid request")
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, err
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	round, err := getUserPlayingRound(user, true, sessionId)
	if err != nil {
		return nil, err
	}

	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	if request.bets != nil {
		round.Bets = *request.bets
	}
	if request.status != nil {
		round.Status = *request.status
	}
	if request.profit != nil {
		round.Profit = request.profit
	}

	if result := session.Save(&round); result.Error != nil {
		return nil, result.Error
	}

	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, err
	}

	return round, nil
}

// @Internal
// Get history rounds
func getHistory(user *db_aggregator.User, offset int, count int) (*[]models.DreamTowerRound, error) {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, err
	}

	session = session.Order("id desc").
		Where("bet_amount > ?", 0).
		Where("status IN ?",
			[]models.DreamTowerStatus{
				models.DreamTowerWin,
				models.DreamTowerLoss,
				models.DreamTowerCashout,
			},
		)

	if user != nil {
		session = session.Where("user_id = ?", user)
	}

	var rounds []models.DreamTowerRound
	if result := session.Offset(offset).Limit(count).Find(&rounds); result.Error != nil {
		return nil, err
	}

	return &rounds, nil
}

func getTempWalletBalance() (*db_aggregator.BalanceLoad, error) {
	tempBalance, err := db_aggregator.GetUserBalance((*db_aggregator.User)(&config.DREAMTOWER_TEMP_ID))
	if err != nil {
		return nil, err
	}
	tempBalanceLoad, err := db_aggregator.GetBalance(tempBalance)
	if err != nil {
		return nil, err
	}

	return tempBalanceLoad, nil
}

func (c *Controller) cashIn(userID uint, betAmount int64) (*[]uint, *models.PaidBalanceForGame, error) {
	var txs []uint
	var paidBalanceType models.PaidBalanceForGame = models.ChipBalanceForGame
	if betAmount <= 0 {
		return &txs, &paidBalanceType, nil
	}
	result, tx, err := coupon.TryBet(coupon.TryBetWithCouponRequest{
		UserID:  userID,
		Balance: betAmount,
		Type:    models.CpTxDreamtowerBet,
	})
	if result == coupon.CouponBetUnavailable {
		tx1, err := transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&userID),
			ToUser:   (*db_aggregator.User)(&config.DREAMTOWER_TEMP_ID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &betAmount,
			},
			Type:          models.TxDreamtowerBet,
			ToBeConfirmed: false,
		})
		if err != nil {
			return nil, nil, err
		}
		fee := betAmount * config.DREAMTOWER_FEE / 100
		tx2, err := transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&config.DREAMTOWER_TEMP_ID),
			ToUser:   (*db_aggregator.User)(&config.DREAMTOWER_FEE_ID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &fee,
			},
			Type:          models.TxDreamtowerFee,
			ToBeConfirmed: false,
			HouseFeeMeta: &transaction.HouseFeeMeta{
				User:        db_aggregator.User(userID),
				WagerAmount: betAmount,
			},
		})
		if err != nil {
			transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx1,
				OwnerID:     userID,
				OwnerType:   models.TransactionUserReferenced,
			})
			return nil, nil, utils.MakeError(
				"dreamtowerCashIn",
				"transfer fee",
				"failed to transfer round fee",
				err,
			)
		}
		txs = []uint{uint(*tx1), uint(*tx2)}
		paidBalanceType = models.ChipBalanceForGame
	} else if result == coupon.CouponBetFailed || result == coupon.CouponBetInsufficientFunds {
		return nil, nil, utils.MakeError(
			"dreamtowerCashIn",
			"coupon bet",
			"failed to bet coupon",
			err,
		)
	} else if result == coupon.CouponBetSucceed {
		txs = []uint{tx}
		paidBalanceType = models.CouponBalanceForGame
	}
	return &txs, &paidBalanceType, nil
}

func cashOut(userID uint, roundID uint, profit int64, paidBalanceType models.PaidBalanceForGame) error {
	if paidBalanceType == models.ChipBalanceForGame {
		_, err := transaction.Transfer(&transaction.TransactionRequest{
			FromUser: (*db_aggregator.User)(&config.DREAMTOWER_TEMP_ID),
			ToUser:   (*db_aggregator.User)(&userID),
			Balance: db_aggregator.BalanceLoad{
				ChipBalance: &profit,
			},
			Type:          models.TxDreamtowerProfit,
			ToBeConfirmed: true,
			OwnerID:       roundID,
			OwnerType:     models.TransactionDreamTowerReferenced,
		})
		return err
	} else if paidBalanceType == models.CouponBalanceForGame {
		_, err := coupon.Perform(coupon.CouponTransactionRequest{
			UserID:        userID,
			Balance:       profit,
			Type:          models.CpTxDreamtowerProfit,
			ToBeConfirmed: true,
		})
		return err
	}
	return utils.MakeError(
		"dreamtowerCashout",
		"cash out",
		"invalid balance type",
		errors.New(fmt.Sprintf("invalid paid balance type: %v", paidBalanceType)),
	)
}

func confirmTransactions(txs []uint, paidBalanceType models.PaidBalanceForGame, ownerID uint, ownerType models.TransactionOwnerType) error {
	var err error
	if len(txs) == 0 {
		return errors.New("transaction array is empty")
	}
	if paidBalanceType == models.ChipBalanceForGame {
		for _, tx := range txs {
			err = transaction.Confirm(transaction.ConfirmRequest{
				Transaction: db_aggregator.Transaction(tx),
				OwnerID:     ownerID,
				OwnerType:   ownerType,
			})
			if err != nil {
				err = utils.MakeError(
					"dreamtower",
					"confirm chip transactions",
					"failed to confirm transaction",
					err,
				)
			}
		}
	} else if paidBalanceType == models.CouponBalanceForGame {
		for _, tx := range txs {
			err = coupon.Confirm(tx)
			if err != nil {
				err = utils.MakeError(
					"dreamtower",
					"confirm coupon transactions",
					"failed to confirm transaction",
					err,
				)
			}
		}
	} else {
		err = utils.MakeError(
			"dreamtower",
			"confirm transactions",
			"invalid balance type",
			nil,
		)
	}
	return err
}

func declineTransactions(txs []uint, paidBalanceType models.PaidBalanceForGame, ownerID uint, ownerType models.TransactionOwnerType) error {
	var err error
	if paidBalanceType == models.ChipBalanceForGame {
		for _, tx := range txs {
			err = transaction.Decline(transaction.DeclineRequest{
				Transaction: db_aggregator.Transaction(tx),
				OwnerID:     ownerID,
				OwnerType:   ownerType,
			})
			if err != nil {
				err = utils.MakeError(
					"dreamtower",
					"decline chip transactions",
					"failed to decline transaction",
					err,
				)
			}
		}
	} else if paidBalanceType == models.CouponBalanceForGame {
		for _, tx := range txs {
			err = coupon.Decline(tx)
			if err != nil {
				err = utils.MakeError(
					"dreamtower",
					"decline coupon transactions",
					"failed to decline transaction",
					err,
				)
			}
		}
	} else {
		err = utils.MakeError(
			"dreamtower",
			"decline transactions",
			"invalid balance type",
			nil,
		)
	}
	return err
}

func byteGenerator(serverSeed string, clientSeed string, nonce int, cursor int) []byte {
	currentRound := cursor / 32
	currentRoundCursor := cursor % 32
	str := fmt.Sprintf("%s:%s:%d:%d", serverSeed, clientSeed, nonce, currentRound)
	sum := sha256.Sum256([]byte(str))
	return sum[currentRoundCursor : currentRoundCursor+4]
}

func generateEvent(rows int, starsInRow int, shuffle []int) []int {
	arr := []int{}
	for i := 0; i < rows; i++ {
		arr = append(arr, i)
	}
	for i := rows - 1; i > 0; i-- {
		arr[i], arr[shuffle[i]] = arr[shuffle[i]], arr[i]
	}
	return arr[:starsInRow]
}

func generateTower(serverSeed string, clientSeed string, nonce uint, rows int, starsInRow int, count int) [][]int {
	var tower [][]int
	cursor := 0
	for i := 0; i < count; i++ {
		shuffle := []int{}
		bytes := byteGenerator(serverSeed, clientSeed, int(nonce), cursor)
		for j := 1; j <= rows; j++ {
			value := int(bytes[j-1]) * j / 256
			shuffle = append(shuffle, value)
		}
		tower = append(tower, generateEvent(rows, starsInRow, shuffle))
		cursor += 4
	}
	return tower
}

func checkResult(tower [][]int, bets []int32, autoMode bool) models.DreamTowerStatus {
	for i := 0; i < len(bets); i++ {
		row := tower[i]
		bet := bets[i]
		var isStar bool
		for j := 0; j < len(row); j++ {
			if int32(row[j]) == bet {
				isStar = true
				break
			}
		}
		if !isStar {
			return models.DreamTowerLoss
		}
	}
	if len(bets) == len(tower) {
		return models.DreamTowerWin
	} else if autoMode && len(bets) != 0 {
		return models.DreamTowerCashout
	}
	return models.DreamTowerPlaying
}

func calculateMutiplier(difficulty models.DreamTowerDifficulty, fee uint, level int) float32 {
	return calculateMutiplierV2(difficulty, fee, level)
}

func calculateMutiplierV1(difficulty models.DreamTowerDifficulty, fee uint, level int) float32 {
	multiplier := math.Pow(float64(difficulty.BlocksInRow), float64(level)) / math.Pow(float64(difficulty.StarsInRow), float64(level))
	multiplier = multiplier * float64(100-fee) / float64(100)
	return float32(math.Floor(multiplier*100) / 100)
}

func calculateMutiplierV2(difficulty models.DreamTowerDifficulty, fee uint, level int) float32 {
	var duel_ev float64 = float64(100-fee*uint(level)) / 100
	original_odd := math.Pow(float64(difficulty.BlocksInRow), float64(level)) / math.Pow(float64(difficulty.StarsInRow), float64(level))
	multiplier := original_odd * duel_ev
	return float32(math.Floor(multiplier*100) / 100)
}
