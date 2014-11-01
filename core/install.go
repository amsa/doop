package core

import (
	"fmt"
	. "github.com/amsa/doop/common"
	"os"
)

func (doop *Doop) install() {
	if _, err := os.Stat(doop.homeDir); err == nil {
		fmt.Println("Doop is ready! You can create a new Doop project.")
		return
	}

	Debug("Doop directory does not exist. Creating doop directory at %s...", doop.homeDir)
	HandleError(os.Mkdir(doop.homeDir, 0755)) // Create Doop home directory

	// Create configuration file
	//defaultConfig := `[database]
	//DSN = ` + dsn + `
	//`
	//handleError(ioutil.WriteFile(doop.getConfigFile(), []byte(defaultConfig), 0644))

	// create database and needed tables
	//adapter := adapter.GetAdapter(dsn)
	//adapter.Exec(`CREATE TABLE doop_mapping (
	//id integer NOT NULL PRIMARY KEY,
	//dsn     text,
	//db_name text,
	//hash    text
	//);`)
	//adapter.Exec(`CREATE UNIQUE INDEX ON doop_mapping (hash);`)
	//adapter.Close()
}
