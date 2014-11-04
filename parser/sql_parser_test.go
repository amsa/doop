package parser

import (
	"testing"

	"github.com/amsa/doop/common"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	parser := MakeSqlParser()
	sql, err := parser.Parse(`SELECT * FROM users;`)
	common.HandleError(err)
	assert.Equal(t, "SELECT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "*", sql.Columns)
	//fmt.Printf("%#v\n", sql)

	sql, err = parser.Parse(`SELECT * FROM users WHERE id=2 AND username='amsa'`)
	common.HandleError(err)
	assert.Equal(t, "SELECT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "*", sql.Columns)
	assert.Equal(t, "id=2 AND username='amsa'", sql.Conditions)
	//fmt.Printf("%#v\n", sql)
}
