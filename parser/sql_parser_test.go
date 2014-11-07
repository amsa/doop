package parser

import (
	"testing"

	"github.com/amsa/doop/common"
	"github.com/stretchr/testify/assert"
)

func TestBasicSelectRewrite(t *testing.T) {
	parser := MakeSqlParser()

	raw := "SELECT * FROM t1 WHERE t1.name == ?"
	rewriter := func(raw string) string {
		prefix := "branch"
		suffix := "logical"
		return common.ConcreteName(raw, prefix, suffix)
	}
	candidates := map[string]string{"t1": "sql1"}

	actual := parser.Rewrite(raw, rewriter, candidates)
	expected := "SELECT * FROM __branch_t1_logical WHERE __branch_t1_logical.name == ?"
	assert.Equal(t, expected, actual)
}

func TestComplexSelectRewrite(t *testing.T) {
	parser := MakeSqlParser()

	origin := `SELECT
			t2.id FROM t1 AS t , 
			t2 WHERE t1.name == t2.name
			GROUP BY t2.age`
	rewriter := func(raw string) string {
		prefix := "branch"
		suffix := "logical"
		return common.ConcreteName(raw, prefix, suffix)
	}
	target := map[string]string{"t1": "sql1", "t2": "sql2"}

	actual := parser.Rewrite(origin, rewriter, target)
	expected := `SELECT
			__branch_t2_logical.id FROM __branch_t1_logical AS t , 
			__branch_t2_logical WHERE __branch_t1_logical.name == __branch_t2_logical.name
			GROUP BY __branch_t2_logical.age`
	assert.Equal(t, expected, actual)
}

func TestSechemaRewrite(t *testing.T) {
	parser := MakeSqlParser()

	origin := `CREATE t1
			(
				c1 INTEGER, 
				c2 char(10)
			)`
	rewriter := func(raw string) string {
		prefix := "branch"
		suffix := "h"
		return common.ConcreteName(raw, prefix, suffix)
	}
	candidates := map[string]string{"t1": "sql1"}

	actual := parser.Rewrite(origin, rewriter, candidates)
	expected := `CREATE __branch_t1_h
			(
				c1 INTEGER, 
				c2 char(10)
			)`
	assert.Equal(t, expected, actual)
}

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
