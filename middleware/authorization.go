package middleware

import (
	// "basic-trade-api/services"
	"database/sql"

	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, _ := ctx.Get("db")
		dbConn, ok := db.(*sql.DB)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to cast database connection to *sql.DB",
			})
			return
		}

		productUUID := ctx.Param("productUUID")
		adminData := ctx.MustGet("adminData").(jwt5.MapClaims)
		adminDataId := uint(adminData["id"].(float64))
		query := "SELECT admin_id FROM products WHERE uuid = $1"
		var adminID uint
		err := dbConn.QueryRow(query, productUUID).Scan(&adminID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   "Data Not Found",
					"message": err.Error(),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		if adminID != adminDataId {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}

func VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, _ := ctx.Get("db")
		dbConn, ok := db.(*sql.DB)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to cast database connection to *sql.DB",
			})
			return
		}

		variantUUID := ctx.Param("variantUUID")
		adminData := ctx.MustGet("adminData").(jwt5.MapClaims)
		adminDataId := uint(adminData["id"].(float64))
		query := `SELECT p.admin_id FROM variants v JOIN products p ON v.product_id = p.id WHERE v.uuid = $1`
		var adminID uint
		err := dbConn.QueryRow(query, variantUUID).Scan(&adminID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error":   "Data Not Found",
					"message": err.Error(),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
			return
		}

		if adminID != adminDataId {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}
