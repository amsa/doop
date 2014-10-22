doop-core
=========

##doop-core package provides following APIs:

0. `DSN_type`
1. `DB_type`
0. `Result_type`
2. `CreateDB(dsn DSN_type) (bool, error)`
3. `RemoveDB(dsn DSN_type) (bool, error)`
3. `CloneDB(dsn DSN_type)`
4. `ListDBs() ([]DSN_type, error)`
5. `GetDB(dsn DSN_type) (DB_type, error)`

##DB_type interface requires to implement:

1. `createBranch(branch_name string, base_branch string) (bool, error)`
2. `removeBranch(branch_name string) (bool, error)`
3. `listBranches() ([]string, error)`
4. `executeSQL(sql_string string, branch_name string, result *Result_type) (bool, error)`

##Result_type interface
