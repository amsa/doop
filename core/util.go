package core

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	doop_dirname = ".doop"
)

var debug bool = true

func SetDebug(val bool) {
	debug = val
}

func Debug(values ...interface{}) {
	if debug {
		if len(values) == 0 {
			log.Println(values[0])
		} else {
			log.Printf(values[0].(string), values[1:]...)
		}
	}
}

func getDoopDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := strings.Join([]string{currentUser.HomeDir, doop_dirname}, string(os.PathSeparator))
	if _, err := os.Stat(homeDir); err != nil {
		Debug("Doop dir does not exist. Creating doop directory at %s...", homeDir)
		os.Mkdir(homeDir, 0755)
	}
	return homeDir
}

// Returns the path to the database directory inside the Doop home directory
func getDbDir(dbId string) string {
	doopHomeDir := getDoopDir()
	return strings.Join([]string{doopHomeDir, dbId}, string(os.PathSeparator))
}

// generaeDbId generates the unique identifier for the given DSN
func generateDbId(dsn string) string {
	h := sha1.New()
	return fmt.Sprintf("%x", h.Sum([]byte(dsn)))
}

// GetDbHash returns the identifier (hash) of the given alias
func getDBHash(alias string) string {
	return ""
}
