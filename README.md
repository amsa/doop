What is Doop?
=============
Doop is a tool that aims to provide branching and versioning for your database, just like what Git does for your code base! 

Doop is:

<img src="http://upload.wikimedia.org/wikipedia/en/d/d6/Doop.png" width="12%" />
- A Marvel character 
- In the dictionary: “A little copper cup in which a diamond is held while being cut.”

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
- doop rm -d <alias>  (e.g. doop rm -d testdb)
- doop rm -b <branch@alias> (e.g. doop rm -b myfork@testdb)
- doop branch <alias> <new> <from> (e.g. doop branch testdb myfork test)
- doop merge  <alias> <to> <from> (e.g. doop merge testdb test myfork)
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

Contributors
============
- Amin Saeidi
- Wenbin Xiao
