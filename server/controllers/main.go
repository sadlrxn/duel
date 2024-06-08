package controllers

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/chat"
	"github.com/Duelana-Team/duelana-v1/controllers/coinflip"
	"github.com/Duelana-Team/duelana-v1/controllers/crash"
	"github.com/Duelana-Team/duelana-v1/controllers/daily_race"
	"github.com/Duelana-Team/duelana-v1/controllers/dreamtower"
	"github.com/Duelana-Team/duelana-v1/controllers/grand_jackpot"
	"github.com/Duelana-Team/duelana-v1/controllers/jackpot"
	"github.com/Duelana-Team/duelana-v1/controllers/payment"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/controllers/weekly_raffle"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Chat chat.Controller
var User user.Controller
var Payment payment.Controller
var Coinflip coinflip.Controller
var JackpotLow jackpot.Controller
var JackpotMedium jackpot.Controller
var JackpotWild jackpot.Controller
var GrandJackpot grand_jackpot.Controller
var Dreamtower dreamtower.Controller
var Crash crash.GameController

func Init(eventEmitter chan types.WSEvent) {
	Chat = chat.Controller{EventEmitter: eventEmitter}
	User = user.Controller{EventEmitter: eventEmitter, Chat: &Chat}
	Payment = payment.Controller{EventEmitter: eventEmitter}
	Coinflip = coinflip.Controller{EventEmitter: eventEmitter}
	JackpotLow = jackpot.Controller{EventEmitter: eventEmitter, Room: types.Jackpot, Type: models.Low}
	JackpotMedium = jackpot.Controller{EventEmitter: eventEmitter, Room: types.Jackpot, Type: models.Medium}
	JackpotWild = jackpot.Controller{EventEmitter: eventEmitter, Room: types.Jackpot, Type: models.Wild}
	GrandJackpot = grand_jackpot.Controller{EventEmitter: eventEmitter}
	Dreamtower = dreamtower.Controller{}
	Crash = crash.GameController{EventEmitter: eventEmitter}
	if err := daily_race.Initialize(eventEmitter); err != nil {
		log.LogMessage(
			"controllers_Init",
			"failed to initialize daily race module",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}
	weekly_raffle.Initialize(eventEmitter)
	if config.CRASH_START_ON_SERVER_STARTUP &&
		config.Get().ENV != "dev" {
		if err := Crash.Start(); err != nil {
			log.LogMessage(
				"controllers_Init",
				"failed to init crash on start up",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
		}
	}
}

func GetServerConfig(ctx *gin.Context) {
	user, _ := ctx.Get(middlewares.SocketAuthMiddleware().IdentityKey)
	var userID *uint
	if user != nil {
		id := user.(gin.H)["id"].(uint)
		userID = &id
	}

	serverConfig := config.GetServerConfig()

	var rakebackRate uint
	var err error
	if userID == nil {
		rakebackRate = serverConfig.BaseRakeBackRate + serverConfig.AdditionalRakeBackRate
		if rakebackRate > config.RAKEBACK_MAX {
			rakebackRate = config.RAKEBACK_MAX
		}
	} else {
		rakebackRate, err = db_aggregator.GetUserRakebackRate(db_aggregator.User(*userID))
		if err != nil {
			log.LogMessage("GetServerConfig", "Failed to get user's rakeback rate.", "error", logrus.Fields{"user": *userID, "error": err.Error()})
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"coinflip": Coinflip.GetMeta(),
			"jackpot": gin.H{
				"low":    JackpotLow.GetMeta(),
				"medium": JackpotMedium.GetMeta(),
				"wild":   JackpotWild.GetMeta(),
			},
			"dreamtower":   Dreamtower.GetMeta(),
			"grandJackpot": GrandJackpot.GetMeta(),
			"crash":        Crash.GetMeta(),
		},
		"config": gin.H{
			"balanceDecimals":  config.BALANCE_DECIMALS,
			"rakebackRate":     rakebackRate,
			"crashClientSeed":  serverConfig.CrashClientSeed,
			"couponWagerTimes": config.COUPON_REQUIRED_WAGER_TIMES,
			"couponLifeTime":   config.COUPON_BALANCE_LIFE_TIME_IN_HOURS,
			"couponMaxClaim":   config.COUPON_MAXIMUM_EXCHANGE,
		},
	})
}
