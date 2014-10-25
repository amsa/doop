package core

//file: common.go
//only this file contains APIs that are exported to doop-core user.

//Package Level
type Doop struct {
	home_dir string
}

type DoopDB struct {
	dbName string
}

type Result struct {
	rawResult string
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

func (doop *Doop) GetDB(dbName string) (*DoopDB, error) {
	return nil, nil
}

func (doop *Doop) UntrackDB(dbName string) (bool, error) {
	return false, nil
}

//Methods of DoopDB
func (db *DoopDB) createBranch(branchName string, baseBranch string) (bool, error) {
	return false, nil
}

func (db *DoopDB) removeBranch(branchName string) (bool, error) {
	return false, nil
}

func (db *DoopDB) listBranches() (bool, error) {
	return false, nil
}

func (db *DoopDB) executeSQL(sqlCommand string, branchName string) (*Result, error) {
	return nil, nil
}
