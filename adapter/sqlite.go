package adapter

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SQLite struct {
	dbPath string
	db     *sql.DB
}

func MakeSQLite(databasePath string, args ...interface{}) (Adapter, error) {
	db, err := sql.Open("sqlite3", databasePath)
	ret := Adapter(&SQLite{databasePath, db})
	return ret, err
}

func (sqliteDb *SQLite) Close() error {
	return sqliteDb.db.Close()
}

func (sqliteDb *SQLite) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return sqliteDb.db.Query(sql, args...)
}

func (sqliteDb *SQLite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return sqliteDb.db.Exec(sql, args...)
}

func (sqliteDb *SQLite) GetSchema(table_name string) (string, error) {
	var ret string
	query := "SELECT sql FROM sqlite_master WHERE name=?"
	row := sqliteDb.db.QueryRow(query, table_name)
	if row == nil {
		log.Fatal("invalid table_name")
	}
	row.Scan(&ret)
	return ret, nil
}

func (sqliteDb *SQLite) GetTables() ([]string, error) {
	ret := make([]string, 0, 64)
	query := "SELECT name FROM sqlite_master WHERE type='table'"
	rows, err := sqliteDb.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		ret = append(ret, name)
	}
	return ret, nil
}
