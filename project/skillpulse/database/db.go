package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "skillpulse")
	password := getEnv("DB_PASSWORD", "skillpulse123")
	dbname := getEnv("DB_NAME", "skillpulse")
	ssl := getEnv("DB_SSL", "false")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	if ssl == "true" {
		dsn += "&tls=true"
	} else if ssl == "skip-verify" {
		dsn += "&tls=skip-verify"
	}

	var err error
	maxAttempts := 30
	if os.Getenv("VERCEL") == "1" {
		maxAttempts = 3
	}

	for i := 0; i < maxAttempts; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				log.Println("Connected to MySQL database")
				DB.SetMaxOpenConns(10)
				DB.SetMaxIdleConns(5)
				DB.SetConnMaxLifetime(5 * time.Minute)
				return nil
			}
		}

		log.Printf("Waiting for database... attempt %d/%d (error: %v)", i+1, maxAttempts, err)
		if i < maxAttempts-1 {
			sleepDur := 2 * time.Second
			if os.Getenv("VERCEL") == "1" {
				sleepDur = 500 * time.Millisecond
			}
			time.Sleep(sleepDur)
		}
	}

	return fmt.Errorf("could not connect to database after %d attempts: %w", maxAttempts, err)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
