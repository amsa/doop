package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/amsa/doop/adapter"
)

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
	dsn := args[0]
	defaultConfig := `[database]
DSN = ` + dsn + `
`
	handleError(ioutil.WriteFile(doop.getConfigFile(), []byte(defaultConfig), 0644))

	// create database and needed tables
	adapter := adapter.GetAdapter(dsn)
	adapter.Exec(`CREATE TABLE doop_mapping (
		id integer NOT NULL PRIMARY KEY, 
		dsn     text,
		db_name text,
		hash    text
    );`)
	adapter.Exec(`CREATE UNIQUE INDEX ON doop_mapping (hash);`)
	adapter.Close()
}
