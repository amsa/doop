if [ $1 ]; then
    OPT="-benchtime $1"
fi

echo Startin MySQL Benchmark...
go test -bench . $OPT ./bench/mysql
echo Startin SQLite Benchmark...
go test -bench . $OPT ./bench/sqlite
rm -rf ./bench/sqlite/test_db*
