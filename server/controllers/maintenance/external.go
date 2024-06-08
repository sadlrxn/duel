package maintenance

import "github.com/gin-gonic/gin"

func Initialize() error {
	return initialize()
}

// func Prepare() error {
// 	return prepare()
// }

func Maintain() error {
	return maintain()
}

func Finish() error {
	return finish()
}

func Current() MaintenanceDetails {
	return current()
}

func AbleToBet() bool {
	return ableToBet()
}

func StartMaintenance(ctx *gin.Context) {
	if err := maintain(); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"msg": "started maintenance status"})
}

func FinishMaintenance(ctx *gin.Context) {
	if err := finish(); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"msg": "finished maintenance status"})
}
