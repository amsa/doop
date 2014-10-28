package core

type DoopDb interface {
	//Methods of DoopDB interface
	CreateBranch(branchName string, baseBranch string) (bool, error)
	RemoveBranch(branchName string) (bool, error)
	ListBranches() ([]string, error)
	MergeBranch(to string, from string) (bool, error)
	Query(branchName string, statement string, args ...interface{}) (*Rows, error)
	Exec(branchName string, statement string, args ...interface{}) (*Result, error)
	Close() error
}