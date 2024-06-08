package duel_bots

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func getUserDuelBots(userID uint) ([]DuelBotMeta, error) {
	db := db.GetDB()

	// 1. Validate parameters.
	if userID == 0 {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"invalid parameter",
			errors.New("provided user ID is invalid"),
		)
	}

	// 2. Retrieve User info.
	userInfo := models.User{}
	if result := db.Preload(
		"Wallet",
	).First(&userInfo, userID); result.Error != nil {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"failed to retrieve user info",
			result.Error,
		)
	}

	// 3. Retrieve Duel Bot nfts owned by user wallet.
	duelBotNfts := []models.DepositedNft{}
	if result := db.Where(
		"wallet_id = ?",
		userInfo.Wallet.ID,
	).Where("collection_id = ?", 97).Order("id").Find(&duelBotNfts); result.Error != nil {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"failed to retrieve duel bot nfts",
			result.Error,
		)
	}
	// if len(duelBotNfts) == 0 {
	// 	return []DuelBotMeta{}, nil
	// }

	duelBotNftIDs := []uint{}
	for _, nft := range duelBotNfts {
		duelBotNftIDs = append(duelBotNftIDs, nft.ID)
	}

	// 4. Retrieve Duel Bot records.
	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"deposited_nft_id in ?",
		duelBotNftIDs,
	).Order("deposited_nft_id").Find(&duelBotRecords); result.Error != nil {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"failed to retrieve duel bot records",
			result.Error,
		)
	}

	// 5. Build Duel Bot Meta.
	if len(duelBotRecords) != len(duelBotNfts) {
		return nil, utils.MakeError(
			"duel_bot",
			"getuserDuelBots",
			"duelbot nfts and records are mismatching",
			fmt.Errorf(
				"nfts: %d, records: %d",
				len(duelBotNfts),
				len(duelBotRecords),
			),
		)
	}
	duelBotMetas := []DuelBotMeta{}
	for i, nft := range duelBotNfts {
		if nft.ID != duelBotRecords[i].DepositedNftID {
			return nil, utils.MakeError(
				"duel_bot",
				"getUserDuelBots",
				"retrieved duelbot nfts' id mismatching with record",
				fmt.Errorf(
					"nftId: %d, recordId: %d",
					nft.ID,
					duelBotRecords[i].DepositedNftID,
				),
			)
		}
		meta := DuelBotMeta{
			Name:          nft.Name,
			MintAddress:   nft.MintAddress,
			Image:         nft.Image,
			Status:        duelBotRecords[i].Status,
			TotalEarned:   duelBotRecords[i].TotalEarned,
			StakingReward: duelBotRecords[i].StakingReward,
		}

		duelBotMetas = append(duelBotMetas, meta)
	}

	// 6. Retrieve staking duel bots.
	duelBotRecords = []models.DuelBot{}
	if result := db.Where(
		"status != ?",
		models.DuelBotNormal,
	).Where(
		"staking_user_id = ?",
		userID,
	).Order(
		"deposited_nft_id",
	).Find(&duelBotRecords); result.Error != nil {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"failed to retrieve staked duel bot records",
			result.Error,
		)
	}
	if len(duelBotRecords) == 0 {
		return duelBotMetas, nil
	}

	// 7. Retrieve staking duelbot nfts.
	duelBotNftIDs = []uint{}
	for _, record := range duelBotRecords {
		duelBotNftIDs = append(duelBotNftIDs, record.DepositedNftID)
	}
	duelBotNfts = []models.DepositedNft{}
	if result := db.Where(
		"id in ?",
		duelBotNftIDs,
	).Order("id").Find(&duelBotNfts); result.Error != nil {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"failed to retrieve staked duel bot nfts",
			result.Error,
		)
	}
	if len(duelBotNfts) != len(duelBotRecords) {
		return nil, utils.MakeError(
			"duel_bot",
			"getUserDuelBots",
			"duelbot nfts and records are mismatching 2",
			fmt.Errorf(
				"nfts: %d, records: %d",
				len(duelBotNfts),
				len(duelBotRecords),
			),
		)
	}

	// 8. Append duel bot metas
	for i, nft := range duelBotNfts {
		if nft.ID != duelBotRecords[i].DepositedNftID {
			return nil, utils.MakeError(
				"duel_bot",
				"getUserDuelBots",
				"retrieved duelbot nfts' id mismatching with record 2",
				fmt.Errorf(
					"nftId: %d, recordId: %d",
					nft.ID,
					duelBotRecords[i].DepositedNftID,
				),
			)
		}
		meta := DuelBotMeta{
			Name:          nft.Name,
			MintAddress:   nft.MintAddress,
			Image:         nft.Image,
			Status:        models.DuelBotStaked,
			TotalEarned:   duelBotRecords[i].TotalEarned,
			StakingReward: duelBotRecords[i].StakingReward,
		}

		duelBotMetas = append(duelBotMetas, meta)
	}

	return duelBotMetas, nil
}
