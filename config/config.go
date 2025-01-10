package config

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)


func InitDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully!")
	return db, nil
}
