@echo off
cd %~dp0src
go get -u
go mod edit -go 1.24.1
go mod tidy