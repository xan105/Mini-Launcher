#!/bin/sh
cd "$(dirname "$0")/src"
go get -u
go mod edit -go 1.24.3
go mod tidy