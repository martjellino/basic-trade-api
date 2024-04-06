package middleware

import (
	"basic-trade-api/helpers"
	"basic-trade-api/models/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerRequest admin.AdminRegisterRequest
		if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
            errors := helpers.GeneralValidator(err)
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error":   errors,
                "message": "Failed to validate request",
            })
            return
        }

        // Validate the request using the Validate struct
        if err := admin.Validate.Struct(registerRequest); err != nil {
            errors := helpers.GeneralValidator(err)
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error":   errors,
                "message": "Validation errors",
            })
            return
        }

		ctx.Set("request", registerRequest)
		ctx.Next()
	}
}
