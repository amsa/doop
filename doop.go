package main

import (
	"fmt"
	"os"

	"github.com/amsa/doop-core/core"
)

func help() {
	fmt.Println("doop [command] [options]")

	fmt.Println(`list of commands:
	init			initialize a new doop environment
	help			print this message`)

	fmt.Println(`list of options:
	-v			enable verbose mode`)
}

func initialize(doop *core.Doop, args []string) {
	if len(args) != 2 {
		fmt.Println("Too few arguments passed. usage: doop init <alias> <DSN> (e.g. doop init mydb sqlite:///usr/local/mydb.db)")
		return
	}
	_, err := doop.TrackDb(args[0], args[1])
	if err != nil {
		fmt.Println(err)
	}
	// TODO: create the database directory inside doop dir (if not exists)
}

func main() {
	doop := core.GetDoop()
	var args []string
	for i := range os.Args {
		if os.Args[i] == "-v" { // enable verbose mode
			core.SetDebug(true)
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
		initialize(doop, args[2:])
	default:
		fmt.Errorf("Invalid command: %s", os.Args[1])
	}

	//result := doop.Query("SELECT * FROM Towns;")
	//fmt.Println(result)
}
