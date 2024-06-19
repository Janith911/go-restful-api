#! /bin/bash

export GOAPI_ENDPOINT=0.0.0.0:8000
export MYSQL_ENDPOINT=192.168.1.171:3306
export MYSQL_PASS=goapipasswd
export MYSQL_SCHEMA=goapi
export MYSQL_USER=goapiuser
export NR_APP_NAME=Go-Users-API
export NR_LICENSE_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
go run main.go