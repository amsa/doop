#!/bin/bash

SLEEP_TIME=5
if [ $1 ]; then
    SLEEP_TIME=$1
fi

rm -rf ~/.doop
cp demo-orig.db demo.db

echo "> doop init demo sqlite://demo.db"
doop init demo sqlite://demo.db
sleep $SLEEP_TIME
echo

echo "> doop run master@demo 'SELECT * FROM products'"
doop run master@demo 'SELECT * FROM products'
sleep $SLEEP_TIME
echo

echo "> doop branch demo test"
doop branch demo test
sleep $SLEEP_TIME
echo

echo "> doop run test@demo 'INSERT INTO products (id, name, type, price) VALUES (7, "Intel 730 240", "SSD", 240.88)'"
doop run test@demo 'INSERT INTO products (id, name, type, price) VALUES (7, "Intel 730 240", "SSD", 240.88)' &> /dev/null
sleep $SLEEP_TIME
echo

echo "> doop run test@demo 'SELECT * FROM products'"
doop run test@demo 'SELECT * FROM products'
sleep $SLEEP_TIME
echo

echo "> doop run master@demo 'SELECT * FROM products'"
doop run master@demo 'SELECT * FROM products'
sleep $SLEEP_TIME

