doop-core
=========

##doop-core package provides following APIs:

0. DSN_type
1. DB_type
2. `CreateDB(dsn DSN_type) (bool, error)`
3. `RemoveDB(dsn DSN_type) (bool, error)`
4. `ListDB() ([]DSN_type, error)`
5. `GetDB(dsn DSN_type) (DB_type, error)`

##DB_Type interface requires to implement:

0. Result_type
1. `createBranch(branch_name string) (bool, error)`
2. `removeBranch(branch_name string) (bool, error)`
3. `listBranch() ([]string, error)`
4. `executeSQL(sql_string string, branch_name string, result *Result_type) (bool, error)`
