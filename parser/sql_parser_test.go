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
	assert.Equal(t, "WHERE id=2 AND username='amsa'", sql.Tail)
	//fmt.Printf("%#v\n", sql)

	sql, err = parser.Parse(`SELECT * FROM users ORDER BY fname`)
	common.HandleError(err)
	assert.Equal(t, "SELECT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "*", sql.Columns)
	assert.Equal(t, "ORDER BY fname", sql.Tail)

	sql, err = parser.Parse(`SELECT * FROM users WHERE id=2 AND username='amsa' ORDER BY id GROUP BY test`)
	common.HandleError(err)
	assert.Equal(t, "SELECT", sql.Op)
	assert.Equal(t, "users", sql.TblName)
	assert.Equal(t, "*", sql.Columns)
	assert.Equal(t, "WHERE id=2 AND username='amsa' ORDER BY id GROUP BY test", sql.Tail)
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
