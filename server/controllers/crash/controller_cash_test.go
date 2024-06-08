package crash

// import (
// 	"testing"

// 	"github.com/Duelana-Team/duelana-v1/config"
// 	"github.com/Duelana-Team/duelana-v1/models"
// 	"github.com/Duelana-Team/duelana-v1/tests"
// 	"gorm.io/gorm"
// )

// func initCrashMockDB() (*GameController, error) {
// 	db := tests.InitMockDB(true, true)

// 	mockUsers := []models.User{
// 		{
// 			Name: "taka",
// 			Wallet: models.Wallet{
// 				Balance: models.Balance{
// 					ChipBalance: &models.ChipBalance{
// 						Balance: 100000000,
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Model: gorm.Model{ID: config.CRASH_TEMP_ID},
// 			Name:  "CH_TEMP",
// 			Wallet: models.Wallet{
// 				Balance: models.Balance{
// 					ChipBalance: &models.ChipBalance{
// 						Balance: 100000000,
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Model: gorm.Model{ID: config.CRASH_FEE_ID},
// 			Name:  "CH_FEE",
// 			Wallet: models.Wallet{
// 				Balance: models.Balance{
// 					ChipBalance: &models.ChipBalance{
// 						Balance: 0,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	c := GameController{}
// 	return &c, nil
// }

// func TestCashIn(t *testing.T) {

// }
