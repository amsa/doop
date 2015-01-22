package doop

/*
	file: doop.go
	Only this file contains exposed API
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
	DOOP_DIRNAME      = ".doop"
	DOOP_CONF_FILE    = "config"
	DOOP_MAPPING_FILE = "doopm"
)

var homeDir string
var dbInfoMap map[string]*DoopDbInfo

type DoopConfig struct {
	Database struct {
		DSN string
	}
}

func init() {
	initDoop()
}

/* Private methods:
************************************
 */

// initDoop initializes Doop home directory and creates it if it does not exist
func initDoop() {
	if homeDir != "" {
		return
	}
	currentUser, err := user.Current()
	HandleError(err)
	homeDir = strings.Join([]string{currentUser.HomeDir, DOOP_DIRNAME}, string(os.PathSeparator))

	if _, err := os.Stat(homeDir); err != nil {
		install()
	}
}

func install() {
	if _, err := os.Stat(homeDir); err == nil {
		fmt.Println("Doop is ready! You can create a new Doop project.")
		return
	}

	Debug("Doop directory does not exist. Creating doop directory at %s...", homeDir)
	HandleError(os.Mkdir(homeDir, 0755)) // Create Doop home directory
}

// getDbInfoByDbName retrieves database info for the given database (name)
func getDbInfoByDbName(dbName string) *DoopDbInfo {
	var info *DoopDbInfo = nil
	for _, dbInfo := range GetDbIdMap() {
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
func getDbMappingFile() string {
	return strings.Join([]string{homeDir, DOOP_MAPPING_FILE}, string(os.PathSeparator))
}

// setDbId returns the identifier (hash) for the given database name
func setDbId(dbName string, dbId string, dsn string) (bool, error) {
	file, err := os.OpenFile(getDbMappingFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	HandleError(err)
	defer file.Close()
	file.WriteString(dbId + "," + dbName + "," + dsn + "\n")
	return true, nil
}

// removeDbId removes database information from the mappin file
func removeDbId(dbId string) (bool, error) {
	newMap := ""
	delete(dbInfoMap, dbId)
	for _, v := range dbInfoMap {
		newMap += v.Hash + "," + v.Name + "," + v.DSN + "\n"
	}
	err := ioutil.WriteFile(getDbMappingFile(), []byte(newMap), 0644)
	return err == nil, err
}

/* Public methods:
************************************
 */
// getDbIdMap returns the mapping of all the database names with their identifiers
func GetDbIdMap() map[string]*DoopDbInfo {
	if dbInfoMap != nil {
		return dbInfoMap
	}
	dbInfoMap = make(map[string]*DoopDbInfo)
	file, err := os.OpenFile(getDbMappingFile(), os.O_CREATE|os.O_RDONLY, 0644)
	HandleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lArray := strings.Split(line, ",")
		dbInfoMap[lArray[0]] = &DoopDbInfo{Hash: lArray[0], Name: lArray[1], DSN: lArray[2]}
	}
	return dbInfoMap
}

// TrackDb initializes the database directory with a given identifier (hash)
func TrackDb(dbName string, dsn string) (bool, error) {
	dbId := GenerateDbId(dsn) // Generate the id (hash) from the given DSN
	if _, ok := GetDbIdMap()[dbId]; !ok {
		// Set mapping for the new database
		_, e := setDbId(dbName, dbId, dsn)
		if e != nil {
			return false, e
		}
		//initialize metadata in that database
		doopdb := MakeDoopDb(&DoopDbInfo{dsn, dbName, dbId})

		err := doopdb.Init()
		if err != nil {
			return false, errors.New("Fail to initialize the database: " + err.Error())
		} else {
			return true, nil
		}
	}
	return false, errors.New("Database already initialized!")
}

// ListDbs returns an array of all the tracked databases
func ListDbs() []string {
	var list []string
	for _, dbInfo := range GetDbIdMap() {
		list = append(list, dbInfo.Name)
	}
	return list
}

func GetDoopDb(dbName string) *DoopDb {
	info := getDbInfoByDbName(dbName)
	doopdb := MakeDoopDb(info)
	return doopdb
}

// UntrackDb untracks a database
func UntrackDb(dbName string) (bool, error) {
	info := getDbInfoByDbName(dbName)
	doopdb := MakeDoopDb(info)
	removeDbId(info.Hash)
	err := doopdb.Clean()
	if err != nil {
		return false, errors.New("Fail to clean metadata in the database")
	}
	return true, nil
}
