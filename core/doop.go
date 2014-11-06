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

	. "github.com/amsa/doop/common"

	//_ "github.com/mattn/go-sqlite3"
)

const (
	DOOP_DIRNAME        = ".doop"
	DOOP_CONF_FILE      = "config"
	DOOP_MAPPING_FILE   = "doopm"
	DOOP_DEFAULT_BRANCH = "master"
	DOOP_TABLE_BRANCH   = "__branch"
	DOOP_MASTER         = "__doop_master"
)

type Doop struct {
	homeDir string
	//config  DoopConfig
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
	dbId := GenerateDbId(dsn) // Generate the id (hash) from the given DSN
	if _, ok := doop.GetDbIdMap()[dbId]; !ok {
		// Set mapping for the new database
		_, e := doop.setDbId(dbName, dbId, dsn)
		if e != nil {
			return false, e
		}
		//initialize metadata in that database
		doopdb := MakeDoopDb(dsn)

		err := doopdb.Init()
		if err != nil {
			return false, errors.New("Fail to initialize the database")
		} else {
			return true, nil
		}
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
	doopdb := MakeDoopDb(info.DSN)
	doop.removeDbId(info.Hash)
	err := doopdb.Clean()
	if err != nil {
		return false, errors.New("Fail to clean metadata in the database")
	}
	return true, nil
}
