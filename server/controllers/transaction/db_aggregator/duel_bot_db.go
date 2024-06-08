package db_aggregator

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/** Record staked duelbots.
* - Update status to `staked`
* - Update staking_user_id
* - Errors for already staked DuelBots
* - Errors for DuelBots not owned by staking user
**/
func recordStakedDuelBotsChecked(from User, duelBots []Nft, sessionId ...UUID) error {
	// 1. Validate parameters.
	if from <= 0 || len(duelBots) == 0 {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"invalid parameter",
			errors.New("0 duelbots, or 0 from user"),
		)
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"failed to get session",
			err,
		)
	}

	// 3. Retrieve Deposited Nft IDs
	duelBotNfts := []models.DepositedNft{}
	if result := session.Where(
		"mint_address in ?",
		duelBots,
	).Find(&duelBotNfts); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"failed to retrieve duel bot nfts",
			result.Error,
		)
	}

	if len(duelBotNfts) != len(duelBots) {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"provided duelbots mismatching with retrieved",
			fmt.Errorf(
				"provided: %v, retrieved: %v",
				len(duelBots),
				len(duelBotNfts),
			),
		)
	}

	duelBotNftIDs := []uint{}
	for _, nft := range duelBotNfts {
		duelBotNftIDs = append(duelBotNftIDs, nft.ID)
	}

	// 4. Update staked DuelBot records.
	duelBotRecords := []models.DuelBot{}
	if result := session.Preload(
		"DepositedNft",
	).Where(
		"status = ?",
		models.DuelBotNormal,
	).Where(
		"staking_user_id is null",
	).Where(
		"deposited_nft_id in ?",
		duelBotNftIDs,
	).Find(&duelBotRecords); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"failed to retrieve duel bot records",
			result.Error,
		)
	}

	if len(duelBotRecords) != len(duelBots) {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"retrieved duel bot records count mismatching",
			fmt.Errorf(
				"requested count: %d, actual retrieved: %d",
				len(duelBots),
				len(duelBotRecords),
			),
		)
	}

	duelStakeWallet := models.Wallet{}
	if result := session.Where(
		"user_id = ?", config.DUEL_BOT_STAKE_ID,
	).Find(&duelStakeWallet); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"failed to get duel stake user's wallet",
			result.Error,
		)
	}

	for i := range duelBotRecords {
		if *duelBotRecords[i].DepositedNft.WalletID != duelStakeWallet.ID {
			return utils.MakeError(
				"duel_bot_db",
				"recordStakedDuelBotsChecked",
				"not staked duelbot",
				errors.New("duel bot not owned by stake user"),
			)
		}
		duelBotRecords[i].Status = models.DuelBotStaked
		duelBotRecords[i].StakingUserID = (*uint)(&from)
	}

	if result := session.Save(&duelBotRecords); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordStakedDuelBotsChecked",
			"failed to update status and staking user",
			result.Error,
		)
	}

	return nil
}

/** Record unstaked duelbots.
* - Update status to `normal`
* - Update staking_user_id
* - Errors for already unstaked DuelBots
* - Errors for DuelBots not owned by user
**/
func recordUnstakedDuelBotsChecked(to User, duelBots []Nft, sessionId ...UUID) error {
	// 1. Validate parameters.
	if to <= 0 || len(duelBots) == 0 {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"invalid parameter",
			errors.New("0 duelbots, or 0 to user"),
		)
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"failed to get session",
			err,
		)
	}

	// 3. Retrieve Deposited Nft IDs
	duelBotNfts := []models.DepositedNft{}
	if result := session.Where(
		"mint_address in ?",
		duelBots,
	).Find(&duelBotNfts); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"failed to retrieve duel bot nfts",
			result.Error,
		)
	}

	if len(duelBotNfts) != len(duelBots) {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"provided duelbots mismatching with retrieved",
			fmt.Errorf(
				"provided: %v, retrieved: %v",
				len(duelBots),
				len(duelBotNfts),
			),
		)
	}

	duelBotNftIDs := []uint{}
	for _, nft := range duelBotNfts {
		duelBotNftIDs = append(duelBotNftIDs, nft.ID)
	}

	// 4. Update unstaked DuelBot records.
	duelBotRecords := []models.DuelBot{}
	if result := session.Preload(
		"DepositedNft",
	).Where(
		"status = ?",
		models.DuelBotStaked,
	).Where(
		"staking_user_id = ?",
		to,
	).Where(
		"deposited_nft_id in ?",
		duelBotNftIDs,
	).Find(&duelBotRecords); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"failed to retrieve duel bot records",
			result.Error,
		)
	}

	if len(duelBotRecords) != len(duelBots) {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"retrieved duel bot records count mismatching",
			fmt.Errorf(
				"requested count: %d, actual retrieved: %d",
				len(duelBots),
				len(duelBotRecords),
			),
		)
	}

	userInfo := models.User{}
	if result := session.Preload("Wallet").First(&userInfo, to); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"failed to retrieve user info",
			result.Error,
		)
	}

	for i := range duelBotRecords {
		if *duelBotRecords[i].DepositedNft.WalletID != userInfo.Wallet.ID {
			return utils.MakeError(
				"duel_bot_db",
				"recordUnstakedDuelBotsChecked",
				"not staked duelbot",
				errors.New("duel bot not owned by user"),
			)
		}
		duelBotRecords[i].Status = models.DuelBotNormal
		duelBotRecords[i].StakingUserID = nil
	}

	if result := session.Save(&duelBotRecords); result.Error != nil {
		return utils.MakeError(
			"duel_bot_db",
			"recordUnstakedDuelBotsChecked",
			"failed to update status and staking user",
			result.Error,
		)
	}

	return nil
}

func stakeDuelBots(from User, duelBots []Nft, sessionId ...UUID) error {
	// 1. Validate parameters.
	if from <= 0 || len(duelBots) == 0 {
		return utils.MakeError(
			"duel_bot_db",
			"stakeDuelBots",
			"invalid parameter",
			errors.New("0 duelbots, or 0 from user"),
		)
	}

	// 2. Try transfer duelbots to duel stake user.
	duelStake := User(config.DUEL_BOT_STAKE_ID)
	if _, err := transfer(
		&from,
		&duelStake,
		&BalanceLoad{
			NftBalance: &duelBots,
		},
		sessionId...,
	); err != nil {
		return utils.MakeError(
			"duel_bot_db",
			"stakeDuelBots",
			"failed to transfer duel bot",
			err,
		)
	}

	// 3. Try update staked user of duel bot.
	if err := recordStakedDuelBotsChecked(from, duelBots, sessionId...); err != nil {
		return utils.MakeError(
			"duel_bot_db",
			"stakeDuelBots",
			"failed to update staked user",
			err,
		)
	}

	return nil
}

func distributeFeeToDuelBots(totalFee int64, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if totalFee <= 0 {
		return 0, utils.MakeError(
			"duel_bot_db",
			"distributeFeeToDuelBots",
			"invalid parameter",
			fmt.Errorf("totalFee: %d", totalFee),
		)
	}

	// 2. Retrieve session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"distributeFeeToDuelBots",
			"failed to get session",
			err,
		)
	}

	// 3. Get total staked count.
	stakedCount := int64(0)
	if result := session.Model(
		&models.DuelBot{},
	).Where(
		"status != ?",
		models.DuelBotNormal,
	).Count(&stakedCount); result.Error != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"distributeFeeToDuelBots",
			"failed to retrieve staked count",
			result.Error,
		)
	}
	if stakedCount == 0 {
		return 0, nil
	}

	// 4. Calculate distribution per staked duel bot.
	rewardPerBot := (totalFee * config.DUEL_BOT_TOTAL_SHARE / 100) / stakedCount
	if rewardPerBot == 0 {
		return 0, nil
	}

	// 5. Distribute rewards.
	if result := session.Model(
		&models.DuelBot{},
	).Where(
		"status != ?",
		models.DuelBotNormal,
	).Updates(
		map[string]interface{}{
			"staking_reward": gorm.Expr("staking_reward + ?", rewardPerBot),
			"total_earned":   gorm.Expr("total_earned + ?", rewardPerBot),
		},
	); result.Error != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"distributeFeeToDuelBots",
			"failed to distribute fee",
			result.Error,
		)
	}

	return rewardPerBot * stakedCount, nil
}

func claimDuelBotsRewards(to User, duelBots []Nft, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if to <= 0 || len(duelBots) == 0 {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"invalid parameter",
			errors.New("0 duelbots, or invalid to user"),
		)
	}

	// 2. Get session.
	session, err := getSession(sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to get session",
			err,
		)
	}

	// 3. Retrieve Deposited Nft IDs
	duelBotNfts := []models.DepositedNft{}
	if result := session.Where(
		"mint_address in ?",
		duelBots,
	).Find(&duelBotNfts); result.Error != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to retrieve duel bot nfts",
			result.Error,
		)
	}

	if len(duelBotNfts) != len(duelBots) {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"provided duelbots mismatching with retrieved",
			fmt.Errorf(
				"provided: %v, retrieved: %v",
				len(duelBots),
				len(duelBotNfts),
			),
		)
	}

	duelBotNftIDs := []uint{}
	for _, nft := range duelBotNfts {
		duelBotNftIDs = append(duelBotNftIDs, nft.ID)
	}

	// 4. Retrieve duelbot records.
	duelBotRecords := []models.DuelBot{}
	if result := session.Clauses(
		clause.Locking{Strength: "UPDATE"},
	).Preload(
		"DepositedNft",
	).Where(
		"status != ?",
		models.DuelBotNormal,
	).Where(
		"staking_user_id = ?",
		to,
	).Where(
		"deposited_nft_id in ?",
		duelBotNftIDs,
	).Find(&duelBotRecords); result.Error != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to retrieve duel bot records",
			result.Error,
		)
	}

	if len(duelBotRecords) != len(duelBots) {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"retrieved duel bot records count mismatching",
			fmt.Errorf(
				"requested count: %d, actual retrieved: %d",
				len(duelBots),
				len(duelBotRecords),
			),
		)
	}

	// 5. Sum up rewards and reset rewards.
	totalRewards := int64(0)
	for i, duelBot := range duelBotRecords {
		totalRewards += duelBot.StakingReward
		duelBotRecords[i].StakingReward = 0
	}

	if result := session.Save(&duelBotRecords); result.Error != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to update status and staking user",
			result.Error,
		)
	}
	if totalRewards == 0 {
		return 0, nil
	}

	// 6. Transfer rewards to user.
	/*transferResult*/
	transferResult, err := transfer(
		nil,
		&to,
		&BalanceLoad{
			ChipBalance: &totalRewards,
		},
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to transfer rewards to user",
			err,
		)
	}

	// ================== Archived for transaction ==================
	// 7. Record transaction as confirmed one.
	toWallet, err := GetUserWallet(&to, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to retrieve user's wallet",
			err,
		)
	}

	transaction, err := RecordTransaction(
		&TransactionLoad{
			FromWallet: nil,
			ToWallet:   toWallet,
			Balance: BalanceLoad{
				ChipBalance: &totalRewards,
			},
			Type: models.TxClaimStakingReward,
		},
		sessionId...,
	)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to record transaction",
			err,
		)
	}

	if err := ConfirmTransaction(
		&TransactionLoad{
			ToWalletPrevID: transferResult.ToPrevBalance,
			ToWalletNextID: transferResult.ToNextBalance,
			OwnerID:        uint(*toWallet),
			OwnerType:      models.TransactionWalletReferenced,
		},
		transaction,
		sessionId...,
	); err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"claimDuelBotsRewards",
			"failed to confirm transaction",
			err,
		)
	}
	// ================== Archived for transaction ==================

	return totalRewards, nil
}

func unstakeDuelBots(to User, duelBots []Nft, sessionId ...UUID) (int64, error) {
	// 1. Validate parameters.
	if to <= 0 || len(duelBots) == 0 {
		return 0, utils.MakeError(
			"duel_bot_db",
			"unstakeDuelBots",
			"invalid parameter",
			errors.New("0 duelbots, or 0 to user"),
		)
	}

	// 2. Claim rewards.
	rewards, err := claimDuelBotsRewards(to, duelBots, sessionId...)
	if err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"unstakeDueBots",
			"failed to claim rewards",
			err,
		)
	}

	// 3. Transfer duelbots from stake user to user.
	duelStake := User(config.DUEL_BOT_STAKE_ID)
	if _, err := transfer(
		&duelStake,
		&to,
		&BalanceLoad{
			NftBalance: &duelBots,
		},
		sessionId...,
	); err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"unstakeDuelBots",
			"failed to transfer duel bot",
			err,
		)
	}

	// 4. Update unstaked duel bot records.
	if err := recordUnstakedDuelBotsChecked(to, duelBots, sessionId...); err != nil {
		return 0, utils.MakeError(
			"duel_bot_db",
			"unstakeDuelBots",
			"failed to update staked user",
			err,
		)
	}

	return rewards, nil
}
