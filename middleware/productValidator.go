package middleware

import (
	"basic-trade-api/helpers"
	"basic-trade-api/models/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProductValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productRequest product.ProductRequest
		if err := ctx.ShouldBind(&productRequest); err != nil {
			errors := helpers.GeneralValidator(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   errors,
				"message": "Failed to validate",
			})
			return
		}

		ctx.Set("request", productRequest)
		ctx.Next()
	}
}
