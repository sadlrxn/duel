package coupon

import (
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
)

func TestCreateCoupon(t *testing.T) {
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
	var userNames []string = []string{"User", "User2"}
	var limit int = 10

	createCouponRequest = CreateCouponRequest{
		Type:    models.CouponForSpecUsers,
		Balance: 100,
	}
	_, err := Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeAccessUserNamesMissingForCreate) {
		t.Fatalf("should return access user names missing error code: %v", err)
	}

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserLimit: &limit,
		Balance:         100,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeAccessUserNamesMissingForCreate) {
		t.Fatalf("should return access user names missing error code: %v", err)
	}

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserNames: &userNames,
		Balance:         100,
	}
	code, err := Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create coupon for specific users: %v", err)
	}

	var coupon models.Coupon
	if result := db.First(&coupon, code); result.Error != nil {
		t.Fatalf("failed to get created coupon from db: %v", result.Error)
	}
	if coupon.Type != models.CouponForSpecUsers ||
		len(coupon.AccessUserIDs) != len(userNames) ||
		int(coupon.AccessUserLimit) != len(userNames) ||
		coupon.BonusBalance != 100 {
		t.Fatalf("coupon code not created properly. %v", coupon)
	}

	userNames = append(userNames, "AAA")
	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForSpecUsers,
		AccessUserNames: &userNames,
		Balance:         100,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeAccessUserNameNotExistForCreate) {
		t.Fatalf("should return access user names not exist error: %v", err)
	}

	createCouponRequest = CreateCouponRequest{
		Type:    models.CouponForLimitUsers,
		Balance: 100,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeLimitUserCountMissingForCreate) {
		t.Fatalf("should return limit user count missing error: %v", err)
	}

	limit = 0
	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &limit,
		Balance:         100,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeLimitUserCountMissingForCreate) {
		t.Fatalf("should return limit user count missing error: %v", err)
	}

	limit = 10
	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &limit,
		Balance:         100,
	}
	code, err = Create(createCouponRequest)
	if err != nil {
		t.Fatalf("failed to create coupon for limited number of users: %v", err)
	}

	coupon = models.Coupon{}
	if result := db.First(&coupon, code); result.Error != nil {
		t.Fatalf("failed to get created coupon from db: %v", result.Error)
	}
	if coupon.Type != models.CouponForLimitUsers ||
		int(coupon.AccessUserLimit) != limit ||
		coupon.BonusBalance != 100 {
		t.Fatalf("coupon code not created properly. %v", coupon)
	}

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &limit,
		Balance:         0,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeInvalidBalanceForCreate) {
		t.Fatalf("should return invalid balance error: %v", err)
	}

	createCouponRequest = CreateCouponRequest{
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: &limit,
		Balance:         -100,
	}
	_, err = Create(createCouponRequest)
	if err == nil || !utils.IsErrorCode(err, ErrCodeInvalidBalanceForCreate) {
		t.Fatalf("should return invalid balance error: %v", err)
	}

}

func TestLockAndRetrieveCoupon(t *testing.T) {
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
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	code := uuid.New()
	coupon := models.Coupon{
		Code:            &code,
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: 10,
		BonusBalance:    100,
		ClaimedCoupons:  []models.ClaimedCoupon{},
	}
	if result := db.Create(&coupon); result.Error != nil {
		t.Fatalf("failed to create coupon: %v", result.Error)
	}

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
		if _, err := lockAndRetrieveCoupon(
			code, sessionId,
		); err != nil {
			fmt.Printf("failed to lock and retrieve coupon: %v", err)
		}
		time.Sleep(time.Second * 5)
		if err := db_aggregator.CommitSession(sessionId); err != nil {
			fmt.Printf("failed to commit session: %v", err)
		}
	}()

	time.Sleep(time.Second)
	startTime := time.Now()

	fmt.Println("outside one ======")
	_, err = lockAndRetrieveCoupon(code, sessionId)
	if err != nil {
		t.Fatalf("failed to lock and retrieve coupon: %v", err)
	}

	endTime := time.Now()

	if endTime.Sub(startTime) < time.Second*4 {
		t.Fatalf(
			"failed to lock properly: left: %d, right: %d",
			endTime.Sub(startTime),
			time.Second*5,
		)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	coupon.CreatedAt = time.UnixMilli(time.Now().UnixMilli() - time.Hour.Milliseconds()*24*15)
	if result := db.Save(&coupon); result.Error != nil {
		t.Fatalf("failed to update coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	_, err = lockAndRetrieveCoupon(code, sessionId)
	if err == nil || !utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		t.Fatalf("should return coupon code not found error: %v", err)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	coupon.CreatedAt = time.Now()
	if result := db.Save(&coupon); result.Error != nil {
		t.Fatalf("failed to update coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	_, err = lockAndRetrieveCoupon(uuid.New(), sessionId)
	if err == nil || !utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		t.Fatalf("should return coupon code not found error: %v", err)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	coupon.CreatedAt = time.Now()
	if result := db.Save(&coupon); result.Error != nil {
		t.Fatalf("failed to update coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	retrievedCoupon, err := lockAndRetrieveCoupon(code, sessionId)
	if err != nil {
		t.Fatalf("failed to lock and retrieve coupon: %v", err)
	}
	if retrievedCoupon.Type != models.CouponForLimitUsers ||
		retrievedCoupon.AccessUserLimit != 10 ||
		retrievedCoupon.BonusBalance != 100 {
		t.Fatalf("retrieved coupon is not valid: %v", retrievedCoupon)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}
}

func TestLockAndRetrieveActiveCoupon(t *testing.T) {
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
	}
	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	code := uuid.New()
	coupon := models.Coupon{
		Code:            &code,
		Type:            models.CouponForLimitUsers,
		AccessUserLimit: 10,
		BonusBalance:    100,
		ClaimedCoupons: []models.ClaimedCoupon{
			{
				ClaimedUserID: users[0].ID,
				Balance:       100,
			},
		},
	}
	if result := db.Create(&coupon); result.Error != nil {
		t.Fatalf("failed to create coupon: %v", result.Error)
	}

	var claimed_coupon models.ClaimedCoupon
	if result := db.Where("coupon_id = ?", coupon.Code).Where("claimed_user_id = ?", users[0].ID).First(&claimed_coupon); result.Error != nil {
		t.Fatalf("failed to get created claimed coupon: %v", result.Error)
	}

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
		if _, err := lockAndRetrieveActiveCoupon(
			users[0].ID, sessionId,
		); err != nil {
			fmt.Printf("failed to lock and retrieve coupon: %v", err)
		}
		time.Sleep(time.Second * 5)
		if err := db_aggregator.CommitSession(sessionId); err != nil {
			fmt.Printf("failed to commit session: %v", err)
		}
	}()

	time.Sleep(time.Second)
	startTime := time.Now()

	fmt.Println("outside one ======")
	_, err = lockAndRetrieveActiveCoupon(users[0].ID, sessionId)
	if err != nil {
		t.Fatalf("failed to lock and retrieve coupon: %v", err)
	}

	endTime := time.Now()

	if endTime.Sub(startTime) < time.Second*4 {
		t.Fatalf(
			"failed to lock properly: left: %d, right: %d",
			endTime.Sub(startTime),
			time.Second*5,
		)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	claimed_coupon.CreatedAt = time.UnixMilli(time.Now().UnixMilli() - time.Hour.Milliseconds()*9)
	if result := db.Save(&claimed_coupon); result.Error != nil {
		t.Fatalf("failed to update claimed coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	_, err = lockAndRetrieveActiveCoupon(users[0].ID, sessionId)
	if err == nil {
		t.Fatalf("should return coupon code not found error: %v", err)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	claimed_coupon.CreatedAt = time.Now()
	claimed_coupon.Exchanged = 100
	if result := db.Save(&claimed_coupon); result.Error != nil {
		t.Fatalf("failed to update claimed coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	_, err = lockAndRetrieveActiveCoupon(users[0].ID, sessionId)
	if err == nil {
		t.Fatalf("should return coupon code not found error: %v", err)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	claimed_coupon.CreatedAt = time.Now()
	claimed_coupon.Exchanged = 0
	if result := db.Save(&claimed_coupon); result.Error != nil {
		t.Fatalf("failed to update claimed coupon created_at: %v", result.Error)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	_, err = lockAndRetrieveActiveCoupon(2, sessionId)
	if err == nil {
		t.Fatalf("should return coupon code not found error: %v", err)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}

	sessionId, err = db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v\n\r", err)
	}
	retrieved_coupon, err := lockAndRetrieveActiveCoupon(users[0].ID, sessionId)
	if err != nil {
		t.Fatalf("failed to lock and retrieve claimed coupon: %v", err)
	}
	if retrieved_coupon.ClaimedUserID != users[0].ID ||
		retrieved_coupon.CouponID != *coupon.Code ||
		retrieved_coupon.Exchanged != 0 ||
		retrieved_coupon.Balance != 100 {
		t.Fatalf("retrieved claimed coupon info is not valid: %v", retrieved_coupon)
	}
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}
}
