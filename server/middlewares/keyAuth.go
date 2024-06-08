package middlewares

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(apiAccessToken string) gin.HandlerFunc {

	return func(c *gin.Context) {
		token, prs := c.Request.Header[http.CanonicalHeaderKey("x-api-key")]

		if !prs || len(token) == 0 {
			utils.RespondWithError(c, 401, "API token required")
			return
		}

		if token[0] != apiAccessToken {
			utils.RespondWithError(c, 401, "Invalid API token")
			return
		}

		c.Next()
	}
}
