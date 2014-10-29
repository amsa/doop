package core

/*
	file: doop.go
	Only this file contains APIs that are exported to doop-core user.
*/

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/amsa/doop/adapter"

	"code.google.com/p/gcfg"

	//_ "github.com/mattn/go-sqlite3"
)

const (
	DOOP_DIRNAME   = ".doop"
	DOOP_CONF_FILE = "config"
)

type Doop struct {
	homeDir string
	config  DoopConfig
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
	d.initConfig()
	d.adapter = adapter.GetAdapter(d.config.Database.DSN)
	return d
}

/* Private methods:
************************************
 */

// initDoopDir returns the path to home directory of Doop
func (doop *Doop) initDoopDir() string {
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

func (doop *Doop) initConfig() DoopConfig {
	if doop.config.Database.DSN != "" {
		return doop.config
	}
	handleError(gcfg.ReadFileInto(&doop.config, doop.getConfigFile()))
	return doop.config
}

func (doop *Doop) getConfigFile() string {
	return strings.Join([]string{doop.homeDir, DOOP_CONF_FILE}, string(os.PathSeparator))
}

// dbIdExists checks whether a database is being tracked or not
func (doop *Doop) dbIdExists(dbId string) bool {
	q := `SELECT id FROM doop_mapping WHERE hash = '%s';`
	rows, err := doop.adapter.Query(fmt.Sprintf(q, dbId))
	if err != nil {
		panic(err)
	}
	return rows.Next()
}

// setDbId returns the identifier (hash) for the given database name
func (doop *Doop) setDbId(dbName string, dbId string) (bool, error) {
	q := `INSERT INTO doop_mapping (dsn, db_name, hash) VALUES ('%s', '%s', '%s');`
	_, err := doop.adapter.Exec(fmt.Sprintf(q, doop.config.Database.DSN, dbName, dbId))
	if err != nil {
		return false, err
	}
	return true, nil
}

/* Public methods:
************************************
 */
// getDbIdMap returns the mapping of all the database names with their identifiers
func (doop *Doop) GetDbIdMap() (map[string]string, error) {
	m := make(map[string]string)
	rows, err := doop.adapter.Query(`SELECT db_name, hash FROM doop_mapping;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var hash string
		rows.Scan(&name, &hash)
		m[name] = hash
	}
	return m, nil
}

// TrackDb initializes the database directory with a given identifier (hash)
func (doop *Doop) TrackDb(dbName string, dsn string) (bool, error) {
	dbId := generateDbId(dsn)
	if !doop.dbIdExists(dbId) {
		_, e := doop.setDbId(dbName, dbId)
		if e != nil {
			return false, e
		}
		return true, nil
	}
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
