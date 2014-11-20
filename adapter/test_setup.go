package adapter

/*
use to setup test enviroment
*/

import "fmt"

func SetupDb(db Adapter) {
	setup := `
		CREATE TABLE t1 (id INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t2 (id INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t3 (id INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
		CREATE TABLE t4 (id INTEGER PRIMARY KEY, c1 INTEGER, c2 varchar(32));
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

	_, err := db.Exec(setup)
	if err != nil {
		fmt.Println("Fail to setup test database: " + err.Error())
		panic(err)
	}
}

func CleanDb(db Adapter) {
	teardown := `
	DROP TABLE t1;
	DROP TABLE t2;
	DROP TABLE t3;
	DROP TABLE t4;
	`
	_, err := db.Exec(teardown)
	//cmd := exec.Command("rm", "-rf", dbName)
	//err := cmd.Run()
	if err != nil {
		fmt.Println("Fail to remove test database...")
		panic(err)
	}
}
