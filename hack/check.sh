#!/usr/bin/env bash

set -o errexit
set -o nounset

timeout=10s
if [ $# -eq 0 ]
  then
        for d in $(go list ./...); do
            go test -timeout $timeout -race $d
        done
        exit
fi

pkg=$1
testname=$2
echo "Running test pkg:" $pkg " name: " $testname
go test -timeout $timeout -v -race $pkg --run $testname
