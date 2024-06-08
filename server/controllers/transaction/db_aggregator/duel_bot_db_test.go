package db_aggregator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/lib/pq"
)

func initMockDuelBots() error {
	db, err := getSession()
	if err != nil {
		return err
	}

	user, duelStake := getMockUser()
	user.Wallet.Balance.NftBalance.Balance = pq.StringArray{
		"Duel Mint Address #1",
		"Duel Mint Address #2",
		"Duel Mint Address #3",
		"Duel Mint Address #4",
		"Duel Mint Address #5",
	}
	duelStake.ID = config.DUEL_BOT_STAKE_ID
	if result := db.Create(&[]models.User{user, duelStake}); result.Error != nil {
		fmt.Println("failed to create mock user")
		return result.Error
	}

	collection := models.NftCollection{
		Name: "Duel Bot",
	}
	if result := db.Create(&collection); result.Error != nil {
		fmt.Println("failed to create mock duel bot collection")
		return result.Error
	}

	wallet := uint(1)
	nfts := []models.DepositedNft{
		{
			Name:         "Duel Bot #1",
			MintAddress:  "Duel Mint Address #1",
			WalletID:     &wallet,
			CollectionID: 1,
		},
		{
			Name:         "Duel Bot #2",
			MintAddress:  "Duel Mint Address #2",
			WalletID:     &wallet,
			CollectionID: 1,
		},
		{
			Name:         "Duel Bot #3",
			MintAddress:  "Duel Mint Address #3",
			WalletID:     &wallet,
			CollectionID: 1,
		},
		{
			Name:         "Duel Bot #4",
			MintAddress:  "Duel Mint Address #4",
			WalletID:     &wallet,
			CollectionID: 1,
		},
		{
			Name:         "Duel Bot #5",
			MintAddress:  "Duel Mint Address #5",
			WalletID:     &wallet,
			CollectionID: 1,
		},
	}
	if result := db.Create(&nfts); result.Error != nil {
		fmt.Println("failed to create mock duel bot nfts")
		return result.Error
	}

	records := []models.DuelBot{
		{
			DepositedNftID: 1,
			Status:         models.DuelBotNormal,
		},
		{
			DepositedNftID: 2,
			Status:         models.DuelBotNormal,
		},
		{
			DepositedNftID: 3,
			Status:         models.DuelBotNormal,
		},
		{
			DepositedNftID: 4,
			Status:         models.DuelBotNormal,
		},
		{
			DepositedNftID: 5,
			Status:         models.DuelBotNormal,
		},
	}
	if result := db.Create(&records); result.Error != nil {
		fmt.Println("failed to create mock duel bot records")
		return result.Error
	}

	return nil
}

func TestInitMockDuelBots(t *testing.T) {
	db := tests.InitMockDB(true, true)
	initialize(db)

	if err := initMockDuelBots(); err != nil {
		t.Fatalf("failed to init mock duel bots: %v", err)
	}

	userID := User(1)
	balanceID, err := getUserBalance(&userID, false)
	if err != nil {
		t.Fatalf("failed to get user's balance id: %v", err)
	}
	balanceLoad, err := getBalance(balanceID)
	if err != nil {
		t.Fatalf("failed to get user's balance load: %v", err)
	}

	if balanceLoad.ChipBalance == nil ||
		*balanceLoad.ChipBalance != 100 {
		t.Fatalf("unexpected user's chip balance")
	}

	if balanceLoad.NftBalance == nil ||
		len(*balanceLoad.NftBalance) != 5 {
		t.Fatalf("unexpected user's nft balance")
	}

	stakeID := User(config.DUEL_BOT_STAKE_ID)
	balanceID, err = getUserBalance(&stakeID, false)
	if err != nil {
		t.Fatalf("failed to get duel stake's balance id: %v", err)
	}
	balanceLoad, err = getBalance(balanceID)
	if err != nil {
		t.Fatalf("failed to get duel stake's balance load: %v", err)
	}

	if balanceLoad.ChipBalance == nil ||
		*balanceLoad.ChipBalance != 50 {
		t.Fatalf("unexpected duel stake's chip balance")
	}

	if balanceLoad.NftBalance == nil ||
		len(*balanceLoad.NftBalance) != 0 {
		t.Fatalf("unexpected duel stake's nft balance")
	}

}

func TestRecordStakedDuelBotsChecked(t *testing.T) {
	db := tests.InitMockDB(true, true)
	initialize(db)

	if err := initMockDuelBots(); err != nil {
		t.Fatalf("failed to init mock duel bots: %v", err)
	}

	fromUser := User(1)
	if err := recordStakedDuelBotsChecked(
		fromUser,
		[]Nft{
			Nft("Duel Mint Address #1"),
			Nft("Duel Mint Address #2"),
			Nft("Duel Mint Address #3"),
		},
	); err == nil || err.Error() != utils.MakeError(
		"duel_bot_db",
		"recordStakedDuelBotsChecked",
		"not staked duelbot",
		errors.New("duel bot not owned by stake user"),
	).Error() {
		t.Fatalf("should require owned by stake user to be recorded: %v", err.Error())
	}

	stakeUser := User(config.DUEL_BOT_STAKE_ID)
	_, err := transfer(
		&fromUser,
		&stakeUser,
		&BalanceLoad{
			NftBalance: &[]Nft{
				Nft("Duel Mint Address #1"),
				Nft("Duel Mint Address #2"),
				Nft("Duel Mint Address #3"),
			},
		},
	)
	if err != nil {
		t.Fatalf("failed to transfer nfts to duel stake user: %v", err)
	}

	if err := recordStakedDuelBotsChecked(
		fromUser,
		[]Nft{
			Nft("Duel Mint Address #1"),
			Nft("Duel Mint Address #2"),
			Nft("Duel Mint Address #3"),
		},
	); err != nil {
		t.Fatalf("failed to record staked duelbots: %v", err.Error())
	}

	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"staking_user_id is not null",
	).Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 3 ||
		duelBotRecords[0].ID != 1 ||
		duelBotRecords[1].ID != 2 ||
		duelBotRecords[2].ID != 3 {
		t.Fatalf("failed to record properly")
	}
}

func TestStake(t *testing.T) {
	db := tests.InitMockDB(true, true)
	initialize(db)

	if err := initMockDuelBots(); err != nil {
		t.Fatalf("failed to init mock duel bots: %v", err)
	}

	fromUser := User(1)
	if err := stakeDuelBots(fromUser, []Nft{
		Nft("Duel Mint Address #1"),
		Nft("Duel Mint Address #2"),
		Nft("Duel Mint Address #3"),
	}); err != nil {
		t.Fatalf("failed to stake duel bots: %v", err)
	}

	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"staking_user_id is not null",
	).Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 3 ||
		duelBotRecords[0].ID != 1 ||
		duelBotRecords[1].ID != 2 ||
		duelBotRecords[2].ID != 3 {
		t.Fatalf("failed to record properly")
	}

	duelBotNfts := []models.DepositedNft{}
	if result := db.Where(
		"wallet_id = ?", 2,
	).Find(&duelBotNfts); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot nfts: %v", result.Error)
	}

	if len(duelBotNfts) != 3 ||
		duelBotNfts[0].ID != 1 ||
		duelBotNfts[1].ID != 2 ||
		duelBotNfts[2].ID != 3 {
		t.Fatalf("staked duel bots should be owned by stake user")
	}
}

func TestDistributeFee(t *testing.T) {
	TestStake(t)

	totalFee := int64(100)
	distributed, err := distributeFeeToDuelBots(totalFee)
	if err != nil {
		t.Fatalf("failed to distribute fee: %v", err)
	}
	if distributed != 78 {
		t.Fatalf("failed to properly distribute fee, actual: %d", distributed)
	}

	db := tests.InitMockDB(false, false)

	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Order(
		"id",
	).Find(
		&duelBotRecords,
	); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 5 ||
		duelBotRecords[0].StakingReward != 26 ||
		duelBotRecords[1].StakingReward != 26 ||
		duelBotRecords[2].StakingReward != 26 ||
		duelBotRecords[3].StakingReward != 0 ||
		duelBotRecords[4].StakingReward != 0 ||
		duelBotRecords[0].TotalEarned != 26 ||
		duelBotRecords[1].TotalEarned != 26 ||
		duelBotRecords[2].TotalEarned != 26 ||
		duelBotRecords[3].TotalEarned != 0 ||
		duelBotRecords[4].TotalEarned != 0 {
		t.Fatalf("failed to distribute fee properly")
	}
}

func TestClaimDuelBotsRewards(t *testing.T) {
	TestDistributeFee(t)

	if _, err := claimDuelBotsRewards(
		User(1),
		[]Nft{
			Nft("Duel Mint Address #3"),
			Nft("Duel Mint Address #4"),
		},
	); err == nil || err.Error() != utils.MakeError(
		"duel_bot_db",
		"claimDuelBotsRewards",
		"retrieved duel bot records count mismatching",
		fmt.Errorf(
			"requested count: 2, actual retrieved: 1",
		),
	).Error() {
		t.Fatalf("failed to check staked duel bots properly")
	}

	rewards, err := claimDuelBotsRewards(
		User(1),
		[]Nft{
			Nft("Duel Mint Address #2"),
			Nft("Duel Mint Address #3"),
		},
	)
	if err != nil {
		t.Fatalf("failed to claim rewards: %v", err)
	}
	if rewards != 52 {
		t.Fatalf("failed to claim exact amount: %d", rewards)
	}

	db := tests.InitMockDB(false, false)

	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Order(
		"id",
	).Find(
		&duelBotRecords,
	); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 5 ||
		duelBotRecords[0].StakingReward != 26 ||
		duelBotRecords[1].StakingReward != 0 ||
		duelBotRecords[2].StakingReward != 0 ||
		duelBotRecords[3].StakingReward != 0 ||
		duelBotRecords[4].StakingReward != 0 ||
		duelBotRecords[0].TotalEarned != 26 ||
		duelBotRecords[1].TotalEarned != 26 ||
		duelBotRecords[2].TotalEarned != 26 ||
		duelBotRecords[3].TotalEarned != 0 ||
		duelBotRecords[4].TotalEarned != 0 {
		t.Fatalf("failed to distribute fee properly")
	}

	user := models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where("id = ?", 1).Find(&user); result.Error != nil {
		t.Fatalf("failed to retrieve user: %v", result.Error)
	}

	if user.Wallet.Balance.ChipBalance.Balance != 152 {
		t.Fatalf(
			"failed to properly claim rewards, expected :%v, actual: %v",
			152,
			user.Wallet.Balance.ChipBalance.Balance,
		)
	}
}

func TestRecordUnstakedDuelBotsChecked(t *testing.T) {
	TestStake(t)

	stakeUser := User(config.DUEL_BOT_STAKE_ID)
	toUser := User(1)

	if err := recordUnstakedDuelBotsChecked(
		toUser,
		[]Nft{
			Nft("Duel Mint Address #2"),
		},
	); err == nil || err.Error() != utils.MakeError(
		"duel_bot_db",
		"recordUnstakedDuelBotsChecked",
		"not staked duelbot",
		errors.New("duel bot not owned by user"),
	).Error() {
		t.Fatalf("failed to check duel bot status properly: %v", err)
	}

	if _, err := transfer(
		&stakeUser,
		&toUser,
		&BalanceLoad{
			NftBalance: &[]Nft{
				Nft("Duel Mint Address #2"),
			},
		},
	); err != nil {
		t.Fatalf("failed to transfer nfts back to user: %v", err)
	}

	if err := recordUnstakedDuelBotsChecked(
		toUser,
		[]Nft{
			Nft("Duel Mint Address #2"),
			Nft("Duel Mint Address #4"),
		},
	); err == nil || err.Error() != utils.MakeError(
		"duel_bot_db",
		"recordUnstakedDuelBotsChecked",
		"retrieved duel bot records count mismatching",
		fmt.Errorf("requested count: 2, actual retrieved: 1"),
	).Error() {
		t.Fatalf("failed to check duel bots ownership properly: %v", err)
	}

	if err := recordUnstakedDuelBotsChecked(
		toUser,
		[]Nft{
			Nft("Duel Mint Address #2"),
		},
	); err != nil {
		t.Fatalf("failed to record unstaked duel bots: %v", err)
	}

	db := tests.InitMockDB(false, false)
	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Order("id").Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duelbot records: %v", result.Error)
	}

	if len(duelBotRecords) != 5 ||
		duelBotRecords[0].Status != models.DuelBotStaked ||
		duelBotRecords[1].Status != models.DuelBotNormal ||
		duelBotRecords[2].Status != models.DuelBotStaked ||
		duelBotRecords[3].Status != models.DuelBotNormal ||
		duelBotRecords[4].Status != models.DuelBotNormal {
		t.Fatalf("failed to properly record unstaked duel bots")
	}
}

func TestUnstake(t *testing.T) {
	TestDistributeFee(t)

	rewards, err := unstakeDuelBots(
		User(1),
		[]Nft{
			Nft("Duel Mint Address #2"),
			Nft("Duel Mint Address #3"),
		},
	)
	if err != nil {
		t.Fatalf("failed to unstake duel bots: %v", err)
	}
	if rewards != 52 {
		t.Fatalf("failed to unstake with exact rewards: %d", rewards)
	}

	db := tests.InitMockDB(false, false)

	duelBotRecords := []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Order(
		"id",
	).Find(
		&duelBotRecords,
	); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 5 ||
		duelBotRecords[0].StakingReward != 26 ||
		duelBotRecords[1].StakingReward != 0 ||
		duelBotRecords[2].StakingReward != 0 ||
		duelBotRecords[3].StakingReward != 0 ||
		duelBotRecords[4].StakingReward != 0 ||
		duelBotRecords[0].TotalEarned != 26 ||
		duelBotRecords[1].TotalEarned != 26 ||
		duelBotRecords[2].TotalEarned != 26 ||
		duelBotRecords[3].TotalEarned != 0 ||
		duelBotRecords[4].TotalEarned != 0 {
		t.Fatalf("failed to distribute fee properly")
	}

	user := models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Preload(
		"Wallet.Balance.NftBalance",
	).Where("id = ?", 1).Find(&user); result.Error != nil {
		t.Fatalf("failed to retrieve user: %v", result.Error)
	}

	if user.Wallet.Balance.ChipBalance.Balance != 152 {
		t.Fatalf(
			"failed to properly claim rewards, expected :%v, actual: %v",
			152,
			user.Wallet.Balance.ChipBalance.Balance,
		)
	}
	if len(user.Wallet.Balance.NftBalance.Balance) != 4 ||
		user.Wallet.Balance.NftBalance.Balance[0] != "Duel Mint Address #4" ||
		user.Wallet.Balance.NftBalance.Balance[1] != "Duel Mint Address #5" ||
		user.Wallet.Balance.NftBalance.Balance[2] != "Duel Mint Address #2" ||
		user.Wallet.Balance.NftBalance.Balance[3] != "Duel Mint Address #3" {
		t.Fatalf("failed to transfer nfts back to user properly")
	}

	totalFee := int64(100)
	distributed, err := distributeFeeToDuelBots(totalFee)
	if err != nil {
		t.Fatalf("failed to distribute fee: %v", err)
	}
	if distributed != 80 {
		t.Fatalf("failed to properly distribute fee, actual: %d", distributed)
	}

	duelBotRecords = []models.DuelBot{}
	if result := db.Where(
		"id is not null",
	).Order(
		"id",
	).Find(
		&duelBotRecords,
	); result.Error != nil {
		t.Fatalf("failed to retrieve duel bot records: %v", result.Error)
	}

	if len(duelBotRecords) != 5 ||
		duelBotRecords[0].StakingReward != 106 ||
		duelBotRecords[1].StakingReward != 0 ||
		duelBotRecords[2].StakingReward != 0 ||
		duelBotRecords[3].StakingReward != 0 ||
		duelBotRecords[4].StakingReward != 0 ||
		duelBotRecords[0].TotalEarned != 106 ||
		duelBotRecords[1].TotalEarned != 26 ||
		duelBotRecords[2].TotalEarned != 26 ||
		duelBotRecords[3].TotalEarned != 0 ||
		duelBotRecords[4].TotalEarned != 0 {
		t.Fatalf("failed to distribute fee properly")
	}
}
