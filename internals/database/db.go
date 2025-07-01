package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, pass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Error while connecting to database:", err)
	} else {
		DB = db
		fmt.Println("Successfully connected to database")
	}
}
