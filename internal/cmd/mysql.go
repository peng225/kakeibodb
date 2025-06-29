package cmd

import (
	"database/sql"
	"fmt"
	"time"
)

func OpenDB(dbName string, dbPort int, user string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:%d)/%s?parseTime=true",
		user, dbPort, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	return db, nil
}
