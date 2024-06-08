package seed

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
)

func initBaseMockDB() (*gorm.DB, error) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		return nil, utils.MakeError(
			"seed_test",
			"initBaseMockDB",
			"failed to initialize mocked db",
			errors.New("retrieved db pointer is nil"),
		)
	}

	if err := db_aggregator.Initialize(db); err != nil {
		return nil, utils.MakeError(
			"seed_test",
			"initBaseMockDB",
			"failed to initialize session",
			err,
		)
	}

	mockUsers := []models.User{
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
			Name:          "User3",
			WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
			Role:          models.UserRole,
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: 100,
					},
				},
			},
			Banned: true,
		},
	}
	if result := db.Create(&mockUsers); result.Error != nil {
		return nil, utils.MakeError(
			"seed_test",
			"initBaseMockDB",
			"failed to create mock users",
			result.Error,
		)
	}

	return db, nil
}

func TestInitUserSeedPair(t *testing.T) {
	if _, err := initBaseMockDB(); err != nil {
		t.Fatalf("failed to initialize base mock db: %v", err)
	}

	newPair, err := initUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to init user seed pair: %v", err)
	}
	if len(newPair.ClientSeed.Seed) == 0 ||
		len(newPair.ServerSeed.Seed) == 0 ||
		len(newPair.ServerSeed.Hash) == 0 ||
		newPair.Nonce != 0 ||
		newPair.UserID != 1 ||
		newPair.IsExpired != false ||
		newPair.UsingCount != 0 {
		t.Fatalf("failed to initiate user seed pair properly: %v", newPair)
	}

	if _, err := initUserSeedPair(db_aggregator.User(1)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeAlreadyExistingPair) {
		t.Fatalf("should be failed for duplicated code creation: %v", err)
	}

	if _, err := initUserSeedPair(db_aggregator.User(4)); err == nil {
		t.Fatalf("should be failed for non user code creation: %v", err)
	}

	if _, err := initUserSeedPair(db_aggregator.User(0)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed for invalid parameter: %v", err)
	}

	initUserSeedPair(db_aggregator.User(3))
}

func TestLockAndRetrieveUserSeedPair(t *testing.T) {
	TestInitUserSeedPair(t)

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}

	go func() {
		sessionId, err := db_aggregator.StartSession()
		if err != nil {
			fmt.Printf("failed to start session: %v\n\r", err)
		}

		fmt.Println("inside one ======")
		if _, err := lockAndRetrieveUserActiveSeedPair(
			db_aggregator.User(1),
			sessionId,
		); err != nil {
			fmt.Printf("failed to lock and retrieve seed pair: %v", err)
		}
		time.Sleep(time.Second * 5)

		if err := db_aggregator.CommitSession(sessionId); err != nil {
			fmt.Printf("failed to commit session: %v", err)
		}
	}()

	time.Sleep(time.Second)
	startTime := time.Now()

	fmt.Println("outside one ======")
	seedPair, err := lockAndRetrieveUserActiveSeedPair(
		db_aggregator.User(1),
		sessionId,
	)
	if err != nil {
		t.Fatalf("failed to lock and retrieve seed pair: %v", err)
	}

	endTime := time.Now()

	if endTime.Sub(startTime) < time.Second*4 {
		t.Fatalf(
			"failed to lock properly: left: %d, right: %d",
			endTime.Sub(startTime),
			time.Second*5,
		)
	}

	if len(seedPair.ClientSeed.Seed) == 0 ||
		len(seedPair.ServerSeed.Seed) == 0 ||
		len(seedPair.ServerSeed.Hash) == 0 ||
		seedPair.Nonce != 0 ||
		seedPair.UserID != 1 ||
		seedPair.IsExpired != false ||
		seedPair.UsingCount != 0 {
		t.Fatalf("failed to initiate user seed pair properly: %v", seedPair)
	}
}

func TestGetUserInfoChecked(t *testing.T) {
	db, err := initBaseMockDB()
	if err != nil {
		t.Fatalf("failed to init base mock db: %v", err)
	}

	userInfo, err := getUserInfoChecked(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get user info: %v", err)
	}
	if userInfo.ID != 1 ||
		userInfo.Name != "User1" {
		t.Fatalf("failed to retrieve user info properly: %v", userInfo)
	}

	userInfo.Banned = true
	if result := db.Save(&userInfo); result.Error != nil {
		t.Fatalf("failed to update user as banned: %v", result.Error)
	}

	if _, err := getUserInfoChecked(db_aggregator.User(1)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeBannedUser) {
		t.Fatalf("should be failed on banned user retrieve: %v", err)
	}

	if _, err := getUserInfoChecked(db_aggregator.User(0)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on invalid parameter: %v", err)
	}

	if _, err := getUserInfoChecked(db_aggregator.User(3)); err == nil {
		t.Fatalf("should be failed on not found user: %v", err)
	}
}

func TestExpireSeedPairUnchecked(t *testing.T) {
	TestInitUserSeedPair(t)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}

	clientSeed, err := utils.GenerateClientSeed(24)
	if err != nil {
		t.Fatalf("failed to generate client seed: %v", err)
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start a new session: %v", err)
	}

	newSeedPair, err := expireSeedPairUnchecked(
		seedPair,
		clientSeed,
		sessionId,
	)
	if err != nil {
		t.Fatalf("failed to expire seedpair: %v", err)
	}

	if newSeedPair == nil ||
		newSeedPair.IsExpired != false ||
		newSeedPair.UserID != 1 ||
		len(newSeedPair.ClientSeed.Seed) == 0 ||
		newSeedPair.ClientSeed.Seed == seedPair.ClientSeed.Seed ||
		newSeedPair.Nonce != 0 ||
		seedPair.IsExpired != true ||
		len(newSeedPair.NextServerSeed.Seed) == 0 ||
		seedPair.ID != 1 ||
		seedPair.NextServerSeedID != newSeedPair.ServerSeedID ||
		newSeedPair.ClientSeed.Seed != clientSeed {
		t.Fatalf(
			"failed to create new seed pair properly: \n\rold:%v\n\rnew:%v",
			seedPair,
			newSeedPair,
		)
	}

	if _, err := expireSeedPairUnchecked(nil, clientSeed, sessionId); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on invalid seed pair parameter: %v", err)
	}
	if _, err := expireSeedPairUnchecked(&models.SeedPair{}, clientSeed, sessionId); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on zero id seed pair parameter: %v", err)
	}
	if _, err := expireSeedPairUnchecked(newSeedPair, "", sessionId); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on empty string client seed parameter: %v", err)
	}
	if _, err := expireSeedPairUnchecked(seedPair, clientSeed, sessionId); err == nil ||
		!utils.IsErrorCode(err, ErrCodeExpiredPair) {
		t.Fatalf("should be failed to try expire already expired pair: %v", err)
	}

	db_aggregator.CommitSession(sessionId)
}

func TestGetActiveUserSeedPair(t *testing.T) {
	TestInitUserSeedPair(t)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}
	if len(seedPair.ClientSeed.Seed) == 0 ||
		len(seedPair.ServerSeed.Seed) == 0 ||
		len(seedPair.ServerSeed.Hash) == 0 ||
		seedPair.Nonce != 0 ||
		seedPair.UserID != 1 ||
		seedPair.IsExpired != false ||
		seedPair.UsingCount != 0 {
		t.Fatalf("failed to get user seed pair properly: %v", seedPair)
	}

	if _, err := getActiveUserSeedPair(db_aggregator.User(2)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotFoundUnexpiredPair) {
		t.Fatalf("should be failed on not found unexpired pair: %v", err)
	}

	db := tests.InitMockDB(false, false)
	if db == nil {
		t.Fatal("retrieved mock db pointer is nil")
	}

	if result := db.Create(&models.SeedPair{
		UserID:         1,
		ClientSeed:     models.ClientSeed{Seed: "client-seed"},
		ServerSeed:     models.ServerSeed{Seed: "server-seed"},
		NextServerSeed: models.ServerSeed{Seed: "next-server-seed"},
	}); result.Error != nil {
		t.Fatalf("failed to create duplicated seed pair: %v", result.Error)
	}

	if _, err := getActiveUserSeedPair(db_aggregator.User(1)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeMultipleUnexpiredPairs) {
		t.Fatalf("should be failed on multiple unexpired pairs retrivial: %v", err)
	}
}

func TestGetExpiredUserSeedPairs(t *testing.T) {
	TestInitUserSeedPair(t)

	firstSeed, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get activer user seed pair: %v", err)
	}

	serverSeeds := []string{firstSeed.ServerSeed.Seed}
	clientSeeds := []string{
		firstSeed.ClientSeed.Seed,
		"client-seed-0", "client-seed-1",
		"client-seed-2", "client-seed-3",
		"client-seed-4", "client-seed-5",
		"client-seed-6", "client-seed-7",
		"client-seed-8", "client-seed-9",
	}

	prevPair := firstSeed
	for i := 1; i <= 10; i++ {
		prevPair, err = rotateUserSeedPair(
			db_aggregator.User(1),
			prevPair.ServerSeed.Hash,
			clientSeeds[i],
		)
		if err != nil {
			t.Fatalf("failed to rotate user seed pair: %v, %d", err, i)
		}
		serverSeeds = append(
			serverSeeds,
			prevPair.ServerSeed.Seed,
		)
	}

	{
		startIndex := 0
		endIndex := 3
		limit := 3
		seedPairs, err := getExpiredUserSeedPairs(
			db_aggregator.User(1),
			uint(startIndex),
			uint(limit),
		)
		if err != nil {
			t.Fatalf(
				"failed to get expired user seed pairs, skip: %d, limit: %d, error: %v",
				startIndex, limit, err,
			)
		}

		expected := []models.SeedPair{}
		for i := startIndex; i < endIndex; i++ {
			expected = append(
				expected,
				models.SeedPair{
					UserID:    1,
					IsExpired: true,
					ServerSeed: models.ServerSeed{
						Seed: serverSeeds[9-i],
					},
					ClientSeed: models.ClientSeed{
						Seed: clientSeeds[9-i],
					},
				},
			)
		}

		if len(seedPairs) != len(expected) {
			t.Fatalf(
				"failed to retrieve exact number of seed pairs. expected: %d, retrieved: %d",
				len(expected),
				len(seedPairs),
			)
		}

		for i, seedPair := range seedPairs {
			if seedPair.ClientSeed.Seed != expected[i].ClientSeed.Seed ||
				seedPair.ServerSeed.Seed != expected[i].ServerSeed.Seed ||
				seedPair.IsExpired != expected[i].IsExpired ||
				seedPair.UserID != expected[i].UserID ||
				seedPair.IsExpired != true {
				t.Fatalf(
					"failed to retrieve expected seed pair properly. index: %d, \n\rexpected: %v, \n\ractual: %v",
					i,
					expected[i].ClientSeed,
					seedPair.ClientSeed,
				)
			}
		}
	}

	{
		startIndex := 3
		endIndex := 7
		limit := 4
		seedPairs, err := getExpiredUserSeedPairs(
			db_aggregator.User(1),
			uint(startIndex),
			uint(limit),
		)
		if err != nil {
			t.Fatalf(
				"failed to get expired user seed pairs, skip: %d, limit: %d, error: %v",
				startIndex, endIndex, err,
			)
		}

		expected := []models.SeedPair{}
		for i := startIndex; i < endIndex; i++ {
			expected = append(
				expected,
				models.SeedPair{
					UserID:    1,
					IsExpired: true,
					ServerSeed: models.ServerSeed{
						Seed: serverSeeds[9-i],
					},
					ClientSeed: models.ClientSeed{
						Seed: clientSeeds[9-i],
					},
				},
			)
		}

		if len(seedPairs) != len(expected) {
			t.Fatalf(
				"failed to retrieve exact number of seed pairs. expected: %d, retrieved: %d",
				len(expected),
				len(seedPairs),
			)
		}

		for i, seedPair := range seedPairs {
			if seedPair.ClientSeed.Seed != expected[i].ClientSeed.Seed ||
				seedPair.ServerSeed.Seed != expected[i].ServerSeed.Seed ||
				seedPair.IsExpired != expected[i].IsExpired ||
				seedPair.UserID != expected[i].UserID ||
				seedPair.IsExpired != true {
				t.Fatalf(
					"failed to retrieve expected seed pair properly. expected: %v, actual: %v",
					expected,
					seedPair,
				)
			}
		}
	}

	{
		startIndex := 5
		endIndex := 10
		limit := 10
		seedPairs, err := getExpiredUserSeedPairs(
			db_aggregator.User(1),
			uint(startIndex),
			uint(limit),
		)
		if err != nil {
			t.Fatalf(
				"failed to get expired user seed pairs, skip: %d, limit: %d, error: %v",
				startIndex, endIndex, err,
			)
		}

		expected := []models.SeedPair{}
		for i := startIndex; i < endIndex; i++ {
			expected = append(
				expected,
				models.SeedPair{
					UserID:    1,
					IsExpired: true,
					ServerSeed: models.ServerSeed{
						Seed: serverSeeds[9-i],
					},
					ClientSeed: models.ClientSeed{
						Seed: clientSeeds[9-i],
					},
				},
			)
		}

		if len(seedPairs) != len(expected) {
			t.Fatalf(
				"failed to retrieve exact number of seed pairs. expected: %d, retrieved: %d",
				len(expected),
				len(seedPairs),
			)
		}

		for i, seedPair := range seedPairs {
			if seedPair.ClientSeed.Seed != expected[i].ClientSeed.Seed ||
				seedPair.ServerSeed.Seed != expected[i].ServerSeed.Seed ||
				seedPair.IsExpired != expected[i].IsExpired ||
				seedPair.UserID != expected[i].UserID ||
				seedPair.IsExpired != true {
				t.Fatalf(
					"failed to retrieve expected seed pair properly. expected: %v, actual: %v",
					expected,
					seedPair,
				)
			}
		}
	}

	{
		startIndex := 10
		endIndex := 10
		limit := 5
		seedPairs, err := getExpiredUserSeedPairs(
			db_aggregator.User(1),
			uint(startIndex),
			uint(limit),
		)
		if err != nil {
			t.Fatalf(
				"failed to get expired user seed pairs, skip: %d, limit: %d, error: %v",
				startIndex, endIndex, err,
			)
		}

		expected := []models.SeedPair{}
		for i := startIndex; i < endIndex; i++ {
			expected = append(
				expected,
				models.SeedPair{
					UserID:    1,
					IsExpired: true,
					ServerSeed: models.ServerSeed{
						Seed: serverSeeds[9-i],
					},
					ClientSeed: models.ClientSeed{
						Seed: clientSeeds[9-i],
					},
				},
			)
		}

		if len(seedPairs) != len(expected) {
			t.Fatalf(
				"failed to retrieve exact number of seed pairs. expected: %d, retrieved: %d",
				len(expected),
				len(seedPairs),
			)
		}

		for i, seedPair := range seedPairs {
			if seedPair.ClientSeed.Seed != expected[i].ClientSeed.Seed ||
				seedPair.ServerSeed.Seed != expected[i].ServerSeed.Seed ||
				seedPair.IsExpired != expected[i].IsExpired ||
				seedPair.UserID != expected[i].UserID ||
				seedPair.IsExpired != true {
				t.Fatalf(
					"failed to retrieve expected seed pair properly. expected: %v, actual: %v",
					expected,
					seedPair,
				)
			}
		}
	}

	{
		startIndex := 20
		endIndex := 10
		limit := 5
		seedPairs, err := getExpiredUserSeedPairs(
			db_aggregator.User(1),
			uint(startIndex),
			uint(limit),
		)
		if err != nil {
			t.Fatalf(
				"failed to get expired user seed pairs, skip: %d, limit: %d, error: %v",
				startIndex, endIndex, err,
			)
		}

		expected := []models.SeedPair{}
		for i := startIndex; i < endIndex; i++ {
			expected = append(
				expected,
				models.SeedPair{
					UserID:    1,
					IsExpired: true,
					ServerSeed: models.ServerSeed{
						Seed: serverSeeds[9-i],
					},
					ClientSeed: models.ClientSeed{
						Seed: clientSeeds[9-i],
					},
				},
			)
		}

		if len(seedPairs) != len(expected) {
			t.Fatalf(
				"failed to retrieve exact number of seed pairs. expected: %d, retrieved: %d",
				len(expected),
				len(seedPairs),
			)
		}

		for i, seedPair := range seedPairs {
			if seedPair.ClientSeed.Seed != expected[i].ClientSeed.Seed ||
				seedPair.ServerSeed.Seed != expected[i].ServerSeed.Seed ||
				seedPair.IsExpired != expected[i].IsExpired ||
				seedPair.UserID != expected[i].UserID ||
				seedPair.IsExpired != true {
				t.Fatalf(
					"failed to retrieve expected seed pair properly. expected: %v, actual: %v",
					expected,
					seedPair,
				)
			}
		}
	}
}

// func TestCanBeExpired(t *testing.T) {
// 	canBeExpired1 := models.SeedPair{
// 		Model: gorm.Model{
// 			ID: 1,
// 		},
// 		UserID: 1,
// 	}
// 	canNotBeExpired1 := models.SeedPair{
// 		Model: gorm.Model{
// 			ID: 0,
// 		},
// 		UserID: 2,
// 	}
// 	canNotBeExpired2 := models.SeedPair{
// 		Model: gorm.Model{
// 			ID: 1,
// 		},
// 		IsExpired: true,
// 	}
// 	canNotBeExpired3 := models.SeedPair{
// 		Model: gorm.Model{
// 			ID: 1,
// 		},
// 		UsingCount: 3,
// 	}

// 	if !canBeExpired(&canBeExpired1) {
// 		t.Fatalf("should can be expired: %v", canBeExpired1)
// 	}
// 	if canBeExpired(&canNotBeExpired1) {
// 		t.Fatalf("should can not be expired: %v", canNotBeExpired1)
// 	}
// 	if canBeExpired(&canNotBeExpired2) {
// 		t.Fatalf("should can not be expired: %v", canNotBeExpired2)
// 	}
// 	if canBeExpired(&canNotBeExpired3) {
// 		t.Fatalf("should can not be expired: %v", canNotBeExpired3)
// 	}
// }

func TestGetUnhashedServerSeed(t *testing.T) {
	TestExpireSeedPairUnchecked(t)

	seedPairs, err := getExpiredUserSeedPairs(
		db_aggregator.User(1),
		0, 1,
	)
	if err != nil {
		t.Fatalf("failed to retrieve expired user seed pairs: %v", err)
	}
	if len(seedPairs) != 1 {
		t.Fatalf("failed to retrieve expired user seed pairs properly: %v", seedPairs)
	}

	unhashedServerSeed, err := getUnhashedServerSeed(seedPairs[0].ServerSeed.Hash)
	if err != nil {
		t.Fatalf("failed to get unhashed server seed: %v", err)
	}
	if unhashedServerSeed != seedPairs[0].ServerSeed.Seed {
		t.Fatalf(
			"failed to get unhashed server seed properly: expected: %s, actual: %s",
			seedPairs[0].ServerSeed.Seed, unhashedServerSeed,
		)
	}

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}
	if _, err := getUnhashedServerSeed(seedPair.ServerSeed.Hash); err == nil ||
		!utils.IsErrorCode(err, ErrCodeUnexpiredPair) {
		t.Fatalf("should be failed on trying to hash unexpired server seed: %v", err)
	}

	if _, err := getUnhashedServerSeed("not-found"); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotFoundServerSeed) {
		t.Fatalf("should be failed on not found server seed: %v", err)
	}
	if _, err := getUnhashedServerSeed(""); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on invalid parameter with empty string: %v", err)
	}
}

func TestRotateUserSeedPair(t *testing.T) {
	TestInitUserSeedPair(t)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}
	if seedPair.IsExpired != false {
		t.Fatalf("already expired seed retrieved: %v", seedPair)
	}

	newSeedPair, err := rotateUserSeedPair(
		db_aggregator.User(1),
		seedPair.ServerSeed.Hash,
		"client-seed",
	)
	if err != nil {
		t.Fatalf("failed to rotate user seed pair: %v", err)
	}
	if seedPair.IsExpired != false ||
		seedPair.ID == newSeedPair.ID ||
		newSeedPair.ClientSeed.Seed != "client-seed" ||
		newSeedPair.UserID != 1 ||
		seedPair.NextServerSeed.Seed != newSeedPair.ServerSeed.Seed {
		t.Fatalf(
			"failed to rotate seed pairs properly: new: %v, old: %v",
			newSeedPair, seedPair,
		)
	}

	if _, err := rotateUserSeedPair(db_aggregator.User(2), "server-seed-hash", "client-seed"); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotFoundUnexpiredPair) {
		t.Fatalf("should be failed on not found unexpired pair: %v", err)
	}

	borrowedPair, err := borrowUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to borrow user seed pair: %v", err)
	}
	if _, err := rotateUserSeedPair(
		db_aggregator.User(1),
		borrowedPair.ServerSeed.Hash,
		"client-seed",
	); err == nil ||
		!utils.IsErrorCode(err, ErrCodePairsStillInUse) {
		t.Fatalf("should be failed on rotating seed being in use: %v", err)
	}

	if err := returnUserSeedPair(db_aggregator.User(1), borrowedPair.ID); err != nil {
		t.Fatalf("failed to return user seed pair: %v", err)
	}
	prevSeed, err := rotateUserSeedPair(
		db_aggregator.User(1),
		borrowedPair.ServerSeed.Hash,
		"client-seed",
	)
	if err != nil {
		t.Fatalf("failed to rotate user seed pair after returning: %v", err)
	}

	remainingCnt := SeedPairRotateLimitPerMinute - 2
	for i := uint(0); i < remainingCnt; i++ {
		prevSeed, err = rotateUserSeedPair(
			db_aggregator.User(1),
			prevSeed.ServerSeed.Hash,
			"client-seed",
		)
		if err != nil {
			t.Fatalf("failed to rotate user seed pair in loop: %d, %v", i, err)
		}
	}

	if _, err := rotateUserSeedPair(
		db_aggregator.User(1),
		prevSeed.ServerSeed.Hash,
		"client-seed"); err == nil ||
		!utils.IsErrorCode(err, ErrCodeTooManyRotates) {
		t.Fatalf("should be failed on too many rotates: %v", err)
	}
}

func TestBorrowUserSeedPair(t *testing.T) {
	TestInitUserSeedPair(t)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}

	borrowedPair, err := borrowUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to borrow user seed pair: %v", err)
	}

	if seedPair.ID != borrowedPair.ID ||
		seedPair.IsExpired != borrowedPair.IsExpired ||
		seedPair.IsExpired != false ||
		seedPair.Nonce+1 != borrowedPair.Nonce ||
		seedPair.UsingCount+1 != borrowedPair.UsingCount ||
		borrowedPair.Nonce != 1 ||
		borrowedPair.UsingCount != 1 {
		t.Fatalf(
			"failed to borrow user seed properly: prev: %v, next: %v",
			seedPair, borrowedPair,
		)
	}

	if _, err := borrowUserSeedPair(db_aggregator.User(2)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotFoundUnexpiredPair) {
		t.Fatalf("should be failed on not found unexpired pair: %v", err)
	}
	if _, err := borrowUserSeedPair(db_aggregator.User(3)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeBannedUser) {
		t.Fatalf("should be failed to be brought by banned user: %v", err)
	}

	remainingCnt := SeedPairBorrowLimit - 1
	for i := uint(0); i < remainingCnt; i++ {
		if _, err := borrowUserSeedPair(db_aggregator.User(1)); err != nil {
			t.Fatalf("failed to borrow user seed pair in loop: %d, %v", i, err)
		}
	}

	if _, err := borrowUserSeedPair(db_aggregator.User(1)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeTooManyBorrows) {
		t.Fatalf("should be failed on too many borrows: %v", err)
	}
	if _, err := borrowUserSeedPair(db_aggregator.User(0)); err == nil ||
		!utils.IsErrorCode(err, ErrCodeInvalidParameter) {
		t.Fatalf("should be failed on invalid parameter: %v", err)
	}
}

func TestReturnUserSeedPair(t *testing.T) {
	TestBorrowUserSeedPair(t)

	for i := uint(0); i < SeedPairBorrowLimit; i++ {
		seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
		if err != nil {
			t.Fatalf("failed to get active user seed pair: %v", err)
		}
		if seedPair.UsingCount != SeedPairBorrowLimit-i {
			t.Fatalf(
				"failed to get seed pair using count properly: expected: %d, actual: %d",
				SeedPairBorrowLimit-i,
				seedPair.UsingCount,
			)
		}

		if err := returnUserSeedPair(db_aggregator.User(1), seedPair.ID); err != nil {
			t.Fatalf("failed to return user seed pair: %v", err)
		}
	}

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}
	if err := returnUserSeedPair(db_aggregator.User(1), seedPair.ID); err == nil ||
		!utils.IsErrorCode(err, ErrCodePairsNotInUse) {
		t.Fatalf("should be failed on returning not used pairs: %v", err)
	}
	if seedPair.Nonce != 10 ||
		seedPair.UsingCount != 0 {
		t.Fatalf(
			"invalid nonce and using count: nonce; %d, usingCount; %d",
			seedPair.Nonce,
			seedPair.UsingCount,
		)
	}
	if err := returnUserSeedPair(db_aggregator.User(1), seedPair.ID+1); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotOwnedSeedPairID) {
		t.Fatalf("should be failed on returning not owned pair id: %v", err)
	}
	if err := returnUserSeedPair(db_aggregator.User(2), seedPair.ID); err == nil ||
		!utils.IsErrorCode(err, ErrCodeNotFoundUnexpiredPair) {
		t.Fatalf("should be failed on not found unexpired seed pair: %v", err)
	}
}

func TestTooManyRotates(t *testing.T) {
	TestInitUserSeedPair(t)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(1))
	if err != nil {
		t.Fatalf("failed to get active user seed pair: %v", err)
	}

	_, err = rotateUserSeedPair(
		db_aggregator.User(1),
		seedPair.ServerSeed.Hash,
		"client-seed",
	)
	if err != nil {
		t.Fatalf("failed to rotate seed pair: %v", err)
	}

	if tooMany, err := checkTooManyRotates(
		db_aggregator.User(1),
		db_aggregator.MainSessionId(),
	); err != nil || tooMany {
		t.Fatalf("should not be too many rotates. err: %v", err)
	}

	TestRotateUserSeedPair(t)

	if tooMany, err := checkTooManyRotates(
		db_aggregator.User(1),
		db_aggregator.MainSessionId(),
	); err != nil || !tooMany {
		t.Fatalf("should be too many rotates. err: %v", err)
	}
}

func TestCheckTooManyBorrows(t *testing.T) {
	if checkTooManyBorrows(&models.SeedPair{
		UsingCount: SeedPairBorrowLimit - 1,
	}) {
		t.Fatalf("should not be too many borrows")
	}

	if !checkTooManyBorrows(&models.SeedPair{
		UsingCount: SeedPairBorrowLimit,
	}) {
		t.Fatalf("should be too many borrows")
	}
}

func rotateUserSeedPairHelper(user db_aggregator.User, clientSeed string) {
	seedPair, err := getActiveUserSeedPair(user)
	if err != nil {
		fmt.Printf("failed to get active user seed pair: %v", err)
	}
	rotateUserSeedPair(user, seedPair.ServerSeed.Hash, clientSeed)
}

func TestConcurrentBorrowRotateReturns1(t *testing.T) {
	TestInitUserSeedPair(t)

	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")

	if _, err := getActiveUserSeedPair(db_aggregator.User(1)); err != nil {
		t.Fatalf("failed to retrieve active user seed pair: %v", err)
	}

}

func TestConcurrentBorrowRotateReturns2(t *testing.T) {
	TestInitUserSeedPair(t)

	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)

	if _, err := getActiveUserSeedPair(db_aggregator.User(1)); err != nil {
		t.Fatalf("failed to retrieve active user seed pair: %v", err)
	}
}

func TestConcurrentBorrowRotateReturns3(t *testing.T) {
	TestInitUserSeedPair(t)

	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go borrowUserSeedPair(db_aggregator.User(1))
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go borrowUserSeedPair(db_aggregator.User(1))
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go returnUserSeedPair(db_aggregator.User(1), 1)
	go rotateUserSeedPairHelper(db_aggregator.User(1), "client-seed")
	go returnUserSeedPair(db_aggregator.User(1), 1)

	if _, err := getActiveUserSeedPair(db_aggregator.User(1)); err != nil {
		t.Fatalf("failed to retrieve active user seed pair: %v", err)
	}
}
