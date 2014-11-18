package bench

/*
tests for doop_db.go
*/
import (
	"testing"

	"github.com/amsa/doop/adapter"
	"github.com/amsa/doop/core"
	"github.com/amsa/doop/test"
)

var db *core.DoopDb
var dbAdapter adapter.Adapter

func init() {
	dbPath1 := "test_db1"
	dbPath2 := "test_db2"

	test.CleanDb(dbPath1)
	test.CleanDb(dbPath2)

	test.SetupDb(dbPath1)
	test.SetupDb(dbPath2)

	db = core.MakeDoopDb(&core.DoopDbInfo{"sqlite://" + dbPath1, dbPath1, ""})
	db.Init()
	db.CreateBranch("branch1", "master")
	db.Exec("branch1", "INSERT INTO t1 VALUES(1827, 8718, 'test branch1')")
	db.Exec("master", "INSERT INTO t1 VALUES(1927, 7718, 'test master')")

	dbAdapter = adapter.GetAdapter("sqlite://" + dbPath2)
	dbAdapter.Exec("INSERT INTO t1 VALUES(2000, 6281, 'test adapter')")
}

func BenchmarkQueryBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db.Query("master", "SELECT * FROM t1")
	}
}

func BenchmarkQueryNoBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dbAdapter.Query("SELECT * FROM t1")
	}
}
