if [ $1 ]; then
    OPT="-benchtime $1"
fi

go test -bench . $OPT ./bench
rm -rf ./bench/test_db
