package adapter

/*
file: adapter_test.go
description: test for adapter adapter

*/
import (
	"fmt"
	"testing"

	"github.com/amsa/doop/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteTester struct {
	suite.Suite
	adapter Adapter
}

func (suite *SuiteTester) SetupSuite() {
	SetupDb(suite.adapter)
}

func (suite *SuiteTester) TearDownSuite() {
	CleanDb(suite.adapter)
}

func (suite *SuiteTester) TestGetTableSchema() {
	tables, err := suite.adapter.GetTableSchema()
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), len(tables), 0)
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

	tables, err := suite.adapter.GetTableSchema()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(tables), 5)

	found := false
	for t, _ := range tables {
		if t == "testDDL" {
			found = true
			break
		}
	}

	assert.True(suite.T(), found)

	statement = `DROP TABLE testDDL`
	_, err = suite.adapter.Exec(statement)
	assert.Nil(suite.T(), err)

	tables, err = suite.adapter.GetTableSchema()
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
	tables, err := suite.adapter.GetTableSchema()
	if err != nil {
		assert.Nil(suite.T(), err)
	}
	for table, _ := range tables {
		query := fmt.Sprintf("SELECT * FROM %s", table)

		rows, err := suite.adapter.Query(query)
		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), rows)

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
	suiteTester.adapter = GetAdapter("sqlite://test_db")

	suite.Run(t, suiteTester)

	suiteTester.adapter.Close()
	suiteTester.adapter.DropDb()
}
