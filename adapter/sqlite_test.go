package adapter

/*
file: sqlite_test.go
description: test for sqlite adapter

*/
import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os/exec"
	"strings"
	"testing"
)

type SuiteTester struct {
	suite.Suite
	sqlite Adapter
	dsn    string
}

func (suite *SuiteTester) SetupSuite() {
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
	suite.dsn = "test_db"
	cmd := exec.Command("sqlite3", suite.dsn)
	cmd.Stdin = strings.NewReader(setup)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to initialized the test db: %s\n", err.Error())
	}

	suite.sqlite, err = MakeSQLite(suite.dsn)
	if err != nil {
		log.Fatalf("Test setup fail: %s\n", err.Error())
	}
	fmt.Printf("Setup test enviroment...\n")
}

func (suite *SuiteTester) TearDownSuite() {
	suite.sqlite.Close()
	cmd := exec.Command("rm", "-rf", suite.dsn)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to rm the test database: %s\n", err.Error())
	}
	fmt.Printf("Teardown test enviroment...\n")
}

func (suite *SuiteTester) TestGetTables() {
	tables, err := suite.sqlite.GetTables()
	if err != nil {
		log.Fatalf("failed to get tables: %s\n", err.Error())
	}
	assert.Equal(suite.T(), len(tables), 4)
}

func (suite *SuiteTester) TestGetSchema() {
	tables, err := suite.sqlite.GetTables()
	if err != nil {
		log.Fatalf("failed to get tables: %s\n", err.Error())
	}
	for _, table := range tables {
		_, err := suite.sqlite.GetSchema(table)
		assert.Nil(suite.T(), err)
	}
}

func (suite *SuiteTester) TestQuery() {
	tables, err := suite.sqlite.GetTables()
	if err != nil {
		assert.Nil(suite.T(), err)
	}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT * FROM %s", table)

		rows, err := suite.sqlite.Query(query)
		assert.Nil(suite.T(), err)

		columns, err := rows.Columns()
		assert.Nil(suite.T(), err)

		assert.Equal(suite.T(), len(columns), 3)

		next, err := rowToStrings(rows)
		assert.Nil(suite.T(), err)

		i := 0
		for true {
			row, err := next()
			if err != nil {
				log.Fatalf("invalid row")
			}
			if row == nil {
				break
			}
			i++
		}
		assert.Equal(suite.T(), i, 4)
		rows.Close()
	}
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suite.Run(t, suiteTester)
}
