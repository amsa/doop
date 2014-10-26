package core

/*
	file: doop.go
	Only this file contains APIs that are exported to doop-core user.
*/

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	DOOP_DIRNAME      = ".doop"
	DOOP_MAPPING_FILE = "doopm"
)

type Doop struct {
	homeDir string
}

func GetDoop() *Doop {
	return &Doop{}
}

/* Private methods:
************************************
 */

// getDoopDir returns the path to home directory of Doop
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

// getDbDir returns the path to the database directory inside the Doop home directory
func (doop *Doop) getDbDir(dbId string) string {
	return strings.Join([]string{doop.getDoopDir(), dbId}, string(os.PathSeparator))
}

// getDbMappingFile returns the path to the file that keeps track of mapping between database name and identifier (hash)
func (doop *Doop) getDbMappingFile() string {
	return strings.Join([]string{doop.getDoopDir(), DOOP_MAPPING_FILE}, string(os.PathSeparator))
}

// getDbId returns the identifier (hash) for the given database name
func (doop *Doop) getDbId(dbName string) (string, error) {
	file, err := os.OpenFile(doop.getDbMappingFile(), os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lArray := strings.Split(line, ",")
		if lArray[0] == dbName {
			return lArray[1], nil
		}
	}
	return "", errors.New("Database not found: " + dbName)
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string) (bool, error) {
	err := ioutil.WriteFile(doop.getDbMappingFile(), []byte(dbName+","+dbId+"\n"), 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}

// getDbIdMap returns the mapping of all the database names with their identifiers
func (doop *Doop) GetDbIdMap() (map[string]string, error) {
	m := make(map[string]string)
	file, err := os.OpenFile(doop.getDbMappingFile(), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lArray := strings.Split(line, ",")
		m[lArray[0]] = lArray[1]
	}
	return m, nil
}

/* Public methods:
************************************
 */
// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	dbId := generateDbId(dsn)
	dbDir := doop.getDbDir(dbId)
	if _, err := os.Stat(dbDir); err != nil {
		_, e := doop.setDbId(dbName, dbId)
		if e != nil {
			return false, e
		}
		os.Mkdir(dbDir, 0755)
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
