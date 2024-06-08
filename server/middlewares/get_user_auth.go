package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
* @External
* Returns authorized userID from middleware.
* If not authorized returns 0.
 */
func GetAuthUserID(
	ctx *gin.Context,
	abortForNotAuth bool,
) uint {
	user, prs := ctx.Get(AuthMiddleware().IdentityKey)
	if !prs {
		if abortForNotAuth {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid permission",
			})
		}
		return 0
	}

	userID, ok := user.(gin.H)["id"].(uint)
	if !ok {
		if abortForNotAuth {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid permission",
			})
		}
		return 0
	}

	return userID
}
