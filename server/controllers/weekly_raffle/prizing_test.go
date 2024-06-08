package weekly_raffle

import (
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestPerformWeeklyPrizing(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to init mock db")
	}
	db_aggregator.Initialize(db)

	users := []models.User{
		{
			Name: "taka1",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: utils.ConvertChipToBalance(1000),
					},
				},
			},
		},
		{
			Name: "taka2",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: utils.ConvertChipToBalance(2000),
					},
				},
			},
		},
		{
			Name: "taka3",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: utils.ConvertChipToBalance(1500),
					},
				},
			},
		},
		{
			Model: gorm.Model{
				ID: config.WEEKLY_RAFFLE_TEMP_ID,
			},
			Name: "WR_TEMP",
			Wallet: models.Wallet{
				Balance: models.Balance{
					ChipBalance: &models.ChipBalance{
						Balance: utils.ConvertChipToBalance(5000),
					},
				},
			},
		},
	}
	if err := db.Create(
		&users,
	).Error; err != nil {
		t.Fatal("failed to create users")
	}

	startedAt := datatypes.Date(time.Date(
		2023,
		time.February,
		22,
		0, 0, 0, 0,
		time.Local,
	))
	endAt := time.Date(
		2023,
		time.March,
		1,
		0, 0, 0, 0,
		time.Local,
	)
	raffle := models.WeeklyRaffle{
		StartedAt: startedAt,
		EndAt:     endAt,
		Prizes: pq.Int64Array{
			utils.ConvertChipToBalance(100),
			utils.ConvertChipToBalance(50),
			utils.ConvertChipToBalance(30),
		},
	}
	if err := db.Create(
		&raffle,
	).Error; err != nil {
		t.Fatal("failed to create weekly raffle record")
	}
	weeklyRaffle = &raffle

	if err := issueWeeklyRaffleTicket(
		1, 1,
	); err != nil {
		t.Fatal("failed to issue weekly raffle tickets")
	}

	if err := issueWeeklyRaffleTicket(
		2, 2,
	); err != nil {
		t.Fatal("failed to issue weekly raffle tickets")
	}

	if err := issueWeeklyRaffleTicket(
		3, 1,
	); err != nil {
		t.Fatal("failed to issue weekly raffle tickets")
	}

	if err := issueWeeklyRaffleTicket(
		1, 2,
	); err != nil {
		t.Fatal("failed to issue weekly raffle tickets")
	}

	if err := issueWeeklyRaffleTicket(
		3, 1,
	); err != nil {
		t.Fatal("failed to issue weekly raffle tickets")
	}

	tickets := []models.WeeklyRaffleTicket{}
	if err := db.Find(&tickets).Error; err != nil {
		t.Fatal("failed to retrieve all tickets")
	}

	if len(tickets) != 7 {
		t.Fatal("mismatching tickets")
	}

	expectedUserIDs := []uint{
		1, 2, 2, 3, 1, 1, 3,
	}
	for i, ticket := range tickets {
		if ticket.TicketID != uint(i+1) ||
			ticket.UserID != expectedUserIDs[i] ||
			ticket.Rank != nil ||
			ticket.Claimed != nil {
			t.Fatal("mismatching expected user ids")
		}
	}

	reviewWinners, err := getPrizingPreviewResult(
		[]uint{5, 2, 7},
		time.Time(startedAt),
	)
	if err != nil {
		t.Fatalf("failed to get prizing review: %v", err)
	}
	if len(reviewWinners.Winners) != 3 ||
		reviewWinners.Winners[0].UserID != 1 ||
		reviewWinners.Winners[1].UserID != 2 ||
		reviewWinners.Winners[2].UserID != 3 ||
		reviewWinners.Winners[0].Rank != 1 ||
		reviewWinners.Winners[1].Rank != 2 ||
		reviewWinners.Winners[2].Rank != 3 ||
		reviewWinners.Winners[0].Prize != utils.ConvertChipToBalance(100) ||
		reviewWinners.Winners[1].Prize != utils.ConvertChipToBalance(50) ||
		reviewWinners.Winners[2].Prize != utils.ConvertChipToBalance(30) {
		t.Fatal("failed to get proper review")
	}

	initSocket(make(chan types.WSEvent, 4096))
	prizingResult, err := performWeeklyPrizing(
		[]uint{5, 2, 7},
		time.Time(startedAt),
	)
	if err != nil {
		t.Fatalf("failed to perform prizing: %v", err)
	}
	if len(prizingResult.Winners) != 3 ||
		prizingResult.Winners[0].UserID != 1 ||
		prizingResult.Winners[1].UserID != 2 ||
		prizingResult.Winners[2].UserID != 3 ||
		prizingResult.Winners[0].Rank != 1 ||
		prizingResult.Winners[1].Rank != 2 ||
		prizingResult.Winners[2].Rank != 3 ||
		prizingResult.Winners[0].Prize != utils.ConvertChipToBalance(100) ||
		prizingResult.Winners[1].Prize != utils.ConvertChipToBalance(50) ||
		prizingResult.Winners[2].Prize != utils.ConvertChipToBalance(30) {
		t.Fatal(
			"failed to perform proper prizing",
			prizingResult.Winners[0].UserID,
			prizingResult.Winners[0].Rank,
			prizingResult.Winners[0].Prize,
		)
	}

	tickets = []models.WeeklyRaffleTicket{}
	if err := db.Order("id asc").Find(&tickets).Error; err != nil {
		t.Fatal("failed to retrieve all tickets")
	}

	if len(tickets) != 7 {
		t.Fatal("mismatching tickets")
	}

	expectedRanks := []uint{
		999,
		1,
		999,
		999,
		0,
		999,
		2,
	}
	for i, ticket := range tickets {
		rank := uint(999)
		if ticket.Rank != nil {
			rank = *ticket.Rank
		}
		if rank != expectedRanks[i] {
			t.Fatal(
				"failed to save proper ranks",
				i,
				rank,
				expectedRanks[i],
				ticket.ID,
			)
		}
	}

	raffle = models.WeeklyRaffle{}
	if err := db.First(&raffle).Error; err != nil {
		t.Fatal(
			"failed to retrieve weekly raffle",
			err,
		)
	}
	if !raffle.Ended {
		t.Fatal("failed to update weekly raffle ended")
	}

	claimed, err := claimReward(
		1,
		[]uint{5},
	)
	if err != nil {
		t.Fatal(
			"failed to claim reward",
			err,
		)
	}
	if claimed != utils.ConvertChipToBalance(100) {
		t.Fatal("failed to claim proper amount")
	}

	user1 := models.User{}
	if err := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"name = ?",
		"taka1",
	).First(&user1).Error; err != nil {
		t.Fatal("failed to retrieve taka1 info")
	}

	wr_temp := models.User{}
	if err := db.Preload(
		"Wallet.Balance.ChipBalance",
	).Where(
		"name = ?",
		"WR_TEMP",
	).First(&wr_temp).Error; err != nil {
		t.Fatal("failed to retrieve wr_temp info")
	}

	if user1.Wallet.Balance.ChipBalance.Balance != utils.ConvertChipToBalance(1100) ||
		wr_temp.Wallet.Balance.ChipBalance.Balance != utils.ConvertChipToBalance(4900) {
		t.Fatal("failed to claim proper amount")
	}

	tx := models.Transaction{}
	if err := db.First(&tx).Error; err != nil {
		t.Fatal("failed to retrieve transaction")
	}
	if tx.Type != models.TxClaimWeeklyRaffleReward ||
		*tx.FromWallet != wr_temp.Wallet.ID ||
		*tx.ToWallet != user1.Wallet.ID {
		t.Fatal("failed to leave transaction properly")
	}
}
