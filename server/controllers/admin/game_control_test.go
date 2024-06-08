package admin

import (
	"strings"
	"testing"
	"time"
)

func TestGameController(t *testing.T) {
	InitGameController()
	gameStatus := gameController.GetTotalGameBlocked()
	gameNames := []string{
		GAME_CONTROLLER_COINFLIP,
		GAME_CONTROLLER_JACKPOT,
		GAME_CONTROLLER_DREAMTOWER,
		GAME_CONTROLLER_CRASH,
		GAME_CONTROLLER_PLINKO,
		GAME_CONTROLLER_BLACKJACK,
		GAME_CONTROLLER_DEPOSIT,
		GAME_CONTROLLER_WITHDRAW,
		GAME_CONTROLLER_SEED,
		GAME_CONTROLLER_GRAND_JACKPOT,
		GAME_CONTROLLER_DUEL_BOT,
		GAME_CONTROLLER_REWARDS,
		GAME_CONTROLLER_AFFILIATE,
	}
	for i, status := range gameStatus {
		if status.GameName != gameNames[i] {
			t.Fatalf(
				"game name mismatching. expected: %s, actual: %s",
				gameNames[i], status.GameName,
			)
		}
		if status.IsBlocked {
			t.Fatalf(
				"game should be not blocked at first time: %s",
				status.GameName,
			)
		}
	}

	if err := gameController.BlockGame(
		GAME_CONTROLLER_DREAMTOWER,
		true,
	); err != nil {
		t.Fatalf("failed to block game: %v", err)
	}
	if err := gameController.BlockGame(
		GAME_CONTROLLER_DREAMTOWER,
		true,
	); err == nil || !strings.Contains(err.Error(), "already blocked") {
		t.Fatalf("should be failed on already blocked game: %v", err)
	}
	if !gameController.GetGameBlocked(GAME_CONTROLLER_DREAMTOWER) {
		t.Fatalf("failed to block game properly")
	}
	gameStatus = gameController.GetTotalGameBlocked()
	if gameStatus[2].GameName != GAME_CONTROLLER_DREAMTOWER ||
		!gameStatus[2].IsBlocked {
		t.Fatalf("failed to block game properly: %v", gameStatus[2])
	}
	if err := gameController.BlockGame(
		GAME_CONTROLLER_DREAMTOWER,
		false,
	); err != nil {
		t.Fatalf("failed to unblock game: %v", err)
	}
	if err := gameController.BlockGame(
		GAME_CONTROLLER_DREAMTOWER,
		false,
	); err == nil || !strings.Contains(err.Error(), "already unblocked") {
		t.Fatalf("should be failed on already unblocked case")
	}
}

func TestTimeStringConversion(t *testing.T) {
	var t1 time.Time

	t.Fatalf("%s", t1.String())
}
