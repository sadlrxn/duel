package db_aggregator

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestGenerateRakebackInfo(t *testing.T) {
	db := tests.InitMockDB(true, true)
	initialize(db)

	user, _ := getMockUser()
	if result := db.Create(&user); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	rakebackInfo, err := generateRakebackInfo(User(user.ID))
	if err != nil {
		t.Fatalf("failed to generate rakeback info: %v", err)
	}

	if rakebackInfo.ID != 1 {
		t.Fatalf("failed to generate rakeback properly: %d", rakebackInfo.ID)
	}

	rakeback := models.Rakeback{}
	if result := db.First(&rakeback); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback: %v", result.Error)
	}

	if rakeback.ID != 1 ||
		rakeback.UserID != 1 ||
		rakeback.TotalEarned != 0 ||
		rakeback.Reward != 0 {
		t.Fatalf("failed to generate rakeback properly: %v", rakeback)
	}
}

func TestDistributeRakeback(t *testing.T) {
	TestGenerateRakebackInfo(t)

	db := tests.InitMockDB(false, false)

	distributed, err := distributeRakeback(User(1), 1000)
	if err != nil {
		t.Fatalf("failed to distribute rakeback properly: %v", err)
	}
	if distributed != 50 {
		t.Fatalf("failed to distribute proper amount of rakeback: %v", distributed)
	}

	rakeback := models.Rakeback{}
	if result := db.First(&rakeback); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback: %v", result.Error)
	}

	if rakeback.ID != 1 ||
		rakeback.UserID != 1 ||
		rakeback.TotalEarned != 50 ||
		rakeback.Reward != 50 {
		t.Fatalf("failed to distribute rakeback properly: %v", rakeback)
	}
}

func TestClaimRakeback(t *testing.T) {
	TestDistributeRakeback(t)

	db := tests.InitMockDB(false, false)

	amount, err := claimRakeback(User(1))
	if err != nil {
		t.Fatalf("failed to claim rakeback: %v", err)
	}

	if amount != 50 {
		t.Fatalf("failed to claim proper amount: %d", amount)
	}

	rakeback := models.Rakeback{}
	if result := db.First(&rakeback); result.Error != nil {
		t.Fatalf("failed to retrieve rakeback: %v", result.Error)
	}

	if rakeback.ID != 1 ||
		rakeback.UserID != 1 ||
		rakeback.TotalEarned != 50 ||
		rakeback.Reward != 0 {
		t.Fatalf("failed to generate rakeback properly: %v", rakeback)
	}

	// ================== Archived for transaction ==================
	transactionLoad := models.Transaction{}
	if result := db.Preload(
		"Balance.ChipBalance",
	).Preload(
		"ToWalletPrev.ChipBalance",
	).Preload(
		"ToWalletNext.ChipBalance",
	).First(
		&transactionLoad,
	); result.Error != nil {
		t.Fatalf("failed to retrieve transaction load: %v", result.Error)
	}

	if transactionLoad.Balance.ChipBalance.Balance != 50 ||
		transactionLoad.FromWallet != nil ||
		transactionLoad.ToWalletPrev.ChipBalance.Balance != 100 ||
		transactionLoad.ToWalletNext.ChipBalance.Balance != 150 ||
		*transactionLoad.ToWallet != 1 ||
		transactionLoad.Type != models.TxClaimRakebackReward ||
		transactionLoad.Status != models.TransactionSucceed {
		t.Fatalf("failed to record transaction properly: %v", transactionLoad)
	}
	// ================== Archived for transaction ==================

	userInfo := models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).First(&userInfo); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if userInfo.Wallet.Balance.ChipBalance.Balance != 150 {
		t.Fatalf(
			"failed to claim and update balance properly: %d",
			userInfo.Wallet.Balance.ChipBalance.Balance,
		)
	}
}
