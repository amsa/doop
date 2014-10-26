package core

/*
file: common.go

	Only this file contains APIs that are exported to doop-core user.

*/

import (
	"errors"
	"os"
)

//Package Level
type Doop struct {
	homeDir string
}

type DoopDb interface {
	//Methods of DoopDB interface
	getInfo()
	createBranch(branchName string, baseBranch string) (bool, error)
	removeBranch(branchName string) (bool, error)
	listBranches() (bool, error)
	executeSQL(sqlCommand string, branchName string) (*Result, error)
}

type Result struct {
	rawResult string
}

func GetDoop() *Doop {
	return &Doop{getDoopDir()}
}

/* Private methods:
************************************
 */
// getDbId returns the identifier (hash) for the given database name
func (doop *Doop) getDbId(dbName string) (string, error) {
	return ""
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string) (bool, error) {

}

/* Public methods:
************************************
 */
// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	dbId := generateDbId(dsn)
	dbDir := getDbDir(dbId)
	if _, err := os.Stat(dbDir); err != nil {
		os.Mkdir(dbDir, 0755)
		doop.setDbId(dbName, dbId)
		return true, nil
	}
	Debug("Database already initialized in %s", dbDir)
	return false, errors.New("Database exists!")
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
