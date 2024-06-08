package coupon

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func TestTryBet(t *testing.T) {
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

	if result, tx, err := TryBet(TryBetWithCouponRequest{
		UserID:  0,
		Balance: 0,
		Type:    models.CpTxCoinflipBet,
	}); result != CouponBetFailed ||
		tx != 0 ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("expected result: %v, actual result: %v, error: %v", CouponBetFailed, result, err)
	}

	if result, tx, err := TryBet(TryBetWithCouponRequest{
		UserID:  users[0].ID,
		Balance: 1000,
		Type:    models.CpTxCoinflipBet,
	}); result != CouponBetUnavailable ||
		tx != 0 ||
		!utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		t.Fatalf("expected result: %v, actual result: %v, error: %v", CouponBetUnavailable, result, err)
	}

	var createCouponRequest CreateCouponRequest
	var userNames []string = []string{"User"}
	var bonusBalance = int64(10000)

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserNames: &userNames,
		Balance:         bonusBalance,
	}
	specCouponCode, err := Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create a coupon for specific users: %v", err)
	}

	balance, err := Claim(users[0].ID, specCouponCode.String())
	if err != nil {
		t.Fatalf("failed to claim coupon code: %v, %v", specCouponCode, err)
	}
	if balance != bonusBalance {
		t.Fatalf("coupon claimed not properly: %d", balance)
	}

	if result, tx, err := TryBet(TryBetWithCouponRequest{
		UserID:  users[0].ID,
		Balance: 100000,
		Type:    models.CpTxCoinflipBet,
	}); result != CouponBetInsufficientFunds ||
		tx != 0 ||
		!utils.IsErrorCode(err, ErrCodeInsufficientBonusBalance) {
		t.Fatalf("expected result: %v, actual result: %v, error: %v", CouponBetUnavailable, result, err)
	}

	result, tx, err := TryBet(TryBetWithCouponRequest{
		UserID:  users[0].ID,
		Balance: 1000,
		Type:    models.CpTxCoinflipBet,
	})
	if result != CouponBetSucceed || err != nil {
		t.Fatalf("failed to bet with coupon balance: result: %v, error: %v", result, err)
	}

	var couponTransaction models.CouponTransaction
	if result := db.First(&couponTransaction, tx); result.Error != nil {
		t.Fatalf("failed to get coupon transaction from db: %v", result.Error)
	}

	if couponTransaction.Type != models.CpTxCoinflipBet ||
		couponTransaction.PrevBalance != bonusBalance ||
		couponTransaction.TxBalance != 1000 ||
		couponTransaction.NextBalance != bonusBalance-1000 ||
		couponTransaction.Status != models.CouponTransactionPending {
		t.Fatalf("coupont transaction recorded not properly: %v", couponTransaction)
	}

	if err := Confirm(tx); err != nil {
		t.Fatalf("failed to confirm transaction: %v", err)
	}

	couponTransaction = models.CouponTransaction{}
	if result := db.First(&couponTransaction, tx); result.Error != nil {
		t.Fatalf("failed to get coupon transaction from db: %v", result.Error)
	}
	if couponTransaction.Status != models.CouponTransactionSucceed {
		t.Fatalf("transaction confirmed not properly: %v", couponTransaction.Status)
	}

	result, tx, err = TryBet(TryBetWithCouponRequest{
		UserID:  users[0].ID,
		Balance: 1000,
		Type:    models.CpTxDreamtowerBet,
	})
	if result != CouponBetSucceed || err != nil {
		t.Fatalf("failed to bet with coupon balance: result: %v, error: %v", result, err)
	}

	couponTransaction = models.CouponTransaction{}
	if result := db.First(&couponTransaction, tx); result.Error != nil {
		t.Fatalf("failed to get coupon transaction from db: %v", result.Error)
	}

	if couponTransaction.Type != models.CpTxDreamtowerBet ||
		couponTransaction.PrevBalance != bonusBalance-1000 ||
		couponTransaction.TxBalance != 1000 ||
		couponTransaction.NextBalance != bonusBalance-2000 ||
		couponTransaction.Status != models.CouponTransactionPending {
		t.Fatalf("coupont transaction recorded not properly: %v", couponTransaction)
	}

	if err := Decline(tx); err != nil {
		t.Fatalf("failed to decline transaction: %v", err)
	}

	couponTransaction = models.CouponTransaction{}
	if result := db.First(&couponTransaction, tx); result.Error != nil {
		t.Fatalf("failed to get coupon transaction from db: %v", result.Error)
	}
	if couponTransaction.Status != models.CouponTransactionFailed ||
		couponTransaction.AfterRefund != bonusBalance-1000 {
		t.Fatalf("transaction confirmed not properly: %v", couponTransaction.Status)
	}

	var claimedCoupon models.ClaimedCoupon
	if result := db.Where("coupon_id = ?", specCouponCode).Where("claimed_user_id = ?", users[0].ID).First(&claimedCoupon); result.Error != nil {
		t.Fatalf("failed to get claimed coupon from db: %v", result.Error)
	}
	if claimedCoupon.Balance != bonusBalance-1000 {
		t.Fatalf("expected: %d, actual: %d", bonusBalance-1000, claimedCoupon.Balance)
	}
	if claimedCoupon.Wagered != 1000 {
		t.Fatalf("expected: %d, actual: %d", 1000, claimedCoupon.Wagered)
	}

	{
		result, tx, err := TryBet(TryBetWithCouponRequest{
			UserID:  users[0].ID,
			Balance: 9000,
			Type:    models.CpTxDreamtowerBet,
		})
		if result != CouponBetSucceed || err != nil {
			t.Fatalf("failed to try bet: result: %v, error: %v", result, err)
		}

		if result, _, err := TryBet(TryBetWithCouponRequest{
			UserID:  users[0].ID,
			Balance: 1000,
			Type:    models.CpTxDreamtowerBet,
		}); result != CouponBetUnavailable ||
			!utils.IsErrorCode(err, ErrCodeZeroBonusBalance) {
			t.Fatalf(
				"failed to try bet properly: result: %v, err: %v",
				result, err,
			)
		}

		if err := Decline(tx); err != nil {
			t.Fatalf("failed to decline tx: err: %v", err)
		}

		if result, _, err := TryBet(TryBetWithCouponRequest{
			UserID:  users[0].ID,
			Balance: 9000,
			Type:    models.CpTxDreamtowerBet,
		}); result != CouponBetSucceed || err != nil {
			t.Fatalf("failed to try bet: result: %v, error: %v", result, err)
		}
	}
}
