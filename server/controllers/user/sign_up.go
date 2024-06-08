package user

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Internal
// Handles sign up action for a new user.
func signUpHandler(ctx *gin.Context, request CreateUserRequest) error {

	// 1. Check whether reached account count limit by IP address.
	reachedLimit, err := checkIpAccountLimit(ctx.ClientIP())
	if err != nil {
		return utils.MakeError(
			"user_sign_up",
			"signUpHandler",
			"failed to check IP account limit",
			err,
		)
	} else if reachedLimit {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Already reached the account count limit.",
		})
		return utils.MakeErrorWithCode(
			"user_sign_up",
			"signUpHandler",
			"already reached the account count limit.",
			ErrCodeAlreadyReachedAccountLimit,
			nil,
		)
	}

	// 2. Retry 5 times on failure.
	for i := 0; i < 5; i++ {
		// 2.1. Generate random user name.
		userName, err := generateUserName()
		if err != nil {
			// 2.2. Set user name as the first 16 characters of wallet address
			log.LogMessage(
				"user_request_nonce",
				"failed to generate random user name for a new user",
				"error",
				logrus.Fields{
					"error": err.Error(),
				},
			)
			userName = request.WalletAddress[:16]
		}
		if isLongerThanMax(userName) {
			log.LogMessage(
				"user_request_nonce",
				"generated user name is longer than max length",
				"error",
				logrus.Fields{
					"name": userName,
				},
			)
			userName = request.WalletAddress[:16]
		}

		request.Name = userName

		// 2. Create a new user to DB.
		if err := createUser(request); err == nil {
			return nil
		} else if utils.IsErrorCode(
			err,
			ErrCodeFailedToCreateNewUser,
		) {
			log.LogMessage(
				"user_sign_up_signUpHandler",
				"retry sign up since failed to create a new user in DB",
				"error",
				logrus.Fields{
					"error": err.Error(),
					"user":  request,
				},
			)
			continue
		} else {
			return utils.MakeError(
				"user_sign_up",
				"signUpHandler",
				"failed to create new user",
				err,
			)
		}
	}
	return utils.MakeError(
		"user_sign_up",
		"signUpHandler",
		"failed all the sign up tries",
		nil,
	)
}

/*
* @Internal
* Checks whether user name is logner that max length.
 */
func isLongerThanMax(name string) bool {
	return len(name) > 16
}
