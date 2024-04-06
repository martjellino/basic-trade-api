package middleware

import (
	"basic-trade-api/helpers"
	"basic-trade-api/models/variant"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
)

func VariantValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var variantRequest variant.VariantRequest
        if err := ctx.ShouldBindJSON(&variantRequest); err != nil {
            errors := helpers.GeneralValidator(err)
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error":   errors,
                "message": "Failed to validate request",
            })
            return
        }

        // Validate the request using the Validate struct
        if err := variant.Validate.Struct(variantRequest); err != nil {
            errors := helpers.GeneralValidator(err)
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error":   errors,
                "message": "Validation errors",
            })
            return
        }

		ctx.Set("request", variantRequest)
		ctx.Next()
	}
}
