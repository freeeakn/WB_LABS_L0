package db

import (
	"database/sql"
	"time"

	"wb_labs_l0_backend/internal/logger"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string, log logger.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(15)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Infof("connected to db")
	return db, nil
}
