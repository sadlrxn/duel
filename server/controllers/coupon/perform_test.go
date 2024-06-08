package coupon

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestPerform(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("retrieved mock db is nil pointer")
	}

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("failed to initialize db aggregator: %v", err)
	}

	users := []models.User{
		{
			Name: "User1",
		},
		{
			Name: "User2",
		},
		{
			Name: "CF_TEMP",
		},
		{
			Name: "DT_TEMP",
		},
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	couponCode, err := Create(
		CreateCouponRequest{
			Type:            models.CouponForSpecUsers,
			AccessUserNames: &[]string{"User1", "User2"},
			Balance:         50000,
		},
	)
	if err != nil {
		t.Fatalf("failed to create coupon code: %v", err)
	}

	if claimed, err := Claim(users[0].ID, couponCode.String()); err != nil ||
		claimed != 50000 {
		t.Fatalf(
			"failed to claim coupon balance: err: %v, claimed: %d",
			err, claimed,
		)
	}
	if claimed, err := Claim(users[1].ID, couponCode.String()); err != nil ||
		claimed != 50000 {
		t.Fatalf(
			"failed to claim for second user: err: %v, claimed: %d",
			err, claimed,
		)
	}

	{
		txId, err := Perform(
			CouponTransactionRequest{
				Type:          models.CpTxCoinflipBet,
				UserID:        users[0].ID,
				Balance:       1000,
				ToBeConfirmed: false,
			},
		)
		if err != nil {
			t.Fatalf(
				"failed to perform coinflip bet: %v",
				err,
			)
		}

		cpTransaction := models.CouponTransaction{}
		if result := db.First(&cpTransaction, txId); result.Error != nil {
			t.Fatalf(
				"failed to retrieve coupon transaction: %v",
				result.Error,
			)
		}

		if cpTransaction.CouponID != couponCode ||
			cpTransaction.Type != models.CpTxCoinflipBet ||
			cpTransaction.ClaimedUserID != users[0].ID ||
			cpTransaction.Status != models.CouponTransactionPending ||
			cpTransaction.PrevBalance != 50000 ||
			cpTransaction.TxBalance != 1000 ||
			cpTransaction.NextBalance != 49000 ||
			cpTransaction.AfterRefund != 0 {
			t.Fatalf("failed to leave coupon transaction properly: %v", cpTransaction)
		}

		user1Coupon := GetActiveUserCoupon(users[0].ID)
		if user1Coupon == nil ||
			user1Coupon.Balance != 49000 || // 49000
			user1Coupon.Code != couponCode ||
			user1Coupon.Wagered != 0 ||
			user1Coupon.Claimed != 50000 || // 50000
			user1Coupon.WagerLimit != 50000*int64(config.COUPON_REQUIRED_WAGER_TIMES) /* 50000*30 */ {
			t.Fatalf("failed to get user active coupon properly: %v", user1Coupon)
		}

		if err := Confirm(txId); err != nil {
			t.Fatalf("failed to confirm transaction: %v", err)
		}

		cpTransaction = models.CouponTransaction{}
		if result := db.First(&cpTransaction, txId); result.Error != nil {
			t.Fatalf(
				"failed to retrieve coupon transaction: %v",
				result.Error,
			)
		}

		if cpTransaction.CouponID != couponCode ||
			cpTransaction.Type != models.CpTxCoinflipBet ||
			cpTransaction.ClaimedUserID != users[0].ID ||
			cpTransaction.Status != models.CouponTransactionSucceed ||
			cpTransaction.PrevBalance != 50000 ||
			cpTransaction.TxBalance != 1000 ||
			cpTransaction.NextBalance != 49000 ||
			cpTransaction.AfterRefund != 0 {
			t.Fatalf("failed to confirm coupon transaction properly: %v", cpTransaction)
		}

		user1Coupon = GetActiveUserCoupon(users[0].ID)
		if user1Coupon == nil ||
			user1Coupon.Balance != 49000 ||
			user1Coupon.Code != couponCode ||
			user1Coupon.Wagered != 1000 ||
			user1Coupon.Claimed != 50000 ||
			user1Coupon.WagerLimit != 50000*int64(config.COUPON_REQUIRED_WAGER_TIMES) {
			t.Fatalf("failed to get user active coupon properly: %v", user1Coupon)
		}
	}

	{
		txId, err := Perform(
			CouponTransactionRequest{
				Type:          models.CpTxDreamtowerBet,
				UserID:        users[1].ID,
				Balance:       1000,
				ToBeConfirmed: false,
			},
		)
		if err != nil {
			t.Fatalf(
				"failed to perform coinflip bet: %v",
				err,
			)
		}

		cpTransaction := models.CouponTransaction{}
		if result := db.First(&cpTransaction, txId); result.Error != nil {
			t.Fatalf(
				"failed to retrieve coupon transaction: %v",
				result.Error,
			)
		}

		if cpTransaction.CouponID != couponCode ||
			cpTransaction.Type != models.CpTxDreamtowerBet ||
			cpTransaction.ClaimedUserID != users[1].ID ||
			cpTransaction.Status != models.CouponTransactionPending ||
			cpTransaction.PrevBalance != 50000 ||
			cpTransaction.TxBalance != 1000 ||
			cpTransaction.NextBalance != 49000 ||
			cpTransaction.AfterRefund != 0 {
			t.Fatalf("failed to leave coupon transaction properly: %v", cpTransaction)
		}

		user1Coupon := GetActiveUserCoupon(users[1].ID)
		if user1Coupon == nil ||
			user1Coupon.Balance != 49000 ||
			user1Coupon.Code != couponCode ||
			user1Coupon.Wagered != 0 ||
			user1Coupon.Claimed != 50000 ||
			user1Coupon.WagerLimit != 50000*int64(config.COUPON_REQUIRED_WAGER_TIMES) {
			t.Fatalf("failed to get user active coupon properly: %v", user1Coupon)
		}

		if err := Decline(txId); err != nil {
			t.Fatalf("failed to confirm transaction: %v", err)
		}

		cpTransaction = models.CouponTransaction{}
		if result := db.First(&cpTransaction, txId); result.Error != nil {
			t.Fatalf(
				"failed to retrieve coupon transaction: %v",
				result.Error,
			)
		}

		if cpTransaction.CouponID != couponCode ||
			cpTransaction.Type != models.CpTxDreamtowerBet ||
			cpTransaction.ClaimedUserID != users[1].ID ||
			cpTransaction.Status != models.CouponTransactionFailed ||
			cpTransaction.PrevBalance != 50000 ||
			cpTransaction.TxBalance != 1000 ||
			cpTransaction.NextBalance != 49000 ||
			cpTransaction.AfterRefund != 50000 {
			t.Fatalf("failed to confirm coupon transaction properly: %v", cpTransaction)
		}

		user1Coupon = GetActiveUserCoupon(users[1].ID)
		if user1Coupon == nil ||
			user1Coupon.Balance != 50000 ||
			user1Coupon.Code != couponCode ||
			user1Coupon.Wagered != 0 ||
			user1Coupon.Claimed != 50000 ||
			user1Coupon.WagerLimit != 50000*int64(config.COUPON_REQUIRED_WAGER_TIMES) {
			t.Fatalf("failed to get user active coupon properly: %v", user1Coupon)
		}
	}
}
