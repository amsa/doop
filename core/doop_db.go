package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/amsa/doop/adapter"
	. "github.com/amsa/doop/common"
)

type DoopDbInfo struct {
	DSN  string
	Name string
	Hash string
}

type DoopDb struct {
	dsn     string
	adapter adapter.Adapter
}

func MakeDoopDb(dsn string) *DoopDb {
	return &DoopDb{dsn, adapter.GetAdapter(dsn)}
}

func (doopdb *DoopDb) Init() error {
	// Create doop_master table to store metadata of all tables
	statement := fmt.Sprintf(`
			CREATE TABLE %s (
				id integer NOT NULL PRIMARY KEY,
				name text,
				type text,
				branch text,
				sql text,
			)	
		`, DOOP_MASTER)
	HandleErrorAny(doopdb.adapter.Exec(statement))

	//Insert tables into doop_master
	tables, err := doopdb.adapter.GetTables()
	HandleErrorAny(tables, err)

	for i, table := range tables {
		schema, err := doopdb.adapter.GetSchema(table)
		HandleErrorAny(schema, err)
		statement = fmt.Sprintf(`
				INSERT INTO %s VALUES (
					%i,		
					%s,
					%s,
					%s,
					%sql,	
				)	
			`, i, table, "logical_table", DOOP_DEFAULT_BRANCH, schema)
	}
	// Create branch table to store the branches
	HandleErrorAny(doopdb.adapter.Exec(`CREATE TABLE ` + DOOP_TABLE_BRANCH + ` (
			id integer NOT NULL PRIMARY KEY,
			name     text,
			parent   text,
			metadata text
		);`))
	HandleErrorAny(doopdb.adapter.Exec(`CREATE UNIQUE INDEX __branch_name_idx ON ` + DOOP_TABLE_BRANCH + ` (name);`))

	// Create default branch
	_, err = doopdb.CreateBranch(DOOP_DEFAULT_BRANCH, "")
	return err

}

func (doopdb *DoopDb) Clean() error {
	_, err := doopdb.adapter.Exec(`DROP INDEX IF EXISTS __branch_name_idx;`)
	if err != nil {
		return err
	}
	_, err = doopdb.adapter.Exec(`DROP TABLE ` + DOOP_TABLE_BRANCH + `;`)
	if err != nil {
		return err
	}
	return nil
}

// initAdapter initializes Doop database adapter based on the given DSN string
func (doopdb *DoopDb) initAdapter(dsn string) {
	doopdb.adapter = adapter.GetAdapter(dsn)
}

// CreateBranch creates a new branch of the database forking from the given parent branch
func (doopdb *DoopDb) CreateBranch(branchName string, parentBranch string) (bool, error) {
	if branchName != DOOP_DEFAULT_BRANCH && parentBranch == "" {
		return false, errors.New("Parent branch name is not specified.")
	}
	// insert a row with branch and its parent name along with metadata (empty json object for now)
	HandleErrorAny(doopdb.adapter.
		Exec(`INSERT INTO `+DOOP_TABLE_BRANCH+` (name, parent, metadata) VALUES (?, ?, '{}')`, branchName, parentBranch))

	//get all physical table
	tables, err := doopdb.adapter.GetTables()

	if err != nil {
		return false, err
	}

	//create companian tables for each logical table
	for tableName, schema := range tables {
		statements := make([]string, 0, 64)

		//need to parse the schema to create h
		//vdel
		//TODO: parse schema to get primary key type,
		//if no primary key, error thrown
		vdel := fmt.Sprintf(`
			CREATE TABLE __%s_%s_vdel (
				%s
			)	
		`, branchName, tableName, schema)
		statements = append(statements, vdel)

		//hdel

		//vsec

		//hsec
	}
	return true, nil
}

func (doopdb *DoopDb) GetTables() map[string]string {
	//find out the name of tables in default branch
	statement := fmt.Sprintf(`
		SELECT name, sql FROM %s WHERE branch=? AND type=?
	`)
	rows, err := doopdb.adapter.Query(statement, DOOP_MASTER, "logical_table")
	HandleErrorAny(rows, err)

	ret := make(map[string]string)
	for rows.Next() {
		var name string
		var sql string
		err := rows.Scan(&name, &sql)
		HandleErrorAny(rows, err)
		ret[name] = sql
	}
	return ret
}

// ListBranches returns the list of all the branches for the given database
func (doopdb *DoopDb) ListBranches() []string {
	rt := make([]string, 1)
	rows, err := doopdb.adapter.Query(`SELECT name FROM ` + DOOP_TABLE_BRANCH + `;`)
	HandleError(err)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		rt = append(rt, name)
	}
	return rt
}

// RemoveBranch deletes a branch
func (doopdb *DoopDb) RemoveBranch(branchName string) (bool, error) {
	return false, nil
}

// MergeBranch merges two branches into the first one (from)
func (doopdb *DoopDb) MergeBranch(from string, to string) (bool, error) {
	return false, nil
}

/*
Query interface
it translates the query, which is on logical tables, to query on concrete tables
*/
func (doopDb *DoopDb) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

/*
Exec interface
it transalate the exec statement, which is on logical tables, to exec statement on concreate tables
*/

func (doopDb *DoopDb) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
