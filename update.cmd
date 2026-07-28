@echo off
cd %~dp0src
go get -u
go mod edit -go 1.26.5
go mod tidy