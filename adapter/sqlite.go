package adapter

import (
	"database/sql"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
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

func (sqliteDb *SQLite) DropDb() (bool, error) {
	err := exec.Command("rm", "-rf", sqliteDb.dbPath).Run()
	if err == nil {
		return true, nil
	}
	return false, err
}

func (sqliteDb *SQLite) CreateDb() (bool, error) {
	return true, nil
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

func (sqliteDb *SQLite) GetTableSchema() (map[string]string, error) {
	ret := make(map[string]string)
	query := "SELECT name, sql FROM sqlite_master WHERE type='table' AND tbl_name <> 'sqlite_sequence';"
	rows, err := sqliteDb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var sql string
		rows.Scan(&name, &sql)
		ret[name] = sql
	}
	return ret, nil
}
