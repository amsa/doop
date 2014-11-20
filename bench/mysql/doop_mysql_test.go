package mysql

/*
tests for doop_db.go
*/
import (
	"fmt"
	"os"
	"testing"

	"github.com/amsa/doop/adapter"
	"github.com/amsa/doop/core"
)

var db *core.DoopDb
var dbAdapter1 adapter.Adapter
var dbAdapter2 adapter.Adapter

func init() {
	db1 := "mysql://root:@/doop_test_bench1"
	db2 := "mysql://root:@/doop_test_bench2"

	dbAdapter1 = adapter.GetAdapter(db1)
	dbAdapter2 = adapter.GetAdapter(db2)

	dbAdapter1.DropDb()
	dbAdapter2.DropDb()

	dbAdapter1.CreateDb()
	dbAdapter2.CreateDb()

	// reset the connection
	dbAdapter1 = adapter.GetAdapter(db1)
	dbAdapter2 = adapter.GetAdapter(db2)

	adapter.SetupDb(dbAdapter1)
	adapter.SetupDb(dbAdapter2)

	db = core.MakeDoopDb(&core.DoopDbInfo{db1, "", ""})
	err := db.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.CreateBranch("branch1", "master")
	db.Exec("branch1", "INSERT INTO t1 VALUES(1827, 8718, 'test branch1')")
	db.Exec("master", "INSERT INTO t1 VALUES(1927, 7718, 'test master')")

	dbAdapter2.Exec("INSERT INTO t1 VALUES(2000, 6281, 'test adapter')")
}

func BenchmarkQueryBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db.Query("master", "SELECT * FROM t1")
	}
}

func BenchmarkQueryNoBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dbAdapter2.Query("SELECT * FROM t1")
	}
}
