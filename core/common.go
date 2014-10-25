package core

//file: common.go
//only this file contains APIs that are exported to doop-core user.

//Package Level
type Doop struct {
	home_dir string
}

type DoopDB struct {
	db_name string
}

type Result struct {
	raw_result string
}

func GetDoop() (*Doop, error) {
	return nil, nil
}

//Methods of Doop
func (doop *Doop) TrackDB(dsn string) (bool, error) {
	return false, nil
}

func (doop *Doop) ListDBs() ([]string, error) {
	return nil, nil
}

func (doop *Doop) GetDB(db_name string) (*DoopDB, error) {
	return nil, nil
}

func (doop *Doop) UntrackDB(db_name string) (bool, error) {
	return false, nil
}

//Methods of DoopDB
func (db *DoopDB) createBranch(branch_name string, base_branch string) (bool, error) {
	return false, nil
}

func (db *DoopDB) removeBranch(branch_name string) (bool, error) {
	return false, nil
}

func (db *DoopDB) listBranches() (bool, error) {
	return false, nil
}

func (db *DoopDB) executeSQL(sql_command string, branch_name string) (*Result, error) {
	return nil, nil
}
