package adapter

/*
file: adapter_test.go
description: test for adapter adapter

*/
import (
	"fmt"
	"github.com/amsa/doop/common"
	"github.com/amsa/doop/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteTester struct {
	suite.Suite
	adapter Adapter
	dsn     string
}

func (suite *SuiteTester) SetupSuite() {
	suite.dsn = "test_db"
	test.SetupDb(suite.dsn)
	adpt, err := MakeSQLite(suite.dsn)
	assert.Nil(suite.T(), err)
	suite.adapter = adpt
}

func (suite *SuiteTester) TearDownSuite() {
	suite.adapter.Close()
	test.CleanDb(suite.dsn)
}

func (suite *SuiteTester) TestGetTables() {
	tables, err := suite.adapter.GetTables()
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), len(tables), 0)
}

func (suite *SuiteTester) TestGetSchema() {
	tables, err := suite.adapter.GetTables()
	assert.Nil(suite.T(), err)
	for _, table := range tables {
		_, err := suite.adapter.GetSchema(table)
		assert.Nil(suite.T(), err)
	}
}

/*
create a table, check number of tables, drop the table, check number of tables
*/
func (suite *SuiteTester) TestDDLTable() {
	statement := `CREATE TABLE testDDL ( 
					id INTEGER PRIMARY KEY, 
					c1 VARCHAR(10), 
					c2 VARCHAR(20) 
				  )`
	_, err := suite.adapter.Exec(statement)
	assert.Nil(suite.T(), err)

	tables, err := suite.adapter.GetTables()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(tables), 5)

	found := false
	for _, t := range tables {
		if t == "testDDL" {
			found = true
			break
		}
	}

	assert.True(suite.T(), found)

	statement = `DROP TABLE testDDL`
	_, err = suite.adapter.Exec(statement)
	assert.Nil(suite.T(), err)

	tables, err = suite.adapter.GetTables()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(tables), 4)
}

/*
add a column, check schema
*/
func (suite *SuiteTester) TestDDLSchema() {
	//sqlite3 only supports ADD COLUMN
}

func (suite *SuiteTester) TestQuery() {
	tables, err := suite.adapter.GetTables()
	if err != nil {
		assert.Nil(suite.T(), err)
	}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT * FROM %s", table)

		rows, err := suite.adapter.Query(query)
		assert.Nil(suite.T(), err)

		columns, err := rows.Columns()
		assert.Nil(suite.T(), err)

		assert.Equal(suite.T(), len(columns), 3)

		next, err := common.RowToStrings(rows)
		assert.Nil(suite.T(), err)

		for true {
			row, err := next()
			assert.Nil(suite.T(), err)
			if row == nil {
				break
			}
		}
		rows.Close()
	}
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suite.Run(t, suiteTester)
}
