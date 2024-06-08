package tests

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetMockDbUrl() string {
	mockDbConf := config.Config{
		DBHost:     "localhost",
		DBUser:     "postgres",
		DBPassword: "2000115",
		DBName:     "mockduel",
		DBPort:     "5432",
	}

	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		mockDbConf.DBHost,
		mockDbConf.DBPort,
		mockDbConf.DBUser,
		mockDbConf.DBName,
		mockDbConf.DBPassword,
	)

	return dbUrl
}

func InitMockDB(drop bool, migrate bool) *gorm.DB {
	dbUrl := GetMockDbUrl()
	database := connectDB(dbUrl)

	if database == nil {
		return database
	}

	if drop {
		if err := dropTables(database); err != nil {
			fmt.Printf("error occurred while dropping tables: %v", err)
			return nil
		}
	}

	if migrate {
		if err := migrateTables(database); err != nil {
			fmt.Printf("error occurred while migrating tables: %v", err)
			return nil
		}
	}

	return database
}

func connectDB(DBUrl string) *gorm.DB {
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
	db, err := gorm.Open(postgres.Open(DBUrl), &gorm.Config{Logger: newLogger})

	if err != nil {
		fmt.Printf("error occurred while connecting to db: %v", err)
		return nil
	}

	return db
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
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
}

func migrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
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
}
