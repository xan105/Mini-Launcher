@echo off
cd %~dp0src
go get -u
go mod edit -go 1.26.4
go mod tidy