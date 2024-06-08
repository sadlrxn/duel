package transaction

import (
	"strings"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestActivationTimeline(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to get mock db")
	}

	db_aggregator.Initialize(db)

	user := models.User{
		Name: "user",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create mock user: %v", err)
	}

	affiliate := models.Affiliate{
		Code: "affiliate-code",
		Creator: models.User{
			Name: "creator",
		},
	}
	if err := db.Create(&affiliate).Error; err != nil {
		t.Fatalf("failed to create mock affiliate: %v", err)
	}

	if err := db.Exec(
		"update users set created_at = ? where name = 'user';",
		time.Now().Add(-time.Hour*time.Duration(config.AFFILIATE_ACTIVATION_TIMELINE_IN_HOURS+1)),
	).Error; err != nil {
		t.Fatalf("failed to perform update created_at query: %v", err)
	}

	if _, err := activateAffiliateCode(
		db_aggregator.User(1),
		"affiliate-code",
	); err == nil ||
		!strings.Contains(
			err.Error(),
			"possibly expired affiliate code activation timeline",
		) {
		t.Fatalf("should be failed on affiliate activation after 24 hours: %v", err)
	}

	if err := db.Exec(
		"update users set created_at = ? where name = 'user';",
		time.Now(),
	).Error; err != nil {
		t.Fatalf("failed to perform update created_at query: %v", err)
	}

	if _, err := activateAffiliateCode(
		db_aggregator.User(1),
		"affiliate-code",
	); err != nil {
		t.Fatalf("failed to activate code: %v", err)
	}
}

func TestApplyForFirstDepositBonus(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to init mock db")
	}

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("failed to initialize db aggregator; %v", err)
	}

	if err := redis.InitializeMockRedis(true); err != nil {
		t.Fatalf("failed to init mock redis:%v", err)
	}

	var userID uint
	var depositAmount int64 = 1000

	{
		userID = 0
		if bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount); bonus != 0 || err != nil {
			t.Fatalf("should return zero bonus since userID is zero")
		}
	}

	{
		userID = 1
		depositAmount = -1000
		if bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount); bonus != 0 || err != nil {
			t.Fatalf("should return zero bonus since depositAmount is less than zero")
		}
	}

	users := []models.User{
		{
			Name:          "User1",
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
		},
		{
			Name:          "User3",
			WalletAddress: "EvPpQ4TQHHFxsXjSaBWKZavvhXXwCLRv25LbMBfYmZGN",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 0,
					},
				},
			},
		},
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	if err := CreateAffiliateCode(db_aggregator.User(users[0].ID), []string{"affiliate_code_1"}); err != nil {
		t.Fatalf("failed to create affiliate code: %v", err)
	}

	if err := CreateAffiliateCode(db_aggregator.User(users[1].ID), []string{"affiliate_code_2"}); err != nil {
		t.Fatalf("failed to create affiliate code: %v", err)
	}

	if err := UpdateAffiliateFirstDepositBonus("affiliate_code_1", true); err != nil {
		t.Fatalf("failed to update affiliate first deposit bonus: %v", err)
	}

	{
		userID = users[1].ID
		depositAmount = 1000
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if bonus != 0 || err != nil {
			t.Fatalf("should return zero bonus and nil error: bonus: %d, err: %v", bonus, err)
		}
	}

	if _, err := ActivateAffiliateCode(db_aggregator.User(users[1].ID), "affiliate_code_1"); err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}

	{
		activeAffiliate, err := db_aggregator.GetActiveAffiliateCode(
			db_aggregator.User(users[1].ID),
		)
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}

		if activeAffiliate.FirstDepositDone || !activeAffiliate.IsFirstDepositBonus {
			t.Fatalf("active affiliate is not set properly: done: %v, isFirst: %v",
				activeAffiliate.FirstDepositDone,
				activeAffiliate.IsFirstDepositBonus,
			)
		}

		userID = users[1].ID
		depositAmount = 1000
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if err != nil {
			t.Fatalf("failed to apply for first deposit bonus: %v", err)
		}
		if bonus != depositAmount {
			t.Fatalf("incorrect bonus balance: depositAmount = %d, bonus = %d",
				depositAmount,
				bonus,
			)
		}

		activeAffiliate, err = db_aggregator.GetActiveAffiliateCode(
			db_aggregator.User(users[1].ID),
		)
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}

		if !activeAffiliate.FirstDepositDone {
			t.Fatalf("active affiliate is not set properly: done: %v",
				activeAffiliate.FirstDepositDone,
			)
		}
	}

	{
		userID = users[1].ID
		depositAmount = 1000
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if err != nil {
			t.Fatalf("failed to apply for first deposit bonus: %v", err)
		}
		if bonus != 0 {
			t.Fatalf("incorrect bonus balance: bonus = %d",
				bonus,
			)
		}
	}

	if _, err := ActivateAffiliateCode(db_aggregator.User(users[0].ID), "affiliate_code_2"); err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}

	{
		userID = users[1].ID
		depositAmount = 1000
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if bonus != 0 || err != nil {
			t.Fatalf("should return zero bonus and nil error: bonus: %d, err: %v", bonus, err)
		}
	}

	if err := UpdateAffiliateFirstDepositBonus("affiliate_code_2", true); err != nil {
		t.Fatalf("failed to update affiliate first deposit bonus: %v", err)
	}

	{
		activeAffiliate, err := db_aggregator.GetActiveAffiliateCode(
			db_aggregator.User(users[0].ID),
		)
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}

		if activeAffiliate.FirstDepositDone || !activeAffiliate.IsFirstDepositBonus {
			t.Fatalf("active affiliate is not set properly: done: %v, isFirst: %v",
				activeAffiliate.FirstDepositDone,
				activeAffiliate.IsFirstDepositBonus,
			)
		}

		userID = users[0].ID
		depositAmount = config.COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS + 1
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if err != nil {
			t.Fatalf("failed to apply for first deposit bonus: %v", err)
		}
		if bonus != config.COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS {
			t.Fatalf("incorrect bonus balance: expected = %d, result = %d",
				config.COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS,
				bonus,
			)
		}

		activeAffiliate, err = db_aggregator.GetActiveAffiliateCode(
			db_aggregator.User(users[0].ID),
		)
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}

		if !activeAffiliate.FirstDepositDone {
			t.Fatalf("active affiliate is not set properly: done: %v",
				activeAffiliate.FirstDepositDone,
			)
		}
	}

	{
		code, err := coupon.Create(
			coupon.CreateCouponRequest{
				Type:          models.CouponForSpecUsers,
				AccessUserIDs: &[]uint{users[2].ID},
				Balance:       1000,
			},
		)
		if err != nil {
			t.Fatalf("failed to create coupon code: %v", err)
		}

		_, err = coupon.Claim(
			users[2].ID,
			code.String(),
		)
		if err != nil {
			t.Fatalf("failed to redeem coupon code: %v", err)
		}
	}

	if _, err := ActivateAffiliateCode(db_aggregator.User(users[2].ID), "affiliate_code_1"); err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}

	{
		activeAffiliate, err := db_aggregator.GetActiveAffiliateCode(
			db_aggregator.User(users[2].ID),
		)
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}

		if activeAffiliate.FirstDepositDone || !activeAffiliate.IsFirstDepositBonus {
			t.Fatalf("active affiliate is not set properly: done: %v, isFirst: %v",
				activeAffiliate.FirstDepositDone,
				activeAffiliate.IsFirstDepositBonus,
			)
		}

		userID = users[2].ID
		depositAmount = 1000
		bonus, err := TryApplyForFirstDepositBonus(userID, depositAmount)
		if err == nil || bonus != 0 {
			t.Fatalf("should return error since already redeemed another coupon: %v, %d",
				err,
				bonus,
			)
		}
	}

}
