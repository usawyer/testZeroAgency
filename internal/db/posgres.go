package database

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initDB() (*sql.DB, error) {
	connectionParams := map[string]string{
		"host":     getEnv("DB_HOST", "localhost"),
		"user":     getEnv("POSTGRES_USER", "postgres"),
		"password": getEnv("POSTGRES_PASSWORD", "postgres"),
		"dbname":   getEnv("POSTGRES_DB", "test"),
		"port":     getEnv("DB_PORT", "5432"),
		"sslmode":  "disable",
		"TimeZone": "Asia/Novosibirsk",
	}

	var dsn string
	for key, value := range connectionParams {
		dsn += fmt.Sprintf("%s=%s ", key, value)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 2)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Println(err)
			continue
		}

		db.SetMaxOpenConns(10)
		db.SetConnMaxLifetime(time.Minute * 5)

		err = db.Ping()
		if err == nil {
			return db, nil
		}
	}

	return nil, errors.New("Error open db")
}
