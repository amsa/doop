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

