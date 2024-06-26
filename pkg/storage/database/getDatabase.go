package database

import (
	"database/sql"
	"fmt"
	"src/env"

	_ "github.com/lib/pq"
)

func GetConnection(e *env.Env) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		e.PostgreHost, e.PostgrePort, e.PostgreUser, e.PostgrePassword, e.PostgreDB, e.PostgreSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
