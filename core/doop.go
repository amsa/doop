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
	"strings"

	"code.google.com/p/gcfg"
	//_ "github.com/mattn/go-sqlite3"
)

const (
	DOOP_DIRNAME      = ".doop"
	DOOP_CONF_FILE    = "config"
	DOOP_MAPPING_FILE = "doopm"
)

type Doop struct {
	homeDir string
	config  DoopConfig
}

type DoopConfig struct {
	Database struct {
		DSN string
	}
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
	doop.homeDir = getDoopHome(DOOP_DIRNAME)

	if _, err := os.Stat(doop.homeDir); err != nil {
		fmt.Println("Doop environment is not installed yet.")
		os.Exit(1)
	}
	return doop.homeDir
}

func (doop *Doop) Install(args []string) {
	if len(args) == 0 {
		fmt.Println("In order to install Doop you need to pass a valid DSN. Doop will use this database to keep its own tables.")
		return
	}
	doop.homeDir = getDoopHome(DOOP_DIRNAME)
	if _, err := os.Stat(doop.homeDir); err == nil {
		fmt.Println("Doop is ready! You can create a new Doop project.")
		return
	}
	Debug("Doop dir does not exist. Creating doop directory at %s...", doop.homeDir)
	handleError(os.Mkdir(doop.homeDir, 0755)) // Create Doop home directory

	// Create configuration file
	dbFilePath := args[0]
	defaultConfig := `[database]
DSN = ` + dbFilePath + `
`
	handleError(ioutil.WriteFile(doop.getConfigFile(), []byte(defaultConfig), 0644))

	// TODO: create database based on the DSN
	//db, err := sql.Open("sqlite3", dbFilePath) // Create default database (sqlite)
	//handleError(err)
	//defer db.Close()
}

func (doop *Doop) getConfig() {
	var cfg DoopConfig
	handleError(gcfg.ReadFileInto(&cfg, doop.getConfigFile()))
	doop.config = cfg
}

func (doop *Doop) getConfigFile() string {
	return strings.Join([]string{doop.getDoopDir(), DOOP_CONF_FILE}, string(os.PathSeparator))
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
