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

	//_ "github.com/mattn/go-sqlite3"
)

const (
	DOOP_DIRNAME        = ".doop"
	DOOP_CONF_FILE      = "config"
	DOOP_MAPPING_FILE   = "doopm"
	DOOP_DEFAULT_BRANCH = "master"
)

type Doop struct {
	homeDir string
	//config  DoopConfig
	adapter adapter.Adapter
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
	handleError(err)
	doop.homeDir = strings.Join([]string{currentUser.HomeDir, DOOP_DIRNAME}, string(os.PathSeparator))

	if _, err := os.Stat(doop.homeDir); err != nil {
		doop.install()
	}
}

func (doop *Doop) initAdapter(dsn string) {
	doop.adapter = adapter.GetAdapter(dsn)
}

// initConfig loads and parses Doop configurations
//func (doop *Doop) initConfig() {
//if doop.config.Database.DSN != "" {
//return
//}
//handleError(gcfg.ReadFileInto(&doop.config, doop.getConfigFile()))
//}

//func (doop *Doop) getConfigFile() string {
//return strings.Join([]string{doop.homeDir, DOOP_CONF_FILE}, string(os.PathSeparator))
//}

func (doop *Doop) getDbMappingFile() string {
	return strings.Join([]string{doop.homeDir, DOOP_MAPPING_FILE}, string(os.PathSeparator))
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string, dsn string) (bool, error) {
	err := ioutil.WriteFile(doop.getDbMappingFile(), []byte(dbId+","+dbName+","+dsn+"\n"), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* Public methods:
************************************
 */
// getDbIdMap returns the mapping of all the database names with their identifiers
func (doop *Doop) GetDbIdMap() map[string]*DoopDbInfo {
	m := make(map[string]*DoopDbInfo)
	file, err := os.OpenFile(doop.getDbMappingFile(), os.O_CREATE|os.O_RDONLY, 0644)
	handleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lArray := strings.Split(line, ",")
		m[lArray[0]] = &DoopDbInfo{Name: lArray[1], DSN: lArray[2]}
	}
	return m
}

// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	doop.initAdapter(dsn)
	dbId := generateDbId(dsn) // Generate the id (hash) from the given DSN
	if _, ok := doop.GetDbIdMap()[dbId]; !ok {
		// Set mapping for the new database
		_, e := doop.setDbId(dbName, dbId, dsn)
		if e != nil {
			return false, e
		}
		// Create branch table to store the branches
		handleErrorAny(doop.adapter.Exec(`CREATE TABLE __branch (
			id integer NOT NULL PRIMARY KEY,
			name     text,
			parent   text,
			metadata text
		);`))
		handleErrorAny(doop.adapter.Exec(`CREATE UNIQUE INDEX __branch_name_idx ON __branch (name);`))

		// Create default branch
		_, e = doop.CreateBranch(DOOP_DEFAULT_BRANCH, "")
		if e != nil {
			return false, e
		}
		return true, nil
	}
	return false, errors.New("Database already initialized!")
}

func (doop *Doop) ListDbs() ([]string, error) {
	return nil, nil
}

func (doop *Doop) GetDb(dbName string) (*DoopDb, error) {
	return nil, nil
}

func (doop *Doop) UntrackDb(dbName string) (bool, error) {
	return false, nil
}

func (doop *Doop) CreateBranch(branchName string, parentBranch string) (bool, error) {
	if branchName != DOOP_DEFAULT_BRANCH && parentBranch == "" {
		return false, errors.New("Parent branch name is not specified.")
	}
	handleErrorAny(doop.adapter.
		Exec(fmt.Sprintf(`INSERT INTO __branch (name, parent, metadata) VALUES ('%s', '%s', '{}')`,
		branchName, parentBranch)))

	// TODO: create branch specific tables: __{branch}_vdel, __{branch}_hdel, __{branch}_v, __{branch}_h
	return true, nil
}

func (doop *Doop) RemoveBranch(branchName string) (bool, error) {
	return false, nil
}

func (doop *Doop) MergeBranch(to string, from string) (bool, error) {
	return false, nil
}

func (doop *Doop) ListBranches(dbName string) []string {
	mapping := doop.GetDbIdMap()
	dsn := ""
	for _, dbInfo := range mapping {
		if dbInfo.Name == dbName {
			dsn = dbInfo.DSN
			break
		}
	}
	if dsn == "" {
		fmt.Println("Could not find database: " + dbName)
		os.Exit(1)
	}
	doop.initAdapter(dsn)

	rt := make([]string, 1)
	rows, err := doop.adapter.Query(`SELECT name FROM __branch;`)
	handleError(err)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		rt = append(rt, name)
	}
	return rt
}
