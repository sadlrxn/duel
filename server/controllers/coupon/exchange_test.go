package coupon

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
)

func TestExchange(t *testing.T) {
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
		{
			Model: gorm.Model{ID: config.COUPON_TEMP_ID},
			Name:  "CP_TEMP",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 1000,
					},
				},
			},
		},
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	var createCouponRequest CreateCouponRequest
	var userNames []string = []string{"User"}
	var limit int = 2

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserNames: &userNames,
		Balance:         1000,
	}
	specCouponCode, err := Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create a coupon for specific users: %v", err)
	}

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &limit,
		Balance:         1000,
	}
	limitCouponCode, err := Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create a coupon for limited number of users: %v", err)
	}

	balance, err := Claim(users[0].ID, specCouponCode.String())
	if err != nil {
		t.Fatalf("failed to claim coupon code: %v, %v", specCouponCode, err)
	}
	if balance != 1000 {
		t.Fatalf("coupon claimed not properly: %d", balance)
	}

	if _, err := Exchange(users[0].ID, limitCouponCode); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotActiveCodeForExchange) {
		t.Fatalf("should return error: %v but got %v", ErrCodeNotActiveCodeForExchange, err)
	}

	if _, err := Exchange(users[0].ID, specCouponCode); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotReachingExchangeWager) {
		t.Fatalf("should return error: %v but got %v", ErrCodeNotReachingExchangeWager, err)
	}

	var claimedCoupon models.ClaimedCoupon
	if result := db.Where("coupon_id = ?", specCouponCode).Where("claimed_user_id = ?", users[0].ID).First(&claimedCoupon); result.Error != nil {
		t.Fatalf("failed to get claimed coupon: %v", result.Error)
	}

	claimedCoupon.Wagered = int64(config.COUPON_REQUIRED_WAGER_TIMES) * balance
	if result := db.Save(&claimedCoupon); result.Error != nil {
		t.Fatalf("failed to save updated claimed coupon: %v", result.Error)
	}

	amount, err := Exchange(users[0].ID, specCouponCode)
	if err != nil {
		t.Fatalf("failed to exchange coupon balance to real chips: %v", err)
	}
	if amount != balance {
		t.Fatalf("failed to exchange coupon: expected: %d, result: %d", balance, amount)
	}

	if result := db.Where("coupon_id = ?", specCouponCode).Where("claimed_user_id = ?", users[0].ID).First(&claimedCoupon); result.Error != nil {
		t.Fatalf("failed to get claimed coupon: %v", result.Error)
	}
	if claimedCoupon.Exchanged != balance {
		t.Fatalf("failed to exchange coupon: expected: %d, exchanged: %d", balance, claimedCoupon.Exchanged)
	}

	var couponTransaction models.CouponTransaction
	if result := db.Preload("RealTransaction.Balance.ChipBalance").Last(&couponTransaction); result.Error != nil {
		t.Fatalf("failed to get coupon transaction: %v", result.Error)
	}
	if couponTransaction.Type != models.CpTxExchangeToChip ||
		couponTransaction.RealTransaction.Type != models.TxExchangeCouponToChips ||
		*couponTransaction.RealTransaction.FromWallet != users[4].Wallet.ID ||
		*couponTransaction.RealTransaction.ToWallet != users[0].Wallet.ID ||
		couponTransaction.RealTransaction.Balance.ChipBalance.Balance != balance ||
		couponTransaction.CouponID != specCouponCode ||
		couponTransaction.ClaimedUserID != users[0].ID ||
		couponTransaction.TxBalance != balance ||
		couponTransaction.PrevBalance-couponTransaction.TxBalance != couponTransaction.NextBalance ||
		couponTransaction.Status != models.CouponTransactionSucceed ||
		couponTransaction.RealTransaction.OwnerID != couponTransaction.ID ||
		couponTransaction.RealTransaction.OwnerType != models.TransactionCouponTransactionReferenced {
		t.Fatalf("exchange transaction not recorded properly: %v", couponTransaction)
	}

	{
		accessUserLimit := int(5)
		code, err := Create(CreateCouponRequest{
			Type:            models.CouponForLimitUsers,
			AccessUserLimit: &accessUserLimit,
			Balance:         1000,
		})
		if err != nil {
			t.Fatalf("failed to create code: %v", err)
		}

		claimed, err := Claim(users[1].ID, code.String())
		if err != nil {
			t.Fatalf("failed to claime code: %v", err)
		}

		activeCoupon, err := lockAndRetrieveActiveCoupon(users[1].ID, db_aggregator.MainSessionId())
		if err != nil {
			t.Fatalf("failed to get active claimed coupon record: %v", err)
		}

		if err := addWagerAmountUnchecked(
			claimed*int64(config.COUPON_REQUIRED_WAGER_TIMES),
			activeCoupon,
			db_aggregator.MainSessionId(),
		); err != nil {
			t.Fatalf("failed to add wager amount to claimed coupon: %v", err)
		}

		if _, err := Exchange(users[1].ID, code); err == nil ||
			!utils.IsErrorCode(err, ErrCodeInsufficientAdminBalance) {
			t.Fatalf("should be failed on insufficient admin balance: %v", err)
		}
	}
}
