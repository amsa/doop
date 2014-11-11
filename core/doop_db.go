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
	"github.com/amsa/doop/parser"
)

const (
	DOOP_DEFAULT_BRANCH  = "master"
	DOOP_TABLE_BRANCH    = "__branch"
	DOOP_BRANCH_NAME_IDX = "__branch_name_idx"
	DOOP_MASTER          = "__doop"
	DOOP_TABLE_TYPE      = "logical_table"
	DOOP_SUFFIX_T        = "t"    //logical table, will be appended to the view of the logical table
	DOOP_SUFFIX_VD       = "vdel" //vertical deletion
	DOOP_SUFFIX_HD       = "hdel" //horizontal deletion
	DOOP_SUFFIX_V        = "v"    //v section
	DOOP_SUFFIX_H        = "h"    //h section
	DOOP_SUFFIX_VIEW     = "view" //view section
)

type DoopDbInfo struct {
	DSN  string
	Name string
	Hash string
}

type DoopDb struct {
	info         *DoopDbInfo
	adapter      adapter.Adapter
	branchTables map[string]map[string]string
	sql_parser   *parser.SqlParser
}

func MakeDoopDb(info *DoopDbInfo) *DoopDb {
	adpt := adapter.GetAdapter(info.DSN)
	parser := parser.MakeSqlParser()
	return &DoopDb{info, adpt, make(map[string]map[string]string, 16), parser}
}

/*

	DooDb Creation and Elimination

*/

func (doopdb *DoopDb) createMaster() error {
	//get tables first of all
	tables, err := doopdb.adapter.GetTableSchema()
	if err != nil {
		return err
	}

	// Create doop_master table to store metadata of all tables
	statement := fmt.Sprintf(`
			CREATE TABLE %s (
				id integer NOT NULL PRIMARY KEY,
				tbl_name text,
				tbl_type text,
				branch text,
				sql text
			);`, DOOP_MASTER)
	_, err = doopdb.adapter.Exec(statement)
	if err != nil {
		return errors.New("failed to create table: " + err.Error())
	}

	//Insert tables into doop_master
	i := 0
	for table, sql := range tables {
		i++
		statement = fmt.Sprintf(`
				INSERT INTO %s VALUES (
					%v,
					'%s',
					'%s',
					'%s',
					'%s'	
				)	
			`, DOOP_MASTER, i, table, DOOP_TABLE_TYPE, DOOP_DEFAULT_BRANCH, sql)
		_, err = doopdb.adapter.Exec(statement)
		if err != nil {
			msg := fmt.Sprintf("failed to insert table: %s", err.Error())
			return errors.New(msg)
		}
	}
	return nil
}

func (doopdb *DoopDb) destroyMaster() error {
	statement := fmt.Sprintf(`DROP TABLE %s`, DOOP_MASTER)
	_, err := doopdb.adapter.Exec(statement)
	return err
}

func (doopdb *DoopDb) createBranchTable() error {
	// Create branch table to store the branches
	_, err := doopdb.adapter.Exec(`CREATE TABLE ` + DOOP_TABLE_BRANCH + ` (
			id integer NOT NULL PRIMARY KEY,
			name     text,
			parent   text,
			metadata text
		);`)
	if err != nil {
		return err
	}
	//_, err = doopdb.adapter.Exec(`CREATE UNIQUE INDEX ` + DOOP_BRANCH_NAME_IDX + ` ON ` + DOOP_TABLE_BRANCH + ` (name);`)
	//if err != nil {
	//return err
	//}
	return nil
}

func (doopdb *DoopDb) destroyBranchTable() error {
	_, err := doopdb.adapter.Exec("DROP TABLE " + DOOP_TABLE_BRANCH)
	if err != nil {
		return err
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
	err := doopdb.createMaster()
	if err != nil {
		msg := fmt.Sprintf("failed to create master: %s", err.Error())
		return errors.New(msg)
	}

	// Create the branch management table
	err = doopdb.createBranchTable()
	if err != nil {
		msg := fmt.Sprintf("failed to create branch management: %s", err.Error())
		return errors.New(msg)
	}

	// Create default branch
	_, err = doopdb.CreateBranch(DOOP_DEFAULT_BRANCH, "")
	if err != nil {
		msg := fmt.Sprintf("failed to create default branch : %s", err.Error())
		return errors.New(msg)
	}

	return nil
}

func (doopdb *DoopDb) Clean() error {
	_, err := doopdb.adapter.Exec(`DROP INDEX IF EXISTS ` + DOOP_BRANCH_NAME_IDX)
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
func (doopdb *DoopDb) Query(branchName string, sql string, args ...interface{}) (*sql.Rows, error) {
	//hsec
	rewriter := func(origin string) string {
		prefix := branchName
		suffix := DOOP_SUFFIX_VIEW
		return ConcreteName(origin, prefix, suffix)
	}
	rewritten_sql := doopdb.sql_parser.Rewrite(sql, rewriter, doopdb.GetTableSchema(branchName))
	rows, err := doopdb.adapter.Query(rewritten_sql, args...)
	return rows, err
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

func (doopdb *DoopDb) GetAllTableSchema() (map[string]string, error) {
	return doopdb.adapter.GetTableSchema()
}

//argument is not supported yet
func (doopdb *DoopDb) GetTableSchema(branchName string) map[string]string {
	tables, ok := doopdb.branchTables[branchName]
	if ok {
		return tables
	}
	//find out the name of tables in default branch
	statement := fmt.Sprintf(`
		SELECT tbl_name, sql FROM %s WHERE branch=? AND tbl_type=?
	`, DOOP_MASTER)
	rows, err := doopdb.adapter.Query(statement, DOOP_DEFAULT_BRANCH, "logical_table")
	defer rows.Close()
	HandleError(err)

	ret := make(map[string]string)
	for rows.Next() {
		var name string
		var sql string
		err := rows.Scan(&name, &sql)
		HandleError(err)
		ret[name] = sql
	}
	doopdb.branchTables[branchName] = ret
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
	rt := make([]string, 0, 16)
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
	//TODO remove branch should check whether the branch is the parent
	//of other branches, if is, return false
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

	//get all tables in parent branch
	//if no parent branch(we are creating master branch), then just say parent branch is current branch
	if branchName == DOOP_DEFAULT_BRANCH {
		parentBranch = branchName
	}
	tables := doopdb.GetTableSchema(parentBranch)

	//create companian tables for each logical table
	for tableName, schema := range tables {
		//vdel
		vdel_name := ConcreteName(tableName, branchName, DOOP_SUFFIX_VD)
		vdel := fmt.Sprintf(`
			CREATE TABLE %s (
				c_name char(128)
			)	
		`, vdel_name)
		_, err := doopdb.adapter.Exec(vdel)
		if err != nil {
			msg := fmt.Sprintf("failed to create vdel: %s", err.Error())
			return false, errors.New(msg)
		}

		//hdel
		//TODO finish schema parsing, them complete this part
		// Currently it defaultly use integer as the type of the key
		hdel_name := ConcreteName(tableName, branchName, DOOP_SUFFIX_HD)
		hdel := fmt.Sprintf(`
			CREATE TABLE %s (
				key integer
			)
		`, hdel_name)
		_, err = doopdb.adapter.Exec(hdel)
		if err != nil {
			msg := fmt.Sprintf("failed to create hdel: %s", err.Error())
			return false, errors.New(msg)
		}

		//vsec
		//sql, err := sql_parser.Parse(schema)
		//TODO finish schema parsing, them complete this part
		//now it assume there is a column call "id" in all table
		vsec_name := ConcreteName(tableName, branchName, DOOP_SUFFIX_V)
		vsec := fmt.Sprintf(`
			CREATE TABLE %s AS SELECT id FROM %s
		`, vsec_name, tableName)
		_, err = doopdb.adapter.Exec(vsec)
		if err != nil {
			msg := fmt.Sprintf("failed to execute \n%s\n: %s", vsec, err.Error())
			return false, errors.New(msg)
		}

		//hsec
		rewriter := func(origin string) string {
			prefix := branchName
			suffix := DOOP_SUFFIX_H
			return ConcreteName(origin, prefix, suffix)
		}
		hsec := doopdb.sql_parser.Rewrite(schema, rewriter, tables)

		_, err = doopdb.adapter.Exec(hsec)

		if err != nil {
			msg := fmt.Sprintf("failed to create hsec: %s", err.Error())
			return false, errors.New(msg)
		}

		// Create logical views
		_, err = doopdb.CreateLogicalView(branchName, tableName)
		if err != nil {
			msg := fmt.Sprintf("failed to create logical views for the default branch : %s", err.Error())
			return false, errors.New(msg)
		}
	}
	return true, nil
}

func (doopdb *DoopDb) CreateLogicalView(branchName string, tableName string) (bool, error) {
	viewName := ConcreteName(tableName, branchName, DOOP_SUFFIX_VIEW)
	hsec := ConcreteName(tableName, branchName, DOOP_SUFFIX_H)

	// TODO: complete the view to incorporate all the elements including vdel, hdel and etc.
	viewCreationSql := fmt.Sprintf(`CREATE VIEW %s AS 
		SELECT * FROM %s UNION SELECT * FROM %s;`, viewName, tableName, hsec)
	_, err := doopdb.adapter.Exec(viewCreationSql)
	if err != nil {
		return false, err
	}
	return true, nil
}
