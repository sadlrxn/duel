package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TxDepositSol              TransactionType = "deposit_sol"
	TxDepositNft              TransactionType = "deposit_nft"
	TxWithdrawSol             TransactionType = "withdraw_sol"
	TxWithdrawNft             TransactionType = "withdraw_nft"
	TxWithdrawSpl             TransactionType = "withdraw_spl"
	TxJackpotBet              TransactionType = "jackpot_bet"
	TxJackpotProfit           TransactionType = "jackpot_profit"
	TxJackpotFee              TransactionType = "jackpot_fee"
	TxGrandJackpotBet         TransactionType = "grand_jackpot_bet"
	TxGrandJackpotProfit      TransactionType = "grand_jackpot_profit"
	TxGrandJackpotFee         TransactionType = "grand_jackpot_fee"
	TxCoinflipBet             TransactionType = "coinflip_bet"
	TxCoinflipCancel          TransactionType = "coinflip_cancel"
	TxCoinflipProfit          TransactionType = "coinflip_profit"
	TxCoinflipFee             TransactionType = "coinflip_fee"
	TxCoinflipRefill          TransactionType = "coinflip_refill"
	TxTip                     TransactionType = "tip"
	TxDreamtowerBet           TransactionType = "dreamtower_bet"
	TxDreamtowerFee           TransactionType = "dreamtower_fee"
	TxDreamtowerProfit        TransactionType = "dreamtower_profit"
	TxClaimStakingReward      TransactionType = "claim_staking_reward"
	TxClaimRakebackReward     TransactionType = "claim_rakeback_reward"
	TxClaimAffiliateReward    TransactionType = "claim_affiliate_reward"
	TxRain                    TransactionType = "rain"
	TxExchangeCouponToChips   TransactionType = "exchange_coupon_to_chips"
	TxCrashBet                TransactionType = "crash_bet"
	TxCrashProfit             TransactionType = "crash_profit"
	TxCrashFee                TransactionType = "crash_fee"
	TxClaimDailyRaceReward    TransactionType = "claim_daily_race_reward"
	TxClaimWeeklyRaffleReward TransactionType = "claim_weekly_raffle_reward"
)

type TransactionStatus string

const (
	TransactionSucceed TransactionStatus = "succeed"
	TransactionFailed  TransactionStatus = "failed"
	TransactionPending TransactionStatus = "pending"
)

type TransactionOwnerType string

const (
	TransactionUserReferenced               TransactionOwnerType = "tx_user_referenced"
	TransactionWalletReferenced             TransactionOwnerType = "tx_wallet_referenced"
	TransactionJackpotReferenced            TransactionOwnerType = "tx_jackpot_referenced"
	TransactionCoinflipReferenced           TransactionOwnerType = "tx_coinflip_referenced"
	TransactionDreamTowerReferenced         TransactionOwnerType = "tx_dream_tower_referenced"
	TransactionPaymentReferenced            TransactionOwnerType = "tx_payment_referenced"
	TransactionCouponTransactionReferenced  TransactionOwnerType = "tx_coupon_transaction_referenced"
	TransactionCrashBetReferencedForCashIn  TransactionOwnerType = "tx_crash_bet_referenced_for_cash_in"
	TransactionCrashBetReferencedForCashOut TransactionOwnerType = "tx_crash_bet_referenced_for_cash_out"
	TransactionCrashRoundReferencedForFee   TransactionOwnerType = "tx_crash_round_referenced_for_fee"
	TransactionCrashRoundReferenced         TransactionOwnerType = "tx_crash_round_referenced"
	TransactionCrashBetReferenced           TransactionOwnerType = "tx_crash_bet_referenced"
	TransactionDailyRaceRewardsReferenced   TransactionOwnerType = "tx_daily_race_rewards_referenced"
	TransactionWeeklyRaffleRewardReferenced TransactionOwnerType = "tx_weekly_raffle_reward_referenced"
)

type Transaction struct {
	gorm.Model

	FromWallet *uint             `json:"fromUser"`
	ToWallet   *uint             `json:"toUser"`
	Balance    Balance           `gorm:"not null;polymorphic:Owner;polymorphicValue:in-transaction" json:"balance"`
	Type       TransactionType   `gorm:"not null" json:"type"`
	Status     TransactionStatus `gorm:"not null;default:pending" json:"status"`

	FromWalletPrevID *uint    `json:"fromWalletPrevId"`
	FromWalletPrev   *Balance `gorm:"foreignKey:FromWalletPrevID" json:"fromWalletPrev"`
	FromWalletNextID *uint    `json:"fromWalletNextId"`
	FromWalletNext   *Balance `gorm:"foreignKey:FromWalletNextID" json:"fromWalletNext"`
	RefundPrevID     *uint    `json:"refundPrevId"`
	RefundPrev       *Balance `gorm:"foreignKey:FromWalletNextID" json:"refundPrev"`
	RefundNextID     *uint    `json:"refundNextId"`
	RefundNext       *Balance `gorm:"foreignKey:FromWalletNextID" json:"refundNext"`

	ToWalletPrevID *uint    `json:"toWalletPrevId"`
	ToWalletPrev   *Balance `gorm:"foreignKey:ToWalletPrevID" json:"toWalletPrev"`
	ToWalletNextID *uint    `json:"toWalletNextId"`
	ToWalletNext   *Balance `gorm:"foreignKey:ToWalletNextID" json:"toWalletNext"`

	Receipients pq.Int64Array        `gorm:"type:bigint[]" json:"recipients"`
	OwnerID     uint                 `json:"ownerId"`
	OwnerType   TransactionOwnerType `json:"ownerType"`
}
