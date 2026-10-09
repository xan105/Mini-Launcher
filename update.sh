#!/bin/sh
cd "$(dirname "$0")/src"
go get -u
go mod edit -go 1.27.2
go mod tidy