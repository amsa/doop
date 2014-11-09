package core

/*
tests for doop_db.go
*/
import (
	//"fmt"
	"github.com/amsa/doop/common"
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

func (suite *SuiteTester) SetupTest() {
	suite.dsn = "sqlite://test_db"
	suite.db_path = "test_db"
	test.SetupDb(suite.db_path)
	db := MakeDoopDb(&DoopDbInfo{suite.dsn, suite.db_path, ""})
	suite.db = db
}

func (suite *SuiteTester) TearDownTest() {
	suite.db.Close()
	test.CleanDb(suite.db_path)
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suite.Run(t, suiteTester)
}

func (suite *SuiteTester) TestMasterTable() {
	all_tables, err := suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)
	err = suite.db.createMaster()
	assert.Nil(suite.T(), err)

	tables := suite.db.GetTableSchema("master") //the argument is not supported yet

	assert.NotEmpty(suite.T(), tables)

	for tableName, sql := range tables {
		assert.Equal(suite.T(), sql, all_tables[tableName])
	}
	err = suite.db.destroyMaster()
	assert.Nil(suite.T(), err)

	//test sucessfully destroy
	all_tables, err = suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)

	_, ok := all_tables[DOOP_MASTER]
	assert.False(suite.T(), ok)
}

func (suite *SuiteTester) TestBranchManagementTable() {
	tables, err := suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)
	_, ok := tables[DOOP_TABLE_BRANCH]
	assert.False(suite.T(), ok)
	//create the table
	err = suite.db.createBranchTable()
	assert.Nil(suite.T(), err)

	tables, err = suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)

	_, ok = tables[DOOP_TABLE_BRANCH]
	assert.True(suite.T(), ok)

	err = suite.db.destroyBranchTable()
	assert.Nil(suite.T(), err)

	//test sucessfully destroy
	tables, err = suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)

	_, ok = tables[DOOP_TABLE_BRANCH]
	assert.False(suite.T(), ok)
}

func (suite *SuiteTester) TestInit() {
	//original setting
	original_tables, err := suite.db.GetAllTableSchema()
	assert.NotEmpty(suite.T(), original_tables)
	assert.Nil(suite.T(), err)

	//init the db
	err = suite.db.Init()
	assert.Nil(suite.T(), err)

	//get logical table
	tables := suite.db.GetTableSchema("")
	assert.NotEmpty(suite.T(), tables)

	//shoud be same as original setting
	for t, sql := range original_tables {
		assert.Equal(suite.T(), sql, tables[t])
	}

	//get all tables
	all_tables, err := suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), all_tables)

	//should have correct number of tables
	//4 companian tables for each table, 1 doop_master table, 1 __branch table
	old_size := len(original_tables)
	expected := old_size + old_size*4 + 1 + 1
	actual := len(all_tables)
	assert.Equal(suite.T(), expected, actual)

	//check naming convention is applied
	suffixes := []string{DOOP_SUFFIX_H, DOOP_SUFFIX_V, DOOP_SUFFIX_HD, DOOP_SUFFIX_VD}
	for t, _ := range original_tables {
		for _, suffix := range suffixes {
			name := common.ConcreteName(t, DOOP_DEFAULT_BRANCH, suffix)
			_, ok := all_tables[name]
			assert.True(suite.T(), ok, "table %s is not correctly created", name)
		}
	}
}
func (suite *SuiteTester) TestInitAndCreateBranch() {

	new_branch := "new_branch"

	//original setting
	original_tables, err := suite.db.GetAllTableSchema()
	assert.NotEmpty(suite.T(), original_tables)
	assert.Nil(suite.T(), err)

	//init
	suite.db.Init()

	//create branch
	//should give error
	ok, err := suite.db.CreateBranch(new_branch, "")
	suite.False(ok)
	suite.NotNil(err)

	//should succeed
	ok, err = suite.db.CreateBranch(new_branch, "master")
	suite.True(ok)
	suite.Nil(err)

	//check number of tables
	//get all tables
	all_tables, err := suite.db.GetAllTableSchema()
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), all_tables)

	//should have correct number of tables
	//4 companian tables for each table, we have 2 branches so it's 8,
	//1 doop_master table, 1 __branch table
	old_size := len(original_tables)
	expected := old_size + old_size*8 + 1 + 1
	actual := len(all_tables)
	assert.Equal(suite.T(), expected, actual)

	//check naming convention is applied to all branches
	suffixes := []string{DOOP_SUFFIX_H, DOOP_SUFFIX_V, DOOP_SUFFIX_HD, DOOP_SUFFIX_VD}
	branches := []string{DOOP_DEFAULT_BRANCH, new_branch}
	for t, _ := range original_tables {
		for _, suffix := range suffixes {
			for _, br := range branches {
				name := common.ConcreteName(t, br, suffix)
				_, ok := all_tables[name]
				suite.True(ok, "table %s is not correctly created", name)
			}
		}
	}
}
