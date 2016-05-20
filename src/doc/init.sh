#!/bin/sh

echo "*************start init TODO*************"

echo "go get github.com/go-martini/martini ..."
go get github.com/go-martini/martini

echo "go get github.com/go-sql-driver/mysql ..."
go get github.com/go-sql-driver/mysql

echo "go get github.com/coopernurse/gorp ..."
go get github.com/coopernurse/gorp

echo "github.com/martini-contrib/binding ..."
github.com/martini-contrib/binding

echo "github.com/martini-contrib/render ..."
github.com/martini-contrib/render

echo "github.com/martini-contrib/sessions ..."
github.com/martini-contrib/sessions

echo "****************complete****************"
