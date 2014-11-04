package parser

import (
	"testing"

	"github.com/amsa/doop/common"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
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

func TestInsert(t *testing.T) {
	parser := MakeSqlParser()
	sql, err := parser.Parse(`INSERT INTO users VALUES ('Amin', 'Saeidi');`)
	common.HandleError(err)
	assert.Equal(t, "INSERT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "'Amin', 'Saeidi'", sql.Values)

	sql, err = parser.Parse(`INSERT INTO users (fname, lname) VALUES ('Amin', 'Saeidi');`)
	common.HandleError(err)
	assert.Equal(t, "INSERT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "fname, lname", sql.Columns)
	assert.Equal(t, "'Amin', 'Saeidi'", sql.Values)
}
