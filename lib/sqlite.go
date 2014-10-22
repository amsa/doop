package doop

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Query(query string) string {
	db, err := sql.Open("sqlite3", "./test-sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: create the result set
	return "something"
}
