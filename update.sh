#!/bin/sh
cd "$(dirname "$0")/src"
go get -u
go mod edit -go 1.25.0
go mod tidy