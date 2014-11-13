package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/amsa/doop/common"
	"github.com/amsa/doop/core"
)

func help() {
	fmt.Println("doop [command] [options]")

	fmt.Println(`list of commands:
	init			initialize a new Doop project
	list			list all the objects (databases/branches)
	run			    runs a given SQL expression on the specified branch
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
	if len(args) == 0 { // show the list of the databases
		fmt.Println("List of databases:")
		for _, db := range doop.ListDbs() {
			fmt.Println("    " + db)
		}
	} else { // show the list of branches for the given database
		fmt.Printf("List of branches for `%s`:", args[0])
		branches, err := doop.GetDoopDb(args[0]).ListBranches()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, branch := range branches {
			fmt.Println("    " + branch)
		}
	}
}

func run(doop *core.Doop, args []string) {
	if len(args) < 2 {
		fmt.Println("Too few arguments passed. usage: run <branch@db> <sql> (e.g. doop run mybranch@mydb SELECT * FROM user)")
		return
	}

	branchInfo := strings.Split(args[0], "@")
	if len(branchInfo) != 2 {
		fmt.Println("Invalid branch info format. It should be like branch@db")
		return
	}
	sql := strings.Join(args[1:], " ")
	sqlOp := strings.SplitN(sql, " ", 2)
	if strings.ToUpper(sqlOp[0]) == "SELECT" {
		fmt.Println(doop.GetDoopDb(branchInfo[1]).Query(branchInfo[0], sql))
	} else {
		fmt.Println(doop.GetDoopDb(branchInfo[1]).Exec(branchInfo[0], sql))
	}
}

func remove(doop *core.Doop, args []string) {
	if len(args) != 2 {
		fmt.Println("Too few arguments passed. usage: rm -d|-b <alias> (e.g. doop rm -d mydb)")
		return
	}
	if args[0] == "-d" {
		doop.UntrackDb(args[1])
		fmt.Println("Database removed from Doop: " + args[1])
	} else if args[0] == "-b" {
		// TODO: handle removing branch
	} else {
		fmt.Println("Invalid option passed to rm: " + args[0])
	}
}

func main() {
	doop := core.GetDoop()
	var args []string
	for i := range os.Args {
		if os.Args[i] == "-v" { // enable verbose mode
			common.SetDebug(true)
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
	case "rm":
		remove(doop, args[2:])
	case "run":
		run(doop, args[2:])
	default:
		fmt.Errorf("Invalid command: %s", os.Args[1])
	}

	//result := doop.Query("SELECT * FROM Towns;")
	//fmt.Println(result)
}
