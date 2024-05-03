package router

import (
	"basic-trade-api/controllers"
	"basic-trade-api/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func StartApp(db *sql.DB) *gin.Engine {
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Next()
	})

	adminRouter := router.Group("/auth")
	{
		adminRouter.POST("/register", middleware.RegisterValidator(), controllers.AdminRegister)
		adminRouter.POST("/login", middleware.LoginValidator(), controllers.AdminLogin)
	}

	productRouter := router.Group("/products")
	{
		productRouter.GET("/", controllers.GetAllProduct)
		productRouter.GET("/:productUUID", controllers.GetProductByID)
		// productRouter.Use(middleware.Authentication())
		productRouter.POST("/", middleware.Authentication(), middleware.ProductValidator(), controllers.CreateProduct)
		productRouter.PUT("/:productUUID", middleware.Authentication(), middleware.ProductAuthorization(), middleware.ProductValidator(), controllers.UpdateProduct)
		productRouter.DELETE("/:productUUID", middleware.Authentication(), middleware.ProductAuthorization(), controllers.DeleteProduct)
	}

	variantRouter := router.Group("/products/variants")
	{
		variantRouter.GET("/", controllers.GetAllVariant)
		variantRouter.GET("/:variantUUID", controllers.GetVariantByID)
		// variantRouter.Use(middleware.Authentication())
		variantRouter.POST("/", middleware.Authentication(), middleware.VariantValidator(), controllers.CreateVariant)
		variantRouter.PUT("/:variantUUID", middleware.Authentication(), middleware.VariantAuthorization(), middleware.VariantValidator(), controllers.UpdateVariant)
		variantRouter.DELETE("/:variantUUID", middleware.Authentication(), middleware.VariantAuthorization(), controllers.DeleteVariant)
	}
	return router
}
