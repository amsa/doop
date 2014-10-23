package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amsa/doop-core/core"
)

var DEBUG bool = true

func help() {
	fmt.Println("doop [command] [options]")

	fmt.Println(`list of commands:
	init			initialize a new doop environment
	help			print this message`)

	fmt.Println(`list of options:
	-v			enable verbose mode`)
}

func initialize(args []string) {
	if len(args) == 0 {
		fmt.Println("DSN string is required to initialize doop environment. (e.g. sqlite://mydb.db)")
		return
	}
	dbId := core.GetDbId(args[0])
	doopDir := core.GetDoopDir()
	if _, err := os.Stat(doopDir); err != nil {
		debug("Doop dir does not exist. Creating doop directory at %s...", doopDir)
		os.Mkdir(doopDir, 0755)
	}
	debug("Database id: %s", dbId)
	// TODO: create the database directory inside doop dir (if not exists)
}

func debug(values ...interface{}) {
	if DEBUG {
		if len(values) == 0 {
			log.Println(values[0])
		} else {
			log.Printf(values[0].(string), values[1:]...)
		}
	}
}

func main() {
	var args []string
	for i := range os.Args {
		if os.Args[i] == "-v" { // enable verbose mode
			DEBUG = true
		} else {
			args = append(args, os.Args[i])
		}
	}
	if len(args) == 1 || args[1] == "help" {
		help()
		return
	}
	switch args[1] {
	case "init":
		fmt.Println("Initializing...")
		initialize(args[2:])
	default:
		fmt.Errorf("Invalid command: %s", os.Args[1])
	}

	//result := doop.Query("SELECT * FROM Towns;")
	//fmt.Println(result)
}
