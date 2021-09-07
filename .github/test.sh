#!/bin/bash

set -e

export LOCAL_TEST=true

module_name=$(cat go.mod | grep module | cut -d ' ' -f 2-2)
echo "module_name is $module_name"

go test -coverprofile=coverage.txt -coverpkg=./... -parallel 1 -p 1 -count=1 -gcflags=-l -v ./...

# go tool cover -func=coverage.txt -o coverage.txt

# go tool cover -html=coverage.txt
