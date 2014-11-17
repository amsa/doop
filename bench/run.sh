if [ $1 ]; then
    OPT="-benchtime $1"
fi

go test -bench BenchmarkQueryNoBranch $OPT ./bench
go test -bench BenchmarkQueryBranch $OPT ./bench
rm -rf ./bench/test_db
