package adapter

import (
	"database/sql"
	"strings"
)

type Adapter interface {
	Close() (bool, error)
	Query(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	GetTables() ([]string, error)
	GetSchema(table_name string) (string, error)
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

func rowToStrings(rows *sql.Rows) (func() ([]string, error), error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	return func() ([]string, error) {
		ptrs := make([]interface{}, len(columns))
		rawResult := make([][]byte, len(columns))
		result := make([]string, len(columns))
		for i, _ := range ptrs {
			ptrs[i] = &rawResult[i]
		}
		if rows.Next() {
			err := rows.Scan(ptrs...)
			if err != nil {
				return nil, err
			}
			for i, raw := range rawResult {
				if raw == nil {
					result[i] = "\\N"
				} else {
					result[i] = string(raw)
				}
			}
			return result, nil
		} else {
			return nil, nil
		}
	}, nil
}
