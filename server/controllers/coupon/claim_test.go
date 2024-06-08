package coupon

import (
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestClaimCoupon(t *testing.T) {
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

	var createCouponRequest CreateCouponRequest
	var userNames []string = []string{"User"}
	var limit int = 2
	var couponTransaction models.CouponTransaction
	var activeCoupon models.ClaimedCoupon

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserNames: &userNames,
		Balance:         1000,
	}
	specCouponCode, err := Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create a coupon for specific users: %v", err)
	}

	if _, err := Claim(users[0].ID, uuid.New().String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		t.Fatalf("should return error: %v but got %v", ErrCodeCouponCodeNotFound, err)
	}

	if _, err := Claim(users[1].ID, specCouponCode.String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeCouponNotAllowedToClaim) {
		t.Fatalf("should return error: %v but got %v", ErrCodeCouponNotAllowedToClaim, err)
	}

	balance, err := Claim(users[0].ID, specCouponCode.String())
	if err != nil {
		t.Fatalf("failed to claim coupon code: %v, %v", specCouponCode, err)
	}
	if balance != 1000 {
		t.Fatalf("coupon claimed not properly: %d", balance)
	}
	if result := db.Where("coupon_id = ? and claimed_user_id = ?", specCouponCode, users[0].ID).First(&activeCoupon); result.Error != nil {
		t.Fatalf("failed to get active coupon: %v", result.Error)
	}
	if activeCoupon.Balance != balance ||
		activeCoupon.CouponID != specCouponCode ||
		activeCoupon.ClaimedUserID != users[0].ID ||
		activeCoupon.Wagered != 0 ||
		activeCoupon.Exchanged != 0 {
		t.Fatalf("coupon claimed not properly: %v", activeCoupon)
	}
	if result := db.Last(&couponTransaction); result.Error != nil {
		t.Fatalf("failed to get coupon transaction history: %v", result.Error)
	}
	if couponTransaction.CouponID != specCouponCode ||
		couponTransaction.ClaimedUserID != users[0].ID ||
		couponTransaction.TxBalance != balance ||
		couponTransaction.NextBalance != couponTransaction.PrevBalance+couponTransaction.TxBalance ||
		couponTransaction.Status != models.CouponTransactionSucceed ||
		couponTransaction.Type != models.CpTxClaimCode {
		t.Fatalf("coupon claimed not properly: %v", couponTransaction)
	}

	if _, err := Claim(users[0].ID, specCouponCode.String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeAlreadyExistingActiveCoupon) {
		t.Fatalf("should return error: %v but got %v", ErrCodeAlreadyExistingActiveCoupon, err)
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

	var limitCoupon models.Coupon
	if result := db.First(&limitCoupon, limitCouponCode); result.Error != nil {
		t.Fatalf("failed to get limit coupon from db: %v", result.Error)
	}
	limitCoupon.CreatedAt = time.UnixMilli(time.Now().UnixMilli() - time.Hour.Milliseconds()*24*15)
	if result := db.Save(&limitCoupon); result.Error != nil {
		t.Fatalf("failed to save updated limit coupon: %v", result.Error)
	}

	if _, err := Claim(users[2].ID, limitCouponCode.String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		t.Fatalf("should return error: %v but got %v", ErrCodeCouponCodeNotFound, err)
	}

	limitCoupon.CreatedAt = time.Now()
	if result := db.Save(&limitCoupon); result.Error != nil {
		t.Fatalf("failed to save updated limit coupon: %v", result.Error)
	}

	balance, err = Claim(users[1].ID, limitCouponCode.String())
	if err != nil {
		t.Fatalf("failed to claim limit coupon code: %v, %v", limitCouponCode, err)
	}
	if balance != 1000 {
		t.Fatalf("limit coupon claimed not properly: %d", balance)
	}

	activeCoupon = models.ClaimedCoupon{}
	if result := db.Where("coupon_id = ? and claimed_user_id = ?", limitCouponCode, users[1].ID).First(&activeCoupon); result.Error != nil {
		t.Fatalf("failed to get active coupon: %v", result.Error)
	}
	if activeCoupon.Balance != balance ||
		activeCoupon.CouponID != limitCouponCode ||
		activeCoupon.ClaimedUserID != users[1].ID ||
		activeCoupon.Wagered != 0 ||
		activeCoupon.Exchanged != 0 {
		t.Fatalf("coupon claimed not properly: %v", activeCoupon)
	}

	couponTransaction = models.CouponTransaction{}
	if result := db.Last(&couponTransaction); result.Error != nil {
		t.Fatalf("failed to get coupon transaction history: %v", result.Error)
	}
	if couponTransaction.CouponID != limitCouponCode ||
		couponTransaction.ClaimedUserID != users[1].ID ||
		couponTransaction.TxBalance != balance ||
		couponTransaction.NextBalance != couponTransaction.PrevBalance+couponTransaction.TxBalance ||
		couponTransaction.Status != models.CouponTransactionSucceed ||
		couponTransaction.Type != models.CpTxClaimCode {
		t.Fatalf(
			"uuid expected: %s, actual: %s, coupon claimed not properly: %v",
			couponTransaction.CouponID,
			limitCouponCode,
			couponTransaction,
		)
	}

	if _, err := Claim(users[0].ID, limitCouponCode.String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeAlreadyExistingActiveCoupon) {
		t.Fatalf("should return error: %v but got %v", ErrCodeAlreadyExistingActiveCoupon, err)
	}

	balance, err = Claim(users[2].ID, limitCouponCode.String())
	if err != nil {
		t.Fatalf("failed to claim limit coupon code: %v, %v", limitCouponCode, err)
	}
	if balance != 1000 {
		t.Fatalf("limit coupon claimed not properly: %d", balance)
	}

	activeCoupon = models.ClaimedCoupon{}
	if result := db.Where("coupon_id = ? and claimed_user_id = ?", limitCouponCode, users[2].ID).First(&activeCoupon); result.Error != nil {
		t.Fatalf("failed to get active coupon: %v", result.Error)
	}
	if activeCoupon.Balance != balance ||
		activeCoupon.CouponID != limitCouponCode ||
		activeCoupon.ClaimedUserID != users[2].ID ||
		activeCoupon.Wagered != 0 ||
		activeCoupon.Exchanged != 0 {
		t.Fatalf("coupon claimed not properly: %v", activeCoupon)
	}

	couponTransaction = models.CouponTransaction{}
	if result := db.Last(&couponTransaction); result.Error != nil {
		t.Fatalf("failed to get coupon transaction history: %v", result.Error)
	}
	if couponTransaction.CouponID != limitCouponCode ||
		couponTransaction.ClaimedUserID != users[2].ID ||
		couponTransaction.TxBalance != balance ||
		couponTransaction.NextBalance != couponTransaction.PrevBalance+couponTransaction.TxBalance ||
		couponTransaction.Status != models.CouponTransactionSucceed ||
		couponTransaction.Type != models.CpTxClaimCode {
		t.Fatalf("coupon claimed not properly: %v", couponTransaction)
	}

	if _, err := Claim(users[3].ID, limitCouponCode.String()); err == nil ||
		!utils.IsErrorCode(err, ErrCodeCouponClaimReachedLimit) {
		t.Fatalf("should return error: %v but got %v", ErrCodeCouponClaimReachedLimit, err)
	}
}

func TestBlockClaimingForPlayingRounds(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to get mock db")
	}
	db_aggregator.Initialize(db)

	accessUserLimit := int(5)
	couponCode, err := Create(CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &accessUserLimit,
		Balance:         100 * 100000,
	})
	if err != nil {
		t.Fatalf("failed to create coupon code: %v", err)
	}

	if err := db.Create(
		&[]models.User{
			{
				Name: "user",
			},
			{
				Model: gorm.Model{
					ID: config.COUPON_TEMP_ID,
				},
				Name: "CP_TEMP",
				Wallet: models.Wallet{
					Balance: models.Balance{
						ChipBalance: &models.ChipBalance{
							Balance: int64(5000 * 100000),
						},
					},
				},
			},
		},
	).Error; err != nil {
		t.Fatalf("failed to create mock user: %v", err)
	}

	if _, err := Claim(
		1, couponCode.String(),
	); err != nil {
		t.Fatalf("failed to claim coupon code: %v", err)
	}

	if err := db.Exec(
		"update claimed_coupons set wagered = ?",
		100*100000*config.COUPON_REQUIRED_WAGER_TIMES,
	).Error; err != nil {
		t.Fatalf("failed to update wagered to required one: %v", err)
	}

	{
		if err := db.Create(
			&models.DreamTowerRound{
				UserID:          1,
				BetAmount:       10 * 100000,
				Status:          models.DreamTowerPlaying,
				PaidBalanceType: models.CouponBalanceForGame,
			},
		).Error; err != nil {
			t.Fatalf("failed to create fake dream tower round: %v", err)
		}
		if _, err := Exchange(
			1, couponCode,
		); !utils.IsErrorCode(
			err,
			ErrCodeExistingPlayingRounds,
		) {
			t.Fatalf("should be failed for existing playing round: %v", err)
		}
		if err := db.Exec(
			"delete from dream_tower_rounds",
		).Error; err != nil {
			t.Fatalf("failed to clean up dream tower rounds: %v", err)
		}
	}

	{
		now := time.Now()
		if err := db.Create(
			&models.CrashRound{
				BetStartedAt: &now,
				Bets: []models.CrashBet{
					{
						UserID:          1,
						BetAmount:       100 * 100000,
						PaidBalanceType: models.CouponBalanceForGame,
					},
				},
			},
		).Error; err != nil {
			t.Fatalf("failed to create fake crash round: %v", err)
		}
		if _, err := Exchange(
			1, couponCode,
		); !utils.IsErrorCode(
			err,
			ErrCodeExistingPlayingRounds,
		) {
			t.Fatalf("should be failed for existing playing round: %v", err)
		}
		if err := db.Exec(
			`
delete from crash_bets;
delete from crash_rounds;
`,
		).Error; err != nil {
			t.Fatalf("failed to clean up crash rounds: %v", err)
		}
	}

	{
		if err := db.Create(
			&models.DreamTowerRound{
				UserID:          1,
				BetAmount:       10 * 100000,
				Status:          models.DreamTowerPlaying,
				PaidBalanceType: models.ChipBalanceForGame,
			},
		).Error; err != nil {
			t.Fatalf("failed to create fake dream tower round: %v", err)
		}

		now := time.Now()
		crashProfit := int64(200 * 100000)
		payoutMul := float64(2)
		if err := db.Create(
			&models.CrashRound{
				BetStartedAt: &now,
				Bets: []models.CrashBet{
					{
						UserID:           1,
						BetAmount:        100 * 100000,
						PaidBalanceType:  models.CouponBalanceForGame,
						Profit:           &crashProfit,
						PayoutMultiplier: &payoutMul,
					},
					{
						UserID:          1,
						BetAmount:       50 * 100000,
						PaidBalanceType: models.ChipBalanceForGame,
					},
				},
			},
		).Error; err != nil {
			t.Fatalf("failed to create fake crash round: %v", err)
		}
		if _, err := Exchange(
			1, couponCode,
		); err != nil {
			t.Fatalf("failed to perform exchange: %v", err)
		}
	}
}
