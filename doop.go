package main

import (
	"fmt"
	"os"

	"github.com/amsa/doop/core"
)

func help() {
	fmt.Println("doop [command] [options]")

	fmt.Println(`list of commands:
	init			initialize a new Doop project
	list			list all the objects (databases/branches)
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
	if err == nil {
		fmt.Println("Initialized new Doop environment for " + args[0])
	} else {
		fmt.Println("Error while initializing: " + err.Error())
	}
}

func list(doop *core.Doop, args []string) {
	m := doop.GetDbIdMap()
	if len(m) == 0 {
		return
	}
	if len(args) == 0 { // show the list of the databases
		fmt.Println("List of databases:")
		for _, dbInfo := range m {
			fmt.Println("    " + dbInfo.Name)
		}
	} else { // show the list of branches for the given database
		fmt.Printf("List of branches for `%s`:", args[0])
		for _, branch := range doop.ListBranches(args[0]) {
			fmt.Println("    " + branch)
		}
	}
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
		initialize(doop, args[2:])
	case "list":
		list(doop, args[2:])
	default:
		fmt.Errorf("Invalid command: %s", os.Args[1])
	}

	//result := doop.Query("SELECT * FROM Towns;")
	//fmt.Println(result)
}
