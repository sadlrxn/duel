package prelude

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestInitDuelBots(t *testing.T) {
	db := tests.InitMockDB(true, true)

	mockCollections := []models.NftCollection{
		{
			Name: "Mock Collection",
		},
		{
			Name: "Duelbots",
		},
	}
	if result := db.Create(&mockCollections); result.Error != nil {
		t.Fatalf("failed to create mock collections: %v", result.Error)
	}

	mockNfts := []models.DepositedNft{
		{
			Name:         "Mock Nft #1",
			CollectionID: 1,
		},
		{
			Name:         "Mock Nft #2",
			CollectionID: 1,
		},
		{
			Name:         "Mock Nft #3",
			CollectionID: 1,
		},
		{
			Name:         "Duel Bot #1",
			CollectionID: 2,
		},
		{
			Name:         "Duel Bot #2",
			CollectionID: 2,
		},
	}
	if result := db.Create(&mockNfts); result.Error != nil {
		t.Fatalf("failed to create mock nfts: %v", result.Error)
	}

	if err := InitDuelBotsV2(db); err != nil {
		t.Fatalf("failed to init duel bots: %v", err)
	}

	duelBotRecords := []models.DuelBot{}
	if result := db.Where("id is not null").Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duelbot records: %v", result.Error)
	}

	if len(duelBotRecords) != 2 ||
		duelBotRecords[0].DepositedNftID == 1 ||
		duelBotRecords[1].DepositedNftID == 2 {
		t.Fatalf(
			"failed to create duelbot records properly, \n\rlen: %d, 0_ID: %d, 1_ID: %d",
			len(duelBotRecords),
			duelBotRecords[0].DepositedNftID,
			duelBotRecords[1].DepositedNftID,
		)
	}

	dbRecordCnt := int64(0)
	if result := db.Model(&models.DuelBot{}).Where(
		"id is not null",
	).Count(
		&dbRecordCnt,
	); result.Error != nil {
		t.Fatalf(
			"failed to retrieve total count of duel bot records: %v",
			result.Error,
		)
	}

	if dbRecordCnt != 2 {
		t.Fatalf("failed to retrieve count of duelbot properly")
	}

	newDbNfts := []models.DepositedNft{
		{
			CollectionID: 2,
			Name:         "Duel Bot #3",
		},
		{
			CollectionID: 2,
			Name:         "Duel Bot #4",
		},
	}
	if result := db.Create(&newDbNfts); result.Error != nil {
		t.Fatalf("failed to create new duel bot nfts: %v", result.Error)
	}
	if err := InitDuelBotsV2(db); err != nil {
		t.Fatalf("failed to init duel bot records again: %v", err)
	}

	duelBotRecords = []models.DuelBot{}
	if result := db.Where("id is not null").Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duelbot records: %v", result.Error)
	}

	if len(duelBotRecords) != 4 ||
		duelBotRecords[0].DepositedNftID == 1 ||
		duelBotRecords[1].DepositedNftID == 2 ||
		duelBotRecords[2].DepositedNftID == 3 ||
		duelBotRecords[3].DepositedNftID == 4 {
		t.Fatalf("failed to create duelbot records properly")
	}

}
