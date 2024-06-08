package db_aggregator

import (
	"strings"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
)

func getMockUsers() []models.User {
	return []models.User{
		{
			Name: "taka1",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 100,
					},
				},
			},
			Statistics: models.Statistics{
				TotalWagered: 1,
			},
		},
		{
			Name: "taka2",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 70,
					},
				},
			},
			Statistics: models.Statistics{
				TotalWagered: int64(config.AFFILIATE_WAGER_LIMIT_FOR_CREATION),
			},
		},
		{
			Name: "taka3",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 50,
					},
				},
			},
		},
	}
}

func initiateTestDB() error {
	db := tests.InitMockDB(true, true)
	if err := initialize(db); err != nil {
		return utils.MakeError(
			"test affiliate db aggregator",
			"initiateTestDB",
			"failed to initialize db session",
			err,
		)
	}

	mockUsers := getMockUsers()
	if result := db.Create(mockUsers); result.Error != nil {
		return utils.MakeError(
			"test affiliate db aggregator",
			"initiateTestDB",
			"failed to create mock users",
			result.Error,
		)
	}

	if _, err := getRakebackInfo(User(1)); err != nil {
		return utils.MakeError(
			"test affiliate db aggregator",
			"initiateTestDB",
			"failed to generate rakeback info for the first user",
			err,
		)
	}
	if _, err := getRakebackInfo(User(2)); err != nil {
		return utils.MakeError(
			"test affiliate db aggregator",
			"initiateTestDB",
			"failed to generate rakeback info for the second user",
			err,
		)
	}

	return nil
}

func TestCreateAffiliateCode(t *testing.T) {
	if err := initiateTestDB(); err != nil {
		t.Fatalf("failed to init mock db: %v", err)
	}

	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}
	session, err := getSession(sessionId)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if err := createAffiliateCode(User(1), []string{"code1", "code2"}, sessionId); err == nil ||
		!strings.Contains(err.Error(), "not enough wager amount") {
		t.Fatalf("should be failed on not enough wager amount: %v", err)
	}

	if err := createAffiliateCode(User(1), []string{"code "}, sessionId); err == nil ||
		!strings.Contains(err.Error(), "space") {
		t.Fatalf("should be failed for creating code with space: %v", err)
	}
	if err := createAffiliateCode(User(1), []string{"cduelt"}, sessionId); err == nil ||
		!strings.Contains(err.Error(), "reserved word") {
		t.Fatalf("should be failed for reserved words: %v", err)
	}

	if result := session.Model(
		&models.Statistics{},
	).Where(
		"user_id = ?",
		1,
	).Update(
		"total_wagered",
		config.AFFILIATE_WAGER_LIMIT_FOR_CREATION,
	); result.Error != nil {
		t.Fatalf("failed to update wager amount of first user: %v", result.Error)
	}

	if err := createAffiliateCode(User(1), []string{"code1", "code2"}, sessionId); err != nil {
		t.Fatalf("failed to create affiliate code: %v", err)
	}

	if err := createAffiliateCode(User(1), []string{"code1"}, sessionId); err == nil ||
		!strings.Contains(err.Error(), "duplicated code") {
		t.Fatalf("should fail with duplicated code: %v", err)
	}

	codes, err := getOwnedAffiliateCode(User(1), sessionId)
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes meta: %v", err)
	}

	if len(codes) != 2 ||
		codes[0].Code != "code1" ||
		codes[0].Reward != 0 ||
		codes[0].TotalEarned != 0 ||
		codes[0].UserCnt != 0 ||
		codes[1].Code != "code2" ||
		codes[1].Reward != 0 ||
		codes[1].TotalEarned != 0 ||
		codes[1].UserCnt != 0 {
		t.Fatalf("failed to create affiliate codes properly: %v", codes)
	}

	if err := createAffiliateCode(User(2), []string{"CoDe1"}, sessionId); err == nil ||
		!strings.Contains(err.Error(), "duplicated code") {
		t.Fatalf("should fail with duplicated code: %v", err)
	}

	if err := commitSession(sessionId); err != nil {
		t.Fatalf("failed to commit session: %v", err)
	}
}

func TestActivateAffiliateCode(t *testing.T) {
	TestCreateAffiliateCode(t)

	{
		if _, err := activateAffiliateCode(User(2), "code"); err == nil ||
			!strings.Contains(err.Error(), "failed to retrieve affiliate") {
			t.Fatalf("should be failed on not found code")
		}

		isFirst, err := activateAffiliateCode(User(2), "CodE1")
		if err != nil {
			t.Fatalf("failed to activate affiliate code: %v", err)
		}
		if !isFirst {
			t.Fatalf("should be the first activation: %v", isFirst)
		}

		codes, err := getOwnedAffiliateCode(User(1))
		if err != nil {
			t.Fatalf("failed to get owned affiliate codes meta: %v", err)
		}

		if len(codes) != 2 ||
			codes[0].Code != "code1" ||
			codes[0].Reward != 0 ||
			codes[0].TotalEarned != 0 ||
			codes[0].UserCnt != 1 ||
			codes[1].Code != "code2" ||
			codes[1].Reward != 0 ||
			codes[1].TotalEarned != 0 ||
			codes[1].UserCnt != 0 {
			t.Fatalf("failed to create affiliate codes properly: %v", codes)
		}

		_, err = activateAffiliateCode(User(1), "cOde2")
		if err == nil || !strings.Contains(err.Error(), "cannot activate your own code") {
			t.Fatalf("should fail in case trying to activate your own code:%v", err)
		}

		activeAffiliate, err := getActiveAffiliateCode(User(2))
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}
		if activeAffiliate.Code != "code1" {
			t.Fatalf("failed to get correct active affiliate code: %v", *activeAffiliate)
		}
	}

	{
		isFirst, err := activateAffiliateCode(User(3), "coDe2")
		if err != nil {
			t.Fatalf("failed to activate affiliate code: %v", err)
		}
		if !isFirst {
			t.Fatalf("should be the first activation: %v", isFirst)
		}

		activeAffiliate, err := getActiveAffiliateCode(User(3))
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}
		if activeAffiliate.Code != "code2" {
			t.Fatalf("failed to get correct active affiliate code: %v", *activeAffiliate)
		}

		codes, err := getOwnedAffiliateCode(User(1))
		if err != nil {
			t.Fatalf("failed to get owned affiliate codes meta: %v", err)
		}

		if len(codes) != 2 ||
			codes[0].Code != "code1" ||
			codes[0].Reward != 0 ||
			codes[0].TotalEarned != 0 ||
			codes[0].UserCnt != 1 ||
			codes[1].Code != "code2" ||
			codes[1].Reward != 0 ||
			codes[1].TotalEarned != 0 ||
			codes[1].UserCnt != 1 {
			t.Fatalf("failed to create affiliate codes properly: %v", codes)
		}
	}

	{
		isFirst, err := activateAffiliateCode(User(3), "code1")
		if err != nil {
			t.Fatalf("failed to activate affiliate code: %v", err)
		}
		if isFirst {
			t.Fatalf("should not be the first activation: %v", isFirst)
		}

		activeAffiliate, err := getActiveAffiliateCode(User(3))
		if err != nil {
			t.Fatalf("failed to get active affiliate code: %v", err)
		}
		if activeAffiliate.Code != "code1" {
			t.Fatalf("failed to get correct active affiliate code: %v", *activeAffiliate)
		}

		codes, err := getOwnedAffiliateCode(User(1))
		if err != nil {
			t.Fatalf("failed to get owned affiliate codes meta: %v", err)
		}

		if len(codes) != 2 ||
			codes[0].Code != "code1" ||
			codes[0].Reward != 0 ||
			codes[0].TotalEarned != 0 ||
			codes[0].UserCnt != 2 ||
			codes[1].Code != "code2" ||
			codes[1].Reward != 0 ||
			codes[1].TotalEarned != 0 ||
			codes[1].UserCnt != 0 {
			t.Fatalf("failed to create affiliate codes properly: %v", codes)
		}
	}

	// {
	// 	err := deactivateAffiliateCode(User(3), "code2")
	// 	if err == nil ||
	// 		!strings.Contains(err.Error(), "mismatching active affiliate id") {
	// 		t.Fatalf("should be failed on mismatching prev code: %v", err)
	// 	}

	// 	if err := deactivateAffiliateCode(User(3), "code1"); err != nil {
	// 		t.Fatalf("failed to deactivate affiliate code: %v", err)
	// 	}

	// 	activeCode, err := getActiveAffiliateCode(User(3))
	// 	if err != nil {
	// 		t.Fatalf("failed to get active affiliate code: %v", err)
	// 	}
	// 	if activeCode != nil {
	// 		t.Fatalf("failed to get correct active affiliate code: %v", *activeCode)
	// 	}

	// 	codes, err := getOwnedAffiliateCode(User(1))
	// 	if err != nil {
	// 		t.Fatalf("failed to get owned affiliate codes meta: %v", err)
	// 	}

	// 	if len(codes) != 2 ||
	// 		codes[0].Code != "code1" ||
	// 		codes[0].Reward != 0 ||
	// 		codes[0].TotalEarned != 0 ||
	// 		codes[0].UserCnt != 1 ||
	// 		codes[1].Code != "code2" ||
	// 		codes[1].Reward != 0 ||
	// 		codes[1].TotalEarned != 0 ||
	// 		codes[1].UserCnt != 0 {
	// 		t.Fatalf("failed to create affiliate codes properly: %v", codes)
	// 	}
	// }
}

func TestDistributeAffiliate(t *testing.T) {
	TestActivateAffiliateCode(t)

	distributed, err := distributeAffiliate(User(3), 1000, 20000)
	if err != nil {
		t.Fatalf("failed to distribute affiliate: %v", err)
	}
	if distributed != 25 {
		t.Fatalf("failed to distribute proper amount: %d", distributed)
	}

	codes, err := getOwnedAffiliateCode(User(1))
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes meta: %v", err)
	}

	if len(codes) != 2 ||
		codes[0].Code != "code1" ||
		codes[0].Reward != 25 ||
		codes[0].TotalEarned != 25 ||
		codes[0].UserCnt != 2 ||
		codes[0].TotalWagered != 20000 ||
		codes[1].Code != "code2" ||
		codes[1].Reward != 0 ||
		codes[1].TotalEarned != 0 ||
		codes[1].UserCnt != 0 ||
		codes[1].TotalWagered != 0 {
		t.Fatalf("failed to create affiliate codes properly: %v", codes)
	}
}

func TestClaimAffiliateRewards(t *testing.T) {
	TestDistributeAffiliate(t)

	{
		claimed, err := claimAffiliateRewards(User(1), []string{
			"code2",
		})
		if err != nil {
			t.Fatalf("failed to claim affiliate rewards: %v", err)
		}
		if claimed != 0 {
			t.Fatalf("failed to claim proper amount: %d", claimed)
		}

		codes, err := getOwnedAffiliateCode(User(1))
		if err != nil {
			t.Fatalf("failed to get owned affiliate codes meta: %v", err)
		}

		if len(codes) != 2 ||
			codes[0].Code != "code1" ||
			codes[0].Reward != 25 ||
			codes[0].TotalEarned != 25 ||
			codes[0].UserCnt != 2 ||
			codes[1].Code != "code2" ||
			codes[1].Reward != 0 ||
			codes[1].TotalEarned != 0 ||
			codes[1].UserCnt != 0 {
			t.Fatalf("failed to create affiliate codes properly: %v", codes)
		}
	}

	{
		claimed, err := claimAffiliateRewards(User(1), []string{
			"code2",
			"code1",
		})
		if err != nil {
			t.Fatalf("failed to claim affiliate rewards: %v", err)
		}
		if claimed != 25 {
			t.Fatalf("failed to claim proper amount: %d", claimed)
		}

		codes, err := getOwnedAffiliateCode(User(1))
		if err != nil {
			t.Fatalf("failed to get owned affiliate codes meta: %v", err)
		}

		if len(codes) != 2 ||
			codes[0].Code != "code1" ||
			codes[0].Reward != 0 ||
			codes[0].TotalEarned != 25 ||
			codes[0].UserCnt != 2 ||
			codes[1].Code != "code2" ||
			codes[1].Reward != 0 ||
			codes[1].TotalEarned != 0 ||
			codes[1].UserCnt != 0 {
			t.Fatalf("failed to create affiliate codes properly: %v", codes)
		}

		user1 := models.User{}
		db := tests.InitMockDB(false, false)
		if result := db.Preload(
			"Wallet.Balance.ChipBalance",
		).First(&user1, 1); result.Error != nil {
			t.Fatalf("failed to retrieve first user: %v", result.Error)
		}
		if user1.Wallet.Balance.ChipBalance.Balance != 125 {
			t.Fatalf(
				"failed to claim affiliate properly: %d",
				user1.Wallet.Balance.ChipBalance.Balance,
			)
		}
	}
}

func TestDeleteAffiliateCode(t *testing.T) {
	TestDistributeAffiliate(t)

	claimed, err := deleteAffiliateCode(User(1), []string{"code1", "code2"})
	if err != nil {
		t.Fatalf("failed to delete affiliate codes: %v", err)
	}
	if claimed != 25 {
		t.Fatalf("failed to claim proper amount affiliate rewards on delete: %d", claimed)
	}

	codes, err := getOwnedAffiliateCode(User(1))
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes meta: %v", err)
	}

	if len(codes) != 0 {
		t.Fatalf("failed to delete affiliate codes properly: %v", codes)
	}
}

func TestAffiliateRate(t *testing.T) {
	if getAffiliateRate(2) != 5 {
		t.Fatalf(
			"failed to get affiliate rate properly, expected: %d, actual: %d",
			getAffiliateRate(2), 5,
		)
	}
	if getAffiliateRate(5) != 5 {
		t.Fatalf(
			"failed to get affiliate rate properly, expected: %d, actual: %d",
			getAffiliateRate(5), 5,
		)
	}
	if getAffiliateRate(7) != 7 {
		t.Fatalf(
			"failed to get affiliate rate properly, expected: %d, actual: %d",
			getAffiliateRate(7), 7,
		)
	}
	if getAffiliateRate(20) != 20 {
		t.Fatalf(
			"failed to get affiliate rate properly, expected: %d, actual: %d",
			getAffiliateRate(20), 20,
		)
	}
	if getAffiliateRate(30) != 20 {
		t.Fatalf(
			"failed to get affiliate rate properly, expected: %d, actual: %d",
			getAffiliateRate(30), 20,
		)
	}
}

func TestSetAffiliateCustomRate(t *testing.T) {
	TestCreateAffiliateCode(t)

	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	if err := SetAffiliateCustomRate("code1", 20, sessionId); err != nil {
		t.Fatalf("failed to set affiliate custom rate: %v", err)
	}
	if err := setAffiliateCustomRate("code", 10, sessionId); err == nil ||
		!strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
		t.Fatalf("should be failed on not found code")
	}

	if _, err := activateAffiliateCode(User(2), "code1", sessionId); err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}

	activeAffiliate, err := getActiveAffiliateCode(User(2), sessionId)
	if err != nil {
		t.Fatalf("failed to get active affiliate for second user: %v", err)
	}
	if activeAffiliate.Rate != 20 {
		t.Fatalf("failed to set custom affiliate rate properly: %d", activeAffiliate.Rate)
	}
	if err := commitSession(sessionId); err != nil {
		t.Fatalf("failed to commit session: %v", err)
	}
}

func TestDistributeAffiliateRewardWithCustomRate(t *testing.T) {
	TestSetAffiliateCustomRate(t)

	distributed, err := distributeAffiliate(User(2), 1000, 20000)
	if err != nil {
		t.Fatalf("failed to distribute affiliate: %v", err)
	}
	if distributed != 100 {
		t.Fatalf("failed to distribute proper amount: %d", distributed)
	}

	codes, err := getOwnedAffiliateCode(User(1))
	if err != nil {
		t.Fatalf("failed to get owned affiliate codes meta: %v", err)
	}

	if len(codes) != 2 ||
		codes[0].Code != "code1" ||
		codes[0].Reward != 100 ||
		codes[0].TotalEarned != 100 ||
		codes[0].UserCnt != 1 ||
		codes[0].TotalWagered != 20000 ||
		codes[1].Code != "code2" ||
		codes[1].Reward != 0 ||
		codes[1].TotalEarned != 0 ||
		codes[1].UserCnt != 0 ||
		codes[1].TotalWagered != 0 {
		t.Fatalf("failed to create affiliate codes properly: %v", codes)
	}
}

func TestAffiliateLifetime(t *testing.T) {
	TestCreateAffiliateCode(t)

	if _, err := activateAffiliateCode(
		User(2),
		"code1",
	); err != nil {
		t.Fatalf("failed to activate affiliate code: %v", err)
	}

	distributed, err := distributeAffiliate(
		User(2),
		1000,
		10000,
	)
	if err != nil {
		t.Fatalf("failed to distribute affiliate rewards: %v", err)
	}

	if result, err := getAffiliateDetail("code1"); err != nil {
		t.Fatalf("failed to get affiliate detail: %v", err)
	} else if result.Code != "code1" {
		t.Fatal("affiliate detail code mismatching")
	} else if len(result.Users) != 1 {
		t.Fatal("failed to get activated affiliate user properly")
	} else if result.Users[0].Name != "taka2" ||
		result.Users[0].Wagered != 10000 ||
		result.Users[0].Reward != distributed {
		t.Fatalf(
			"failed to get user detail properly: %d, %d",
			result.Users[0].Wagered,
			result.Users[0].Reward,
		)
	}

	time.Sleep(time.Second * 3)

	deactivateAffiliateCode(User(2), "code1")
	if result, err := getAffiliateDetail("code1"); err != nil {
		t.Fatal("failed")
	} else if len(result.Users) != 0 {
		t.Fatal("failed")
	}

	activateAffiliateCode(User(2), "code1")
	if result, err := getAffiliateDetail("code1"); err != nil {
		t.Fatal("failed")
	} else if len(result.Users) != 1 ||
		result.Users[0].Lifetime < 3 ||
		result.Users[0].Lifetime > 4 ||
		result.Users[0].Wagered != 10000 ||
		result.Users[0].Reward != distributed {
		t.Fatal("failed")
	}
}

func TestUniqueConstraintDelete(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to get mock db")
	}
	initialize(db)

	if err := db.Create(
		&models.Affiliate{
			Code: "code",
			Creator: models.User{
				Name: "taka",
			},
		},
	).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Where(
		"code = ?",
		"code",
	).Delete(
		&models.Affiliate{},
	).Error; err != nil {
		t.Fatal(err)
	}

	if err := db.Create(
		&models.Affiliate{
			Code:      "code",
			CreatorID: 1,
		},
	).Error; err != nil {
		t.Fatal(err)
	}
}
