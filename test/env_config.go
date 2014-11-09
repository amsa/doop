package test

/*
use to setup test enviroment
*/

import (
	"fmt"
	"os/exec"
	"strings"
)

func SetupDb(dbName string) {
	setup := `
		CREATE TABLE t1 (pkey INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t2 (pkey INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t3 (pkey INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t4 (pkey INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		INSERT INTO t1 VALUES(1,1,'HEY');
		INSERT INTO t1 VALUES(2,1,'HEY');
		INSERT INTO t1 VALUES(3,1,'HEY');
		INSERT INTO t1 VALUES(4,1,'HEY');
		INSERT INTO t2 VALUES(1,2,'HEY');
		INSERT INTO t2 VALUES(2,2,'HEY');
		INSERT INTO t2 VALUES(3,2,'HEY');
		INSERT INTO t2 VALUES(4,2,'HEY');
		INSERT INTO t3 VALUES(1,3,'HEY');
		INSERT INTO t3 VALUES(2,3,'HEY');
		INSERT INTO t3 VALUES(3,3,'HEY');
		INSERT INTO t3 VALUES(4,3,'HEY');
		INSERT INTO t4 VALUES(1,4,'HEY');
		INSERT INTO t4 VALUES(2,4,'HEY');
		INSERT INTO t4 VALUES(3,4,'HEY');
		INSERT INTO t4 VALUES(4,4,'HEY');
		`
	cmd := exec.Command("sqlite3", dbName)
	cmd.Stdin = strings.NewReader(setup)
	err := cmd.Run()

	if err != nil {
		fmt.Println("Fail to setup test database...")
		panic(err)
	}
}

func CleanDb(dbName string) {
	cmd := exec.Command("rm", "-rf", dbName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Fail to remove test database...")
		panic(err)
	}
}
