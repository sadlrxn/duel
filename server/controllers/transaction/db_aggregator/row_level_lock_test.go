package db_aggregator

import (
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func commitSessionAfterTime(sessionId UUID, waitTime time.Duration) {
	time.Sleep(waitTime)

	if err := commitSession(sessionId); err != nil {
		fmt.Printf("failed to commit session: %v", err)
	}
}

func GetMockTxSession(t *testing.T) (UUID, *gorm.DB) {
	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start first session: %v", err)
	}
	session, err := getSession(sessionId)
	if err != nil {
		t.Fatalf("failed to get first session: %v", err)
	}

	return sessionId, session
}

func TestRowLevelLock(t *testing.T) {
	database := tests.InitMockDB(true, true)

	mockUser := models.User{
		Name:          "Mock User",
		WalletAddress: "EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB",
	}

	if result := database.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	if err := initialize(database); err != nil {
		t.Fatalf("failed to initialize transaction db")
	}

	sessionId1, session1 := GetMockTxSession(t)
	_, session2 := GetMockTxSession(t)

	user := models.User{}
	if result := session1.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user); result.Error != nil {
		t.Fatalf("failed to get & lock from first session: %v", result.Error)
	}

	if user.Name != mockUser.Name ||
		user.WalletAddress != mockUser.WalletAddress {
		t.Fatalf("retrieved user information mismatched")
	}

	var waitTime time.Duration = 5 * time.Second
	go commitSessionAfterTime(sessionId1, waitTime)

	t1 := time.Now()
	user2 := models.User{}
	session2.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).First(&user2)
	t2 := time.Now()

	if diff := t2.Sub(t1); diff < waitTime {
		t.Fatalf("retrieved user earlier than other session's commitment: %v", diff)
	}

	if user2.Name != mockUser.Name ||
		user2.WalletAddress != mockUser.WalletAddress {
		t.Fatalf("retrieved user information mismatched")
	}
}

func TestGetUserWalletRowLock(t *testing.T) {
	database := tests.InitMockDB(true, true)

	if err := initialize(database); err != nil {
		t.Fatalf("failed to initialize transaction db")
	}

	mockUsers := []models.User{
		{
			Name:          "Mock User1",
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
		},
		{
			Name:          "Mock User2",
			WalletAddress: "H79ykEibBW2R8UXbaNaRJEjzDqbLKKWQsg4WfLxScE45",
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
		},
	}

	if result := database.Create(&mockUsers); result.Error != nil {
		t.Fatalf("failed to create mock users: %v", result.Error)
	}

	sessionId1, _ := GetMockTxSession(t)
	sessionId2, _ := GetMockTxSession(t)

	user1 := User(1)
	wallet1, err := GetUserWallet(&user1, sessionId1)
	if err != nil || *wallet1 != Wallet(1) {
		t.Fatalf("failed to get proper wallet from first session: %v", err)
	}

	waitTime := 5 * time.Second
	go commitSessionAfterTime(sessionId1, waitTime)

	t1 := time.Now()
	wallet1, err = GetUserWallet(&user1, sessionId2)
	if err != nil || *wallet1 != Wallet(1) {
		t.Fatalf("failed to get proper wallet from second session: %v", err)
	}
	t2 := time.Now()

	if diff := t2.Sub(t1); diff < waitTime {
		t.Fatalf("retrieved user earlier than other session's commitment: %v", diff)
	}
}

func TestUpdateLockWithoutTx(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session")
	}

	mockUser, _ := getMockUser()
	userIdt := User(1)

	if result := db.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	_, err := GetUserWallet(&userIdt)
	if err != nil {
		t.Fatalf("failed to get user wallet: %v", err)
	}

	_, err = GetUserWallet(&userIdt)
	if err != nil {
		t.Fatalf("failed to get user wallet: %v", err)
	}
}

func TestDoubleLockInTransaction(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatalf("failed to initialize mock db")
	}

	if err := initialize(db); err != nil {
		t.Fatalf("failed to initialize session")
	}

	mockUser, _ := getMockUser()
	userId := User(1)

	if result := db.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start session %v", err)
	}

	_, err = getUserWallet(&userId, true, sessionId)
	if err != nil {
		t.Fatalf("failed to get user wallet: %v", err)
	}

	_, err = getUserWallet(&userId, true, sessionId)
	if err != nil {
		t.Fatalf("failed to get user wallet: %v", err)
	}

	// sessionId1, err := startSession()
	// if err != nil {
	// 	t.Fatalf("failed to start second session %v", err)
	// }

	// _, err = getUserWallet(&userId, true, sessionId1)
	// if err != nil {
	// 	t.Fatalf("failed to get user wallet: %v", err)
	// }
}
