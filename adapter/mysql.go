package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dbName string
	db     *sql.DB
}

func MakeMySQL(dsn string, args ...interface{}) (Adapter, error) {
	db, err := sql.Open("mysql", dsn)
	dsnChunks := strings.Split(dsn, "/")
	ret := Adapter(&MySQL{dbName: dsnChunks[len(dsnChunks)-1], db: db})
	return ret, err
}

func (mysqlDb *MySQL) CreateDb() (bool, error) {
	_, err := mysqlDb.Exec("CREATE DATABASE " + mysqlDb.dbName)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (mysqlDb *MySQL) DropDb() (bool, error) {
	_, err := mysqlDb.Exec("DROP DATABASE IF EXISTS " + mysqlDb.dbName)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (mysqlDb *MySQL) Close() error {
	return mysqlDb.db.Close()
}

func (mysqlDb *MySQL) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return mysqlDb.db.Query(sql, args...)
}

func (mysqlDb *MySQL) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return mysqlDb.db.Exec(sql, args...)
}

func (mysqlDb *MySQL) GetTableSchema() (map[string]string, error) {
	ret := make(map[string]string)
	query := "SELECT `TABLE_NAME` FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA` = ?;"
	rows, err := mysqlDb.db.Query(query, mysqlDb.dbName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tbl string
	var schema string
	for rows.Next() {
		var tblname string
		rows.Scan(&tblname)
		sqlRow := mysqlDb.db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`;", mysqlDb.dbName, tblname))
		sqlRow.Scan(&tbl, &schema)
		ret[tblname] = schema
	}
	return ret, nil
}
