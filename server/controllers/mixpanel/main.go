package mixpanel

import (
	"math"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/sirupsen/logrus"
	"github.com/sjhitchner/go-mixpanel"
)

var mixpanelClient *mixpanel.MixpanelClient

// @External
// Initializer
func Init(token string, url string) {
	mixpanelClient = mixpanel.NewMixpanelWithUrl(token, url)
}

// @External
// Get MixpanelClient object
func Get() *mixpanel.MixpanelClient {
	return mixpanelClient
}

// @External
// Track event on every transaction.
func TrackTransaction(transactionRequest TransactionBody) error {
	db := db.GetDB()
	var from, to models.User
	db.First(&from, transactionRequest.FromUser)
	db.First(&to, transactionRequest.ToUser)
	switch transactionRequest.Type {
	case models.TxDepositSol:
		trackDepositSol(transactionRequest, to)
	case models.TxWithdrawSol:
		trackWithdrawSOL(transactionRequest, from)
	case models.TxCoinflipBet:
		trackCoinflipBet(transactionRequest, from)
	case models.TxCoinflipCancel:
		trackCoinflipCancel(transactionRequest, to)
	case models.TxCoinflipProfit:
		trackCoinflipProfit(transactionRequest, to)
	case models.TxCoinflipFee:
		trackCoinflipFee(transactionRequest, to)
	case models.TxTip:
		trackTip(transactionRequest, from, to)
	case models.TxDreamtowerBet:
		trackDreamtowerBet(transactionRequest, from)
	case models.TxDreamtowerProfit:
		trackDreamtowerProfit(transactionRequest, to)
	case models.TxDreamtowerFee:
		trackDreamtowerFee(transactionRequest, to)
	}
	return nil
}

// @External
// Tracker for jackpot events
func TrackJackpot(transactionRequest TransactionBody) {
	db := db.GetDB()
	var from, to models.User
	db.First(&from, transactionRequest.FromUser)
	db.First(&to, transactionRequest.ToUser)
	switch transactionRequest.Type {
	case models.TxJackpotBet:
		trackJackpotBet(transactionRequest, from)
	case models.TxJackpotProfit:
		trackJackpotProfit(transactionRequest, to)
	case models.TxJackpotFee:
		trackJackpotFee(transactionRequest, to)
	case models.TxGrandJackpotBet:
		trackGrandJackpotBet(transactionRequest, from)
	case models.TxGrandJackpotProfit:
		trackGrandJackpotProfit(transactionRequest, to)
	case models.TxGrandJackpotFee:
		trackGrandJackpotFee(transactionRequest, to)
	}
}

// @External
// Set user profile fields & track user sign in event.
func TrackUserSignIn(user models.User, ip string) error {
	var update = mixpanel.NewUpdate(user.WalletAddress)
	var set = mixpanel.NewSet()
	set.AddProperty("name", user.Name)
	set.AddProperty("Role", user.Role)
	set.AddProperty("Avatar", user.Avatar)
	set.AddProperty("Nonce", user.Nonce)
	set.AddProperty("Status", "Signed In")
	update.Set = set
	update.Ip = ip
	if err := mixpanelClient.Update(update); err != nil {
		log.LogMessage("sign in user", "update user failed", "error", logrus.Fields{"error": err.Error})
		return err
	}

	var event = mixpanel.NewEvent("Sign In")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("sign in user", "tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @External
// Set user profile fields & track user sign out event.
func TrackUserSignOut(user models.User, ip string) error {
	var update = mixpanel.NewUpdate(user.WalletAddress)
	var set = mixpanel.NewSet()
	set.AddProperty("Status", "Signed Out")
	update.Ip = ip
	update.Set = set
	if err := mixpanelClient.Update(update); err != nil {
		log.LogMessage("sign out user", "update user failed", "error", logrus.Fields{"error": err.Error})
		return err
	}

	var event = mixpanel.NewEvent("Sign Out")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("sign out user", "tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @External
// Set user profile fields & track user sign up event.
func TrackUserSignUp(user models.User, ip string) error {
	var update = mixpanel.NewUpdate(user.WalletAddress)
	var set = mixpanel.NewSet()
	set.AddProperty("name", user.Name)
	set.AddProperty("Role", user.Role)
	set.AddProperty("Avatar", user.Avatar)
	set.AddProperty("Created At", time.Now().String())
	update.Set = set
	update.Ip = ip
	if err := mixpanelClient.Update(update); err != nil {
		log.LogMessage("sign up user", "update user failed", "error", logrus.Fields{"error": err.Error})
		return err
	}

	var event = mixpanel.NewEvent("Sign Up")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("sign up user", "tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @External
// Update user profile on update user.
func TrackUserUpdate(user models.User, ip string) error {
	var update = mixpanel.NewUpdate(user.WalletAddress)
	var set = mixpanel.NewSet()
	set.AddProperty("name", user.Name)
	set.AddProperty("Avatar", user.Avatar)
	update.Set = set
	update.Ip = ip
	if err := mixpanelClient.Update(update); err != nil {
		log.LogMessage("update user", "update user failed", "error", logrus.Fields{"error": err.Error})
		return err
	}

	var event = mixpanel.NewEvent("Update Profile")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Name", user.Name)
	event.AddProperty("Avatar", user.Avatar)
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("update user", "tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @External
// Get current total chip amount
func TrackChipAmount() {
	db := db.GetDB()
	var users []models.User
	var total int64
	db.Preload("Wallet.Balance.ChipBalance").Where("id <= 1000 OR id >= 1010").Find(&users)
	for _, user := range users {
		total += user.Wallet.Balance.ChipBalance.Balance
	}

	var event = mixpanel.NewEvent("Total Chips")
	event.SetTime(time.Now())
	event.SetDistinctId("duel")
	event.AddProperty("Amount", float64(total)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("chip balance", "tracking event failed", "error", logrus.Fields{"error": err.Error})
	}
}
