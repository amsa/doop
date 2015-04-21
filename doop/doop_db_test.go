package doop

/*
tests for doop_db.go
*/
import (
	"fmt"

	"strings"
	"testing"

	"github.com/amsa/doop/adapter"
	"github.com/amsa/doop/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteTester struct {
	suite.Suite
	db      *DoopDb
	adapter adapter.Adapter
	dbName  string
}

func (suite *SuiteTester) SetupTest() {
	suite.dbName = "test_db"
	suite.adapter = adapter.GetAdapter("sqlite://" + suite.dbName)
	suite.db = MakeDoopDb(&DoopDbInfo{"sqlite://" + suite.dbName, suite.dbName, ""})
	adapter.SetupDb(suite.adapter)
}

func (suite *SuiteTester) TearDownTest() {
	suite.adapter.Close()
	suite.adapter.DropDb()
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

	tables, err := suite.db.GetTableSchema(DOOP_DEFAULT_BRANCH) //the argument is not supported yet
	suite.Nil(err)
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
	tables, err := suite.db.GetTableSchema(DOOP_DEFAULT_BRANCH)
	suite.Nil(err)
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
	//1 __doop table, 1 __branch table
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

	// check views
	query := "SELECT name, sql FROM sqlite_master WHERE type='view';"
	rows, err := suite.db.adapter.Query(query)
	views := 0
	for rows.Next() {
		views++
	}
	// we have 4 tables and 2 branches
	assert.Equal(suite.T(), 4*2, views)

	//check query on view
	//since no modification on the tables yet, queries on view and the original table should return same values
	//this test involves tests on Query.
	for _, br := range branches {
		for t, _ := range original_tables {
			query := fmt.Sprintf("SELECT * FROM %s", t)
			rows_br, err := suite.db.Query(br, query)
			rows_default, err := suite.db.adapter.Query(query)

			defer rows_br.Close()
			defer rows_default.Close()

			//common.RowToStrings turns a row into an array of strings
			rows_br_next, err := common.RowToStrings(rows_br)
			suite.Nil(err)
			rows_default_next, err := common.RowToStrings(rows_default)
			suite.Nil(err)

			//compare all rows
			for true {
				r1, e1 := rows_br_next()
				r2, e2 := rows_default_next()
				suite.Nil(e1)
				suite.Nil(e2)
				if r1 == nil && r2 == nil {
					break
				} else if r1 == nil || r2 == nil {
					suite.Fail("views and the raw tables should have same number of rows.")
					break
				}
				s1 := strings.Join(r1, " ")
				s2 := strings.Join(r2, " ")
				suite.Equal(s2, s1, "%s \n\t!= %s", s1, s2)
			}
		}
	}
}

func (suite *SuiteTester) TestListBranches() {
	new_branch := "new_branch"

	//original setting
	original_tables, err := suite.db.GetAllTableSchema()
	assert.NotEmpty(suite.T(), original_tables)
	assert.Nil(suite.T(), err)

	//init
	suite.db.Init()

	//should have 1 branch
	branches, err := suite.db.ListBranches()
	suite.Nil(err)
	suite.Equal(1, len(branches))
	suite.Equal(DOOP_DEFAULT_BRANCH, branches[0])

	//create branch
	//should give error
	ok, err := suite.db.CreateBranch(new_branch, "")
	suite.False(ok)
	suite.NotNil(err)

	//should succeed
	ok, err = suite.db.CreateBranch(new_branch, "master")
	suite.True(ok)
	suite.Nil(err)

	//should have 2 branches
	branches, err = suite.db.ListBranches()
	suite.Nil(err)
	suite.Equal(2, len(branches))
	suite.Equal(new_branch, branches[1])
}

func (suite *SuiteTester) TestQueryExecBasic() {
	//original setting
	original_tables, err := suite.db.GetAllTableSchema()
	assert.NotEmpty(suite.T(), original_tables)
	assert.Nil(suite.T(), err)

	//init
	suite.db.Init()

	//test error handling
	for t, _ := range original_tables {
		//it's query a non-exist branch, it shoud give an error indicating the table does not exist
		statement := fmt.Sprintf("SELECT * FROM %s", t)
		rows, err := suite.db.Query("non_exist", statement)
		suite.NotNil(err)
		if rows != nil {
			rows.Close()
		}

		//INSERT should give an error
		statement = fmt.Sprintf("INSERT INTO %s VALUES (AA, BB, CC)") //values do not matter here
		rows, err = suite.db.Query(DOOP_DEFAULT_BRANCH, statement)
		suite.NotNil(err)
		if rows != nil {
			rows.Close()
		}
	}

	//create more branch
	branch_name := "branch1"
	suite.db.CreateBranch(branch_name, DOOP_DEFAULT_BRANCH)
	//Insert new rows into each branch
	suite.db.Exec(branch_name, "INSERT INTO t1 VALUES(1827, 8718, 'test branch1')")
	suite.db.Exec(DOOP_DEFAULT_BRANCH, "INSERT INTO t1 VALUES(1927, 7718, 'test master')")

	q := `select * from t1 where id1 = 1827` //shold only succeed in new branch

	rows1, err := suite.db.Query(DOOP_DEFAULT_BRANCH, q)
	defer rows1.Close()
	suite.Nil(err)
	rows2, err := suite.db.Query(branch_name, q)
	defer rows2.Close()
	suite.Nil(err)

	rows1_itr, err := common.RowToStrings(rows1)
	suite.Nil(err)
	suite.Nil(rows1_itr())

	rows2_itr, err := common.RowToStrings(rows2)
	suite.Nil(err)
	suite.NotNil(rows2_itr())
	suite.Nil(rows2_itr())

	q = `select * from t1 where id1 = 1927` //shold only succeed in DEFAULT_BRANCH
	rows3, err := suite.db.Query(DOOP_DEFAULT_BRANCH, q)
	suite.Nil(err)

	rows3_itr, err := common.RowToStrings(rows3)
	suite.NotNil(rows3_itr())
	suite.Nil(rows3_itr())
}

func (suite *SuiteTester) TestQueryExec() {
	suite.True(false)
}

func (suite *SuiteTester) TestExecErrorHandling() {
	tables, _ := suite.db.GetAllTableSchema()
	suite.db.Init()

	for t, _ := range tables {
		//select is not valid operation for Exec
		statement := fmt.Sprintf("SELECT * FROM %s", t)
		_, err := suite.db.Exec("master", statement)
		suite.NotNil(err, "failed to exec: %s", statement)

		statement = fmt.Sprintf("INSERT INTO %s VALUES (8, 10, 'baz')", t)
		_, err = suite.db.Exec("master", statement)
		suite.Nil(err, "failed to exec: %s", statement)
	}
}

func (suite *SuiteTester) TestGetParentBranch() {
	//init db
	suite.db.Init()

	//test default branch creation
	parentBranch, err := suite.db.getParentBranch(DOOP_DEFAULT_BRANCH)
	if err != nil {
		suite.Fail(err.Error())
		return
	}
	suite.Equal(DOOP_NULL_BRANCH, parentBranch)

	//new branch
	branchName := "newBranch"
	ok, err := suite.db.CreateBranch(branchName, DOOP_DEFAULT_BRANCH)
	if err != nil {
		suite.Fail(err.Error())
	}
	suite.True(ok)
	parentBranch, err = suite.db.getParentBranch(branchName)
	if err != nil {
		suite.Fail(err.Error())
		return
	}
	suite.Equal(DOOP_DEFAULT_BRANCH, parentBranch)

	//newer branch
	branchNameAgain := "newerBranch"
	ok, err = suite.db.CreateBranch(branchNameAgain, branchName)
	if err != nil {
		suite.Fail(err.Error())
	}
	suite.True(ok)
	parentBranch, err = suite.db.getParentBranch(branchNameAgain)
	if err != nil {
		suite.Fail(err.Error())
		return
	}
	suite.Equal(branchName, parentBranch)
}

func (suite *SuiteTester) TestRemoveBranch() {
	suite.True(false)
}
