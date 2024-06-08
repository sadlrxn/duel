package wager

import "github.com/Duelana-Team/duelana-v1/models"

type PlayerInPerformAfterWagerParams struct {
	UserID uint
	Bet    int64
	Profit int64
}

type PerformAfterWagerParams struct {
	Players     []PlayerInPerformAfterWagerParams
	Type        models.GameType
	IsHouseGame bool
}
