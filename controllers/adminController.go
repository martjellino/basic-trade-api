package controllers

import (
	"basic-trade-api/helpers"
	// "basic-trade-api/middleware"
	"basic-trade-api/models/admin"
	"basic-trade-api/services"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRegister(ctx *gin.Context) {
	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	requestInterface, ok := ctx.Get("request")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Parsed data not found in context",
		})
		return
	}

	// Get the request data
	adminRequest, ok := requestInterface.(admin.AdminRegisterRequest)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast request to *adminRequest",
		})
		return
	}

	newAdmin, err := services.AdminRegisterService(dbConn, adminRequest)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"message": "Failed to register admin",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	responseData := gin.H{
		"message": "Successfully created user!",
		"data": gin.H{
			"id":        newAdmin.ID,
			"uuid":      newAdmin.UUID,
			"name":      newAdmin.Name,
			"email":     newAdmin.Email,
			"createdAt": newAdmin.CreatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func AdminLogin(ctx *gin.Context) {
	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	requestInterface, ok := ctx.Get("request")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Parsed data not found in context",
		})
		return
	}

	// Get the request data
	adminRequest, ok := requestInterface.(admin.AdminLoginRequest)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast request to *adminRequest",
		})
		return
	}

	// Call the AdminRegisterService
	adminResponse, err := services.AdminLoginService(dbConn, adminRequest)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "user not found":
			statusCode = http.StatusNotFound
		case "invalid password":
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(adminResponse.ID, adminResponse.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	responseData := gin.H{
		"message": "User logged successfully",
		"data": gin.H{
			"email":       adminResponse.Email,
			"name":        adminResponse.Name,
			"accessToken": token,
		},
	}

	ctx.JSON(http.StatusOK, responseData)

}
