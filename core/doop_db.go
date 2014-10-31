package core

type DoopDbInfo struct {
	DSN  string
	Name string
	Hash string
}

//type DoopDb interface {
////Methods of DoopDB interface
//CreateBranch(branchName string, baseBranch string) (bool, error)
//RemoveBranch(branchName string) (bool, error)
//ListBranches() ([]string, error)
//MergeBranch(to string, from string) (bool, error)
//Query(branchName string, statement string) (*Rows, error)
//Exec(branchName string, statement string) (*Result, error)
//Close() error
//}
