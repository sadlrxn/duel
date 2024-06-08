package crash

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestGetFirstUnplayedRoundID(t *testing.T) {
	db := tests.InitMockDB(true, true)

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("%v", err)
	}
	if err := autoMigrateCrashRound(); err != nil {
		t.Fatalf("%v", err)
	}

	getRoundHistory(10)
	t.Fatal("asdf")
}
