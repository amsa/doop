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

func getDoopDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join([]string{currentUser.HomeDir, doop_dirname}, string(os.PathSeparator))
}

func getDbId(dsn string) string {
	h := sha1.New()
	return fmt.Sprintf("%x", h.Sum([]byte(dsn)))
}

func getDBHash(alias string) string {
	//return the hash by given alias
	return ""
}
