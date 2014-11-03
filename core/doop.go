package core

/*
	file: doop.go
	Only this file contains APIs that are exported to doop-core user.
*/

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/amsa/doop/adapter"
	. "github.com/amsa/doop/common"

	//_ "github.com/mattn/go-sqlite3"
)

const (
	DOOP_DIRNAME        = ".doop"
	DOOP_CONF_FILE      = "config"
	DOOP_MAPPING_FILE   = "doopm"
	DOOP_DEFAULT_BRANCH = "master"

	DOOP_TABLE_BRANCH = "__branch"
	DOOP_MASTER       = "__doop_master"
)

type Doop struct {
	homeDir string
	//config  DoopConfig
	adapter   adapter.Adapter
	dbInfoMap map[string]*DoopDbInfo
}

type DoopConfig struct {
	Database struct {
		DSN string
	}
}

func GetDoop() *Doop {
	d := &Doop{}
	// initialize Doop object
	d.initDoopDir()
	//d.initConfig()
	//d.adapter = adapter.GetAdapter(d.config.Database.DSN)
	return d
}

/* Private methods:
************************************
 */

// initDoopDir initializes Doop home directory and creates it if it does not exist
func (doop *Doop) initDoopDir() {
	if doop.homeDir != "" {
		return
	}
	currentUser, err := user.Current()
	HandleError(err)
	doop.homeDir = strings.Join([]string{currentUser.HomeDir, DOOP_DIRNAME}, string(os.PathSeparator))

	if _, err := os.Stat(doop.homeDir); err != nil {
		doop.install()
	}
}

// initAdapter initializes Doop database adapter based on the given DSN string
func (doop *Doop) initAdapter(dsn string) {
	doop.adapter = adapter.GetAdapter(dsn)
}

//initConfig loads and parses Doop configurations
/*func (doop *Doop) initConfig() {
	if doop.config.Database.DSN != "" {
	return
	}
	handleError(gcfg.ReadFileInto(&doop.config, doop.getConfigFile()))
}

func (doop *Doop) getConfigFile() string {
	return strings.Join([]string{doop.homeDir, DOOP_CONF_FILE}, string(os.PathSeparator))
}*/

func (doop *Doop) getLogicalTables() map[string]string {
	//find out the name of tables in default branch
	statement := fmt.Sprintf(`
		SELECT name, sql FROM %s WHERE branch=? AND type=?
	`)
	rows, err := doop.adapter.Query(statement, DOOP_MASTER, "logical_table")
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

// getDbInfoByDbName retrieves database info for the given database (name)
func (doop *Doop) getDbInfoByDbName(dbName string) *DoopDbInfo {
	var info *DoopDbInfo = nil
	for _, dbInfo := range doop.GetDbIdMap() {
		if dbInfo.Name == dbName {
			info = dbInfo
			break
		}
	}
	if info == nil {
		fmt.Println("Could not find database: " + dbName)
		os.Exit(1)
	}

	return info
}

// getDbMappingFile returns the path to Doop mapping file
func (doop *Doop) getDbMappingFile() string {
	return strings.Join([]string{doop.homeDir, DOOP_MAPPING_FILE}, string(os.PathSeparator))
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string, dsn string) (bool, error) {
	file, err := os.OpenFile(doop.getDbMappingFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	HandleError(err)
	defer file.Close()
	file.WriteString(dbId + "," + dbName + "," + dsn + "\n")
	return true, nil
}

// removeDbId removes database information from the mappin file
func (doop *Doop) removeDbId(dbId string) (bool, error) {
	newMap := ""
	delete(doop.dbInfoMap, dbId)
	for _, v := range doop.dbInfoMap {
		newMap += v.Hash + "," + v.Name + "," + v.DSN + "\n"
	}
	err := ioutil.WriteFile(doop.getDbMappingFile(), []byte(newMap), 0644)
	return err == nil, err
}

/* Public methods:
************************************
 */
// getDbIdMap returns the mapping of all the database names with their identifiers
func (doop *Doop) GetDbIdMap() map[string]*DoopDbInfo {
	if doop.dbInfoMap != nil {
		return doop.dbInfoMap
	}
	doop.dbInfoMap = make(map[string]*DoopDbInfo)
	file, err := os.OpenFile(doop.getDbMappingFile(), os.O_CREATE|os.O_RDONLY, 0644)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lArray := strings.Split(line, ",")
		doop.dbInfoMap[lArray[0]] = &DoopDbInfo{Hash: lArray[0], Name: lArray[1], DSN: lArray[2]}
	}
	return doop.dbInfoMap
}

// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	doop.initAdapter(dsn)
	dbId := GenerateDbId(dsn) // Generate the id (hash) from the given DSN
	if _, ok := doop.GetDbIdMap()[dbId]; !ok {
		// Set mapping for the new database
		_, e := doop.setDbId(dbName, dbId, dsn)
		if e != nil {
			return false, e
		}
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
		HandleErrorAny(doop.adapter.Exec(statement))

		//Insert tables into doop_master
		tables, err := doop.adapter.GetTables()
		HandleErrorAny(tables, err)

		for i, table := range tables {
			schema, err := doop.adapter.GetSchema(table)
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
		HandleErrorAny(doop.adapter.Exec(`CREATE TABLE ` + DOOP_TABLE_BRANCH + ` (
			id integer NOT NULL PRIMARY KEY,
			name     text,
			parent   text,
			metadata text
		);`))
		HandleErrorAny(doop.adapter.Exec(`CREATE UNIQUE INDEX __branch_name_idx ON ` + DOOP_TABLE_BRANCH + ` (name);`))

		// Create default branch
		_, e = doop.CreateBranch(DOOP_DEFAULT_BRANCH, "")
		if e != nil {
			return false, e
		}
		return true, nil
	}
	return false, errors.New("Database already initialized!")
}

// ListDbs returns an array of all the tracked databases
func (doop *Doop) ListDbs() []string {
	var list []string
	for _, dbInfo := range doop.GetDbIdMap() {
		list = append(list, dbInfo.Name)
	}
	return list
}

// UntrackDb untracks a database
func (doop *Doop) UntrackDb(dbName string) (bool, error) {
	info := doop.getDbInfoByDbName(dbName)
	doop.initAdapter(info.DSN)
	doop.removeDbId(info.Hash)
	HandleErrorAny(doop.adapter.Exec(`DROP INDEX IF EXISTS __branch_name_idx;`))
	HandleErrorAny(doop.adapter.Exec(`DROP TABLE ` + DOOP_TABLE_BRANCH + `;`))
	return false, nil
}

// CreateBranch creates a new branch of the database forking from the given parent branch
func (doop *Doop) CreateBranch(branchName string, parentBranch string) (bool, error) {
	if branchName != DOOP_DEFAULT_BRANCH && parentBranch == "" {
		return false, errors.New("Parent branch name is not specified.")
	}
	HandleErrorAny(doop.adapter.
		Exec(`INSERT INTO `+DOOP_TABLE_BRANCH+` (name, parent, metadata) VALUES (?, ?, '{}')`, branchName, parentBranch))

	//get all tables
	tables := doop.getLogicalTables()

	//create companion tables for each logical table
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

// RemoveBranch deletes a branch
func (doop *Doop) RemoveBranch(branchName string) (bool, error) {
	return false, nil
}

// MergeBranch merges two branches into the first one (from)
func (doop *Doop) MergeBranch(from string, to string) (bool, error) {
	return false, nil
}

// ListBranches returns the list of all the branches for the given database
func (doop *Doop) ListBranches(dbName string) []string {
	doop.initAdapter(doop.getDbInfoByDbName(dbName).DSN)

	rt := make([]string, 1)
	rows, err := doop.adapter.Query(`SELECT name FROM ` + DOOP_TABLE_BRANCH + `;`)
	HandleError(err)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		rt = append(rt, name)
	}
	return rt
}
