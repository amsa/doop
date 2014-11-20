/*
file: common_test.go
description: test for common.go
*/

package common

import (
	//"fmt"

	"testing"

	"github.com/amsa/doop/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteTester struct {
	suite.Suite
	adapter adapter.Adapter
	dsn     string
	dbPath  string
}

func (suite *SuiteTester) SetupSuite() {
	adapter.SetupDb(suite.adapter)
}

func (suite *SuiteTester) TearDownSuite() {
	adapter.CleanDb(suite.adapter)
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suiteTester.adapter = adapter.GetAdapter("sqlite://test_db")

	suite.Run(t, suiteTester)

	suiteTester.adapter.Close()
	suiteTester.adapter.DropDb()
}

func (suite *SuiteTester) TestAddPrefix() {
	origin := "table"
	prefix := "branch1"
	expect := "branch1_table"
	result := AddPrefix(origin, prefix)

	assert.Equal(suite.T(), result, expect)
}

func (suite *SuiteTester) TestAddSuffix() {
	origin := "table"
	suffix := "v"
	expect := "table_v"
	result := AddSuffix(origin, suffix)

	assert.Equal(suite.T(), result, expect)
}

func (suite *SuiteTester) TestConcreteName() {
	origin := "table"
	prefix := "branch1"
	suffix := "v"
	expect := "branch1_table_v"
	result := ConcreteName(origin, prefix, suffix)

	assert.Equal(suite.T(), result, expect)
}

func (suite *SuiteTester) TestRowToStrings() {
	statement := "SELECT * FROM t1"
	rows1, err := suite.adapter.Query(statement)
	assert.Nil(suite.T(), err)
	rows2, err := suite.adapter.Query(statement)
	assert.Nil(suite.T(), err)

	next, err := RowToStrings(rows1)
	assert.Nil(suite.T(), err)

	rows2.Next()
	for true {
		rows2.Next()
		r, err := next()
		assert.Nil(suite.T(), err)
		if r == nil {
			break
		}
	}
	rows1.Close()
	rows2.Close()

}
