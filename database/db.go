package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func StartDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	dbname := os.Getenv("PG_DBNAME")
	sslmode := os.Getenv("PG_SSLMODE")

	// config := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// db, err := sql.Open("postgres", config)
	// if err != nil {
	// 	panic(err)
	// }

	// Constructing connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", user, password, host, dbname, sslmode)

	// Opening a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
	  log.Fatal(err)
	}
	defer db.Close()
  
	rows, err := db.Query("select version()")
	if err != nil {
	  log.Fatal(err)
	}
	defer rows.Close()
  
	var version string
	for rows.Next() {
	  err := rows.Scan(&version)
	  if err != nil {
		log.Fatal(err)
	  }
	}
	fmt.Printf("version=%s\n", version)

	return db
}