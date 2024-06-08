package prelude

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestInitDuelUsers(t *testing.T) {
	db := tests.InitMockDB(true, true)

	if err := InitDuelMainUsers(db); err != nil {
		t.Fatalf("failed to init duel main users %v", err)
	}

	expectedUsers := getInitialUsers()

	mainUsers := []models.User{}
	if result := db.Find(&mainUsers); result.Error != nil {
		t.Fatalf("failed to get main duel users: %v", result.Error)
	}

	if len(expectedUsers) != len(mainUsers) {
		t.Fatalf("mismatching number of users, expected: %d, recorded: %d", len(expectedUsers), len(mainUsers))
	}

	for i, user := range mainUsers {
		if user.ID != expectedUsers[i].ID ||
			user.Name != expectedUsers[i].Name ||
			user.WalletAddress != expectedUsers[i].WalletAddress {
			t.Fatalf("mismatching data \n\rexpected: %v\n\rrecorded: %v", expectedUsers[i], user)
		}
	}
}

func TestInitDuplication(t *testing.T) {
	TestInitDuelUsers(t)

	db := tests.InitMockDB(false, false)

	if result := db.Model(&models.ChipBalance{}).Where("id is not null").Update("balance", 100); result.Error != nil {
		t.Fatalf("failed to update chip balances: %v", result.Error)
	}

	balances := []models.ChipBalance{}
	if result := db.Where("balance != ?", 100).Find(&balances); result.Error != nil || len(balances) > 0 {
		t.Fatalf("failed to update chip balances properly: %v, %v", result.Error, balances)
	}

	if err := InitDuelMainUsers(db); err != nil {
		t.Fatalf("failed to init duel main users %v", err)
	}

	balances = []models.ChipBalance{}
	if result := db.Where("balance != ?", 100).Find(&balances); result.Error != nil || len(balances) > 0 {
		t.Fatalf("balances are changed while duplicated initialization: %v, %v", result.Error, balances)
	}
}

func TestMoreInits(t *testing.T) {
	db := tests.InitMockDB(true, true)

	initUsers := getInitialUsers()

	firstInit := []models.User{}
	for i := 0; i < 5; i++ {
		firstInit = append(firstInit, buildUserAccount(initUsers[i]))
	}

	if result := db.Create(&firstInit); result.Error != nil {
		t.Fatalf("failed to init first users: %v", result.Error)
	}

	db.Model(&models.ChipBalance{}).Where("id is not null").Update("balance", 100)

	firstInit = []models.User{}
	if result := db.Find(&firstInit); result.Error != nil || len(firstInit) != 5 {
		t.Fatalf("failed to init first users properly: %v, %v", result.Error, firstInit)
	}

	if err := InitDuelMainUsers(db); err != nil {
		t.Fatalf("failed to perform 2nd init: %v", err)
	}

	secondInit := []models.User{}
	if result := db.Find(&secondInit); result.Error != nil || len(secondInit) != len(initUsers) {
		t.Fatalf("failed to init 2nd users properly: %v, %v", secondInit, len(secondInit))
	}
}
