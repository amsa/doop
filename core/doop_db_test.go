package core

/*
tests for doop_db.go
*/
import (
	//"fmt"
	//"github.com/amsa/doop/common"
	"github.com/amsa/doop/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteTester struct {
	suite.Suite
	db      *DoopDb
	dsn     string
	db_path string
}

func (suite *SuiteTester) SetupSuite() {
	suite.dsn = "sqlite://test_db"
	suite.db_path = "test_db"
	test.SetupDb(suite.db_path)
	db := MakeDoopDb(&DoopDbInfo{suite.dsn, suite.db_path, ""})
	suite.db = db
}

func (suite *SuiteTester) TearDownSuite() {
	suite.db.Close()
	test.CleanDb(suite.db_path)
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suite.Run(t, suiteTester)
}

func (suite *SuiteTester) TestMasterTable() {
	all_tables, err := suite.db.GetAllTables()
	assert.Nil(suite.T(), err)
	err = suite.db.createDoopMaster()
	assert.Nil(suite.T(), err)

	tables := suite.db.GetTables("master") //the argument is not supported yet
	for tableName, sql := range tables {
		assert.Equal(suite.T(), sql, all_tables[tableName])
	}
}
