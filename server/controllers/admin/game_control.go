package admin

import (
	"fmt"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

type GameController struct {
	gameBlocked syncmap.Map
	gameNames   []string
}

const GAME_CONTROLLER_COINFLIP = "Coinflip"
const GAME_CONTROLLER_JACKPOT = "Jackpot"
const GAME_CONTROLLER_DREAMTOWER = "Dreamtower"
const GAME_CONTROLLER_CRASH = "Crash"
const GAME_CONTROLLER_PLINKO = "Plinko"
const GAME_CONTROLLER_BLACKJACK = "Blackjack"
const GAME_CONTROLLER_DEPOSIT = "Deposit"
const GAME_CONTROLLER_WITHDRAW = "Withdraw"
const GAME_CONTROLLER_SEED = "Seed"
const GAME_CONTROLLER_GRAND_JACKPOT = "GrandJackpot"
const GAME_CONTROLLER_DUEL_BOT = "DuelBot"
const GAME_CONTROLLER_REWARDS = "Rewards"
const GAME_CONTROLLER_AFFILIATE = "Affiliate"

func (c *GameController) Init() {
	c.gameNames = []string{
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
	c.gameBlocked = syncmap.Map{}
	for _, gameName := range c.gameNames {
		c.gameBlocked.Store(gameName, false)
	}
	log.LogMessage(
		"admin_game_controller",
		"successfully initialized",
		"success",
		logrus.Fields{
			"gameNames": c.gameNames,
		},
	)
}

func (c *GameController) BlockGame(gameName string, block bool) error {
	blocked, prs := c.gameBlocked.Load(gameName)
	if !prs {
		return utils.MakeError(
			"admin-game-control",
			"BlockGame",
			"invalid parameter",
			fmt.Errorf("%s: is not known game", gameName),
		)
	}
	isBlocked, ok := blocked.(bool)
	if !ok {
		return utils.MakeError(
			"admin-game-control",
			"BlockGame",
			"failed to get bool type",
			fmt.Errorf("actual value: %v", blocked),
		)
	}
	if isBlocked == block {
		prefix := ""
		if !block {
			prefix = "un"
		}
		return utils.MakeError(
			"admin-game-control",
			"BlockGame",
			"game already (un)blocked",
			fmt.Errorf(
				"%s game already %sblocked",
				gameName,
				prefix,
			),
		)
	}

	c.gameBlocked.Store(gameName, block)
	return nil
}

func (c *GameController) GetGameBlocked(gameName string) bool {
	blocked, prs := c.gameBlocked.Load(gameName)
	isBlocked, ok := blocked.(bool)
	return !prs || !ok || isBlocked
}

type GameStatus struct {
	GameName  string `json:"gameName"`
	IsBlocked bool   `json:"isBlocked"`
}

func (c *GameController) GetTotalGameBlocked() []GameStatus {
	result := []GameStatus{}
	for _, gameName := range c.gameNames {
		result = append(result, GameStatus{
			GameName:  gameName,
			IsBlocked: c.GetGameBlocked(gameName),
		})
	}

	return result
}

var gameController GameController

func InitGameController() {
	gameController.Init()
}

func GameControllerMiddleware(gameName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gameController.GetGameBlocked(gameName) {
			c.AbortWithStatusJSON(
				http.StatusServiceUnavailable,
				gin.H{
					"message": fmt.Sprintf(
						"%s is blocked  by admin.",
						gameName,
					),
				},
			)
			return
		}
		c.Next()
	}
}

func GetGameBlocked(gameName string) bool {
	return gameController.GetGameBlocked(gameName)
}
