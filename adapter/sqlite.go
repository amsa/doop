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

func (sqliteDb *SQLite) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	//no branch involved, simply return the query on base table
	return sqliteDb.db.Query(sql)
}

func (sqliteDb *SQLite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return sqliteDb.db.Exec(sql)
}

// TODO: The following methods should be part of Doop not SQLite
//func (sqliteDb *SQLite) CreateBranch(branchName string, baseBranch string) (bool, error) {
//return false, nil
//}

//func (sqliteDb *SQLite) RemoveBranch(branchName string) (bool, error) {
//return false, nil
//}

//func (sqliteDb *SQLite) MergeBranch(to string, from string) (bool, error) {
//return false, nil
//}

//func (sqliteDb *SQLite) ListBranches() ([]string, error) {
//return nil, nil
//}
