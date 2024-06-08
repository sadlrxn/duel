package middlewares

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
)

func IPRestrict(blackList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.CheckIpAddress(blackList, c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Blacklisted IP address",
			})
			return
		}
	}
}

func CheckCountryCodeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var countryCode string
		var isLegal bool

		code, prs := c.Request.Header[http.CanonicalHeaderKey("CF-IPCountry")]

		if !prs || len(code) == 0 {
			countryCode, isLegal = utils.CheckIpRegionWithAbstractAPI(c.ClientIP(), config.BLACK_LIST)
		} else {
			countryCode = code[0]
			isLegal = !config.BLACK_LIST[countryCode]
		}

		if !isLegal {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"reason":      "IP address restricted.",
				"countryCode": countryCode,
			})
			return
		}

		c.Next()

	}
}
