@echo off
cd %~dp0src
go get -u
go mod edit -go 1.27.2
go mod tidy