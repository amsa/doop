package mysql

/*
file: adapter_test.go
description: test for adapter adapter

*/
import (
	"testing"

	"github.com/amsa/doop/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteTester struct {
	suite.Suite
	adapter adapter.Adapter
}

func (suite *SuiteTester) SetupSuite() {
	adapter.SetupDb(suite.adapter)
}

func (suite *SuiteTester) TearDownSuite() {
	adapter.CleanDb(suite.adapter)
}

func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	suiteTester.adapter = adapter.GetAdapter("mysql://root:@/doop_test")

	suite.Run(t, suiteTester)

	suiteTester.adapter.Close()
}

func (suite *SuiteTester) TestGetTableSchema() {
	tables, err := suite.adapter.GetTableSchema()
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), len(tables), 0)
}
