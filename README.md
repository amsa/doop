doop-core
=========

##doop-core package provides following APIs:

1. `Doop` -- struct that represents the Doop instance, manages one or multiple DoopDBs.
1. `DoopDB` -- struct that represents Doop-wrapping database 
0. `Result_type`  -- struct that the result of query on database
1. `GetDoop() Doop`  -- return a Doop instance

##`Doop struct` implements following functions
1. `TrackDB(dsn string) (bool, error)` -- start managing the database by given dsn.
2. `ListDBs() ([]string, error)` -- list the names of databases this Doop instance is managing.
3. `GetDB(db_name string) (*DoopDB, error)` -- given the name, return the corresponding DoopDB object 
4. `RemoveDB(db_name string) (bool, error)` -- remove the DoopDB from current Doop instance

##`DoopDB struct` implements following functions
1. `createBranch(branch_name string, base_branch string) (bool, error)`
2. `removeBranch(branch_name string) (bool, error)`
3. `listBranches() ([]string, error)` 
4. `executeSQL(sql_string string, branch_name string, result *Result_type) (bool, error)`


Command-line interface
=======================
Administration commands:
```
- doop init <alias> <DSN> (e.g. doop init testdb sqlite://mytestdb.db)
- doop list  [<alias>] (e.g. doop list [testdb])  -- if <alias> is passed the list of branches will be returned
- doop branch <alias> <new> <from> (e.g. doop branch testdb myfork test)
- doop merge  <alias> <to> <from> (e.g. doop merge testdb test myfork)
- doop rm -d <alias>  (e.g. doop rm -d testdb)
- doop rm -b <branch@alias> (e.g. doop rm -b myfork@testdb)
- doop stats [<alias>] (e.g. doop stats [testdb])
- doop export [<branch@alias>] (e.g. doop export testdb@myfork)
```

Query command:
```
- doopq <branch/db> <query> (e.g. doopq myfork@testdb "SELECT * FROM user)
```

DSN Format
===========
There is a standard format for DSN to connect to different databases:
```
- sqlite:     sqlite://<path to db file> (e.g. sqlite:///usr/local/myapp/mydb.db)
- mysql:      mysql://username:password@address/dbname?param=value
- postgres:   postgresql://username:password@address/dbname?param=value
```
