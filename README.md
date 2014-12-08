What is Doop?
=============
Doop is a tool that aims to provide branching and versioning for your database, just like what Git does for your code base! 

Doop is:

- A Marvel character 
- In the dictionary: “A little copper cup in which a diamond is held while being cut.”

##`Doop struct` implements following functions
1. `TrackDB(dsn string) (bool, error)` -- start managing the database by given dsn.
2. `ListDBs() ([]string)` -- list the names of databases this Doop instance is managing.
3. `GetDB(db_name string) (*DoopDB, error)` -- given the name, return the corresponding DoopDB object 
4. `UntrackDb(db_name string) (bool, error)` -- untrack the database 

##`DoopDB struct` implements following functions
1. `Init() error`
2. `CreateBranch(branch_name string, base_branch string) (bool, error)`
3. `ListBranches() ([]string, error)` 
4. `Exec(branchName string, sql string, args ...interface{}) (sql.Result, error)`
5. `Query(branchName string, sql string, args ...interface{}) (*sql.Rows, error)`


Command-line interface
=======================
Available commands:
```
- doop init <alias> <DSN> (e.g. doop init testdb sqlite://mytestdb.db)
- doop list  [<alias>] (e.g. doop list [testdb])  -- if <alias> is passed the list of branches will be returned
- doop rm -d <alias>  (e.g. doop rm -d testdb)
- doop rm -b <branch@alias> (e.g. doop rm -b myfork@testdb)
- doop run <branch@alias> <sql> (e.g. doop run myfork@testdb "SELECT * FROM users")
- doop branch <alias> <new> [<from>] (e.g. doop branch testdb myfork test) -- if from is not given the default value is master
- doop merge  <alias> <to> <from> (e.g. doop merge testdb test myfork)
- doop stats [<alias>] (e.g. doop stats [testdb])
- doop export [<branch@alias>] (e.g. doop export testdb@myfork)
```

DSN Format
===========
There is a standard format for DSN to connect to different databases:
```
- sqlite:     sqlite://<path to db file> (e.g. sqlite:///usr/local/myapp/mydb.db)
- mysql:      mysql://username:password@address/dbname?param=value
- postgres:   postgresql://username:password@address/dbname?param=value
```

Required dependencies
==================
- [TableWriter](https://github.com/olekukonko/tablewriter)
- [Testify](https://github.com/stretchr/testify) (only for running tests)

Installation
=============
1. Follow the instruction to [install Go distribution](http://golang.org/doc/install) for your OS
2. Install required dependencies as mentioned in the section above
3. If you have set Go environment variables properly you can run `go install` inside Doop root directory
4. Now if you have `$GOROOT/bin` in your `$PATH` you should be able to run `doop`

Alternatively, if you don't want to install it, you can simply run `go run doop.go` inside Doop root directory

Contributors
============
- Amin Saeidi
- Wenbin Xiao
