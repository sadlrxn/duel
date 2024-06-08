package db

import (
	"log"
	"os"
	"time"

	Logger "github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDB(DBUrl string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	var err error
	db, err = gorm.Open(postgres.Open("postgresql://duel:Backtohome1111@db:5433/duel?sslmode=disable"), &gorm.Config{Logger: newLogger, SkipDefaultTransaction: true})

	if err != nil {
		Logger.LogMessage("connect db", "failed", "error", logrus.Fields{"error": err.Error()})
		logrus.Fatal(err.Error())
	}

	Logger.LogMessage("connect db", "success", "info", logrus.Fields{})

	err = db.AutoMigrate(
		&models.Payment{},
		&models.NftCollection{},
		&models.User{},
		&models.Statistics{},
		&models.Wallet{},
		&models.Transaction{},
		&models.ChipBalance{},
		&models.NftBalance{},
		&models.Balance{},
		&models.DepositedNft{},
		&models.JackpotRound{}, &models.JackpotPlayer{}, &models.JackpotBet{},
		&models.CoinflipRound{},
		&models.NftInGame{},
		&models.ClientSeed{}, &models.ServerSeed{}, &models.SeedPair{},
		&models.DreamTowerRound{},
		&models.DuelBot{},
		&models.Rakeback{},
		&models.ServerConfig{},
		&models.Affiliate{},
		&models.ActiveAffiliate{},
		&models.AffiliateLifetime{},
		&models.Coupon{},
		&models.ClaimedCoupon{},
		&models.CouponTransaction{},
		&models.CrashRound{},
		&models.CrashBet{},
		&models.SelfExclusion{},
		&models.DailyRaceRewards{},
		&models.CouponShortcut{},
		&models.WeeklyRaffleTicket{},
		&models.WeeklyRaffle{},
	)

	if err != nil {
		Logger.LogMessage("migrate db", "failed", "error", logrus.Fields{"error": err.Error()})
		logrus.Fatal(err.Error())
	}

	return db
}

// GetDB ...
func GetDB() *gorm.DB {
	return db
}
