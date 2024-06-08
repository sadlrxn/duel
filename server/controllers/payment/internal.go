package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/solana"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	userController "github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

func (c *Controller) depositChip(walletAddress string, decodedResult *solana.DecodedTransaction, txId string) {
	db := db.GetDB()
	var user models.User
	if result := db.Where("wallet_address = ?", walletAddress).First(&user); result.Error != nil {
		log.LogMessage("payment chip deposit handler", "can not find User from db", "error", logrus.Fields{"wallet": walletAddress})
		return
	}

	var pay models.Payment
	if err := db.Where("tx_hash = ?", txId).First(&pay).Error; err == nil {
		log.LogMessage("payment chip deposit handler", "duplicated transaction", "error", logrus.Fields{"tx": txId})
		return
	}

	conf := config.Get()
	var splAmount = float32(float64(decodedResult.Lamports) / math.Pow10(decodedResult.SplToken.Decimals))
	var cashAmount int64
	if conf.Network == "mainnet" && config.USDC_SPL_ADDRESS != decodedResult.SplToken.MintAddress.String() {
		usdAmount, err := utils.SwapTokens(splAmount, decodedResult.SplToken.MintAddress.String(), config.USDC_SPL_ADDRESS)
		if err != nil {
			log.LogMessage("payment chip deposit handler", "failed to swap SPL to USDC via Jupiter instance", "error", logrus.Fields{"err": err, "amount": splAmount})
			tm := map[string]string{decodedResult.SplToken.Keyword: decodedResult.SplToken.MintAddress.String()}
			prices := utils.FetchTokenPrice(tm)
			cashAmount = int64(float64(decodedResult.Lamports) * float64(prices[decodedResult.SplToken.Keyword]) * math.Pow10(config.BALANCE_DECIMALS) / math.Pow10(decodedResult.SplToken.Decimals))
		} else {
			log.LogMessage("payment chip deposit handler", "swap SPL => USDC succeed.", "success", logrus.Fields{"spl": splAmount, "USDC": usdAmount})
			cashAmount = int64(float64(usdAmount) * math.Pow10(config.BALANCE_DECIMALS))
		}
	} else {
		tm := map[string]string{decodedResult.SplToken.Keyword: decodedResult.SplToken.MintAddress.String()}
		prices := utils.FetchTokenPrice(tm)
		cashAmount = int64(float64(decodedResult.Lamports) * float64(prices[decodedResult.SplToken.Keyword]) * math.Pow10(config.BALANCE_DECIMALS) / math.Pow10(decodedResult.SplToken.Decimals))
	}

	payment := models.Payment{
		UserID: user.ID,
		Type:   "deposit_" + strings.ToLower(decodedResult.SplToken.Keyword),
		Status: models.Success,
		SolDetail: models.SolDetail{
			SolAmount: int64(decodedResult.Lamports),
			UsdAmount: cashAmount,
		},
		TxHash: txId,
	}
	if result := db.Create(&payment); result.Error != nil {
		log.LogMessage("payment chip deposit handler", "failed to create record on payments", "error", logrus.Fields{})
		return
	}

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: nil,
		ToUser:   (*db_aggregator.User)(&user.ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &cashAmount,
		},
		Type:          models.TxDepositSol,
		ToBeConfirmed: true,
		OwnerID:       user.ID,
		OwnerType:     models.TransactionUserReferenced,
	})
	if err != nil {
		log.LogMessage(
			"payment chip deposit handler",
			"can not transfer balances",
			"error",
			logrus.Fields{
				"error":  err.Error(),
				"wallet": walletAddress},
		)
		return
	}

	payment.TransactionID = (*uint)(tx)
	if result := db.Save(&payment); result.Error != nil {
		log.LogMessage("payment chip deposit handler", "failed to save payment with transaction id", "error", logrus.Fields{})
		return
	}

	// Try to applying first deposit bonus.
	bonusBalance, err := transaction.TryApplyForFirstDepositBonus(
		user.ID,
		cashAmount,
	)
	if err != nil {
		log.LogMessage(
			"payment_internal_depositChips",
			"failed to perform first deposit bonus",
			"error",
			logrus.Fields{
				"userID":        user.ID,
				"depositAmount": cashAmount,
				"error":         err.Error(),
			},
		)
	}

	b, _ := json.Marshal(struct {
		EventType string `json:"eventType"`
		TxID      string `json:"txId"`
		Amount    uint   `json:"amount"`
	}{EventType: "deposit_sol", TxID: txId, Amount: uint(cashAmount)})
	c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}

	b, _ = json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     cashAmount,
			BalanceType: models.ChipBalanceForGame,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}

	// Send coupon balance updating event if any redeemed.
	if bonusBalance > 0 {
		b, _ = json.Marshal(types.WSMessage{
			EventType: "balance_update",
			Payload: types.BalanceUpdatePayload{
				UpdateType:  types.Increase,
				Balance:     bonusBalance,
				BalanceType: models.CouponBalanceForGame,
				Delay:       0,
			}})
		c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
	}

	log.LogMessage("payment chip deposit handler", "deposit chip succeed", "success", logrus.Fields{"user": user.ID, "amount": cashAmount})
}

func (c *Controller) depositNfts(walletAddress string, mintAddresses []string, txId string) {
	db := db.GetDB()
	var user models.User
	if result := db.Where("wallet_address = ?", walletAddress).First(&user); result.Error != nil {
		log.LogMessage("payment nft deposit handler", "failed to get User data from db", "error", logrus.Fields{"wallet": walletAddress})
		return
	}

	var pay models.Payment
	if err := db.Where("tx_hash = ?", txId).First(&pay).Error; err == nil {
		log.LogMessage("payment nft deposit handler", "duplicated transaction", "error", logrus.Fields{"tx": txId, "error": err.Error()})
		return
	}

	payment := models.Payment{
		UserID: user.ID,
		Type:   "deposit_nft",
		Status: models.Success,
		NftDetail: models.NftDetail{
			Mints: mintAddresses,
		},
		TxHash: txId,
	}
	if result := db.Create(&payment); result.Error != nil {
		log.LogMessage("payment nft deposit handler", "failed to create record on payments", "error", logrus.Fields{"error": result.Error.Error()})
		return
	}

	tx, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: nil,
		ToUser:   (*db_aggregator.User)(&user.ID),
		Balance: db_aggregator.BalanceLoad{
			NftBalance: db_aggregator.ConvertStringArrayToNftArray(&mintAddresses),
		},
		Type:          models.TxDepositNft,
		ToBeConfirmed: true,
		OwnerID:       user.ID,
		OwnerType:     models.TransactionUserReferenced,
	})
	if err != nil {
		log.LogMessage("payment nft deposit handler", "Can not transfer balances", "error", logrus.Fields{"error": err.Error()})
		return
	}

	payment.TransactionID = (*uint)(tx)
	if result := db.Save(&payment); result.Error != nil {
		log.LogMessage("payment nft deposit handler", "failed to save payment with transaction id", "error", logrus.Fields{"error": result.Error.Error()})
		return
	}

	b, _ := json.Marshal(struct {
		EventType     string   `json:"eventType"`
		TxID          string   `json:"txId"`
		MintAddresses []string `json:"mintAddresses"`
	}{EventType: "deposit_nft", TxID: txId, MintAddresses: mintAddresses})
	c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}

	_, nfts := userController.GetNftDetailsFromMintAddresses(mintAddresses)
	b, _ = json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType: types.Increase,
			Nfts:       nfts,
			Delay:      0,
		}})
	c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}

	log.LogMessage("payment nft deposit handler", "deposit nft succeed", "success", logrus.Fields{"user": user.ID, "nfts": mintAddresses, "tx": txId})
}

func parseSolanaTransaction(byteArray []byte) types.TransactionDataResult {
	var result types.TransactionDataResult
	err := json.Unmarshal(byteArray, &result)
	if err != nil {
		log.LogMessage("parse transaction", "invalid transaction", "error", logrus.Fields{})
	}
	return result
}

func getSolanaTransaction(txID string) (types.TransactionDataResult, error) {
	config := config.Get()
	reqUrl := "https://api-eu1.tatum.io/v3/solana/transaction/" + txID
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Add("x-api-key", config.TatumApiKey)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.LogMessage("payment handler", "failed to get transaction data", "error", logrus.Fields{"error": err.Error()})
		return types.TransactionDataResult{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if err != nil {
		log.LogMessage("payment handler", "unable to read response body", "error", logrus.Fields{"error": err.Error()})
		return types.TransactionDataResult{}, err
	}
	result := parseSolanaTransaction(body)

	return result, nil
}

// Temporary function to handle migeration from `MintAddresses` to `Mints`
func MigratePaymentModel() error {
	const MIGRATE_DONE = "migrate done"
	db := db.GetDB()
	var allPayments []models.Payment
	if result := db.Find(&allPayments); result.Error != nil {
		return fmt.Errorf("failed to migrate Payment Model. \r\nError: %v", result.Error)
	}

	totalError := errors.New("")
	for _, payment := range allPayments {
		if payment.NftDetail.MintAddresses == "" || payment.NftDetail.MintAddresses == MIGRATE_DONE {
			continue
		}
		if err := json.Unmarshal([]byte(payment.NftDetail.MintAddresses), &payment.NftDetail.Mints); err != nil {
			totalError = fmt.Errorf("%v\r\n%v", totalError.Error(), err.Error())
			continue
		}
		payment.NftDetail.MintAddresses = MIGRATE_DONE
		db.Save(&payment)
	}

	if totalError.Error() == "" {
		return nil
	}
	return totalError
}
