package seed

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetActiveSeed(ctx *gin.Context) {
	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = user.(gin.H)["id"].(uint)

	seedPair, err := getActiveUserSeedPair(db_aggregator.User(userID))
	if err != nil {
		log.LogMessage(
			"GetActiveSeed",
			"Failed to get active seed",
			"error",
			logrus.Fields{
				"user":  userID,
				"error": err.Error(),
			},
		)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"clientSeed":         seedPair.ClientSeed.Seed,
		"serverSeedHash":     seedPair.ServerSeed.Hash,
		"nonce":              seedPair.Nonce,
		"nextServerSeedHash": seedPair.NextServerSeed.Hash,
	})
}

func RotateSeed(ctx *gin.Context) {
	var params struct {
		ServerSeedHash string `json:"serverSeedHash"`
		ClientSeed     string `json:"clientSeed"`
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameter."})
		return
	}

	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = user.(gin.H)["id"].(uint)

	newSeedPair, err := rotateUserSeedPair(
		db_aggregator.User(userID),
		params.ServerSeedHash,
		params.ClientSeed,
	)
	if err != nil {
		log.LogMessage("RotateSeed", "Failed to rotate seed", "error", logrus.Fields{"user": userID, "error": err.Error()})
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to rotate seed."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"clientSeed":         newSeedPair.ClientSeed.Seed,
		"serverSeedHash":     newSeedPair.ServerSeed.Hash,
		"nonce":              newSeedPair.Nonce,
		"nextServerSeedHash": newSeedPair.NextServerSeed.Hash,
	})
}

func UnhashServerSeed(ctx *gin.Context) {
	var params struct {
		Hash string `form:"hash"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameter."})
		return
	}

	serverSeed, err := getUnhashedServerSeed(params.Hash)
	if err != nil {
		log.LogMessage("UnhashServerSeed", "Failed to get expired server seed", "error", logrus.Fields{"hash": params.Hash, "error": err.Error()})
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot unhash server seed."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"seed": serverSeed})
}

func GetExpiredSeeds(ctx *gin.Context) {
	var params struct {
		Offset int `form:"offset"`
		Count  int `form:"count"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("seed history", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = user.(gin.H)["id"].(uint)

	seedPairs, err := getExpiredUserSeedPairs(
		db_aggregator.User(userID),
		uint(params.Offset),
		uint(params.Count),
	)
	if err != nil {
		log.LogMessage("GetExpiredSeeds", "Failed to get expired seed pairs", "error", logrus.Fields{"user": userID, "error": err.Error()})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := []gin.H{}
	for _, seedPair := range seedPairs {
		result = append(
			result,
			gin.H{
				"clientSeed":     seedPair.ClientSeed.Seed,
				"serverSeed":     seedPair.ServerSeed.Seed,
				"serverSeedHash": seedPair.ServerSeed.Hash,
				"nonce":          seedPair.Nonce,
			},
		)
	}

	ctx.JSON(http.StatusOK, result)
}
