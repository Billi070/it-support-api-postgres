package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite", "tickets.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			priority TEXT NOT NULL,
			status TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}