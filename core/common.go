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
	DOOP_DIRNAME = ".doop"
)

func GetDoopDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join([]string{currentUser.HomeDir, DOOP_DIRNAME}, string(os.PathSeparator))
}

func GetDbId(dsn string) string {
	h := sha1.New()
	return fmt.Sprintf("%x", h.Sum([]byte(dsn)))
}
