#!/bin/bash
trap "exit" ERR

go get github.com/denisenkom/go-mssqldb
go get github.com/lib/pq
go get github.com/go-sql-driver/mysql
go get golang.org/x/crypto/md4 # required if cross building from windows
