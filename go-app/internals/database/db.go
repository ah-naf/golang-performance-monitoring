package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables directly.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, pass)
	db, err := sql.Open("postgres", connStr)

	if err == nil {
		if pingErr := db.Ping(); pingErr == nil {
			DB = db
			log.Println("Connected to existing database:", dbName)
			return
		}
	}

	log.Printf("Database %q not found, creating...", dbName)

	adminDSN := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=postgres password=%s sslmode=disable",
		host, port, user, pass,
	)
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		log.Fatalf("Could not connect to postgres database: %v", err)
	}
	defer adminDB.Close()

	if err = adminDB.Ping(); err != nil {
		log.Fatalf("Could not ping postgres database: %v", err)
	}

	safeName := pq.QuoteIdentifier(dbName)
	if _, err := adminDB.Exec("CREATE DATABASE " + safeName); err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}
	log.Printf("Database %q created.", dbName)

	newDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to new database %s: %v", dbName, err)
	}
	if err := newDB.Ping(); err != nil {
		log.Fatalf("New database %s ping failed: %v", dbName, err)
	}

	initPath := filepath.Join("init.sql") // adjust if needed
	sqlBytes, err := os.ReadFile(initPath)
	if err != nil {
		log.Fatalf("Could not read init.sql: %v", err)
	}

	if _, err := newDB.Exec(string(sqlBytes)); err != nil {
		log.Fatalf("Failed to execute init.sql: %v", err)
	}
	log.Println("init.sql applied successfully.")

	DB = newDB
	log.Println("Database initialization complete and connected.")
}
