package core

import (
	"errors"
	"os"
)

//file: common.go
//only this file contains APIs that are exported to doop-core user.

//Package Level
type Doop struct {
	homeDir string
}

type DoopDb struct {
	dbName string
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
		return true, nil
	}
	Debug("Database already initialized in %s...", dbDir)
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

//Methods of DoopDB
func (db *DoopDb) createBranch(branchName string, baseBranch string) (bool, error) {
	return false, nil
}

func (db *DoopDb) removeBranch(branchName string) (bool, error) {
	return false, nil
}

func (db *DoopDb) listBranches() (bool, error) {
	return false, nil
}

func (db *DoopDb) executeSQL(sqlCommand string, branchName string) (*Result, error) {
	return nil, nil
}
