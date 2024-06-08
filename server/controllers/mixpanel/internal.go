package mixpanel

import (
	"math"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/sirupsen/logrus"
	"github.com/sjhitchner/go-mixpanel"
)

// @Interanl
// Update user profile & track event on deposit sol.
func trackDepositSol(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Deposit SOL")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("deposit sol", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Interanl
// Update user profile & track event on withdraw sol.
func trackWithdrawSOL(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Withdraw SOL")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("withdraw sol", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on jackpot bet.
func trackJackpotBet(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Jackpot Bet")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("jackpot bet", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on jackpot win.
func trackJackpotProfit(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Jackpot Win")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("jackpot profit", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event for jackpot fee.
func trackJackpotFee(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Jackpot Fee")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("jackpot fee", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on grand jackpot bet.
func trackGrandJackpotBet(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Grand Jackpot Bet")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("grand jackpot bet", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on grand jackpot win.
func trackGrandJackpotProfit(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Grand Jackpot Win")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("grand jackpot profit", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event for grand jackpot fee.
func trackGrandJackpotFee(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Grand Jackpot Fee")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Chip Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("NFT Amount", float64(*transactionRequest.NftBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("grand jackpot fee", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on coinflip bet.
func trackCoinflipBet(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Coinflip Bet")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("coinflip bet", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on coinflip win.
func trackCoinflipProfit(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Coinflip Win")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("coinflip profit", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event for coinflip fee.
func trackCoinflipFee(transactionRequest TransactionBody, user models.User) error {

	var event = mixpanel.NewEvent("Coinflip Fee")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("coinflip fee", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on coinflip cancel.
func trackCoinflipCancel(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Coinflip Cancel")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("coinflip cancel", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on dream tower bet.
func trackDreamtowerBet(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Dreamtower Bet")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("dreamtower bet", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on dream tower win.
func trackDreamtowerProfit(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Dreamtower Win")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("dreamtower profit", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event for dream tower fee.
func trackDreamtowerFee(transactionRequest TransactionBody, user models.User) error {
	var event = mixpanel.NewEvent("Dreamtower Fee")
	event.SetTime(time.Now())
	event.SetDistinctId(user.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("dreamtower fee", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}

// @Internal
// Update user profile & track event on send tip.
func trackTip(transactionRequest TransactionBody, from models.User, to models.User) error {
	var event = mixpanel.NewEvent("Tip")
	event.SetTime(time.Now())
	event.SetDistinctId(from.WalletAddress)
	event.AddProperty("Amount", float64(*transactionRequest.ChipBalance)/math.Pow10(config.BALANCE_DECIMALS))
	event.AddProperty("Reciepient", to.WalletAddress)
	if err := mixpanelClient.Track(event); err != nil {
		log.LogMessage("send tip", "mixpanel: tracking event failed", "error", logrus.Fields{"error": err.Error})
		return err
	}
	return nil
}
