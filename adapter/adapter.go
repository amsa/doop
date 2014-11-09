package adapter

import (
	"database/sql"
	"strings"
)

type Adapter interface {
	Close() error
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	GetTableSchema() (map[string]string, error) //return all tables
}

func GetAdapter(dsn string) Adapter {
	DsnChunks := strings.Split(dsn, "://")
	if len(DsnChunks) != 2 {
		panic("Invalid DSN string: " + dsn)
	}
	var db Adapter = nil
	//var err error = nil
	switch strings.ToLower(DsnChunks[0]) {
	case "sqlite":
		db, _ = MakeSQLite(DsnChunks[1])
	default:
		panic("Invalid database connection: " + dsn)
	}
	return db
}
