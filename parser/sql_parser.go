/*
A simple parser to parse SQL
*/

package parser

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

const (
	SQL_TYPE_DDL = 1
	SQL_TYPE_DML = 2
)

type Sql struct {
	Raw     string // Raw sql query
	Pos     int    // Current cursor position
	Type    int    // DDL or DML
	Op      string // SELECT or INSERT or ...
	TblName string // Table name
	Columns string // Columns (e.g. in SELECT specified before FROM)
	Values  string // Values used in INSERT or UPDATE
	Tail    string // The rest of the query (which we don't care about)
}

type SqlParser struct {
}

type Rewriter func(string) string

func MakeSqlParser() *SqlParser {
	return new(SqlParser)
}

func (sqlParser *SqlParser) Rewrite(origin string, rewriter Rewriter, target map[string]string) string {
	var buffer bytes.Buffer
	var lastToken bytes.Buffer
	var lastChar rune
	for _, c := range origin {
		//only keep one space
		if c == ' ' || c == '.' || c == ',' {
			if lastChar != ' ' {
				word := lastToken.String()
				if _, ok := target[word]; ok {
					word = rewriter(word)
				}
				buffer.WriteString(word)
				buffer.WriteRune(c)
				lastChar = c
				lastToken.Reset()
			}
		} else {
			lastToken.WriteRune(c)
			lastChar = c
		}
	}
	buffer.WriteString(lastToken.String())
	return buffer.String()
}

func (sqlParser *SqlParser) Parse(query string) (*Sql, error) {
	op := strings.ToUpper(strings.Trim(strings.SplitAfterN(query, " ", 2)[0], " "))
	switch op {
	case "SELECT":
		return sqlParser.parseSelect(query)
	case "INSERT":
		return sqlParser.parseInsert(query)
	default:
		return nil, errors.New("Invalid SQL operation: " + op)
	}
}

func (sqlParser *SqlParser) parseSelect(query string) (*Sql, error) {
	//selectRe := regexp.MustCompile(`(?i)(SELECT)(?:\s+DISTINCT)?\s+(?:INTO\s+\w+\s+)?([A-Za-z,*_]+)\s+FROM\s+(\w+)(?:\s+WHERE\s+(.+))?`)
	selectRe := regexp.MustCompile(`(?i)(SELECT)(?:\s+DISTINCT)?\s+(?:INTO\s+\w+\s+)?([A-Za-z,*_]+)\s+FROM\s+(\w+)\s*(.*)?`)
	matches := selectRe.FindStringSubmatch(query)
	if matches == nil {
		return nil, errors.New("Query did not match select pattern: " + query)
	}

	sql := new(Sql)
	sql.Raw = query
	sql.Type = SQL_TYPE_DML
	sql.Op = matches[1]
	sql.Columns = matches[2]
	sql.TblName = matches[3]
	sql.Tail = matches[4]
	return sql, nil
}

func (sqlParser *SqlParser) parseInsert(query string) (*Sql, error) {
	selectRe := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)(?:\s+\(([^)]+)\))?\s+VALUES\s*\(([^)]+)\)`)
	matches := selectRe.FindStringSubmatch(query)
	if matches == nil {
		return nil, errors.New("Query did not match INSERT pattern: " + query)
	}
	sql := new(Sql)
	sql.Raw = query
	sql.Type = SQL_TYPE_DML
	sql.Op = "INSERT"
	sql.TblName = matches[1]
	sql.Columns = matches[2]
	sql.Values = matches[3]
	//sql.Tail = matches[5]
	return sql, nil
}
