package adapter

import (
	"database/sql"
	"log"

	"github.com/amsa/doop-core/core"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	metaInfo string
}

func (sqliteDb *SQLite) getInfo() string {
	return ""
}

func (sqliteDb *SQLite) executeSQL(query string) (*core.Result, error) {
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
	return nil, nil
}

func (sqliteDb *SQLite) createBranch(branchName string, baseBranch string) (bool, error) {
	return false, nil
}

func (sqliteDb *SQLite) removeBranch(branchName string) (bool, error) {
	return false, nil
}

func (sqliteDb *SQLite) listBranches() ([]string, error) {
	return nil, nil
}
