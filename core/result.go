package core

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Columns() ([]string, error)
	Err() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
