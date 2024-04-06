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
		productRouter.POST("/", middleware.Authentication(), middleware.ProductValidator(), controllers.CreateProduct)
		productRouter.GET("/", controllers.GetAllProduct)
		productRouter.GET("/:productUUID", controllers.GetProductByID)
		productRouter.PUT("/:productUUID", middleware.Authentication(), middleware.ProductValidator(), controllers.UpdateProduct)
		productRouter.DELETE("/:productUUID", middleware.Authentication(), controllers.DeleteProduct)
	}

	variantRouter := router.Group("/products/variants")
	{
		variantRouter.POST("/", middleware.Authentication(), middleware.VariantValidator(), controllers.CreateVariant)
		variantRouter.GET("/", controllers.GetAllVariant)
		variantRouter.GET("/:variantUUID", controllers.GetVariantByID)
		variantRouter.PUT("/:variantUUID", middleware.Authentication(), middleware.VariantValidator(), controllers.UpdateVariant)
		variantRouter.DELETE("/:variantUUID", middleware.Authentication(), controllers.DeleteVariant)
	}

	return router
}
