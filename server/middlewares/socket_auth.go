package middlewares

import (
	b64 "encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"

	"github.com/gagliardetto/solana-go"
)

func SocketAuthMiddleware() (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "duelana zone",
		Key:           []byte("JWT_SECRET_KEY"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		DisabledAbort: true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(gin.H); ok {
				return jwt.MapClaims{
					"id":            v["id"],
					"name":          v["name"],
					"walletAddress": v["walletAddress"],
					"role":          v["role"],
					"avatar":        v["avatar"],
					"ipAddress":     v["ipAddress"],
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return gin.H{
				"id":            uint(claims["id"].(float64)),
				"name":          claims["name"].(string),
				"walletAddress": claims["walletAddress"].(string),
				"role":          models.Role(claims["role"].(string)),
				"avatar":        claims["avatar"].(string),
				"ipAddress":     claims["ipAddress"].(string),
			}

		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var logInVals LogIn
			if err := c.ShouldBind(&logInVals); err != nil {
				log.LogMessage("User Authenticator", "Invalid login parameter", "error", logrus.Fields{"Error": err.Error()})
				return "", errors.New("missing walletAddress or signature")
			}
			db := db.GetDB()

			walletAddress := logInVals.WalletAddress
			sigOrTx := logInVals.SigOrTx

			var user models.User
			if result := db.Where("wallet_address = ?", walletAddress).First(&user); result.Error != nil {
				return nil, errors.New("not exist user")
			}

			result := db.Find(&user)

			if result.Error != nil {
				return nil, errors.New("can't find user with walletAddress")
			}

			signMsg := []byte("Sign in Duelana with nonce: " + user.Nonce)

			// ledger support

			if len(sigOrTx) > 300 {
				return nil, errors.New("invalid signature length")
			}

			var tx solana.Transaction
			err := tx.UnmarshalBase64(b64.StdEncoding.EncodeToString(sigOrTx))

			if err != nil {

				pubKey, _ := solana.PublicKeyFromBase58(walletAddress)

				sig := solana.SignatureFromBytes(sigOrTx)

				if sig.Verify(pubKey, signMsg) {
					// struct {ID uint, Name string, WalletAddress string, Role string}
					return gin.H{
						"id":            user.ID,
						"name":          user.Name,
						"walletAddress": walletAddress,
						"role":          user.Role,
						"avatar":        user.Avatar,
						"ipAddress":     c.ClientIP(),
					}, nil
				}
			}

			duel, err := solana.PublicKeyFromBase58("DUELLrBB96snTu3Wn3Cjyj7s2pRRqiG5LpPCC1fmw2Wm")
			bytes, _ := base58.Decode(tx.Message.Instructions[0].Data.String())

			if err != nil {
				return nil, errors.New("failed to decode base58 string")
			}

			userWallet, err := solana.PublicKeyFromBase58(walletAddress)
			if err != nil {
				return nil, errors.New("failed to decode base58 string")
			}
			if len(tx.Message.AccountKeys) != 2 {
				return nil, errors.New("mismatching account counts")
			}
			if tx.Message.AccountKeys[0] != userWallet {
				return nil, errors.New("invalid walletAddress param")
			}
			if tx.Message.AccountKeys[len(tx.Message.AccountKeys)-1] != duel {
				return nil, errors.New("invalid account")
			}
			if string(bytes) != string(signMsg) {
				return nil, errors.New("invalid nonce")
			}

			if err := tx.VerifySignatures(); err != nil {
				return nil, errors.New("invalid Signature")
			}

			// ledger end

			return gin.H{
				"id":            user.ID,
				"name":          user.Name,
				"walletAddress": walletAddress,
				"role":          user.Role,
				"avatar":        user.Avatar,
				"ipAddress":     c.ClientIP(),
			}, nil

		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.Next()
		},
		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		CookieName:     "duelana",
		CookieSameSite: http.SameSiteStrictMode,
		TokenLookup:    "cookie:duelana",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
	})

	if err != nil {
		log.LogMessage("Auth Middleware", "JWT Error occured", "error", logrus.Fields{"Error": err.Error()})
	}

	return
}
