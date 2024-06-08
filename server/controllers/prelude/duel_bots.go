package prelude

import (
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
)

func InitDuelBots(db *gorm.DB) error {
	// 0. Retrieve DuelBot Collection.
	dbName := "Duelbots"
	dbCollection := models.NftCollection{}
	if result := db.Where(
		"name = ?",
		dbName,
	).First(&dbCollection); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to retrieve duelbot collection",
			result.Error,
		)
	}

	// 1. Retrieve DuelBot nfts count.
	totalDbRecordCnt := int64(0)
	if result := db.Model(&models.DepositedNft{}).Where(
		"collection_id = ?",
		dbCollection.ID,
	).Count(
		&totalDbRecordCnt,
	); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to retrieve duelbot nfts count",
			result.Error,
		)
	}
	if totalDbRecordCnt == 0 {
		return nil
	}

	// 2. Retrieve DuelBot records count.
	dbRecordCnt := int64(0)
	if result := db.Model(&models.DuelBot{}).Where(
		"id is not null",
	).Count(
		&dbRecordCnt,
	); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to retrieve duelbot records count",
			result.Error,
		)
	}

	if dbRecordCnt >= int64(totalDbRecordCnt) {
		return nil
	}

	// 3. Retrieve DuelBot records.
	currentDbRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Find(&currentDbRecords); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to retrieve duelbot records",
			result.Error,
		)
	}
	currentDbNftIDs := []uint{0}
	for _, dbRecord := range currentDbRecords {
		currentDbNftIDs = append(currentDbNftIDs, dbRecord.DepositedNftID)
	}

	// 4. Retrieve DuelBot Nfts.
	dbNfts := []models.DepositedNft{}
	if result := db.Where(
		"collection_id = ?",
		dbCollection.ID,
	).Where(
		"id not in ?",
		currentDbNftIDs,
	).Find(&dbNfts); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to retrieve duelbot nfts",
			result.Error,
		)
	}
	if len(dbNfts) == 0 {
		return nil
	}

	// 5. Create DuelBot Records.
	dbRecords := []models.DuelBot{}
	for i := 0; i < len(dbNfts); i++ {
		dbRecords = append(dbRecords, models.DuelBot{
			DepositedNftID: dbNfts[i].ID,
		})
	}

	if result := db.Create(&dbRecords); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBots",
			"failed to create duelbot nfts",
			result.Error,
		)
	}

	return nil
}

func InitDuelBotsV2(db *gorm.DB) error {
	// 1. Retrieve DuelBot Collection.
	dbName := "Duelbots"
	dbCollection := models.NftCollection{}
	if result := db.Where(
		"name = ?",
		dbName,
	).First(&dbCollection); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBotsV2",
			"failed to retrieve duelbot collection",
			result.Error,
		)
	}

	// 2. Retrieve Duelbot records.
	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Find(&duelBotRecords); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBotsV2",
			"failed to retrieve duelbot records",
			result.Error,
		)
	}

	// 3. Get initialized deposited nft ids.
	duelBotNftIDs := []uint{0}
	for _, duelBotRecord := range duelBotRecords {
		duelBotNftIDs = append(duelBotNftIDs, duelBotRecord.DepositedNftID)
	}

	// 4. Get uninitialized deposited nft records.
	duelBotNfts := []models.DepositedNft{}
	if result := db.Where(
		"collection_id = ?",
		dbCollection.ID,
	).Where(
		"id not in ?",
		duelBotNftIDs,
	).Find(&duelBotNfts); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBotsV2",
			"failed to retrieve uninitialized duelbot nfts",
			result.Error,
		)
	}
	if len(duelBotNfts) == 0 {
		return nil
	}

	// 5. Create new duelbot records.
	newDuelBotRecords := []models.DuelBot{}
	for _, duelBotNft := range duelBotNfts {
		newDuelBotRecords = append(
			newDuelBotRecords,
			models.DuelBot{
				DepositedNftID: duelBotNft.ID,
			},
		)
	}
	if len(newDuelBotRecords) == 0 {
		return nil
	}
	if result := db.Create(&newDuelBotRecords); result.Error != nil {
		return utils.MakeError(
			"prelude_duel_bot",
			"InitDuelBotsV2",
			"failed to create new duelbot records",
			result.Error,
		)
	}

	return nil
}
