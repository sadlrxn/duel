package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/mixpanel"
	"github.com/Duelana-Team/duelana-v1/controllers/payment"
	"github.com/Duelana-Team/duelana-v1/controllers/prelude"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/solana"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/routes"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initConfig() *config.Config {
	config.Init()
	config := config.Get()
	// log.LogMessage("main thread", "initializing config done...", "success", logrus.Fields{})
	return &config
}

func initLog() {
	log.Init()
	log.LogMessage("main thread", "initializing log done...", "success", logrus.Fields{})
}

func initDB(config *config.Config) *gorm.DB {
	DBUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", config.DBHost, config.DBPort, config.DBUser, config.DBName, config.DBPassword)
	database := db.ConnectDB(DBUrl)
	log.LogMessage("main thread", "initializing db done...", "success", logrus.Fields{})
	return database
}

func initTxModule(database *gorm.DB) {
	transaction.Initialize(database)
	log.LogMessage("main thread", "initializing tx module done...", "success", logrus.Fields{})
}

func initRoute(config *config.Config) {
	r := routes.Init()
	if config.ENV == "dev" {
		r.Use(static.Serve("/", static.LocalFile("../build", true)))
		r.NoRoute(func(c *gin.Context) {
			c.File("../build/index.html")
		})
	}

	if err := r.Run(fmt.Sprintf(":%v", config.AppPort)); err != nil {
		log.LogMessage("main thread", fmt.Sprintf("failed to run server %v", err), "error", logrus.Fields{})
		return
	}

	log.LogMessage("main thread", "initializing router done...", "success", logrus.Fields{})
}

func initSolana(config *config.Config) {
	initParam := solana.InitParam{
		TreasuryBs58: config.MasterWalletPriKey,
		Cluster:      solana.ClusterDevNet,
		RpcUrl:       config.SolanaRpcUrl,
	}
	if config.Network == "mainnet" {
		initParam.Cluster = solana.ClusterMainNetBeta
	}
	solana.Initialize(&initParam)
}

func initMixpanel(config *config.Config) {
	mixpanel.Init(config.MixpanelToken, config.MixpanelServerUrl)
	log.LogMessage("main thread", "initializing mixpanel done...", "success", logrus.Fields{})

	c := cron.New()
	c.AddFunc("0 * * * *", mixpanel.TrackChipAmount)
	c.Start()
}

func initGlobalTimezone() {
	time.Local = time.UTC
}

func initServerConfig() {
	db := db.GetDB()
	if db == nil {
		log.LogMessage(
			"init server config",
			"failed to retrieve db pointer",
			"error",
			logrus.Fields{},
		)
		return
	}

	var serverConfig models.ServerConfig
	if result := db.First(&serverConfig); result.Error != nil {
		log.LogMessage(
			"init server config",
			"failed to retrieve server config",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
		return
	}

	config.SetServerConfig(serverConfig)
}

func migrations() {
	db := db.GetDB()
	var serverConfig models.ServerConfig
	if result := db.First(&serverConfig); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		serverConfig.ShouldNotReset = true
		serverConfig.BaseRakeBackRate = config.BASE_RAKEBACK_RATE
		serverConfig.AdditionalRakeBackRate = config.ADDITIONAL_RAKEBACK_RATE
		serverConfig.NextGrandJackpotStartAt = config.GRAND_JACKPOT_NEXT_ROUND
		byte, err := json.Marshal(config.API_RATE_LIMIT_CONFIGURATION)
		if err != nil {
			log.LogMessage("migrate server configs", "falied to msrshal api rate limit configuratioins", "error", logrus.Fields{"error": err.Error()})
			return
		}
		serverConfig.ApiRateLimitConfiguration = string(byte)
		byte, err = json.Marshal(config.WEBSOCKET_RATE_LIMIT_CONFIGURATION)
		if err != nil {
			log.LogMessage("migrate server configs", "falied to msrshal websocket rate limit configuratioins", "error", logrus.Fields{"error": err.Error()})
			return
		}
		serverConfig.WebsocketRateLimitConfiguration = string(byte)
		if result := db.Create(&serverConfig); result.Error != nil {
			log.LogMessage("migrate server configs", "falied to create server config", "error", logrus.Fields{"error": err.Error()})
			return
		}
		log.LogMessage("migrate server configs", "successfully created", "success", logrus.Fields{})
		return
	}
	if !serverConfig.ShouldNotReset {
		serverConfig.ShouldNotReset = true
		serverConfig.BaseRakeBackRate = config.BASE_RAKEBACK_RATE
		serverConfig.AdditionalRakeBackRate = config.ADDITIONAL_RAKEBACK_RATE
		serverConfig.NextGrandJackpotStartAt = config.GRAND_JACKPOT_NEXT_ROUND
		byte, err := json.Marshal(config.API_RATE_LIMIT_CONFIGURATION)
		if err != nil {
			log.LogMessage("migrate server configs", "falied to msrshal api rate limit configuratioins", "error", logrus.Fields{"error": err.Error()})
			return
		}
		serverConfig.ApiRateLimitConfiguration = string(byte)
		byte, err = json.Marshal(config.WEBSOCKET_RATE_LIMIT_CONFIGURATION)
		if err != nil {
			log.LogMessage("migrate server configs", "falied to msrshal websocket rate limit configuratioins", "error", logrus.Fields{"error": err.Error()})
			return
		}
		serverConfig.WebsocketRateLimitConfiguration = string(byte)
		if result := db.Save(&serverConfig); result.Error != nil {
			log.LogMessage("migrate server configs", "falied to save server config", "error", logrus.Fields{"error": err.Error()})
			return
		}
		log.LogMessage("migrate server configs", "successfully saved", "success", logrus.Fields{})
	}
	log.LogMessage("migrate server configs", "done...", "success", logrus.Fields{})
}

func initRedis() {
	if err := redis.InitRedis(
		config.Get().RedisUrl,
		// config.Get().RedisPwd,
		"",
	); err != nil {
		log.LogMessage(
			"init_redis",
			"failed to initialize",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}
}

func initialize() {
	// log.LogMessage("main thread", "initializing...", "info", logrus.Fields{})
	config := initConfig()
	initLog()
	database := initDB(config)
	initRedis()
	if err := prelude.InitDuelMainUsers(database); err != nil {
		log.LogMessage("init duel main users", "failed", "error", logrus.Fields{"error": err.Error()})
	}
//	if err := prelude.InitDuelBotsV2(database); err != nil {
//		log.LogMessage("init duel bots", "failed", "error", logrus.Fields{"error": err.Error()})
//	}
	initTxModule(database)
	// if err := prelude.PerformRemoveOdes(); err != nil {
	// 	log.LogMessage("remove odes", "failed", "error", logrus.Fields{"error": err.Error()})
	// }
	// if err := prelude.MigrateRakebackRecords(); err != nil {
	// 	log.LogMessage("migrate rakeback records", "failed", "error", logrus.Fields{"error": err.Error()})
	// }
	if err := prelude.InitAffiliateLifetime(); err != nil {
		log.LogMessage(
			"initialize affiliate lifetime",
			"failed to initialize",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}
	initSolana(config)
	// initMixpanel(config)
	if err := payment.MigratePaymentModel(); err != nil {
		log.LogMessage("payment migration", "failed", "error", logrus.Fields{"error": err.Error()})
	}
	migrations()
	initServerConfig()
	// initGlobalTimezone()
	admin.InitGameController()
	initRoute(config)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.LogMessage(
				"main thread",
				"recovered",
				"info",
				logrus.Fields{
					"recover": r,
					"stack":   string(debug.Stack()),
				},
			)
		}
	}()

	initialize()
}
