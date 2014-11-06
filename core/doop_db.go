package core

/*
DoopDb provides the same interface as Adapter.
However, DoopDb is logical database adpater, which means that the tables in DoopDb are not arranged
as how they look like.

Thus DoopDb is responsible for translating operations on (logical) tables to those on concrete tables.

Also, DoopDb provides another set of APIs for branch.
*/
import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/amsa/doop/adapter"
	. "github.com/amsa/doop/common"
)

const (
	DOOP_DEFAULT_BRANCH = "master"
	DOOP_TABLE_BRANCH   = "__branch"
	DOOP_MASTER         = "__doop_master"
)

type DoopDbInfo struct {
	DSN  string
	Name string
	Hash string
}

type DoopDb struct {
	info    *DoopDbInfo
	adapter adapter.Adapter
}

func MakeDoopDb(info *DoopDbInfo) *DoopDb {
	return &DoopDb{info, adapter.GetAdapter(info.DSN)}
}

/*

	DooDb Creation and Elimination

*/

func (doopdb *DoopDb) createDoopMaster() error {
	// Create doop_master table to store metadata of all tables
	statement := fmt.Sprintf(`
			CREATE TABLE %s (
				id integer NOT NULL PRIMARY KEY,
				name text,
				type text,
				branch text,
				sql text
			)	
		`, DOOP_MASTER)
	_, err := doopdb.adapter.Exec(statement)
	if err != nil {
		return err
	}

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
					%s	
				)	
			`, i, table, "logical_table", DOOP_DEFAULT_BRANCH, schema)
	}
	return nil
}

/*
Initialize the database to doopDb, it:

* create doop_master table
* insert existing tables info to doop_master table
* create branch table to store information of branches
* create a default(master) branch

*/
func (doopdb *DoopDb) Init() error {

	// Create the doop_master table
	err := doopdb.createDoopMaster()
	HandleError(err)

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

/*


				Standard SQL adapter interace

	it translates the operation, which is on branch on logical tables, to corresponding ones on concrete tables

*/

/*
Query interface
*/
func (doopDb *DoopDb) Query(branchName string, sql string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

/*
Exec interface
*/

func (doopDb *DoopDb) Exec(branchName string, sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (doopdb *DoopDb) GetSchema(branchName string, tableName string) ([]string, error) {
	return nil, nil
}

func (doopdb *DoopDb) GetAllTables() ([]string, error) {
	return doopdb.adapter.GetTables()
}
func (doopdb *DoopDb) GetTables(branchName string) map[string]string {
	//find out the name of tables in default branch
	statement := fmt.Sprintf(`
		SELECT name, sql FROM %s WHERE branch=? AND type=?
	`, DOOP_MASTER)
	rows, err := doopdb.adapter.Query(statement, DOOP_DEFAULT_BRANCH, "logical_table")
	HandleError(err)

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

func (doopdb *DoopDb) Close() error {
	return doopdb.adapter.Close()
}

/*

					Branch Interface
	It provides an interface to interact on branches


*/
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
