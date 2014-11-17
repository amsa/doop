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
	dbPath := "test_db"
	test.CleanDb(dbPath)
	test.SetupDb(dbPath)
	db = core.MakeDoopDb(&core.DoopDbInfo{"sqlite://" + dbPath, dbPath, ""})
	dbAdapter = adapter.GetAdapter("sqlite://" + dbPath)
	dbAdapter.Exec("INSERT INTO t1 VALUES(2000, 6281, 'test adapter')")

	db.Init()
	db.CreateBranch("branch1", "master")
	db.Exec("branch1", "INSERT INTO t1 VALUES(1827, 8718, 'test branch1')")
	db.Exec("master", "INSERT INTO t1 VALUES(1927, 7718, 'test master')")
}

func BenchmarkQueryNoBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dbAdapter.Query("SELECT * FROM t1")
	}
}

func BenchmarkQueryBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db.Query("master", "SELECT * FROM t1")
	}
}
