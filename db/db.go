package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func New(cfg Config) *sql.DB {
	conn, err := sql.Open("mysql", cfg.Source())
	if err != nil {
		log.Panicln(err)
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)
	return conn
}
