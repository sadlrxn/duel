package admin

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCouponHandler(ctx *gin.Context) {
	var params coupon.CreateCouponRequest
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
	}
	params.Balance = utils.ConvertChipToBalance(params.Balance)

	if code, err := coupon.Create(params); err == nil {
		if params.Shortcut != nil &&
			*params.Shortcut != "" {
			if err := coupon.CreateShortcut(
				code,
				*params.Shortcut,
			); err == nil {
				ctx.JSON(
					http.StatusOK,
					gin.H{
						"message":  "succeed to create coupon code with shortcut",
						"code":     code,
						"shortcut": *params.Shortcut,
					},
				)
			} else {
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"message": "succeed to create coupon but failed for shortcut",
						"error":   err.Error(),
						"code":    code,
					},
				)
			}
			return
		}

		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "succeed to create coupon code",
				"code":    code,
			},
		)
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to create coupon code",
				"error":   err.Error(),
			},
		)
	}
}

func CreateCouponShortcutHandler(ctx *gin.Context) {
	var params struct {
		Code     uuid.UUID `json:"code"`
		Shortcut string    `json:"shortcut"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
	}

	if err := coupon.CreateShortcut(
		params.Code,
		params.Shortcut,
	); err == nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message":  "succeed to create coupon shortcut",
				"code":     params.Code,
				"shortcut": params.Shortcut,
			},
		)
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to set shortcut",
				"error":   err.Error(),
			},
		)
	}
}

func DeleteCouponShortcutHandler(ctx *gin.Context) {
	var params struct {
		Shortcut string `json:"shortcut"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
	}

	if err := coupon.DeleteShortcut(
		params.Shortcut,
	); err == nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message":  "succeed to delete coupon shortcut",
				"shortcut": params.Shortcut,
			},
		)
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to delete shortcut",
				"error":   err.Error(),
			},
		)
	}
}
