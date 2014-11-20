if [ $1 ]; then
    OPT="-benchtime $1"
fi

go test -bench . $OPT ./bench/mysql
#rm -rf ./bench/sqlite/test_db*
