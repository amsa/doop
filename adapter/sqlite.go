package adapter

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	dbPath string
	db     *sql.DB
}

func MakeSQLite(databasePath string) (Adapter, error) {
	db, err := sql.Open("sqlite3", databasePath)
	ret := Adapter(&SQLite{databasePath, db})
	return ret, err
}

func (sqliteDb *SQLite) Close() (bool, error) {
	return true, sqliteDb.db.Close()
}

func (sqliteDb *SQLite) Query(sql string) (*sql.Rows, error) {
	//no branch involved, simply return the query on base table
	return sqliteDb.db.Query(sql)
}

func (sqliteDb *SQLite) Exec(sql string) (sql.Result, error) {
	return sqliteDb.db.Exec(sql)
}
