package controllers

import (
	"basic-trade-api/models/variant"
	"basic-trade-api/services"
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func CreateVariant(ctx *gin.Context) {
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
	variantRequest, ok := requestInterface.(variant.VariantRequest)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast request to *variantRequest",
		})
		return
	}

	adminData, ok := ctx.MustGet("adminData").(jwt5.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to extract admin data",
		})
		return
	}

	adminIDFloat64, ok := adminData["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}
	adminID := int(adminIDFloat64)

	newVariant, err := services.CreateVariantService(dbConn, variantRequest, adminID)
	if err != nil {
		if err.Error() == "product does not belong to the admin" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	responseData := gin.H{
		"message": "Successfully created variant!",
		"data": gin.H{
			"id":          newVariant.ID,
			"uuid":        newVariant.UUID,
			"variantName": newVariant.VariantName,
			"quantity":    newVariant.Quantity,
			"productId":   newVariant.ProductID,
			"createdAt":   newVariant.CreatedAt,
			"updatedAt":   newVariant.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func GetAllVariant(ctx *gin.Context) {
	variantName := ctx.Query("variantName")
	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	pageSize := 10
	pageNum, _ := strconv.Atoi(ctx.Param("pageNum"))
	offset := (pageNum - 1) * pageSize

	if offset < 0 {
		offset = 0
	}

	getVariants, total, err := services.GetAllVariantService(dbConn, pageSize, offset, variantName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	responseData := gin.H{
		"message": "Successfully fetch variants!",
		"data":    getVariants,
		"meta": gin.H{
			"limit":     pageSize,
			"offset":    offset,
			"total":     total,
			"totalPage": totalPages,
		},
	}
	ctx.JSON(http.StatusOK, responseData)
}

func GetVariantByID(ctx *gin.Context) {
	variantUUID := ctx.Param("variantUUID")

	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	getVariant, err := services.GetVariantByIDService(dbConn, variantUUID)
	if err != nil {
		// Check if the error is due to product not found
		if err.Error() == "variant not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Variant not found",
				"message": "Variant with the specified UUID does not exist",
			})
			return
		}
		// If the error is not related to product not found, return internal server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseData := gin.H{
		"message": "Successfully fetched specific variant!",
		"data":    getVariant,
	}
	ctx.JSON(http.StatusOK, responseData)
}

func UpdateVariant(ctx *gin.Context) {
	variantUUID := ctx.Param("variantUUID")

	var variantRequest variant.VariantRequest
	if err := ctx.ShouldBind(&variantRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate the variant request data
	if err := variant.Validate.Struct(variantRequest); err != nil {
		// If validation fails, extract validation errors and return them
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation errors", "details": validationErrors})
		return
	}

	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	adminData, ok := ctx.MustGet("adminData").(jwt5.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to extract admin data",
		})
		return
	}
	adminIdFloat64, ok := adminData["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}
	adminId := int(adminIdFloat64)

	editVariant, err := services.UpdateVariantService(dbConn, variantRequest, variantUUID, adminId)
	if err != nil {
		// Check if the error is due to product not found
		if err.Error() == "variant not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Variant not found",
				"message": "Variant with the specified UUID does not exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseData := gin.H{
		"message": "Successfully update the variant!",
		"data": gin.H{
			"id":          editVariant.ID,
			"uuid":        editVariant.UUID,
			"variantName": editVariant.VariantName,
			"quantity":    editVariant.Quantity,
			"productId":   editVariant.ProductID,
			"createdAt":   editVariant.CreatedAt,
			"updatedAt":   editVariant.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusOK, responseData)
}

func DeleteVariant(ctx *gin.Context) {
	variantUUID := ctx.Param("variantUUID")

	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	adminData, ok := ctx.MustGet("adminData").(jwt5.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to extract admin data",
		})
		return
	}
	adminIdFloat64, ok := adminData["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}
	adminId := int(adminIdFloat64)

	removeVariant, err := services.DeleteVariantService(dbConn, variantUUID, adminId)
	if err != nil {
		if err.Error() == "variant not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Variant not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseData := gin.H{
		"message": "Successfully delete the variant!",
		"data": gin.H{
			"id":          removeVariant.ID,
			"uuid":        removeVariant.UUID,
			"variantName": removeVariant.VariantName,
			"quantity":    removeVariant.Quantity,
			"productId":   removeVariant.ProductID,
			"createdAt":   removeVariant.CreatedAt,
			"updatedAt":   removeVariant.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusOK, responseData)
}
