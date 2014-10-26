package core

import (
	"crypto/sha1"
	"fmt"
	"log"
)

const (
	DOOP_DIRNAME      = ".doop"
	DOOP_MAPPING_FILE = "doopm"
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

// generaeDbId generates the unique identifier for the given DSN
func generateDbId(dsn string) string {
	h := sha1.New()
	return fmt.Sprintf("%x", h.Sum([]byte(dsn)))
}
