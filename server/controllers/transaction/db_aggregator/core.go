package db_aggregator

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @External
// Initializes session pointer and other data variable.
func initialize(db *gorm.DB) error {
	if db == nil {
		return utils.MakeError("db_aggregator", "initialize", "", errors.New("invalid db parameter"))
	}

	sessions.Store(UUID_NIL, db)

	return nil
}

// @Internal
// A helpler function that converts pq.int32Array to []Nft.
func convertPqStringArrayToNftArray(pqArray *pq.StringArray) *[]Nft {
	if pqArray == nil {
		return &[]Nft{}
	}
	var nftArray = make([]Nft, len(*pqArray))
	for i, n := 0, len(*pqArray); i < n; i++ {
		nftArray[i] = Nft(string((*pqArray)[i]))
	}
	return &nftArray
}

// @Internal
// A helpler function that converts []Nft to pq.int32Array.
func convertNftArrayToPqStringArray(nftArray *[]Nft) *pq.StringArray {
	if nftArray == nil {
		return &pq.StringArray{}
	}
	var pqArray = make(pq.StringArray, len(*nftArray))
	for i, n := 0, len(*nftArray); i < n; i++ {
		pqArray[i] = string((*nftArray)[i])
	}
	return &pqArray
}

// @External
// A helper function that converts []string to []Nft.
func convertStringArrayToNftArray(strArray *[]string) *[]Nft {
	if strArray == nil {
		return &[]Nft{}
	}
	var nftArray = make([]Nft, len(*strArray))
	for i, n := 0, len(*strArray); i < n; i++ {
		nftArray[i] = Nft((*strArray)[i])
	}
	return &nftArray
}

// @External
// A helper function that converts []string to []Nft.
func convertNftArrayToStringArray(nftArray *[]Nft) *[]string {
	if nftArray == nil {
		return &[]string{}
	}
	var strArray = make([]string, len(*nftArray))
	for i, n := 0, len(*nftArray); i < n; i++ {
		strArray[i] = string((*nftArray)[i])
	}
	return &strArray
}

// @Internal
// This is a helper function to get query args for nft identifier.
func getNftQuerier(nft *Nft) []interface{} {
	return []interface{}{"mint_address = ?", *nft}
}

// @Internal
// This is a helper function to get query args for user identifier.
func getUserQuerier(user *User) []interface{} {
	return []interface{}{*user}
}

// @Internal
// This is a helper function to get query args for wallet identifier.
func getWalletQuerier(wallet *Wallet) []interface{} {
	return []interface{}{*wallet}
}

// @Internal
// This is a helper function to get query args for balance identifier.
func getBalanceQuerier(balance *Balance) []interface{} {
	return []interface{}{*balance}
}

// @Internal
// This is a helper function to get query args for balance identifier.
func getTransactionQuerier(transaction *Transaction) []interface{} {
	return []interface{}{*transaction}
}

// @External
// Gets balance.
func getBalance(balance *Balance, sessionId ...UUID) (*BalanceLoad, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getBalance", "failed to find session", err)
	}
	if balance == nil {
		return nil, utils.MakeError("db_aggregator", "getBalance", "invalid parameter", fmt.Errorf("balance is null pointer"))
	}

	var balanceInfo models.Balance
	if result := session.Preload("ChipBalance").Preload("NftBalance").First(&balanceInfo, getBalanceQuerier(balance)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getBalance", "failed to get balance info", result.Error)
	}

	var result = BalanceLoad{
		Balance: (*Balance)(&balanceInfo.ID),
	}

	if balanceInfo.ChipBalance != nil {
		result.ChipBalance = &balanceInfo.ChipBalance.Balance
	}
	if balanceInfo.NftBalance != nil {
		result.NftBalance = convertPqStringArrayToNftArray(&balanceInfo.NftBalance.Balance)
	}

	return &result, nil
}

// @External
// Gets wallet identifier bound to a user.
func getUserWallet(user *User, lockWallet bool, sessionId ...UUID) (*Wallet, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getUserWallet", "failed to find session", err)
	}
	if user == nil {
		return nil, nil
	}

	var userInfo models.User
	if result := session.Preload("Wallet").First(&userInfo, getUserQuerier(user)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getUserWallet", "failed to get user info", result.Error)
	}
	if userInfo.Wallet.ID == 0 {
		return nil, utils.MakeError("db_aggregator", "getUserWallet", "user's wallet id is 0", errors.New("failed to find user's wallet"))
	}
	if lockWallet {
		if result := session.Clauses(
			clause.Locking{Strength: "UPDATE"},
		).First(&models.Wallet{}, userInfo.Wallet.ID); result.Error != nil {
			return nil, utils.MakeError("db_aggregator", "getUserWallet", "failed to lock wallet info", result.Error)
		}
	}

	return (*Wallet)(&userInfo.Wallet.ID), nil
}

// @External
// Gets wallet's user.
func getWalletUser(wallet *Wallet, sessionId ...UUID) (*User, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getWalletUser", "failed to find session", err)
	}
	if wallet == nil {
		return nil, utils.MakeError("db_aggregator", "getWalletUser", "invalid parameter", fmt.Errorf("wallet is null pointer"))
	}

	var walletInfo models.Wallet
	if result := session.First(&walletInfo, getWalletQuerier(wallet)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getWalletUser", "failed to get wallet info", result.Error)
	}
	if walletInfo.UserID == 0 {
		return nil, utils.MakeError("db_aggregator", "getWalletUser", "wallet's user id is 0", errors.New("failed to find wallet's user"))
	}

	return (*User)(&walletInfo.UserID), nil
}

// @External
// Gets balance identifier bound to a wallet.
func getWalletBalance(wallet *Wallet, sessionId ...UUID) (*Balance, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getWalletBalance", "failed to find session", err)
	}
	if wallet == nil {
		return nil, utils.MakeError("db_aggregator", "getWalletBalance", "invalid parameter", fmt.Errorf("wallet is null pointer"))
	}

	var walletInfo models.Wallet
	if result := session.Preload("Balance").First(&walletInfo, getWalletQuerier(wallet)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getWalletBalance", "failed to get wallet info", result.Error)
	}
	if walletInfo.Balance.ID == 0 {
		return nil, utils.MakeError("db_aggregator", "getWalletBalance", "wallet's balance id is 0", errors.New("failed to find wallet's balance"))
	}

	return (*Balance)(&walletInfo.Balance.ID), nil
}

// @External
// Gets balance identifier bound to user's wallet.
func getUserBalance(user *User, lockWallet bool, sessionId ...UUID) (*Balance, error) {
	wallet, err := getUserWallet(user, lockWallet, sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getUserBalance", "failed to get user wallet", err)
	}

	balance, err := getWalletBalance(wallet, sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getUserBalance", "failed to get wallet balance", err)
	}

	return balance, nil
}

// @External
// Gets wallet in which the nft is stored.
func getNftWallet(nft *Nft, sessionId ...UUID) (*Wallet, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getNftWallet", "failed to find session", err)
	}
	if nft == nil {
		return nil, utils.MakeError("db_aggregator", "getNftWallet", "invalid parameter", fmt.Errorf("nft is null pointer"))
	}

	var nftInfo models.DepositedNft
	if result := session.First(&nftInfo, getNftQuerier(nft)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getNftWallet", "failed to get nft info", result.Error)
	}

	return (*Wallet)(nftInfo.WalletID), nil
}

// @Internal
// This is a helper function to remove Nfts from Nft balance array.
func removeNftsFromBalanceArray(nfts *[]Nft, balanceArray *pq.StringArray) (*pq.StringArray, error) {
	if nfts == nil || len(*nfts) == 0 {
		return balanceArray, nil
	}
	if balanceArray == nil || len(*nfts) > len(*balanceArray) {
		return nil, utils.MakeError("db_aggregator", "removeNftsFromBalanceArray", "invalid parameter", fmt.Errorf("balanceArray is null pointer or too small"))
	}
	totalRemoval := 0
	isRemoval := make([]bool, len(*balanceArray))
	for _, nft := range *nfts {
		for i, nftInBalance := range *balanceArray {
			if (Nft)(nftInBalance) == nft {
				if isRemoval[i] {
					return nil, utils.MakeError("db_aggregator", "removeNftsFromBalanceArray", "", fmt.Errorf("%v nft is not in the ledger", nft))
				}
				isRemoval[i] = true
				totalRemoval += 1
				break
			}
		}
	}

	if totalRemoval != len(*nfts) {
		return nil, utils.MakeError("db_aggregator", "removeNftsFromBalanceArray", "", fmt.Errorf("number of nfts are mismatching"))
	}

	var result = make(pq.StringArray, len(*balanceArray)-len(*nfts))
	if len(result) == 0 {
		return &result, nil
	}
	for i, j, n := 0, 0, len(*balanceArray); i < n; i++ {
		if j == len(result) {
			break
		}
		if !isRemoval[i] {
			result[j] = (*balanceArray)[i]
			j++
		}
	}

	return &result, nil
}

// @Internal
// This is a helper function to remove nfts from ledger of wallet balance.
func removeNftsFromBalanceLedger(nfts *[]Nft, balance *Balance, sessionId ...UUID) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError("db_aggregator", "removeNftsFromBalanceLedger", "invalid parameter", fmt.Errorf("balance is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "removeNftsFromBalanceLedger", "failed to find session", err)
	}

	var balanceInfo models.Balance
	if result := session.Preload("NftBalance").First(&balanceInfo, getBalanceQuerier(balance)...); result.Error != nil {
		return utils.MakeError("db_aggregator", "removeNftsFromBalanceLedger", "failed to get balance info", result.Error)
	}

	updatedBalance, err := removeNftsFromBalanceArray(nfts, &balanceInfo.NftBalance.Balance)
	if err != nil {
		return utils.MakeError("db_aggregator", "removeNftsFromBalanceLedger", "failed to remove nfts from balance array", err)
	}

	var nftBalance models.NftBalance
	if result := session.Model(&nftBalance).Select("balance").Where("id = ?", balanceInfo.NftBalanceID).Updates(models.NftBalance{Balance: *updatedBalance}); result.Error != nil {
		return utils.MakeError("db_aggregator", "removeNftsFromBalanceLedger", "failed to update nft balance", result.Error)
	}

	return nil
}

// @Internal
// This is a helper function to add Nfts to Nft balance array.
func addNftsToBalanceArray(nfts *[]Nft, balanceArray *pq.StringArray) (*pq.StringArray, error) {
	if nfts == nil || len(*nfts) == 0 {
		return balanceArray, nil
	}
	var containedIndex []int
	for _, nft := range *nfts {
		for i, nftInBalance := range *balanceArray {
			if (Nft)(nftInBalance) == nft {
				containedIndex = append(containedIndex, i)
				break
			}
		}
	}

	if len(containedIndex) > 0 {
		return nil, utils.MakeError("db_aggregator", "addNftsToBalanceArray", "", fmt.Errorf("%v nfts are already in the ledger", len(containedIndex)))
	}

	var result = make(pq.StringArray, len(*balanceArray)+len(*nfts))
	copy(result, *balanceArray)
	for i, nft := range *nfts {
		result[i+len(*balanceArray)] = (string)(nft)
	}

	return &result, nil
}

// @Internal
// This is a helper function to add nfts to ledger of wallet balance.
func addNftsToBalanceLedger(nfts *[]Nft, balance *Balance, sessionId ...UUID) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError("db_aggregator", "addNftsToBalanceLedger", "invalid parameter", fmt.Errorf("balance is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "addNftsToBalanceLedger", "failed to find session", err)
	}

	var balanceInfo models.Balance
	if result := session.Preload("NftBalance").First(&balanceInfo, getBalanceQuerier(balance)...); result.Error != nil {
		return utils.MakeError("db_aggregator", "addNftsToBalanceLedger", "failed to get balance info", result.Error)
	}

	updatedBalance, err := addNftsToBalanceArray(nfts, &balanceInfo.NftBalance.Balance)
	if err != nil {
		return utils.MakeError("db_aggregator", "addNftsToBalanceLedger", "add nfts to balance array", err)
	}

	var nftBalance models.NftBalance
	if result := session.Model(&nftBalance).Select("balance").Where("id = ?", balanceInfo.NftBalanceID).Updates(models.NftBalance{Balance: *updatedBalance}); result.Error != nil {
		return utils.MakeError("db_aggregator", "addNftsToBalanceLedger", "failed to update nft balance", result.Error)
	}

	return nil
}

// @Internal
// Removes wallet id from nft. In other word, takes an nft out of the stored wallet.
func removeNftsFromWallet(nfts *[]Nft, wallet *Wallet, sessionId ...UUID) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if wallet == nil {
		return utils.MakeError("db_aggregator", "removeNftsFromWallet", "invalid parameter", fmt.Errorf("wallet is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "removeNftsFromWallet", "failed to find session", err)
	}

	for _, nft := range *nfts {
		var nftInfo models.DepositedNft
		if result := session.First(&nftInfo, getNftQuerier(&nft)...); result.Error != nil {
			return utils.MakeError("db_aggregator", "removeNftsFromWallet", "failed to get nft info", result.Error)
		}
		if nftInfo.WalletID == nil {
			return utils.MakeError("db_aggregator", "removeNftsFromWallet", "", fmt.Errorf("something wrong with nft(%v) ownership", nft))
		}
		if (Wallet)(*nftInfo.WalletID) != *wallet {
			return utils.MakeError("db_aggregator", "removeNftsFromWallet", "", fmt.Errorf("the nft(%v) is not stored in wallet(%v)", nft, *wallet))
		}

		nftInfo.WalletID = nil
		if result := session.Save(&nftInfo); result.Error != nil {
			return utils.MakeError("db_aggregator", "removeNftsFromWallet", "failed to update nft wallet", result.Error)
		}
	}

	return nil
}

// @Internal
// Adds wallet id to nft. In other word, puts an nft in to the wallet.
func addNftsToWallet(nfts *[]Nft, wallet *Wallet, sessionId ...UUID) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if wallet == nil {
		return utils.MakeError("db_aggregator", "addNftsToWallet", "invalid parameter", fmt.Errorf("wallet is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "addNftsToWallet", "failed to find session", err)
	}

	for _, nft := range *nfts {
		var nftInfo models.DepositedNft
		if result := session.First(&nftInfo, getNftQuerier(&nft)...); result.Error != nil {
			return utils.MakeError("db_aggregator", "addNftsToWallet", "failed to get nft info", result.Error)
		}
		if nftInfo.WalletID != nil {
			return utils.MakeError("db_aggregator", "addNftsToWallet", "", fmt.Errorf("the nft(%T) is stored in another wallet", nft))
		}

		nftInfo.WalletID = (*uint)(wallet)
		if result := session.Save(&nftInfo); result.Error != nil {
			return utils.MakeError("db_aggregator", "addNftsToWallet", "failed to update nft wallet", result.Error)
		}
	}

	return nil
}

// @Internal
// Removes nfts from user.
func removeNftsFromUser(user *User, nfts *[]Nft, sessionId ...UUID) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if user == nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromUser",
			"invalid parameter",
			fmt.Errorf("user is a null pointer"),
		)
	}
	wallet, err := getUserWallet(user, true, sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromUser",
			"failed to get user wallet",
			err,
		)
	}

	balance, err := getWalletBalance(wallet, sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromUser",
			"failed to get wallet balance",
			err,
		)
	}

	if err := removeNftsFromBalanceAndWalletUnchecked(
		balance,
		wallet,
		nfts,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromUser",
			"failed to remove nfts from balance and wallet",
			err,
		)
	}

	return nil
}

// @Internal
// Removes nfts from user.
func removeNftsFromBalanceAndWalletUnchecked(
	balance *Balance,
	wallet *Wallet,
	nfts *[]Nft,
	sessionId ...UUID,
) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromBalanceAndWalletUnchecked",
			"invalid parameter",
			fmt.Errorf("balance is a null pointer"),
		)
	}
	if wallet == nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromBalanceAndWalletUnchecked",
			"invalid parameter",
			fmt.Errorf("wallet is a null pointer"),
		)
	}

	if err := removeNftsFromBalanceLedger(
		nfts,
		balance,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromBalanceAndWalletUnchecked",
			"failed to remove nfts from balance",
			err,
		)
	}

	if err := removeNftsFromWallet(
		nfts,
		wallet,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeNftsFromBalanceAndWalletUnchecked",
			"failed to remove nfts from wallet",
			err,
		)
	}

	return nil
}

// @Internal
// add nfts to user.
func addNftsToUser(
	user *User,
	nfts *[]Nft,
	sessionId ...UUID,
) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if user == nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToUser",
			"invalid parameter",
			fmt.Errorf("user is a null pointer"),
		)
	}
	wallet, err := getUserWallet(
		user,
		true,
		sessionId...,
	)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToUser",
			"failed to get user wallet",
			err,
		)
	}

	balance, err := getWalletBalance(wallet, sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToUser",
			"failed to get wallet balance",
			err,
		)
	}

	if err := addNftsToBalanceAndWalletUnchecked(
		balance,
		wallet,
		nfts,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToUser",
			"failed to add nfts to balance and walelt",
			err,
		)
	}

	return nil
}

// @Internal
// add nfts to user.
func addNftsToBalanceAndWalletUnchecked(
	balance *Balance,
	wallet *Wallet,
	nfts *[]Nft,
	sessionId ...UUID,
) error {
	if nfts == nil || len(*nfts) == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToBalanceAndWalletUnchecked",
			"invalid parameter",
			fmt.Errorf("balance is a null pointer"),
		)
	}
	if wallet == nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToBalanceAndWalletUnchecked",
			"invalid parameter",
			fmt.Errorf("wallet is a null pointer"),
		)
	}

	if err := addNftsToBalanceLedger(
		nfts,
		balance,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToBalanceAndWalletUnchecked",
			"failed to add nfts to balance",
			err,
		)
	}

	if err := addNftsToWallet(
		nfts,
		wallet,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addNftsToBalanceAndWalletUnchecked",
			"failed to add nfts to wallet",
			err,
		)
	}

	return nil
}

// @Internal
// Adds user chip balance.
func addChipsToUser(user *User, chips int64, sessionId ...UUID) error {
	if user == nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToUser",
			"invalid parameter",
			fmt.Errorf("user is a null pointer"),
		)
	}

	balance, err := getUserBalance(user, true, sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToUser",
			"failed to get user balance",
			err,
		)
	}

	if err := addChipsToBalance(
		balance,
		chips,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToUser",
			"failed to add chips to balance",
			err,
		)
	}

	return nil
}

// @Internal
// Adds user chip balance.
func addChipsToBalance(
	balance *Balance,
	chips int64,
	sessionId ...UUID,
) error {
	if chips < 0 {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToBalance",
			"",
			errors.New("negative chips amount"),
		)
	}
	if chips == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToBalance",
			"invalid parameter",
			fmt.Errorf("user is a null pointer"),
		)
	}

	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToBalance",
			"failed to find session",
			err,
		)
	}

	var balanceInfo models.Balance
	if result := session.First(
		&balanceInfo,
		getBalanceQuerier(balance)...,
	); result.Error != nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToBalance",
			"get balance info",
			result.Error,
		)
	}

	var chipBalanceInfo models.ChipBalance
	if result := session.Model(
		&chipBalanceInfo,
	).Where(
		"id = ?",
		balanceInfo.ChipBalanceID,
	).Update(
		"balance",
		gorm.Expr("balance + ?", chips),
	); result.Error != nil {
		return utils.MakeError(
			"db_aggregator",
			"addChipsToBalance",
			"failed to update chip balance",
			result.Error,
		)
	}

	return nil
}

// @Internal
// Removes user chip balance
func removeChipsFromUser(user *User, chips int64, sessionId ...UUID) error {
	if user == nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromUser",
			"invalid parameter",
			fmt.Errorf("user is a null pointer"),
		)
	}

	balance, err := getUserBalance(user, true, sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromUser",
			"failed to get user balance",
			err,
		)
	}

	if err := removeChipsFromBalance(
		balance,
		chips,
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromUser",
			"failed to remove chips from balance",
			err,
		)
	}

	return nil
}

// @Internal
// Removes user chip balance
func removeChipsFromBalance(
	balance *Balance,
	chips int64,
	sessionId ...UUID,
) error {
	if chips < 0 {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"",
			errors.New("negative chips amount"),
		)
	}
	if chips == 0 {
		return nil
	}
	if balance == nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"invalid parameter",
			fmt.Errorf("balance is a null pointer"),
		)
	}

	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"failed to find session",
			err,
		)
	}

	var balanceInfo models.Balance
	if result := session.First(
		&balanceInfo,
		getBalanceQuerier(balance)...,
	); result.Error != nil {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"failed to get balance info",
			result.Error,
		)
	}

	var chipBalanceInfo = models.ChipBalance{
		Model: gorm.Model{
			ID: *balanceInfo.ChipBalanceID,
		},
	}
	if result := session.Model(
		&chipBalanceInfo,
	).Clauses(
		clause.Returning{},
	).Update(
		"balance",
		gorm.Expr("balance - ?", chips),
	); result.Error != nil ||
		result.RowsAffected == 0 {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"failed to update chip balance",
			result.Error,
		)
	}
	if chipBalanceInfo.Balance < 0 {
		return utils.MakeError(
			"db_aggregator",
			"removeChipsFromBalance",
			"",
			fmt.Errorf(
				"insufficient funds. balanceID: %d, removingChips: %d, afterRemoving: %d",
				*balance, chips, chipBalanceInfo.Balance,
			),
		)
	}

	return nil
}

// @External
// Record a pending transaction with the provided details.
// Details required are listed below.
// - FromWallet
// - ToWallet
// - Type
// - Balance
// - FromWalletPrevID
// - FromWalletNextID
// ~ Status is set as "pending"
func recordTransaction(transactionLoad *TransactionLoad, sessionId ...UUID) (*Transaction, error) {
	if transactionLoad == nil {
		return nil, utils.MakeError("db_aggregator", "recordTransaction", "invalid parameter", fmt.Errorf("transactionLoad is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "recordTransaction", "failed to find session", err)
	}

	var receipients pq.Int64Array
	if transactionLoad.Receipients != nil {
		for _, receipient := range *transactionLoad.Receipients {
			receipients = append(receipients, int64(receipient))
		}
	}

	var transactionInfo = models.Transaction{
		FromWallet:       (*uint)(transactionLoad.FromWallet),
		ToWallet:         (*uint)(transactionLoad.ToWallet),
		Type:             transactionLoad.Type,
		Status:           models.TransactionStatus(models.TransactionPending),
		FromWalletPrevID: (*uint)(transactionLoad.FromWalletPrevID),
		FromWalletNextID: (*uint)(transactionLoad.FromWalletNextID),
		Receipients:      receipients,
	}
	if transactionLoad.Balance.ChipBalance != nil {
		transactionInfo.Balance.ChipBalance = &models.ChipBalance{
			Balance: *transactionLoad.Balance.ChipBalance,
		}
	}
	if transactionLoad.Balance.NftBalance != nil {
		transactionInfo.Balance.NftBalance = &models.NftBalance{
			Balance: *convertNftArrayToPqStringArray(transactionLoad.Balance.NftBalance),
		}
	}

	if result := session.Create(&transactionInfo); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "recordTransaction", "failed to create transaction record", result.Error)
	}

	return (*Transaction)(&transactionInfo.ID), nil
}

// @External
// Make a pending transaction as succeed one with the provided details.
// Details required are listed below.
// - ToWalletPrevID
// - ToWalletNextID
// ~ Status is set as "succeed"
func confirmTransaction(transactionLoad *TransactionLoad, transaction *Transaction, sessionId ...UUID) error {
	if transactionLoad == nil {
		return utils.MakeError("db_aggregator", "confirmTransaction", "invalid parameter", fmt.Errorf("transactionLoad is a null pointer"))
	}
	if transaction == nil {
		return utils.MakeError("db_aggregator", "confirmTransaction", "invalid parameter", fmt.Errorf("transaction is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "confirmTransaction", "failed to find session", err)
	}

	var transactionInfo models.Transaction
	if result := session.First(&transactionInfo, getTransactionQuerier(transaction)...); result.Error != nil {
		return utils.MakeError("db_aggregator", "confirmTransaction", "failed to find transaction info", result.Error)
	}
	if transactionInfo.Status != models.TransactionStatus(models.TransactionPending) {
		return utils.MakeError("db_aggregator", "confirmTransaction", "", errors.New("transaction's status is not pending"))
	}
	transactionInfo.Status = models.TransactionStatus(models.TransactionSucceed)
	transactionInfo.ToWalletPrevID = (*uint)(transactionLoad.ToWalletPrevID)
	transactionInfo.ToWalletNextID = (*uint)(transactionLoad.ToWalletNextID)
	transactionInfo.OwnerID = transactionLoad.OwnerID
	transactionInfo.OwnerType = transactionLoad.OwnerType
	if result := session.Model(&transactionInfo).Select(
		"status",
		"to_wallet_prev_id",
		"to_wallet_next_id",
		"owner_id",
		"owner_type",
	).Save(&transactionInfo); result.Error != nil {
		return utils.MakeError("db_aggregator", "confirmTransaction", "failed to update transaction", result.Error)
	}

	return nil
}

// @External
// Make a pending transaction as failed one with the provided details.
// For the failed transaction, ToWalletPrevID/ToWalletNextID save the balances before/after refunding.
// Details required are listed below.
// - ToWalletPrevID
// - ToWalletNextID
// ~ Status is set as "failed"
func declineTransaction(transactionLoad *TransactionLoad, transaction *Transaction, sessionId ...UUID) error {
	if transactionLoad == nil {
		return utils.MakeError("db_aggregator", "declineTransaction", "invalid parameter", fmt.Errorf("transactionLoad is a null pointer"))
	}
	if transaction == nil {
		return utils.MakeError("db_aggregator", "declineTransaction", "invalid parameter", fmt.Errorf("transaction is a null pointer"))
	}
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError("db_aggregator", "declineTransaction", "failed to find session", err)
	}

	var transactionInfo models.Transaction
	if result := session.First(&transactionInfo, getTransactionQuerier(transaction)...); result.Error != nil {
		return result.Error
	}
	if transactionInfo.Status != models.TransactionStatus(models.TransactionPending) {
		return utils.MakeError("db_aggregator", "declineTransaction", "failed to get transaction info", errors.New("transaction's status is not pending"))
	}
	transactionInfo.Status = models.TransactionStatus(models.TransactionFailed)
	transactionInfo.RefundPrevID = (*uint)(transactionLoad.RefundPrevID)
	transactionInfo.RefundNextID = (*uint)(transactionLoad.RefundNextID)
	transactionInfo.OwnerID = transactionLoad.OwnerID
	transactionInfo.OwnerType = transactionLoad.OwnerType
	if result := session.Model(&transactionInfo).Select(
		"status",
		"refund_prev_id",
		"refund_next_id",
		"owner_id",
		"owner_type",
	).Save(&transactionInfo); result.Error != nil {
		return utils.MakeError("db_aggregator", "declineTransaction", "failed to update transaction", result.Error)
	}

	return nil
}

// @External
// Makes a copy of a balance. The balance should be attached to a wallet.
// For the created one, Owner is the wallet and type is in-history.
func recordBalanceHistory(
	balance *Balance,
	changedBalance ChangedBalance,
	sessionId ...UUID,
) (*BalanceHistoryChain, error) {
	if balance == nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"invalid parameter",
			fmt.Errorf("balance is a null pointer"),
		)
	}

	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"failed to find session",
			err,
		)
	}

	var balanceInfo models.Balance
	if result := session.Preload(
		"ChipBalance",
	).Preload(
		"NftBalance",
	).First(
		&balanceInfo,
		getBalanceQuerier(balance)...,
	); result.Error != nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"failed to get balance info",
			result.Error,
		)
	}
	if balanceInfo.OwnerType != models.BalanceOwnerType(models.InWallet) {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"",
			errors.New("not a balance attached to a wallet"),
		)
	}

	balanceInfo.OwnerType = models.BalanceOwnerType(models.InHistory)
	if result := session.Save(&balanceInfo); result.Error != nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"failed to update balance owner",
			result.Error,
		)
	}

	var chipBalanceID = balanceInfo.ChipBalanceID
	var nftBalanceID = balanceInfo.NftBalanceID
	var newBalanceInfo = models.Balance{
		OwnerID:   balanceInfo.OwnerID,
		OwnerType: models.BalanceOwnerType(models.InWallet),
	}

	if changedBalance.ChipBalance {
		newBalanceInfo.ChipBalance = &models.ChipBalance{
			Balance: balanceInfo.ChipBalance.Balance,
		}
	} else {
		newBalanceInfo.ChipBalanceID = chipBalanceID
	}

	if changedBalance.NftBalance {
		newBalanceInfo.NftBalance = &models.NftBalance{
			Balance: balanceInfo.NftBalance.Balance,
		}
	} else {
		newBalanceInfo.NftBalanceID = nftBalanceID
	}

	if result := session.Create(&newBalanceInfo); result.Error != nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"recordBalanceHistory",
			"failed to update balance",
			result.Error,
		)
	}

	return &BalanceHistoryChain{
		Prev: (*Balance)(&balanceInfo.ID),
		Next: (*Balance)(&newBalanceInfo.ID),
	}, nil
}

// @Internal
// Removes balances from a user.
func removeBalanceFromUser(from *User, balanceLoad *BalanceLoad, sessionId ...UUID) error {
	if balanceLoad.ChipBalance != nil {
		if err := removeChipsFromUser(from, *balanceLoad.ChipBalance, sessionId...); err != nil {
			return utils.MakeError("db_aggregator", "removeBalanceFromUser", "failed to remove chips from user", err)
		}
	}
	if balanceLoad.NftBalance != nil {
		if err := removeNftsFromUser(from, balanceLoad.NftBalance, sessionId...); err != nil {
			return utils.MakeError("db_aggregator", "removeBalanceFromUser", "failed to remove nfts from user", err)
		}
	}
	return nil
}

// @Internal
// Removes balances from balance.
func removeBalanceFromBalanceAndWalletUnchecked(
	balance *Balance,
	wallet *Wallet,
	balanceLoad *BalanceLoad,
	sessionId ...UUID,
) error {
	if balanceLoad.ChipBalance != nil {
		if err := removeChipsFromBalance(
			balance,
			*balanceLoad.ChipBalance,
			sessionId...,
		); err != nil {
			return utils.MakeError(
				"db_aggregator",
				"removeBalanceFromUser",
				"failed to remove chips from user",
				err,
			)
		}
	}
	if balanceLoad.NftBalance != nil {
		if err := removeNftsFromBalanceAndWalletUnchecked(
			balance,
			wallet,
			balanceLoad.NftBalance,
			sessionId...,
		); err != nil {
			return utils.MakeError(
				"db_aggregator",
				"removeBalanceFromUser",
				"failed to remove nfts from user",
				err,
			)
		}
	}
	return nil
}

// @Internal
// Adds balances to a user.
func addBalanceToUser(to *User, balanceLoad *BalanceLoad, sessionId ...UUID) error {
	if balanceLoad.ChipBalance != nil {
		if err := addChipsToUser(to, *balanceLoad.ChipBalance, sessionId...); err != nil {
			return utils.MakeError("db_aggregator", "addBalanceToUser", "failed to add chips to user", err)
		}
	}
	if balanceLoad.NftBalance != nil {
		if err := addNftsToUser(to, balanceLoad.NftBalance, sessionId...); err != nil {
			return utils.MakeError("db_aggregator", "addBalanceToUser", "failed to add nfts to user", err)
		}
	}
	return nil
}

// @Internal
// Adds balances to a user.
func addBalanceToBalanceAndWalletUnchecked(
	balance *Balance,
	wallet *Wallet,
	balanceLoad *BalanceLoad,
	sessionId ...UUID,
) error {
	if balanceLoad.ChipBalance != nil {
		if err := addChipsToBalance(
			balance,
			*balanceLoad.ChipBalance,
			sessionId...,
		); err != nil {
			return utils.MakeError(
				"db_aggregator",
				"addBalanceToBalanceAndWalletUnchecked",
				"failed to add chips to user",
				err,
			)
		}
	}
	if balanceLoad.NftBalance != nil {
		if err := addNftsToBalanceAndWalletUnchecked(
			balance,
			wallet,
			balanceLoad.NftBalance,
			sessionId...,
		); err != nil {
			return utils.MakeError(
				"db_aggregator",
				"addBalanceToBalanceAndWalletUnchecked",
				"failed to add nfts to user",
				err,
			)
		}
	}
	return nil
}

// @Internal
// Check balance presence. Checks for only `chipBalance`, and `nftBalance`.
func balancePresence(balance *BalanceLoad) (*ChangedBalance, *ChangedBalance, bool, error) {
	var fromChangedBalance = ChangedBalance{
		ChipBalance: false,
		NftBalance:  false,
	}
	var toChangedBalance = ChangedBalance{
		ChipBalance: false,
		NftBalance:  false,
	}
	var totalPresence = false

	if balance == nil {
		return &fromChangedBalance, &toChangedBalance, totalPresence, nil
	}

	if balance.ChipBalance != nil && *(balance.ChipBalance) > 0 {
		fromChangedBalance.ChipBalance = true
		toChangedBalance.ChipBalance = true
		totalPresence = true
	}

	if balance.NftBalance != nil && len(*(balance.NftBalance)) > 0 {
		fromChangedBalance.NftBalance = true
		toChangedBalance.NftBalance = true
		totalPresence = true
	}

	return &fromChangedBalance, &toChangedBalance, totalPresence, nil
}

// @External
// Transfer Balances from user to another
func transfer(
	from *User,
	to *User,
	balanceLoad *BalanceLoad,
	sessionId ...UUID,
) (*TransferResult, error) {
	fromChangedBalance, toChangedBalance, totalPresence, err := balancePresence(balanceLoad)
	if err != nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"transfer",
			"failed to get balance presence",
			err,
		)
	}
	if !totalPresence {
		return nil, utils.MakeError(
			"db_aggregator",
			"transfer",
			"invalid parameter",
			errors.New("empty balance to be transfered"),
		)
	}
	if from == nil && to == nil {
		return nil, utils.MakeError(
			"db_aggregator",
			"transfer",
			"invalid parameter",
			errors.New("both sender/recipient are null"),
		)
	}

	var result = TransferResult{}
	// Lock from/to wallets before move funds.
	if from != nil {
		wallet, err := getUserWallet(
			from,
			true,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to lock and retrieve from user wallet",
				fmt.Errorf(
					"userId: %d, err: %v",
					*from, err,
				),
			)
		}
		result.FromWallet = wallet
	}
	if to != nil {
		wallet, err := getUserWallet(
			to,
			true,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to lock and retrieve to user wallet",
				fmt.Errorf(
					"userId: %d, err: %v",
					*to, err,
				),
			)
		}
		result.ToWallet = wallet
	}

	if from != nil {
		balance, err := getWalletBalance(
			result.FromWallet,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to get from user balance",
				err,
			)
		}

		balanceChain, err := recordBalanceHistory(
			balance,
			*fromChangedBalance,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to record balance of from user",
				err,
			)
		}
		balance = balanceChain.Next
		result.FromPrevBalance = balanceChain.Prev
		result.FromNextBalance = balanceChain.Next

		if err := removeBalanceFromBalanceAndWalletUnchecked(
			balance,
			result.FromWallet,
			(*BalanceLoad)(balanceLoad),
			sessionId...,
		); err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to remove from user's balance",
				fmt.Errorf(
					"user ID: %d, error: %v",
					from, err,
				),
			)
		}
	}

	if to != nil {
		balance, err := getWalletBalance(
			result.ToWallet,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to get to user balance",
				err,
			)
		}

		balanceChain, err := recordBalanceHistory(
			balance,
			*toChangedBalance,
			sessionId...,
		)
		if err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to record balance of to user",
				err,
			)
		}
		balance = balanceChain.Next
		result.ToPrevBalance = balanceChain.Prev
		result.ToNextBalance = balanceChain.Next

		if err := addBalanceToBalanceAndWalletUnchecked(
			balance,
			result.ToWallet,
			(*BalanceLoad)(balanceLoad),
			sessionId...,
		); err != nil {
			return nil, utils.MakeError(
				"db_aggregator",
				"transfer",
				"failed to remove to user's balance",
				err,
			)
		}
	}

	return &result, nil
}

// @External
// Burn chips from user's wallet.
func burn(from *User, amount int64, sessionId ...UUID) error {
	if from == nil {
		return utils.MakeError(
			"db_aggregator",
			"burn",
			"invalid parameter",
			errors.New("provided from wallet is nil pointer"),
		)
	}

	if amount == 0 {
		return nil
	}

	if err := removeChipsFromUser(from, amount, sessionId...); err != nil {
		return utils.MakeError(
			"db_aggregator",
			"burn",
			"failed to remove chips from user",
			err,
		)
	}

	return nil
}

// External
// Transfer Balances from user to others
func rain(from *User, receipients *[]User, balanceLoad *BalanceLoad, sessionId ...UUID) (*RainResult, error) {
	fromChangedBalance, toChangedBalance, totalPresence, err := balancePresence(balanceLoad)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "rain", "failed to get balance presence", err)
	}
	if !totalPresence {
		return nil, utils.MakeError("db_aggregator", "rain", "", errors.New("empty balance to be transfered"))
	}
	if from == nil || receipients == nil || len(*receipients) == 0 {
		return nil, utils.MakeError(
			"db_aggregator",
			"rain",
			"invalid parameter",
			fmt.Errorf("from: %v, to: %v", from, receipients),
		)
	}
	if balanceLoad == nil || balanceLoad.ChipBalance == nil || *balanceLoad.ChipBalance == 0 {
		return nil, utils.MakeError(
			"db_aggregator",
			"rain",
			"invalid parameter",
			errors.New("chip balance should not be null"),
		)
	}

	var result = RainResult{}

	balance, err := getUserBalance(from, true, sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "rain", "failed to get from user balance", err)
	}
	balanceChain, err := recordBalanceHistory(balance, ChangedBalance(*fromChangedBalance), sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "rain", "failed to record balance of from user", err)
	}
	result.FromPrevBalance = balanceChain.Prev
	result.FromNextBalance = balanceChain.Next

	totalChipBalance := *balanceLoad.ChipBalance * int64(len(*receipients))

	if err := removeBalanceFromUser(from, &BalanceLoad{ChipBalance: &totalChipBalance}, sessionId...); err != nil {
		return nil, utils.MakeError("db_aggregator", "rain", "failed to remove from user's balance", err)
	}

	for _, to := range *receipients {
		balance, err := getUserBalance(&to, true, sessionId...)
		if err != nil {
			return nil, utils.MakeError("db_aggregator", "transfer", "failed to get to user balance", err)
		}

		_, err = recordBalanceHistory(balance, ChangedBalance(*toChangedBalance), sessionId...)
		if err != nil {
			return nil, utils.MakeError("db_aggregator", "transfer", "failed to record balance of to user", err)
		}

		if err := addBalanceToUser(&to, (*BalanceLoad)(balanceLoad), sessionId...); err != nil {
			return nil, utils.MakeError("db_aggregator", "transfer", "failed to remove to user's balance", err)
		}
	}

	return &result, nil
}

// @External
// Get transaction load.
func getTransactionLoad(transaction *Transaction, sessionId ...UUID) (*TransactionLoad, error) {
	session, err := getSession(sessionId...)
	if err != nil {
		return nil, utils.MakeError("db_aggregator", "getTransactionLoad", "failed to find session", err)
	}

	var transactionInfo models.Transaction
	if result := session.Preload("Balance.ChipBalance").Preload("Balance.NftBalance").First(&transactionInfo, getTransactionQuerier(transaction)...); result.Error != nil {
		return nil, utils.MakeError("db_aggregator", "getTransactionLoad", "failed to get transaction info", result.Error)
	}

	transactionLoad := TransactionLoad{
		FromWallet: (*Wallet)(transactionInfo.FromWallet),
		ToWallet:   (*Wallet)(transactionInfo.ToWallet),
		Balance: BalanceLoad{
			Balance: (*Balance)(&transactionInfo.Balance.ID),
		},
		Type:             transactionInfo.Type,
		Status:           transactionInfo.Status,
		FromWalletPrevID: (*Balance)(transactionInfo.FromWalletPrevID),
		FromWalletNextID: (*Balance)(transactionInfo.FromWalletNextID),
		ToWalletPrevID:   (*Balance)(transactionInfo.ToWalletPrevID),
		ToWalletNextID:   (*Balance)(transactionInfo.ToWalletNextID),
	}
	if transactionInfo.Balance.ChipBalance != nil {
		transactionLoad.Balance.ChipBalance = &transactionInfo.Balance.ChipBalance.Balance
	}
	if transactionInfo.Balance.NftBalance != nil {
		transactionLoad.Balance.NftBalance = convertPqStringArrayToNftArray(&transactionInfo.Balance.NftBalance.Balance)
	}
	return &transactionLoad, nil
}
