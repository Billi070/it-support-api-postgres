package main

import (
	"database/sql"
	"log"
	"os"


	//_ "modernc.org/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func initDB() {
	var err error

	//db, err = sql.Open("sqlite", "tickets.db")

db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed", err)
	}
log.Println("Connected to postgreSQL")

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS tickets (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		title TEXT NOT NULL,
	// 		description TEXT NOT NULL,
	// 		priority TEXT NOT NULL,
	// 		status TEXT NOT NULL
	// 	)
	// `)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
}