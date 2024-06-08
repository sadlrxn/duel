package db_aggregator

import (
	"errors"
	"testing"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func getMockUser() (models.User, models.User) {
	return models.User{
			Name:          "Mock User",
			WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 100,
					},
					NftBalance: &models.NftBalance{
						Balance: pq.StringArray{},
					},
				},
				History:   []models.Balance{},
				IncomeTx:  []models.Transaction{},
				OutcomeTx: []models.Transaction{},
			},
		}, models.User{
			Name:          "Mock User2",
			WalletAddress: "J8FgYgkGrwnapBAaPTkvZUfWGykXTgBzhmbUuo6ahLcD",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 50,
					},
					NftBalance: &models.NftBalance{
						Balance: pq.StringArray{},
					},
				},
				History:   []models.Balance{},
				IncomeTx:  []models.Transaction{},
				OutcomeTx: []models.Transaction{},
			},
		}
}

func TestShouldCommitUpdateDatabase(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session module: %v", err)
	}

	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start session, err: %v", err)
	}

	session, err := getSession(sessionId)
	if err != nil {
		t.Fatalf("failed to get session, err: %v", err)
	}

	mockUser := models.User{
		Name:          "Mock User",
		WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
		Role:          models.UserRole,
		Wallet: models.Wallet{
			Balance: models.Balance{
				ChipBalance: &models.ChipBalance{
					Balance: 100,
				},
				NftBalance: &models.NftBalance{
					Balance: pq.StringArray{},
				},
			},
			History:   []models.Balance{},
			IncomeTx:  []models.Transaction{},
			OutcomeTx: []models.Transaction{},
		},
	}

	if result := session.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create a mock user: %v", result.Error)
	}

	if err := commitSession(sessionId); err != nil {
		t.Fatalf("failed to commit session: %v", err)
	}

	if err := removeSession(sessionId); err != nil {
		t.Fatalf("failed to remove session: %v", err)
	}

	user := models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Preload(
		"Wallet.Balance.NftBalance",
	).First(&user); result.Error != nil {
		t.Fatalf("failed to get user from mock db: %v", result.Error)
	}

	if user.Name != mockUser.Name ||
		user.WalletAddress != mockUser.WalletAddress ||
		user.Role != mockUser.Role ||
		user.Wallet.Balance.ChipBalance.Balance != mockUser.Wallet.Balance.ChipBalance.Balance ||
		len(user.Wallet.Balance.NftBalance.Balance) != len(mockUser.Wallet.Balance.NftBalance.Balance) {
		t.Fatalf("data mismatching with saved and got user info")
	}
}

func TestShouldRemoveNotUpdateDatabase(t *testing.T) {
	TestShouldCommitUpdateDatabase(t)

	db, err := getSession()
	if err != nil || db == nil {
		t.Fatalf("failed to get main db: %v", err)
	}

	sessionId1, err := startSession()
	if err != nil {
		t.Fatalf("failed to start first session: %v", err)
	}

	session1, err := getSession(sessionId1)
	if err != nil {
		t.Fatalf("failed to get first session: %v", err)
	}

	mockUser2 := models.User{
		Name:          "Mock User2",
		WalletAddress: "J8FgYgkGrwnapBAaPTkvZUfWGykXTgBzhmbUuo6ahLcD",
		Role:          models.UserRole,
		Wallet: models.Wallet{
			Balance: models.Balance{
				ChipBalance: &models.ChipBalance{
					Balance: 50,
				},
				NftBalance: &models.NftBalance{
					Balance: pq.StringArray{},
				},
			},
			History:   []models.Balance{},
			IncomeTx:  []models.Transaction{},
			OutcomeTx: []models.Transaction{},
		},
	}

	if result := session1.Create(&mockUser2); result.Error != nil {
		t.Fatalf("failed to create second mock user: %v", err)
	}

	user1 := models.User{}

	if result := session1.Preload(
		"Wallet.Balance.ChipBalance",
	).First(
		&user1, "wallet_address = ?", "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
	); result.Error != nil {
		t.Fatalf("failed to get first mock user: %v", result.Error)
	}

	user1.Wallet.Balance.ChipBalance.Balance += 100
	if result := session1.Save(&user1); result.Error != nil {
		t.Fatalf("failed to update first mock user: %v", result.Error)
	}

	user2 := models.User{}
	if result := db.First(
		&user2, "wallet_address = ?", "J8FgYgkGrwnapBAaPTkvZUfWGykXTgBzhmbUuo6ahLcD",
	); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("should not insert second mock user, but found")
	}

	user1 = models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Preload(
		"Wallet.Balance.NftBalance",
	).First(
		&user1, "wallet_address = ?", "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
	); result.Error != nil {
		t.Fatalf("failed to get updated first mock user: %v", result.Error)
	}

	mockUser1 := models.User{
		Name:          "Mock User",
		WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
		Role:          models.UserRole,
		Wallet: models.Wallet{
			Balance: models.Balance{
				ChipBalance: &models.ChipBalance{
					Balance: 100,
				},
				NftBalance: &models.NftBalance{
					Balance: pq.StringArray{},
				},
			},
			History:   []models.Balance{},
			IncomeTx:  []models.Transaction{},
			OutcomeTx: []models.Transaction{},
		},
	}

	if user1.Name != mockUser1.Name ||
		user1.WalletAddress != mockUser1.WalletAddress ||
		user1.Role != mockUser1.Role ||
		user1.Wallet.Balance.ChipBalance.Balance != mockUser1.Wallet.Balance.ChipBalance.Balance ||
		len(user1.Wallet.Balance.NftBalance.Balance) != len(mockUser1.Wallet.Balance.NftBalance.Balance) {
		t.Fatalf("data mismatching with saved and got user info")
	}

	if err := removeSession(sessionId1); err != nil {
		t.Fatalf("failed to remove session: %v", err)
	}
}

func TestRemoveSession(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session module: %v", err)
	}

	mockUser, _ := getMockUser()

	sessionId, session := GetMockTxSession(t)

	if result := session.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	user := models.User{}
	if result := session.First(&user); result.Error != nil {
		t.Fatalf("failed to get mock user: %v", result.Error)
	}

	if user.Name != mockUser.Name ||
		user.WalletAddress != mockUser.WalletAddress {
		t.Fatalf("failed to get proper mock user")
	}

	if err := removeSession(sessionId); err != nil {
		t.Fatalf("failed to remove session: %v", err)
	}

	if result := db.First(&models.User{}); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("should rollback session, but found record: %v", result.Error)
	}
}

func TestTransfer(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session module: %v", err)
	}

	mockUser1, mockUser2 := getMockUser()

	if result := db.Create(&[]models.User{mockUser1, mockUser2}); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	sessionId, _ := GetMockTxSession(t)

	userId1, userId2, amount := User(1), User(2), int64(20)
	if _, err := transfer(&userId1, &userId2, &BalanceLoad{
		ChipBalance: &amount,
	}, sessionId); err != nil {
		t.Fatalf("failed to transfer chip balance: %v", err)
	}

	if err := commitSession(sessionId); err != nil {
		t.Fatalf("failed to commit transfer tx: %v", err)
	}

	users := []models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Find(&users); result.Error != nil {
		t.Fatalf("failed to retrieve users: %v", result.Error)
	}

	if mockUser1.Wallet.Balance.ChipBalance.Balance-amount != users[0].Wallet.Balance.ChipBalance.Balance ||
		mockUser2.Wallet.Balance.ChipBalance.Balance+amount != users[1].Wallet.Balance.ChipBalance.Balance {
		t.Fatalf("failed to transfer amount properly")
	}
}

func TestNftTransfer(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session module: %v", err)
	}

	mockUser1, mockUser2 := getMockUser()

	mockUser1.Wallet.Balance.NftBalance = &models.NftBalance{
		Balance: pq.StringArray{
			"ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE",
			"6mFiFUCDLLiSxTikNtAHtpE2JCgbB7VSkcKv76a7zjTa",
			"EyKa1EXfZk4Fa9kDSSYjjDxDSjYYVbZ6NafH8to3UFQb",
			"7BbLsCh44gqNVyhjXF3umHKUed58Ukkv9nSSDknJZXmZ",
			"EbavJqfhPTcbuErZmGaDuwRBnHjUGFYBZwKyAHbrX8aE",
			"caauoXogcDYVAjGyaEyrKAak7kPyok8KCK59fGgqmnm",
			"CsqXem8HLXnuKaToJcRXXdMbXZgPa1fBdNQMmugwuiCk",
			"Eby6ikiWRevcgWxTDYTA9rc4T9XZQKBnqbiEz69Sh6Y6",
			"371ijDTrtUwgbP5y6nkzCEHjMP3thSsFohr4FJeugPim",
			"DNMKCVGzB6YpoeSn64LinqTFSn51cu8sUPDcEz795tiN",
			"3fPoG48ERaCWxkawQpWsJLmuoRJwzxdrxavrHDGatcrk",
			"6wopTGJvQhCTvG96A15E4Fo7GR2whecRhaSKCPQ2ty2e",
			"BkUf5x3FoZ9FkhatxcrcbZjN34oWV2A7CTTD7sKvsjMR",
			"9z6x4Gm9d4P3j6CF87pZDJtahz5X9aBopyKfTVQVbDQW",
			"WJQZM1vFwbS3tHD9F9awCGEac5BJCyjzDY8g2gyVQQD",
		},
	}

	wallet1 := uint(1)
	mockCollection := models.NftCollection{
		Nfts: []models.DepositedNft{
			{
				MintAddress: "ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "6mFiFUCDLLiSxTikNtAHtpE2JCgbB7VSkcKv76a7zjTa",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "EyKa1EXfZk4Fa9kDSSYjjDxDSjYYVbZ6NafH8to3UFQb",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "7BbLsCh44gqNVyhjXF3umHKUed58Ukkv9nSSDknJZXmZ",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "EbavJqfhPTcbuErZmGaDuwRBnHjUGFYBZwKyAHbrX8aE",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "caauoXogcDYVAjGyaEyrKAak7kPyok8KCK59fGgqmnm",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "CsqXem8HLXnuKaToJcRXXdMbXZgPa1fBdNQMmugwuiCk",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "Eby6ikiWRevcgWxTDYTA9rc4T9XZQKBnqbiEz69Sh6Y6",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "371ijDTrtUwgbP5y6nkzCEHjMP3thSsFohr4FJeugPim",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "DNMKCVGzB6YpoeSn64LinqTFSn51cu8sUPDcEz795tiN",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "3fPoG48ERaCWxkawQpWsJLmuoRJwzxdrxavrHDGatcrk",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "6wopTGJvQhCTvG96A15E4Fo7GR2whecRhaSKCPQ2ty2e",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "BkUf5x3FoZ9FkhatxcrcbZjN34oWV2A7CTTD7sKvsjMR",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "9z6x4Gm9d4P3j6CF87pZDJtahz5X9aBopyKfTVQVbDQW",
				WalletID:    &wallet1,
			},
			{
				MintAddress: "WJQZM1vFwbS3tHD9F9awCGEac5BJCyjzDY8g2gyVQQD",
				WalletID:    &wallet1,
			},
		},
	}

	if result := db.Create(&[]models.User{mockUser1, mockUser2}); result.Error != nil {
		t.Fatalf("failed to initialize mock users: %v", result.Error)
	}

	if result := db.Create(&mockCollection); result.Error != nil {
		t.Fatalf("failed to initialize mock collection and nfts: %v", result.Error)
	}

	fromUserId, toUserId := User(1), User(2)
	_, err := transfer(&fromUserId, &toUserId, &BalanceLoad{
		NftBalance: &[]Nft{"ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE"},
	})

	if err != nil {
		t.Fatalf("failed to transfer nft: %v", err)
	}

	user1, user2 := models.User{}, models.User{}
	if result := db.Preload("Wallet.Balance.NftBalance").First(&user1, 1); result.Error != nil {
		t.Fatalf("failed to retrieve first mock user: %v", result.Error)
	}
	if result := db.Preload("Wallet.Balance.NftBalance").First(&user2, 2); result.Error != nil {
		t.Fatalf("failed to retrieve second mock user: %v", result.Error)
	}

	user1Exp := pq.StringArray{
		"6mFiFUCDLLiSxTikNtAHtpE2JCgbB7VSkcKv76a7zjTa",
		"EyKa1EXfZk4Fa9kDSSYjjDxDSjYYVbZ6NafH8to3UFQb",
		"7BbLsCh44gqNVyhjXF3umHKUed58Ukkv9nSSDknJZXmZ",
		"EbavJqfhPTcbuErZmGaDuwRBnHjUGFYBZwKyAHbrX8aE",
		"caauoXogcDYVAjGyaEyrKAak7kPyok8KCK59fGgqmnm",
		"CsqXem8HLXnuKaToJcRXXdMbXZgPa1fBdNQMmugwuiCk",
		"Eby6ikiWRevcgWxTDYTA9rc4T9XZQKBnqbiEz69Sh6Y6",
		"371ijDTrtUwgbP5y6nkzCEHjMP3thSsFohr4FJeugPim",
		"DNMKCVGzB6YpoeSn64LinqTFSn51cu8sUPDcEz795tiN",
		"3fPoG48ERaCWxkawQpWsJLmuoRJwzxdrxavrHDGatcrk",
		"6wopTGJvQhCTvG96A15E4Fo7GR2whecRhaSKCPQ2ty2e",
		"BkUf5x3FoZ9FkhatxcrcbZjN34oWV2A7CTTD7sKvsjMR",
		"9z6x4Gm9d4P3j6CF87pZDJtahz5X9aBopyKfTVQVbDQW",
		"WJQZM1vFwbS3tHD9F9awCGEac5BJCyjzDY8g2gyVQQD",
	}

	allMatch := true
	result := &user1.Wallet.Balance.NftBalance.Balance
	for i, n := 0, len(*result); i < n; i++ {
		allMatch = allMatch && ((*result)[i] == user1Exp[i])
	}

	if !allMatch {
		t.Fatalf("failed to user1 result balance: %v %d", result, len(*result))
	}

	user2Exp := pq.StringArray{
		"ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE",
	}

	allMatch = true
	result = &user2.Wallet.Balance.NftBalance.Balance
	for i, n := 0, len(*result); i < n; i++ {
		allMatch = allMatch && ((*result)[i] == user2Exp[i])
	}

	if !allMatch {
		t.Fatalf("failed to user2 result balance: %v %d", result, len(*result))
	}
}
