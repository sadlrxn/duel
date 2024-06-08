package transaction

import (
	"strings"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/lib/pq"
)

func TestTransfer(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to init mock db")
	}

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("failed to initialize db aggregator: %v", err)
	}

	users := []models.User{
		{
			Name:          "User",
			WalletAddress: "EvPpQ4TQHHFxsXjSaBWKZavvhXXwCLRv25LbMBfYmZGN",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 1000,
					},
				},
			},
		},
		{
			Name:          "Temp",
			WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 100,
					},
				},
			},
		},
		{
			Name:          "Fee",
			WalletAddress: "J8FgYgkGrwnapBAaPTkvZUfWGykXTgBzhmbUuo6ahLcD",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 50,
					},
				},
			},
		},
		{
			Name:          "User2",
			WalletAddress: "EvPpQ4TQHHFxsXjSaBWKZavvhXXwCLRv25LbMBfYmZGN",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 0,
					},
				},
			},
			Statistics: models.Statistics{
				TotalWagered: 250000000,
			},
		},
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	duelBotCollection := models.NftCollection{
		Name: "Duel Bots",
		Nfts: []models.DepositedNft{
			{
				Name:        "DuelBot #1",
				MintAddress: "DuelBot Mint Address #1",
			},
			{
				Name:        "DuelBot #2",
				MintAddress: "DuelBot Mint Address #2",
			},
			{
				Name:        "DuelBot #3",
				MintAddress: "DuelBot Mint Address #3",
			},
		},
	}
	if result := db.Create(&duelBotCollection); result.Error != nil {
		t.Fatalf("failed to create duelbot collections and nfts: %v", result.Error)
	}

	userID := uint(1)
	duelBotRecords := []models.DuelBot{
		{
			DepositedNftID: 1,
			Status:         models.DuelBotStaked,
			StakingUserID:  &userID,
		},
		{
			DepositedNftID: 2,
		},
		{
			DepositedNftID: 3,
			Status:         models.DuelBotStaked,
			StakingUserID:  &userID,
		},
	}
	if result := db.Create(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to create duelbot records: %v", result.Error)
	}

	if err := createAffiliateCode(db_aggregator.User(4), []string{"code1", "code2"}); err != nil {
		t.Fatalf("failed to create affiliate code properly: %v", err)
	}

	isFirst, err := db_aggregator.ActivateAffiliateCode(db_aggregator.User(1), "code2")
	if err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}
	if !isFirst {
		t.Fatalf("should be the first activation: %v", isFirst)
	}

	transferAmount := int64(100)
	tx, err := transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[1].ID),
		ToUser:   (*db_aggregator.User)(&users[2].ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &transferAmount,
		},
		Type:          models.TxCoinflipFee,
		ToBeConfirmed: true,
		HouseFeeMeta: &HouseFeeMeta{
			User: db_aggregator.User(1),
		},
	})
	if err != nil {
		t.Fatalf("failed to perform transfer transaction: %v", err)
	}
	if uint(*tx) != 1 {
		t.Fatalf("failed to create transfer record properly: %d", *tx)
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 1000 ||
		users[1].Wallet.Balance.ChipBalance.Balance != 0 ||
		users[2].Wallet.Balance.ChipBalance.Balance != 58 {
		t.Fatalf(
			"failed to distribute fee properly: %d, %d, %d",
			users[0].Wallet.Balance.ChipBalance.Balance,
			users[1].Wallet.Balance.ChipBalance.Balance,
			users[2].Wallet.Balance.ChipBalance.Balance,
		)
	}

	duelBotRecords = []models.DuelBot{}
	if result := db.Where("id is not null").Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duelbot records properly: %v", result.Error)
	}

	if duelBotRecords[0].TotalEarned != 40 ||
		duelBotRecords[1].TotalEarned != 0 ||
		duelBotRecords[2].TotalEarned != 40 ||
		duelBotRecords[0].StakingReward != 40 ||
		duelBotRecords[1].StakingReward != 0 ||
		duelBotRecords[2].StakingReward != 40 {
		t.Fatal("failed to distribute staking reward properly")
	}

	rakebacks := []models.Rakeback{}
	if result := db.Where("id is not null").Find(&rakebacks); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback records: %v", result.Error)
	}

	if len(rakebacks) != 1 ||
		rakebacks[0].TotalEarned != 10 ||
		rakebacks[0].UserID != 1 ||
		rakebacks[0].Reward != 10 {
		t.Fatal("failed to distribute rakeback reward properly")
	}

	rakeback, err := claimRakeback((*db_aggregator.User)(&userID))
	if err != nil {
		t.Fatalf("failed to claim rakeback: %v", err)
	}
	if rakeback != 10 {
		t.Fatalf("failed to claim rakeback properly: %d", rakeback)
	}

	codes, err := db_aggregator.GetOwnedAffiliateCode(db_aggregator.User(4))
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes: %v", err)
	}
	if len(codes) != 2 ||
		codes[0].Code != "code1" ||
		codes[0].Reward != 0 ||
		codes[0].TotalEarned != 0 ||
		codes[0].UserCnt != 0 ||
		codes[1].Code != "code2" ||
		codes[1].Reward != 2 ||
		codes[1].TotalEarned != 2 ||
		codes[1].UserCnt != 1 {
		t.Fatalf("failed to distribute affiliate rewards properly: %v", codes)
	}

	transaction := models.Transaction{}
	if result := db.Preload(
		"ToWalletPrev.ChipBalance",
	).Preload(
		"ToWalletNext.ChipBalance",
	).First(&transaction, 2); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback transaction: %v", result.Error)
	}

	if transaction.FromWallet != nil ||
		*transaction.ToWallet != 1 ||
		transaction.ToWalletPrev.ChipBalance.Balance != 1000 ||
		transaction.ToWalletNext.ChipBalance.Balance != 1010 ||
		transaction.Type != models.TxClaimRakebackReward ||
		transaction.Status != models.TransactionSucceed {
		t.Fatalf("failed to record rakeback transaction properly: %v", transaction)
	}

	rakebacks = []models.Rakeback{}
	if result := db.Where("id is not null").Find(&rakebacks); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback records after claiming: %v", result.Error)
	}

	if len(rakebacks) != 1 ||
		rakebacks[0].TotalEarned != 10 ||
		rakebacks[0].UserID != 1 ||
		rakebacks[0].Reward != 0 {
		t.Fatal("failed to update amount rakeback reward properly after claiming")
	}

	stakingReward, err := claimDuelBotsRewards(DuelBotsRequest{
		FromUser: db_aggregator.User(userID),
		DuelBots: []db_aggregator.Nft{
			"DuelBot Mint Address #1",
			"DuelBot Mint Address #3",
		},
	})
	if err != nil {
		t.Fatalf("failed to claim duelbot rewards: %v", err)
	}
	if stakingReward != 80 {
		t.Fatalf("failed to claim duelbot rewards properly: %d", stakingReward)
	}

	transaction = models.Transaction{}
	if result := db.Preload(
		"ToWalletPrev.ChipBalance",
	).Preload(
		"ToWalletNext.ChipBalance",
	).First(&transaction, 3); result.Error != nil {
		t.Fatalf("failed to retrieve claim staking transaction: %v", result.Error)
	}

	if transaction.FromWallet != nil ||
		*transaction.ToWallet != 1 ||
		transaction.ToWalletPrev.ChipBalance.Balance != 1010 ||
		transaction.ToWalletNext.ChipBalance.Balance != 1090 ||
		transaction.Type != models.TxClaimStakingReward ||
		transaction.Status != models.TransactionSucceed {
		t.Fatalf("failed to record claim staking transaction properly: %v", transaction)
	}

	duelBotRecords = []models.DuelBot{}
	if result := db.Where("id is not null").Find(&duelBotRecords); result.Error != nil {
		t.Fatalf("failed to retrieve duelbot records properly after claiming: %v", result.Error)
	}

	if duelBotRecords[0].TotalEarned != 40 ||
		duelBotRecords[1].TotalEarned != 0 ||
		duelBotRecords[2].TotalEarned != 40 ||
		duelBotRecords[0].StakingReward != 0 ||
		duelBotRecords[1].StakingReward != 0 ||
		duelBotRecords[2].StakingReward != 0 {
		t.Fatal("failed to update staking reward properly after claiming")
	}

	claimed, err := claimAffiliateRewards(db_aggregator.User(4), []string{"code1", "code2"})
	if err != nil {
		t.Fatalf("failed to claim affiliate rewards: %v", err)
	}
	if claimed != 2 {
		t.Fatalf("failed to claim affiliate rewards properly: %d", claimed)
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info after all actions: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 1090 ||
		users[1].Wallet.Balance.ChipBalance.Balance != 0 ||
		users[2].Wallet.Balance.ChipBalance.Balance != 58 ||
		users[3].Wallet.Balance.ChipBalance.Balance != 2 {
		t.Fatalf(
			"failed to update balances after all actions: %d, %d, %d, %d",
			users[0].Wallet.Balance.ChipBalance.Balance,
			users[1].Wallet.Balance.ChipBalance.Balance,
			users[2].Wallet.Balance.ChipBalance.Balance,
			users[3].Wallet.Balance.ChipBalance.Balance,
		)
	}

	transferAmount = 200
	transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[0].ID),
		ToUser:   (*db_aggregator.User)(&users[1].ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &transferAmount,
		},
		Type:          models.TxCoinflipBet,
		ToBeConfirmed: true,
	})

	rakebackInfo := models.Rakeback{}
	if result := db.Where(
		"user_id = ?", 1,
	).First(&rakebackInfo); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback info: %v", result.Error)
	}

	rakebackInfo.AdditionalRakebackExpired = time.Now()
	if result := db.Save(&rakebackInfo); result.Error != nil {
		t.Fatalf(
			"failed to update additional rakeback expired to be expired: %v",
			result.Error,
		)
	}

	transferAmount = 200
	transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[1].ID),
		ToUser:   (*db_aggregator.User)(&users[2].ID),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &transferAmount,
		},
		Type:          models.TxCoinflipFee,
		ToBeConfirmed: true,
		HouseFeeMeta: &HouseFeeMeta{
			User: db_aggregator.User(1),
		},
	})

	user1ID := uint(1)
	totalEarned, reward, err := getRakebackRewards((*db_aggregator.User)(&user1ID))
	if err != nil {
		t.Fatalf("failed to get rakeback rewards: %v", err)
	}
	if totalEarned != 20 ||
		reward != 10 {
		t.Fatalf(
			"failed to distribute rakeback properly: %d, %d",
			totalEarned,
			reward,
		)
	}

	codes, err = db_aggregator.GetOwnedAffiliateCode(db_aggregator.User(4))
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes: %v", err)
	}
	if len(codes) != 2 ||
		codes[0].Code != "code1" ||
		codes[0].Reward != 0 ||
		codes[0].TotalEarned != 0 ||
		codes[0].UserCnt != 0 ||
		codes[1].Code != "code2" ||
		codes[1].Reward != 5 ||
		codes[1].TotalEarned != 7 ||
		codes[1].UserCnt != 1 {
		t.Fatalf("failed to distribute affiliate rewards properly: %v", codes)
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info after all actions: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 890 ||
		users[1].Wallet.Balance.ChipBalance.Balance != 0 ||
		users[2].Wallet.Balance.ChipBalance.Balance != 83 ||
		users[3].Wallet.Balance.ChipBalance.Balance != 2 {
		t.Fatalf(
			"failed to update balances after all actions: %d, %d, %d, %d",
			users[0].Wallet.Balance.ChipBalance.Balance,
			users[1].Wallet.Balance.ChipBalance.Balance,
			users[2].Wallet.Balance.ChipBalance.Balance,
			users[3].Wallet.Balance.ChipBalance.Balance,
		)
	}
}

func TestWithdraw(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to init mock db")
	}

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("failed to initialize db aggregator: %v", err)
	}

	users := []models.User{
		{
			Name:          "User",
			WalletAddress: "User wallet address",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 1000000,
					},
					NftBalance: &models.NftBalance{
						Balance: pq.StringArray{
							"Mintaddress #1",
							"Mintaddress #2",
							"Mintaddress #3",
							"Mintaddress #4",
							"Mintaddress #5",
						},
					},
				},
			},
		},
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	nftCollection := models.NftCollection{
		Name: "Collection name",
		Nfts: []models.DepositedNft{
			{
				Name:        "NFT #1",
				MintAddress: "Mintaddress #1",
				WalletID:    &users[0].Wallet.ID,
			},
			{
				Name:        "NFT #2",
				MintAddress: "Mintaddress #2",
				WalletID:    &users[0].Wallet.ID,
			},
			{
				Name:        "NFT #3",
				MintAddress: "Mintaddress #3",
				WalletID:    &users[0].Wallet.ID,
			},
			{
				Name:        "NFT #4",
				MintAddress: "Mintaddress #4",
				WalletID:    &users[0].Wallet.ID,
			},
			{
				Name:        "NFT #5",
				MintAddress: "Mintaddress #5",
				WalletID:    &users[0].Wallet.ID,
			},
		},
	}
	if result := db.Create(&nftCollection); result.Error != nil {
		t.Fatalf("failed to create nft collection and nfts: %v", result.Error)
	}

	withdrawAmount := int64(1000000)
	_, err := transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[0].ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &withdrawAmount,
		},
		Type:          models.TxWithdrawSpl,
		ToBeConfirmed: false,
	})
	if err == nil ||
		!strings.Contains(err.Error(), "insufficient funds") {
		t.Fatalf("should be failed with inssuficiant balance error.")
	}

	withdrawAmount = int64(100000)
	_, err = transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[0].ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &withdrawAmount,
		},
		Type:          models.TxWithdrawSol,
		ToBeConfirmed: false,
	})
	if err != nil {
		t.Fatalf("failed to withdraw balance")
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 900000 {
		t.Fatalf(
			"Failed to burn balance properly: %d", users[0].Wallet.Balance.ChipBalance.Balance,
		)
	}

	withdrawAmount = int64(100000)
	_, err = transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[0].ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &withdrawAmount,
		},
		Type:          models.TxWithdrawSpl,
		ToBeConfirmed: false,
	})
	if err != nil {
		t.Fatalf("failed to withdraw balance")
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 790000 {
		t.Fatalf(
			"Failed to burn balance+fee properly: %d", users[0].Wallet.Balance.ChipBalance.Balance,
		)
	}

	nftsToWithdraw := []string{
		"Mintaddress #1",
		"Mintaddress #2",
		"Mintaddress #3",
		"Mintaddress #4",
		"Mintaddress #5",
	}
	_, err = transfer(&TransactionRequest{
		FromUser: (*db_aggregator.User)(&users[0].ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			NftBalance: db_aggregator.ConvertStringArrayToNftArray(&nftsToWithdraw),
		},
		Type:          models.TxWithdrawNft,
		ToBeConfirmed: false,
	})
	if err != nil {
		t.Fatalf("failed to withdraw nft %v", err)
	}

	users = []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Preload("Wallet.Balance.NftBalance").Where(
		"id is not null",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if users[0].Wallet.Balance.ChipBalance.Balance != 790000-config.WITHDRAW_FEE_PER_SPL*int64(len(nftsToWithdraw)) {
		t.Fatalf(
			"Failed to burn balance+fee properly: %d", users[0].Wallet.Balance.ChipBalance.Balance,
		)
	}
	if len(users[0].Wallet.Balance.NftBalance.Balance) != 0 {
		t.Fatalf(
			"Failed to withdraw nfts properly: %v", users[0].Wallet.Balance.NftBalance.Balance,
		)
	}
}
