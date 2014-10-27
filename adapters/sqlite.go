package adapter

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	dsn     string
	db_name string
	db      *sql.DB
}

func MakeSQLite(dsn string, db_name string, connection_path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", connection_path)
	ret := &SQLite{dsn, db_name, db}
	return ret, err
}

func (sqliteDb *SQLite) Close() error {
	return sqliteDb.db.Close()
}

func (sqliteDb *SQLite) Query(branchName string, statement string, args ...interface{}) (*sql.Rows, error) {
	//no branch involved, simply return the query on base table
	rows, err := sqliteDb.db.Query(statement, args)
	return rows, err
}

func (sqliteDb *SQLite) Exec(branchName string, statement string, args ...interface{}) (sql.Result, error) {
	result, err := sqliteDb.db.Exec(statement, args)
	return result, err
}

func (sqliteDb *SQLite) CreateBranch(branchName string, baseBranch string) (bool, error) {
	return false, nil
}

func (sqliteDb *SQLite) RemoveBranch(branchName string) (bool, error) {
	return false, nil
}

func (sqliteDb *SQLite) MergeBranch(to string, from string) (bool, error) {
	return false, nil
}

func (sqliteDb *SQLite) ListBranches() ([]string, error) {
	return nil, nil
}
