package db_aggregator

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) error {
	return initialize(db)
}

func StartSession() (UUID, error) {
	return startSession()
}

func GetSession(sessionId ...UUID) (*gorm.DB, error) {
	return getSession(sessionId...)
}

func CommitSession(session UUID) error {
	return commitSession(session)
}

func RemoveSession(session UUID) error {
	return removeSession(session)
}

func GetBalance(balance *Balance, sessionId ...UUID) (*BalanceLoad, error) {
	return getBalance(balance, sessionId...)
}

func GetUserWallet(user *User, sessionId ...UUID) (*Wallet, error) {
	return getUserWallet(user, false, sessionId...)
}

func GetWalletBalance(wallet *Wallet, sessionId ...UUID) (*Balance, error) {
	return getWalletBalance(wallet, sessionId...)
}

func GetUserBalance(user *User, sessionId ...UUID) (*Balance, error) {
	return getUserBalance(user, false, sessionId...)
}

func GetNftWallet(nft *Nft, sessionId ...UUID) (*Wallet, error) {
	return getNftWallet(nft, sessionId...)
}

func Transfer(from *User, to *User, balanceLoad *BalanceLoad, sessionId ...UUID) (*TransferResult, error) {
	return transfer(from, to, balanceLoad, sessionId...)
}

func Burn(from *User, amount int64, sessionId ...UUID) error {
	return burn(from, amount, sessionId...)
}

func Rain(from *User, to *[]User, balanceLoad *BalanceLoad, sessionId ...UUID) (*RainResult, error) {
	return rain(from, to, balanceLoad, sessionId...)
}

func RecordTransaction(transactionLoad *TransactionLoad, sessionId ...UUID) (*Transaction, error) {
	return recordTransaction(transactionLoad, sessionId...)
}

func ConfirmTransaction(transactionLoad *TransactionLoad, transaction *Transaction, sessionId ...UUID) error {
	return confirmTransaction(transactionLoad, transaction, sessionId...)
}

func DeclineTransaction(transactionLoad *TransactionLoad, transaction *Transaction, sessionId ...UUID) error {
	return declineTransaction(transactionLoad, transaction, sessionId...)
}

func RecordBalanceHistory(balance *Balance, changedBalance ChangedBalance, sessionId ...UUID) (*BalanceHistoryChain, error) {
	return recordBalanceHistory(balance, changedBalance, sessionId...)
}

func GetTransactionLoad(transaction *Transaction, sessionId ...UUID) (*TransactionLoad, error) {
	return getTransactionLoad(transaction, sessionId...)
}

func GetWalletUser(wallet *Wallet, sessionId ...UUID) (*User, error) {
	return getWalletUser(wallet, sessionId...)
}

func ConvertStringArrayToNftArray(strArray *[]string) *[]Nft {
	return convertStringArrayToNftArray(strArray)
}

func ConvertNftArrayToStringArray(nftArray *[]Nft) *[]string {
	return convertNftArrayToStringArray(nftArray)
}

func StakeDuelBots(from User, duelBots []Nft, sesssionId ...UUID) error {
	return stakeDuelBots(from, duelBots, sesssionId...)
}

func DistributeFeeToDuelBots(totalFee int64, sessionId ...UUID) (int64, error) {
	return distributeFeeToDuelBots(totalFee, sessionId...)
}

func ClaimDuelBotsRewards(to User, duelBots []Nft, sessionId ...UUID) (int64, error) {
	return claimDuelBotsRewards(to, duelBots, sessionId...)
}

func UnstakeDuelBots(to User, duelBots []Nft, sessionId ...UUID) (int64, error) {
	return unstakeDuelBots(to, duelBots, sessionId...)
}

func DistributeRakeback(to User, feeAmount int64, sessionId ...UUID) (int64, error) {
	return distributeRakeback(to, feeAmount, sessionId...)
}

func ClaimRakeback(user User, sessionId ...UUID) (int64, error) {
	return claimRakeback(user, sessionId...)
}

func GetRakebackInfo(user User, sessionId ...UUID) (models.Rakeback, error) {
	return getRakebackInfo(user, sessionId...)
}

func GenerateRakebackInfo(user User, sessionId ...UUID) (*models.Rakeback, error) {
	return generateRakebackInfo(user, sessionId...)
}

func GetUserRakebackRate(user User) (uint, error) {
	_, rate, err := retrieveRakebackInfoAndRate(user)
	return rate, err
}

func SetActivateAffiliateOnceForRakeback(user User, sessionId ...UUID) (bool, error) {
	return setActivateAffiliateOnceForRakeback(user, sessionId...)
}

func CreateAffiliateCode(user User, codes []string, sessionId ...UUID) error {
	return createAffiliateCode(user, codes, sessionId...)
}

func DeleteAffiliateCode(user User, codes []string, sessionId ...UUID) (int64, error) {
	return deleteAffiliateCode(user, codes, sessionId...)
}

func ClaimAffiliateRewards(user User, codes []string, sessionId ...UUID) (int64, error) {
	return claimAffiliateRewards(user, codes, sessionId...)
}

func ActivateAffiliateCode(user User, code string, sessionId ...UUID) (bool, error) {
	return activateAffiliateCode(user, code, sessionId...)
}

func DeactivateAffiliateCode(user User, code string, sessionId ...UUID) error {
	return deactivateAffiliateCode(user, code, sessionId...)
}

func GetActiveAffiliateCode(user User, sessionId ...UUID) (*ActiveAffiliateMeta, error) {
	return getActiveAffiliateCode(user, sessionId...)
}

func GetOwnedAffiliateCode(user User, sessionId ...UUID) ([]AffiliateMeta, error) {
	return getOwnedAffiliateCode(user, sessionId...)
}

func DistributeAffiliate(from User, feeAmount int64, wagerAmount int64, sessionId ...UUID) (int64, error) {
	return distributeAffiliate(from, feeAmount, wagerAmount, sessionId...)
}

func SetAffiliateCustomRate(code string, customRate uint, sessionId ...UUID) error {
	return setAffiliateCustomRate(code, customRate, sessionId...)
}

func GetAffiliateDetail(code string, sessionId ...UUID) (*AffiliateDetail, error) {
	return getAffiliateDetail(code, sessionId...)
}

// @Internal
// Perform real chip transaction.
func LeaveRealTransaction(
	transaction *models.Transaction,
	sessionId UUID,
) error {
	// 1. Validate parameter
	if transaction == nil {
		return utils.MakeError(
			"coupon_db",
			"leaveRealTransaction",
			"invalid parameter",
			errors.New("provided transaction is nil pointer"),
		)
	}

	// 2. Retrieve session.
	session, err := GetSession(sessionId)
	if err != nil {
		return utils.MakeError(
			"coupon_db",
			"leaveRealTransaction",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Save transaction record.
	if result := session.Create(transaction); result.Error != nil {
		return utils.MakeError(
			"coupon_db",
			"leaveRealTransaction",
			"failed to create a new transaction record",
			result.Error,
		)
	}

	return nil
}
