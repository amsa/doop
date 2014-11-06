package common

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
)

var debug bool = true

func SetDebug(val bool) {
	debug = val
}

func AddPrefix(origin string, prefix string) string {
	return fmt.Sprintf("__%s_%s", prefix, origin)
}

func Debug(values ...interface{}) {
	if debug {
		if len(values) == 0 {
			log.Println(values[0])
		} else {
			log.Printf(values[0].(string), values[1:]...)
		}
	}
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleErrorAny(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}

// generaeDbId generates the unique identifier for the given DSN
func GenerateDbId(dsn string) string {
	h := sha1.New()
	h.Write([]byte(dsn))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RowToStrings(rows *sql.Rows) (func() ([]string, error), error) {
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
