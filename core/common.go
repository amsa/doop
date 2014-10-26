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
	getInfo() (string, error)
	createBranch(branchName string, baseBranch string) (bool, error)
	removeBranch(branchName string) (bool, error)
	listBranches() ([]string, error)
	executeSQL(sqlCommand string, branchName string) (*Result, error)
}

type Result struct {
	rawResult string
}

func GetDoop() *Doop {
	return &Doop{getDoopDir()}
}

// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	dbDir := getDbDir(generateDbId(dsn))
	if _, err := os.Stat(dbDir); err != nil {
		os.Mkdir(dbDir, 0755)
		// TODO: save the mapping of dbName and the db identifier (hash)
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
