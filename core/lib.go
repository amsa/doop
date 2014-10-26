package core

/*
file: common.go

	Only this file contains APIs that are exported to doop-core user.

*/

import (
	"errors"
	"log"
	"os"
	"os/user"
	"strings"
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

type Result interface {
	getRaw() string
}

func GetDoop() *Doop {
	return &Doop{}
}

/* Private methods:
************************************
 */

func (doop *Doop) getDoopDir() string {
	if doop.homeDir != "" {
		return doop.homeDir
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := strings.Join([]string{currentUser.HomeDir, DOOP_DIRNAME}, string(os.PathSeparator))
	if _, err := os.Stat(homeDir); err != nil {
		Debug("Doop dir does not exist. Creating doop directory at %s...", homeDir)
		os.Mkdir(homeDir, 0755)
	}
	doop.homeDir = homeDir
	return homeDir
}

// Returns the path to the database directory inside the Doop home directory
func (doop *Doop) getDbDir(dbId string) string {
	return strings.Join([]string{doop.homeDir, dbId}, string(os.PathSeparator))
}

// getDbId returns the identifier (hash) for the given database name
func (doop *Doop) getDbId(dbName string) string {
	return ""
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string) (bool, error) {
	return true, nil
	//file, err := os.Open(doop.homeDir)
}

/* Public methods:
************************************
 */
// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	dbId := generateDbId(dsn)
	dbDir := doop.getDbDir(dbId)
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
