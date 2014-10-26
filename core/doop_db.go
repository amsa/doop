package core

type DoopDb interface {
	//Methods of DoopDB interface
	getInfo() (string, error)
	createBranch(branchName string, baseBranch string) (bool, error)
	removeBranch(branchName string) (bool, error)
	listBranches() ([]string, error)
	executeSQL(sqlCommand string, branchName string) (*Result, error)
}
