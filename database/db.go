package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func StartDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// user := os.Getenv("PG_USER")
	// password := os.Getenv("PG_PASSWORD")
	// host := os.Getenv("PG_HOST")
	// dbname := os.Getenv("PG_DBNAME")
	// sslmode := os.Getenv("PG_SSLMODE")
	// port := os.Getenv("PG_PORT")

	// config := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	// db, err := sql.Open("postgres", config)
	// if err != nil {
	// 	panic(err)
	// }

	// Constructing connection string
	// connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", user, password, host, dbname, sslmode)

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Opening a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// rows, err := db.Query("select version()")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// var version string
	// for rows.Next() {
	// 	err := rows.Scan(&version)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// fmt.Printf("version=%s\n", version)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
