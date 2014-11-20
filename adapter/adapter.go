package adapter

import (
	"database/sql"
	"strings"
)

/*type Database struct {
	db *sql.DB
}

func (database *Database) SetDb(db *sql.DB) {
	database.db = db
}

func (database *Database) Close() error {
	return database.db.Close()
}

func (database *Database) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return database.db.Query(sql, args...)
}

func (database *Database) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return database.db.Exec(sql, args...)
}*/

type Adapter interface {
	Close() error
	DropDb() (bool, error)
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
