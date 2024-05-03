package controllers

import (
	"basic-trade-api/helpers"
	"basic-trade-api/models/product"
	"basic-trade-api/services"
	"database/sql"
	"math"
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func CreateProduct(ctx *gin.Context) {
	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	// Get the request data from context
	var productRequest product.ProductRequest
	if err := ctx.ShouldBind(&productRequest); err != nil {
		errors := helpers.GeneralValidator(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errors,
			"message": "Failed to validate",
		})
		return
	}

	// Validate the product request data
	if err := product.Validate.Struct(productRequest); err != nil {
		// If validation fails, extract validation errors and return them
		errors := helpers.GeneralValidator(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation errors", "details": errors})
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

	// Declare uploadResult outside the conditional block
	var uploadResult string

	// Check if productRequest.ImageFile is not nil
	if productRequest.ImageFile != nil {
		// Check if the uploaded file is an image (JPG, JPEG, PNG)
		contentType := productRequest.ImageFile.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/jpeg") &&
			!strings.HasPrefix(contentType, "image/jpg") &&
			!strings.HasPrefix(contentType, "image/png") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Only JPG, JPEG, and PNG images are allowed."})
			return
		}

		fileName := helpers.RemoveExtension(productRequest.ImageFile.Filename)
		// Assign the result of UploadFile to uploadResult
		var err error
		uploadResult, err = helpers.UploadFile(productRequest.ImageFile, fileName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Set uploaded file URL in the product request
		productRequest.ImageURL = uploadResult
	}

	newProduct, err := services.CreateProductService(dbConn, productRequest, adminId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseData := gin.H{
		"message": "Successfully created product!",
		"data": gin.H{
			"id":   newProduct.ID,
			"uuid": newProduct.UUID,
			"name": newProduct.Name,
			"imageUrl":  uploadResult,
			"adminId":   adminId,
			"createdAt": newProduct.CreatedAt,
			"updatedAt": newProduct.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func GetAllProduct(ctx *gin.Context) {
    name := ctx.Query("name")
    db, _ := ctx.Get("db")
    dbConn, ok := db.(*sql.DB)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to cast database connection to *sql.DB",
        })
        return
    }

    pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
    pageNum, _ := strconv.Atoi(ctx.Query("pageNum"))

    if pageNum < 1 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid page number",
        })
        return
    }

    offset := (pageNum - 1) * pageSize

    getProducts, total, err := services.GetAllProductService(dbConn, pageSize, offset, name)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
        })
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
    responseData := gin.H{
        "message": "Successfully fetch products!",
        "data":    getProducts,
        "meta": gin.H{
            "limit":     pageSize,
            "offset":    offset,
            "total":     total,
            "totalPage": totalPages,
        },
    }

    ctx.JSON(http.StatusOK, responseData)
}

func GetProductByID(ctx *gin.Context) {
	productUUID := ctx.Param("productUUID")

	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	getProduct, err := services.GetProductByIDService(dbConn, productUUID)
	if err != nil {
		// Check if the error is due to product not found
		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Product not found",
				"message": "Product with the specified UUID does not exist",
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
		"message": "Successfully fetched specific product!",
		"data":    getProduct,
	}
	ctx.JSON(http.StatusOK, responseData)
}

func UpdateProduct(ctx *gin.Context) {
	productUUID := ctx.Param("productUUID")

	db, _ := ctx.Get("db")
	dbConn, ok := db.(*sql.DB)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to cast database connection to *sql.DB",
		})
		return
	}

	// Get the request data from context
	var productRequest product.ProductRequest
	if err := ctx.ShouldBind(&productRequest); err != nil {
		errors := helpers.GeneralValidator(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errors,
			"message": "Failed to validate",
		})
		return
	}

	// Validate the product request data
	if err := product.Validate.Struct(productRequest); err != nil {
		errors := helpers.GeneralValidator(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation errors", "details": errors})
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

	// Declare uploadResult outside the conditional block
	var uploadResult string

	// Check if productRequest.ImageFile is not nil
	if productRequest.ImageFile != nil {
		fileName := helpers.RemoveExtension(productRequest.ImageFile.Filename)
		// Assign the result of UploadFile to uploadResult
		var err error
		uploadResult, err = helpers.UploadFile(productRequest.ImageFile, fileName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Set uploaded file URL in the product request
		productRequest.ImageURL = uploadResult
	}

	editProduct, err := services.UpdateProductService(dbConn, productRequest, productUUID, adminId)
	if err != nil {
		// Check if the error is due to product not found
		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Product not found",
				"message": "Product with the specified UUID does not exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	responseData := gin.H{
		"message": "Successfully update the product!",
		"data": gin.H{
			"id":        editProduct.ID,
			"uuid":      editProduct.UUID,
			"name":      editProduct.Name,
			"imageUrl": editProduct.ImageURL,
			"adminId":   editProduct.AdminID,
			"createdAt": editProduct.CreatedAt,
			"updatedAt": editProduct.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusOK, responseData)
}

func DeleteProduct(ctx *gin.Context) {
	productUUID := ctx.Param("productUUID")

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

	_, err := services.DeleteProductService(dbConn, productUUID, adminId)
	if err != nil {
		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	responseData := gin.H{
		"message": "Successfully delete the product!",
	}
	ctx.JSON(http.StatusOK, responseData)
}
