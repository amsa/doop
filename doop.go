package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/amsa/doop/common"
	"github.com/amsa/doop/doop"
	"github.com/olekukonko/tablewriter"
)

func help() {
	fmt.Println("doop [command] [options]")

	fmt.Println(`list of commands:
	init			initialize a new Doop project
	branch			create a new branch
	list			list all the objects (databases/branches)
	run			    runs a given SQL expression on the specified branch
	help			print this message`)

	fmt.Println(`list of options:
	-v			enable verbose mode`)
}

func branch(args []string) {
	if len(args) < 2 { // at least two arguments should be passed
		fmt.Println("Too few arguments passed. usage: doop branch <alias> <new> [<from>] (e.g. doop branch mydb newbranch)")
		return
	}

	parentBranch := "master"
	if len(args) == 3 {
		parentBranch = args[2]
	}
	_, err := doop.GetDoopDb(args[0]).CreateBranch(args[1], parentBranch)
	if err == nil {
		fmt.Printf("New branch '%s' has been created successfully.\n", args[1])
	} else {
		fmt.Println(err.Error())
	}
}

func initialize(args []string) {
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

func list(args []string) {
	if len(args) == 0 { // show the list of the databases
		fmt.Println("List of databases:")
		for _, db := range doop.ListDbs() {
			fmt.Println("    " + db)
		}
	} else { // show the list of branches for the given database
		fmt.Printf("List of branches for `%s`: \n", args[0])
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

func run(args []string) {
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
		results, err := doop.GetDoopDb(branchInfo[1]).Query(branchInfo[0], sql)
		if err != nil {
			fmt.Println(err)
		}
		cols, _ := results.Columns()

		// create table writer to write to stdout
		table := tablewriter.NewWriter(os.Stdout)

		// table header
		table.SetHeader(cols)

		// print the table contents
		next, _ := common.RowToStrings(results)
		row, _ := next()
		for row != nil {
			table.Append(row)
			row, _ = next()
		}
		results.Close()

		table.Render() // print out the table
	} else {
		results, err := doop.GetDoopDb(branchInfo[1]).Exec(branchInfo[0], sql)
		if err == nil {
			rowsAffected, _ := results.RowsAffected()
			fmt.Println("Number of affected rows: " + string(rowsAffected))

			lastId, _ := results.LastInsertId()
			fmt.Println("Last row ID: " + string(lastId))
		} else {
			fmt.Println(err.Error())
		}
	}
}

func remove(args []string) {
	if len(args) != 1 {
		fmt.Println("Too few arguments passed. usage: rm <[branch@]alias> (e.g. doop rm test@mydb)")
		return
	}
	branchInfo := strings.Split(args[0], "@")
	if len(branchInfo) == 1 {
		doop.UntrackDb(branchInfo[0])
		fmt.Println("Database removed from Doop: " + branchInfo[0])
	} else {
		// TODO: handle removing branch
	}
}

func main() {
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
		initialize(args[2:])
	case "list":
		list(args[2:])
	case "branch":
		branch(args[2:])
	case "rm":
		remove(args[2:])
	case "run":
		run(args[2:])
	default:
		fmt.Errorf("Invalid command: %s", os.Args[1])
	}
}
