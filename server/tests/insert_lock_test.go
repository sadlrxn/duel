package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"gorm.io/gorm/clause"
)

func TestInsertingAfterLock(t *testing.T) {
	db := InitMockDB(true, true)
	if db == nil {
		t.Fatalf("nil db pointer retrieved")
	}

	mockUser := models.User{}
	if result := db.Create(&mockUser); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	mockCrashRound := models.CrashRound{}
	if result := db.Create(&mockCrashRound); result.Error != nil {
		t.Fatalf("failed to create mock crash round: %v", result.Error)
	}

	mockCrashBets := []models.CrashBet{
		{
			UserID:  1,
			RoundID: 1,
		},
		{
			UserID:  1,
			RoundID: 1,
		},
		{
			UserID:  1,
			RoundID: 1,
		},
	}
	if result := db.Create(&mockCrashBets); result.Error != nil {
		t.Fatalf("failed to create mock crash bets: %v", result.Error)
	}

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("failed to initialize db session: %v", err)
	}

	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		t.Fatalf("faile to get session: %v", err)
	}

	go func() {
		mockCrashBets = []models.CrashBet{}
		if result := session.Clauses(
			clause.Locking{
				Strength: "UPDATE",
			},
		).Where(
			"user_id = ? and round_id = ?",
			1, 1,
		).Find(&mockCrashBets); result.Error != nil {
			fmt.Printf("failed to lock and retrieve crash bets: %v", result.Error)
		}
		time.Sleep(time.Second * 5)

		if err := db_aggregator.CommitSession(sessionId); err != nil {
			fmt.Printf("failed to commit session: %v", err)
		}
	}()

	time.Sleep(time.Second)
	startTime := time.Now()
	newCrashBet := models.CrashBet{
		UserID:  1,
		RoundID: 1,
	}
	if result := session.Create(
		&newCrashBet,
	); result.Error != nil {
		fmt.Println(
			"failed to create new crash bet",
			result.Error,
		)
	}
	endTime := time.Now()

	if endTime.Sub(startTime) < time.Second*4 {
		t.Fatalf(
			"failed to lock properly: left: %d, right: %d",
			endTime.Sub(startTime),
			time.Second*4,
		)
	}

}
