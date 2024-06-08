package dreamtower

import (
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/lib/pq"
)

type saveRoundRequest struct {
	bets   *pq.Int32Array
	status *models.DreamTowerStatus
	profit *int64
}
