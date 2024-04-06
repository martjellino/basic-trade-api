package main

import (
	"basic-trade-api/database"
	"basic-trade-api/router"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	PORT = ":8000"
	DB   *sql.DB
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Start the database connection
	DB := database.StartDB()
	defer DB.Close()

	// Initialize the router
	r := router.StartApp(DB)
	fmt.Println("Server is running on", PORT)

	// Start the server
	r.Run(PORT)
}
